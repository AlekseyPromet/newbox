// Package widgets предоставляет реализации конкретных виджетов панели управления
package widgets

import (
	"context"
	"html/template"
)

// ObjectCountsWidget виджет для отображения количества объектов по моделям
type ObjectCountsWidget struct {
	BaseWidget
}

// ObjectCountsConfigForm форма конфигурации для ObjectCountsWidget
type ObjectCountsConfigForm struct {
	BaseConfigForm
}

// Validate проверяет корректность данных формы
func (f *ObjectCountsConfigForm) Validate() error {
	data := f.GetData()
	models, ok := data["models"].([]string)
	if !ok || len(models) == 0 {
		return &ValidationError{Field: "models", Message: "At least one model must be selected"}
	}

	if filters, exists := data["filters"]; exists && filters != nil {
		// Validate that filters is a map
		if _, ok := filters.(map[string]interface{}); !ok {
			return &ValidationError{Field: "filters", Message: "Invalid format. Object filters must be passed as a dictionary."}
		}
	}

	return nil
}

// GetName возвращает имя виджета
func (w *ObjectCountsWidget) GetName() string {
	return "extras.ObjectCountsWidget"
}

// GetDescription возвращает описание виджета
func (w *ObjectCountsWidget) GetDescription() string {
	return "Display a set of NetBox models and the number of objects created for each type."
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (w *ObjectCountsWidget) GetDefaultTitle() string {
	return "Object Counts"
}

// GetConfigForm возвращает форму конфигурации
func (w *ObjectCountsWidget) GetConfigForm() ConfigForm {
	return &ObjectCountsConfigForm{}
}

// Render рендерит содержимое виджета
func (w *ObjectCountsWidget) Render(ctx context.Context, request *Request) (template.HTML, error) {
	// TODO: Implement object counts rendering
	// This would need to:
	// 1. Get models from config
	// 2. Check user permissions
	// 3. Count objects for each model
	// 4. Apply filters if specified
	// 5. Render template with counts
	return template.HTML("ObjectCountsWidget content"), nil
}

// NewObjectCountsWidget создаёт новый ObjectCountsWidget
func NewObjectCountsWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *ObjectCountsWidget {
	widget := &ObjectCountsWidget{
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
