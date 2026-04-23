package vault

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
)

// Client wraps HashiCorp Vault client with caching
type Client struct {
	client   *api.Client
	cache    map[string]*cachedSecret
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
	timeout  time.Duration
}

// Config holds Vault client configuration
type Config struct {
	Address  string
	RoleID   string
	SecretID string
	Timeout  time.Duration
	CacheTTL time.Duration
}

// Option configures the Vault client
type Option func(*Client)

// WithAddress sets the Vault address
func WithAddress(addr string) Option {
	return func(c *Client) {
		if c.client != nil {
			c.client.SetAddress(addr)
		}
	}
}

// WithAppRole sets the AppRole authentication
func WithAppRole(roleID, secretID string) Option {
	return func(c *Client) {
		if c.client != nil {
			c.client.SetToken("")
			// AppRole login will be handled in authentication
		}
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithCacheTTL sets the cache TTL
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *Client) {
		c.cacheTTL = ttl
	}
}

// cachedSecret holds a cached secret with expiration
type cachedSecret struct {
	Data      map[string]interface{}
	ExpiresAt time.Time
}

// NewClient creates a new Vault client
func NewClient(opts ...Option) (*Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: "http://localhost:8200",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	c := &Client{
		client:   client,
		cache:    make(map[string]*cachedSecret),
		cacheTTL: 5 * time.Minute,
		timeout:  10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Authenticate performs AppRole authentication
func (c *Client) Authenticate(ctx context.Context, roleID, secretID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Write secret ID to the auth method
	secretData := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// Mount point for AppRole - adjust if using different mount
	resp, err := c.client.Logical().Write("auth/approle/login", secretData)
	if err != nil {
		return fmt.Errorf("vault approle login failed: %w", err)
	}

	if resp == nil || resp.Auth == nil {
		return fmt.Errorf("vault approle login: empty response")
	}

	c.client.SetToken(resp.Auth.ClientToken)
	return nil
}

// GetSecret retrieves a secret from Vault with caching
func (c *Client) GetSecret(ctx context.Context, path string) (map[string]interface{}, error) {
	// Check cache first
	c.cacheMu.RLock()
	if cached, ok := c.cache[path]; ok && time.Now().Before(cached.ExpiresAt) {
		c.cacheMu.RUnlock()
		return cached.Data, nil
	}
	c.cacheMu.RUnlock()

	// Fetch from Vault
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	secret, err := c.client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault read failed for %s: %w", path, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault secret not found: %s", path)
	}

	// Extract the "data" wrapper if present (KV v2)
	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		secret.Data = data
	}

	// Cache the result
	c.cacheMu.Lock()
	c.cache[path] = &cachedSecret{
		Data:      secret.Data,
		ExpiresAt: time.Now().Add(c.cacheTTL),
	}
	c.cacheMu.Unlock()

	return secret.Data, nil
}

// GetDeviceCredentials retrieves credentials for a specific device
func (c *Client) GetDeviceCredentials(ctx context.Context, deviceID string) (username, password string, err error) {
	path := fmt.Sprintf("secret/data/netbox/gnmi/devices/%s", deviceID)

	data, err := c.GetSecret(ctx, path)
	if err != nil {
		return "", "", err
	}

	if u, ok := data["username"].(string); ok {
		username = u
	}
	if p, ok := data["password"].(string); ok {
		password = p
	}

	return username, password, nil
}

// InvalidateCache invalidates the cache for a specific path
func (c *Client) InvalidateCache(path string) {
	c.cacheMu.Lock()
	delete(c.cache, path)
	c.cacheMu.Unlock()
}

// InvalidateAllCache clears the entire cache
func (c *Client) InvalidateAllCache() {
	c.cacheMu.Lock()
	c.cache = make(map[string]*cachedSecret)
	c.cacheMu.Unlock()
}

// Close closes the Vault client
func (c *Client) Close() error {
	c.client = nil
	c.InvalidateAllCache()
	return nil
}
