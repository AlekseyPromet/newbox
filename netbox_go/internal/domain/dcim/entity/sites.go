// Package entity содержит сущности домена DCIM
package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/enum"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// Region представляет регион - географическую коллекцию сайтов
type Region struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	ParentID    *types.ID       `json:"parent_id,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность региона
func (r *Region) Validate() error {
	if err := r.Slug.Validate(); err != nil {
		return err
	}
	if len(r.Name) == 0 || len(r.Name) > 100 {
		return types.ErrValidationFailed
	}
	return nil
}

// SiteGroup представляет группу сайтов - произвольное объединение сайтов
type SiteGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	ParentID    *types.ID       `json:"parent_id,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы сайтов
func (sg *SiteGroup) Validate() error {
	if err := sg.Slug.Validate(); err != nil {
		return err
	}
	if len(sg.Name) == 0 || len(sg.Name) > 100 {
		return types.ErrValidationFailed
	}
	return nil
}

// Site представляет сайт - географическое расположение в сети
type Site struct {
	ID             types.ID              `json:"id"`
	Name           string                `json:"name"`
	Slug           types.Slug            `json:"slug"`
	Status         enum.SiteStatus       `json:"status"`
	RegionID       *types.ID             `json:"region_id,omitempty"`
	GroupID        *types.ID             `json:"group_id,omitempty"`
	TenantID       *types.ID             `json:"tenant_id,omitempty"`
	Facility       types.Facility        `json:"facility,omitempty"`
	ASNIDs         []types.ID            `json:"asn_ids,omitempty"`
	TimeZone       types.TimeZone        `json:"time_zone,omitempty"`
	PhysicalAddress types.Address        `json:"physical_address,omitempty"`
	ShippingAddress types.Address        `json:"shipping_address,omitempty"`
	Latitude       *float64              `json:"latitude,omitempty"`
	Longitude      *float64              `json:"longitude,omitempty"`
	Description    types.Description     `json:"description,omitempty"`
	Comments       types.Comments        `json:"comments,omitempty"`
	Created        time.Time             `json:"created"`
	Updated        time.Time             `json:"updated"`
}

// Validate проверяет корректность сайта
func (s *Site) Validate() error {
	if err := s.Slug.Validate(); err != nil {
		return err
	}
	if len(s.Name) == 0 || len(s.Name) > 100 {
		return types.ErrValidationFailed
	}
	if err := s.Status.Validate(); err != nil {
		return err
	}
	// Проверка координат если они указаны
	if s.Latitude != nil || s.Longitude != nil {
		coord := types.Coordinate{}
		if s.Latitude != nil {
			coord.Latitude = *s.Latitude
		}
		if s.Longitude != nil {
			coord.Longitude = *s.Longitude
		}
		if err := coord.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetStatusColor возвращает цвет для статуса сайта
func (s *Site) GetStatusColor() string {
	colors := map[enum.SiteStatus]string{
		enum.SiteStatusPlanned:  "#9e9e9e",
		enum.SiteStatusStaging:  "#ff9800",
		enum.SiteStatusActive:   "#4caf50",
		enum.SiteStatusRetired:  "#f44336",
	}
	return colors[s.Status]
}

// Location представляет локацию - подгруппу стоек и/или устройств внутри сайта
type Location struct {
	ID          types.ID              `json:"id"`
	Name        string                `json:"name"`
	Slug        types.Slug            `json:"slug"`
	SiteID      types.ID              `json:"site_id"`
	Status      enum.LocationStatus   `json:"status"`
	ParentID    *types.ID             `json:"parent_id,omitempty"`
	TenantID    *types.ID             `json:"tenant_id,omitempty"`
	Facility    types.Facility        `json:"facility,omitempty"`
	Description types.Description     `json:"description,omitempty"`
	Comments    types.Comments        `json:"comments,omitempty"`
	Created     time.Time             `json:"created"`
	Updated     time.Time             `json:"updated"`
}

// Validate проверяет корректность локации
func (l *Location) Validate() error {
	if err := l.Slug.Validate(); err != nil {
		return err
	}
	if len(l.Name) == 0 || len(l.Name) > 100 {
		return types.ErrValidationFailed
	}
	if err := l.Status.Validate(); err != nil {
		return err
	}
	if l.SiteID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// GetStatusColor возвращает цвет для статуса локации
func (l *Location) GetStatusColor() string {
	colors := map[enum.LocationStatus]string{
		enum.LocationStatusPlanned:  "#9e9e9e",
		enum.LocationStatusStaging:  "#ff9800",
		enum.LocationStatusActive:   "#4caf50",
		enum.LocationStatusRetired:  "#f44336",
	}
	return colors[l.Status]
}
