// Package validator wraps go-playground/validator/v10 and produces
// standardized, user-friendly validation error responses that fit the
// AppError shape from internal/errors.
package validator

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	apperrors "github.com/mmycin/goforge/core/errors"
)

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

var (
	once      sync.Once
	validate  *validator.Validate
	trans     ut.Translator
)

// init registers the validator with Gin's binding engine and sets up English
// translations once on package import.
func init() {
	once.Do(func() {
		enLocale := en.New()
		uni := ut.New(enLocale, enLocale)
		trans, _ = uni.GetTranslator("en")

		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			validate = v
			_ = enTranslations.RegisterDefaultTranslations(validate, trans)
		} else {
			validate = validator.New()
			_ = enTranslations.RegisterDefaultTranslations(validate, trans)
		}

		// Use JSON tag names in error messages instead of struct field names
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}
			return name
		})
	})
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// FieldError represents a single field-level validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct and returns a 422 AppError whose Details field
// contains []FieldError when validation fails, or nil on success.
func Validate(s any) *apperrors.AppError {
	if err := validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if ok := isValidationErrors(err, &ve); ok {
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

// RegisterRule registers a custom validation rule under the given tag.
//
// Example:
//
//	validator.RegisterRule("strong_password", func(fl validator.FieldLevel) bool {
//	    return len(fl.Field().String()) >= 12
//	})
func RegisterRule(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isValidationErrors(err error, out *validator.ValidationErrors) bool {
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		*out = ve
	}
	return ok
}
