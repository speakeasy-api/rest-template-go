package errors

import (
	"errors"
	"faceittechtest/internal/app"
)

const (
	// ErrUnknown is returned when an unexpected error occurs
	ErrUnknown = app.Error("err_unknown: unknown error occured")
	// ErrInvalidRequest is returned when either the parameters or the request body is invalid
	ErrInvalidRequest = app.Error("err_invalid_request: invalid request received")
	// ErrValidation is returned when the parameters don't pass validation
	ErrValidation = app.Error("err_validation: failed validation")
	// ErrNotFound is returned when the requested resource is not found
	ErrNotFound = app.Error("err_not_found: not found")
)

// Is just wraps errors.Is as we don't want to alias the errors package everywhere to use it
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As just wraps errors.As as we don't want to alias the errors package everywhere to use it
func As(err error, target any) bool {
	return errors.As(err, target)
}
