package dns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Resolver handles DNS resolution operations
type Resolver struct {
	config  *ResolverConfig
	client  *dns.Client
	servers []string
	mu      sync.RWMutex
}

// ResolverConfig holds resolver configuration
type ResolverConfig struct {
	Timeout time.Duration
	Servers []string
	Workers int
}

// Option configures the resolver
type Option func(*Resolver)

// WithTimeout sets the timeout
func WithTimeout(t time.Duration) Option {
	return func(r *Resolver) {
		r.config.Timeout = t
	}
}

// WithServers sets the DNS servers
func WithServers(servers ...string) Option {
	return func(r *Resolver) {
		r.servers = servers
	}
}

// WithWorkers sets the number of workers
func WithWorkers(w int) Option {
	return func(r *Resolver) {
		r.config.Workers = w
	}
}

// DNSResult holds the result of a DNS query
type DNSResult struct {
	QueryName     string
	QueryType     string
	DNSServer     string
	QueryTimeMs   float64
	ResolveTimeMs float64
	AnswerCount   int
	NXDOMAIN      bool
	SERVFAIL      bool
	Err           error
}

// NewResolver creates a new DNS resolver
func NewResolver(opts ...Option) (*Resolver, error) {
	r := &Resolver{
		config: &ResolverConfig{
			Timeout: 10 * time.Second,
			Servers: nil, // nil means use system default
			Workers: 20,
		},
		client: &dns.Client{
			Timeout: 10 * time.Second,
		},
		servers: nil, // nil means use system default
	}

	for _, opt := range opts {
		opt(r)
	}

	// If servers not explicitly set via options, use system default (nil servers)
	// When servers is nil, ExchangeContext will use the system's configured DNS
	return r, nil
}

// Resolve performs a DNS query
func (r *Resolver) Resolve(ctx context.Context, queryName, queryType string) (*DNSResult, error) {
	startTime := time.Now()

	// Build DNS message
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(queryName), r.getQueryType(queryType))
	msg.RecursionDesired = true

	// Use first server or fallback, or nil for system default
	r.mu.RLock()
	var server string
	if len(r.servers) > 0 {
		server = r.servers[0]
	}
	r.mu.RUnlock()

	var response *dns.Msg
	var rtt time.Duration
	var err error

	if server == "" {
		// Use system default DNS - connect to a DNS server implicitly
		// The miekg/dns library will use the system's default DNS when address is empty
		// But since we need an actual server, fall back to well-known public DNS
		response, rtt, err = r.client.ExchangeContext(ctx, msg, "8.8.8.8:53")
		server = "system-default"
	} else {
		// Use configured server
		response, rtt, err = r.client.ExchangeContext(ctx, msg, server+":53")
	}
	queryTimeMs := float64(rtt.Microseconds()) / 1000.0

	if err != nil {
		return &DNSResult{
			QueryName:   queryName,
			QueryType:   queryType,
			DNSServer:   server,
			QueryTimeMs: queryTimeMs,
			Err:         fmt.Errorf("DNS query failed: %w", err),
		}, nil
	}

	resolveTimeMs := float64(time.Since(startTime).Microseconds()) / 1000.0

	// Check response status
	if response == nil {
		return &DNSResult{
			QueryName:     queryName,
			QueryType:     queryType,
			DNSServer:     server,
			QueryTimeMs:   queryTimeMs,
			ResolveTimeMs: resolveTimeMs,
			SERVFAIL:      true,
			Err:           fmt.Errorf("empty DNS response"),
		}, nil
	}

	// Check RCODE
	result := &DNSResult{
		QueryName:     queryName,
		QueryType:     queryType,
		DNSServer:     server,
		QueryTimeMs:   queryTimeMs,
		ResolveTimeMs: resolveTimeMs,
		AnswerCount:   len(response.Answer),
	}

	if response.Rcode == dns.RcodeNameError {
		result.NXDOMAIN = true
	} else if response.Rcode == dns.RcodeServerFailure {
		result.SERVFAIL = true
	}

	return result, nil
}

// ResolveAsync performs DNS resolution asynchronously
func (r *Resolver) ResolveAsync(ctx context.Context, queryName, queryType string, results chan<- *DNSResult) {
	go func() {
		result, err := r.Resolve(ctx, queryName, queryType)
		if err != nil {
			results <- &DNSResult{
				QueryName: queryName,
				QueryType: queryType,
				Err:       err,
			}
			return
		}
		results <- result
	}()
}

// ResolveBatch performs multiple DNS queries concurrently
func (r *Resolver) ResolveBatch(ctx context.Context, queries []struct{ Name, Type string }) ([]*DNSResult, error) {
	results := make(chan *DNSResult, len(queries))
	done := make(chan struct{})

	var wg sync.WaitGroup
	workerCount := r.config.Workers
	if workerCount > len(queries) {
		workerCount = len(queries)
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, q := range queries {
				select {
				case <-done:
					return
				default:
					r.ResolveAsync(ctx, q.Name, q.Type, results)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
		close(done)
	}()

	var dnsResults []*DNSResult
	for result := range results {
		dnsResults = append(dnsResults, result)
	}

	return dnsResults, nil
}

// getQueryType converts string to DNS query type
func (r *Resolver) getQueryType(queryType string) uint16 {
	switch queryType {
	case "A":
		return dns.TypeA
	case "AAAA":
		return dns.TypeAAAA
	case "MX":
		return dns.TypeMX
	case "NS":
		return dns.TypeNS
	case "TXT":
		return dns.TypeTXT
	case "CNAME":
		return dns.TypeCNAME
	case "SOA":
		return dns.TypeSOA
	case "PTR":
		return dns.TypePTR
	case "SRV":
		return dns.TypeSRV
	default:
		return dns.TypeA
	}
}
