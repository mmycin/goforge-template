package errors

import (
	"errors"
	"net/http"
)

// Is reports whether any error in err's chain is an *AppError with the given code.
func Is(err error, code string) bool {
	var e *AppError
	return errors.As(err, &e) && e.Code == code
}

// IsNotFound reports whether err (or any error in its chain) is a 404 AppError.
func IsNotFound(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusNotFound
}

// IsBadRequest reports whether err is a 400 AppError.
func IsBadRequest(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusBadRequest
}

// IsUnauthorized reports whether err is a 401 AppError.
func IsUnauthorized(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusUnauthorized
}

// IsForbidden reports whether err is a 403 AppError.
func IsForbidden(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusForbidden
}

// IsConflict reports whether err is a 409 AppError.
func IsConflict(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusConflict
}

// IsUnprocessable reports whether err is a 422 AppError.
func IsUnprocessable(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusUnprocessableEntity
}

// IsInternalServer reports whether err is a 500 AppError.
func IsInternalServer(err error) bool {
	var e *AppError
	return errors.As(err, &e) && e.HTTPStatus == http.StatusInternalServerError
}

// HTTPStatus extracts the HTTP status code from err if it is an *AppError.
// Returns 500 for unknown errors.
func HTTPStatus(err error) int {
	var e *AppError
	if errors.As(err, &e) {
		return e.HTTPStatus
	}
	return http.StatusInternalServerError
}
