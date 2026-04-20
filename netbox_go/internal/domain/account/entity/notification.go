// Package entity contains Account domain entities for notifications
package entity

import (
	"errors"
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// Notification represents a user-facing notification entry
// Simplified projection of NetBox notification list for the current user.
type Notification struct {
	ID      types.ID  `json:"id"`
	UserID  types.ID  `json:"user_id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Level   string    `json:"level"` // info, warning, error
	Created time.Time `json:"created"`
	Read    bool      `json:"read"`
}

// Validate checks notification fields
func (n *Notification) Validate() error {
	if n.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	if len(n.Title) == 0 {
		return errors.New("title is required")
	}
	if len(n.Title) > 200 {
		return errors.New("title too long (max 200 characters)")
	}
	if len(n.Message) > 2000 {
		return errors.New("message too long (max 2000 characters)")
	}
	return nil
}
