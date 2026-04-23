package telemetry

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TelemetryMetrics holds all Prometheus metrics for telemetry
type TelemetryMetrics struct {
	// Device metrics
	DevicesTotal          prometheus.Gauge
	DevicesEnabled        prometheus.Gauge
	DevicesByCollector    *prometheus.GaugeVec
	CollectionErrorsTotal prometheus.Counter

	// Collection metrics
	CollectionJobsTotal      *prometheus.CounterVec
	CollectionJobsInProgress prometheus.Gauge
	CollectionDuration       *prometheus.HistogramVec
	RecordsCollectedTotal    *prometheus.CounterVec

	// Collector metrics
	CollectorsTotal       prometheus.Gauge
	ActiveCollectors      prometheus.Gauge
	CollectorDeviceCount  *prometheus.GaugeVec
	CollectorActiveJobs   *prometheus.GaugeVec
	CollectorHeartbeatAge *prometheus.GaugeVec

	// Network metrics
	GNMIRequestsTotal   *prometheus.CounterVec
	GNMIRequestDuration prometheus.Histogram
	GNMIErrorsTotal     *prometheus.CounterVec

	// Circuit breaker metrics
	CircuitBreakerState *prometheus.GaugeVec
	CircuitBreakerCalls *prometheus.CounterVec

	// InfluxDB metrics
	InfluxDBWritesTotal   prometheus.Counter
	InfluxDBWriteErrors   prometheus.Counter
	InfluxDBWriteDuration prometheus.Histogram

	// NetFlow metrics
	NetFlowRecordsTotal prometheus.Counter
	NetFlowErrors       prometheus.Counter

	mu sync.RWMutex
}

var (
	globalMetrics *TelemetryMetrics
	once          sync.Once
)

// NewTelemetryMetrics creates and registers telemetry metrics
func NewTelemetryMetrics() *TelemetryMetrics {
	once.Do(func() {
		globalMetrics = &TelemetryMetrics{
			DevicesTotal: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "telemetry_devices_total",
				Help: "Total number of configured telemetry devices",
			}),
			DevicesEnabled: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "telemetry_devices_enabled",
				Help: "Number of enabled telemetry devices",
			}),
			DevicesByCollector: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "telemetry_devices_by_collector",
				Help: "Number of devices assigned to each collector",
			}, []string{"collector_id"}),
			CollectionErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
				Name: "telemetry_collection_errors_total",
				Help: "Total number of collection errors",
			}),

			CollectionJobsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "telemetry_collection_jobs_total",
				Help: "Total number of collection jobs by status",
			}, []string{"status"}),
			CollectionJobsInProgress: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "telemetry_collection_jobs_in_progress",
				Help: "Number of collection jobs currently in progress",
			}),
			CollectionDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name:    "telemetry_collection_duration_seconds",
				Help:    "Duration of collection jobs in seconds",
				Buckets: prometheus.DefBuckets,
			}, []string{"device_id", "collection_type"}),
			RecordsCollectedTotal: promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "telemetry_records_collected_total",
				Help: "Total number of telemetry records collected",
			}, []string{"device_id", "telemetry_type"}),

			CollectorsTotal: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "telemetry_collectors_total",
				Help: "Total number of registered collectors",
			}),
			ActiveCollectors: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "telemetry_active_collectors",
				Help: "Number of active (healthy) collectors",
			}),
			CollectorDeviceCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "telemetry_collector_device_count",
				Help: "Number of devices assigned to each collector",
			}, []string{"collector_id", "zone"}),
			CollectorActiveJobs: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "telemetry_collector_active_jobs",
				Help: "Number of active jobs on each collector",
			}, []string{"collector_id"}),
			CollectorHeartbeatAge: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "telemetry_collector_heartbeat_age_seconds",
				Help: "Age of last heartbeat from each collector in seconds",
			}, []string{"collector_id"}),

			GNMIRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "telemetry_gnmi_requests_total",
				Help: "Total number of gNMI requests",
			}, []string{"device_address", "status"}),
			GNMIRequestDuration: promauto.NewHistogram(prometheus.HistogramOpts{
				Name:    "telemetry_gnmi_request_duration_seconds",
				Help:    "Duration of gNMI requests in seconds",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			}),
			GNMIErrorsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "telemetry_gnmi_errors_total",
				Help: "Total number of gNMI errors by type",
			}, []string{"error_type"}),

			CircuitBreakerState: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "telemetry_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
			}, []string{"name"}),
			CircuitBreakerCalls: promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "telemetry_circuit_breaker_calls_total",
				Help: "Total circuit breaker calls by result",
			}, []string{"name", "result"}),

			InfluxDBWritesTotal: promauto.NewCounter(prometheus.CounterOpts{
				Name: "telemetry_influxdb_writes_total",
				Help: "Total number of InfluxDB writes",
			}),
			InfluxDBWriteErrors: promauto.NewCounter(prometheus.CounterOpts{
				Name: "telemetry_influxdb_write_errors_total",
				Help: "Total number of InfluxDB write errors",
			}),
			InfluxDBWriteDuration: promauto.NewHistogram(prometheus.HistogramOpts{
				Name:    "telemetry_influxdb_write_duration_seconds",
				Help:    "Duration of InfluxDB writes in seconds",
				Buckets: prometheus.DefBuckets,
			}),

			NetFlowRecordsTotal: promauto.NewCounter(prometheus.CounterOpts{
				Name: "telemetry_netflow_records_total",
				Help: "Total number of NetFlow records processed",
			}),
			NetFlowErrors: promauto.NewCounter(prometheus.CounterOpts{
				Name: "telemetry_netflow_errors_total",
				Help: "Total number of NetFlow processing errors",
			}),
		}
	})

	return globalMetrics
}

// GetMetrics returns the global metrics instance
func GetMetrics() *TelemetryMetrics {
	if globalMetrics == nil {
		return NewTelemetryMetrics()
	}
	return globalMetrics
}

// RecordDeviceMetrics updates device-related metrics
func (m *TelemetryMetrics) RecordDeviceMetrics(total, enabled int) {
	m.DevicesTotal.Set(float64(total))
	m.DevicesEnabled.Set(float64(enabled))
}

// RecordDeviceByCollector records device count per collector
func (m *TelemetryMetrics) RecordDeviceByCollector(collectorID string, count int) {
	m.DevicesByCollector.WithLabelValues(collectorID).Set(float64(count))
}

// RecordCollectionError increments the collection error counter
func (m *TelemetryMetrics) RecordCollectionError() {
	m.CollectionErrorsTotal.Inc()
}

// RecordJobStarted increments in-progress jobs
func (m *TelemetryMetrics) RecordJobStarted() {
	m.CollectionJobsInProgress.Inc()
}

// RecordJobCompleted records a completed job
func (m *TelemetryMetrics) RecordJobCompleted(status string, duration time.Duration, deviceID, collectionType string) {
	m.CollectionJobsTotal.WithLabelValues(status).Inc()
	m.CollectionJobsInProgress.Dec()
	m.CollectionDuration.WithLabelValues(deviceID, collectionType).Observe(duration.Seconds())
}

// RecordRecordsCollected records the number of records collected
func (m *TelemetryMetrics) RecordRecordsCollected(deviceID, telemetryType string, count int) {
	m.RecordsCollectedTotal.WithLabelValues(deviceID, telemetryType).Add(float64(count))
}

// RecordCollectorMetrics updates collector metrics
func (m *TelemetryMetrics) RecordCollectorMetrics(collectors []*CollectorInfo) {
	m.CollectorsTotal.Set(float64(len(collectors)))

	active := 0
	for _, c := range collectors {
		if c.Status == CollectorStatusHealthy {
			active++
		}
		m.CollectorDeviceCount.WithLabelValues(c.ID, c.Zone).Set(float64(c.DeviceCount))
		m.CollectorActiveJobs.WithLabelValues(c.ID).Set(float64(c.ActiveJobs))
		m.CollectorHeartbeatAge.WithLabelValues(c.ID).Set(time.Since(c.LastHeartbeat).Seconds())
	}
	m.ActiveCollectors.Set(float64(active))
}

// RecordGNMIRequest records a gNMI request
func (m *TelemetryMetrics) RecordGNMIRequest(deviceAddress, status string, duration time.Duration) {
	m.GNMIRequestsTotal.WithLabelValues(deviceAddress, status).Inc()
	m.GNMIRequestDuration.Observe(duration.Seconds())
}

// RecordGNMIError records a gNMI error
func (m *TelemetryMetrics) RecordGNMIError(errorType string) {
	m.GNMIErrorsTotal.WithLabelValues(errorType).Inc()
}

// RecordCircuitBreakerState records circuit breaker state
func (m *TelemetryMetrics) RecordCircuitBreakerState(name string, state CircuitState) {
	var stateValue float64
	switch state {
	case CircuitStateClosed:
		stateValue = 0
	case CircuitStateOpen:
		stateValue = 1
	case CircuitStateHalfOpen:
		stateValue = 2
	}
	m.CircuitBreakerState.WithLabelValues(name).Set(stateValue)
}

// RecordCircuitBreakerCall records a circuit breaker call
func (m *TelemetryMetrics) RecordCircuitBreakerCall(name, result string) {
	m.CircuitBreakerCalls.WithLabelValues(name, result).Inc()
}

// RecordInfluxDBWrite records an InfluxDB write
func (m *TelemetryMetrics) RecordInfluxDBWrite(duration time.Duration, err error) {
	m.InfluxDBWritesTotal.Inc()
	m.InfluxDBWriteDuration.Observe(duration.Seconds())
	if err != nil {
		m.InfluxDBWriteErrors.Inc()
	}
}

// RecordNetFlowRecord records a NetFlow record
func (m *TelemetryMetrics) RecordNetFlowRecord(err error) {
	m.NetFlowRecordsTotal.Inc()
	if err != nil {
		m.NetFlowErrors.Inc()
	}
}

// MetricsCollector is a background collector that updates metrics periodically
type MetricsCollector struct {
	coordinator *Coordinator
	metrics     *TelemetryMetrics
	interval    time.Duration
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(coordinator *Coordinator, interval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		coordinator: coordinator,
		metrics:     GetMetrics(),
		interval:    interval,
		stopCh:      make(chan struct{}),
	}
}

// Start begins collecting metrics
func (m *MetricsCollector) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.collectLoop(ctx)
}

// Stop stops the metrics collector
func (m *MetricsCollector) Stop() {
	close(m.stopCh)
	m.wg.Wait()
}

// collectLoop periodically collects and updates metrics
func (m *MetricsCollector) collectLoop(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.collect()
		}
	}
}

// collect gathers metrics from various sources
func (m *MetricsCollector) collect() {
	if m.coordinator == nil {
		return
	}

	// Update collector metrics
	ctx := context.Background()
	collectors := m.coordinator.GetCollectors(ctx)

	activeCollectors := 0
	for _, c := range collectors {
		if c.Status == CollectorStatusHealthy || c.Status == CollectorStatusDegraded {
			activeCollectors++
		}
		m.metrics.CollectorDeviceCount.WithLabelValues(c.ID, c.Zone).Set(float64(c.DeviceCount))
		m.metrics.CollectorActiveJobs.WithLabelValues(c.ID).Set(float64(c.ActiveJobs))
		m.metrics.CollectorHeartbeatAge.WithLabelValues(c.ID).Set(time.Since(c.LastHeartbeat).Seconds())
	}

	m.metrics.CollectorsTotal.Set(float64(len(collectors)))
	m.metrics.ActiveCollectors.Set(float64(activeCollectors))
}

// AtomicInt32 is a wrapper for atomic operations on int32
type AtomicInt32 struct {
	value int32
}

// NewAtomicInt32 creates a new atomic int32
func NewAtomicInt32(val int32) *AtomicInt32 {
	return &AtomicInt32{value: val}
}

// Add adds a value and returns the new value
func (a *AtomicInt32) Add(delta int32) int32 {
	return atomic.AddInt32(&a.value, delta)
}

// Get returns the current value
func (a *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&a.value)
}

// Set sets the value
func (a *AtomicInt32) Set(val int32) {
	atomic.StoreInt32(&a.value, val)
}
