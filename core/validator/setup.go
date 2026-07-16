// Package validator wraps go-playground/validator/v10 and produces
// standardized, user-friendly validation error responses that fit the
// AppError shape from core/errors.
package validator

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
)

// init wires the validator into Gin's binding engine and configures
// English translations. Runs exactly once on package import.
func init() {
	once.Do(func() {
		enLocale := en.New()
		uni := ut.New(enLocale, enLocale)
		trans, _ = uni.GetTranslator("en")

		// Prefer the validator instance already bound to Gin so
		// ShouldBindJSON and Validate share the same rule set.
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			validate = v
		} else {
			validate = validator.New()
		}

		_ = enTranslations.RegisterDefaultTranslations(validate, trans)

		// Use json tags in error messages instead of Go field names.
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}
			return name
		})
	})
}
