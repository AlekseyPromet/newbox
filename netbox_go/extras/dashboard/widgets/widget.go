// Package widgets предоставляет типы виджетов для панелей управления
package widgets

import (
	"context"
	"fmt"
	"html/template"
)

// Widget интерфейс для всех виджетов панели управления
type Widget interface {
	// GetName возвращает имя виджета в формате "app.WidgetName"
	GetName() string
	// GetDescription возвращает описание функции виджета
	GetDescription() string
	// GetDefaultTitle возвращает заголовок по умолчанию
	GetDefaultTitle() string
	// GetDefaultWidth возвращает ширину по умолчанию (1-12)
	GetDefaultWidth() int
	// GetDefaultHeight возвращает высоту по умолчанию
	GetDefaultHeight() int
	// GetConfigForm возвращает форму конфигурации виджета
	GetConfigForm() ConfigForm
	// Render рендерит содержимое виджета
	Render(ctx context.Context, request *Request) (template.HTML, error)
}

// Request представляет запрос на рендеринг виджета
type Request struct {
	Context context.Context
	User    *User
}

// User представляет пользователя
type User struct {
	ID          int64
	Username    string
	IsAnonymous bool
	// Add other user fields as needed
}

// ConfigForm интерфейс для форм конфигурации виджетов
type ConfigForm interface {
	// Validate проверяет корректность данных формы
	Validate() error
	// GetData возвращает данные конфигурации
	GetData() map[string]interface{}
}

// BaseWidget базовая реализация виджета
type BaseWidget struct {
	ID            string
	Title         string
	Color         string
	Config        map[string]interface{}
	Width         int
	Height        int
	X             *int
	Y             *int
	Description   string
	DefaultTitle  string
	DefaultConfig map[string]interface{}
}

// GetName возвращает имя виджета в формате "app.WidgetName"
func (b *BaseWidget) GetName() string {
	// Implementation will be provided by concrete types
	return ""
}

// GetDescription возвращает описание функции виджета
func (b *BaseWidget) GetDescription() string {
	return b.Description
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (b *BaseWidget) GetDefaultTitle() string {
	return b.DefaultTitle
}

// GetDefaultWidth возвращает ширину по умолчанию (1-12)
func (b *BaseWidget) GetDefaultWidth() int {
	if b.Width > 0 {
		return b.Width
	}
	return 4 // default width
}

// GetDefaultHeight возвращает высоту по умолчанию
func (b *BaseWidget) GetDefaultHeight() int {
	if b.Height > 0 {
		return b.Height
	}
	return 3 // default height
}

// GetConfigForm возвращает форму конфигурации виджета
func (b *BaseWidget) GetConfigForm() ConfigForm {
	return &BaseConfigForm{}
}

// SetLayout устанавливает параметры макета из grid item
func (b *BaseWidget) SetLayout(gridItem map[string]interface{}) {
	if w, ok := gridItem["w"].(int); ok {
		b.Width = w
	}
	if h, ok := gridItem["h"].(int); ok {
		b.Height = h
	}
	if x, ok := gridItem["x"]; ok {
		if xv, ok := x.(int); ok {
			b.X = &xv
		}
	}
	if y, ok := gridItem["y"]; ok {
		if yv, ok := y.(int); ok {
			b.Y = &yv
		}
	}
}

// GetFormData возвращает данные формы виджета
func (b *BaseWidget) GetFormData() map[string]interface{} {
	return map[string]interface{}{
		"title":  b.Title,
		"color":  b.Color,
		"config": b.Config,
	}
}

// BaseConfigForm базовая форма конфигурации
type BaseConfigForm struct {
	data map[string]interface{}
}

// Validate проверяет корректность данных формы
func (f *BaseConfigForm) Validate() error {
	return nil
}

// GetData возвращает данные конфигурации
func (f *BaseConfigForm) GetData() map[string]interface{} {
	if f.data == nil {
		f.data = make(map[string]interface{})
	}
	return f.data
}

// NewBaseWidget создаёт новый базовый виджет
func NewBaseWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *BaseWidget {
	return &BaseWidget{
		ID:     id,
		Title:  title,
		Color:  color,
		Config: config,
		Width:  width,
		Height: height,
		X:      x,
		Y:      y,
	}
}

// String возвращает строковое представление виджета
func (b *BaseWidget) String() string {
	if b.Title != "" {
		return b.Title
	}
	return fmt.Sprintf("%T", b)
}
