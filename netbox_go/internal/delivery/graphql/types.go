package graphql

import (
	"netbox_go/internal/domain/core/entity"
)

// DataFileType represents the GraphQL type for DataFile
type DataFileType struct {
	ID          string           `json:"id"`
	Path        string           `json:"path"`
	Size        int64            `json:"size"`
	Hash        string           `json:"hash"`
	Created     string           `json:"created"`
	LastUpdated string           `json:"last_updated"`
	Source      *DataSourceType  `json:"source"`
}

// DataSourceType represents the GraphQL type for DataSource
type DataSourceType struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	SourceURL   string         `json:"source_url"`
	Status      string         `json:"status"`
	Enabled     bool           `json:"enabled"`
	IgnoreRules string         `json:"ignore_rules"`
	Parameters  interface{}    `json:"parameters"`
	LastSynced  string         `json:"last_synced"`
	DataFiles   []*DataFileType `json:"datafiles"`
}

// ObjectChangeType represents the GraphQL type for ObjectChange
type ObjectChangeType struct {
	ID                  string `json:"id"`
	Time                string `json:"time"`
	UserName            string `json:"user_name"`
	RequestID           string `json:"request_id"`
	Action              string `json:"action"`
	ChangedObjectTypeID string `json:"changed_object_type_id"`
	ChangedObjectID     string `json:"changed_object_id"`
	RelatedObjectTypeID string `json:"related_object_type_id"`
	RelatedObjectID     string `json:"related_object_id"`
	ObjectRepr          string `json:"object_repr"`
	PrechangeData       string `json:"prechange_data"`
	PostchangeData      string `json:"postchange_data"`
}

// ContentType represents the GraphQL type for ContentType
type ContentType struct {
	ID        string `json:"id"`
	AppLabel  string `json:"app_label"`
	Model     string `json:"model"`
}
