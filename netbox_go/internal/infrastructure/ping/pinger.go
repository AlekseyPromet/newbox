package ping

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// Pinger handles ICMP ping operations
type Pinger struct {
	config *PingerConfig
	mu     sync.RWMutex
}

// PingerConfig holds pinger configuration
type PingerConfig struct {
	Timeout    time.Duration
	Count      int
	PacketSize int
	Workers    int
}

// Option configures the pinger
type Option func(*Pinger)

// WithTimeout sets the timeout
func WithTimeout(t time.Duration) Option {
	return func(p *Pinger) {
		p.config.Timeout = t
	}
}

// WithCount sets the packet count
func WithCount(c int) Option {
	return func(p *Pinger) {
		p.config.Count = c
	}
}

// WithPacketSize sets the packet size
func WithPacketSize(s int) Option {
	return func(p *Pinger) {
		p.config.PacketSize = s
	}
}

// WithWorkers sets the number of workers
func WithWorkers(w int) Option {
	return func(p *Pinger) {
		p.config.Workers = w
	}
}

// PingResult holds the result of a ping operation
type PingResult struct {
	TargetAddress   string
	RTTMs           float64
	RTTMinMs        float64
	RTTMaxMs        float64
	RTTAvgMs        float64
	PacketLossPct   float64
	PacketsSent     int
	PacketsReceived int
	TTL             int
	Err             error
}

// NewPinger creates a new pinger
func NewPinger(opts ...Option) (*Pinger, error) {
	p := &Pinger{
		config: &PingerConfig{
			Timeout:    5 * time.Second,
			Count:      5,
			PacketSize: 64,
			Workers:    50,
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

// Ping performs a ping to the target address
func (p *Pinger) Ping(ctx context.Context, target string) (*PingResult, error) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create pinger: %w", err)
	}

	pinger.SetTimeout(p.config.Timeout)
	pinger.Count = p.config.Count
	pinger.Size = p.config.PacketSize

	// Use privileged mode if available
	pinger.Privileged = true

	// Run ping
	err = pinger.Run()
	if err != nil {
		return &PingResult{
			TargetAddress: target,
			Err:           fmt.Errorf("ping failed: %w", err),
		}, nil
	}

	stats := pinger.Statistics()

	return &PingResult{
		TargetAddress:   target,
		RTTMs:           float64(stats.Rtt.Milliseconds()),
		RTTMinMs:        float64(stats.MinRtt.Milliseconds()),
		RTTMaxMs:        float64(stats.MaxRtt.Milliseconds()),
		RTTAvgMs:        float64(stats.AvgRtt.Milliseconds()),
		PacketsSent:     stats.PacketsSent,
		PacketsReceived: stats.PacketsRecv,
		PacketLossPct:   float64(stats.PacketsSent-stats.PacketsRecv) / float64(stats.PacketsSent) * 100,
		TTL:             stats.TTL,
	}, nil
}

// PingAsync performs a ping asynchronously
func (p *Pinger) PingAsync(ctx context.Context, target string, results chan<- *PingResult) {
	go func() {
		result, err := p.Ping(ctx, target)
		if err != nil {
			results <- &PingResult{
				TargetAddress: target,
				Err:           err,
			}
			return
		}
		results <- result
	}()
}

// PingBatch performs pings to multiple targets concurrently
func (p *Pinger) PingBatch(ctx context.Context, targets []string) ([]*PingResult, error) {
	results := make(chan *PingResult, len(targets))
	done := make(chan struct{})

	// Worker pool
	var wg sync.WaitGroup
	workerCount := p.config.Workers
	if workerCount > len(targets) {
		workerCount = len(targets)
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range targets {
				select {
				case <-done:
					return
				default:
					p.PingAsync(ctx, target, results)
				}
			}
		}()
	}

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
		close(done)
	}()

	// Collect results
	var pingResults []*PingResult
	for result := range results {
		pingResults = append(pingResults, result)
	}

	return pingResults, nil
}
