// Package dashboard предоставляет функциональность панелей управления
package dashboard

import (
	"context"
	"fmt"

	"github.com/AlekseyPromet/netbox_go/extras/dashboard/widgets"
	"github.com/google/uuid"
)

// Dashboard представляет панель управления пользователя
type Dashboard struct {
	ID     int64
	UserID *int64
	Layout []DashboardLayoutItem
	Config map[string]WidgetConfig
}

// DashboardLayoutItem представляет элемент макета панели управления
type DashboardLayoutItem struct {
	ID string `json:"id"`
	W  int    `json:"w"`
	H  int    `json:"h"`
	X  *int   `json:"x,omitempty"`
	Y  *int   `json:"y,omitempty"`
}

// WidgetConfig представляет конфигурацию виджета
type WidgetConfig struct {
	Class  string                 `json:"class"`
	Title  string                 `json:"title,omitempty"`
	Color  string                 `json:"color,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// DefaultWidgetConfig представляет конфигурацию виджета по умолчанию
type DefaultWidgetConfig struct {
	Widget string                 `json:"widget"`
	Width  int                    `json:"width"`
	Height int                    `json:"height"`
	X      *int                   `json:"x,omitempty"`
	Y      *int                   `json:"y,omitempty"`
	Title  string                 `json:"title,omitempty"`
	Color  string                 `json:"color,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// GetDashboard возвращает панель управления для данного пользователя или создаёт панель по умолчанию
func GetDashboard(ctx context.Context, user *widgets.User) (*Dashboard, error) {
	if user.IsAnonymous {
		return GetDefaultDashboard(nil)
	}

	// Try to get user's dashboard from database
	dashboard, err := getUserDashboard(ctx, user.ID)
	if err != nil {
		if err == ErrDashboardNotFound {
			// Create a dashboard for this user
			dashboard, err = GetDefaultDashboard(nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create default dashboard: %w", err)
			}
			dashboard.UserID = &user.ID
			if err = saveDashboard(ctx, dashboard); err != nil {
				return nil, fmt.Errorf("failed to save dashboard: %w", err)
			}
			return dashboard, nil
		}
		return nil, fmt.Errorf("failed to get user dashboard: %w", err)
	}

	return dashboard, nil
}

// getUserDashboard получает панель управления пользователя из базы данных
// This is a placeholder - actual implementation would query the database
func getUserDashboard(ctx context.Context, userID int64) (*Dashboard, error) {
	// TODO: Implement database query
	return nil, ErrDashboardNotFound
}

// saveDashboard сохраняет панель управления в базу данных
// This is a placeholder - actual implementation would save to database
func saveDashboard(ctx context.Context, dashboard *Dashboard) error {
	// TODO: Implement database save
	return nil
}

// ErrDashboardNotFound ошибка при отсутствии панели управления
var ErrDashboardNotFound = fmt.Errorf("dashboard not found")

// GetDefaultDashboard создаёт панель управления по умолчанию
func GetDefaultDashboard(config []DefaultWidgetConfig) (*Dashboard, error) {
	dashboard := &Dashboard{
		Layout: make([]DashboardLayoutItem, 0),
		Config: make(map[string]WidgetConfig),
	}

	if config == nil {
		config = getDefaultDashboardConfig()
	}

	for _, widget := range config {
		id := generateUUID()
		dashboard.Layout = append(dashboard.Layout, DashboardLayoutItem{
			ID: id,
			W:  widget.Width,
			H:  widget.Height,
			X:  widget.X,
			Y:  widget.Y,
		})
		dashboard.Config[id] = WidgetConfig{
			Class:  widget.Widget,
			Title:  widget.Title,
			Color:  widget.Color,
			Config: widget.Config,
		}
	}

	return dashboard, nil
}

// getDefaultDashboardConfig возвращает конфигурацию панели управления по умолчанию
func getDefaultDashboardConfig() []DefaultWidgetConfig {
	return []DefaultWidgetConfig{
		{
			Widget: "extras.BookmarksWidget",
			Width:  4,
			Height: 5,
			Title:  "Bookmarks",
			Color:  "orange",
		},
		{
			Widget: "extras.ObjectCountsWidget",
			Width:  4,
			Height: 2,
			Title:  "Organization",
			Color:  "",
			Config: map[string]interface{}{
				"models": []string{
					"dcim.Site",
					"tenancy.Tenant",
					"tenancy.Contact",
				},
			},
		},
		{
			Widget: "extras.NoteWidget",
			Width:  4,
			Height: 2,
			Title:  "Welcome!",
			Color:  "green",
			Config: map[string]interface{}{
				"content": "This is your personal dashboard. Feel free to customize it by rearranging, resizing, or removing widgets. You can also add new widgets using the \"add widget\" button below. Any changes affect only _your_ dashboard, so feel free to experiment!",
			},
		},
		{
			Widget: "extras.ObjectCountsWidget",
			Width:  4,
			Height: 3,
			Title:  "IPAM",
			Color:  "",
			Config: map[string]interface{}{
				"models": []string{
					"ipam.VRF",
					"ipam.Aggregate",
					"ipam.Prefix",
					"ipam.IPRange",
					"ipam.IPAddress",
					"ipam.VLAN",
				},
			},
		},
		{
			Widget: "extras.ObjectCountsWidget",
			Width:  4,
			Height: 3,
			Title:  "Circuits",
			Color:  "",
			Config: map[string]interface{}{
				"models": []string{
					"circuits.Provider",
					"circuits.Circuit",
				},
			},
		},
		{
			Widget: "extras.ObjectCountsWidget",
			Width:  4,
			Height: 3,
			Title:  "DCIM",
			Color:  "",
			Config: map[string]interface{}{
				"models": []string{
					"dcim.Device",
					"dcim.Rack",
					"dcim.Cable",
				},
			},
		},
	}
}

// generateUUID генерирует новый UUID
func generateUUID() string {
	return uuid.New().String()
}

// InstantiateWidget создаёт экземпляр виджета из конфигурации
func InstantiateWidget(config WidgetConfig, layoutItem DashboardLayoutItem) (widgets.Widget, error) {
	// Create a new instance based on the widget type
	var newWidget widgets.Widget
	switch config.Class {
	case "extras.NoteWidget":
		newWidget = widgets.NewNoteWidget(
			layoutItem.ID,
			config.Title,
			config.Color,
			config.Config,
			layoutItem.W,
			layoutItem.H,
			layoutItem.X,
			layoutItem.Y,
		)
	case "extras.ObjectCountsWidget":
		newWidget = widgets.NewObjectCountsWidget(
			layoutItem.ID,
			config.Title,
			config.Color,
			config.Config,
			layoutItem.W,
			layoutItem.H,
			layoutItem.X,
			layoutItem.Y,
		)
	case "extras.ObjectListWidget":
		newWidget = widgets.NewObjectListWidget(
			layoutItem.ID,
			config.Title,
			config.Color,
			config.Config,
			layoutItem.W,
			layoutItem.H,
			layoutItem.X,
			layoutItem.Y,
		)
	case "extras.RSSFeedWidget":
		newWidget = widgets.NewRSSFeedWidget(
			layoutItem.ID,
			config.Title,
			config.Color,
			config.Config,
			layoutItem.W,
			layoutItem.H,
			layoutItem.X,
			layoutItem.Y,
		)
	case "extras.BookmarksWidget":
		newWidget = widgets.NewBookmarksWidget(
			layoutItem.ID,
			config.Title,
			config.Color,
			config.Config,
			layoutItem.W,
			layoutItem.H,
			layoutItem.X,
			layoutItem.Y,
		)
	default:
		return nil, fmt.Errorf("unknown widget class: %s", config.Class)
	}

	return newWidget, nil
}
