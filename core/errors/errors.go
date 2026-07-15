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
//
// Use the constructor helpers (ErrNotFound, ErrBadRequest, etc.) to build
// errors, and handle them at the transport layer with Respond.
package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// Core error type
// ---------------------------------------------------------------------------

// AppError is the canonical application error.
// It wraps an underlying cause and carries an HTTP status code, a machine-
// readable code, a human-readable message, and optional structured details.
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

// Unwrap allows errors.Is / errors.As to traverse the chain.
func (e *AppError) Unwrap() error { return e.cause }

// Wrap attaches an underlying error to an AppError and returns the same
// pointer so calls can be chained.
func (e *AppError) Wrap(cause error) *AppError {
	e.cause = cause
	return e
}

// WithDetails attaches arbitrary structured detail to the error.
func (e *AppError) WithDetails(d any) *AppError {
	e.Details = d
	return e
}

// ---------------------------------------------------------------------------
// Constructor helpers
// ---------------------------------------------------------------------------

// New creates a bare AppError. Prefer the named constructors below.
func New(status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
	}
}

// ErrBadRequest creates a 400 Bad Request error.
func ErrBadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

// ErrUnauthorized creates a 401 Unauthorized error.
func ErrUnauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// ErrForbidden creates a 403 Forbidden error.
func ErrForbidden(message string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message)
}

// ErrNotFound creates a 404 Not Found error.
func ErrNotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

// ErrConflict creates a 409 Conflict error.
func ErrConflict(message string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message)
}

// ErrUnprocessable creates a 422 Unprocessable Entity error,
// typically used for validation failures.
func ErrUnprocessable(message string) *AppError {
	return New(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message)
}

// ErrInternalServer creates a 500 Internal Server Error.
func ErrInternalServer(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

// ---------------------------------------------------------------------------
// Inspection helpers
// ---------------------------------------------------------------------------

// IsNotFound reports whether err (or any error in its chain) is a NOT_FOUND
// AppError.
func IsNotFound(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusNotFound
}

// IsBadRequest reports whether err is a BAD_REQUEST AppError.
func IsBadRequest(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusBadRequest
}

// ---------------------------------------------------------------------------
// Gin response helper
// ---------------------------------------------------------------------------

// Respond writes the error as a JSON response to the Gin context.
// If err is an *AppError it uses its status; otherwise it falls back to 500.
func Respond(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Unknown error – hide the internals from the client
	c.JSON(http.StatusInternalServerError, &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "an unexpected error occurred",
	})
}
