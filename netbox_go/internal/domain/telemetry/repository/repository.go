package repository

import (
	"context"

	"netbox_go/internal/domain/telemetry/entity"

	"github.com/google/uuid"
)

// TelemetryRepository defines the interface for telemetry data access
type TelemetryRepository interface {
	// Device operations
	CreateDevice(ctx context.Context, device *entity.TelemetryDevice) error
	GetDevice(ctx context.Context, id uuid.UUID) (*entity.TelemetryDevice, error)
	ListDevices(ctx context.Context) ([]*entity.TelemetryDevice, error)
	UpdateDevice(ctx context.Context, device *entity.TelemetryDevice) error
	DeleteDevice(ctx context.Context, id uuid.UUID) error

	// Collection operations
	CreateCollection(ctx context.Context, collection *entity.TelemetryCollection) error
	GetCollection(ctx context.Context, id uuid.UUID) (*entity.TelemetryCollection, error)
	ListCollections(ctx context.Context) ([]*entity.TelemetryCollection, error)
	UpdateCollection(ctx context.Context, collection *entity.TelemetryCollection) error
	DeleteCollection(ctx context.Context, id uuid.UUID) error

	// Job operations
	CreateJob(ctx context.Context, job *entity.CollectionJob) error
	GetJob(ctx context.Context, id uuid.UUID) (*entity.CollectionJob, error)
	ListJobs(ctx context.Context, limit int) ([]*entity.CollectionJob, error)
	UpdateJob(ctx context.Context, job *entity.CollectionJob) error

	// Ping target operations
	CreatePingTarget(ctx context.Context, target *entity.PingTarget) error
	GetPingTarget(ctx context.Context, id uuid.UUID) (*entity.PingTarget, error)
	ListPingTargets(ctx context.Context) ([]*entity.PingTarget, error)
	UpdatePingTarget(ctx context.Context, target *entity.PingTarget) error
	DeletePingTarget(ctx context.Context, id uuid.UUID) error

	// DNS query operations
	CreateDNSQuery(ctx context.Context, query *entity.DNSQuery) error
	GetDNSQuery(ctx context.Context, id uuid.UUID) (*entity.DNSQuery, error)
	ListDNSQueries(ctx context.Context) ([]*entity.DNSQuery, error)
	UpdateDNSQuery(ctx context.Context, query *entity.DNSQuery) error
	DeleteDNSQuery(ctx context.Context, id uuid.UUID) error

	// Flow collector operations
	CreateFlowCollector(ctx context.Context, collector *entity.FlowCollector) error
	GetFlowCollector(ctx context.Context, id uuid.UUID) (*entity.FlowCollector, error)
	ListFlowCollectors(ctx context.Context) ([]*entity.FlowCollector, error)
	UpdateFlowCollector(ctx context.Context, collector *entity.FlowCollector) error
	DeleteFlowCollector(ctx context.Context, id uuid.UUID) error
}
