-- Миграция 006: Создание таблиц DCIM для устройств (Devices, Manufacturers, DeviceTypes, etc.)
-- Версия: 1.0.0
-- Дата: 2024-01-15

-- ============================================
-- Таблицы для Manufacturers, DeviceTypes, Platforms, DeviceRoles
-- ============================================

CREATE TABLE dcim_manufacturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_manufacturers_slug ON dcim_manufacturers(slug);

CREATE TABLE dcim_device_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manufacturer_id UUID NOT NULL REFERENCES dcim_manufacturers(id) ON DELETE CASCADE,
    model VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    part_number VARCHAR(100),
    default_platform_id UUID REFERENCES dcim_platforms(id) ON DELETE SET NULL,
    u_height NUMERIC(4,1) NOT NULL DEFAULT 1 CHECK (u_height >= 0),
    full_depth BOOLEAN DEFAULT FALSE,
    subdevice_role VARCHAR(20) CHECK (subdevice_role IN ('parent', 'child', 'true')),
    airflow VARCHAR(20) CHECK (airflow IN ('front-to-rear', 'rear-to-front', 'left-to-right', 'right-to-left', 'both', 'front', 'rear')),
    front_image VARCHAR(100),
    rear_image VARCHAR(100),
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_device_types_slug ON dcim_device_types(slug);
CREATE INDEX idx_dcim_device_types_manufacturer_id ON dcim_device_types(manufacturer_id);

CREATE TABLE dcim_platforms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    manufacturer_id UUID REFERENCES dcim_manufacturers(id) ON DELETE SET NULL,
    napalm_driver VARCHAR(100),
    napalm_args TEXT,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_platforms_slug ON dcim_platforms(slug);
CREATE INDEX idx_dcim_platforms_manufacturer_id ON dcim_platforms(manufacturer_id);

CREATE TABLE dcim_device_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(6) NOT NULL DEFAULT '9e9e9e',
    vm_role BOOLEAN DEFAULT FALSE,
    config_template_id UUID,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_device_roles_slug ON dcim_device_roles(slug);

CREATE TABLE dcim_config_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    data_source VARCHAR(100),
    data_path VARCHAR(255),
    environment JSONB,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- ============================================
-- Таблица для устройств
-- ============================================

CREATE TYPE dcim_device_status AS ENUM ('planned', 'staged', 'preprovisioned', 'offline', 'active', 'failed', 'inventory', 'decommissioning');

CREATE TABLE dcim_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE RESTRICT,
    device_role_id UUID NOT NULL REFERENCES dcim_device_roles(id) ON DELETE RESTRICT,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    platform_id UUID REFERENCES dcim_platforms(id) ON DELETE SET NULL,
    site_id UUID NOT NULL REFERENCES dcim_sites(id) ON DELETE CASCADE,
    location_id UUID REFERENCES dcim_locations(id) ON DELETE SET NULL,
    rack_id UUID REFERENCES dcim_racks(id) ON DELETE SET NULL,
    position NUMERIC(4,1) CHECK (position > 0 AND position <= 100),
    face VARCHAR(20) CHECK (face IN ('front', 'rear', 'inline')),
    status dcim_device_status DEFAULT 'active',
    airflow VARCHAR(20),
    serial VARCHAR(50),
    asset_tag VARCHAR(50) UNIQUE,
    virtual_chassis_id UUID REFERENCES dcim_virtual_chassis(id) ON DELETE SET NULL,
    vc_position SMALLINT CHECK (vc_position >= 0 AND vc_position <= 255),
    vc_priority SMALLINT CHECK (vc_priority >= 0 AND vc_priority <= 255),
    primary_ip4_id UUID REFERENCES ipam_ip_addresses(id) ON DELETE SET NULL,
    primary_ip6_id UUID REFERENCES ipam_ip_addresses(id) ON DELETE SET NULL,
    cluster_id UUID REFERENCES virtualization_clusters(id) ON DELETE SET NULL,
    config_template_id UUID REFERENCES dcim_config_templates(id) ON DELETE SET NULL,
    local_context_data JSONB,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_devices_name ON dcim_devices(name);
CREATE INDEX idx_dcim_devices_site_id ON dcim_devices(site_id);
CREATE INDEX idx_dcim_devices_rack_id ON dcim_devices(rack_id);
CREATE INDEX idx_dcim_devices_device_type_id ON dcim_devices(device_type_id);
CREATE INDEX idx_dcim_devices_device_role_id ON dcim_devices(device_role_id);
CREATE INDEX idx_dcim_devices_status ON dcim_devices(status);
CREATE INDEX idx_dcim_devices_tenant_id ON dcim_devices(tenant_id);
CREATE UNIQUE INDEX idx_dcim_devices_rack_position ON dcim_devices(rack_id, position, face) WHERE position IS NOT NULL;

-- ============================================
-- Таблицы для виртуальных шасси
-- ============================================

CREATE TABLE dcim_virtual_chassis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    site_id UUID REFERENCES dcim_sites(id) ON DELETE SET NULL,
    domain VARCHAR(100),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_virtual_chassis_slug ON dcim_virtual_chassis(slug);

-- ============================================
-- Таблицы для кластеров виртуализации
-- ============================================

CREATE TABLE dcim_cluster_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_cluster_types_slug ON dcim_cluster_types(slug);

CREATE TABLE dcim_cluster_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    site_id UUID REFERENCES dcim_sites(id) ON DELETE SET NULL,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_cluster_groups_slug ON dcim_cluster_groups(slug);

CREATE TABLE dcim_clusters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    cluster_type_id UUID NOT NULL REFERENCES dcim_cluster_types(id) ON DELETE RESTRICT,
    cluster_group_id UUID REFERENCES dcim_cluster_groups(id) ON DELETE SET NULL,
    site_id UUID REFERENCES dcim_sites(id) ON DELETE SET NULL,
    tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_dcim_clusters_name_group ON dcim_clusters(name, cluster_group_id) WHERE cluster_group_id IS NOT NULL;
CREATE INDEX idx_dcim_clusters_site_id ON dcim_clusters(site_id);

-- ============================================
-- Таблицы для модулей
-- ============================================

CREATE TABLE dcim_module_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manufacturer_id UUID NOT NULL REFERENCES dcim_manufacturers(id) ON DELETE CASCADE,
    model VARCHAR(100) NOT NULL,
    part_number VARCHAR(100),
    weight NUMERIC(10,3),
    weight_unit VARCHAR(20) DEFAULT 'kg',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_module_types_manufacturer_id ON dcim_module_types(manufacturer_id);

CREATE TABLE dcim_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_type_id UUID NOT NULL REFERENCES dcim_module_types(id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    module_bay_id UUID REFERENCES dcim_module_bays(id) ON DELETE SET NULL,
    serial VARCHAR(50),
    asset_tag VARCHAR(50) UNIQUE,
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_modules_device_id ON dcim_modules(device_id);

CREATE TABLE dcim_module_bay_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    position VARCHAR(20),
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_dcim_module_bay_templates_device_type_id ON dcim_module_bay_templates(device_type_id);

CREATE TABLE dcim_module_bays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    module_bay_template_id UUID REFERENCES dcim_module_bay_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    position VARCHAR(20),
    description TEXT,
    installed_module_id UUID REFERENCES dcim_modules(id) ON DELETE SET NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_module_bays_device_id ON dcim_module_bays(device_id);
CREATE UNIQUE INDEX idx_dcim_module_bays_name_device ON dcim_module_bays(device_id, name);

-- ============================================
-- Таблицы для компонентов устройств (Interfaces, Console Ports, Power Ports, etc.)
-- ============================================

-- Interface Templates
CREATE TABLE dcim_interface_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL,
    mgmt_only BOOLEAN DEFAULT FALSE,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_dcim_interface_templates_device_type_id ON dcim_interface_templates(device_type_id);

-- Interfaces
CREATE TABLE dcim_interfaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    interface_template_id UUID REFERENCES dcim_interface_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    enabled BOOLEAN DEFAULT TRUE,
    mgmt_only BOOLEAN DEFAULT FALSE,
    mtu INTEGER,
    speed INTEGER,
    duplex VARCHAR(20),
    port_security_enabled BOOLEAN DEFAULT FALSE,
    vn_tenant_id UUID REFERENCES tenancy_tenants(id) ON DELETE SET NULL,
    mode VARCHAR(20),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_interfaces_device_id ON dcim_interfaces(device_id);
CREATE UNIQUE INDEX idx_dcim_interfaces_name_device ON dcim_interfaces(device_id, name);

-- Console Port Templates
CREATE TABLE dcim_console_port_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Console Ports
CREATE TABLE dcim_console_ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    console_port_template_id UUID REFERENCES dcim_console_port_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    speed VARCHAR(20),
    mark_connected BOOLEAN DEFAULT FALSE,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_console_ports_device_id ON dcim_console_ports(device_id);

-- Console Server Port Templates
CREATE TABLE dcim_console_server_port_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Console Server Ports
CREATE TABLE dcim_console_server_ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    console_server_port_template_id UUID REFERENCES dcim_console_server_port_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    speed VARCHAR(20),
    mark_connected BOOLEAN DEFAULT FALSE,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_console_server_ports_device_id ON dcim_console_server_ports(device_id);

-- Power Port Templates
CREATE TABLE dcim_power_port_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    maximum_draw INTEGER,
    allocated_draw INTEGER,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Power Ports
CREATE TABLE dcim_power_ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    power_port_template_id UUID REFERENCES dcim_power_port_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    maximum_draw INTEGER,
    allocated_draw INTEGER,
    mark_connected BOOLEAN DEFAULT FALSE,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_power_ports_device_id ON dcim_power_ports(device_id);

-- Power Outlet Templates
CREATE TABLE dcim_power_outlet_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    power_port_id UUID REFERENCES dcim_power_port_templates(id) ON DELETE SET NULL,
    feed_leg VARCHAR(10),
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Power Outlets
CREATE TABLE dcim_power_outlets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    power_outlet_template_id UUID REFERENCES dcim_power_outlet_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL DEFAULT 'other',
    power_port_id UUID REFERENCES dcim_power_ports(id) ON DELETE SET NULL,
    feed_leg VARCHAR(10),
    mark_connected BOOLEAN DEFAULT FALSE,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_power_outlets_device_id ON dcim_power_outlets(device_id);

-- Front Port Templates
CREATE TABLE dcim_front_port_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL,
    color VARCHAR(6),
    rear_port_id UUID NOT NULL REFERENCES dcim_rear_port_templates(id) ON DELETE CASCADE,
    rear_port_position INTEGER NOT NULL DEFAULT 1,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Front Ports
CREATE TABLE dcim_front_ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    front_port_template_id UUID REFERENCES dcim_front_port_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL,
    color VARCHAR(6),
    rear_port_id UUID NOT NULL REFERENCES dcim_rear_ports(id) ON DELETE CASCADE,
    rear_port_position INTEGER NOT NULL DEFAULT 1,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_front_ports_device_id ON dcim_front_ports(device_id);

-- Rear Port Templates
CREATE TABLE dcim_rear_port_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL,
    positions INTEGER NOT NULL DEFAULT 1,
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Rear Ports
CREATE TABLE dcim_rear_ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    rear_port_template_id UUID REFERENCES dcim_rear_port_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    type VARCHAR(50) NOT NULL,
    positions INTEGER NOT NULL DEFAULT 1,
    cable_id UUID REFERENCES dcim_cables(id) ON DELETE SET NULL,
    connection_status VARCHAR(20) DEFAULT 'connected',
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_rear_ports_device_id ON dcim_rear_ports(device_id);

-- Device Bay Templates
CREATE TABLE dcim_device_bay_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    description TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Device Bays
CREATE TABLE dcim_device_bays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    device_bay_template_id UUID REFERENCES dcim_device_bay_templates(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    installed_device_id UUID REFERENCES dcim_devices(id) ON DELETE SET NULL,
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_device_bays_device_id ON dcim_device_bays(device_id);

-- Inventory Item Templates
CREATE TABLE dcim_inventory_item_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type_id UUID NOT NULL REFERENCES dcim_device_types(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES dcim_inventory_item_templates(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    manufacturer_id UUID REFERENCES dcim_manufacturers(id) ON DELETE SET NULL,
    part_id VARCHAR(100),
    description TEXT,
    component_id UUID,
    component_type VARCHAR(100),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_dcim_inventory_item_templates_device_type_id ON dcim_inventory_item_templates(device_type_id);

-- Inventory Items
CREATE TABLE dcim_inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dcim_devices(id) ON DELETE CASCADE,
    inventory_item_template_id UUID REFERENCES dcim_inventory_item_templates(id) ON DELETE SET NULL,
    parent_id UUID REFERENCES dcim_inventory_items(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(100),
    manufacturer_id UUID REFERENCES dcim_manufacturers(id) ON DELETE SET NULL,
    part_id VARCHAR(100),
    serial VARCHAR(50),
    asset_tag VARCHAR(50),
    description TEXT,
    comments TEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_dcim_inventory_items_device_id ON dcim_inventory_items(device_id);
CREATE UNIQUE INDEX idx_dcim_inventory_items_serial ON dcim_inventory_items(serial) WHERE serial IS NOT NULL AND serial != '';
