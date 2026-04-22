// Package entity содержит сущности домена Users для управления ролями и правами
package entity

import (
	"errors"
	"time"

	"netbox_go/pkg/types"
)

// Role представляет роль пользователя в системе
type Role struct {
	ID          types.ID   `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Permissions []string   `json:"permissions,omitempty"` // Список кодов разрешений
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность роли
func (r *Role) Validate() error {
	if len(r.Name) == 0 || len(r.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if len(r.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}
	return nil
}

// HasPermission проверяет, есть ли у роли указанное разрешение
func (r *Role) HasPermission(permissionCode string) bool {
	for _, perm := range r.Permissions {
		if perm == permissionCode {
			return true
		}
	}
	return false
}

// AddPermission добавляет разрешение к роли (если ещё не добавлено)
func (r *Role) AddPermission(permissionCode string) {
	if !r.HasPermission(permissionCode) {
		r.Permissions = append(r.Permissions, permissionCode)
	}
}

// RemovePermission удаляет разрешение из роли
func (r *Role) RemovePermission(permissionCode string) {
	for i, perm := range r.Permissions {
		if perm == permissionCode {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			return
		}
	}
}

// Permission представляет разрешение на выполнение действия
type Permission struct {
	ID          types.ID  `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"` // Уникальный код разрешения, например "job.create", "job.assign"
	Description string    `json:"description,omitempty"`
	ObjectType  string    `json:"object_type,omitempty"` // Тип объекта, к которому относится разрешение
	Action      string    `json:"action"`                // Действие: create, read, update, delete, assign
	Created     time.Time `json:"created"`
}

// Validate проверяет корректность разрешения
func (p *Permission) Validate() error {
	if len(p.Name) == 0 || len(p.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if p.Code == "" {
		return errors.New("code is required")
	}
	if p.Action == "" {
		return errors.New("action is required")
	}
	if len(p.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}
	return nil
}

// UserRole связывает пользователя с ролью
type UserRole struct {
	ID        types.ID  `json:"id"`
	UserID    types.ID  `json:"user_id"`
	RoleID    types.ID  `json:"role_id"`
	AssignedBy *types.ID `json:"assigned_by,omitempty"` // Кто назначил роль
	Created   time.Time `json:"created"`
	Expires   *time.Time `json:"expires,omitempty"`     // Срок действия роли (опционально)
}

// Validate проверяет корректность связи пользователь-роль
func (ur *UserRole) Validate() error {
	if ur.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	if ur.RoleID.String() == "" {
		return errors.New("role_id is required")
	}
	return nil
}

// IsExpired проверяет, истек ли срок действия роли
func (ur *UserRole) IsExpired() bool {
	if ur.Expires == nil {
		return false
	}
	return time.Now().After(*ur.Expires)
}

// IsActive проверяет, активна ли роль у пользователя
func (ur *UserRole) IsActive() bool {
	return !ur.IsExpired()
}

// JobAssignment представляет назначение задачи ответственному (пользователю или группе)
type JobAssignment struct {
	ID           types.ID   `json:"id"`
	JobID        types.ID   `json:"job_id"`
	AssigneeType string     `json:"assignee_type"` // "user" или "group"
	UserID       *types.ID  `json:"user_id,omitempty"`
	GroupID      *types.ID  `json:"group_id,omitempty"`
	RoleID       *types.ID  `json:"role_id,omitempty"` // Роль, определяющая права ответственного
	AssignedBy   types.ID   `json:"assigned_by"`
	AssignedAt   time.Time  `json:"assigned_at"`
	Comments     string     `json:"comments,omitempty"`
}

// Validate проверяет корректность назначения задачи
func (ja *JobAssignment) Validate() error {
	if ja.JobID.String() == "" {
		return errors.New("job_id is required")
	}
	if ja.AssigneeType != "user" && ja.AssigneeType != "group" {
		return errors.New("assignee_type must be 'user' or 'group'")
	}
	if ja.AssigneeType == "user" && ja.UserID == nil {
		return errors.New("user_id is required for user assignee")
	}
	if ja.AssigneeType == "group" && ja.GroupID == nil {
		return errors.New("group_id is required for group assignee")
	}
	if ja.AssignedBy.String() == "" {
		return errors.New("assigned_by is required")
	}
	return nil
}

// GetAssigneeID возвращает ID ответственного (пользователя или группы)
func (ja *JobAssignment) GetAssigneeID() *types.ID {
	if ja.AssigneeType == "user" {
		return ja.UserID
	}
	return ja.GroupID
}
