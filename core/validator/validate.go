package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	apperrors "github.com/mmycin/goforge/core/errors"
)

// FieldError represents a single field-level validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates any struct and returns a 422 AppError whose Details
// field contains []FieldError on failure, or nil on success.
//
// Usage in a handler:
//
//	if appErr := validator.Validate(req); appErr != nil {
//	    _ = c.Error(appErr)
//	    return
//	}
func Validate(s any) *apperrors.AppError {
	if err := validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if isValidationErrors(err, &ve) {
			fields := make([]FieldError, 0, len(ve))
			for _, fe := range ve {
				fields = append(fields, FieldError{
					Field:   fe.Field(),
					Message: fe.Translate(trans),
				})
			}
			return apperrors.ErrUnprocessable("validation failed").
				WithDetails(fields)
		}
		return apperrors.ErrUnprocessable(fmt.Sprintf("validation error: %v", err))
	}
	return nil
}

func isValidationErrors(err error, out *validator.ValidationErrors) bool {
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		*out = ve
	}
	return ok
}
