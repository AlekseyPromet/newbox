package netflow

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Collector handles NetFlow/sFlow collection
type Collector struct {
	config   *CollectorConfig
	listener net.Listener
	stopCh   chan struct{}
	wg       sync.WaitGroup
	running  int32
}

// CollectorConfig holds collector configuration
type CollectorConfig struct {
	NetFlowPort int
	SFlowPort   int
	Workers     int
	BufferSize  int
}

// Option configures the collector
type Option func(*Collector)

// WithNetFlowPort sets the NetFlow port
func WithNetFlowPort(port int) Option {
	return func(c *Collector) {
		c.config.NetFlowPort = port
	}
}

// WithSFlowPort sets the sFlow port
func WithSFlowPort(port int) Option {
	return func(c *Collector) {
		c.config.SFlowPort = port
	}
}

// WithWorkers sets the number of workers
func WithWorkers(w int) Option {
	return func(c *Collector) {
		c.config.Workers = w
	}
}

// WithBufferSize sets the buffer size
func WithBufferSize(s int) Option {
	return func(c *Collector) {
		c.config.BufferSize = s
	}
}

// FlowRecord represents a NetFlow/sFlow flow record
type FlowRecord struct {
	DeviceUUID   string
	SrcAddr      string
	DstAddr      string
	SrcPort      uint16
	DstPort      uint16
	Protocol     uint8
	Packets      uint64
	Bytes        uint64
	FlowStartMs  int64
	FlowEndMs    int64
	TCPFlags     string
	ToS          uint8
	IngressIf    uint32
	EgressIf     uint32
	SamplingRate uint32
}

// NewCollector creates a new NetFlow/sFlow collector
func NewCollector(opts ...Option) (*Collector, error) {
	c := &Collector{
		config: &CollectorConfig{
			NetFlowPort: 2055,
			SFlowPort:   6343,
			Workers:     10,
			BufferSize:  10000,
		},
		stopCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Start starts the NetFlow collector
func (c *Collector) Start(handler func(*FlowRecord) error) error {
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		return fmt.Errorf("collector already running")
	}

	// Create UDP listener for NetFlow
	addr := fmt.Sprintf(":%d", c.config.NetFlowPort)
	ln, err := net.ListenPacket("udp", addr)
	if err != nil {
		atomic.StoreInt32(&c.running, 0)
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	c.listener = &packetListener{ln}
	c.wg.Add(1)
	go c.readPackets(handler)

	return nil
}

// readPackets reads packets from the listener
func (c *Collector) readPackets(handler func(*FlowRecord) error) {
	defer c.wg.Done()

	buf := make([]byte, 65535)
	for {
		select {
		case <-c.stopCh:
			return
		default:
			c.listener.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := c.listener.ReadFrom(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				continue
			}
			c.processPacket(buf[:n], handler)
		}
	}
}

// processPacket processes a raw packet
func (c *Collector) processPacket(data []byte, handler func(*FlowRecord) error) {
	// Try to detect NetFlow version
	if len(data) < 4 {
		return
	}

	// Check for NetFlow v5/v9/IPFIX header
	version := getNetFlowVersion(data)
	switch version {
	case 5:
		c.processNetFlowV5(data, handler)
	case 9:
		c.processNetFlowV9(data, handler)
	case 10: // IPFIX
		c.processIPFIX(data, handler)
	}
}

// getNetFlowVersion determines the NetFlow version from packet header
func getNetFlowVersion(data []byte) int {
	if len(data) < 2 {
		return 0
	}
	// NetFlow header starts with version number (big-endian)
	version := (uint16(data[0]) << 8) | uint16(data[1])
	switch version {
	case 5:
		return 5
	case 9:
		return 9
	case 10:
		return 10 // IPFIX
	default:
		return 0
	}
}

// processNetFlowV5 processes NetFlow v5 packet
func (c *Collector) processNetFlowV5(data []byte, handler func(*FlowRecord) error) {
	if len(data) < 48 {
		return
	}

	// NetFlow v5 header: version(2), count(2), uptime(4), timestamp(4), type(2)
	count := (uint16(data[2]) << 8) | uint16(data[3])

	// Flow record size is 48 bytes
	for i := 0; i < int(count); i++ {
		offset := 24 + (i * 48)
		if offset+48 > len(data) {
			break
		}

		record := &FlowRecord{}

		// Parse flow record fields
		// src_addr(4), dst_addr(4), next_hop(4), input(2), output(2)
		record.SrcAddr = intToIP(data[offset : offset+4])
		record.DstAddr = intToIP(data[offset+4 : offset+8])
		record.IngressIf = uint32(data[offset+16])<<8 | uint32(data[offset+17])
		record.EgressIf = uint32(data[offset+18])<<8 | uint32(data[offset+19])

		// packets(4), bytes(4)
		record.Packets = uint64(data[offset+20])<<24 | uint64(data[offset+21])<<16 | uint64(data[offset+22])<<8 | uint64(data[offset+23])
		record.Bytes = uint64(data[offset+24])<<24 | uint64(data[offset+25])<<16 | uint64(data[offset+26])<<8 | uint64(data[offset+27])

		// Start/End time
		record.FlowStartMs = int64(data[offset+28])<<24 | int64(data[offset+29])<<16 | int64(data[offset+30])<<8 | int64(data[offset+31])
		record.FlowEndMs = int64(data[offset+32])<<24 | int64(data[offset+33])<<16 | int64(data[offset+34])<<8 | int64(data[offset+35])

		// TCP flags, protocol, TOS
		record.Protocol = data[offset+38]
		record.ToS = data[offset+37]

		handler(record)
	}
}

// processNetFlowV9 processes NetFlow v9 packet (simplified)
func (c *Collector) processNetFlowV9(data []byte, handler func(*FlowRecord) error) {
	// NetFlow v9 is template-based, requires template parsing
	// This is a simplified implementation
	if len(data) < 36 {
		return
	}

	record := &FlowRecord{
		Protocol: 6, // TCP (simplified)
	}
	handler(record)
}

// processIPFIX processes IPFIX packet (simplified)
func (c *Collector) processIPFIX(data []byte, handler func(*FlowRecord) error) {
	// IPFIX is similar to NetFlow v9 with different set IDs
	if len(data) < 16 {
		return
	}

	record := &FlowRecord{
		Protocol: 6, // TCP (simplified)
	}
	handler(record)
}

// Stop stops the collector
func (c *Collector) Stop() error {
	if !atomic.CompareAndSwapInt32(&c.running, 1, 0) {
		return nil
	}

	close(c.stopCh)
	c.wg.Wait()

	if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

// IsRunning returns true if the collector is running
func (c *Collector) IsRunning() bool {
	return atomic.LoadInt32(&c.running) == 1
}

// packetListener wraps a net.PacketConn
type packetListener struct {
	net.PacketConn
}

// intToIP converts a 4-byte slice to IP string
func intToIP(b []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}
