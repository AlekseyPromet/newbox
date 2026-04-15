// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"
)

// ImageAttachment представляет вложение изображения
type ImageAttachment struct {
	ID          string     `json:"id"`
	URL         string     `json:"url"`
	Display     string     `json:"display"`
	ObjectType  string     `json:"object_type"`
	ObjectID    string     `json:"object_id"`
	Parent      interface{} `json:"parent,omitempty"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Description string     `json:"description,omitempty"`
	ImageHeight int        `json:"image_height,omitempty"`
	ImageWidth  int        `json:"image_width,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
}

// ImageAttachmentBrief краткое представление вложения изображения
type ImageAttachmentBrief struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Display     string `json:"display"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description,omitempty"`
}
