// Package serializers содержит DTO и функции сериализации для API
package serializers

// Dashboard представляет панель управления
type Dashboard struct {
	Layout interface{} `json:"layout"`
	Config interface{} `json:"config"`
}
