// Package dashboard предоставляет функциональность панелей управления
package dashboard

// DashboardWidgetColorChoices представляет варианты цветов для виджетов панели управления
type DashboardWidgetColorChoice string

const (
	DashboardWidgetColorBlue   DashboardWidgetColorChoice = "blue"
	DashboardWidgetColorIndigo DashboardWidgetColorChoice = "indigo"
	DashboardWidgetColorPurple DashboardWidgetColorChoice = "purple"
	DashboardWidgetColorPink   DashboardWidgetColorChoice = "pink"
	DashboardWidgetColorRed    DashboardWidgetColorChoice = "red"
	DashboardWidgetColorOrange DashboardWidgetColorChoice = "orange"
	DashboardWidgetColorYellow DashboardWidgetColorChoice = "yellow"
	DashboardWidgetColorGreen  DashboardWidgetColorChoice = "green"
	DashboardWidgetColorTeal   DashboardWidgetColorChoice = "teal"
	DashboardWidgetColorCyan   DashboardWidgetColorChoice = "cyan"
	DashboardWidgetColorGray   DashboardWidgetColorChoice = "gray"
	DashboardWidgetColorBlack  DashboardWidgetColorChoice = "black"
	DashboardWidgetColorWhite  DashboardWidgetColorChoice = "white"
)

// DashboardWidgetColorChoices содержит все доступные варианты цветов
var DashboardWidgetColorChoices = []DashboardWidgetColorChoice{
	DashboardWidgetColorBlue,
	DashboardWidgetColorIndigo,
	DashboardWidgetColorPurple,
	DashboardWidgetColorPink,
	DashboardWidgetColorRed,
	DashboardWidgetColorOrange,
	DashboardWidgetColorYellow,
	DashboardWidgetColorGreen,
	DashboardWidgetColorTeal,
	DashboardWidgetColorCyan,
	DashboardWidgetColorGray,
	DashboardWidgetColorBlack,
	DashboardWidgetColorWhite,
}

// ColorLabel возвращает человеко-читаемое название цвета
func (c DashboardWidgetColorChoice) ColorLabel() string {
	labels := map[DashboardWidgetColorChoice]string{
		DashboardWidgetColorBlue:   "Blue",
		DashboardWidgetColorIndigo: "Indigo",
		DashboardWidgetColorPurple: "Purple",
		DashboardWidgetColorPink:   "Pink",
		DashboardWidgetColorRed:    "Red",
		DashboardWidgetColorOrange: "Orange",
		DashboardWidgetColorYellow: "Yellow",
		DashboardWidgetColorGreen:  "Green",
		DashboardWidgetColorTeal:   "Teal",
		DashboardWidgetColorCyan:   "Cyan",
		DashboardWidgetColorGray:   "Gray",
		DashboardWidgetColorBlack:  "Black",
		DashboardWidgetColorWhite:  "White",
	}
	if label, ok := labels[c]; ok {
		return label
	}
	return string(c)
}
