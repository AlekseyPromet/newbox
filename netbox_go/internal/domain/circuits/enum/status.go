// Package enum содержит перечисления для домена Circuits
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// CircuitStatus представляет статусы цепи
type CircuitStatus string

const (
	CircuitStatusPlanned     CircuitStatus = "planned"
	CircuitStatusProvisioning CircuitStatus = "provisioning"
	CircuitStatusActive      CircuitStatus = "active"
	CircuitStatusOffline     CircuitStatus = "offline"
	CircuitStatusDeprovisioning CircuitStatus = "deprovisioning"
	CircuitStatusDecommissioning CircuitStatus = "decommissioning"
)

// GetAllCircuitStatuses возвращает все возможные статусы цепи
func GetAllCircuitStatuses() []CircuitStatus {
	return []CircuitStatus{
		CircuitStatusPlanned,
		CircuitStatusProvisioning,
		CircuitStatusActive,
		CircuitStatusOffline,
		CircuitStatusDeprovisioning,
		CircuitStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса цепи
func (s CircuitStatus) Validate() error {
	switch s {
	case CircuitStatusPlanned, CircuitStatusProvisioning, CircuitStatusActive,
		CircuitStatusOffline, CircuitStatusDeprovisioning, CircuitStatusDecommissioning:
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
		return "#9e9e9e"
	case CircuitStatusProvisioning:
		return "#ff9800"
	case CircuitStatusActive:
		return "#4caf50"
	case CircuitStatusOffline:
		return "#f44336"
	case CircuitStatusDeprovisioning:
		return "#ff9800"
	case CircuitStatusDecommissioning:
		return "#f44336"
	default:
		return "#9e9e9e"
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
