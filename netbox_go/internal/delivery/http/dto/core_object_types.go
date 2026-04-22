package dto

import (
	"github.com/google/uuid"
)

// ObjectTypeResponse is the representation of an ObjectType
type ObjectTypeResponse struct {
	ID               uuid.UUID `json:"id"`
	AppLabel         string     `json:"app_label"`
	AppName          string     `json:"app_name"`
	Model            string     `json:"model"`
	ModelName        string     `json:"model_name"`
	ModelNamePlural  string     `json:"model_name_plural"`
	Public           bool       `json:"public"`
	Features         []string   `json:"features"`
	IsPluginModel    bool       `json:"is_plugin_model"`
	RestAPIEndpoint  string     `json:"rest_api_endpoint"`
	Description      string     `json:"description"`
}

// ObjectTypeBriefResponse is a condensed version for lists
type ObjectTypeBriefResponse struct {
	ID       uuid.UUID `json:"id"`
	AppLabel string    `json:"app_label"`
	Model    string    `json:"model"`
}
