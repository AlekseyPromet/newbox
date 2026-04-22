// Package entity contains Account domain entities for user configuration
package entity

import (
	"encoding/json"
	"errors"

	"netbox_go/pkg/types"
)

// UserConfig represents user preferences/config stored as JSON blob
// This is a simplified analogue of NetBox UserConfig (users.preferences)
type UserConfig struct {
	UserID types.ID        `json:"user_id"`
	Data   json.RawMessage `json:"data"`
}

// Validate checks basic constraints
func (uc *UserConfig) Validate() error {
	if uc.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	// Data can be empty JSON object; schema validation is out of scope here
	return nil
}
