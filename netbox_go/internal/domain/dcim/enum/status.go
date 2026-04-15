// Package enum содержит перечисления для домена DCIM
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// ClusterStatus представляет статусы кластера
type ClusterStatus string

const (
	ClusterStatusPlanned     ClusterStatus = "planned"
	ClusterStatusStaging     ClusterStatus = "staging"
	ClusterStatusActive      ClusterStatus = "active"
	ClusterStatusOffline     ClusterStatus = "offline"
	ClusterStatusDecommissioning ClusterStatus = "decommissioning"
)

// GetAllClusterStatuses возвращает все возможные статусы кластера
func GetAllClusterStatuses() []ClusterStatus {
	return []ClusterStatus{
		ClusterStatusPlanned,
		ClusterStatusStaging,
		ClusterStatusActive,
		ClusterStatusOffline,
		ClusterStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса кластера
func (s ClusterStatus) Validate() error {
	switch s {
	case ClusterStatusPlanned, ClusterStatusStaging, ClusterStatusActive,
		ClusterStatusOffline, ClusterStatusDecommissioning:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса кластера
func (s ClusterStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса кластера для UI
func (s ClusterStatus) Color() string {
	switch s {
	case ClusterStatusPlanned:
		return "#9e9e9e"
	case ClusterStatusStaging:
		return "#ff9800"
	case ClusterStatusActive:
		return "#4caf50"
	case ClusterStatusOffline:
		return "#f44336"
	case ClusterStatusDecommissioning:
		return "#ff9800"
	default:
		return "#9e9e9e"
	}
}

// VirtualMachineStatus представляет статусы виртуальной машины
type VirtualMachineStatus string

const (
	VirtualMachineStatusOffline      VirtualMachineStatus = "offline"
	VirtualMachineStatusActive       VirtualMachineStatus = "active"
	VirtualMachineStatusPlanned      VirtualMachineStatus = "planned"
	VirtualMachineStatusStaged       VirtualMachineStatus = "staged"
	VirtualMachineStatusFailed       VirtualMachineStatus = "failed"
	VirtualMachineStatusDecommissioning VirtualMachineStatus = "decommissioning"
)

// GetAllVirtualMachineStatuses возвращает все возможные статусы VM
func GetAllVirtualMachineStatuses() []VirtualMachineStatus {
	return []VirtualMachineStatus{
		VirtualMachineStatusOffline,
		VirtualMachineStatusActive,
		VirtualMachineStatusPlanned,
		VirtualMachineStatusStaged,
		VirtualMachineStatusFailed,
		VirtualMachineStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса VM
func (s VirtualMachineStatus) Validate() error {
	switch s {
	case VirtualMachineStatusOffline, VirtualMachineStatusActive,
		VirtualMachineStatusPlanned, VirtualMachineStatusStaged,
		VirtualMachineStatusFailed, VirtualMachineStatusDecommissioning:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса VM
func (s VirtualMachineStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса VM для UI
func (s VirtualMachineStatus) Color() string {
	switch s {
	case VirtualMachineStatusOffline:
		return "#f44336"
	case VirtualMachineStatusActive:
		return "#4caf50"
	case VirtualMachineStatusPlanned:
		return "#9e9e9e"
	case VirtualMachineStatusStaged:
		return "#2196f3"
	case VirtualMachineStatusFailed:
		return "#d32f2f"
	case VirtualMachineStatusDecommissioning:
		return "#ff9800"
	default:
		return "#9e9e9e"
	}
}

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

// ModuleStatus представляет статусы модуля
type ModuleStatus string

const (
	ModuleStatusOffline      ModuleStatus = "offline"
	ModuleStatusActive       ModuleStatus = "active"
	ModuleStatusPlanned      ModuleStatus = "planned"
	ModuleStatusStaged       ModuleStatus = "staged"
	ModuleStatusFailed       ModuleStatus = "failed"
	ModuleStatusInventory    ModuleStatus = "inventory"
	ModuleStatusDecommissioning ModuleStatus = "decommissioning"
)

// GetAllModuleStatuses возвращает все возможные статусы модуля
func GetAllModuleStatuses() []ModuleStatus {
	return []ModuleStatus{
		ModuleStatusOffline,
		ModuleStatusActive,
		ModuleStatusPlanned,
		ModuleStatusStaged,
		ModuleStatusFailed,
		ModuleStatusInventory,
		ModuleStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса модуля
func (s ModuleStatus) Validate() error {
	switch s {
	case ModuleStatusOffline, ModuleStatusActive, ModuleStatusPlanned,
		ModuleStatusStaged, ModuleStatusFailed, ModuleStatusInventory,
		ModuleStatusDecommissioning:
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

// SubdeviceRole представляет роли под-устройств
type SubdeviceRole string

const (
	SubdeviceRoleParent SubdeviceRole = "parent"
	SubdeviceRoleChild  SubdeviceRole = "child"
)

// GetAllSubdeviceRoles возвращает все возможные роли под-устройств
func GetAllSubdeviceRoles() []SubdeviceRole {
	return []SubdeviceRole{
		SubdeviceRoleParent,
		SubdeviceRoleChild,
	}
}

// Validate проверяет корректность роли под-устройства
func (r SubdeviceRole) Validate() error {
	switch r {
	case SubdeviceRoleParent, SubdeviceRoleChild:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// AirflowDirection представляет направление воздушного потока
type AirflowDirection string

const (
	AirflowFrontToRear   AirflowDirection = "front-to-rear"
	AirflowRearToFront   AirflowDirection = "rear-to-front"
	AirflowLeftToRight   AirflowDirection = "left-to-right"
	AirflowRightToLeft   AirflowDirection = "right-to-left"
	AirflowSideToRear    AirflowDirection = "side-to-rear"
	AirflowPassive       AirflowDirection = "passive"
	AirflowMixed         AirflowDirection = "mixed"
)

// GetAllAirflowDirections возвращает все возможные направления воздушного потока
func GetAllAirflowDirections() []AirflowDirection {
	return []AirflowDirection{
		AirflowFrontToRear,
		AirflowRearToFront,
		AirflowLeftToRight,
		AirflowRightToLeft,
		AirflowSideToRear,
		AirflowPassive,
		AirflowMixed,
	}
}

// Validate проверяет корректность направления воздушного потока
func (a AirflowDirection) Validate() error {
	switch a {
	case AirflowFrontToRear, AirflowRearToFront, AirflowLeftToRight,
		AirflowRightToLeft, AirflowSideToRear, AirflowPassive, AirflowMixed:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// WeightUnit представляет единицы измерения веса
type WeightUnit string

const (
	WeightUnitKilogram     WeightUnit = "kg"
	WeightUnitGram         WeightUnit = "g"
	WeightUnitPound        WeightUnit = "lb"
	WeightUnitOunce        WeightUnit = "oz"
)

// GetAllWeightUnits возвращает все возможные единицы измерения веса
func GetAllWeightUnits() []WeightUnit {
	return []WeightUnit{
		WeightUnitKilogram,
		WeightUnitGram,
		WeightUnitPound,
		WeightUnitOunce,
	}
}

// Validate проверяет корректность единицы измерения веса
func (w WeightUnit) Validate() error {
	switch w {
	case WeightUnitKilogram, WeightUnitGram, WeightUnitPound, WeightUnitOunce:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// PhaseType определяет тип электрической фазы.
type PhaseType string

const (
	PhaseSingle PhaseType = "single-phase"
	PhaseThree  PhaseType = "three-phase"
)

// GetAllPhaseTypes возвращает все возможные типы фаз
func GetAllPhaseTypes() []PhaseType {
	return []PhaseType{
		PhaseSingle,
		PhaseThree,
	}
}

// Validate проверяет корректность типа фазы
func (p PhaseType) Validate() error {
	switch p {
	case PhaseSingle, PhaseThree:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// PowerUnit определяет единицы измерения мощности.
type PowerUnit string

const (
	PowerUnitW  PowerUnit = "W"
	PowerUnitKW PowerUnit = "kW"
)

// GetAllPowerUnits возвращает все возможные единицы измерения мощности
func GetAllPowerUnits() []PowerUnit {
	return []PowerUnit{
		PowerUnitW,
		PowerUnitKW,
	}
}

// Validate проверяет корректность единицы измерения мощности
func (u PowerUnit) Validate() error {
	switch u {
	case PowerUnitW, PowerUnitKW:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// PowerFeedStatus определяет статус фидера питания.
type PowerFeedStatus string

const (
	PowerFeedPlanned    PowerFeedStatus = "planned"
	PowerFeedActive     PowerFeedStatus = "active"
	PowerFeedOffline    PowerFeedStatus = "offline"
	PowerFeedFailed     PowerFeedStatus = "failed"
)

// GetAllPowerFeedStatuses возвращает все возможные статусы фидера
func GetAllPowerFeedStatuses() []PowerFeedStatus {
	return []PowerFeedStatus{
		PowerFeedPlanned,
		PowerFeedActive,
		PowerFeedOffline,
		PowerFeedFailed,
	}
}

// Validate проверяет корректность статуса фидера
func (s PowerFeedStatus) Validate() error {
	switch s {
	case PowerFeedPlanned, PowerFeedActive, PowerFeedOffline, PowerFeedFailed:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// Color возвращает цвет статуса для UI
func (s PowerFeedStatus) Color() string {
	switch s {
	case PowerFeedPlanned:
		return "gray"
	case PowerFeedActive:
		return "green"
	case PowerFeedOffline:
		return "yellow"
	case PowerFeedFailed:
		return "red"
	default:
		return "gray"
	}
}

// PowerFeedType определяет тип фидера (основной или резервный).
type PowerFeedType string

const (
	PowerFeedPrimary   PowerFeedType = "primary"
	PowerFeedRedundant PowerFeedType = "redundant"
)

// GetAllPowerFeedTypes возвращает все возможные типы фидеров
func GetAllPowerFeedTypes() []PowerFeedType {
	return []PowerFeedType{
		PowerFeedPrimary,
		PowerFeedRedundant,
	}
}

// Validate проверяет корректность типа фидера
func (t PowerFeedType) Validate() error {
	switch t {
	case PowerFeedPrimary, PowerFeedRedundant:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// PowerSupply определяет тип тока.
type PowerSupply string

const (
	PowerSupplyAC PowerSupply = "AC"
	PowerSupplyDC PowerSupply = "DC"
)

// GetAllPowerSupplies возвращает все возможные типы питания
func GetAllPowerSupplies() []PowerSupply {
	return []PowerSupply{
		PowerSupplyAC,
		PowerSupplyDC,
	}
}

// Validate проверяет корректность типа питания
func (s PowerSupply) Validate() error {
	switch s {
	case PowerSupplyAC, PowerSupplyDC:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}
