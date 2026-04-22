package entity

import (
	"time"

	"github.com/google/uuid"
)

// EventRule defines an action to be taken automatically in response to a specific set of events.
type EventRule struct {
	ID                uuid.UUID
	Name              string
	Description       string
	EventTypes        []string
	Enabled           bool
	Conditions        map[string]interface{}
	ActionType        string
	ActionObjectType  string
	ActionObjectID    *int64
	ActionData        map[string]interface{}
	Comments          string
	OwnerID           uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ObjectTypes       []string // Many-to-Many relationship
}

// Webhook defines a request that will be sent to a remote application.
type Webhook struct {
	ID               uuid.UUID
	Name             string
	Description      string
	PayloadURL       string
	HTTPMethod       string
	HTTPContentType  string
	AdditionalHeaders string
	BodyTemplate     string
	Secret           string
	SSLVerification  bool
	CAFilePath       string
	OwnerID          uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CustomLink is a custom link to an external representation of a NetBox object.
type CustomLink struct {
	ID           uuid.UUID
	Name         string
	Enabled      bool
	LinkText     string
	LinkURL      string
	Weight       int16
	GroupName    string
	ButtonClass  string
	NewWindow    bool
	OwnerID      uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ObjectTypes  []string // Many-to-Many relationship
}

// ExportTemplate is a template for exporting object data.
type ExportTemplate struct {
	ID            uuid.UUID
	Name          string
	Description   string
	TemplateCode  string
	MimeType      string
	FileName      string
	FileExtension string
	AsAttachment  bool
	OwnerID       uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ObjectTypes   []string // Many-to-Many relationship
}

// SavedFilter is a set of predefined keyword parameters for filtering.
type SavedFilter struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	UserID      uuid.UUID
	Weight      int16
	Enabled     bool
	Shared      bool
	Parameters  map[string]interface{}
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ObjectTypes []string // Many-to-Many relationship
}

// TableConfig is a saved configuration of columns and ordering for a table.
type TableConfig struct {
	ID          uuid.UUID
	ObjectType  string
	Table       string
	Name        string
	Description string
	UserID      uuid.UUID
	Weight      int16
	Enabled     bool
	Shared      bool
	Columns     []string
	Ordering    []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ImageAttachment is an uploaded image associated with an object.
type ImageAttachment struct {
	ID          uuid.UUID
	ObjectType  string
	ObjectID    int64
	Image       string
	ImageHeight int16
	ImageWidth  int16
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// JournalEntry is a historical remark concerning an object.
type JournalEntry struct {
	ID                uuid.UUID
	AssignedObjectType string
	AssignedObjectID   int64
	CreatedBy         uuid.UUID
	Kind              string
	Comments          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Bookmark is an object bookmarked by a user.
type Bookmark struct {
	ID         uuid.UUID
	ObjectType string
	ObjectID   int64
	UserID     uuid.UUID
	CreatedAt  time.Time
}
