package dto

import (
	"time"

	"github.com/google/uuid"
)

// DataSourceRequest is used for creating and updating DataSources
type DataSourceRequest struct {
	Name         string                 `json:"name" validate:"required,max=100"`
	Type         string                 `json:"type" validate:"required"`
	SourceURL    string                 `json:"source_url" validate:"required,max=500"`
	Enabled      *bool                  `json:"enabled"`
	SyncInterval *int                   `json:"sync_interval"`
	IgnoreRules  []string               `json:"ignore_rules"`
	Parameters   map[string]interface{} `json:"parameters"`
	Description  string                 `json:"description"`
}

// DataSourceResponse is the full representation of a DataSource
type DataSourceResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	SourceURL    string                 `json:"source_url"`
	Enabled      bool                   `json:"enabled"`
	Status       string                 `json:"status"`
	Description  string                 `json:"description"`
	SyncInterval int                    `json:"sync_interval"`
	Parameters   map[string]interface{} `json:"parameters"`
	IgnoreRules  []string               `json:"ignore_rules"`
	LastSynced   *time.Time             `json:"last_synced"`
	Created      time.Time              `json:"created"`
	Updated      time.Time              `json:"updated"`
	FileCount    int                    `json:"file_count"`
}

// DataSourceBriefResponse is a condensed version for lists
type DataSourceBriefResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// DataFileResponse is the representation of a DataFile
type DataFileResponse struct {
	ID        uuid.UUID           `json:"id"`
	Source    *DataSourceResponse `json:"source"`
	Path      string              `json:"path"`
	Size      int64               `json:"size"`
	Hash      string              `json:"hash"`
	Updated   time.Time           `json:"updated"`
	Created   time.Time           `json:"created"`
}

// DataFileBriefResponse is a condensed version for lists
type DataFileBriefResponse struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"path"`
}
