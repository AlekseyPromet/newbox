package gnmi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// Client wraps gNMI client functionality
type Client struct {
	conn   *grpc.ClientConn
	client gpb.GNMIClient
	config *Config
	mu     sync.RWMutex
}

// Config holds gNMI client configuration
type Config struct {
	Address    string
	Username   string
	Password   string
	Timeout    time.Duration
	TLSEnabled bool
	TLSCert    string
	TLSKey     string
	TLSCACert  string
}

// Option configures the gNMI client
type Option func(*Client)

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.config.Timeout = timeout
	}
}

// WithTLS configures TLS
func WithTLS(enabled bool, cert, key, ca string) Option {
	return func(c *Client) {
		c.config.TLSEnabled = enabled
		c.config.TLSCert = cert
		c.config.TLSKey = key
		c.config.TLSCACert = ca
	}
}

// WithCredentials sets username and password
func WithCredentials(username, password string) Option {
	return func(c *Client) {
		c.config.Username = username
		c.config.Password = password
	}
}

// NewClient creates a new gNMI client
func NewClient(opts ...Option) (*Client, error) {
	cfg := &Config{
		Timeout: 30 * time.Second,
	}

	c := &Client{
		config: cfg,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Connect establishes connection to the device
func (c *Client) Connect(ctx context.Context, address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var opts []grpc.DialOption

	if c.config.TLSEnabled {
		tlsConfig, err := c.getTLSConfig()
		if err != nil {
			return fmt.Errorf("failed to load TLS config: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	c.conn = conn
	c.client = gpb.NewGNMIClient(conn)
	c.config.Address = address

	return nil
}

// getTLSConfig returns TLS configuration
func (c *Client) getTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{}

	if c.config.TLSCert != "" && c.config.TLSKey != "" {
		cert, err := tls.LoadX509KeyPair(c.config.TLSCert, c.config.TLSKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if c.config.TLSCACert != "" {
		caCert, err := os.ReadFile(c.config.TLSCACert)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA cert")
		}
		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = false

	return tlsConfig, nil
}

// Capabilities queries device capabilities
func (c *Client) Capabilities(ctx context.Context) (*gpb.CapabilityResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil, fmt.Errorf("client not connected")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	return c.client.Capabilities(ctx, &gpb.CapabilityRequest{})
}

// Get retrieves data from the device using a gNMI Get request
func (c *Client) Get(ctx context.Context, prefix *gpb.Path, path []*gpb.Path, encoding gpb.Encoding) (*gpb.GetResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil, fmt.Errorf("client not connected")
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &gpb.GetRequest{
		Prefix:   prefix,
		Path:     path,
		Encoding: encoding,
	}

	return c.client.Get(ctx, req)
}

// Subscribe establishes a subscription to the device
func (c *Client) Subscribe(ctx context.Context, stream gpb.GNMI_SubscribeClient) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return fmt.Errorf("client not connected")
	}

	return c.client.Subscribe(ctx, stream)
}

// SubscribeStream creates a subscription stream
func (c *Client) SubscribeStream(ctx context.Context, sub *gpb.SubscriptionList) (<-chan *gpb.SubscribeResponse, error) {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("client not connected")
	}

	stream, err := client.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription stream: %w", err)
	}

	// Send subscription request
	if err := stream.Send(&gpb.SubscribeRequest{
		Request: &gpb.SubscribeRequest_Subscribe{
			Subscribe: sub,
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to send subscription: %w", err)
	}

	// Receive response to confirm subscription
	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to receive subscription confirmation: %w", err)
	}

	// Create response channel
	respCh := make(chan *gpb.SubscribeResponse)

	go func() {
		defer close(respCh)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			select {
			case respCh <- resp:
			case <-ctx.Done():
				return
			}
		}
	}()

	return respCh, nil
}

// Close closes the gNMI client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil
}
