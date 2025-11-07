package apperror

import "errors"

// Standard application errors.
var (
	// ErrNotFound indicates a requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput indicates invalid input from the user.
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict indicates a resource conflict, like a duplicate key.
	ErrConflict = errors.New("resource conflict")

	// ErrInternal is a generic fallback for server-side errors.
	ErrInternal = errors.New("internal error")
)
