// Package types содержит ошибки бизнес-логики для всего приложения
package types

import "errors"

// Ошибки бизнес-логики
var (
	ErrNotFound            = errors.New("entity not found")
	ErrAlreadyExists       = errors.New("entity already exists")
	ErrValidationFailed    = errors.New("validation failed")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrInvalidOperation    = errors.New("invalid operation")
	ErrConstraintViolation = errors.New("constraint violation")
)
