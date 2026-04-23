// Package types содержит общие типы данных и ошибки для всего приложения
package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Ошибки валидации и бизнес-логики
var (
	ErrInvalidID              = errors.New("invalid ID")
	ErrInvalidStatus          = errors.New("invalid status")
	ErrInvalidSlug            = errors.New("invalid slug")
	ErrInvalidLatitude        = errors.New("invalid latitude: must be between -90 and 90")
	ErrInvalidLongitude       = errors.New("invalid longitude: must be between -180 and 180")
	ErrNameRequired           = errors.New("name is required")
	ErrDeviceNameRequired     = errors.New("device name is required")
	ErrDeviceModelRequired    = errors.New("device model is required")
	ErrManufacturerRequired   = errors.New("manufacturer is required")
	ErrPlatformNameRequired   = errors.New("platform name is required")
	ErrDeviceRoleNameRequired = errors.New("device role name is required")
	ErrColorRequired          = errors.New("color is required")
	ErrInvalidUHeight         = errors.New("invalid u_height: must be non-negative")
	ErrDeviceTypeRequired     = errors.New("device type is required")
	ErrDeviceRoleRequired     = errors.New("device role is required")
	ErrSiteRequired           = errors.New("site is required")
	ErrInvalidRackPosition    = errors.New("invalid rack position")
	ErrModuleModelRequired    = errors.New("module model is required")
	ErrModuleBayRequired      = errors.New("module bay is required")
	ErrModuleTypeRequired     = errors.New("module type is required")
	ErrDeviceRequired         = errors.New("device is required")
)

// ID представляет идентификатор сущности
type ID uuid.UUID

// NewID генерирует новый UUID
func NewID() ID {
	return ID(uuid.New())
}

func (id ID) IsNuul() bool {
	return id.String() != ""
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
	Created   time.Time  `json:"created"`
	Updated   time.Time  `json:"updated"`
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
	ID         ID        `json:"id"`
	Name       string    `json:"name"`
	Image      []byte    `json:"image"`
	Uploaded   time.Time `json:"uploaded"`
	UploadedBy *ID       `json:"uploaded_by,omitempty"`
}

// Contact представляет контакт
type Contact struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

// Tenant представляет арендатора
type Tenant struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
	Slug Slug   `json:"slug"`
}

// ASN представляет автономную систему
type ASN struct {
	ID   ID     `json:"id"`
	ASN  uint32 `json:"asn"`
	Name string `json:"name"`
	Slug Slug   `json:"slug"`
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

// StringArray представляет массив строк для работы с PostgreSQL array
type StringArray struct {
	Elements []string
}

// Scan реализует sql.Scanner для StringArray
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		a.Elements = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	str := string(bytes)
	// Простая реализация парсинга PostgreSQL array формата {elem1,elem2}
	if str == "{}" {
		a.Elements = []string{}
		return nil
	}
	if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
		return errors.New("invalid array format")
	}
	content := str[1 : len(str)-1]
	if content == "" {
		a.Elements = []string{}
		return nil
	}
	a.Elements = splitArray(content)
	return nil
}

// Value реализует driver.Valuer для StringArray
func (a StringArray) Value() ([]byte, error) {
	if a.Elements == nil {
		return nil, nil
	}
	if len(a.Elements) == 0 {
		return []byte("{}"), nil
	}
	result := "{"
	for i, elem := range a.Elements {
		if i > 0 {
			result += ","
		}
		result += elem
	}
	result += "}"
	return []byte(result), nil
}

func splitArray(s string) []string {
	var result []string
	current := ""
	inQuotes := false
	for _, c := range s {
		switch c {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if inQuotes {
				current += string(c)
			} else {
				result = append(result, current)
				current = ""
			}
		default:
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func MarshalJSON(v any) (data []byte, e error) {
	return json.Marshal(v)
}
