// Package repository содержит интерфейсы репозиториев для всех доменов
package repository

import (
	"context"

	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
)

// SiteRepository определяет интерфейс для работы с сайтами
type SiteRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Site, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Site, error)
	List(ctx context.Context, filter SiteFilter) ([]*entity.Site, int64, error)
	Create(ctx context.Context, site *entity.Site) error
	Update(ctx context.Context, site *entity.Site) error
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
	GetByID(ctx context.Context, id string) (*entity.Rack, error)
	List(ctx context.Context, filter RackFilter) ([]*entity.Rack, int64, error)
	Create(ctx context.Context, rack *entity.Rack) error
	Update(ctx context.Context, rack *entity.Rack) error
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
	GetByID(ctx context.Context, id string) (*entity.Device, error)
	GetByName(ctx context.Context, name string) (*entity.Device, error)
	List(ctx context.Context, filter DeviceFilter) ([]*entity.Device, int64, error)
	Create(ctx context.Context, device *entity.Device) error
	Update(ctx context.Context, device *entity.Device) error
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
	GetByID(ctx context.Context, id string) (*entity.Cable, error)
	List(ctx context.Context, filter CableFilter) ([]*entity.Cable, int64, error)
	Create(ctx context.Context, cable *entity.Cable) error
	Update(ctx context.Context, cable *entity.Cable) error
	Delete(ctx context.Context, id string) error
	GetTerminations(ctx context.Context, terminationType string, terminationID string) ([]*entity.Cable, error)
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
	GetByID(ctx context.Context, id string) (*entity.PowerPanel, error)
	List(ctx context.Context, filter PowerPanelFilter) ([]*entity.PowerPanel, int64, error)
	Create(ctx context.Context, panel *entity.PowerPanel) error
	Update(ctx context.Context, panel *entity.PowerPanel) error
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
	GetByID(ctx context.Context, id string) (*entity.PowerFeed, error)
	List(ctx context.Context, filter PowerFeedFilter) ([]*entity.PowerFeed, int64, error)
	Create(ctx context.Context, feed *entity.PowerFeed) error
	Update(ctx context.Context, feed *entity.PowerFeed) error
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
