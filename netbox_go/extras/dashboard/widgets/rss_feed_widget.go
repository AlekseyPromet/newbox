// Package widgets предоставляет реализации конкретных виджетов панели управления
package widgets

import (
	"context"
	"html/template"
	"time"
)

// RSSFeedWidget виджет для отображения RSS ленты
type RSSFeedWidget struct {
	BaseWidget
}

// RSSFeedConfigForm форма конфигурации для RSSFeedWidget
type RSSFeedConfigForm struct {
	BaseConfigForm
}

// Validate проверяет корректность данных формы
func (f *RSSFeedConfigForm) Validate() error {
	data := f.GetData()

	feedURL, ok := data["feed_url"].(string)
	if !ok || feedURL == "" {
		return &ValidationError{Field: "feed_url", Message: "Feed URL is required"}
	}

	if maxEntries, exists := data["max_entries"]; exists && maxEntries != nil {
		maxEntriesInt, ok := maxEntries.(int)
		if !ok || maxEntriesInt < 1 || maxEntriesInt > 1000 {
			return &ValidationError{Field: "max_entries", Message: "Max entries must be between 1 and 1000"}
		}
	}

	if cacheTimeout, exists := data["cache_timeout"]; exists && cacheTimeout != nil {
		cacheTimeoutInt, ok := cacheTimeout.(int)
		if !ok || cacheTimeoutInt < 600 || cacheTimeoutInt > 86400 {
			return &ValidationError{Field: "cache_timeout", Message: "Cache timeout must be between 600 and 86400 seconds"}
		}
	}

	if requestTimeout, exists := data["request_timeout"]; exists && requestTimeout != nil {
		requestTimeoutInt, ok := requestTimeout.(int)
		if !ok || requestTimeoutInt < 1 || requestTimeoutInt > 60 {
			return &ValidationError{Field: "request_timeout", Message: "Request timeout must be between 1 and 60 seconds"}
		}
	}

	return nil
}

// GetName возвращает имя виджета
func (w *RSSFeedWidget) GetName() string {
	return "extras.RSSFeedWidget"
}

// GetDescription возвращает описание виджета
func (w *RSSFeedWidget) GetDescription() string {
	return "Embed an RSS feed from an external website."
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (w *RSSFeedWidget) GetDefaultTitle() string {
	return "RSS Feed"
}

// GetDefaultWidth возвращает ширину по умолчанию
func (w *RSSFeedWidget) GetDefaultWidth() int {
	return 6
}

// GetDefaultHeight возвращает высоту по умолчанию
func (w *RSSFeedWidget) GetDefaultHeight() int {
	return 4
}

// GetDefaultConfig возвращает конфигурацию по умолчанию
func (w *RSSFeedWidget) GetDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"max_entries":       10,
		"cache_timeout":     3600, // seconds
		"request_timeout":   3,    // seconds
		"requires_internet": true,
	}
}

// GetConfigForm возвращает форму конфигурации
func (w *RSSFeedWidget) GetConfigForm() ConfigForm {
	return &RSSFeedConfigForm{}
}

// Render рендерит содержимое виджета
func (w *RSSFeedWidget) Render(ctx context.Context, request *Request) (template.HTML, error) {
	// TODO: Implement RSS feed rendering
	// This would need to:
	// 1. Fetch RSS feed from URL
	// 2. Cache the feed content
	// 3. Parse feed entries
	// 4. Render template with feed data
	return template.HTML("RSSFeedWidget content"), nil
}

// NewRSSFeedWidget создаёт новый RSSFeedWidget
func NewRSSFeedWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *RSSFeedWidget {
	widget := &RSSFeedWidget{
		BaseWidget: BaseWidget{
			ID:     id,
			Title:  title,
			Color:  color,
			Config: config,
			Width:  width,
			Height: height,
			X:      x,
			Y:      y,
		},
	}
	if widget.Title == "" {
		widget.Title = widget.GetDefaultTitle()
	}
	if widget.Config == nil {
		widget.Config = widget.GetDefaultConfig()
	}
	return widget
}

// RSSFeedEntry представляет запись RSS ленты
type RSSFeedEntry struct {
	Title       string
	Link        string
	Published   time.Time
	Description string
	Author      string
}

// RSSFeed представляет parsed RSS ленту
type RSSFeed struct {
	Title       string
	Link        string
	Description string
	Entries     []RSSFeedEntry
	Error       error
}
