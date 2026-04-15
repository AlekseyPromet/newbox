// Package entity содержит сущности домена Tenancy
package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// TenantGroup представляет группу арендаторов - иерархическую организацию tenants
type TenantGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	ParentID    *types.ID       `json:"parent_id,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы арендаторов
func (tg *TenantGroup) Validate() error {
	if tg.Name == "" {
		return types.ErrNameRequired
	}
	if err := tg.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Tenant представляет арендатора - организацию или клиента
type Tenant struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	GroupID     *types.ID       `json:"group_id,omitempty"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность арендатора
func (t *Tenant) Validate() error {
	if t.Name == "" {
		return types.ErrNameRequired
	}
	if err := t.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// ContactGroup представляет группу контактов
type ContactGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	ParentID    *types.ID       `json:"parent_id,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы контактов
func (cg *ContactGroup) Validate() error {
	if cg.Name == "" {
		return types.ErrNameRequired
	}
	if err := cg.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// ContactRole представляет роль контакта
type ContactRole struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность роли контакта
func (cr *ContactRole) Validate() error {
	if cr.Name == "" {
		return types.ErrNameRequired
	}
	if err := cr.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Contact представляет контакт
type Contact struct {
	ID          types.ID        `json:"id"`
	GroupID     *types.ID       `json:"group_id,omitempty"`
	Name        string          `json:"name"`
	Title       string          `json:"title,omitempty"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Address     string          `json:"address,omitempty"`
	Link        string          `json:"link,omitempty"` // URL to profile or website
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность контакта
func (c *Contact) Validate() error {
	if c.Name == "" {
		return types.ErrNameRequired
	}
	return nil
}

// ContactAssignment представляет назначение контакта к объекту
type ContactAssignment struct {
	ID           types.ID   `json:"id"`
	ObjectType   string     `json:"object_type"` // e.g., "dcim.Site", "dcim.Rack"
	ObjectID     types.ID   `json:"object_id"`
	ContactID    types.ID   `json:"contact_id"`
	RoleID       types.ID   `json:"role_id"`
	Priority     uint8      `json:"priority,omitempty"`
	Created      time.Time  `json:"created"`
	Updated      time.Time  `json:"updated"`
}

// Validate проверяет корректность назначения контакта
func (ca *ContactAssignment) Validate() error {
	if ca.ObjectType == "" {
		return types.ErrValidationFailed
	}
	if ca.ObjectID.String() == "" {
		return types.ErrValidationFailed
	}
	if ca.ContactID.String() == "" {
		return types.ErrValidationFailed
	}
	if ca.RoleID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}
