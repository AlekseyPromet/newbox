package dto

import (
	"time"

	"github.com/google/uuid"
)

// ObjectChangeResponse is the representation of an ObjectChange (Change Log entry)
type ObjectChangeResponse struct {
	ID                uuid.UUID              `json:"id"`
	Time              time.Time              `json:"time"`
	User              *UserBriefResponse     `json:"user"`
	UserName          string                 `json:"user_name"`
	RequestID         string                 `json:"request_id"`
	Action            string                 `json:"action"` // 'create', 'update', 'delete'
	ChangedObjectType string                 `json:"changed_object_type"`
	ChangedObjectID   string                 `json:"changed_object_id"`
	ChangedObject     interface{}            `json:"changed_object"` // Dynamic representation
	ObjectRepr        string                 `json:"object_repr"`
	Message           string                 `json:"message"`
	PrechangeData     map[string]interface{} `json:"prechange_data"`
	PostchangeData    map[string]interface{} `json:"postchange_data"`
}

// ObjectChangeBriefResponse is a condensed version for lists
type ObjectChangeBriefResponse struct {
	ID                uuid.UUID `json:"id"`
	Time              time.Time `json:"time"`
	User              string    `json:"user"`
	Action            string    `json:"action"`
	ChangedObjectType string    `json:"changed_object_type"`
	ObjectRepr        string    `json:"object_repr"`
}
