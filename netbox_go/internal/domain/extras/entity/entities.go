package entity

import (
	"time"

	"github.com/google/uuid"
)

// GenericForeignKey represents a reference to any object in the system.
type GenericForeignKey struct {
	ContentType string `json:"content_type"`
	ObjectID    int64  `json:"object_id"`
}

type EventRule struct {
	ID                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	EventTypes        []string          `json:"event_types"`
	Enabled           bool              `json:"enabled"`
	Conditions        map[string]any    `json:"conditions,omitempty"`
	ActionType        string            `json:"action_type"`
	ActionObject      GenericForeignKey `json:"action_object"`
	ActionData        map[string]any    `json:"action_data,omitempty"`
	Comments          string            `json:"comments,omitempty"`
	OwnerID           *uuid.UUID        `json:"owner_id,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type Webhook struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description,omitempty"`
	PayloadURL        string    `json:"payload_url"`
	HTTPMethod        string    `json:"http_method"`
	HTTPContentType   string    `json:"http_content_type"`
	AdditionalHeaders string    `json:"additional_headers,omitempty"`
	BodyTemplate      string    `json:"body_template,omitempty"`
	Secret            string    `json:"secret,omitempty"`
	SSLVerification   bool      `json:"ssl_verification"`
	CAFilePath        string    `json:"ca_file_path,omitempty"`
	OwnerID           *uuid.UUID `json:"owner_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CustomLink struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Enabled      bool      `json:"enabled"`
	LinkText     string    `json:"link_text"`
	LinkURL      string    `json:"link_url"`
	Weight       int       `json:"weight"`
	GroupName    string    `json:"group_name,omitempty"`
	ButtonClass  string    `json:"button_class"`
	NewWindow    bool      `json:"new_window"`
	OwnerID      *uuid.UUID `json:"owner_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ExportTemplate struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	TemplateCode   string    `json:"template_code"`
	MimeType       string    `json:"mime_type"`
	FileName       string    `json:"file_name"`
	FileExtension  string    `json:"file_extension"`
	AsAttachment   bool      `json:"as_attachment"`
	OwnerID        *uuid.UUID `json:"owner_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SavedFilter struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description,omitempty"`
	UserID      *uuid.UUID    `json:"user_id,omitempty"`
	Weight      int           `json:"weight"`
	Enabled     bool          `json:"enabled"`
	Shared      bool          `json:"shared"`
	Parameters  map[string]any `json:"parameters"`
	OwnerID     *uuid.UUID    `json:"owner_id,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type TableConfig struct {
	ID          uuid.UUID `json:"id"`
	ObjectType  string    `json:"object_type"`
	Table       string    `json:"table"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	Weight      int       `json:"weight"`
	Enabled     bool      `json:"enabled"`
	Shared      bool      `json:"shared"`
	Columns     []string  `json:"columns"`
	Ordering    []string  `json:"ordering,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ImageAttachment struct {
	ID          uuid.UUID         `json:"id"`
	Object       GenericForeignKey `json:"object"`
	Image       string            `json:"image"`
	ImageHeight  int               `json:"image_height"`
	ImageWidth   int               `json:"image_width"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type JournalEntry struct {
	ID                uuid.UUID         `json:"id"`
	AssignedObject    GenericForeignKey `json:"assigned_object"`
	CreatedBy         *uuid.UUID        `json:"created_by,omitempty"`
	Kind              string            `json:"kind"`
	Comments          string            `json:"comments"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type Bookmark struct {
	ID        uuid.UUID         `json:"id"`
	Object    GenericForeignKey `json:"object"`
	UserID    uuid.UUID         `json:"user_id"`
	CreatedAt time.Time         `json:"created_at"`
}
