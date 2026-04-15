package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/enum"
)

// PowerPanel представляет панель распределения питания в стойке или помещении.
// Используется для группировки фидеров питания.
type PowerPanel struct {
	ID           types.ID         `json:"id"`
	Created      time.Time        `json:"created"`
	LastUpdated  time.Time        `json:"last_updated"`
	Site         *Site            `json:"site"` // Ссылка на сайт
	Location     *Location        `json:"location,omitempty"` // Опционально: конкретная локация
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	CustomFields map[string]any   `json:"custom_fields"`
}

// Validate проверяет корректность данных PowerPanel.
func (p *PowerPanel) Validate() error {
	if p.ID == (types.ID{}) {
		return &types.ValidationError{Field: "id", Message: "ID is required"}
	}
	if p.Site == nil || p.Site.ID == (types.ID{}) {
		return &types.ValidationError{Field: "site", Message: "Site is required"}
	}
	if p.Name == "" {
		return &types.ValidationError{Field: "name", Message: "Name is required"}
	}
	return nil
}

// PowerFeed представляет физический ввод питания (фидер) от панели до устройства.
// Содержит электрические параметры (вольты, амперы, фаза).
type PowerFeed struct {
	ID              types.ID             `json:"id"`
	Created         time.Time            `json:"created"`
	LastUpdated     time.Time            `json:"last_updated"`
	PowerPanel      *PowerPanel          `json:"power_panel"` // Ссылка на панель
	Rack            *Rack                `json:"rack,omitempty"` // Опционально: стойка назначения
	Name            string               `json:"name"`
	Status          enum.PowerFeedStatus `json:"status"`
	Type            enum.PowerFeedType   `json:"type"` // Primary/Redundant
	Supply          enum.PowerSupply     `json:"supply"` // AC/DC
	Phase           enum.PhaseType       `json:"phase"` // 1-phase/3-phase
	Voltage         int32                `json:"voltage"` // Вольты
	Amperage        int32                `json:"amperage"` // Амперы
	MaxUtilization  int32                `json:"max_utilization"` // Процент максимальной нагрузки
	AvailablePower  int32                `json:"available_power"` // Расчетное значение (Вт)
	Unit            enum.PowerUnit       `json:"unit"` // W/kW
	Description     string               `json:"description"`
	CustomFields    map[string]any       `json:"custom_fields"`
	
	// Связь с кабелем (трассировка питания) - пока как заглушка, т.к. Cable ещё не определён
	CableID         *types.ID            `json:"cable_id,omitempty"`
	CableEnd        string               `json:"cable_end"` // A/B сторона кабеля
}

// Validate проверяет корректность данных PowerFeed.
func (f *PowerFeed) Validate() error {
	if f.ID == (types.ID{}) {
		return &types.ValidationError{Field: "id", Message: "ID is required"}
	}
	if f.PowerPanel == nil || f.PowerPanel.ID == (types.ID{}) {
		return &types.ValidationError{Field: "power_panel", Message: "PowerPanel is required"}
	}
	if f.Name == "" {
		return &types.ValidationError{Field: "name", Message: "Name is required"}
	}
	if f.Voltage <= 0 {
		return &types.ValidationError{Field: "voltage", Message: "Voltage must be positive"}
	}
	if f.Amperage <= 0 {
		return &types.ValidationError{Field: "amperage", Message: "Amperage must be positive"}
	}
	if f.MaxUtilization < 0 || f.MaxUtilization > 100 {
		return &types.ValidationError{Field: "max_utilization", Message: "Max utilization must be between 0 and 100"}
	}
	return nil
}

// GetAvailablePower рассчитывает доступную мощность в Ваттах.
// Если поле AvailablePower не заполнено вручную, вычисляется автоматически.
func (f *PowerFeed) GetAvailablePower() int32 {
	if f.AvailablePower > 0 {
		return f.AvailablePower
	}
	// P = V * A * (sqrt(3) для 3-phase) * (Utilization / 100)
	var power float64
	power = float64(f.Voltage) * float64(f.Amperage)
	
	if f.Phase == enum.PhaseThree {
		power *= 1.732 // sqrt(3)
	}
	
	power *= float64(f.MaxUtilization) / 100.0
	return int32(power)
}

// GetStatusColor возвращает цвет статуса для UI.
func (f *PowerFeed) GetStatusColor() string {
	return f.Status.Color()
}
