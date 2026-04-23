package gnmi

import (
	"context"
	"fmt"
	"sync"
	"time"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// SubscriptionMode defines the type of subscription
type SubscriptionMode string

const (
	// ModeOnce requests a single snapshot
	ModeOnce SubscriptionMode = "once"
	// ModeStream provides ongoing updates
	ModeStream SubscriptionMode = "stream"
	// ModePoll provides updates on demand
	ModePoll SubscriptionMode = "poll"
)

// Subscriber implements subscription-based gNMI collection
type Subscriber struct {
	client *Client
	config *SubscriberConfig
	mu     sync.RWMutex
	stopCh chan struct{}
	subs   map[string]*gpb.SubscriptionList
}

// SubscriberConfig holds subscriber configuration
type SubscriberConfig struct {
	Mode              SubscriptionMode
	SampleInterval    time.Duration // for STREAM mode
	Heartbeat         time.Duration // for STREAM mode
	UpdatesOnly       bool
	SuppressRedundant bool
}

// NewSubscriber creates a new gNMI subscriber
func NewSubscriber(client *Client, cfg *SubscriberConfig) *Subscriber {
	return &Subscriber{
		client: client,
		config: cfg,
		stopCh: make(chan struct{}),
		subs:   make(map[string]*gpb.SubscriptionList),
	}
}

// AddSubscription adds a subscription
func (s *Subscriber) AddSubscription(name string, sub *gpb.SubscriptionList) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[name] = sub
}

// RemoveSubscription removes a subscription
func (s *Subscriber) RemoveSubscription(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subs, name)
}

// Subscribe starts a subscription and returns a channel of updates
func (s *Subscriber) Subscribe(ctx context.Context, name string) (<-chan *gpb.SubscribeResponse, error) {
	s.mu.RLock()
	sub, ok := s.subs[name]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("subscription %s not found", name)
	}

	// Set subscription mode
	switch s.config.Mode {
	case ModeOnce:
		sub.Mode = gpb.SubscriptionList_ONCE
	case ModeStream:
		sub.Mode = gpb.SubscriptionList_STREAM
	case ModePoll:
		sub.Mode = gpb.SubscriptionList_POLL
	}

	// Configure stream settings
	if s.config.Mode == ModeStream {
		for _, subscription := range sub.Subscription {
			if subscription.Mode == gpb.Subscription_SAMPLE {
				subscription.SampleInterval = uint64(s.config.SampleInterval.Nanoseconds())
			}
		}
	}

	stream, err := s.client.SubscribeStream(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to start subscription: %w", err)
	}

	return stream, nil
}

// SubscribeAll starts all subscriptions
func (s *Subscriber) SubscribeAll(ctx context.Context, handler func(string, *gpb.SubscribeResponse) error) error {
	s.mu.RLock()
	subs := make(map[string]*gpb.SubscriptionList, len(s.subs))
	for k, v := range s.subs {
		subs[k] = v
	}
	s.mu.RUnlock()

	errCh := make(chan error, len(subs))
	respCh := make(chan struct{})

	for name, sub := range subs {
		go func(n string, subscription *gpb.SubscriptionList) {
			if err := s.subscribeAndHandle(ctx, n, subscription, handler); err != nil {
				errCh <- fmt.Errorf("subscription %s failed: %w", n, err)
			}
		}(name, sub)
	}

	// Wait for completion or context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// subscribeAndHandle handles a single subscription
func (s *Subscriber) subscribeAndHandle(ctx context.Context, name string, sub *gpb.SubscriptionList, handler func(string, *gpb.SubscribeResponse) error) error {
	streamCh, err := s.Subscribe(ctx, name)
	if err != nil {
		return err
	}

	for {
		select {
		case resp, ok := <-streamCh:
			if !ok {
				return fmt.Errorf("stream closed for %s", name)
			}
			if err := handler(name, resp); err != nil {
				return fmt.Errorf("handler error for %s: %w", name, err)
			}
		case <-s.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Stop stops all subscriptions
func (s *Subscriber) Stop() {
	close(s.stopCh)
}

// NewInterfaceSubscription creates a subscription for interface stats
func NewInterfaceSubscription() *gpb.SubscriptionList {
	return &gpb.SubscriptionList{
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Path:   "interfaces",
		},
		Subscription: []*gpb.Subscription{
			{
				Path: &gpb.Path{
					Path: "interface[name=*]/state",
					Elem: []*gpb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "*"}},
						{Name: "state"},
					},
				},
				Mode: gpb.Subscription_SAMPLE,
			},
		},
		Encoding: gpb.Encoding_JSON_IETF,
	}
}

// NewBGPSubscription creates a subscription for BGP metrics
func NewBGPSubscription() *gpb.SubscriptionList {
	return &gpb.SubscriptionList{
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Path:   "network-instances/network-instance[name=default]/protocols/protocol[name=BGP]",
		},
		Subscription: []*gpb.Subscription{
			{
				Path: &gpb.Path{
					Path: "bgp/neighbor[neighbor-address=*]/state",
					Elem: []*gpb.PathElem{
						{Name: "network-instances"},
						{Name: "network-instance", Key: map[string]string{"name": "default"}},
						{Name: "protocols"},
						{Name: "protocol", Key: map[string]string{"name": "BGP"}},
						{Name: "bgp"},
						{Name: "neighbor", Key: map[string]string{"neighbor-address": "*"}},
						{Name: "state"},
					},
				},
				Mode: gpb.Subscription_SAMPLE,
			},
		},
		Encoding: gpb.Encoding_JSON_IETF,
	}
}

// NewSystemSubscription creates a subscription for system metrics
func NewSystemSubscription() *gpb.SubscriptionList {
	return &gpb.SubscriptionList{
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Path:   "system",
		},
		Subscription: []*gpb.Subscription{
			{
				Path: &gpb.Path{
					Path: "state",
					Elem: []*gpb.PathElem{
						{Name: "system"},
						{Name: "state"},
					},
				},
				Mode: gpb.Subscription_SAMPLE,
			},
		},
		Encoding: gpb.Encoding_JSON_IETF,
	}
}
