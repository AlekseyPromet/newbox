// Package widgets предоставляет реализации конкретных виджетов панели управления
package widgets

import (
	"context"
	"html/template"
)

// BookmarksWidget виджет для отображения закладок пользователя
type BookmarksWidget struct {
	BaseWidget
}

// BookmarksConfigForm форма конфигурации для BookmarksWidget
type BookmarksConfigForm struct {
	BaseConfigForm
}

// Validate проверяет корректность данных формы
func (f *BookmarksConfigForm) Validate() error {
	data := f.GetData()

	if objectTypes, exists := data["object_types"]; exists && objectTypes != nil {
		// Validate that object_types is a slice of strings
		if _, ok := objectTypes.([]string); !ok {
			return &ValidationError{Field: "object_types", Message: "Object types must be a list of strings"}
		}
	}

	if orderBy, exists := data["order_by"]; exists && orderBy != nil {
		orderByStr, ok := orderBy.(string)
		if !ok {
			return &ValidationError{Field: "order_by", Message: "Order by must be a string"}
		}
		// Validate order_by value (newest, oldest, az, za)
		validOrderings := map[string]bool{
			"created":  true,
			"-created": true,
			"az":       true,
			"za":       true,
		}
		if !validOrderings[orderByStr] {
			return &ValidationError{Field: "order_by", Message: "Invalid ordering value"}
		}
	}

	if maxItems, exists := data["max_items"]; exists && maxItems != nil {
		maxItemsInt, ok := maxItems.(int)
		if !ok || maxItemsInt < 1 {
			return &ValidationError{Field: "max_items", Message: "Max items must be greater than 0"}
		}
	}

	return nil
}

// GetName возвращает имя виджета
func (w *BookmarksWidget) GetName() string {
	return "extras.BookmarksWidget"
}

// GetDescription возвращает описание виджета
func (w *BookmarksWidget) GetDescription() string {
	return "Show your personal bookmarks"
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (w *BookmarksWidget) GetDefaultTitle() string {
	return "Bookmarks"
}

// GetDefaultConfig возвращает конфигурацию по умолчанию
func (w *BookmarksWidget) GetDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"order_by": "created", // newest first
	}
}

// GetConfigForm возвращает форму конфигурации
func (w *BookmarksWidget) GetConfigForm() ConfigForm {
	return &BookmarksConfigForm{}
}

// Render рендерит содержимое виджета
func (w *BookmarksWidget) Render(ctx context.Context, request *Request) (template.HTML, error) {
	// TODO: Implement bookmarks rendering
	// This would need to:
	// 1. Check if user is anonymous
	// 2. Fetch bookmarks for the user
	// 3. Filter by object types if specified
	// 4. Sort by order_by setting
	// 5. Limit to max_items if specified
	// 6. Render template with bookmarks
	return template.HTML("BookmarksWidget content"), nil
}

// NewBookmarksWidget создаёт новый BookmarksWidget
func NewBookmarksWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *BookmarksWidget {
	widget := &BookmarksWidget{
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

// BookmarkOrdering варианты сортировки закладок
type BookmarkOrdering string

const (
	// BookmarkOrderingNewest сортировка по убыванию даты создания
	BookmarkOrderingNewest BookmarkOrdering = "created"
	// BookmarkOrderingOldest сортировка по возрастанию даты создания
	BookmarkOrderingOldest BookmarkOrdering = "-created"
	// BookmarkOrderingAZ сортировка по алфавиту A-Z
	BookmarkOrderingAZ BookmarkOrdering = "az"
	// BookmarkOrderingZA сортировка по алфавиту Z-A
	BookmarkOrderingZA BookmarkOrdering = "za"
)
