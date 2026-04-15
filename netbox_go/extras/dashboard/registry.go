// Package dashboard предоставляет функциональность панелей управления
package dashboard

import (
	"github.com/netbox-go/extras/dashboard/widgets"
)

// WidgetRegistry реестр зарегистрированных виджетов
type WidgetRegistry map[string]widgets.Widget

// globalRegistry глобальный реестр виджетов
var globalRegistry = make(WidgetRegistry)

// RegisterWidget регистрирует класс виджета панели управления
func RegisterWidget(w widgets.Widget) {
	label := w.GetName()
	globalRegistry[label] = w
}

// GetWidgetClass возвращает зарегистрированный класс виджета по имени
func GetWidgetClass(name string) (widgets.Widget, error) {
	widget, ok := globalRegistry[name]
	if !ok {
		return nil, &WidgetNotFoundError{Name: name}
	}
	return widget, nil
}

// GetAllWidgets возвращает все зарегистрированные виджеты
func GetAllWidgets() WidgetRegistry {
	return globalRegistry
}

// WidgetNotFoundError ошибка при отсутствии зарегистрированного виджета
type WidgetNotFoundError struct {
	Name string
}

func (e *WidgetNotFoundError) Error() string {
	return "Unregistered widget class: " + e.Name
}
