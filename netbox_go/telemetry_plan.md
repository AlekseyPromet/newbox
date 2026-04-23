# Telemetry Module Implementation Plan

## Progress

### Phase 1: Core Infrastructure
- [x] Analyze existing project structure
- [x] Create comprehensive plan (newbox_gnmi_plan.md)
- [x] Create module_gnmi.go with fx dependencies
- [x] Create domain entities (Device, CollectionJob, TelemetryPoint, PingTarget, DNSQuery, FlowCollector)
- [x] Create repository interfaces
- [x] Create PostgreSQL migration 007_telemetry_schema.up.sql
- [x] Implement Vault client wrapper (infrastructure/vault/client.go)
- [x] Implement InfluxDB writer with batching (infrastructure/influxdb/writer.go)
- [x] Update go.mod with telemetry dependencies

### Phase 2: gNMI Client
- [x] Implement gNMI client with TLS support (infrastructure/gnmi/client.go)
- [x] Implement poll-based collector (infrastructure/gnmi/poller.go)
- [x] Implement path builders for telemetry types (Interface, BGP, System, OSPF)
- [x] Implement subscription-based collector (infrastructure/gnmi/subscriber.go)
- [x] Implement ICMP pinger (infrastructure/ping/pinger.go)
- [x] Implement DNS resolver (infrastructure/dns/resolver.go)
- [x] Implement NetFlow/sFlow collector (infrastructure/netflow/collector.go)

### Phase 3: Subscription Support
- [ ] Implement gNMI subscriber
- [ ] Handle stream reconnection
- [ ] Implement sample/heartbeat logic
- [ ] Unit tests for subscriber

### Phase 4: Distributed Coordination
- [ ] Implement coordinator service
- [ ] Implement collector service
- [ ] Add health monitoring and heartbeats
- [ ] Implement device assignment/rebalancing
- [ ] Integration tests

### Phase 5: API and Integration
- [ ] REST API for telemetry configuration
- [ ] REST API for telemetry data queries
- [ ] GraphQL API for telemetry
- [ ] Webhook/notification support
- [ ] Integration with existing DCIM module

### Phase 6: Production Hardening
- [ ] Add Prometheus metrics
- [ ] Implement circuit breakers
- [ ] Add retry logic with exponential backoff
- [ ] Performance testing
- [ ] Documentation

## Implementation Notes

- PostgreSQL: localhost:5432
- InfluxDB: localhost:8086
- Use sqlc with Docker for code generation
- Default locale: Russian (ru)
- All strings in code: English
