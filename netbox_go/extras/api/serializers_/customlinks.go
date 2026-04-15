// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"
)

// CustomLink представляет пользовательскую ссылку
type CustomLink struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	DisplayURL  string   `json:"display_url,omitempty"`
	Display     string   `json:"display"`
	ObjectTypes []string `json:"object_types"`
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	LinkText    string   `json:"link_text"`
	LinkURL     string   `json:"link_url"`
	Weight      int      `json:"weight"`
	GroupName   string   `json:"group_name,omitempty"`
	ButtonClass string   `json:"button_class,omitempty"`
	NewWindow   bool     `json:"new_window"`
	Owner       interface{} `json:"owner,omitempty"`
	Created     *time.Time  `json:"created,omitempty"`
	LastUpdated *time.Time  `json:"last_updated,omitempty"`
}

// CustomLinkBrief краткое представление пользовательской ссылки
type CustomLinkBrief struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Display string `json:"display"`
	Name    string `json:"name"`
}
