// Package entity содержит доменные сущности DCIM модуля
package entity

import (
	"time"

	"netbox_go/internal/domain/dcim/enum"
	"netbox_go/pkg/types"
)

// Manufacturer представляет производителя устройств
type Manufacturer struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность производителя
func (m *Manufacturer) Validate() error {
	if m.Name == "" {
		return types.ErrValidationFailed
	}
	if err := m.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// DeviceType представляет тип устройства
type DeviceType struct {
	ID              types.ID         `json:"id"`
	Manufacturer    *Manufacturer    `json:"manufacturer"`
	Model           string           `json:"model"`
	Slug            types.Slug       `json:"slug"`
	DefaultPlatform *Platform        `json:"default_platform,omitempty"`
	PartNumber      string           `json:"part_number,omitempty"`
	UHeight         float64          `json:"u_height"`
	FullDepth       bool             `json:"full_depth"`
	SubdeviceRole   *enum.SubdeviceRole    `json:"subdevice_role,omitempty"`
	Airflow         *enum.AirflowDirection `json:"airflow,omitempty"`
	FrontImage      *types.Image     `json:"front_image,omitempty"`
	RearImage       *types.Image     `json:"rear_image,omitempty"`
	Comments        types.Comments   `json:"comments,omitempty"`
	Created         time.Time        `json:"created"`
	Updated         time.Time        `json:"updated"`
}

// Validate проверяет корректность типа устройства
func (dt *DeviceType) Validate() error {
	if dt.Model == "" {
		return types.ErrValidationFailed
	}
	if dt.Manufacturer == nil {
		return types.ErrValidationFailed
	}
	if err := dt.Slug.Validate(); err != nil {
		return err
	}
	if dt.UHeight < 0 {
		return types.ErrValidationFailed
	}
	return nil
}

// Platform представляет платформу устройства (ОС/ПО)
type Platform struct {
	ID           types.ID        `json:"id"`
	Name         string          `json:"name"`
	Slug         types.Slug      `json:"slug"`
	Manufacturer *Manufacturer   `json:"manufacturer,omitempty"`
	NapalmDriver string          `json:"napalm_driver,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
}

// Validate проверяет корректность платформы
func (p *Platform) Validate() error {
	if p.Name == "" {
		return types.ErrValidationFailed
	}
	if err := p.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// DeviceRole представляет роль устройства
type DeviceRole struct {
	ID                 types.ID        `json:"id"`
	Name               string          `json:"name"`
	Slug               types.Slug      `json:"slug"`
	Color              string          `json:"color"` // HEX color code
	VMRole             bool            `json:"vm_role"`
	ConfigTemplate     *ConfigTemplate `json:"config_template,omitempty"`
	Description        types.Description `json:"description,omitempty"`
	Created            time.Time       `json:"created"`
	Updated            time.Time       `json:"updated"`
}

// ConfigTemplate представляет шаблон конфигурации
type ConfigTemplate struct {
	ID          types.ID   `json:"id"`
	Name        string     `json:"name"`
	DataSource  *string    `json:"data_source,omitempty"`
	DataPath    *string    `json:"data_path,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Environment interface{} `json:"environment,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность роли устройства
func (dr *DeviceRole) Validate() error {
	if dr.Name == "" {
		return types.ErrValidationFailed
	}
	if err := dr.Slug.Validate(); err != nil {
		return err
	}
	if dr.Color == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// Device представляет физическое или виртуальное устройство
type Device struct {
	ID              types.ID         `json:"id"`
	Name            string           `json:"name"`
	DeviceType      *DeviceType      `json:"device_type"`
	Role            *DeviceRole      `json:"role"`
	Tenant          *types.Tenant    `json:"tenant,omitempty"`
	Platform        *Platform        `json:"platform,omitempty"`
	Serial          string           `json:"serial,omitempty"`
	AssetTag        string           `json:"asset_tag,omitempty"`
	Site            *Site            `json:"site"`
	Location        *Location        `json:"location,omitempty"`
	Rack            *Rack            `json:"rack,omitempty"`
	Position        *int             `json:"position,omitempty"` // Rack position
	Face            *RackFace        `json:"face,omitempty"`     // Front/Rear
	Status          enum.DeviceStatus `json:"status"`
	Airflow         *enum.AirflowDirection `json:"airflow,omitempty"`
	PrimaryIPv4     *IPAddress       `json:"primary_ipv4,omitempty"`
	PrimaryIPv6     *IPAddress       `json:"primary_ipv6,omitempty"`
	Cluster         *Cluster         `json:"cluster,omitempty"`
	VirtualChassis  *VirtualChassis  `json:"virtual_chassis,omitempty"`
	VcPosition      *int             `json:"vc_position,omitempty"`
	VcPriority      *int             `json:"vc_priority,omitempty"`
	ConfigTemplate  *ConfigTemplate  `json:"config_template,omitempty"`
	Comments        types.Comments   `json:"comments,omitempty"`
	LocalContextData interface{}      `json:"local_context_data,omitempty"`
	Created         time.Time        `json:"created"`
	Updated         time.Time        `json:"updated"`
}

// Validate проверяет корректность устройства
func (d *Device) Validate() error {
	if d.Name == "" {
		return types.ErrValidationFailed
	}
	if d.DeviceType == nil {
		return types.ErrValidationFailed
	}
	if d.Role == nil {
		return types.ErrValidationFailed
	}
	if d.Site == nil {
		return types.ErrValidationFailed
	}
	if err := d.Status.Validate(); err != nil {
		return err
	}
	if d.Position != nil && (*d.Position < 1 || *d.Position > int(d.DeviceType.UHeight)) {
		return types.ErrValidationFailed
	}
	return nil
}

// GetStatusColor возвращает цвет статуса устройства
func (d *Device) GetStatusColor() string {
	switch d.Status {
	case enum.DeviceStatusActive:
		return "#5cb85c" // green
	case enum.DeviceStatusOffline:
		return "#d9534f" // red
	case enum.DeviceStatusPlanned:
		return "#f0ad4e" // orange
	case enum.DeviceStatusStaged:
		return "#428bca" // blue
	case enum.DeviceStatusFailed:
		return "#d9534f" // red
	case enum.DeviceStatusInventory:
		return "#777777" // gray
	case enum.DeviceStatusDecommissioning:
		return "#f0ad4e" // orange
	default:
		return "#777777"
	}
}

// RackFace представляет сторону стойки (передняя/задняя)
type RackFace string

const (
	RackFaceFront RackFace = "front"
	RackFaceRear  RackFace = "rear"
)

// IPAddress представляет IP адрес
type IPAddress struct {
	ID        types.ID   `json:"id"`
	Address   string     `json:"address"`
	AssignedObjectID *types.ID `json:"assigned_object_id,omitempty"`
}

// Cluster представляет кластер виртуализации
type Cluster struct {
	ID       types.ID   `json:"id"`
	Name     string     `json:"name"`
	ClusterType *ClusterType `json:"cluster_type"`
	Site     *Site      `json:"site,omitempty"`
}

// ClusterType представляет тип кластера
type ClusterType struct {
	ID   types.ID   `json:"id"`
	Name string     `json:"name"`
	Slug types.Slug `json:"slug"`
}

// VirtualChassis представляет виртуальное шасси
type VirtualChassis struct {
	ID       types.ID   `json:"id"`
	Name     string     `json:"name"`
	Domain   *string    `json:"domain,omitempty"`
	Master   *Device    `json:"master,omitempty"`
}

// ModuleType представляет тип модуля
type ModuleType struct {
	ID           types.ID        `json:"id"`
	Manufacturer *Manufacturer   `json:"manufacturer"`
	Model        string          `json:"model"`
	PartNumber   string          `json:"part_number,omitempty"`
	Weight       *float64        `json:"weight,omitempty"`
	WeightUnit   *enum.WeightUnit     `json:"weight_unit,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments   `json:"comments,omitempty"`
	Created      time.Time        `json:"created"`
	Updated      time.Time        `json:"updated"`
}

// Validate проверяет корректность типа модуля
func (mt *ModuleType) Validate() error {
	if mt.Model == "" {
		return types.ErrValidationFailed
	}
	if mt.Manufacturer == nil {
		return types.ErrValidationFailed
	}
	return nil
}

// ModuleBayTemplate представляет шаблон отсека для модулей
type ModuleBayTemplate struct {
	ID           types.ID   `json:"id"`
	DeviceType   *DeviceType `json:"device_type"`
	Name         string     `json:"name"`
	Label        string     `json:"label,omitempty"`
	Position     *string    `json:"position,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Created      time.Time  `json:"created"`
	Updated      time.Time  `json:"updated"`
}

// Validate проверяет корректность шаблона отсека для модулей
func (mbt *ModuleBayTemplate) Validate() error {
	if mbt.Name == "" {
		return types.ErrValidationFailed
	}
	if mbt.DeviceType == nil {
		return types.ErrValidationFailed
	}
	return nil
}

// ModuleBay представляет отсек для модулей в устройстве
type ModuleBay struct {
	ID          types.ID   `json:"id"`
	Device      *Device    `json:"device"`
	Name        string     `json:"name"`
	Label       string     `json:"label,omitempty"`
	Position    *string    `json:"position,omitempty"`
	InstalledModule *Module `json:"installed_module,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность отсека для модулей
func (mb *ModuleBay) Validate() error {
	if mb.Name == "" {
		return types.ErrValidationFailed
	}
	if mb.Device == nil {
		return types.ErrValidationFailed
	}
	return nil
}

// Module представляет установленный модуль
type Module struct {
	ID           types.ID        `json:"id"`
	Device       *Device         `json:"device"`
	ModuleBay    *ModuleBay      `json:"module_bay"`
	ModuleType   *ModuleType     `json:"module_type"`
	Status       enum.ModuleStatus `json:"status"`
	Serial       string          `json:"serial,omitempty"`
	AssetTag     string          `json:"asset_tag,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments   `json:"comments,omitempty"`
	Created      time.Time        `json:"created"`
	Updated      time.Time        `json:"updated"`
}

// Validate проверяет корректность модуля
func (m *Module) Validate() error {
	if m.Device == nil {
		return types.ErrValidationFailed
	}
	if m.ModuleBay == nil {
		return types.ErrValidationFailed
	}
	if m.ModuleType == nil {
		return types.ErrValidationFailed
	}
	if err := m.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetStatusColor возвращает цвет статуса модуля
func (m *Module) GetStatusColor() string {
	switch m.Status {
	case enum.ModuleStatusActive:
		return "#5cb85c"
	case enum.ModuleStatusOffline:
		return "#d9534f"
	case enum.ModuleStatusPlanned:
		return "#f0ad4e"
	case enum.ModuleStatusStaged:
		return "#428bca"
	case enum.ModuleStatusFailed:
		return "#d9534f"
	case enum.ModuleStatusInventory:
		return "#777777"
	case enum.ModuleStatusDecommissioning:
		return "#f0ad4e"
	default:
		return "#777777"
	}
}
