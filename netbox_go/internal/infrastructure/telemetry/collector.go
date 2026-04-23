package telemetry

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"netbox_go/internal/domain/telemetry/entity"
	"netbox_go/internal/domain/telemetry/service"
	"netbox_go/internal/infrastructure/gnmi"
	"netbox_go/internal/infrastructure/influxdb"
	"netbox_go/internal/infrastructure/netflow"
)

// Collector handles telemetry data collection from devices
type Collector struct {
	id             string
	config         *CollectorConfig
	service        *service.TelemetryService
	gnmiClient     *gnmi.Client
	gnmiPoller     *gnmi.Poller
	gnmiSub        *gnmi.Subscriber
	netflowColl    *netflow.Collector
	influxWriter   *influxdb.Writer
	logger         *zap.Logger
	stopCh         chan struct{}
	wg             sync.WaitGroup
	running        atomic.Bool
	activeJobs     int32
	processedCount int32
	errorCount     int32
	lastHeartbeat  time.Time
	mu             sync.RWMutex

	// Coordinator client for registration
	coordinator *Coordinator
}

// CollectorConfig holds collector configuration
type CollectorConfig struct {
	ID                string
	Name              string
	Address           string
	Port              int
	Weight            int
	Zone              string
	Region            string
	HeartbeatInterval time.Duration
	WorkerCount       int
	JobQueueSize      int
	RetryAttempts     int
	RetryDelay        time.Duration
}

// DefaultCollectorConfig returns default collector configuration
func DefaultCollectorConfig() *CollectorConfig {
	return &CollectorConfig{
		HeartbeatInterval: 10 * time.Second,
		WorkerCount:       4,
		JobQueueSize:      100,
		RetryAttempts:     3,
		RetryDelay:        5 * time.Second,
		Weight:            1,
	}
}

// NewCollector creates a new telemetry collector
func NewCollector(
	cfg *CollectorConfig,
	svc *service.TelemetryService,
	gnmiClient *gnmi.Client,
	influxWriter *influxdb.Writer,
	logger *zap.Logger,
) *Collector {
	if cfg == nil {
		cfg = DefaultCollectorConfig()
	}

	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}

	c := &Collector{
		id:           cfg.ID,
		config:       cfg,
		service:      svc,
		gnmiClient:   gnmiClient,
		influxWriter: influxWriter,
		logger:       logger,
		stopCh:       make(chan struct{}),
	}

	return c
}

// SetCoordinator sets the coordinator for this collector
func (c *Collector) SetCoordinator(coordinator *Coordinator) {
	c.coordinator = coordinator
}

// SetNetFlowCollector sets the NetFlow collector
func (c *Collector) SetNetFlowCollector(coll *netflow.Collector) {
	c.netflowColl = coll
}

// SetGNMISubscriber sets the gNMI subscriber
func (c *Collector) SetGNMISubscriber(sub *gnmi.Subscriber) {
	c.gnmiSub = sub
}

// SetGNMIPoller sets the gNMI poller
func (c *Collector) SetGNMIPoller(poller *gnmi.Poller) {
	c.gnmiPoller = poller
}

// Start begins the collector's background work
func (c *Collector) Start(ctx context.Context) error {
	if c.running.Load() {
		return fmt.Errorf("collector already running")
	}

	// Register with coordinator
	if c.coordinator != nil {
		info := &CollectorInfo{
			ID:      c.id,
			Name:    c.config.Name,
			Address: c.config.Address,
			Port:    c.config.Port,
			Weight:  c.config.Weight,
			Zone:    c.config.Zone,
			Region:  c.config.Region,
		}
		if err := c.coordinator.RegisterCollector(ctx, info); err != nil {
			c.logger.Warn("failed to register with coordinator", zap.Error(err))
		}
	}

	c.running.Store(true)
	c.stopCh = make(chan struct{})
	c.lastHeartbeat = time.Now()

	// Start heartbeat loop
	c.wg.Add(1)
	go c.heartbeatLoop(ctx)

	c.logger.Info("collector started",
		zap.String("id", c.id),
		zap.String("name", c.config.Name),
	)

	return nil
}

// Stop stops the collector
func (c *Collector) Stop(ctx context.Context) error {
	if !c.running.Load() {
		return nil
	}

	// Unregister from coordinator
	if c.coordinator != nil {
		if err := c.coordinator.UnregisterCollector(ctx, c.id); err != nil {
			c.logger.Warn("failed to unregister from coordinator", zap.Error(err))
		}
	}

	close(c.stopCh)
	c.wg.Wait()
	c.running.Store(false)

	c.logger.Info("collector stopped", zap.String("id", c.id))
	return nil
}

// heartbeatLoop sends periodic heartbeats to the coordinator
func (c *Collector) heartbeatLoop(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.coordinator != nil {
				if err := c.coordinator.Heartbeat(ctx, c.id); err != nil {
					c.logger.Warn("heartbeat failed", zap.Error(err))
				}
			}
			c.updateMetrics(ctx)
		}
	}
}

// updateMetrics updates collector metrics in the coordinator
func (c *Collector) updateMetrics(ctx context.Context) {
	if c.coordinator == nil {
		return
	}

	stats := c.GetStats()
	c.coordinator.UpdateCollectorMetrics(ctx, c.id, stats.DeviceCount, stats.ActiveJobs)
}

// CollectDevice performs telemetry collection for a single device
func (c *Collector) CollectDevice(ctx context.Context, deviceID uuid.UUID) error {
	if !c.running.Load() {
		return fmt.Errorf("collector not running")
	}

	atomic.AddInt32(&c.activeJobs, 1)
	defer atomic.AddInt32(&c.activeJobs, -1)

	// Get device configuration
	device, err := c.service.GetDevice(ctx, deviceID)
	if err != nil {
		atomic.AddInt32(&c.errorCount, 1)
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Create collection job
	job := &entity.CollectionJob{
		ID:          uuid.New(),
		JobID:       uuid.New().String(),
		CollectorID: c.id,
		DeviceID:    deviceID,
		Status:      "running",
	}

	if err := c.service.CreateJob(ctx, job); err != nil {
		c.logger.Warn("failed to create job record", zap.Error(err))
	}

	// Perform collection based on device type
	var recordsCollected int
	var collectionErr error

	switch device.CollectionType {
	case entity.CollectionTypePoll, entity.CollectionTypeBoth:
		recordsCollected, collectionErr = c.pollDevice(ctx, device)
	case entity.CollectionTypeSubscribe, entity.CollectionTypeBoth:
		recordsCollected, collectionErr = c.subscribeDevice(ctx, device)
	default:
		collectionErr = fmt.Errorf("unknown collection type: %s", device.CollectionType)
	}

	// Update job status
	if collectionErr != nil {
		c.service.MarkJobFailed(ctx, job.ID, collectionErr.Error())
		atomic.AddInt32(&c.errorCount, 1)
		device.CollectionErrorsCount++
		device.LastCollectionStatus = "failed"
	} else {
		c.service.MarkJobCompleted(ctx, job.ID, recordsCollected)
		atomic.AddInt32(&c.processedCount, 1)
		now := time.Now()
		device.LastCollectionAt = &now
		device.LastCollectionStatus = "success"
		device.CollectionErrorsCount = 0
	}

	device.LastCollectionAt = &time.Now()
	c.service.UpdateDevice(ctx, device)

	return collectionErr
}

// pollDevice performs poll-based collection
func (c *Collector) pollDevice(ctx context.Context, device *entity.TelemetryDevice) (int, error) {
	if c.gnmiPoller == nil {
		return 0, fmt.Errorf("gNMI poller not configured")
	}

	c.logger.Debug("polling device",
		zap.String("device_id", device.ID.String()),
		zap.String("gnmi_address", device.GNMIAddress),
	)

	// Get credentials from vault if needed
	creds, err := c.getCredentials(ctx, device.VaultSecretPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Perform gNMI poll
	data, err := c.gnmiPoller.Poll(ctx, device.GNMIAddress, device.GNmiPort, creds)
	if err != nil {
		return 0, fmt.Errorf("gNMI poll failed: %w", err)
	}

	// Write to InfluxDB
	if c.influxWriter != nil {
		if err := c.influxWriter.Write(ctx, device.ID.String(), data); err != nil {
			return 0, fmt.Errorf("failed to write to influxdb: %w", err)
		}
	}

	return len(data), nil
}

// subscribeDevice sets up subscription-based collection
func (c *Collector) subscribeDevice(ctx context.Context, device *entity.TelemetryDevice) (int, error) {
	if c.gnmiSub == nil {
		return 0, fmt.Errorf("gNMI subscriber not configured")
	}

	c.logger.Debug("subscribing to device",
		zap.String("device_id", device.ID.String()),
		zap.String("gnmi_address", device.GNMIAddress),
	)

	// This would typically set up a persistent subscription
	// For now, just mark as subscribed
	return 0, nil
}

// CollectPing performs ICMP ping collection
func (c *Collector) CollectPing(ctx context.Context, target *entity.PingTarget) error {
	// Ping collection would be handled by the ping infrastructure
	c.logger.Debug("ping collection not yet implemented",
		zap.String("target_id", target.ID.String()),
	)
	return nil
}

// CollectDNS performs DNS query collection
func (c *Collector) CollectDNS(ctx context.Context, query *entity.DNSQuery) error {
	// DNS collection would be handled by the DNS infrastructure
	c.logger.Debug("DNS collection not yet implemented",
		zap.String("query_id", query.ID.String()),
	)
	return nil
}

// CollectNetFlow processes NetFlow data
func (c *Collector) CollectNetFlow(ctx context.Context, record *netflow.FlowRecord) error {
	if c.influxWriter == nil {
		return fmt.Errorf("influxdb writer not configured")
	}

	c.logger.Debug("processing netflow record",
		zap.String("device_uuid", record.DeviceUUID),
		zap.String("src_addr", record.SrcAddr),
		zap.String("dst_addr", record.DstAddr),
	)

	// Convert to telemetry data and write
	data := map[string]interface{}{
		"device_uuid": record.DeviceUUID,
		"src_addr":    record.SrcAddr,
		"dst_addr":    record.DstAddr,
		"src_port":    record.SrcPort,
		"dst_port":    record.DstPort,
		"protocol":    record.Protocol,
		"packets":     record.Packets,
		"bytes":       record.Bytes,
		"flow_start":  record.FlowStartMs,
		"flow_end":    record.FlowEndMs,
	}

	return c.influxWriter.Write(ctx, record.DeviceUUID, []map[string]interface{}{data})
}

// getCredentials retrieves credentials from vault
func (c *Collector) getCredentials(ctx context.Context, secretPath string) (*gnmi.Credentials, error) {
	if secretPath == "" {
		return &gnmi.Credentials{}, nil
	}

	// This would use the vault client to retrieve credentials
	// For now, return empty credentials (no auth)
	return &gnmi.Credentials{}, nil
}

// GetStats returns collector statistics
func (c *Collector) GetStats() CollectorStats {
	return CollectorStats{
		ID:             c.id,
		Name:           c.config.Name,
		ActiveJobs:     atomic.LoadInt32(&c.activeJobs),
		ProcessedCount: atomic.LoadInt32(&c.processedCount),
		ErrorCount:     atomic.LoadInt32(&c.errorCount),
		Running:        c.running.Load(),
		LastHeartbeat:  c.lastHeartbeat,
	}
}

// CollectorStats holds collector statistics
type CollectorStats struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	ActiveJobs     int32     `json:"active_jobs"`
	ProcessedCount int32     `json:"processed_count"`
	ErrorCount     int32     `json:"error_count"`
	Running        bool      `json:"running"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
	DeviceCount    int32     `json:"device_count"`
}
