// Package entity содержит сущности домена Circuits
package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/circuits/enum"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// Provider представляет провайдера телекоммуникационных услуг
type Provider struct {
	ID           types.ID        `json:"id"`
	Name         string          `json:"name"`
	Slug         types.Slug      `json:"slug"`
	ASN          *uint32         `json:"asn,omitempty"`
	Account      string          `json:"account,omitempty"`
	PortalURL    string          `json:"portal_url,omitempty"`
	NOCContact   string          `json:"noc_contact,omitempty"`
	AdminContact string          `json:"admin_contact,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments  `json:"comments,omitempty"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
}

// Validate проверяет корректность провайдера
func (p *Provider) Validate() error {
	if p.Name == "" {
		return types.ErrNameRequired
	}
	if err := p.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// ProviderAccount представляет аккаунт внутри провайдера
type ProviderAccount struct {
	ID          types.ID        `json:"id"`
	ProviderID  types.ID        `json:"provider_id"`
	Account     string          `json:"account"`
	Name        string          `json:"name,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность аккаунта провайдера
func (pa *ProviderAccount) Validate() error {
	if pa.Account == "" {
		return types.ErrValidationFailed
	}
	if pa.ProviderID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// ProviderNetwork представляет сеть провайдера
type ProviderNetwork struct {
	ID          types.ID        `json:"id"`
	ProviderID  types.ID        `json:"provider_id"`
	Name        string          `json:"name"`
	ServiceID   string          `json:"service_id,omitempty"` // Provider's service ID
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность сети провайдера
func (pn *ProviderNetwork) Validate() error {
	if pn.Name == "" {
		return types.ErrNameRequired
	}
	if pn.ProviderID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// CircuitType представляет тип цепи
type CircuitType struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Color       string          `json:"color,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность типа цепи
func (ct *CircuitType) Validate() error {
	if ct.Name == "" {
		return types.ErrNameRequired
	}
	if err := ct.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Circuit представляет телекоммуникационную цепь
type Circuit struct {
	ID                types.ID         `json:"id"`
	CID               string           `json:"cid"` // Circuit ID (unique identifier)
	ProviderID        types.ID         `json:"provider_id"`
	ProviderAccountID *types.ID        `json:"provider_account_id,omitempty"`
	TypeID            types.ID         `json:"type_id"`
	Status            enum.CircuitStatus `json:"status"`
	TenantID          *types.ID        `json:"tenant_id,omitempty"`
	InstallDate       *time.Time       `json:"install_date,omitempty"`
	TerminationDate   *time.Time       `json:"termination_date,omitempty"`
	CommitRate        *int32           `json:"commit_rate,omitempty"` // Kbps
	Distance          *float64         `json:"distance,omitempty"`
	DistanceUnit      *string          `json:"distance_unit,omitempty"` // km, m, mi, ft
	Description       types.Description `json:"description,omitempty"`
	Comments          types.Comments   `json:"comments,omitempty"`
	Created           time.Time        `json:"created"`
	Updated           time.Time        `json:"updated"`
	TerminationAID    *types.ID        `json:"termination_a_id,omitempty"`
	TerminationZID    *types.ID        `json:"termination_z_id,omitempty"`
}

// Validate проверяет корректность цепи
func (c *Circuit) Validate() error {
	if c.CID == "" {
		return types.ErrValidationFailed
	}
	if c.ProviderID.String() == "" {
		return types.ErrValidationFailed
	}
	if c.TypeID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := c.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetStatusColor возвращает цвет статуса цепи
func (c *Circuit) GetStatusColor() string {
	return c.Status.Color()
}

// Clean валидирует связь между provider и provider_account
func (c *Circuit) Clean(providerAccount *ProviderAccount) error {
	if c.ProviderAccountID != nil && providerAccount != nil {
		if c.ProviderID != providerAccount.ProviderID {
			return types.ErrValidationFailed
		}
	}
	return nil
}

// CircuitTermination представляет точку завершения цепи
type CircuitTermination struct {
	ID                types.ID        `json:"id"`
	CircuitID         types.ID        `json:"circuit_id"`
	TermSide          enum.CircuitTermSide `json:"term_side"` // A or Z side
	TerminationType   string          `json:"termination_type,omitempty"`
	TerminationID     *types.ID       `json:"termination_id,omitempty"`
	PortSpeed         *int32          `json:"port_speed,omitempty"` // Kbps
	UpstreamSpeed     *int32          `json:"upstream_speed,omitempty"` // Kbps for asymmetric circuits
	XConnectID        string          `json:"xconnect_id,omitempty"` // Cross-connect ID
	PPInfo            string          `json:"pp_info,omitempty"` // Patch panel/port info
	Description       types.Description `json:"description,omitempty"`
	Created           time.Time       `json:"created"`
	Updated           time.Time       `json:"updated"`
	// Cached associations for filtering
	ProviderNetworkID *types.ID `json:"provider_network_id,omitempty"`
	RegionID          *types.ID `json:"region_id,omitempty"`
	SiteGroupID       *types.ID `json:"site_group_id,omitempty"`
	SiteID            *types.ID `json:"site_id,omitempty"`
	LocationID        *types.ID `json:"location_id,omitempty"`
}

// Validate проверяет корректность точки завершения цепи
func (ct *CircuitTermination) Validate() error {
	if ct.CircuitID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := ct.TermSide.Validate(); err != nil {
		return err
	}
	// Termination must be set
	if ct.TerminationType == "" || ct.TerminationID == nil {
		return types.ErrValidationFailed
	}
	return nil
}

// CacheRelatedObjects кэширует связанные объекты для фильтрации
func (ct *CircuitTermination) CacheRelatedObjects(terminationType string, terminationData interface{}) {
	// Reset all cached fields
	ct.ProviderNetworkID = nil
	ct.RegionID = nil
	ct.SiteGroupID = nil
	ct.SiteID = nil
	ct.LocationID = nil

	// Cache based on termination type
	switch terminationType {
	case "dcim.region":
		if data, ok := terminationData.(*Region); ok {
			ct.RegionID = &data.ID
		}
	case "dcim.sitegroup":
		if data, ok := terminationData.(*SiteGroup); ok {
			ct.SiteGroupID = &data.ID
		}
	case "dcim.site":
		if data, ok := terminationData.(*Site); ok {
			ct.RegionID = data.RegionID
			ct.SiteGroupID = data.GroupID
			ct.SiteID = &data.ID
		}
	case "dcim.location":
		if data, ok := terminationData.(*Location); ok {
			ct.RegionID = data.RegionID
			ct.SiteGroupID = data.GroupID
			ct.SiteID = data.SiteID
			ct.LocationID = &data.ID
		}
	case "circuits.providernetwork":
		if data, ok := terminationData.(*ProviderNetwork); ok {
			ct.ProviderNetworkID = &data.ID
		}
	}
}

// CircuitGroup представляет административную группу цепей
type CircuitGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы цепей
func (cg *CircuitGroup) Validate() error {
	if cg.Name == "" {
		return types.ErrNameRequired
	}
	if err := cg.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// CircuitGroupAssignment представляет привязку цепи к группе
type CircuitGroupAssignment struct {
	ID        types.ID              `json:"id"`
	MemberType string               `json:"member_type"`
	MemberID  types.ID              `json:"member_id"`
	GroupID   types.ID              `json:"group_id"`
	Priority  *enum.CircuitPriority `json:"priority,omitempty"`
	Created   time.Time             `json:"created"`
	Updated   time.Time             `json:"updated"`
}

// Validate проверяет корректность привязки
func (cga *CircuitGroupAssignment) Validate() error {
	if cga.MemberType == "" {
		return types.ErrValidationFailed
	}
	if cga.MemberID.String() == "" {
		return types.ErrValidationFailed
	}
	if cga.GroupID.String() == "" {
		return types.ErrValidationFailed
	}
	if cga.Priority != nil {
		if err := cga.Priority.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// VirtualCircuitType представляет тип виртуальной цепи
type VirtualCircuitType struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Color       string          `json:"color,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность типа виртуальной цепи
func (vct *VirtualCircuitType) Validate() error {
	if vct.Name == "" {
		return types.ErrNameRequired
	}
	if err := vct.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// VirtualCircuit представляет виртуальную цепь
type VirtualCircuit struct {
	ID                types.ID         `json:"id"`
	CID               string           `json:"cid"`
	ProviderNetworkID types.ID         `json:"provider_network_id"`
	ProviderAccountID *types.ID        `json:"provider_account_id,omitempty"`
	TypeID            types.ID         `json:"type_id"`
	Status            enum.CircuitStatus `json:"status"`
	TenantID          *types.ID        `json:"tenant_id,omitempty"`
	Description       types.Description `json:"description,omitempty"`
	Comments          types.Comments   `json:"comments,omitempty"`
	Created           time.Time        `json:"created"`
	Updated           time.Time        `json:"updated"`
}

// Validate проверяет корректность виртуальной цепи
func (vc *VirtualCircuit) Validate() error {
	if vc.CID == "" {
		return types.ErrValidationFailed
	}
	if vc.ProviderNetworkID.String() == "" {
		return types.ErrValidationFailed
	}
	if vc.TypeID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := vc.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetStatusColor возвращает цвет статуса виртуальной цепи
func (vc *VirtualCircuit) GetStatusColor() string {
	return vc.Status.Color()
}

// Clean валидирует связь между provider_network и provider_account
func (vc *VirtualCircuit) Clean(providerNetwork *ProviderNetwork, providerAccount *ProviderAccount) error {
	if vc.ProviderAccountID != nil && providerAccount != nil && providerNetwork != nil {
		if providerNetwork.ProviderID != providerAccount.ProviderID {
			return types.ErrValidationFailed
		}
	}
	return nil
}

// Provider возвращает провайдера из provider_network
func (vc *VirtualCircuit) Provider(providerNetwork *ProviderNetwork) *Provider {
	if providerNetwork == nil {
		return nil
	}
	return &Provider{ID: providerNetwork.ProviderID}
}

// VirtualCircuitTermination представляет точку завершения виртуальной цепи
type VirtualCircuitTermination struct {
	ID             types.ID                        `json:"id"`
	VirtualCircuitID types.ID                      `json:"virtual_circuit_id"`
	Role           enum.VirtualCircuitTerminationRole `json:"role"`
	InterfaceID    types.ID                        `json:"interface_id"`
	Description    types.Description               `json:"description,omitempty"`
	Created        time.Time                       `json:"created"`
	Updated        time.Time                       `json:"updated"`
}

// Validate проверяет корректность точки завершения виртуальной цепи
func (vct *VirtualCircuitTermination) Validate() error {
	if vct.VirtualCircuitID.String() == "" {
		return types.ErrValidationFailed
	}
	if vct.InterfaceID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := vct.Role.Validate(); err != nil {
		return err
	}
	return nil
}

// GetRoleColor возвращает цвет роли
func (vct *VirtualCircuitTermination) GetRoleColor() string {
	return vct.Role.Color()
}

// Вспомогательные типы для кэширования связанных объектов
// Эти типы представляют заглушки для реальных сущностей из других модулей

// Region - заглушка для dcim.Region
type Region struct {
	ID       types.ID  `json:"id"`
	RegionID *types.ID `json:"region_id,omitempty"`
}

// SiteGroup - заглушка для dcim.SiteGroup
type SiteGroup struct {
	ID      types.ID `json:"id"`
	GroupID *types.ID `json:"group_id,omitempty"`
}

// Site - заглушка для dcim.Site
type Site struct {
	ID       types.ID  `json:"id"`
	RegionID *types.ID `json:"region_id,omitempty"`
	GroupID  *types.ID `json:"group_id,omitempty"`
}

// Location - заглушка для dcim.Location
type Location struct {
	ID       types.ID  `json:"id"`
	RegionID *types.ID `json:"region_id,omitempty"`
	GroupID  *types.ID `json:"group_id,omitempty"`
	SiteID   *types.ID `json:"site_id,omitempty"`
}
