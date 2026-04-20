// Package repository содержит интерфейсы репозиториев для всех доменов
package repository

import (
	"context"

	dcim_entity "github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
	extras_entity "github.com/AlekseyPromet/netbox_go/internal/domain/extras/entity"
)

// SiteRepository определяет интерфейс для работы с сайтами
type SiteRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Site, error)
	GetBySlug(ctx context.Context, slug string) (*dcim_entity.Site, error)
	List(ctx context.Context, filter SiteFilter) ([]*dcim_entity.Site, int64, error)
	Create(ctx context.Context, site *dcim_entity.Site) error
	Update(ctx context.Context, site *dcim_entity.Site) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// SiteFilter представляет фильтры для поиска сайтов
type SiteFilter struct {
	Status   *string
	RegionID *string
	GroupID  *string
	TenantID *string
	Limit    int
	Offset   int
}

// RackRepository определяет интерфейс для работы со стойками
type RackRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Rack, error)
	List(ctx context.Context, filter RackFilter) ([]*dcim_entity.Rack, int64, error)
	Create(ctx context.Context, rack *dcim_entity.Rack) error
	Update(ctx context.Context, rack *dcim_entity.Rack) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// RackFilter представляет фильтры для поиска стоек
type RackFilter struct {
	SiteID     *string
	LocationID *string
	Status     *string
	TypeID     *string
	Limit      int
	Offset     int
}

// DeviceRepository определяет интерфейс для работы с устройствами
type DeviceRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Device, error)
	GetByName(ctx context.Context, name string) (*dcim_entity.Device, error)
	List(ctx context.Context, filter DeviceFilter) ([]*dcim_entity.Device, int64, error)
	Create(ctx context.Context, device *dcim_entity.Device) error
	Update(ctx context.Context, device *dcim_entity.Device) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// DeviceFilter представляет фильтры для поиска устройств
type DeviceFilter struct {
	SiteID     *string
	RackID     *string
	TypeID     *string
	RoleID     *string
	TenantID   *string
	Status     *string
	PlatformID *string
	Limit      int
	Offset     int
}

// CableRepository определяет интерфейс для работы с кабелями
type CableRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Cable, error)
	List(ctx context.Context, filter CableFilter) ([]*dcim_entity.Cable, int64, error)
	Create(ctx context.Context, cable *dcim_entity.Cable) error
	Update(ctx context.Context, cable *dcim_entity.Cable) error
	Delete(ctx context.Context, id string) error
	GetTerminations(ctx context.Context, terminationType string, terminationID string) ([]*dcim_entity.Cable, error)
}

// CableFilter представляет фильтры для поиска кабелей
type CableFilter struct {
	Status *string
	Type   *string
	Limit  int
	Offset int
}

// PowerPanelRepository определяет интерфейс для работы с панелями питания
type PowerPanelRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.PowerPanel, error)
	List(ctx context.Context, filter PowerPanelFilter) ([]*dcim_entity.PowerPanel, int64, error)
	Create(ctx context.Context, panel *dcim_entity.PowerPanel) error
	Update(ctx context.Context, panel *dcim_entity.PowerPanel) error
	Delete(ctx context.Context, id string) error
}

// PowerPanelFilter представляет фильтры для поиска панелей питания
type PowerPanelFilter struct {
	SiteID     *string
	LocationID *string
	Limit      int
	Offset     int
}

// PowerFeedRepository определяет интерфейс для работы с фидерами питания
type PowerFeedRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.PowerFeed, error)
	List(ctx context.Context, filter PowerFeedFilter) ([]*dcim_entity.PowerFeed, int64, error)
	Create(ctx context.Context, feed *dcim_entity.PowerFeed) error
	Update(ctx context.Context, feed *dcim_entity.PowerFeed) error
	Delete(ctx context.Context, id string) error
}

// PowerFeedFilter представляет фильтры для поиска фидеров питания
type PowerFeedFilter struct {
	PowerPanelID *string
	RackID       *string
	Status       *string
	Limit        int
	Offset       int
}

// DashboardRepository определяет интерфейс для работы с панелями управления
type DashboardRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*extras_entity.Dashboard, error)
	Create(ctx context.Context, dashboard *extras_entity.Dashboard) error
	Update(ctx context.Context, dashboard *extras_entity.Dashboard) error
	Delete(ctx context.Context, id int64) error
}
