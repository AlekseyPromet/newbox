// Package enum содержит перечисления для домена Circuits
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// CircuitStatus представляет статусы цепи
type CircuitStatus string

const (
	CircuitStatusPlanned        CircuitStatus = "planned"
	CircuitStatusProvisioning   CircuitStatus = "provisioning"
	CircuitStatusActive         CircuitStatus = "active"
	CircuitStatusOffline        CircuitStatus = "offline"
	CircuitStatusDeprovisioning CircuitStatus = "deprovisioning"
	CircuitStatusDecommissioned CircuitStatus = "decommissioned"
)

// GetAllCircuitStatuses возвращает все возможные статусы цепи
func GetAllCircuitStatuses() []CircuitStatus {
	return []CircuitStatus{
		CircuitStatusPlanned,
		CircuitStatusProvisioning,
		CircuitStatusActive,
		CircuitStatusOffline,
		CircuitStatusDeprovisioning,
		CircuitStatusDecommissioned,
	}
}

// Validate проверяет корректность статуса цепи
func (s CircuitStatus) Validate() error {
	switch s {
	case CircuitStatusPlanned, CircuitStatusProvisioning, CircuitStatusActive,
		CircuitStatusOffline, CircuitStatusDeprovisioning, CircuitStatusDecommissioned:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса цепи
func (s CircuitStatus) String() string {
	return string(s)
}

// Color возвращает цвет статуса цепи для UI
func (s CircuitStatus) Color() string {
	switch s {
	case CircuitStatusPlanned:
		return "#17a2b8" // cyan
	case CircuitStatusProvisioning:
		return "#007bff" // blue
	case CircuitStatusActive:
		return "#28a745" // green
	case CircuitStatusOffline:
		return "#dc3545" // red
	case CircuitStatusDeprovisioning:
		return "#ffc107" // yellow
	case CircuitStatusDecommissioned:
		return "#6c757d" // gray
	default:
		return "#6c757d"
	}
}

// CircuitTermSide представляет сторону завершения цепи
type CircuitTermSide string

const (
	CircuitTermSideA CircuitTermSide = "A"
	CircuitTermSideZ CircuitTermSide = "Z"
)

// GetAllCircuitTermSides возвращает все возможные стороны завершения
func GetAllCircuitTermSides() []CircuitTermSide {
	return []CircuitTermSide{
		CircuitTermSideA,
		CircuitTermSideZ,
	}
}

// Validate проверяет корректность стороны завершения
func (s CircuitTermSide) Validate() error {
	switch s {
	case CircuitTermSideA, CircuitTermSideZ:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление стороны завершения
func (s CircuitTermSide) String() string {
	return string(s)
}

// CircuitPriority представляет приоритеты назначения цепей в группах
type CircuitPriority string

const (
	CircuitPriorityPrimary   CircuitPriority = "primary"
	CircuitPrioritySecondary CircuitPriority = "secondary"
	CircuitPriorityTertiary  CircuitPriority = "tertiary"
	CircuitPriorityInactive  CircuitPriority = "inactive"
)

// GetAllCircuitPriorities возвращает все возможные приоритеты
func GetAllCircuitPriorities() []CircuitPriority {
	return []CircuitPriority{
		CircuitPriorityPrimary,
		CircuitPrioritySecondary,
		CircuitPriorityTertiary,
		CircuitPriorityInactive,
	}
}

// Validate проверяет корректность приоритета
func (p CircuitPriority) Validate() error {
	switch p {
	case CircuitPriorityPrimary, CircuitPrioritySecondary,
		CircuitPriorityTertiary, CircuitPriorityInactive:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление приоритета
func (p CircuitPriority) String() string {
	return string(p)
}

// VirtualCircuitTerminationRole представляет роли завершения виртуальных цепей
type VirtualCircuitTerminationRole string

const (
	VirtualCircuitTerminationRolePeer  VirtualCircuitTerminationRole = "peer"
	VirtualCircuitTerminationRoleHub   VirtualCircuitTerminationRole = "hub"
	VirtualCircuitTerminationRoleSpoke VirtualCircuitTerminationRole = "spoke"
)

// GetAllVirtualCircuitTerminationRoles возвращает все возможные роли
func GetAllVirtualCircuitTerminationRoles() []VirtualCircuitTerminationRole {
	return []VirtualCircuitTerminationRole{
		VirtualCircuitTerminationRolePeer,
		VirtualCircuitTerminationRoleHub,
		VirtualCircuitTerminationRoleSpoke,
	}
}

// Validate проверяет корректность роли
func (r VirtualCircuitTerminationRole) Validate() error {
	switch r {
	case VirtualCircuitTerminationRolePeer,
		VirtualCircuitTerminationRoleHub,
		VirtualCircuitTerminationRoleSpoke:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление роли
func (r VirtualCircuitTerminationRole) String() string {
	return string(r)
}

// Color возвращает цвет роли для UI
func (r VirtualCircuitTerminationRole) Color() string {
	switch r {
	case VirtualCircuitTerminationRolePeer:
		return "#28a745" // green
	case VirtualCircuitTerminationRoleHub:
		return "#007bff" // blue
	case VirtualCircuitTerminationRoleSpoke:
		return "#fd7e14" // orange
	default:
		return "#6c757d"
	}
}
