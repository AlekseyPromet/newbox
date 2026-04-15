// Package types содержит ошибки бизнес-логики для всего приложения
package types

import (
	"errors"
	"fmt"
)

// Ошибки бизнес-логики
var (
	ErrNotFound            = errors.New("entity not found")
	ErrAlreadyExists       = errors.New("entity already exists")
	ErrValidationFailed    = errors.New("validation failed")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrInvalidOperation    = errors.New("invalid operation")
	ErrConstraintViolation = errors.New("constraint violation")
)

// ValidationError представляет ошибку валидации поля
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
