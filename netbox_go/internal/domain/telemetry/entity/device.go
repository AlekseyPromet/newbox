package entity

import (
	"time"

	"github.com/google/uuid"
)

// TelemetryType represents the type of telemetry data
type TelemetryType string

const (
	TelemetryTypeInterface TelemetryType = "interface"
	TelemetryTypeBGP       TelemetryType = "bgp"
	TelemetryTypeOSPF      TelemetryType = "ospf"
	TelemetryTypeSystem    TelemetryType = "system"
	TelemetryTypePing      TelemetryType = "ping"
	TelemetryTypeDNS       TelemetryType = "dns"
	TelemetryTypeNetFlow   TelemetryType = "netflow"
)

// CollectionType represents how telemetry is collected
type CollectionType string

const (
	CollectionTypePoll      CollectionType = "poll"
	CollectionTypeSubscribe CollectionType = "subscribe"
	CollectionTypeBoth      CollectionType = "both"
)

// TelemetryDevice represents a device configured for telemetry collection
type TelemetryDevice struct {
	ID                    uuid.UUID      `json:"id"`
	DeviceID              uuid.UUID      `json:"device_id"`
	CollectionType        CollectionType `json:"collection_type"`
	AssignedCollectorID   string         `json:"assigned_collector_id"`
	LastCollectionAt      *time.Time     `json:"last_collection_at"`
	LastCollectionStatus  string         `json:"last_collection_status"`
	CollectionErrorsCount int            `json:"collection_errors_count"`
	VaultSecretPath       string         `json:"vault_secret_path"`
	GNMIAddress           string         `json:"gnmi_address"`
	GNmiPort              int            `json:"gnmi_port"`
	Enabled               bool           `json:"enabled"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

// TelemetryCollection represents a telemetry collection configuration
type TelemetryCollection struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CollectionType  CollectionType `json:"collection_type"`
	TelemetryType   TelemetryType  `json:"telemetry_type"`
	TargetPath      string         `json:"target_path"`
	IntervalSeconds int            `json:"interval_seconds"`
	Enabled         bool           `json:"enabled"`
	Filters         string         `json:"filters"` // JSON string
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// CollectionJob represents a telemetry collection job
type CollectionJob struct {
	ID               uuid.UUID  `json:"id"`
	JobID            string     `json:"job_id"`
	CollectorID      string     `json:"collector_id"`
	DeviceID         uuid.UUID  `json:"device_id"`
	CollectionID     uuid.UUID  `json:"collection_id"`
	Status           string     `json:"status"`
	StartedAt        *time.Time `json:"started_at"`
	CompletedAt      *time.Time `json:"completed_at"`
	ErrorMessage     string     `json:"error_message"`
	RecordsCollected int        `json:"records_collected"`
	CreatedAt        time.Time  `json:"created_at"`
}

// PingTarget represents an ICMP ping target
type PingTarget struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        uuid.UUID `json:"device_id"`
	TargetAddress   string    `json:"target_address"`
	TargetType      string    `json:"target_type"` // icmp, udp, tcp
	IntervalSeconds int       `json:"interval_seconds"`
	PacketCount     int       `json:"packet_count"`
	PacketSize      int       `json:"packet_size"`
	TimeoutSeconds  int       `json:"timeout_seconds"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DNSQuery represents a DNS query configuration
type DNSQuery struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        uuid.UUID `json:"device_id"`
	QueryName       string    `json:"query_name"`
	QueryType       string    `json:"query_type"` // A, AAAA, MX, NS, TXT, CNAME
	DNSServer       string    `json:"dns_server"`
	IntervalSeconds int       `json:"interval_seconds"`
	TimeoutSeconds  int       `json:"timeout_seconds"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FlowCollector represents a NetFlow/sFlow collector configuration
type FlowCollector struct {
	ID                  uuid.UUID `json:"id"`
	DeviceID            uuid.UUID `json:"device_id"`
	CollectorType       string    `json:"collector_type"` // netflow_v5, netflow_v9, ipfix, sflow
	ListeningPort       int       `json:"listening_port"`
	SamplingRate        int       `json:"sampling_rate"`
	AggregationInterval int       `json:"aggregation_interval"`
	Enabled             bool      `json:"enabled"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// TelemetryPoint represents a single telemetry data point
type TelemetryPoint struct {
	Timestamp time.Time
	Device    string
	Site      string
	Type      TelemetryType
	Tags      map[string]string
	Fields    map[string]interface{}
}
