-- Миграция 001: Создание таблиц DCIM (Sites, Racks, Devices, Components, Cables, Power)
-- Версия: 1.0.0
-- Дата: 2024-01-15

-- ============================================
-- Таблицы для Sites
-- ============================================

CREATE TABLE dcim_regions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES dcim_regions(id) ON DELETE SET NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_regions_slug ON dcim_regions(slug);
CREATE INDEX idx_dcim_regions_parent_id ON dcim_regions(parent_id);

CREATE TABLE dcim_site_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES dcim_site_groups(id) ON DELETE SET NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_site_groups_slug ON dcim_site_groups(slug);

CREATE TYPE dcim_site_status AS ENUM ('planned', 'staging', 'active', 'retired');

CREATE TABLE dcim_sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    status dcim_site_status DEFAULT 'active',
    region_id UUID REFERENCES dcim_regions(id) ON DELETE SET NULL,
    group_id UUID REFERENCES dcim_site_groups(id) ON DELETE SET NULL,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    facility VARCHAR(50),
    asn_ids UUID[],
    time_zone VARCHAR(50),
    physical_address TEXT,
    shipping_address TEXT,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_sites_slug ON dcim_sites(slug);
CREATE INDEX idx_dcim_sites_status ON dcim_sites(status);
CREATE INDEX idx_dcim_sites_region_id ON dcim_sites(region_id);
CREATE INDEX idx_dcim_sites_tenant_id ON dcim_sites(tenant_id);

CREATE TYPE dcim_location_status AS ENUM ('planned', 'staging', 'active', 'retired');

CREATE TABLE dcim_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    site_id UUID NOT NULL REFERENCES dcim_sites(id) ON DELETE CASCADE,
    status dcim_location_status DEFAULT 'active',
    parent_id UUID REFERENCES dcim_locations(id) ON DELETE SET NULL,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    facility VARCHAR(50),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_locations_slug ON dcim_locations(site_id, slug);
CREATE INDEX idx_dcim_locations_site_id ON dcim_locations(site_id);

-- ============================================
-- Таблицы для Racks
-- ============================================

CREATE TYPE dcim_rack_type_form_factor AS ENUM ('wallframe', 'wallcabinet', '4postframe', '4postcabinet', '2postframe', '2postcabinet');
CREATE TYPE dcim_rack_dimension_unit AS ENUM ('mm', 'in');

CREATE TABLE dcim_rack_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manufacturer_id UUID NOT NULL REFERENCES dcim_manufacturers(id) ON DELETE CASCADE,
    model VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    form_factor dcim_rack_type_form_factor NOT NULL,
    width SMALLINT NOT NULL CHECK (width IN (19, 23)),
    u_height SMALLINT NOT NULL CHECK (u_height > 0 AND u_height <= 1000),
    starting_unit SMALLINT DEFAULT 1,
    desc_units BOOLEAN DEFAULT FALSE,
    outer_width SMALLINT,
    outer_height SMALLINT,
    outer_depth SMALLINT,
    outer_unit dcim_rack_dimension_unit,
    mounting_depth SMALLINT,
    weight INTEGER,
    max_weight INTEGER,
    weight_unit VARCHAR(10),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TYPE dcim_rack_status AS ENUM ('reserved', 'available', 'planned', 'active', 'deprecated');

CREATE TABLE dcim_racks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    facility_id VARCHAR(50),
    site_id UUID NOT NULL REFERENCES dcim_sites(id) ON DELETE CASCADE,
    location_id UUID REFERENCES dcim_locations(id) ON DELETE SET NULL,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    status dcim_rack_status DEFAULT 'active',
    role_id UUID REFERENCES dcim_rack_roles(id) ON DELETE SET NULL,
    rack_type_id UUID REFERENCES dcim_rack_types(id) ON DELETE SET NULL,
    form_factor dcim_rack_type_form_factor,
    width SMALLINT NOT NULL DEFAULT 19 CHECK (width IN (19, 23)),
    serial VARCHAR(50),
    asset_tag VARCHAR(50),
    airflow VARCHAR(20),
    u_height SMALLINT NOT NULL DEFAULT 42 CHECK (u_height > 0 AND u_height <= 1000),
    starting_unit SMALLINT DEFAULT 1,
    desc_units BOOLEAN DEFAULT FALSE,
    outer_width SMALLINT,
    outer_height SMALLINT,
    outer_depth SMALLINT,
    outer_unit dcim_rack_dimension_unit,
    mounting_depth SMALLINT,
    weight INTEGER,
    max_weight INTEGER,
    weight_unit VARCHAR(10),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_racks_site_id ON dcim_racks(site_id);
CREATE INDEX idx_dcim_racks_location_id ON dcim_racks(location_id);
CREATE INDEX idx_dcim_racks_status ON dcim_racks(status);

CREATE TABLE dcim_rack_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rack_id UUID NOT NULL REFERENCES dcim_racks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    units SMALLINT[] NOT NULL,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- Таблицы для Cables
-- ============================================

CREATE TYPE dcim_cable_type AS ENUM (
    'cat3', 'cat5', 'cat5e', 'cat6', 'cat6a', 'cat7', 'cat8',
    'dac', 'mrj21', 'mtp', 'multimode-fiber', 'multimode-fiber-om1',
    'multimode-fiber-om2', 'multimode-fiber-om3', 'multimode-fiber-om4',
    'multimode-fiber-om5', 'singlemode-fiber', 'singlemode-fiber-os1',
    'singlemode-fiber-os2', 'coaxial', 'power'
);

CREATE TYPE dcim_cable_status AS ENUM ('connected', 'planned', 'decommissioning', 'failed');

CREATE TABLE dcim_cables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type dcim_cable_type NOT NULL,
    status dcim_cable_status DEFAULT 'connected',
    label VARCHAR(50),
    color VARCHAR(20),
    length INTEGER,
    length_unit VARCHAR(10) DEFAULT 'm',
    description TEXT,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_cables_status ON dcim_cables(status);
CREATE INDEX idx_dcim_cables_type ON dcim_cables(type);

CREATE TABLE dcim_cable_terminations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cable_id UUID NOT NULL REFERENCES dcim_cables(id) ON DELETE CASCADE,
    termination_type VARCHAR(50) NOT NULL,
    termination_id UUID NOT NULL,
    cable_end VARCHAR(1) NOT NULL CHECK (cable_end IN ('A', 'B')),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cable_id, termination_type, termination_id)
);

CREATE INDEX idx_dcim_cable_terminations_cable_id ON dcim_cable_terminations(cable_id);
CREATE INDEX idx_dcim_cable_terminations_termination ON dcim_cable_terminations(termination_type, termination_id);

-- ============================================
-- Таблицы для Power
-- ============================================

CREATE TYPE dcim_power_feed_status AS ENUM ('offline', 'active', 'planned', 'failed');
CREATE TYPE dcim_power_feed_type AS ENUM ('primary', 'redundant');
CREATE TYPE dcim_power_supply AS ENUM ('ac', 'dc');
CREATE TYPE dcim_phase_type AS ENUM ('1-phase', '3-phase');
CREATE TYPE dcim_power_unit AS ENUM ('w', 'kw');

CREATE TABLE dcim_power_panels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID NOT NULL REFERENCES dcim_sites(id) ON DELETE CASCADE,
    location_id UUID REFERENCES dcim_locations(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_power_panels_site_id ON dcim_power_panels(site_id);

CREATE TABLE dcim_power_feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    power_panel_id UUID NOT NULL REFERENCES dcim_power_panels(id) ON DELETE CASCADE,
    rack_id UUID REFERENCES dcim_racks(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    status dcim_power_feed_status DEFAULT 'active',
    type dcim_power_feed_type DEFAULT 'primary',
    supply dcim_power_supply DEFAULT 'ac',
    phase dcim_phase_type DEFAULT '1-phase',
    voltage INTEGER NOT NULL CHECK (voltage > 0),
    amperage INTEGER NOT NULL CHECK (amperage > 0),
    max_utilization INTEGER DEFAULT 80 CHECK (max_utilization >= 0 AND max_utilization <= 100),
    available_power INTEGER,
    unit dcim_power_unit DEFAULT 'w',
    description TEXT,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    cable_end VARCHAR(1) CHECK (cable_end IN ('A', 'B')),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_power_feeds_power_panel_id ON dcim_power_feeds(power_panel_id);
CREATE INDEX idx_dcim_power_feeds_rack_id ON dcim_power_feeds(rack_id);
