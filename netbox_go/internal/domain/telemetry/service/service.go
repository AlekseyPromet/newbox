package service

import (
	"context"
	"time"

	"netbox_go/internal/domain/telemetry/entity"
	"netbox_go/internal/domain/telemetry/repository"

	"github.com/google/uuid"
)

// TelemetryService handles telemetry collection operations
type TelemetryService struct {
	telemetryRepo repository.TelemetryRepository
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(repo repository.TelemetryRepository) *TelemetryService {
	return &TelemetryService{
		telemetryRepo: repo,
	}
}

// Device management

// CreateDevice creates a new telemetry device
func (s *TelemetryService) CreateDevice(ctx context.Context, device *entity.TelemetryDevice) error {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	return s.telemetryRepo.CreateDevice(ctx, device)
}

// GetDevice retrieves a telemetry device by ID
func (s *TelemetryService) GetDevice(ctx context.Context, id uuid.UUID) (*entity.TelemetryDevice, error) {
	return s.telemetryRepo.GetDevice(ctx, id)
}

// ListDevices retrieves all telemetry devices
func (s *TelemetryService) ListDevices(ctx context.Context) ([]*entity.TelemetryDevice, error) {
	return s.telemetryRepo.ListDevices(ctx)
}

// UpdateDevice updates a telemetry device
func (s *TelemetryService) UpdateDevice(ctx context.Context, device *entity.TelemetryDevice) error {
	device.UpdatedAt = time.Now()
	return s.telemetryRepo.UpdateDevice(ctx, device)
}

// DeleteDevice deletes a telemetry device
func (s *TelemetryService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	return s.telemetryRepo.DeleteDevice(ctx, id)
}

// Collection management

// CreateCollection creates a new telemetry collection
func (s *TelemetryService) CreateCollection(ctx context.Context, collection *entity.TelemetryCollection) error {
	if collection.ID == uuid.Nil {
		collection.ID = uuid.New()
	}
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()
	return s.telemetryRepo.CreateCollection(ctx, collection)
}

// GetCollection retrieves a telemetry collection by ID
func (s *TelemetryService) GetCollection(ctx context.Context, id uuid.UUID) (*entity.TelemetryCollection, error) {
	return s.telemetryRepo.GetCollection(ctx, id)
}

// ListCollections retrieves all telemetry collections
func (s *TelemetryService) ListCollections(ctx context.Context) ([]*entity.TelemetryCollection, error) {
	return s.telemetryRepo.ListCollections(ctx)
}

// UpdateCollection updates a telemetry collection
func (s *TelemetryService) UpdateCollection(ctx context.Context, collection *entity.TelemetryCollection) error {
	collection.UpdatedAt = time.Now()
	return s.telemetryRepo.UpdateCollection(ctx, collection)
}

// DeleteCollection deletes a telemetry collection
func (s *TelemetryService) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	return s.telemetryRepo.DeleteCollection(ctx, id)
}

// Job management

// CreateJob creates a new collection job
func (s *TelemetryService) CreateJob(ctx context.Context, job *entity.CollectionJob) error {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}
	job.CreatedAt = time.Now()
	return s.telemetryRepo.CreateJob(ctx, job)
}

// GetJob retrieves a collection job by ID
func (s *TelemetryService) GetJob(ctx context.Context, id uuid.UUID) (*entity.CollectionJob, error) {
	return s.telemetryRepo.GetJob(ctx, id)
}

// ListJobs retrieves recent collection jobs
func (s *TelemetryService) ListJobs(ctx context.Context, limit int) ([]*entity.CollectionJob, error) {
	return s.telemetryRepo.ListJobs(ctx, limit)
}

// UpdateJob updates a collection job
func (s *TelemetryService) UpdateJob(ctx context.Context, job *entity.CollectionJob) error {
	return s.telemetryRepo.UpdateJob(ctx, job)
}

// MarkJobStarted marks a job as started
func (s *TelemetryService) MarkJobStarted(ctx context.Context, id uuid.UUID) error {
	job, err := s.telemetryRepo.GetJob(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	job.StartedAt = &now
	job.Status = "running"
	return s.telemetryRepo.UpdateJob(ctx, job)
}

// MarkJobCompleted marks a job as completed
func (s *TelemetryService) MarkJobCompleted(ctx context.Context, id uuid.UUID, recordsCollected int) error {
	job, err := s.telemetryRepo.GetJob(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	job.CompletedAt = &now
	job.Status = "completed"
	job.RecordsCollected = recordsCollected
	return s.telemetryRepo.UpdateJob(ctx, job)
}

// MarkJobFailed marks a job as failed
func (s *TelemetryService) MarkJobFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	job, err := s.telemetryRepo.GetJob(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	job.CompletedAt = &now
	job.Status = "failed"
	job.ErrorMessage = errMsg
	return s.telemetryRepo.UpdateJob(ctx, job)
}

// Ping target management

// CreatePingTarget creates a new ping target
func (s *TelemetryService) CreatePingTarget(ctx context.Context, target *entity.PingTarget) error {
	if target.ID == uuid.Nil {
		target.ID = uuid.New()
	}
	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()
	return s.telemetryRepo.CreatePingTarget(ctx, target)
}

// GetPingTarget retrieves a ping target by ID
func (s *TelemetryService) GetPingTarget(ctx context.Context, id uuid.UUID) (*entity.PingTarget, error) {
	return s.telemetryRepo.GetPingTarget(ctx, id)
}

// ListPingTargets retrieves all ping targets
func (s *TelemetryService) ListPingTargets(ctx context.Context) ([]*entity.PingTarget, error) {
	return s.telemetryRepo.ListPingTargets(ctx)
}

// UpdatePingTarget updates a ping target
func (s *TelemetryService) UpdatePingTarget(ctx context.Context, target *entity.PingTarget) error {
	target.UpdatedAt = time.Now()
	return s.telemetryRepo.UpdatePingTarget(ctx, target)
}

// DeletePingTarget deletes a ping target
func (s *TelemetryService) DeletePingTarget(ctx context.Context, id uuid.UUID) error {
	return s.telemetryRepo.DeletePingTarget(ctx, id)
}

// DNS query management

// CreateDNSQuery creates a new DNS query
func (s *TelemetryService) CreateDNSQuery(ctx context.Context, query *entity.DNSQuery) error {
	if query.ID == uuid.Nil {
		query.ID = uuid.New()
	}
	query.CreatedAt = time.Now()
	query.UpdatedAt = time.Now()
	return s.telemetryRepo.CreateDNSQuery(ctx, query)
}

// GetDNSQuery retrieves a DNS query by ID
func (s *TelemetryService) GetDNSQuery(ctx context.Context, id uuid.UUID) (*entity.DNSQuery, error) {
	return s.telemetryRepo.GetDNSQuery(ctx, id)
}

// ListDNSQueries retrieves all DNS queries
func (s *TelemetryService) ListDNSQueries(ctx context.Context) ([]*entity.DNSQuery, error) {
	return s.telemetryRepo.ListDNSQueries(ctx)
}

// UpdateDNSQuery updates a DNS query
func (s *TelemetryService) UpdateDNSQuery(ctx context.Context, query *entity.DNSQuery) error {
	query.UpdatedAt = time.Now()
	return s.telemetryRepo.UpdateDNSQuery(ctx, query)
}

// DeleteDNSQuery deletes a DNS query
func (s *TelemetryService) DeleteDNSQuery(ctx context.Context, id uuid.UUID) error {
	return s.telemetryRepo.DeleteDNSQuery(ctx, id)
}

// Flow collector management

// CreateFlowCollector creates a new flow collector
func (s *TelemetryService) CreateFlowCollector(ctx context.Context, collector *entity.FlowCollector) error {
	if collector.ID == uuid.Nil {
		collector.ID = uuid.New()
	}
	collector.CreatedAt = time.Now()
	collector.UpdatedAt = time.Now()
	return s.telemetryRepo.CreateFlowCollector(ctx, collector)
}

// GetFlowCollector retrieves a flow collector by ID
func (s *TelemetryService) GetFlowCollector(ctx context.Context, id uuid.UUID) (*entity.FlowCollector, error) {
	return s.telemetryRepo.GetFlowCollector(ctx, id)
}

// ListFlowCollectors retrieves all flow collectors
func (s *TelemetryService) ListFlowCollectors(ctx context.Context) ([]*entity.FlowCollector, error) {
	return s.telemetryRepo.ListFlowCollectors(ctx)
}

// UpdateFlowCollector updates a flow collector
func (s *TelemetryService) UpdateFlowCollector(ctx context.Context, collector *entity.FlowCollector) error {
	collector.UpdatedAt = time.Now()
	return s.telemetryRepo.UpdateFlowCollector(ctx, collector)
}

// DeleteFlowCollector deletes a flow collector
func (s *TelemetryService) DeleteFlowCollector(ctx context.Context, id uuid.UUID) error {
	return s.telemetryRepo.DeleteFlowCollector(ctx, id)
}
