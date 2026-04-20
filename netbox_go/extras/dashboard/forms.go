// Package dashboard предоставляет функциональность панелей управления
package dashboard

// WidgetForm представляет форму для настройки виджета
type WidgetForm struct {
	Title string
	Color string
}

// WidgetAddForm представляет форму для добавления нового виджета
type WidgetAddForm struct {
	WidgetForm
	WidgetClass string
}

// Validate проверяет корректность данных формы добавления виджета
func (f *WidgetAddForm) Validate() error {
	if f.WidgetClass == "" {
		return &ValidationError{Field: "widget_class", Message: "Widget type is required"}
	}

	// Check if widget class is registered
	_, err := GetWidgetClass(f.WidgetClass)
	if err != nil {
		return &ValidationError{Field: "widget_class", Message: "Invalid widget type"}
	}

	return nil
}

// Validate проверяет корректность данных формы виджета
func (f *WidgetForm) Validate() error {
	// Color validation is optional, so we just check if it's a valid choice if provided
	if f.Color != "" {
		validColor := false
		for _, choice := range DashboardWidgetColorChoices {
			if string(choice) == f.Color {
				validColor = true
				break
			}
		}
		if !validColor {
			return &ValidationError{Field: "color", Message: "Invalid color choice"}
		}
	}

	return nil
}

// ValidationError ошибка валидации формы
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// GetWidgetChoices возвращает список зарегистрированных виджетов
func GetWidgetChoices() []WidgetChoice {
	choices := make([]WidgetChoice, 0)
	for name, widget := range globalRegistry {
		choices = append(choices, WidgetChoice{
			Value:       name,
			Label:       widget.GetDescription(),
			Description: widget.GetDescription(),
		})
	}
	return choices
}

// WidgetChoice представляет вариант выбора виджета
type WidgetChoice struct {
	Value       string
	Label       string
	Description string
}

// NewWidgetForm создаёт новую форму виджета
func NewWidgetForm(title, color string) *WidgetForm {
	return &WidgetForm{
		Title: title,
		Color: color,
	}
}

// NewWidgetAddForm создаёт новую форму добавления виджета
func NewWidgetAddForm(widgetClass, title, color string) *WidgetAddForm {
	return &WidgetAddForm{
		WidgetForm: WidgetForm{
			Title: title,
			Color: color,
		},
		WidgetClass: widgetClass,
	}
}

// ToWidgetConfig преобразует данные формы в конфигурацию виджета
func (f *WidgetForm) ToWidgetConfig() WidgetConfig {
	return WidgetConfig{
		Title: f.Title,
		Color: f.Color,
	}
}

// GetFormData возвращает данные формы в виде карты
func (f *WidgetForm) GetFormData() map[string]interface{} {
	return map[string]interface{}{
		"title": f.Title,
		"color": f.Color,
	}
}

// GetFormData возвращает данные формы добавления виджета в виде карты
func (f *WidgetAddForm) GetFormData() map[string]interface{} {
	data := f.WidgetForm.GetFormData()
	data["widget_class"] = f.WidgetClass
	return data
}
