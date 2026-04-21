// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"
)

// Bookmark представляет закладку
type Bookmark struct {
	ID         string      `json:"id"`
	URL        string      `json:"url"`
	Display    string      `json:"display"`
	ObjectType string      `json:"object_type"`
	ObjectID   string      `json:"object_id"`
	Object     interface{} `json:"object,omitempty"`
	User       *User       `json:"user"`
	Created    *time.Time  `json:"created,omitempty"`
}

// BookmarkBrief краткое представление закладки
type BookmarkBrief struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Display    string `json:"display"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
}

// User представляет пользователя
type User struct {
	ID       string `json:"id"`
	URL      string `json:"url,omitempty"`
	Display  string `json:"display,omitempty"`
	Username string `json:"username,omitempty"`
}
