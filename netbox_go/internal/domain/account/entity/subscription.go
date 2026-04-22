// Package entity contains Account domain entities for subscriptions
package entity

import (
	"errors"
	"time"

	"netbox_go/pkg/types"
)

// Subscription represents a user subscription to object changes/notifications
// Simplified projection to support listing for current user.
type Subscription struct {
	ID         types.ID  `json:"id"`
	UserID     types.ID  `json:"user_id"`
	ObjectID   types.ID  `json:"object_id"`
	ObjectType string    `json:"object_type"`
	Created    time.Time `json:"created"`
}

// Validate checks subscription fields
func (s *Subscription) Validate() error {
	if s.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	if s.ObjectID.String() == "" {
		return errors.New("object_id is required")
	}
	if len(s.ObjectType) == 0 {
		return errors.New("object_type is required")
	}
	if len(s.ObjectType) > 100 {
		return errors.New("object_type too long (max 100 characters)")
	}
	return nil
}
