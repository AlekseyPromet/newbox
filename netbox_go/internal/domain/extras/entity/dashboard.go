// Package entity содержит сущности домена extras
package entity

import (
	"time"
)

// Dashboard представляет панель управления пользователя
type Dashboard struct {
	ID     int64              `json:"id"`
	UserID *int64             `json:"user_id,omitempty"`
	Layout []DashboardLayout  `json:"layout"`
	Config map[string]WidgetConfig `json:"config"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
}

// DashboardLayout представляет элемент макета панели управления
type DashboardLayout struct {
	ID string  `json:"id"`
	W  int     `json:"w"`
	H  int     `json:"h"`
	X  *int    `json:"x,omitempty"`
	Y  *int    `json:"y,omitempty"`
}

// WidgetConfig представляет конфигурацию виджета
type WidgetConfig struct {
	Class  string                 `json:"class"`
	Title  string                 `json:"title,omitempty"`
	Color  string                 `json:"color,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}
