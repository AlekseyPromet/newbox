package influxdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// Writer handles writing telemetry data to InfluxDB
type Writer struct {
	client        influxdb.Client
	org           string
	bucket        string
	batchSize     int
	flushInterval time.Duration
	buffer        []*write.Point
	bufferMu      sync.Mutex
	closeCh       chan struct{}
}

// Config holds InfluxDB writer configuration
type Config struct {
	URL           string
	Token         string
	Org           string
	Bucket        string
	BatchSize     int
	FlushInterval time.Duration
	Timeout       time.Duration
}

// Option configures the InfluxDB writer
type Option func(*Writer)

// WithBatchSize sets the batch size
func WithBatchSize(size int) Option {
	return func(w *Writer) {
		w.batchSize = size
	}
}

// WithFlushInterval sets the flush interval
func WithFlushInterval(interval time.Duration) Option {
	return func(w *Writer) {
		w.flushInterval = interval
	}
}

// NewWriter creates a new InfluxDB writer
func NewWriter(cfg *Config) (*Writer, error) {
	client := influxdb.NewClientWithOptions(
		cfg.URL,
		cfg.Token,
		influxdb.DefaultOptions().
			SetBatchSize(cfg.BatchSize).
			SetFlushInterval(cfg.FlushInterval).
			SetTimeout(cfg.Timeout),
	)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to influxdb: %w", err)
	}

	w := &Writer{
		client:        client,
		org:           cfg.Org,
		bucket:        cfg.Bucket,
		batchSize:     cfg.BatchSize,
		flushInterval: cfg.FlushInterval,
		buffer:        make([]*write.Point, 0, cfg.BatchSize),
		closeCh:       make(chan struct{}),
	}

	// Start background flusher
	go w.backgroundFlusher()

	return w, nil
}

// WritePoint writes a single point to InfluxDB
func (w *Writer) WritePoint(ctx context.Context, measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	pt := influxdb.NewPoint(measurement, tags, fields, timestamp)

	w.bufferMu.Lock()
	w.buffer = append(w.buffer, pt)
	shouldFlush := len(w.buffer) >= w.batchSize
	w.bufferMu.Unlock()

	if shouldFlush {
		return w.flush(ctx)
	}

	return nil
}

// WritePoints writes multiple points to InfluxDB
func (w *Writer) WritePoints(ctx context.Context, points []*write.Point) error {
	w.bufferMu.Lock()
	w.buffer = append(w.buffer, points...)
	shouldFlush := len(w.buffer) >= w.batchSize
	w.bufferMu.Unlock()

	if shouldFlush {
		return w.flush(ctx)
	}

	return nil
}

// flush writes all buffered points to InfluxDB
func (w *Writer) flush(ctx context.Context) error {
	w.bufferMu.Lock()
	if len(w.buffer) == 0 {
		w.bufferMu.Unlock()
		return nil
	}

	points := w.buffer
	w.buffer = make([]*write.Point, 0, w.batchSize)
	w.bufferMu.Unlock()

	writeAPI := w.client.WriteAPI(w.org, w.bucket)
	for _, pt := range points {
		writeAPI.WritePoint(ctx, pt)
	}

	// Flush and check for errors
	writeAPI.Flush()

	// Get errors from the channel
	errCh := writeAPI.Errors()
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("influxdb write error: %w", err)
		}
	default:
	}

	return nil
}

// backgroundFlusher periodically flushes the buffer
func (w *Writer) backgroundFlusher() {
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			w.flush(ctx)
			cancel()
		case <-w.closeCh:
			return
		}
	}
}

// Query executes a query against InfluxDB
func (w *Writer) Query(ctx context.Context, query string) (*QueryResult, error) {
	queryAPI := w.client.QueryAPI(w.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("influxdb query failed: %w", err)
	}

	qr := &QueryResult{
		Result: result,
	}

	return qr, nil
}

// QueryResult holds query results
type QueryResult struct {
	Result interface{}
}

// Close closes the InfluxDB writer
func (w *Writer) Close() error {
	close(w.closeCh)

	// Final flush
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	w.flush(ctx)

	w.client.Close()
	return nil
}
