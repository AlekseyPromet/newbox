// Package types содержит общие типы данных для всего приложения
package types

import (
	"time"

	"github.com/google/uuid"
)

// ID представляет идентификатор сущности
type ID uuid.UUID

// NewID генерирует новый UUID
func NewID() ID {
	return ID(uuid.New())
}

// ParseID парсит строку в ID
func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID(id), nil
}

// String возвращает строковое представление ID
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// Status представляет статус сущности
type Status string

// TimeStamp представляет метки времени создания и обновления
type TimeStamp struct {
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AuditInfo содержит информацию об аудите
type AuditInfo struct {
	CreatedBy   *ID       `json:"created_by,omitempty"`
	UpdatedBy   *ID       `json:"updated_by,omitempty"`
	DeletedBy   *ID       `json:"deleted_by,omitempty"`
	LastChanged time.Time `json:"last_changed"`
	ChangeID    *string   `json:"change_id,omitempty"`
}

// Slug представляет slug-строку для URL
type Slug string

// Validate проверяет корректность slug
func (s Slug) Validate() error {
	// Реализация валидации slug
	if len(s) == 0 || len(s) > 100 {
		return ErrInvalidSlug
	}
	return nil
}

// Description представляет описание сущности
type Description string

// Comments представляет комментарии
type Comments string

// Image представляет изображение
type Image struct {
	ID        ID       `json:"id"`
	Name      string   `json:"name"`
	Image     []byte   `json:"image"`
	Uploaded  time.Time `json:"uploaded"`
	UploadedBy *ID      `json:"uploaded_by,omitempty"`
}

// Contact представляет контакт
type Contact struct {
	ID   ID   `json:"id"`
	Name string `json:"name"`
}

// Tenant представляет арендатора
type Tenant struct {
	ID   ID   `json:"id"`
	Name string `json:"name"`
	Slug Slug `json:"slug"`
}

// ASN представляет автономную систему
type ASN struct {
	ID    ID    `json:"id"`
	ASN   uint32 `json:"asn"`
	Name  string `json:"name"`
	Slug  Slug  `json:"slug"`
}

// Coordinate представляет GPS координаты
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Validate проверяет корректность координат
func (c Coordinate) Validate() error {
	if c.Latitude < -90.0 || c.Latitude > 90.0 {
		return ErrInvalidLatitude
	}
	if c.Longitude < -180.0 || c.Longitude > 180.0 {
		return ErrInvalidLongitude
	}
	return nil
}

// Address представляет физический или почтовый адрес
type Address string

// TimeZone представляет часовой пояс
type TimeZone string

// Facility представляет идентификатор или описание объекта
type Facility string
