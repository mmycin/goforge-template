package errors

import "net/http"

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

// ErrUnprocessable creates a 422 Unprocessable Entity error.
// Typically used for validation failures.
func ErrUnprocessable(message string) *AppError {
	return New(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message)
}

// ErrTooManyRequests creates a 429 Too Many Requests error.
func ErrTooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

// ErrInternalServer creates a 500 Internal Server Error.
func ErrInternalServer(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

// ErrNotImplemented creates a 501 Not Implemented error.
func ErrNotImplemented(message string) *AppError {
	return New(http.StatusNotImplemented, "NOT_IMPLEMENTED", message)
}

// ErrServiceUnavailable creates a 503 Service Unavailable error.
func ErrServiceUnavailable(message string) *AppError {
	return New(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}
