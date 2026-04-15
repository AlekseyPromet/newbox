// Package enum содержит перечисления (enums) для домена DCIM
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// SiteStatus представляет статусы сайта
type SiteStatus string

const (
	SiteStatusPlanned     SiteStatus = "planned"
	SiteStatusStaging     SiteStatus = "staging"
	SiteStatusActive      SiteStatus = "active"
	SiteStatusRetired     SiteStatus = "retired"
)

// GetAllSiteStatuses возвращает все возможные статусы сайта
func GetAllSiteStatuses() []SiteStatus {
	return []SiteStatus{
		SiteStatusPlanned,
		SiteStatusStaging,
		SiteStatusActive,
		SiteStatusRetired,
	}
}

// Validate проверяет корректность статуса сайта
func (s SiteStatus) Validate() error {
	switch s {
	case SiteStatusPlanned, SiteStatusStaging, SiteStatusActive, SiteStatusRetired:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// LocationStatus представляет статусы локации
type LocationStatus string

const (
	LocationStatusPlanned     LocationStatus = "planned"
	LocationStatusStaging     LocationStatus = "staging"
	LocationStatusActive      LocationStatus = "active"
	LocationStatusRetired     LocationStatus = "retired"
)

// GetAllLocationStatuses возвращает все возможные статусы локации
func GetAllLocationStatuses() []LocationStatus {
	return []LocationStatus{
		LocationStatusPlanned,
		LocationStatusStaging,
		LocationStatusActive,
		LocationStatusRetired,
	}
}

// Validate проверяет корректность статуса локации
func (s LocationStatus) Validate() error {
	switch s {
	case LocationStatusPlanned, LocationStatusStaging, LocationStatusActive, LocationStatusRetired:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// RackStatus представляет статусы стойки
type RackStatus string

const (
	RackStatusReserved    RackStatus = "reserved"
	RackStatusAvailable   RackStatus = "available"
	RackStatusPlanned     RackStatus = "planned"
	RackStatusActive      RackStatus = "active"
	RackStatusDeprecated  RackStatus = "deprecated"
)

// GetAllRackStatuses возвращает все возможные статусы стойки
func GetAllRackStatuses() []RackStatus {
	return []RackStatus{
		RackStatusReserved,
		RackStatusAvailable,
		RackStatusPlanned,
		RackStatusActive,
		RackStatusDeprecated,
	}
}

// Validate проверяет корректность статуса стойки
func (s RackStatus) Validate() error {
	switch s {
	case RackStatusReserved, RackStatusAvailable, RackStatusPlanned, RackStatusActive, RackStatusDeprecated:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// DeviceStatus представляет статусы устройства
type DeviceStatus string

const (
	DeviceStatusOffline      DeviceStatus = "offline"
	DeviceStatusActive       DeviceStatus = "active"
	DeviceStatusPlanned      DeviceStatus = "planned"
	DeviceStatusStaged       DeviceStatus = "staged"
	DeviceStatusFailed       DeviceStatus = "failed"
	DeviceStatusInventory    DeviceStatus = "inventory"
	DeviceStatusDecommissioning DeviceStatus = "decommissioning"
)

// GetAllDeviceStatuses возвращает все возможные статусы устройства
func GetAllDeviceStatuses() []DeviceStatus {
	return []DeviceStatus{
		DeviceStatusOffline,
		DeviceStatusActive,
		DeviceStatusPlanned,
		DeviceStatusStaged,
		DeviceStatusFailed,
		DeviceStatusInventory,
		DeviceStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса устройства
func (s DeviceStatus) Validate() error {
	switch s {
	case DeviceStatusOffline, DeviceStatusActive, DeviceStatusPlanned,
		DeviceStatusStaged, DeviceStatusFailed, DeviceStatusInventory,
		DeviceStatusDecommissioning:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// RackType представляет типы стоек
type RackType string

const (
	RackTypeCabinet4Post RackType = "4-post-frame-cabinet"
	RackTypeCabinet2Post RackType = "2-post-frame-cabinet"
	RackTypeOpenFrame    RackType = "4-post-open-frame"
	Enclosure            RackType = "enclosure"
)

// GetAllRackTypes возвращает все возможные типы стоек
func GetAllRackTypes() []RackType {
	return []RackType{
		RackTypeCabinet4Post,
		RackTypeCabinet2Post,
		RackTypeOpenFrame,
		Enclosure,
	}
}

// Validate проверяет корректность типа стойки
func (r RackType) Validate() error {
	switch r {
	case RackTypeCabinet4Post, RackTypeCabinet2Post, RackTypeOpenFrame, Enclosure:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// RackDimensionUnit представляет единицы измерения размеров стойки
type RackDimensionUnit string

const (
	RackDimensionUnitMillimeter RackDimensionUnit = "mm"
	RackDimensionUnitInch       RackDimensionUnit = "in"
)

// GetAllRackDimensionUnits возвращает все возможные единицы измерения
func GetAllRackDimensionUnits() []RackDimensionUnit {
	return []RackDimensionUnit{
		RackDimensionUnitMillimeter,
		RackDimensionUnitInch,
	}
}

// Validate проверяет корректность единицы измерения
func (u RackDimensionUnit) Validate() error {
	switch u {
	case RackDimensionUnitMillimeter, RackDimensionUnitInch:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}
