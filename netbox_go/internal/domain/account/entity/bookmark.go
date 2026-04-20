// Package entity contains Account domain entities for bookmarks
package entity

import (
	"errors"
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// Bookmark represents a user bookmark to an object or URL
// Simplified version of NetBox bookmark model used in account views
// Fields are intentionally minimal to support listing for current user
// without cross-app dependencies.
type Bookmark struct {
	ID      types.ID  `json:"id"`
	UserID  types.ID  `json:"user_id"`
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

// Validate checks bookmark fields
func (b *Bookmark) Validate() error {
	if b.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	if len(b.Title) == 0 {
		return errors.New("title is required")
	}
	if len(b.Title) > 200 {
		return errors.New("title too long (max 200 characters)")
	}
	if len(b.URL) == 0 {
		return errors.New("url is required")
	}
	return nil
}
