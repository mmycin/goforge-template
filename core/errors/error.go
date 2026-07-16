// Package errors provides standardized error types and HTTP error responses
// for GoForge applications.
//
// All API errors share a common JSON shape:
//
//	{
//	  "code":    "NOT_FOUND",
//	  "message": "todo not found",
//	  "details": { ... }   // optional, omitted when nil
//	}
package errors

import "fmt"

// AppError is the canonical application error.
// It carries an HTTP status code, a machine-readable code,
// a human-readable message, an optional structured details payload,
// and an optional underlying cause.
type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	cause      error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows errors.Is / errors.As to traverse the error chain.
func (e *AppError) Unwrap() error { return e.cause }

// Wrap attaches an underlying error to this AppError and returns the same
// pointer so calls can be chained:
//
//	apperrors.ErrNotFound("todo not found").Wrap(err)
func (e *AppError) Wrap(cause error) *AppError {
	e.cause = cause
	return e
}

// WithDetails attaches arbitrary structured detail to the error.
// Useful for validation failure payloads.
func (e *AppError) WithDetails(d any) *AppError {
	e.Details = d
	return e
}

// New creates a bare AppError. Prefer the named constructors in constructors.go.
func New(status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
	}
}
