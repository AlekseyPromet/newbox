package graphql

import (
	"time"
)

// DataFileFilter defines filters for DataFile
type DataFileFilter struct {
	Created     *TimeFilter `json:"created"`
	LastUpdated *TimeFilter `json:"last_updated"`
	SourceID    *string     `json:"source_id"`
	Path        *StringFilter `json:"path"`
	Size        *IntFilter    `json:"size"`
	Hash        *StringFilter `json:"hash"`
}

// DataSourceFilter defines filters for DataSource
type DataSourceFilter struct {
	Name        *StringFilter `json:"name"`
	Type        *StringFilter `json:"type"`
	SourceURL   *StringFilter `json:"source_url"`
	Status      *string       `json:"status"`
	Enabled     *bool         `json:"enabled"`
	IgnoreRules *StringFilter `json:"ignore_rules"`
	LastSynced  *TimeFilter   `json:"last_synced"`
}

// ObjectChangeFilter defines filters for ObjectChange
type ObjectChangeFilter struct {
	Time                *TimeFilter   `json:"time"`
	UserName            *StringFilter `json:"user_name"`
	RequestID           *StringFilter `json:"request_id"`
	Action              *string       `json:"action"`
	ChangedObjectTypeID *string       `json:"changed_object_type_id"`
	ChangedObjectID     *string       `json:"changed_object_id"`
	RelatedObjectTypeID *string       `json:"related_object_type_id"`
	RelatedObjectID     *string       `json:"related_object_id"`
	ObjectRepr          *StringFilter `json:"object_repr"`
}

// ContentTypeFilter defines filters for ContentType
type ContentTypeFilter struct {
	AppLabel *StringFilter `json:"app_label"`
	Model    *StringFilter `json:"model"`
}

// Generic Filter Lookups

type StringFilter struct {
	Exact      *string `json:"exact"`
	Contains   *string `json:"contains"`
	IContains  *string `json:"icontains"`
	StartsWith *string `json:"startswith"`
	EndsWith   *string `json:"endswith"`
}

type IntFilter struct {
	Exact    *int `json:"exact"`
	Gt       *int `json:"gt"`
	Gte      *int `json:"gte"`
	Lt       *int `json:"lt"`
	Lte      *int `json:"lte"`
}

type TimeFilter struct {
	Exact    *time.Time `json:"exact"`
	Gt       *time.Time `json:"gt"`
	Gte      *time.Time `json:"gte"`
	Lt       *time.Time `json:"lt"`
	Lte      *time.Time `json:"lte"`
}
