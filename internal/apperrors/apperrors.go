package apperrors

import "errors"

const (
	ErrInternalServer     = "internal server error"
	ErrInvalidRequestBody = "invalid request body"
)

var (
	ErrRedisNotFound = errors.New("key not found")
)
