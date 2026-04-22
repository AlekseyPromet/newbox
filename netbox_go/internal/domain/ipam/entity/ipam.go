// Package entity содержит сущности домена IPAM
package entity

import (
	"net/netip"
	"time"

	"netbox_go/internal/domain/ipam/enum"
	"netbox_go/pkg/types"
)

// VRF представляет Virtual Routing and Forwarding instance
type VRF struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Description types.Description `json:"description,omitempty"`
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	ImportTargets []string      `json:"import_targets,omitempty"` // Route targets for import
	ExportTargets []string      `json:"export_targets,omitempty"` // Route targets for export
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность VRF
func (v *VRF) Validate() error {
	if v.Name == "" {
		return types.ErrNameRequired
	}
	return nil
}

// RIR представляет Regional Internet Registry
type RIR struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	IsPrivate   bool            `json:"is_private"` // Private address space
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность RIR
func (r *RIR) Validate() error {
	if r.Name == "" {
		return types.ErrNameRequired
	}
	if err := r.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Aggregate представляет агрегат IP адресов
type Aggregate struct {
	ID          types.ID        `json:"id"`
	Prefix      netip.Prefix    `json:"prefix"`
	RIRID       types.ID        `json:"rir_id"`
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	DateAdded   *time.Time      `json:"date_added,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность агрегата
func (a *Aggregate) Validate() error {
	if !a.Prefix.IsValid() {
		return types.ErrValidationFailed
	}
	if a.RIRID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// Prefix представляет IP префикс (подсеть)
type Prefix struct {
	ID           types.ID         `json:"id"`
	Prefix       netip.Prefix     `json:"prefix"`
	SiteID       *types.ID        `json:"site_id,omitempty"`
	VRFID        *types.ID        `json:"vrf_id,omitempty"`
	TenantID     *types.ID        `json:"tenant_id,omitempty"`
	VLANID       *types.ID        `json:"vlan_id,omitempty"`
	Status       enum.PrefixStatus `json:"status"`
	RoleID       *types.ID        `json:"role_id,omitempty"`
	IsPool       bool             `json:"is_pool"` // Marks prefix as available for IP allocation
	Marks        []string         `json:"marks,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments   `json:"comments,omitempty"`
	Created      time.Time        `json:"created"`
	Updated      time.Time        `json:"updated"`
}

// Validate проверяет корректность префикса
func (p *Prefix) Validate() error {
	if !p.Prefix.IsValid() {
		return types.ErrValidationFailed
	}
	if err := p.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetIPVersion возвращает версию IP префикса
func (p *Prefix) GetIPVersion() enum.IPVersion {
	if p.Prefix.Addr().Is4() {
		return enum.IPVersion4
	}
	return enum.IPVersion6
}

// GetNetworkSize возвращает размер сети в количестве адресов
func (p *Prefix) GetNetworkSize() uint64 {
	return 1 << (uint64(p.Prefix.Bits()) ^ 128)
}

// IPAddressRole представляет роль IP адреса
type IPAddressRole struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Color       string          `json:"color"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность роли IP адреса
func (r *IPAddressRole) Validate() error {
	if r.Name == "" {
		return types.ErrNameRequired
	}
	if err := r.Slug.Validate(); err != nil {
		return err
	}
	if r.Color == "" {
		return types.ErrColorRequired
	}
	return nil
}

// IPAddress представляет IP адрес
type IPAddress struct {
	ID            types.ID            `json:"id"`
	Address       netip.Addr          `json:"address"`
	PrefixLength  uint8               `json:"prefix_length"`
	VRFID         *types.ID           `json:"vrf_id,omitempty"`
	TenantID      *types.ID           `json:"tenant_id,omitempty"`
	Status        enum.IPAddressStatus `json:"status"`
	RoleID        *types.ID           `json:"role_id,omitempty"`
	AssignedObjectType *string         `json:"assigned_object_type,omitempty"` // e.g., "dcim.Interface", "virtualization.VMInterface"
	AssignedObjectID *types.ID         `json:"assigned_object_id,omitempty"`
	NATInsideID   *types.ID           `json:"nat_inside_id,omitempty"` // For NAT, points to inside address
	DNSName       string              `json:"dns_name,omitempty"`
	Description   types.Description   `json:"description,omitempty"`
	Comments      types.Comments      `json:"comments,omitempty"`
	Created       time.Time           `json:"created"`
	Updated       time.Time           `json:"updated"`
}

// Validate проверяет корректность IP адреса
func (ip *IPAddress) Validate() error {
	if !ip.Address.IsValid() {
		return types.ErrValidationFailed
	}
	if ip.PrefixLength == 0 || ip.PrefixLength > 128 {
		return types.ErrValidationFailed
	}
	if err := ip.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetIPVersion возвращает версию IP адреса
func (ip *IPAddress) GetIPVersion() enum.IPVersion {
	if ip.Address.Is4() {
		return enum.IPVersion4
	}
	return enum.IPVersion6
}

// IsPrimary проверяет, является ли адрес первичным для устройства/VM
func (ip *IPAddress) IsPrimary() bool {
	return ip.AssignedObjectID != nil && ip.Status == enum.IPAddressStatusActive
}

// VLANGroup представляет группу VLAN
type VLANGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	ScopeType   *string         `json:"scope_type,omitempty"` // e.g., "dcim.Site", "dcim.Location"
	ScopeID     *types.ID       `json:"scope_id,omitempty"`
	VLANMin     *uint16         `json:"vlan_min,omitempty"`
	VLANMax     *uint16         `json:"vlan_max,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы VLAN
func (vg *VLANGroup) Validate() error {
	if vg.Name == "" {
		return types.ErrNameRequired
	}
	if err := vg.Slug.Validate(); err != nil {
		return err
	}
	if vg.VLANMin != nil && vg.VLANMax != nil && *vg.VLANMin > *vg.VLANMax {
		return types.ErrValidationFailed
	}
	return nil
}

// VLAN представляет Virtual LAN
type VLAN struct {
	ID          types.ID        `json:"id"`
	SiteID      *types.ID       `json:"site_id,omitempty"`
	GroupID     *types.ID       `json:"group_id,omitempty"`
	VID         uint16          `json:"vid"` // VLAN ID (1-4094)
	Name        string          `json:"name"`
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	Status      enum.VLANStatus `json:"status"`
	RoleID      *types.ID       `json:"role_id,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность VLAN
func (v *VLAN) Validate() error {
	if v.Name == "" {
		return types.ErrNameRequired
	}
	if v.VID < 1 || v.VID > 4094 {
		return types.ErrValidationFailed
	}
	if err := v.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetStatusColor возвращает цвет статуса VLAN
func (v *VLAN) GetStatusColor() string {
	return v.Status.Color()
}

// ServiceRole представляет роль сервиса
type ServiceRole struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность роли сервиса
func (sr *ServiceRole) Validate() error {
	if sr.Name == "" {
		return types.ErrNameRequired
	}
	if err := sr.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Service представляет сервис (L4-L7 service)
type Service struct {
	ID          types.ID            `json:"id"`
	DeviceID    *types.ID           `json:"device_id,omitempty"`
	VirtualMachineID *types.ID      `json:"virtual_machine_id,omitempty"`
	Name        string              `json:"name"`
	Protocol    enum.ServiceProtocol `json:"protocol"`
	Ports       []uint16            `json:"ports"`
	IPAddresses []types.ID          `json:"ip_addresses,omitempty"`
	Description types.Description   `json:"description,omitempty"`
	Comments    types.Comments      `json:"comments,omitempty"`
	Created     time.Time           `json:"created"`
	Updated     time.Time           `json:"updated"`
}

// Validate проверяет корректность сервиса
func (s *Service) Validate() error {
	if s.Name == "" {
		return types.ErrNameRequired
	}
	if s.DeviceID == nil && s.VirtualMachineID == nil {
		return types.ErrValidationFailed
	}
	if len(s.Ports) == 0 {
		return types.ErrValidationFailed
	}
	if err := s.Protocol.Validate(); err != nil {
		return err
	}
	return nil
}

// Role представляет роль префикса/VLAN
type Role struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Weight      int16           `json:"weight,omitempty"` // Weight for sorting
	Color       string          `json:"color"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность роли
func (r *Role) Validate() error {
	if r.Name == "" {
		return types.ErrNameRequired
	}
	if err := r.Slug.Validate(); err != nil {
		return err
	}
	if r.Color == "" {
		return types.ErrColorRequired
	}
	return nil
}

// RouteTarget представляет BGP Route Target
type RouteTarget struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"` // e.g., "65000:100" or "1.2.3.4:100"
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность Route Target
func (rt *RouteTarget) Validate() error {
	if rt.Name == "" {
		return types.ErrNameRequired
	}
	// Basic validation for AS:NN or IP:NN format
	return nil
}

// FHRPGroup представляет First Hop Redundancy Protocol group
type FHRPGroup struct {
	ID          types.ID        `json:"id"`
	GroupID     int32           `json:"group_id"` // VRRP VRID, HSRP group number, etc.
	Name        string          `json:"name"`
	Protocol    string          `json:"protocol"` // vrrp, hsrp, glbp, carp, clusterxl, other
	AuthType    *string         `json:"auth_type,omitempty"`
	AuthKey     *string         `json:"auth_key,omitempty"`
	IPAddresses []IPAddress     `json:"ip_addresses,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность FHRP группы
func (f *FHRPGroup) Validate() error {
	if f.GroupID <= 0 || f.GroupID > 255 {
		return types.ErrValidationFailed
	}
	if f.Name == "" {
		return types.ErrNameRequired
	}
	return nil
}

// FHRPGroupAssignment представляет назначение интерфейса к FHRP группе
type FHRPGroupAssignment struct {
	ID        types.ID `json:"id"`
	GroupID   types.ID `json:"group_id"`
	Interface types.ID `json:"interface_id"` // dcim.Interface or virtualization.VMInterface
	Priority  uint8    `json:"priority"`
}

// Validate проверяет корректность назначения FHRP
func (fa *FHRPGroupAssignment) Validate() error {
	if fa.GroupID.String() == "" {
		return types.ErrValidationFailed
	}
	if fa.Interface.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}
