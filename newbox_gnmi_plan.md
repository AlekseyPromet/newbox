# gNMI Telemetry Module Implementation Plan

## 1. Overview

A distributed telemetry collection module for network devices supporting:
- **Comprehensive telemetry**: Interface statistics, BGP/OSPF routing metrics, device operational state, ICMP ping, DNS, NetFlow
- **Collection modes**: Both poll-based and subscribe-based (gNMI subscriptions)
- **Protocol support**: gNMI for network devices, ICMP for reachability, DNS for resolution, NetFlow/sFlow for flow analytics
- **Credential management**: HashiCorp Vault integration
- **Storage**: InfluxDB with separate measurements per telemetry type
- **Scale**: 1000+ devices with distributed collectors and central coordination

---

## 2. Architecture

### 2.1 High-Level Design

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Central Coordinator                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ Task        │  │ Device      │  │ Collection  │  │ InfluxDB    │     │
│  │ Scheduler   │  │ Registry    │  │ Coordinator │  │ Writer      │     │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────────────────────┘
         │                   │                   │                 │
         ▼                   ▼                   ▼                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Distributed Collectors                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                      │
│  │ Collector 1 │  │ Collector 2 │  │ Collector N │  ...                │
│  │ (Region A)  │  │ (Region B)  │  │ (Region N)  │                      │
│  └─────────────┘  └─────────────┘  └─────────────┘                      │
└─────────────────────────────────────────────────────────────────────────┘
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     External Dependencies                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ HashiCorp   │  │ InfluxDB    │  │ PostgreSQL  │  │ gNMI        │     │
│  │ Vault       │  │ Cluster     │  │ (NetBox DB) │  │ Devices     │     │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Module Structure

Following domain-driven architecture pattern from `netbox_go`:

```
netbox_go/
├── internal/
│   ├── app/
│   │   └── module_gnmi.go           # fx module definition
│   ├── domain/
│   │   └── telemetry/
│   │       ├── entity/
│   │       │   ├── device.go        # Telemetry-enabled device entity
│   │       │   ├── collection_job.go # Collection job entity
│   │       │   └── telemetry_point.go # Raw telemetry data point
│   │       ├── repository/
│   │       │   └── repository.go    # Repository interfaces
│   │       └── service/
│   │           └── telemetry_service.go
│   ├── delivery/
│   │   └── http/
│   │       └── handlers/
│   │           └── telemetry_handler.go
│   ├── infrastructure/
│   │   ├── gnmi/
│   │   │   ├── client.go            # gNMI client wrapper
│   │   │   ├── collector.go         # Collector implementations
│   │   │   ├── poller.go            # Poll-based collector
│   │   │   └── subscriber.go        # Subscription-based collector
│   │   ├── ping/
│   │   │   ├── pinger.go            # ICMP ping executor
│   │   │   └── parser.go            # Ping result parser
│   │   ├── dns/
│   │   │   └── resolver.go         # DNS resolver and metrics collector
│   │   ├── netflow/
│   │   │   ├── collector.go         # NetFlow/sFlow collector
│   │   │   ├── decoder.go           # NetFlow v5/v9/IPFIX decoder
│   │   │   └── aggregator.go        # Flow aggregation logic
│   │   ├── vault/
│   │   │   └── vault_client.go      # HashiCorp Vault client
│   │   └── influxdb/
│   │       └── writer.go            # InfluxDB writer
│   └── repository/
│       └── postgres/
│           └── telemetry_repository.go
├── migrations/
│   └── 007_telemetry_schema.up.sql
└── cmd/
    ├── gnmi_collector/
    │   └── main.go                  # Standalone collector binary
    └── telemetry_collector/
        └── main.go                  # Unified telemetry collector binary
```

---

## 3. InfluxDB Schema Design

### 3.1 Measurements

Separate measurements per telemetry type for query performance:

| Measurement | Description | Key Tags | Fields |
|-------------|-------------|----------|--------|
| `interface_stats` | Interface counters and statistics | `device`, `interface`, `device_type`, `site` | `in_octets`, `out_octets`, `in_pkts`, `out_pkts`, `in_errors`, `out_errors`, `oper_status`, `admin_status`, `speed` |
| `bgp_metrics` | BGP neighbor and route metrics | `device`, `neighbor`, `remote_as`, `device_type`, `site` | `Established_state`, `prefixes_received`, `prefixes_sent`, `uptime_seconds`, `messages_received`, `messages_sent` |
| `ospf_metrics` | OSPF neighbor and area metrics | `device`, `neighbor`, `area_id`, `device_type`, `site` | `neighbor_state`, `dead_timer`, `lsdb_count`, `adjacency_uptime` |
| `system_metrics` | Device operational state | `device`, `device_type`, `vendor`, `site` | `cpu_percent`, `memory_percent`, `memory_total`, `memory_used`, `uptime_seconds`, `temperature`, `fan_status`, `power_supply_status` |
| `icmp_ping` | ICMP ping reachability metrics | `device`, `target`, `source`, `device_type`, `site` | `rtt_ms`, `rtt_min_ms`, `rtt_max_ms`, `rtt_avg_ms`, `packet_loss_percent`, `packets_sent`, `packets_received`, `ttl` |
| `dns_metrics` | DNS resolution metrics | `device`, `query_name`, `query_type`, `dns_server`, `device_type`, `site` | `query_time_ms`, `resolve_time_ms`, `answer_count`, ` NXDOMAIN`, `SERVFAIL`, `timeout_count` |
| `netflow_records` | NetFlow/sFlow flow data | `device`, `src_addr`, `dst_addr`, `src_port`, `dst_port`, `protocol`, `device_type`, `site` | `packets`, `bytes`, `flow_start_ms`, `flow_end_ms`, `tcp_flags`, `tos`, `ingress_if`, `egress_if` |

### 3.2 Tag Design for High Cardinality

```
Tags: device, device_type, vendor, model, site, region, interface, neighbor, remote_as, area_id
Fields: all numeric values and status enums
Timestamp: nanosecond precision from gNMI timestamp
```

---

## 4. Collection Modes

### 4.1 Poll-Based Collection

```
┌──────────────────────────────────────────────────────┐
│ Poll Scheduler (configurable intervals)              │
│   - Interface stats: 60s default                     │
│   - BGP metrics: 300s default                        │
│   - System metrics: 30s default                      │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Target Selector                                      │
│   - Filter by device tags                            │
│   - Filter by site/region                            │
│   - Filter by collection mode = poll                 │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ gNMI Client Pool                                     │
│   - Connection复用 (TLS, auth)                       │
│   - Rate limiting per device                         │
│   - Timeout: 30s default                            │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Path Builder                                         │
│   - /interfaces/interface[name={if}]/state/         │
│   - /network-instances/network-instance[name=default]/protocols/ │
│   - /system/resources/                               │
└──────────────────────────────────────────────────────┘
```

### 4.2 Subscribe-Based Collection

```
┌──────────────────────────────────────────────────────┐
│ Subscription Manager                                 │
│   - ONCE: single collection                          │
│   - STREAM: ongoing updates                          │
│   - POLL: on-demand per subscription                │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Stream Handler                                       │
│   - gRPC stream management                          │
│   - Reconnection logic                              │
│   - Heartbeat monitoring                            │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Sample Collector                                     │
│   - Sample interval: 10s default                     │
│   - Heartbeat interval: 60s                         │
└──────────────────────────────────────────────────────┘
```

### 4.3 ICMP Ping Collection

```
┌──────────────────────────────────────────────────────┐
│ Ping Scheduler                                       │
│   - Interval: 30s default (configurable per target) │
│   - Burst mode for initial discovery                │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Ping Executor                                        │
│   - Raw socket ICMP (requires root/CAP_NET_RAW)     │
│   - Fallback to UDP ping for restricted envs       │
│   - Parallel execution with worker pool            │
│   - Configurable packet count, payload size, TTL    │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Statistics Calculator                                │
│   - RTT min/max/avg/stddev                          │
│   - Packet loss percentage                          │
│   - Jitter calculation                               │
│   - MOS score estimation                            │
└──────────────────────────────────────────────────────┘
```

### 4.4 DNS Resolution Collection

```
┌──────────────────────────────────────────────────────┐
│ DNS Query Scheduler                                  │
│   - Configurable interval per query                 │
│   - Query types: A, AAAA, MX, NS, TXT, CNAME        │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ DNS Resolver                                         │
│   - Parallel resolution with worker pool           │
│   - Custom DNS servers support                      │
│   - EDNS0 support                                   │
│   - DNSSEC validation (optional)                   │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Metrics Extractor                                    │
│   - Query time (DNS server latency)                 │
│   - Resolution time (total time)                    │
│   - Answer count per record type                    │
│   - Error classification (NXDOMAIN, SERVFAIL, etc) │
└──────────────────────────────────────────────────────┘
```

### 4.5 NetFlow/sFlow Collection

```
┌──────────────────────────────────────────────────────┐
│ Flow Collector Listener                              │
│   - UDP ports: 2055 (NetFlow), 6343 (sFlow)         │
│   - Multiple collectors per device                   │
│   - Template-based decoding (NetFlow v9/IPFIX)      │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Flow Decoder                                          │
│   - NetFlow v5 (fixed format)                       │
│   - NetFlow v9 (template-based)                     │
│   - IPFIX (RFC 7011)                                │
│   - sFlow v5                                       │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Flow Aggregator                                      │
│   - Time-based aggregation (1min, 5min, 15min)    │
│   - Key-based aggregation (src/dst AS, prefix)     │
│   - Traffic matrix calculation                       │
└──────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│ Flow Exporter                                        │
│   - Batching for InfluxDB writes                    │
│   - Circuit breaker on write failures               │
└──────────────────────────────────────────────────────┘
```

---

## 5. Credential Management (HashiCorp Vault)

### 5.1 Vault Integration

```
┌──────────────────────────────────────────────────────┐
│ Vault Client                                          │
│   - AppRole authentication                           │
│   - Token caching with renewal                       │
│   - Connection pooling                              │
└──────────────────────────────────────────────────────┘
```

### 5.2 Secret Paths Structure

```
secret/data/netbox/gnmi/devices/{device_id}
  ├── username
  ├── password
  └── priv_key (optional, for cert-based auth)

secret/data/netbox/gnmi/collectors/{collector_id}
  └── api_key
```

### 5.3 Credential Caching

- Local in-memory cache with TTL (5 minutes default)
- Cache invalidation on Vault lease renewal
- Graceful fallback to Vault on cache miss

---

## 6. Distributed Collector Design

### 6.1 Central Coordinator Responsibilities

| Component | Responsibility |
|-----------|----------------|
| Task Scheduler | Schedule collection jobs, distribute across collectors |
| Device Registry | Track device locations, collector assignments |
| Collection Coordinator | Route collection requests to appropriate collector |
| Health Monitor | Monitor collector heartbeats, detect failures |

### 6.2 Collector Responsibilities

| Component | Responsibility |
|-----------|----------------|
| gNMI Client Pool | Manage device connections |
| Local Scheduler | Execute scheduled collections |
| Data Buffer | Queue data during InfluxDB outages |
| Heartbeat Reporter | Report health to coordinator |

### 6.3 Communication Protocol

```
Collector <-> Coordinator:
  - gRPC for control plane
  - HTTP/REST for metrics push to InfluxDB

Coordinator -> Collector:
  - "Execute collection job" command
  - "Update device list" command
  - "Shutdown" command

Collector -> Coordinator:
  - Heartbeat every 30s
  - Collection results/status
  - Health metrics
```

### 6.4 Device Assignment Strategy

```
Assignment by site/region:
  - Collector registers with sites it can reach
  - Coordinator assigns devices by site proximity

Load balancing:
  - Max devices per collector: 500 (configurable)
  - Rebalance on collector addition/removal
  - Anti-affinity: critical devices spread across collectors
```

---

## 7. Implementation Phases

### Phase 1: Core Infrastructure (Weeks 1-2)

- [ ] Create `module_gnmi.go` with fx dependencies
- [ ] Implement Vault client wrapper
- [ ] Implement InfluxDB writer with batching
- [ ] Create domain entities (`Device`, `CollectionJob`, `TelemetryPoint`)
- [ ] Create PostgreSQL schema for telemetry configuration
- [ ] Implement repository interfaces and basic repository

### Phase 2: gNMI Client (Weeks 3-4)

- [ ] Implement gNMI client with TLS support
- [ ] Implement poll-based collector
- [ ] Implement path builders for each telemetry type
- [ ] Add connection pooling and rate limiting
- [ ] Unit tests for gNMI client

### Phase 3: Subscription Support (Week 5)

- [ ] Implement gNMI subscriber
- [ ] Handle stream reconnection
- [ ] Implement sample/heartbeat logic
- [ ] Unit tests for subscriber

### Phase 4: Distributed Coordination (Weeks 6-7)

- [ ] Implement coordinator service
- [ ] Implement collector service
- [ ] Add health monitoring and heartbeats
- [ ] Implement device assignment/rebalancing
- [ ] Integration tests

### Phase 5: API and Integration (Week 8)

- [ ] REST API for telemetry configuration
- [ ] Webhook/notification support for alerts
- [ ] Grafana dashboard templates
- [ ] Integration with existing DCIM module

### Phase 6: Production Hardening (Weeks 9-10)

- [ ] Add Prometheus metrics
- [ ] Implement circuit breakers
- [ ] Add retry logic with exponential backoff
- [ ] Performance testing with 1000+ devices
- [ ] Documentation

---

## 8. Dependencies

### New Go Modules

```go
require (
    // Existing
    go.uber.org/fx v1.24.0
    github.com/labstack/echo/v4 v4.15.1
    
    // gNMI
    github.com/openconfig/gnmi v0.0.0-20240601-xxxx  # gNMI protocol
    
    // InfluxDB
    github.com/influxdata/influxdb-client-go/v2 v2.x  # InfluxDB client
    
    // Vault
    github.com/hashicorp/vault/api v1.12.0            # Vault client
    
    // gRPC
    google.golang.org/grpc v1.60.0                    # gRPC for collector comms
    
    // Scheduling
    github.com/robfig/cron/v3 v3.0.1                  # Cron scheduling
    
    // ICMP Ping
    github.com/go-ping/ping v1.11.0                   # ICMP ping implementation
    golang.org/x/net icmp                             # Raw ICMP socket
    
    // DNS
    github.com/miekg/dns v1.11.1                      # DNS client library
    
    // NetFlow
    github.com/netsampler/goflow2 v1.x                # NetFlow/sFlow collector
    
    // Utilities
    github.com/google/gopacket v1.1.19                # Packet processing
)
```

---

## 9. Configuration

### Environment Variables

```bash
# Vault
VAULT_ADDR=https://vault.example.com:8200
VAULT_APPROLE_ROLE_ID=xxx
VAULT_APPROLE_SECRET_ID=xxx
VAULT_TIMEOUT=10s
VAULT_CACHE_TTL=5m

# InfluxDB
INFLUXDB_URL=https://influxdb.example.com:8086
INFLUXDB_TOKEN=xxx
INFLUXDB_ORG=netbox
INFLUXDB_BUCKET=telemetry
INFLUXDB_BATCH_SIZE=5000
INFLUXDB_FLUSH_INTERVAL=10s

# Collector
COLLECTOR_ID=collector-01
COORDINATOR_ADDR=coordinator.example.com:8080
COLLECTOR_PORT=9090
COLLECTION_POOL_SIZE=100

# gNMI Defaults
GNMI_TIMEOUT=30s
GNMI_TLS_ENABLED=true
GNMI_TLS_CERT_PATH=/etc/gnmi/cert.pem
GNMI_TLS_KEY_PATH=/etc/gnmi/key.pem
GNMI_TLS_CA_PATH=/etc/gnmi/ca.pem

# ICMP Ping
PING_ENABLED=true
PING_TIMEOUT=5s
PING_INTERVAL=30s
PING_COUNT=5
PING_PACKET_SIZE=64
PING_WORKERS=50

# DNS
DNS_ENABLED=true
DNS_TIMEOUT=10s
DNS_INTERVAL=300s
DNS_WORKERS=20
DNS_SERVERS=8.8.8.8,1.1.1.1

# NetFlow/sFlow
NETFLOW_ENABLED=true
NETFLOW_PORT=2055
SFLOW_PORT=6343
NETFLOW_WORKERS=10
NETFLOW_BUFFER_SIZE=10000
```

---

## 10. Database Schema (PostgreSQL)

### Table: `telemetry_collections`

```sql
CREATE TABLE telemetry_collections (
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
```

### Table: `telemetry_devices`

```sql
CREATE TABLE telemetry_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id),
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
```

### Table: `telemetry_ping_targets`

```sql
CREATE TABLE telemetry_ping_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id),
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
```

### Table: `telemetry_dns_queries`

```sql
CREATE TABLE telemetry_dns_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id),
    query_name VARCHAR(255) NOT NULL, -- Domain to resolve
    query_type VARCHAR(10) NOT NULL, -- 'A', 'AAAA', 'MX', 'NS', 'TXT', 'CNAME'
    dns_server VARCHAR(255), -- Optional specific DNS server
    interval_seconds INTEGER DEFAULT 300,
    timeout_seconds INTEGER DEFAULT 10,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Table: `telemetry_flow_collectors`

```sql
CREATE TABLE telemetry_flow_collectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id),
    collector_type VARCHAR(50) NOT NULL, -- 'netflow_v5', 'netflow_v9', 'ipfix', 'sflow'
    listening_port INTEGER NOT NULL, -- 2055 (NetFlow), 6343 (sFlow)
    sampling_rate INTEGER DEFAULT 1, -- Flow sampling rate
    aggregation_interval INTEGER DEFAULT 60, -- Seconds to aggregate flows
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Table: `telemetry_collection_jobs`

```sql
CREATE TABLE telemetry_collection_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id VARCHAR(255) NOT NULL,
    collector_id VARCHAR(255) NOT NULL,
    device_id UUID NOT NULL REFERENCES dcim_devices(id),
    collection_id UUID NOT NULL REFERENCES telemetry_collections(id),
    status VARCHAR(50) NOT NULL, -- 'pending', 'running', 'completed', 'failed'
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    records_collected INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 11. API Endpoints

### Telemetry Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/collections` | List all collection configs |
| POST | `/api/v1/telemetry/collections` | Create collection config |
| GET | `/api/v1/telemetry/collections/:id` | Get collection details |
| PUT | `/api/v1/telemetry/collections/:id` | Update collection |
| DELETE | `/api/v1/telemetry/collections/:id` | Delete collection |
| POST | `/api/v1/telemetry/collections/:id/test` | Test collection on device |

### Device Assignment

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/devices` | List telemetry-enabled devices |
| POST | `/api/v1/telemetry/devices` | Enable device for telemetry |
| GET | `/api/v1/telemetry/devices/:id` | Get device config |
| PUT | `/api/v1/telemetry/devices/:id` | Update device config |
| DELETE | `/api/v1/telemetry/devices/:id` | Disable device telemetry |

### Job Status

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/jobs` | List recent collection jobs |
| GET | `/api/v1/telemetry/jobs/:id` | Get job details |
| POST | `/api/v1/telemetry/jobs/:id/retry` | Retry failed job |

### Collector Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/collectors` | List active collectors |
| GET | `/api/v1/telemetry/collectors/:id/status` | Get collector health |

### Ping Targets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/ping` | List ping targets |
| POST | `/api/v1/telemetry/ping` | Add ping target |
| GET | `/api/v1/telemetry/ping/:id` | Get ping target details |
| PUT | `/api/v1/telemetry/ping/:id` | Update ping target |
| DELETE | `/api/v1/telemetry/ping/:id` | Delete ping target |
| POST | `/api/v1/telemetry/ping/:id/test` | Test ping to target |

### DNS Queries

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/dns` | List DNS queries |
| POST | `/api/v1/telemetry/dns` | Add DNS query |
| GET | `/api/v1/telemetry/dns/:id` | Get DNS query details |
| PUT | `/api/v1/telemetry/dns/:id` | Update DNS query |
| DELETE | `/api/v1/telemetry/dns/:id` | Delete DNS query |
| POST | `/api/v1/telemetry/dns/:id/test` | Test DNS resolution |

### Flow Collectors

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/flows` | List flow collectors |
| POST | `/api/v1/telemetry/flows` | Add flow collector |
| GET | `/api/v1/telemetry/flows/:id` | Get flow collector details |
| PUT | `/api/v1/telemetry/flows/:id` | Update flow collector |
| DELETE | `/api/v1/telemetry/flows/:id` | Delete flow collector |
| GET | `/api/v1/telemetry/flows/:id/stats` | Get flow collection stats |

### Telemetry Data (REST API)

Query stored telemetry data from InfluxDB:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/telemetry/data/interface` | Query interface statistics |
| GET | `/api/v1/telemetry/data/bgp` | Query BGP metrics |
| GET | `/api/v1/telemetry/data/ospf` | Query OSPF metrics |
| GET | `/api/v1/telemetry/data/system` | Query system metrics |
| GET | `/api/v1/telemetry/data/ping` | Query ICMP ping results |
| GET | `/api/v1/telemetry/data/dns` | Query DNS resolution results |
| GET | `/api/v1/telemetry/data/flows` | Query NetFlow/sFlow records |

#### Query Parameters (all data endpoints)

| Parameter | Type | Description |
|-----------|------|-------------|
| `device` | string | Filter by device name or ID |
| `site` | string | Filter by site |
| `start` | timestamp | Start time (RFC3339) |
| `end` | timestamp | End time (RFC3339) |
| `interval` | string | Aggregation interval (e.g., "1m", "5m", "1h") |
| `limit` | integer | Max records to return (default 1000) |

---

## 11b. GraphQL API

GraphQL schema for flexible telemetry queries:

```graphql
type Query {
  # Interface Statistics
  interfaceStats(
    device: String
    interface: String
    site: String
    start: Time!
    end: Time!
    interval: String
  ): [InterfaceStat!]!

  # BGP Metrics
  bgpMetrics(
    device: String
    neighbor: String
    site: String
    start: Time!
    end: Time!
  ): [BGPMetric!]!

  # System Metrics
  systemMetrics(
    device: String!
    start: Time!
    end: Time!
    interval: String
  ): [SystemMetric!]!

  # Ping Results
  pingResults(
    device: String
    target: String
    start: Time!
    end: Time!
  ): [PingResult!]!

  # DNS Results
  dnsResults(
    device: String
    queryName: String
    start: Time!
    end: Time!
  ): [DNSResult!]!

  # Flow Records
  flowRecords(
    device: String
    srcAddress: String
    dstAddress: String
    start: Time!
    end: Time!
    limit: Int
  ): [FlowRecord!]!

  # Device Health Summary
  deviceHealth(device: String!): DeviceHealth!

  # Telemetry Alerts
  telemetryAlerts(
    severity: AlertSeverity
    start: Time!
    end: Time!
  ): [TelemetryAlert!]!
}

# Types
type InterfaceStat {
  timestamp: Time!
  device: String!
  interface: String!
  site: String
  inOctets: Uint64!
  outOctets: Uint64!
  inPkts: Uint64!
  outPkts: Uint64!
  inErrors: Int!
  outErrors: Int!
  operStatus: String!
  adminStatus: String!
  speed: Int!
}

type BGPMetric {
  timestamp: Time!
  device: String!
  neighbor: String!
  remoteAs: Int!
  establishedState: String!
  prefixesReceived: Int!
  prefixesSent: Int!
  uptimeSeconds: Int!
}

type SystemMetric {
  timestamp: Time!
  device: String!
  cpuPercent: Float!
  memoryPercent: Float!
  memoryUsed: Uint64!
  memoryTotal: Uint64!
  uptimeSeconds: Int!
  temperature: Float
}

type PingResult {
  timestamp: Time!
  device: String!
  target: String!
  rttMs: Float!
  rttMinMs: Float!
  rttMaxMs: Float!
  rttAvgMs: Float!
  packetLossPercent: Float!
  packetsSent: Int!
  packetsReceived: Int!
  ttl: Int!
}

type DNSResult {
  timestamp: Time!
  device: String!
  queryName: String!
  queryType: String!
  dnsServer: String
  queryTimeMs: Float!
  resolveTimeMs: Float!
  answerCount: Int!
  NXDOMAIN: Boolean!
  SERVFAIL: Boolean!
}

type FlowRecord {
  timestamp: Time!
  device: String!
  srcAddr: String!
  dstAddr: String!
  srcPort: Int!
  dstPort: Int!
  protocol: String!
  packets: Uint64!
  bytes: Uint64!
  flowStartMs: Time!
  flowEndMs: Time!
  tcpFlags: String
  tos: Int
}

type DeviceHealth {
  device: String!
  overallStatus: String!
  lastSeen: Time!
  gnmiStatus: String
  pingStatus: String
  dnsStatus: String
  flowStatus: String
  activeAlerts: [TelemetryAlert!]!
}

type TelemetryAlert {
  id: ID!
  timestamp: Time!
  device: String!
  severity: AlertSeverity!
  metric: String!
  message: String!
  value: Float
  threshold: Float
}

enum AlertSeverity {
  INFO
  WARNING
  CRITICAL
}

# Mutations for real-time operations
type Mutation {
  # Trigger immediate collection
  triggerCollection(device: String!, collectionType: String!): CollectionJob!

  # Acknowledge alert
  acknowledgeAlert(alertId: ID!): TelemetryAlert!

  # Update collection interval
  updateCollectionInterval(device: String!, collectionType: String!, intervalSeconds: Int!): Boolean!
}
```

---

## 12. Error Handling Strategy

### Retry Logic

| Error Type | Retry Strategy |
|------------|----------------|
| Network timeout | Exponential backoff: 1s, 2s, 4s, 8s, max 60s |
| Authentication failure | Retry after Vault token refresh |
| Device unreachable | Mark device, retry after 5 min, alert after 3 failures |
| InfluxDB write failure | Buffer in memory, retry with backoff, alert if > 10k points |

### Circuit Breaker

- Open after 5 consecutive failures to a device
- Half-open after 30 seconds
- Close after 3 successful requests

---

## 13. Monitoring and Observability

### Prometheus Metrics

```
# Collection metrics
gnmi_collection_total{collector, device, type, status}
gnmi_collection_duration_seconds{collector, device, type}
gnmi_collection_records{collector, device, type}

# Device metrics
gnmi_device_up{collector, device}
gnmi_device_last_collection_age_seconds{collector, device}

# ICMP Ping metrics
ping_total{collector, device, target, status}
ping_rtt_seconds{collector, device, target, quantile}
ping_packet_loss_percent{collector, device, target}

# DNS metrics
dns_query_total{collector, device, query_name, query_type, status}
dns_query_duration_seconds{collector, device, query_name, query_type, quantile}
dns_resolution_failure_total{collector, device, query_name, error_type}

# NetFlow metrics
netflow_records_total{collector, device, version}
netflow_flows_aggregated_total{collector, device}
netflow_buffer_size{collector}
netflow_decoder_errors_total{collector, device, error_type}

# InfluxDB metrics
influxdb_write_total{status}
influxdb_write_duration_seconds
influxdb_buffer_size

# Vault metrics
vault_cache_hits_total
vault_cache_misses_total
vault_request_duration_seconds
```

### Health Checks

- `/healthz` - Liveness probe
- `/readyz` - Readiness probe (checks InfluxDB + Vault connectivity)

---

## 14. Security Considerations

1. **TLS everywhere** - gNMI, InfluxDB, Vault, inter-collector communication
2. **Vault for secrets** - No plaintext credentials in config or DB
3. **Role-based access** - Collectors get minimal Vault permissions
4. **Audit logging** - Log all credential access
5. **Network segmentation** - Collectors can only reach managed devices

---

## 15. Migration Strategy

1. **Phase 1**: Deploy collectors without data collection (registration only)
2. **Phase 2**: Enable polling for 10% of devices, validate data quality
3. **Phase 3**: Gradually migrate remaining devices
4. **Phase 4**: Enable subscriptions for real-time telemetry

---

## 16. Open Questions / TODOs

### gNMI Related
- [ ] Confirm gNMI dialect (OpenConfig, Cisco NSO, Juniper) - we need path mappings
- [ ] Define gNMI encoding preference JSON_IETF

### ICMP Ping Related
- [ ] Determine if raw sockets are available or UDP fallback is needed
- [ ] Define packet size strategy (default 64 bytes vs jumbo frames)
- [ ] Plan for distributed ping execution (ping from collectors vs centralized)

### DNS Related
- [ ] Determine default DNS servers to use
- [ ] Plan for DNSSEC validation enablement
- [ ] Define cache strategy for repeated queries

### NetFlow Related
- [ ] Confirm NetFlow version support (v5, v9, IPFIX)
- [ ] Define flow aggregation strategy (time-based, count-based)
- [ ] Plan for NetFlow template handling and caching
- [ ] Determine maximum flows to buffer during InfluxDB outages

### General
- [ ] Determine InfluxDB retention policy 7d hot
- [ ] Define alert thresholds for telemetry anomalies
- [ ] Plan for data aggregation/rollup in InfluxDB
- [ ] Consider gRPC for collector-coordinator communication
- [ ] Define collector discovery mechanism (static config vs dynamic registration)
