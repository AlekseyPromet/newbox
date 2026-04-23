-- Migration: 007_telemetry_schema.up.sql
-- Description: Telemetry collection schema for gNMI, ICMP, DNS, and NetFlow

-- Table: telemetry_collections
CREATE TABLE IF NOT EXISTS telemetry_collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    collection_type VARCHAR(50) NOT NULL, -- 'poll' or 'subscribe'
    telemetry_type VARCHAR(50) NOT NULL, -- 'interface', 'bgp', 'ospf', 'system', 'ping', 'dns', 'netflow'
    target_path TEXT NOT NULL, -- gNMI path pattern
    interval_seconds INTEGER DEFAULT 60,
    enabled BOOLEAN DEFAULT true,
    filters JSONB, -- device/site filters
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_collections_type ON telemetry_collections(telemetry_type);
CREATE INDEX idx_telemetry_collections_enabled ON telemetry_collections(enabled);

-- Table: telemetry_devices
CREATE TABLE IF NOT EXISTS telemetry_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    collection_type VARCHAR(50) NOT NULL, -- 'poll' or 'subscribe' or 'both'
    assigned_collector_id VARCHAR(255),
    last_collection_at TIMESTAMPTZ,
    last_collection_status VARCHAR(50),
    collection_errors_count INTEGER DEFAULT 0,
    vault_secret_path VARCHAR(500),
    gnmi_address VARCHAR(255) NOT NULL,
    gnmi_port INTEGER DEFAULT 57400,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_devices_device_id ON telemetry_devices(device_id);
CREATE INDEX idx_telemetry_devices_collector ON telemetry_devices(assigned_collector_id);
CREATE INDEX idx_telemetry_devices_enabled ON telemetry_devices(enabled);

-- Table: telemetry_collection_jobs
CREATE TABLE IF NOT EXISTS telemetry_collection_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id VARCHAR(255) NOT NULL,
    collector_id VARCHAR(255) NOT NULL,
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    collection_id UUID REFERENCES telemetry_collections(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL, -- 'pending', 'running', 'completed', 'failed'
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    records_collected INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_jobs_collector ON telemetry_collection_jobs(collector_id);
CREATE INDEX idx_telemetry_jobs_device ON telemetry_collection_jobs(device_id);
CREATE INDEX idx_telemetry_jobs_status ON telemetry_collection_jobs(status);
CREATE INDEX idx_telemetry_jobs_created ON telemetry_collection_jobs(created_at DESC);

-- Table: telemetry_ping_targets
CREATE TABLE IF NOT EXISTS telemetry_ping_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    target_address VARCHAR(255) NOT NULL, -- IP or hostname
    target_type VARCHAR(50) DEFAULT 'icmp', -- 'icmp', 'udp', 'tcp'
    interval_seconds INTEGER DEFAULT 30,
    packet_count INTEGER DEFAULT 5,
    packet_size INTEGER DEFAULT 64,
    timeout_seconds INTEGER DEFAULT 5,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_ping_device ON telemetry_ping_targets(device_id);
CREATE INDEX idx_telemetry_ping_enabled ON telemetry_ping_targets(enabled);

-- Table: telemetry_dns_queries
CREATE TABLE IF NOT EXISTS telemetry_dns_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    query_name VARCHAR(255) NOT NULL, -- Domain to resolve
    query_type VARCHAR(10) NOT NULL, -- 'A', 'AAAA', 'MX', 'NS', 'TXT', 'CNAME'
    dns_server VARCHAR(255), -- Optional specific DNS server
    interval_seconds INTEGER DEFAULT 300,
    timeout_seconds INTEGER DEFAULT 10,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_dns_device ON telemetry_dns_queries(device_id);
CREATE INDEX idx_telemetry_dns_enabled ON telemetry_dns_queries(enabled);

-- Table: telemetry_flow_collectors
CREATE TABLE IF NOT EXISTS telemetry_flow_collectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    collector_type VARCHAR(50) NOT NULL, -- 'netflow_v5', 'netflow_v9', 'ipfix', 'sflow'
    listening_port INTEGER NOT NULL, -- 2055 (NetFlow), 6343 (sFlow)
    sampling_rate INTEGER DEFAULT 1, -- Flow sampling rate
    aggregation_interval INTEGER DEFAULT 60, -- Seconds to aggregate flows
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_telemetry_flows_device ON telemetry_flow_collectors(device_id);
CREATE INDEX idx_telemetry_flows_enabled ON telemetry_flow_collectors(enabled);
CREATE UNIQUE INDEX idx_telemetry_flows_device_port ON telemetry_flow_collectors(device_id, listening_port);

-- Comments
COMMENT ON TABLE telemetry_collections IS 'Telemetry collection configurations';
COMMENT ON TABLE telemetry_devices IS 'Devices enabled for telemetry collection';
COMMENT ON TABLE telemetry_collection_jobs IS 'Historical record of collection jobs';
COMMENT ON TABLE telemetry_ping_targets IS 'ICMP ping targets';
COMMENT ON TABLE telemetry_dns_queries IS 'DNS query configurations';
COMMENT ON TABLE telemetry_flow_collectors IS 'NetFlow/sFlow collector configurations';
