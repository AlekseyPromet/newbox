// Package enum содержит перечисления для домена IPAM
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// PrefixStatus представляет статусы префикса
type PrefixStatus string

const (
	PrefixStatusContainer   PrefixStatus = "container"
	PrefixStatusActive      PrefixStatus = "active"
	PrefixStatusReserved    PrefixStatus = "reserved"
	PrefixStatusDeprecated  PrefixStatus = "deprecated"
)

// GetAllPrefixStatuses возвращает все возможные статусы префикса
func GetAllPrefixStatuses() []PrefixStatus {
	return []PrefixStatus{
		PrefixStatusContainer,
		PrefixStatusActive,
		PrefixStatusReserved,
		PrefixStatusDeprecated,
	}
}

// Validate проверяет корректность статуса префикса
func (s PrefixStatus) Validate() error {
	switch s {
	case PrefixStatusContainer, PrefixStatusActive, PrefixStatusReserved, PrefixStatusDeprecated:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса префикса
func (s PrefixStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса префикса для UI
func (s PrefixStatus) Color() string {
	switch s {
	case PrefixStatusContainer:
		return "#9e9e9e"
	case PrefixStatusActive:
		return "#4caf50"
	case PrefixStatusReserved:
		return "#ff9800"
	case PrefixStatusDeprecated:
		return "#f44336"
	default:
		return "#9e9e9e"
	}
}

// IPAddressStatus представляет статусы IP адреса
type IPAddressStatus string

const (
	IPAddressStatusActive     IPAddressStatus = "active"
	IPAddressStatusReserved   IPAddressStatus = "reserved"
	IPAddressStatusDeprecated IPAddressStatus = "deprecated"
	IPAddressStatusDHCP       IPAddressStatus = "dhcp"
	IPAddressStatusSLAAC      IPAddressStatus = "slaac"
)

// GetAllIPAddressStatuses возвращает все возможные статусы IP адреса
func GetAllIPAddressStatuses() []IPAddressStatus {
	return []IPAddressStatus{
		IPAddressStatusActive,
		IPAddressStatusReserved,
		IPAddressStatusDeprecated,
		IPAddressStatusDHCP,
		IPAddressStatusSLAAC,
	}
}

// Validate проверяет корректность статуса IP адреса
func (s IPAddressStatus) Validate() error {
	switch s {
	case IPAddressStatusActive, IPAddressStatusReserved, IPAddressStatusDeprecated,
		IPAddressStatusDHCP, IPAddressStatusSLAAC:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса IP адреса
func (s IPAddressStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса IP адреса для UI
func (s IPAddressStatus) Color() string {
	switch s {
	case IPAddressStatusActive:
		return "#4caf50"
	case IPAddressStatusReserved:
		return "#ff9800"
	case IPAddressStatusDeprecated:
		return "#f44336"
	case IPAddressStatusDHCP:
		return "#2196f3"
	case IPAddressStatusSLAAC:
		return "#9c27b0"
	default:
		return "#9e9e9e"
	}
}

// VLANStatus представляет статусы VLAN
type VLANStatus string

const (
	VLANStatusActive     VLANStatus = "active"
	VLANStatusReserved   VLANStatus = "reserved"
	VLANStatusDeprecated VLANStatus = "deprecated"
)

// GetAllVLANStatuses возвращает все возможные статусы VLAN
func GetAllVLANStatuses() []VLANStatus {
	return []VLANStatus{
		VLANStatusActive,
		VLANStatusReserved,
		VLANStatusDeprecated,
	}
}

// Validate проверяет корректность статуса VLAN
func (s VLANStatus) Validate() error {
	switch s {
	case VLANStatusActive, VLANStatusReserved, VLANStatusDeprecated:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса VLAN
func (s VLANStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса VLAN для UI
func (s VLANStatus) Color() string {
	switch s {
	case VLANStatusActive:
		return "#4caf50"
	case VLANStatusReserved:
		return "#ff9800"
	case VLANStatusDeprecated:
		return "#f44336"
	default:
		return "#9e9e9e"
	}
}

// VLANQinQRole представляет роли QinQ VLAN
type VLANQinQRole string

const (
	VLANQinQRoleCustomer VLANQinQRole = "customer"
	VLANQinQRoleService  VLANQinQRole = "service"
)

// GetAllVLANQinQRoles возвращает все возможные роли QinQ
func GetAllVLANQinQRoles() []VLANQinQRole {
	return []VLANQinQRole{
		VLANQinQRoleCustomer,
		VLANQinQRoleService,
	}
}

// Validate проверяет корректность роли QinQ
func (r VLANQinQRole) Validate() error {
	switch r {
	case VLANQinQRoleCustomer, VLANQinQRoleService:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// ServiceProtocol представляет протоколы сервисов
type ServiceProtocol string

const (
	ServiceProtocolTCP  ServiceProtocol = "tcp"
	ServiceProtocolUDP  ServiceProtocol = "udp"
	ServiceProtocolSCTP ServiceProtocol = "sctp"
)

// GetAllServiceProtocols возвращает все возможные протоколы сервисов
func GetAllServiceProtocols() []ServiceProtocol {
	return []ServiceProtocol{
		ServiceProtocolTCP,
		ServiceProtocolUDP,
		ServiceProtocolSCTP,
	}
}

// Validate проверяет корректность протокола сервиса
func (p ServiceProtocol) Validate() error {
	switch p {
	case ServiceProtocolTCP, ServiceProtocolUDP, ServiceProtocolSCTP:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление протокола сервиса
func (p ServiceProtocol) String() string {
	return string(p)
}

// IPVersion представляет версии IP протокола
type IPVersion int

const (
	IPVersion4 IPVersion = 4
	IPVersion6 IPVersion = 6
)

// Validate проверяет корректность версии IP
func (v IPVersion) Validate() error {
	switch v {
	case IPVersion4, IPVersion6:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление версии IP
func (v IPVersion) String() string {
	if v == IPVersion4 {
		return "IPv4"
	}
	return "IPv6"
}
