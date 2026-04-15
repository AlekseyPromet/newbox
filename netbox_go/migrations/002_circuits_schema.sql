-- Миграция 002: Создание таблиц Circuits (Providers, Circuits, Terminations, Virtual Circuits)
-- Версия: 1.0.0
-- Дата: 2024-01-15

-- ============================================
-- Таблицы для Providers
-- ============================================

CREATE TABLE circuits_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    asn_ids UUID[],
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_providers_slug ON circuits_providers(slug);
CREATE INDEX idx_circuits_providers_name ON circuits_providers(name);

CREATE TABLE circuits_provider_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES circuits_providers(id) ON DELETE PROTECT,
    account VARCHAR(100) NOT NULL,
    name VARCHAR(100),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_provider_accounts_provider_id ON circuits_provider_accounts(provider_id);
CREATE UNIQUE INDEX idx_circuits_provider_accounts_unique ON circuits_provider_accounts(provider_id, account);
CREATE UNIQUE INDEX idx_circuits_provider_accounts_unique_name ON circuits_provider_accounts(provider_id, name) WHERE name != '';

CREATE TABLE circuits_provider_networks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES circuits_providers(id) ON DELETE PROTECT,
    name VARCHAR(100) NOT NULL,
    service_id VARCHAR(100),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_provider_networks_provider_id ON circuits_provider_networks(provider_id);
CREATE UNIQUE INDEX idx_circuits_provider_networks_unique ON circuits_provider_networks(provider_id, name);

-- ============================================
-- Таблицы для Circuit Types и Circuits
-- ============================================

CREATE TYPE circuits_circuit_status AS ENUM ('planned', 'provisioning', 'active', 'offline', 'deprovisioning', 'decommissioned');

CREATE TABLE circuits_circuit_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(6),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_circuit_types_slug ON circuits_circuit_types(slug);

CREATE TABLE circuits_circuits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cid VARCHAR(100) NOT NULL,
    provider_id UUID NOT NULL REFERENCES circuits_providers(id) ON DELETE PROTECT,
    provider_account_id UUID REFERENCES circuits_provider_accounts(id) ON DELETE PROTECT,
    type_id UUID NOT NULL REFERENCES circuits_circuit_types(id) ON DELETE PROTECT,
    status circuits_circuit_status DEFAULT 'active',
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE PROTECT,
    install_date DATE,
    termination_date DATE,
    commit_rate INTEGER CHECK (commit_rate > 0),
    distance NUMERIC(10,2) CHECK (distance >= 0),
    distance_unit VARCHAR(10) DEFAULT 'km',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    termination_a_id UUID,
    termination_z_id UUID
);

CREATE INDEX idx_circuits_circuits_provider_id ON circuits_circuits(provider_id);
CREATE INDEX idx_circuits_circuits_provider_account_id ON circuits_circuits(provider_account_id);
CREATE INDEX idx_circuits_circuits_type_id ON circuits_circuits(type_id);
CREATE INDEX idx_circuits_circuits_status ON circuits_circuits(status);
CREATE INDEX idx_circuits_circuits_tenant_id ON circuits_circuits(tenant_id);
CREATE UNIQUE INDEX idx_circuits_circuits_unique_provider_cid ON circuits_circuits(provider_id, cid);
CREATE UNIQUE INDEX idx_circuits_circuits_unique_provideraccount_cid ON circuits_circuits(provider_account_id, cid) WHERE provider_account_id IS NOT NULL;

-- ============================================
-- Таблицы для Circuit Terminations
-- ============================================

CREATE TYPE circuits_termination_side AS ENUM ('A', 'Z');

CREATE TABLE circuits_circuit_terminations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    circuit_id UUID NOT NULL REFERENCES circuits_circuits(id) ON DELETE CASCADE,
    term_side circuits_termination_side NOT NULL,
    termination_type VARCHAR(50),
    termination_id UUID,
    port_speed INTEGER CHECK (port_speed > 0),
    upstream_speed INTEGER CHECK (upstream_speed > 0),
    xconnect_id VARCHAR(50),
    pp_info VARCHAR(100),
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    -- Cached associations for filtering
    _provider_network_id UUID REFERENCES circuits_provider_networks(id) ON DELETE PROTECT,
    _region_id UUID REFERENCES dcim_regions(id) ON DELETE CASCADE,
    _site_group_id UUID REFERENCES dcim_site_groups(id) ON DELETE CASCADE,
    _site_id UUID REFERENCES dcim_sites(id) ON DELETE CASCADE,
    _location_id UUID REFERENCES dcim_locations(id) ON DELETE CASCADE
);

CREATE INDEX idx_circuits_circuit_terminations_circuit_id ON circuits_circuit_terminations(circuit_id);
CREATE INDEX idx_circuits_circuit_terminations_term_side ON circuits_circuit_terminations(term_side);
CREATE INDEX idx_circuits_circuit_terminations_termination ON circuits_circuit_terminations(termination_type, termination_id);
CREATE INDEX idx_circuits_circuit_terminations_cached_site ON circuits_circuit_terminations(_site_id);
CREATE INDEX idx_circuits_circuit_terminations_cached_region ON circuits_circuit_terminations(_region_id);
CREATE INDEX idx_circuits_circuit_terminations_cached_location ON circuits_circuit_terminations(_location_id);
CREATE INDEX idx_circuits_circuit_terminations_cached_provider_network ON circuits_circuit_terminations(_provider_network_id);
CREATE UNIQUE INDEX idx_circuits_circuit_terminations_unique_circuit_term_side ON circuits_circuit_terminations(circuit_id, term_side);

-- ============================================
-- Таблицы для Circuit Groups
-- ============================================

CREATE TABLE circuits_circuit_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE PROTECT,
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_circuit_groups_slug ON circuits_circuit_groups(slug);
CREATE INDEX idx_circuits_circuit_groups_tenant_id ON circuits_circuit_groups(tenant_id);

CREATE TYPE circuits_circuit_priority AS ENUM ('primary', 'secondary', 'tertiary', 'inactive');

CREATE TABLE circuits_circuit_group_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_type VARCHAR(50) NOT NULL,
    member_id UUID NOT NULL,
    group_id UUID NOT NULL REFERENCES circuits_circuit_groups(id) ON DELETE CASCADE,
    priority circuits_circuit_priority,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_circuits_circuit_group_assignments_member ON circuits_circuit_group_assignments(member_type, member_id);
CREATE INDEX idx_circuits_circuit_group_assignments_group_id ON circuits_circuit_group_assignments(group_id);
CREATE UNIQUE INDEX idx_circuits_circuit_group_assignments_unique_member_group ON circuits_circuit_group_assignments(member_type, member_id, group_id);

-- ============================================
-- Таблицы для Virtual Circuits
-- ============================================

CREATE TABLE circuits_virtual_circuit_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(6),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_virtual_circuit_types_slug ON circuits_virtual_circuit_types(slug);

CREATE TABLE circuits_virtual_circuits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cid VARCHAR(100) NOT NULL,
    provider_network_id UUID NOT NULL REFERENCES circuits_provider_networks(id) ON DELETE PROTECT,
    provider_account_id UUID REFERENCES circuits_provider_accounts(id) ON DELETE PROTECT,
    type_id UUID NOT NULL REFERENCES circuits_virtual_circuit_types(id) ON DELETE PROTECT,
    status circuits_circuit_status DEFAULT 'active',
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE PROTECT,
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_virtual_circuits_provider_network_id ON circuits_virtual_circuits(provider_network_id);
CREATE INDEX idx_circuits_virtual_circuits_provider_account_id ON circuits_virtual_circuits(provider_account_id);
CREATE INDEX idx_circuits_virtual_circuits_type_id ON circuits_virtual_circuits(type_id);
CREATE INDEX idx_circuits_virtual_circuits_status ON circuits_virtual_circuits(status);
CREATE INDEX idx_circuits_virtual_circuits_tenant_id ON circuits_virtual_circuits(tenant_id);
CREATE UNIQUE INDEX idx_circuits_virtual_circuits_unique_provider_network_cid ON circuits_virtual_circuits(provider_network_id, cid);
CREATE UNIQUE INDEX idx_circuits_virtual_circuits_unique_provideraccount_cid ON circuits_virtual_circuits(provider_account_id, cid) WHERE provider_account_id IS NOT NULL;

CREATE TYPE circuits_virtual_circuit_termination_role AS ENUM ('peer', 'hub', 'spoke');

CREATE TABLE circuits_virtual_circuit_terminations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    virtual_circuit_id UUID NOT NULL REFERENCES circuits_virtual_circuits(id) ON DELETE CASCADE,
    role circuits_virtual_circuit_termination_role DEFAULT 'peer',
    interface_id UUID NOT NULL REFERENCES dcim_interfaces(id) ON DELETE CASCADE,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_circuits_virtual_circuit_terminations_virtual_circuit_id ON circuits_virtual_circuit_terminations(virtual_circuit_id);
CREATE INDEX idx_circuits_virtual_circuit_terminations_role ON circuits_virtual_circuit_terminations(role);
CREATE UNIQUE INDEX idx_circuits_virtual_circuit_terminations_unique_interface ON circuits_virtual_circuit_terminations(interface_id);

-- ============================================
-- Комментарии к таблицам
-- ============================================

COMMENT ON TABLE circuits_providers IS 'Провайдеры телекоммуникационных услуг';
COMMENT ON TABLE circuits_provider_accounts IS 'Аккаунты внутри провайдеров';
COMMENT ON TABLE circuits_provider_networks IS 'Сети провайдеров вне NetBox';
COMMENT ON TABLE circuits_circuit_types IS 'Типы физических цепей';
COMMENT ON TABLE circuits_circuits IS 'Физические телекоммуникационные цепи';
COMMENT ON TABLE circuits_circuit_terminations IS 'Точки завершения физических цепей';
COMMENT ON TABLE circuits_circuit_groups IS 'Административные группы цепей';
COMMENT ON TABLE circuits_circuit_group_assignments IS 'Привязки цепей к группам';
COMMENT ON TABLE circuits_virtual_circuit_types IS 'Типы виртуальных цепей';
COMMENT ON TABLE circuits_virtual_circuits IS 'Виртуальные соединения между эндпоинтами';
COMMENT ON TABLE circuits_virtual_circuit_terminations IS 'Точки завершения виртуальных цепей';

COMMENT ON COLUMN circuits_circuits.commit_rate IS 'Committed rate в Kbps';
COMMENT ON COLUMN circuits_circuits.distance IS 'Длина цепи';
COMMENT ON COLUMN circuits_circuits.distance_unit IS 'Единица измерения длины (km, m, mi, ft)';
COMMENT ON COLUMN circuits_circuit_terminations.port_speed IS 'Физическая скорость порта в Kbps';
COMMENT ON COLUMN circuits_circuit_terminations.upstream_speed IS 'Upstream скорость для асимметричных цепей в Kbps';
COMMENT ON COLUMN circuits_circuit_terminations.xconnect_id IS 'ID локального кросс-коннекта';
COMMENT ON COLUMN circuits_circuit_terminations.pp_info IS 'ID патч-панели и номера портов';
