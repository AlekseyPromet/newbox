# Migration Plan: NetBox Core to NetBox Go

This document outlines the strategy and task list for migrating the `core` module from the Python-based NetBox to the Go-based `netbox_go` project.

## 🎯 Objective
Port the functionality of `netbox/core` to `netbox_go/core` (and corresponding internal layers), maintaining the business logic while adapting to Go's architecture (Domain-Driven Design).

## 🏗️ Target Architecture in `netbox_go`
The migration will follow the existing project structure:
- **Domain Layer**: `netbox_go/internal/domain/core` (Entities, Interfaces, Enums)
- **Infrastructure Layer**: `netbox_go/internal/infrastructure` (Postgres/sqlc, Mocks)
- **Delivery Layer**: `netbox_go/internal/delivery/http` (Handlers, Middleware)

---

## 📋 Migration Task List

### 1. Database & Schema (Migrations)
- [x] Analyze `netbox/core/migrations/` and map to `netbox_go/migrations/`
- [x] Verify existing SQL schemas for Core (DataSources, Jobs, ObjectTypes, ConfigRevisions)
- [x] Implement missing migrations in `netbox_go/migrations/`

### 2. API Layer (`netbox/core/api` - `netbox_go/internal/delivery/http`)
- [x] **Serializers**: Map `api/serializers/` to Go DTOs (Request/Response structs)
- [x] **Views/URLs**: Port `api/views.py` and `api/urls.py` to `core_handler.go`
    - [x] Data Sources API
    - [x] Jobs API
    - [x] Object Types API
    - [x] Config Revisions API
    - [x] Change Logging API
- [ ] **Validation**: Port logic from `api/serializers/` to Go validation logic

### 3. Forms & UI Logic (`netbox/core/forms` & `netbox/core/ui`)
- [ ] **Forms**: Convert `forms/` logic (bulk edit, import) into Service layer logic in Go
- [ ] **UI Panels**: Analyze `ui/panels.py` to ensure API responses provide necessary data for frontend panels

### 4. GraphQL API (`netbox/core/graphql` - `netbox_go/internal/delivery/graphql`)
- [ ] Define GraphQL Types based on `graphql/types.py`
- [ ] Implement GraphQL Resolvers based on `graphql/schema.py`
- [ ] Port Filters and Mixins from `graphql/filters.py` and `graphql/filter_mixins.py`

### 5. Management Commands (`netbox/core/management` - `netbox_go/cmd/...`)
- [x] Port `syncdatasource.py` to a Go CLI command or background worker
- [x] Port `rqworker.py` logic to Go worker implementation
- [x] Implement `nbshell` equivalent if needed for debugging

### 6. Tables & Presentation (`netbox/core/tables` - API Responses)
- [ ] Map `tables/` column definitions to API response fields to ensure compatibility with the UI

### 7. Core Logic & Utils (`netbox/core/` root)
- [ ] Port `search.py` logic to Postgres FTS (tsvector/tsquery) in Repositories
- [ ] Port `filtersets.py` to type-safe Specification/Query Builder pattern
- [ ] Port `jobs.py` and `events.py` to native Go worker pool with DB-backed queue

---

## 🛠️ Technical Mapping Table

| Python Path (`netbox/core`) | Go Path (`netbox_go`) | Note |
| :--- | :--- | :--- |
| `api/` | `internal/delivery/http/handlers/` | REST API implementation |
| `forms/` | `internal/domain/core/services/` | Business logic for data entry |
| `graphql/` | `internal/delivery/graphql/` | GraphQL API implementation |
| `management/` | `cmd/` or `internal/worker/` | CLI and Background tasks |
| `migrations/` | `migrations/` | SQL schema files |
| `tables/` | `internal/delivery/http/dto/` | Data structure for tables |
| `ui/` | `internal/delivery/http/dto/` | UI-specific data requirements |
| `models/` | `internal/domain/core/entity/` | Core domain entities |
