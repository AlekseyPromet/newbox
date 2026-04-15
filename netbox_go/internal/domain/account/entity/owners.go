// Package entity содержит сущности домена Account (продолжение)
package entity

import (
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// OwnerGroup представляет группу владельцев объектов
type OwnerGroup struct {
	ID          types.ID  `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Validate проверяет корректность группы владельцев
func (og *OwnerGroup) Validate() error {
	if len(og.Name) == 0 || len(og.Name) > 100 {
		return &ValidationError{"name is required and must be <= 100 characters"}
	}
	if len(og.Description) > 200 {
		return &ValidationError{"description too long (max 200 characters)"}
	}
	return nil
}

// Owner представляет владельца объекта (пользователь или группа)
type Owner struct {
	ID          types.ID   `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	GroupID     *types.ID  `json:"group_id,omitempty"`
	UserIDs     []types.ID `json:"user_ids,omitempty"`
	GroupIDs    []types.ID `json:"group_ids,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность владельца
func (o *Owner) Validate() error {
	if len(o.Name) == 0 || len(o.Name) > 100 {
		return &ValidationError{"name is required and must be <= 100 characters"}
	}
	if len(o.Description) > 200 {
		return &ValidationError{"description too long (max 200 characters)"}
	}
	return nil
}

// IsGroupOwner проверяет, является ли владелец групповым
func (o *Owner) IsGroupOwner() bool {
	return o.GroupID != nil
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
