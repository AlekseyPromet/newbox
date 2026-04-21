// Package dashboard предоставляет функциональность панелей управления
package dashboard

import (
	"github.com/AlekseyPromet/netbox_go/extras/dashboard/widgets"
)

// Export widget types for convenience
type (
	Widget                 = widgets.Widget
	Request                = widgets.Request
	User                   = widgets.User
	ConfigForm             = widgets.ConfigForm
	BaseWidget             = widgets.BaseWidget
	BaseConfigForm         = widgets.BaseConfigForm
	NoteWidget             = widgets.NoteWidget
	NoteConfigForm         = widgets.NoteConfigForm
	ObjectCountsWidget     = widgets.ObjectCountsWidget
	ObjectCountsConfigForm = widgets.ObjectCountsConfigForm
	ObjectListWidget       = widgets.ObjectListWidget
	ObjectListConfigForm   = widgets.ObjectListConfigForm
	RSSFeedWidget          = widgets.RSSFeedWidget
	RSSFeedConfigForm      = widgets.RSSFeedConfigForm
	BookmarksWidget        = widgets.BookmarksWidget
	BookmarksConfigForm    = widgets.BookmarksConfigForm
)

// Export constructors
var (
	NewBaseWidget         = widgets.NewBaseWidget
	NewNoteWidget         = widgets.NewNoteWidget
	NewObjectCountsWidget = widgets.NewObjectCountsWidget
	NewObjectListWidget   = widgets.NewObjectListWidget
	NewRSSFeedWidget      = widgets.NewRSSFeedWidget
	NewBookmarksWidget    = widgets.NewBookmarksWidget
)

// Export constants
var (
	BookmarkOrderingNewest = widgets.BookmarkOrderingNewest
	BookmarkOrderingOldest = widgets.BookmarkOrderingOldest
	BookmarkOrderingAZ     = widgets.BookmarkOrderingAZ
	BookmarkOrderingZA     = widgets.BookmarkOrderingZA
)
