// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/extras/enum"
)

// CustomFieldChoiceSet представляет набор выборов для пользовательского поля
type CustomFieldChoiceSet struct {
	ID                 string                `json:"id"`
	URL                string                `json:"url"`
	DisplayURL         string                `json:"display_url,omitempty"`
	Display            string                `json:"display"`
	Name               string                `json:"name"`
	Description        string                `json:"description,omitempty"`
	BaseChoices        enum.CustomFieldChoiceSetBaseChoices `json:"base_choices,omitempty"`
	ExtraChoices       [][]string            `json:"extra_choices"`
	OrderAlphabetically bool                 `json:"order_alphabetically"`
	ChoicesCount       int                   `json:"choices_count"`
	Owner              interface{}           `json:"owner,omitempty"`
	Created            *time.Time            `json:"created,omitempty"`
	LastUpdated        *time.Time            `json:"last_updated,omitempty"`
}

// CustomField представляет пользовательское поле
type CustomField struct {
	ID                  string                    `json:"id"`
	URL                 string                    `json:"url"`
	DisplayURL          string                    `json:"display_url,omitempty"`
	Display             string                    `json:"display"`
	ObjectTypes         []string                  `json:"object_types"`
	Type                enum.CustomFieldType      `json:"type"`
	RelatedObjectType   *string                   `json:"related_object_type,omitempty"`
	DataType            string                    `json:"data_type"`
	Name                string                    `json:"name"`
	Label               string                    `json:"label,omitempty"`
	GroupName           string                    `json:"group_name,omitempty"`
	Description         string                    `json:"description,omitempty"`
	Required            bool                      `json:"required"`
	Unique              bool                      `json:"unique"`
	SearchWeight        int                       `json:"search_weight"`
	FilterLogic         enum.CustomFieldFilterLogic `json:"filter_logic,omitempty"`
	UIVisible           enum.CustomFieldUIVisible `json:"ui_visible,omitempty"`
	UIEditable          enum.CustomFieldUIEditable `json:"ui_editable,omitempty"`
	IsCloneable         bool                      `json:"is_cloneable"`
	Default             interface{}               `json:"default,omitempty"`
	RelatedObjectFilter interface{}               `json:"related_object_filter,omitempty"`
	Weight              int                       `json:"weight"`
	ValidationMinimum   *int                      `json:"validation_minimum,omitempty"`
	ValidationMaximum   *int                      `json:"validation_maximum,omitempty"`
	ValidationRegex     string                    `json:"validation_regex,omitempty"`
	ChoiceSet           *CustomFieldChoiceSet     `json:"choice_set,omitempty"`
	Owner               interface{}               `json:"owner,omitempty"`
	Comments            string                    `json:"comments,omitempty"`
	Created             *time.Time                `json:"created,omitempty"`
	LastUpdated         *time.Time                `json:"last_updated,omitempty"`
}

// CustomFieldChoiceSetBrief краткое представление набора выборов
type CustomFieldChoiceSetBrief struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Display      string `json:"display"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	ChoicesCount int    `json:"choices_count"`
}

// CustomFieldBrief краткое представление пользовательского поля
type CustomFieldBrief struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Display     string `json:"display"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
