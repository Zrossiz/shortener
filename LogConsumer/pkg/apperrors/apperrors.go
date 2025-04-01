package apperrors

import "errors"

const (
	ErrInternalServer     = "internal server error"
	ErrInvalidRequestBody = "invalid request body"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)
