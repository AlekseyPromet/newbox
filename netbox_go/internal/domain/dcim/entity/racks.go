package entity

import (
	"time"

	"netbox_go/internal/domain/dcim/enum"
	"netbox_go/pkg/types"
)

// RackType представляет тип стойки - шаблон для стоек
type RackType struct {
	ID             types.ID              `json:"id"`
	ManufacturerID types.ID              `json:"manufacturer_id"`
	Model          string                `json:"model"`
	Slug           types.Slug            `json:"slug"`
	Description    types.Description     `json:"description,omitempty"`
	FormFactor     enum.RackType         `json:"form_factor"`
	Width          int16                 `json:"width"` // 19 or 23 inches
	UHeight        int16                 `json:"u_height"`
	StartingUnit   int16                 `json:"starting_unit"`
	DescUnits      bool                  `json:"desc_units"`
	OuterWidth     *int16                `json:"outer_width,omitempty"`
	OuterHeight    *int16                `json:"outer_height,omitempty"`
	OuterDepth     *int16                `json:"outer_depth,omitempty"`
	OuterUnit      *enum.RackDimensionUnit `json:"outer_unit,omitempty"`
	MountingDepth  *int16                `json:"mounting_depth,omitempty"`
	Weight         *int32                `json:"weight,omitempty"`
	MaxWeight      *int32                `json:"max_weight,omitempty"`
	WeightUnit     *string               `json:"weight_unit,omitempty"`
	Created        time.Time             `json:"created"`
	Updated        time.Time             `json:"updated"`
}

// Validate проверяет корректность типа стойки
func (rt *RackType) Validate() error {
	if err := rt.Slug.Validate(); err != nil {
		return err
	}
	if len(rt.Model) == 0 || len(rt.Model) > 100 {
		return types.ErrValidationFailed
	}
	if rt.ManufacturerID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := rt.FormFactor.Validate(); err != nil {
		return err
	}
	// Проверка высоты стойки (1-1000U)
	if rt.UHeight < 1 || rt.UHeight > 1000 {
		return types.ErrValidationFailed
	}
	// Проверка ширины (19 или 23 дюйма)
	if rt.Width != 19 && rt.Width != 23 {
		return types.ErrValidationFailed
	}
	return nil
}

// FullName возвращает полное название типа стойки
func (rt *RackType) FullName(manufacturerName string) string {
	return manufacturerName + " " + rt.Model
}

// Units возвращает список единиц стойки сверху вниз
func (rt *RackType) Units() []float64 {
	units := make([]float64, 0, rt.UHeight*2)
	start := float64(rt.StartingUnit)
	if rt.StartingUnit == 0 {
		start = 1.0
	}
	
	if rt.DescUnits {
		// Сверху вниз
		for i := start; i <= start+float64(rt.UHeight); i += 0.5 {
			units = append(units, i)
		}
	} else {
		// Снизу вверх
		for i := start + float64(rt.UHeight) - 0.5; i >= start; i -= 0.5 {
			units = append(units, i)
		}
	}
	return units
}

// RackRole представляет роль стойки - функциональное назначение
type RackRole struct {
	ID          types.ID          `json:"id"`
	Name        string            `json:"name"`
	Slug        types.Slug        `json:"slug"`
	Color       string            `json:"color"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

// Validate проверяет корректность роли стойки
func (rr *RackRole) Validate() error {
	if err := rr.Slug.Validate(); err != nil {
		return err
	}
	if len(rr.Name) == 0 || len(rr.Name) > 100 {
		return types.ErrValidationFailed
	}
	return nil
}

// Rack представляет стойку - физическую стойку для оборудования
type Rack struct {
	ID             types.ID              `json:"id"`
	Name           string                `json:"name"`
	FacilityID     *string               `json:"facility_id,omitempty"`
	SiteID         types.ID              `json:"site_id"`
	LocationID     *types.ID             `json:"location_id,omitempty"`
	TenantID       *types.ID             `json:"tenant_id,omitempty"`
	Status         enum.RackStatus       `json:"status"`
	RoleID         *types.ID             `json:"role_id,omitempty"`
	RackTypeID     *types.ID             `json:"rack_type_id,omitempty"`
	FormFactor     *enum.RackType        `json:"form_factor,omitempty"`
	Width          int16                 `json:"width"`
	Serial         string                `json:"serial,omitempty"`
	AssetTag       *string               `json:"asset_tag,omitempty"`
	Airflow        *string               `json:"airflow,omitempty"`
	UHeight        int16                 `json:"u_height"`
	StartingUnit   int16                 `json:"starting_unit"`
	DescUnits      bool                  `json:"desc_units"`
	OuterWidth     *int16                `json:"outer_width,omitempty"`
	OuterHeight    *int16                `json:"outer_height,omitempty"`
	OuterDepth     *int16                `json:"outer_depth,omitempty"`
	OuterUnit      *enum.RackDimensionUnit `json:"outer_unit,omitempty"`
	MountingDepth  *int16                `json:"mounting_depth,omitempty"`
	Weight         *int32                `json:"weight,omitempty"`
	MaxWeight      *int32                `json:"max_weight,omitempty"`
	WeightUnit     *string               `json:"weight_unit,omitempty"`
	Description    types.Description     `json:"description,omitempty"`
	Comments       types.Comments        `json:"comments,omitempty"`
	Created        time.Time             `json:"created"`
	Updated        time.Time             `json:"updated"`
}

// Validate проверяет корректность стойки
func (r *Rack) Validate() error {
	if len(r.Name) == 0 || len(r.Name) > 100 {
		return types.ErrValidationFailed
	}
	if r.SiteID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := r.Status.Validate(); err != nil {
		return err
	}
	// Проверка высоты стойки
	if r.UHeight < 1 || r.UHeight > 1000 {
		return types.ErrValidationFailed
	}
	// Проверка ширины
	if r.Width != 19 && r.Width != 23 {
		return types.ErrValidationFailed
	}
	// Проверка внешних размеров и единицы измерения
	hasOuterDim := r.OuterWidth != nil || r.OuterHeight != nil || r.OuterDepth != nil
	if hasOuterDim && r.OuterUnit == nil {
		return types.ErrValidationFailed
	}
	// Проверка максимального веса и единицы веса
	if r.MaxWeight != nil && r.WeightUnit == nil {
		return types.ErrValidationFailed
	}
	return nil
}

// GetStatusColor возвращает цвет для статуса стойки
func (r *Rack) GetStatusColor() string {
	colors := map[enum.RackStatus]string{
		enum.RackStatusReserved:   "#9e9e9e",
		enum.RackStatusAvailable:  "#2196f3",
		enum.RackStatusPlanned:    "#ff9800",
		enum.RackStatusActive:     "#4caf50",
		enum.RackStatusDeprecated: "#f44336",
	}
	return colors[r.Status]
}

// CopyRackTypeAttrs копирует атрибуты из типа стойки
func (r *Rack) CopyRackTypeAttrs(rackType *RackType) {
	r.FormFactor = &rackType.FormFactor
	r.Width = rackType.Width
	r.UHeight = rackType.UHeight
	r.StartingUnit = rackType.StartingUnit
	r.DescUnits = rackType.DescUnits
	r.OuterWidth = rackType.OuterWidth
	r.OuterHeight = rackType.OuterHeight
	r.OuterDepth = rackType.OuterDepth
	r.OuterUnit = rackType.OuterUnit
	r.MountingDepth = rackType.MountingDepth
	r.Weight = rackType.Weight
	r.MaxWeight = rackType.MaxWeight
	r.WeightUnit = rackType.WeightUnit
}

// RackReservation представляет резервирование места в стойке
type RackReservation struct {
	ID          types.ID          `json:"id"`
	RackID      types.ID          `json:"rack_id"`
	UserID      types.ID          `json:"user_id"`
	TenantID    *types.ID         `json:"tenant_id,omitempty"`
	Units       []int16           `json:"units"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

// Validate проверяет корректность резервирования
func (rr *RackReservation) Validate() error {
	if rr.RackID.String() == "" {
		return types.ErrValidationFailed
	}
	if rr.UserID.String() == "" {
		return types.ErrValidationFailed
	}
	if len(rr.Units) == 0 {
		return types.ErrValidationFailed
	}
	return nil
}
