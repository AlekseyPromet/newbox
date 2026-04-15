package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/enum"
)

// CableTermination представляет точку подключения кабеля к компоненту устройства.
// Может быть подключена к любому типу компонента (Interface, PowerPort, ConsolePort и т.д.)
type CableTermination struct {
	ID              types.ID     `json:"id"`
	Cable           *Cable       `json:"cable"`
	TerminationType string       `json:"termination_type"` // Тип объекта (e.g., "dcim.Interface", "dcim.PowerPort")
	TerminationID   types.ID     `json:"termination_id"`
}

// Validate проверяет корректность данных CableTermination.
func (t *CableTermination) Validate() error {
	if t.ID == (types.ID{}) {
		return &types.ValidationError{Field: "id", Message: "ID is required"}
	}
	if t.Cable == nil || t.Cable.ID == (types.ID{}) {
		return &types.ValidationError{Field: "cable", Message: "Cable is required"}
	}
	if t.TerminationType == "" {
		return &types.ValidationError{Field: "termination_type", Message: "Termination type is required"}
	}
	if t.TerminationID == (types.ID{}) {
		return &types.ValidationError{Field: "termination_id", Message: "Termination ID is required"}
	}
	return nil
}

// Cable представляет физическое кабельное соединение между двумя точками.
// Поддерживает различные типы кабелей и трассировку соединений.
type Cable struct {
	ID              types.ID         `json:"id"`
	Created         time.Time        `json:"created"`
	LastUpdated     time.Time        `json:"last_updated"`
	Type            enum.CableType   `json:"type"`
	Status          enum.CableStatus `json:"status"`
	Label           string           `json:"label,omitempty"`
	Color           string           `json:"color,omitempty"`
	Length          *int32           `json:"length,omitempty"` // Длина кабеля
	LengthUnit      string           `json:"length_unit,omitempty"` // Единица измерения (m, ft)
	Description     string           `json:"description"`
	CustomFields    map[string]any   `json:"custom_fields"`
	
	// Сторона A соединения
	A_Terminations  []CableTermination `json:"a_terminations"`
	
	// Сторона B соединения
	B_Terminations  []CableTermination `json:"b_terminations"`
	
	// Tenant (опционально)
	Tenant          *types.Tenant    `json:"tenant,omitempty"`
}

// Validate проверяет корректность данных Cable.
func (c *Cable) Validate() error {
	if c.ID == (types.ID{}) {
		return &types.ValidationError{Field: "id", Message: "ID is required"}
	}
	if len(c.A_Terminations) == 0 {
		return &types.ValidationError{Field: "a_terminations", Message: "At least one A termination is required"}
	}
	if len(c.B_Terminations) == 0 {
		return &types.ValidationError{Field: "b_terminations", Message: "At least one B termination is required"}
	}
	if c.Length != nil && *c.Length < 0 {
		return &types.ValidationError{Field: "length", Message: "Length must be non-negative"}
	}
	return nil
}

// GetStatusColor возвращает цвет статуса кабеля для UI.
func (c *Cable) GetStatusColor() string {
	return c.Status.Color()
}

// IsComplete проверяет, полностью ли подключен кабель (обе стороны имеют терминации).
func (c *Cable) IsComplete() bool {
	return len(c.A_Terminations) > 0 && len(c.B_Terminations) > 0
}

// AddATermination добавляет терминацию на сторону A.
func (c *Cable) AddATermination(term CableTermination) {
	c.A_Terminations = append(c.A_Terminations, term)
}

// AddBTermination добавляет терминацию на сторону B.
func (c *Cable) AddBTermination(term CableTermination) {
	c.B_Terminations = append(c.B_Terminations, term)
}
