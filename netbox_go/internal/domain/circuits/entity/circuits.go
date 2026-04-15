// Package entity содержит сущности домена Circuits
package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/circuits/enum"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// Provider представляет провайдера телекоммуникационных услуг
type Provider struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	ASN         *uint32         `json:"asn,omitempty"`
	Account     string          `json:"account,omitempty"`
	PortalURL   string          `json:"portal_url,omitempty"`
	NOCContact  string          `json:"noc_contact,omitempty"`
	AdminContact string         `json:"admin_contact,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
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
	Description types.Description `json:"description,omitempty"`
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
	ID           types.ID         `json:"id"`
	CID          string           `json:"cid"` // Circuit ID (unique identifier)
	ProviderID   types.ID         `json:"provider_id"`
	ProviderNetworkID *types.ID   `json:"provider_network_id,omitempty"`
	TypeID       types.ID         `json:"type_id"`
	Status       enum.CircuitStatus `json:"status"`
	TenantID     *types.ID        `json:"tenant_id,omitempty"`
	InstallDate  *time.Time       `json:"install_date,omitempty"`
	TerminationDate *time.Time    `json:"termination_date,omitempty"`
	CommitRate   *int32           `json:"commit_rate,omitempty"` // Kbps
	CommitUnit   *string          `json:"commit_unit,omitempty"` // kbps, mbps, gbps
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments   `json:"comments,omitempty"`
	Created      time.Time        `json:"created"`
	Updated      time.Time        `json:"updated"`
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

// CircuitTermination представляет точку завершения цепи
type CircuitTermination struct {
	ID            types.ID        `json:"id"`
	CircuitID     types.ID        `json:"circuit_id"`
	TermSide      enum.CircuitTermSide `json:"term_side"` // A or Z side
	SiteID        *types.ID       `json:"site_id,omitempty"`
	ProviderNetworkID *types.ID   `json:"provider_network_id,omitempty"`
	PortSpeed     *int32          `json:"port_speed,omitempty"` // Kbps
	UpstreamSpeed *int32          `json:"upstream_speed,omitempty"` // Kbps for asymmetric circuits
	XConnectID    string          `json:"xconnect_id,omitempty"` // Cross-connect ID
	PatchPanel    string          `json:"patch_panel,omitempty"`
	Port          string          `json:"port,omitempty"`
	FrontPortID   *types.ID       `json:"front_port_id,omitempty"`
	RearPortID    *types.ID       `json:"rear_port_id,omitempty"`
	Description   types.Description `json:"description,omitempty"`
	Created       time.Time       `json:"created"`
	Updated       time.Time       `json:"updated"`
}

// Validate проверяет корректность точки завершения цепи
func (ct *CircuitTermination) Validate() error {
	if ct.CircuitID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := ct.TermSide.Validate(); err != nil {
		return err
	}
	// Either SiteID or ProviderNetworkID must be set
	if ct.SiteID == nil && ct.ProviderNetworkID == nil {
		return types.ErrValidationFailed
	}
	if ct.SiteID != nil && ct.ProviderNetworkID != nil {
		return types.ErrValidationFailed
	}
	return nil
}
