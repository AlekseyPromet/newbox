// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"
)

// ConfigTemplate представляет шаблон конфигурации
type ConfigTemplate struct {
	ID               string      `json:"id"`
	URL              string      `json:"url"`
	DisplayURL       string      `json:"display_url,omitempty"`
	Display          string      `json:"display"`
	Name             string      `json:"name"`
	Description      string      `json:"description,omitempty"`
	EnvironmentParams interface{} `json:"environment_params,omitempty"`
	TemplateCode     string      `json:"template_code"`
	MimeType         string      `json:"mime_type,omitempty"`
	FileName         string      `json:"file_name,omitempty"`
	FileExtension    string      `json:"file_extension,omitempty"`
	AsAttachment     bool        `json:"as_attachment"`
	DataSource       *DataSource `json:"data_source,omitempty"`
	DataPath         string      `json:"data_path,omitempty"`
	DataFile         *DataFile   `json:"data_file,omitempty"`
	AutoSyncEnabled  bool        `json:"auto_sync_enabled"`
	DataSynced       *time.Time  `json:"data_synced,omitempty"`
	Owner            interface{} `json:"owner,omitempty"`
	Tags             []Tag       `json:"tags,omitempty"`
	Created          *time.Time  `json:"created,omitempty"`
	LastUpdated      *time.Time  `json:"last_updated,omitempty"`
}

// ConfigTemplateBrief краткое представление шаблона конфигурации
type ConfigTemplateBrief struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Display     string `json:"display"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
