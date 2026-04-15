package types

import "errors"

// Ошибки валидации
var (
	ErrInvalidLatitude  = errors.New("invalid latitude: must be between -90.0 and 90.0")
	ErrInvalidLongitude = errors.New("invalid longitude: must be between -180.0 and 180.0")
	ErrInvalidSlug      = errors.New("invalid slug: must be between 1 and 100 characters")
	ErrInvalidStatus    = errors.New("invalid status value")
)

// Ошибки бизнес-логики
var (
	ErrNotFound           = errors.New("entity not found")
	ErrAlreadyExists      = errors.New("entity already exists")
	ErrValidationFailed   = errors.New("validation failed")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrConstraintViolation = errors.New("constraint violation")
)
