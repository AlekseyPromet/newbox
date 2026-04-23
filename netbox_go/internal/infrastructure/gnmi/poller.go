package gnmi

import (
	"context"
	"fmt"
	"sync"
	"time"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// Poller implements poll-based gNMI collection
type Poller struct {
	client *Client
	config *PollerConfig
	paths  []*PathBuilder
	mu     sync.RWMutex
	stopCh chan struct{}
}

// PollerConfig holds poller configuration
type PollerConfig struct {
	Interval  time.Duration
	Timeout   time.Duration
	RateLimit int // requests per second
}

// PathBuilder builds gNMI paths for different telemetry types
type PathBuilder struct {
	TelemetryType string
	Path          *gpb.Path
	Fields        []string
}

// NewPoller creates a new gNMI poller
func NewPoller(client *Client, cfg *PollerConfig) *Poller {
	return &Poller{
		client: client,
		config: cfg,
		paths:  make([]*PathBuilder, 0),
		stopCh: make(chan struct{}),
	}
}

// AddPath adds a path to poll
func (p *Poller) AddPath(pb *PathBuilder) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paths = append(p.paths, pb)
}

// RemovePath removes a path from polling
func (p *Poller) RemovePath(telemetryType string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, pb := range p.paths {
		if pb.TelemetryType == telemetryType {
			p.paths = append(p.paths[:i], p.paths[i+1:]...)
			return
		}
	}
}

// Start starts the poller
func (p *Poller) Start(ctx context.Context, handler func(*PathBuilder, *gpb.GetResponse) error) {
	ticker := time.NewTicker(p.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.pollOnce(ctx, handler)
		case <-p.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

// pollOnce performs a single poll iteration
func (p *Poller) pollOnce(ctx context.Context, handler func(*PathBuilder, *gpb.GetResponse) error) {
	p.mu.RLock()
	paths := make([]*PathBuilder, len(p.paths))
	copy(paths, p.paths)
	p.mu.RUnlock()

	for _, pb := range paths {
		if err := p.pollPath(ctx, pb, handler); err != nil {
			// Log error but continue with other paths
			fmt.Printf("poll error for %s: %v\n", pb.TelemetryType, err)
		}
	}
}

// pollPath polls a single path
func (p *Poller) pollPath(ctx context.Context, pb *PathBuilder, handler func(*PathBuilder, *gpb.GetResponse) error) error {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()

	// Rate limiting
	time.Sleep(time.Second / time.Duration(p.config.RateLimit))

	resp, err := p.client.Get(ctx, nil, []*gpb.Path{pb.Path}, gpb.Encoding_JSON_IETF)
	if err != nil {
		return fmt.Errorf("gNMI get failed: %w", err)
	}

	return handler(pb, resp)
}

// Stop stops the poller
func (p *Poller) Stop() {
	close(p.stopCh)
}

// InterfacePathBuilder creates path for interface statistics
func InterfacePathBuilder() *PathBuilder {
	return &PathBuilder{
		TelemetryType: "interface",
		Path: &gpb.Path{
			Origin: "openconfig",
			Path:   "interfaces",
			Elem: []*gpb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
			},
		},
		Fields: []string{
			"state/counters/in-octets",
			"state/counters/out-octets",
			"state/counters/in-pkts",
			"state/counters/out-pkts",
			"state/counters/in-errors",
			"state/counters/out-errors",
			"state/oper-status",
			"state/admin-status",
			"state/speed",
		},
	}
}

// BGPPathBuilder creates path for BGP metrics
func BGPPathBuilder() *PathBuilder {
	return &PathBuilder{
		TelemetryType: "bgp",
		Path: &gpb.Path{
			Origin: "openconfig",
			Path:   "network-instances/network-instance[name=default]/protocols/protocol[name=BGP]",
			Elem: []*gpb.PathElem{
				{Name: "network-instances"},
				{Name: "network-instance", Key: map[string]string{"name": "default"}},
				{Name: "protocols"},
				{Name: "protocol", Key: map[string]string{"name": "BGP"}},
			},
		},
		Fields: []string{
			"state/peer[neighbor-address=*]/state",
		},
	}
}

// SystemPathBuilder creates path for system metrics
func SystemPathBuilder() *PathBuilder {
	return &PathBuilder{
		TelemetryType: "system",
		Path: &gpb.Path{
			Origin: "openconfig",
			Path:   "system",
			Elem: []*gpb.PathElem{
				{Name: "system"},
			},
		},
		Fields: []string{
			"state/cpu",
			"state/memory/physical",
			"state/uptime",
			"state/temperature",
		},
	}
}

// OSPFPathBuilder creates path for OSPF metrics
func OSPFPathBuilder() *PathBuilder {
	return &PathBuilder{
		TelemetryType: "ospf",
		Path: &gpb.Path{
			Origin: "openconfig",
			Path:   "network-instances/network-instance[name=default]/protocols/protocol[name=OSPF]",
			Elem: []*gpb.PathElem{
				{Name: "network-instances"},
				{Name: "network-instance", Key: map[string]string{"name": "default"}},
				{Name: "protocols"},
				{Name: "protocol", Key: map[string]string{"name": "OSPF"}},
			},
		},
		Fields: []string{
			"state/area",
			"state/neighbor",
		},
	}
}
