# Migration Plan: NetBox Extras to NetBox Go

This document tracks the migration and rewrite of the `extras` module from Python (Django) to Go.

## Goals
- Migrate `netbox/extras` functionality to `netbox_go/internal/domain/extras`.
- Implement using SOLID principles.
- Ensure high performance and type safety using Go 1.24+.
- Use SQLC for database interactions.

## Task List

### Phase 1: Analysis & Schema
- [x] Analyze `netbox/extras` models and logic
- [x] Define PostgreSQL schema for extras (`netbox_go/migrations/005_extras_schema.sql`)

### Phase 2: Domain Layer
- [ ] Define Go domain entities in `netbox_go/internal/domain/extras/entity`
- [ ] Implement Repository interfaces for `extras`

### Phase 3: Infrastructure Layer
- [ ] Implement SQLC queries in `netbox_go/internal/infrastructure/storage/sqlc/extras/queries.sql`
- [ ] Implement Repository implementations using SQLC

### Phase 4: Application Layer
- [ ] Implement Service layer for business logic (Event rules evaluation, Webhook triggering, etc.)
- [ ] Implement HTTP Handlers for REST API

### Phase 5: Integration & Testing
- [ ] Integrate `extras` module into `netbox_go/cmd/api/main.go`
- [ ] Write integration tests for `extras` functionality

## Progress
- **Current Status**: Starting Domain Layer implementation.
- **Completed**: Schema definition and initial analysis.
