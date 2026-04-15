// Package widgets предоставляет реализации конкретных виджетов панели управления
package widgets

import (
	"context"
	"html/template"
)

// ObjectListWidget виджет для отображения произвольного списка объектов
type ObjectListWidget struct {
	BaseWidget
}

// ObjectListConfigForm форма конфигурации для ObjectListWidget
type ObjectListConfigForm struct {
	BaseConfigForm
}

// Validate проверяет корректность данных формы
func (f *ObjectListConfigForm) Validate() error {
	data := f.GetData()
	
	model, ok := data["model"].(string)
	if !ok || model == "" {
		return &ValidationError{Field: "model", Message: "Model is required"}
	}
	
	if pageSize, exists := data["page_size"]; exists && pageSize != nil {
		pageSizeInt, ok := pageSize.(int)
		if !ok || pageSizeInt < 1 || pageSizeInt > 100 {
			return &ValidationError{Field: "page_size", Message: "Page size must be between 1 and 100"}
		}
	}
	
	if urlParams, exists := data["url_params"]; exists && urlParams != nil {
		// Validate that url_params is a map
		if _, ok := urlParams.(map[string]interface{}); !ok {
			return &ValidationError{Field: "url_params", Message: "Invalid format. URL parameters must be passed as a dictionary."}
		}
	}
	
	return nil
}

// GetName возвращает имя виджета
func (w *ObjectListWidget) GetName() string {
	return "extras.ObjectListWidget"
}

// GetDescription возвращает описание виджета
func (w *ObjectListWidget) GetDescription() string {
	return "Display an arbitrary list of objects."
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (w *ObjectListWidget) GetDefaultTitle() string {
	return "Object List"
}

// GetDefaultWidth возвращает ширину по умолчанию
func (w *ObjectListWidget) GetDefaultWidth() int {
	return 12
}

// GetDefaultHeight возвращает высоту по умолчанию
func (w *ObjectListWidget) GetDefaultHeight() int {
	return 4
}

// GetConfigForm возвращает форму конфигурации
func (w *ObjectListWidget) GetConfigForm() ConfigForm {
	return &ObjectListConfigForm{}
}

// Render рендерит содержимое виджета
func (w *ObjectListWidget) Render(ctx context.Context, request *Request) (template.HTML, error) {
	// TODO: Implement object list rendering
	// This would need to:
	// 1. Get model from config
	// 2. Check user permissions
	// 3. Build HTMX URL with parameters
	// 4. Render template with HTMX element
	return template.HTML("ObjectListWidget content"), nil
}

// NewObjectListWidget создаёт новый ObjectListWidget
func NewObjectListWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *ObjectListWidget {
	widget := &ObjectListWidget{
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
		widget.Config = make(map[string]interface{})
	}
	return widget
}
