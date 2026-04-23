package telemetry

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"netbox_go/internal/domain/telemetry/entity"
	"netbox_go/internal/domain/telemetry/service"
)

// CollectorInfo holds information about a registered collector
type CollectorInfo struct {
	ID            string
	Name          string
	Address       string
	Port          int
	Weight        int // capacity weight for load balancing
	Zone          string
	Region        string
	RegisteredAt  time.Time
	LastHeartbeat time.Time
	Status        CollectorStatus
	DeviceCount   int32
	ActiveJobs    int32
}

// CollectorStatus represents the health status of a collector
type CollectorStatus string

const (
	CollectorStatusHealthy   CollectorStatus = "healthy"
	CollectorStatusDegraded  CollectorStatus = "degraded"
	CollectorStatusUnhealthy CollectorStatus = "unhealthy"
	CollectorStatusOffline   CollectorStatus = "offline"
)

// Coordinator manages distributed telemetry collectors
type Coordinator struct {
	config     *CoordinatorConfig
	service    *service.TelemetryService
	logger     *zap.Logger
	collectors map[string]*CollectorInfo
	devices    map[uuid.UUID]*entity.TelemetryDevice
	mu         sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
	running    atomic.Bool

	// Callbacks
	onDeviceAssigned   func(deviceID uuid.UUID, collectorID string)
	onDeviceUnassigned func(deviceID uuid.UUID)
	onRebalanceNeeded  func()
}

// CoordinatorConfig holds coordinator configuration
type CoordinatorConfig struct {
	HeartbeatTimeout       time.Duration
	RebalanceThreshold     float64 // trigger rebalance when load variance exceeds this
	MinCollectorWeight     int
	MaxCollectorWeight     int
	HealthCheckInterval    time.Duration
	MaxDevicesPerCollector int
	ZoneAwareness          bool
}

// DefaultCoordinatorConfig returns default coordinator configuration
func DefaultCoordinatorConfig() *CoordinatorConfig {
	return &CoordinatorConfig{
		HeartbeatTimeout:       30 * time.Second,
		RebalanceThreshold:     0.5, // 50% variance triggers rebalance
		MinCollectorWeight:     1,
		MaxCollectorWeight:     10,
		HealthCheckInterval:    10 * time.Second,
		MaxDevicesPerCollector: 100,
		ZoneAwareness:          true,
	}
}

// NewCoordinator creates a new telemetry coordinator
func NewCoordinator(
	cfg *CoordinatorConfig,
	svc *service.TelemetryService,
	logger *zap.Logger,
) *Coordinator {
	if cfg == nil {
		cfg = DefaultCoordinatorConfig()
	}

	c := &Coordinator{
		config:     cfg,
		service:    svc,
		logger:     logger,
		collectors: make(map[string]*CollectorInfo),
		devices:    make(map[uuid.UUID]*entity.TelemetryDevice),
		stopCh:     make(chan struct{}),
	}

	return c
}

// RegisterCollector registers a new collector
func (c *Coordinator) RegisterCollector(ctx context.Context, info *CollectorInfo) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if info.ID == "" {
		info.ID = uuid.New().String()
	}
	info.RegisteredAt = time.Now()
	info.LastHeartbeat = time.Now()
	info.Status = CollectorStatusHealthy

	c.collectors[info.ID] = info
	c.logger.Info("collector registered",
		zap.String("collector_id", info.ID),
		zap.String("name", info.Name),
		zap.Int("weight", info.Weight),
	)

	return nil
}

// UnregisterCollector unregisters a collector
func (c *Coordinator) UnregisterCollector(ctx context.Context, collectorID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	collector, ok := c.collectors[collectorID]
	if !ok {
		return fmt.Errorf("collector %s not found", collectorID)
	}

	// Unassign all devices from this collector
	c.unassignDevicesFromCollector(ctx, collectorID)

	delete(c.collectors, collectorID)
	c.logger.Info("collector unregistered", zap.String("collector_id", collectorID))

	return nil
}

// Heartbeat updates collector heartbeat
func (c *Coordinator) Heartbeat(ctx context.Context, collectorID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	collector, ok := c.collectors[collectorID]
	if !ok {
		return fmt.Errorf("collector %s not found", collectorID)
	}

	collector.LastHeartbeat = time.Now()
	if collector.Status == CollectorStatusOffline {
		collector.Status = CollectorStatusHealthy
	}

	return nil
}

// GetCollectors returns all registered collectors
func (c *Coordinator) GetCollectors(ctx context.Context) []*CollectorInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*CollectorInfo, 0, len(c.collectors))
	for _, col := range c.collectors {
		result = append(result, col)
	}
	return result
}

// GetCollector returns a collector by ID
func (c *Coordinator) GetCollector(ctx context.Context, collectorID string) (*CollectorInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	collector, ok := c.collectors[collectorID]
	if !ok {
		return nil, fmt.Errorf("collector %s not found", collectorID)
	}
	return collector, nil
}

// AssignDevice assigns a device to an appropriate collector
func (c *Coordinator) AssignDevice(ctx context.Context, deviceID uuid.UUID) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	device, err := c.service.GetDevice(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("failed to get device: %w", err)
	}

	// Find best collector for this device
	collectorID, err := c.findBestCollector(ctx, device)
	if err != nil {
		return "", err
	}

	// Update device
	device.AssignedCollectorID = collectorID
	if err := c.service.UpdateDevice(ctx, device); err != nil {
		return "", fmt.Errorf("failed to update device: %w", err)
	}

	// Update collector device count
	if collector, ok := c.collectors[collectorID]; ok {
		atomic.AddInt32(&collector.DeviceCount, 1)
	}

	c.devices[deviceID] = device
	c.logger.Info("device assigned to collector",
		zap.String("device_id", deviceID.String()),
		zap.String("collector_id", collectorID),
	)

	if c.onDeviceAssigned != nil {
		c.onDeviceAssigned(deviceID, collectorID)
	}

	return collectorID, nil
}

// UnassignDevice unassigns a device from its collector
func (c *Coordinator) UnassignDevice(ctx context.Context, deviceID uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	device, ok := c.devices[deviceID]
	if !ok {
		// Try to get from service
		var err error
		device, err = c.service.GetDevice(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("device %s not found", deviceID)
		}
	}

	if device.AssignedCollectorID == "" {
		return nil // already unassigned
	}

	collectorID := device.AssignedCollectorID

	// Update device
	device.AssignedCollectorID = ""
	if err := c.service.UpdateDevice(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	// Update collector device count
	if collector, ok := c.collectors[collectorID]; ok {
		atomic.AddInt32(&collector.DeviceCount, -1)
	}

	c.logger.Info("device unassigned from collector",
		zap.String("device_id", deviceID.String()),
		zap.String("collector_id", collectorID),
	)

	if c.onDeviceUnassigned != nil {
		c.onDeviceUnassigned(deviceID)
	}

	return nil
}

// Rebalance performs load rebalancing across collectors
func (c *Coordinator) Rebalance(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.collectors) == 0 {
		return fmt.Errorf("no collectors available")
	}

	// Calculate target distribution
	totalDevices := int32(0)
	for _, col := range c.collectors {
		totalDevices += col.DeviceCount
	}

	if totalDevices == 0 {
		return nil // nothing to rebalance
	}

	// Calculate per-collector targets based on weight
	type rebalanceAction struct {
		collectorID string
		deviceCount int32
	}

	actions := make([]rebalanceAction, 0, len(c.collectors))
	totalWeight := 0
	for _, col := range c.collectors {
		totalWeight += col.Weight
	}

	if totalWeight == 0 {
		// Equal distribution
		perCollector := int(math.Ceil(float64(totalDevices) / float64(len(c.collectors))))
		for id := range c.collectors {
			actions = append(actions, rebalanceAction{collectorID: id, deviceCount: int32(perCollector)})
		}
	} else {
		for id, col := range c.collectors {
			targetCount := int32(math.Round(float64(totalDevices) * float64(col.Weight) / float64(totalWeight)))
			actions = append(actions, rebalanceAction{collectorID: id, deviceCount: targetCount})
		}
	}

	// Check if rebalance is needed (variance exceeds threshold)
	currentVariance := c.calculateVariance()
	if currentVariance < c.config.RebalanceThreshold {
		c.logger.Debug("rebalance not needed, variance below threshold",
			zap.Float64("current_variance", currentVariance),
			zap.Float64("threshold", c.config.RebalanceThreshold),
		)
		return nil
	}

	c.logger.Info("starting rebalance",
		zap.Int32("total_devices", totalDevices),
		zap.Float64("variance", currentVariance),
	)

	// Perform reassignments
	for _, action := range actions {
		collector := c.collectors[action.collectorID]
		diff := action.deviceCount - collector.DeviceCount

		if diff > 0 {
			// Need to assign more devices
			c.rebalanceAssignDevices(ctx, action.collectorID, int(diff))
		} else if diff < 0 {
			// Need to unassign devices
			c.rebalanceUnassignDevices(ctx, action.collectorID, int(-diff))
		}
	}

	if c.onRebalanceNeeded != nil {
		c.onRebalanceNeeded()
	}

	return nil
}

// Start begins the coordinator's background work
func (c *Coordinator) Start(ctx context.Context) error {
	if c.running.Load() {
		return fmt.Errorf("coordinator already running")
	}

	c.running.Store(true)
	c.stopCh = make(chan struct{})

	c.wg.Add(1)
	go c.healthCheckLoop(ctx)

	c.logger.Info("coordinator started")
	return nil
}

// Stop stops the coordinator
func (c *Coordinator) Stop(ctx context.Context) error {
	if !c.running.Load() {
		return nil
	}

	close(c.stopCh)
	c.wg.Wait()
	c.running.Store(false)

	c.logger.Info("coordinator stopped")
	return nil
}

// healthCheckLoop periodically checks collector health
func (c *Coordinator) healthCheckLoop(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.checkCollectorHealth(ctx)
		}
	}
}

// checkCollectorHealth checks and updates collector status
func (c *Coordinator) checkCollectorHealth(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for id, collector := range c.collectors {
		timeSinceHeartbeat := now.Sub(collector.LastHeartbeat)

		switch {
		case timeSinceHeartbeat > c.config.HeartbeatTimeout*3:
			if collector.Status != CollectorStatusOffline {
				collector.Status = CollectorStatusOffline
				c.logger.Warn("collector marked offline",
					zap.String("collector_id", id),
					zap.Duration("time_since_heartbeat", timeSinceHeartbeat),
				)
				// Unassign devices for offline collector
				c.unassignDevicesFromCollector(ctx, id)
			}
		case timeSinceHeartbeat > c.config.HeartbeatTimeout:
			if collector.Status == CollectorStatusHealthy {
				collector.Status = CollectorStatusDegraded
				c.logger.Warn("collector marked degraded",
					zap.String("collector_id", id),
					zap.Duration("time_since_heartbeat", timeSinceHeartbeat),
				)
			}
		}
	}
}

// findBestCollector finds the best collector for a device using zone awareness
func (c *Coordinator) findBestCollector(ctx context.Context, device *entity.TelemetryDevice) (string, error) {
	var candidates []*CollectorInfo

	if c.config.ZoneAwareness {
		// Find collectors in the same zone
		for _, col := range c.collectors {
			if col.Status == CollectorStatusHealthy || col.Status == CollectorStatusDegraded {
				if col.Zone == "" || col.Zone == "default" {
					candidates = append(candidates, col)
				}
			}
		}
	}

	if len(candidates) == 0 {
		// Fall back to all healthy collectors
		for _, col := range c.collectors {
			if col.Status == CollectorStatusHealthy || col.Status == CollectorStatusDegraded {
				candidates = append(candidates, col)
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no healthy collectors available")
	}

	// Select collector with lowest device count (least loaded)
	var best *CollectorInfo
	minCount := int32(math.MaxInt32)
	for _, col := range candidates {
		// Factor in weight: effective_load = device_count / weight
		effectiveLoad := float64(col.DeviceCount) / float64(col.Weight)
		if col.DeviceCount < minCount {
			minCount = col.DeviceCount
			best = col
		}
	}

	return best.ID, nil
}

// calculateVariance calculates device distribution variance
func (c *Coordinator) calculateVariance() float64 {
	if len(c.collectors) == 0 {
		return 0
	}

	var total int32
	for _, col := range c.collectors {
		total += col.DeviceCount
	}

	if total == 0 {
		return 0
	}

	avg := float64(total) / float64(len(c.collectors))
	var variance float64

	for _, col := range c.collectors {
		diff := float64(col.DeviceCount) - avg
		variance += diff * diff
	}

	return variance / float64(len(c.collectors))
}

// unassignDevicesFromCollector unassigns all devices from a collector
func (c *Coordinator) unassignDevicesFromCollector(ctx context.Context, collectorID string) {
	for _, col := range c.collectors {
		if col.ID == collectorID {
			count := col.DeviceCount
			col.DeviceCount = 0
			c.logger.Info("unassigned devices from offline collector",
				zap.String("collector_id", collectorID),
				zap.Int32("device_count", count),
			)
		}
	}
}

// rebalanceAssignDevices assigns devices to a collector
func (c *Coordinator) rebalanceAssignDevices(ctx context.Context, collectorID string, count int) {
	collector := c.collectors[collectorID]
	currentCount := collector.DeviceCount

	// Find unassigned devices
	for _, device := range c.devices {
		if device.AssignedCollectorID == "" && count > 0 {
			device.AssignedCollectorID = collectorID
			if err := c.service.UpdateDevice(ctx, device); err == nil {
				collector.DeviceCount++
				count--
			}
		}
	}

	c.logger.Debug("rebalance assigned devices",
		zap.String("collector_id", collectorID),
		zap.Int32("previous_count", currentCount),
		zap.Int32("new_count", collector.DeviceCount),
	)
}

// rebalanceUnassignDevices unassigns devices from a collector
func (c *Coordinator) rebalanceUnassignDevices(ctx context.Context, collectorID string, count int) {
	collector := c.collectors[collectorID]
	currentCount := collector.DeviceCount

	for _, device := range c.devices {
		if device.AssignedCollectorID == collectorID && count > 0 {
			device.AssignedCollectorID = ""
			if err := c.service.UpdateDevice(ctx, device); err == nil {
				collector.DeviceCount--
				count--
			}
		}
	}

	c.logger.Debug("rebalance unassigned devices",
		zap.String("collector_id", collectorID),
		zap.Int32("previous_count", currentCount),
		zap.Int32("new_count", collector.DeviceCount),
	)
}

// UpdateCollectorMetrics updates collector metrics (device count, active jobs)
func (c *Coordinator) UpdateCollectorMetrics(ctx context.Context, collectorID string, deviceCount, activeJobs int32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	collector, ok := c.collectors[collectorID]
	if !ok {
		return fmt.Errorf("collector %s not found", collectorID)
	}

	collector.DeviceCount = deviceCount
	collector.ActiveJobs = activeJobs
	return nil
}

// GetCollectorStats returns statistics for all collectors
func (c *Coordinator) GetCollectorStats(ctx context.Context) map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := make(map[string]interface{})
	collectors := make([]map[string]interface{}, 0, len(c.collectors))

	totalDevices := int32(0)
	totalActiveJobs := int32(0)
	healthyCount := 0

	for _, col := range c.collectors {
		totalDevices += col.DeviceCount
		totalActiveJobs += col.ActiveJobs
		if col.Status == CollectorStatusHealthy {
			healthyCount++
		}

		collectors = append(collectors, map[string]interface{}{
			"id":           col.ID,
			"name":         col.Name,
			"status":       col.Status,
			"device_count": col.DeviceCount,
			"active_jobs":  col.ActiveJobs,
			"weight":       col.Weight,
			"zone":         col.Zone,
		})
	}

	stats["collectors"] = collectors
	stats["total_collectors"] = len(c.collectors)
	stats["healthy_collectors"] = healthyCount
	stats["total_devices"] = totalDevices
	stats["total_active_jobs"] = totalActiveJobs

	return stats
}
