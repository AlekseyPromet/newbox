// Package widgets предоставляет реализации конкретных виджетов панели управления
package widgets

import (
	"context"
	"html/template"
)

// NoteWidget виджет для отображения произвольного контента с поддержкой Markdown
type NoteWidget struct {
	BaseWidget
}

// NoteConfigForm форма конфигурации для NoteWidget
type NoteConfigForm struct {
	BaseConfigForm
}

// Validate проверяет корректность данных формы
func (f *NoteConfigForm) Validate() error {
	data := f.GetData()
	content, ok := data["content"].(string)
	if !ok || content == "" {
		return &ValidationError{Field: "content", Message: "Content is required"}
	}
	return nil
}

// GetName возвращает имя виджета
func (w *NoteWidget) GetName() string {
	return "extras.NoteWidget"
}

// GetDescription возвращает описание виджета
func (w *NoteWidget) GetDescription() string {
	return "Display some arbitrary custom content. Markdown is supported."
}

// GetDefaultTitle возвращает заголовок по умолчанию
func (w *NoteWidget) GetDefaultTitle() string {
	return "Note"
}

// GetConfigForm возвращает форму конфигурации
func (w *NoteWidget) GetConfigForm() ConfigForm {
	return &NoteConfigForm{}
}

// Render рендерит содержимое виджета
func (w *NoteWidget) Render(ctx context.Context, request *Request) (template.HTML, error) {
	content, ok := w.Config["content"].(string)
	if !ok {
		return "", nil
	}
	// TODO: Implement markdown rendering
	return template.HTML(content), nil
}

// NewNoteWidget создаёт новый NoteWidget
func NewNoteWidget(id, title, color string, config map[string]interface{}, width, height int, x, y *int) *NoteWidget {
	widget := &NoteWidget{
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

// ValidationError ошибка валидации формы
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
