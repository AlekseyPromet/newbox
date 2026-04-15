# NetBox Go - Implementation Plan

## Overview

This document describes the implementation plan for the NetBox Go application - a complete rewrite of NetBox in Go with full support for all NetBox domains.

## Architecture

The application follows **Hexagonal (Ports & Adapters) / Hive architecture** pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│  ┌──────────────────────┐  ┌────────────────────────────┐   │
│  │   Echo REST API      │  │   gqlgen GraphQL API       │   │
│  └──────────────────────┘  └────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Use Cases / Services                     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌───────────────┐    │
│  │  DCIM   │ │  IPAM   │ │ Circuits │ │ Virtualization│    │
│  └─────────┘ └─────────┘ └──────────┘ └───────────────┘    │
│  ┌─────────┐ ┌─────────┐                                   │
│  │ Tenancy │ │  VPN    │                                   │
│  └─────────┘ └─────────┘                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                         │
│  ┌──────────────────────┐  ┌────────────────────────────┐   │
│  │  PostgreSQL (sqlc)   │  │   Etcd (Cache/Locks)       │   │
│  └──────────────────────┘  └────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Implemented Domains

### 1. DCIM (Data Center Infrastructure Management) ✅

**Location:** `internal/domain/dcim/`

#### Entities:
- **Sites**: Geographic locations with regions, groups, status tracking
- **Racks**: Physical racks with types, roles, reservations, unit positions
- **Devices**: Physical and virtual devices with types, roles, platforms
- **Components**: 
  - Console ports/ports
  - Power ports/feeds  
  - Network interfaces (all types from 100M to 800G, fiber, wireless)
  - Module bays and modules
- **Cables**: Physical cable connections with termination tracking
- **Power**: Power panels and feeds with electrical parameters

#### Features:
- Full status lifecycle management
- Hierarchical organization (regions, groups, locations)
- Rack elevation and unit tracking
- Cable tracing and connectivity
- Power capacity calculation
- Airflow direction tracking
- Serial number and asset tag management

### 2. IPAM (IP Address Management) ✅

**Location:** `internal/domain/ipam/`

#### Entities:
- **VRFs**: Virtual Routing and Forwarding instances with route targets
- **Prefixes**: IPv4/IPv6 subnets with hierarchical nesting
- **IP Addresses**: Individual addresses with assignment tracking
- **VLANs**: Virtual LANs with groups and Q-in-Q support
- **ASNs**: Autonomous System Numbers
- **Services**: L4-L7 services (TCP/UDP/SCTP) on devices/VMs
- **FHRP Groups**: First Hop Redundancy Protocol groups (VRRP, HSRP, etc.)
- **RIRs**: Regional Internet Registries
- **Aggregates**: Top-level IP aggregates
- **Roles**: Functional roles for prefixes, VLANs, IPs

#### Features:
- **Full IPv4 and IPv6 support** using `net/netip` package
- Prefix utilization calculation
- IP address assignment to interfaces
- NAT inside/outside relationships
- DNS name tracking
- Status-based color coding
- Route target import/export for VRFs

### 3. Circuits ✅

**Location:** `internal/domain/circuits/`

#### Entities:
- **Providers**: Telecommunications providers with ASN, contacts
- **Provider Networks**: Provider's network infrastructure
- **Circuit Types**: Categorization of circuits
- **Circuits**: Telecom circuits with installation dates, commit rates
- **Circuit Terminations**: A/Z side terminations with cross-connect info

#### Features:
- Circuit status lifecycle (planned → provisioning → active → decommissioning)
- Port speed and upstream speed tracking
- Cross-connect ID and patch panel documentation
- Provider portal URL integration

### 4. Virtualization ✅

**Location:** `internal/domain/virtualization/`

#### Entities:
- **Cluster Types**: Hypervisor types (VMware, KVM, Hyper-V, etc.)
- **Cluster Groups**: Logical groupings of clusters
- **Clusters**: Virtualization clusters with sites and status
- **Virtual Machines**: VMs with vCPU, memory, disk allocation
- **VM Interfaces**: Virtual network interfaces with VLAN support
- **VM Disks**: Virtual disk definitions

#### Features:
- VM status management (offline, active, planned, staged, failed)
- Resource allocation tracking (vCPU, RAM, disk)
- Primary IP assignment (IPv4/IPv6)
- Interface VLAN configuration (access/tagged modes)
- Cluster device membership

### 5. Tenancy ✅

**Location:** `internal/domain/tenancy/`

#### Entities:
- **Tenant Groups**: Hierarchical tenant organization
- **Tenants**: Organizations/customers
- **Contact Groups**: Contact hierarchy
- **Contact Roles**: Contact function roles
- **Contacts**: Individual contact information
- **Contact Assignments**: Linking contacts to objects

#### Features:
- Hierarchical tenant groups
- Multi-level contact assignments
- Priority-based contact ordering
- Full contact information (phone, email, address, links)

## Repository Layer

### PostgreSQL (sqlc)

**Location:** `internal/repository/postgres/`

- Type-safe SQL queries generated by sqlc
- Connection pooling and transaction management
- Migration support
- Prepared statements for performance

### Etcd (Cache/Distributed Locks)

**Location:** `internal/repository/etcd/`

- Distributed caching for frequently accessed data
- Distributed locks for concurrent operations
- Lease management for cache expiration
- Watch mechanisms for cache invalidation

## Delivery Layer

### REST API (Echo)

**Location:** `internal/delivery/http/`

- Built with Echo framework
- OpenAPI 3.0 specification
- JWT authentication
- Rate limiting
- Request/response logging
- Pagination and filtering support
- CRUD operations for all entities

### GraphQL API (gqlgen)

**Location:** `internal/delivery/graphql/`

- Built with gqlgen
- Schema-first development
- Resolvers for all domains
- Subscription support for real-time updates
- DataLoader for N+1 query optimization

## SOLID Principles Applied

### Single Responsibility Principle (SRP)
- Each entity has its own repository interface
- Separate handlers for each API endpoint
- Domain logic isolated from infrastructure

### Open/Closed Principle (OCP)
- Repository interfaces allow adding new implementations
- Handler middleware chain for extensibility
- Enum validation through interface

### Liskov Substitution Principle (LSP)
- All repository implementations satisfy the same interface
- Entity validate methods follow consistent contract

### Interface Segregation Principle (ISP)
- Granular repository interfaces per aggregate root
- Specific filter structs per entity type

### Dependency Inversion Principle (DIP)
- Domain layer has no dependencies on infrastructure
- Dependency injection via constructors
- Interface-based abstractions

## Project Structure

```
netbox_go/
├── cmd/
│   ├── api/              # REST API entry point
│   └── graphql/          # GraphQL API entry point
├── config/               # Configuration files
├── docs/                 # Documentation
├── internal/
│   ├── domain/
│   │   ├── dcim/         # DCIM domain
│   │   ├── ipam/         # IPAM domain
│   │   ├── circuits/     # Circuits domain
│   │   ├── virtualization/ # Virtualization domain
│   │   └── tenancy/      # Tenancy domain
│   ├── repository/
│   │   ├── postgres/     # PostgreSQL implementation
│   │   └── etcd/         # Etcd implementation
│   ├── delivery/
│   │   ├── http/         # REST API handlers
│   │   └── graphql/      # GraphQL resolvers
│   └── app/              # Application services
├── pkg/
│   └── types/            # Shared types and errors
└── migrations/           # Database migrations
```

## Next Steps

1. **Repository Implementations**
   - Implement sqlc queries for all entities
   - Add Etcd caching layer
   - Implement distributed locks

2. **Application Services**
   - Create use cases for business logic
   - Implement validation pipelines
   - Add audit logging

3. **Delivery Layer**
   - Implement REST API handlers with Echo
   - Create GraphQL schema and resolvers
   - Add authentication/authorization

4. **Testing**
   - Unit tests for domain entities
   - Integration tests for repositories
   - E2E tests for APIs

5. **Documentation**
   - API documentation (OpenAPI/Swagger)
   - GraphQL schema documentation
   - Deployment guides

## Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Cache/Distributed Locks**: Etcd 3.5+
- **REST Framework**: Echo v4
- **GraphQL**: gqlgen v0.17+
- **SQL Generator**: sqlc v1.20+
- **Migrations**: golang-migrate
- **Testing**: testify, testcontainers-go
