package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterRule adds a custom validation rule under the given tag.
// Call this during application startup, before any requests are handled.
//
// Example — require a password of at least 12 characters:
//
//	validator.RegisterRule("strong_password", func(fl validator.FieldLevel) bool {
//	    return len(fl.Field().String()) >= 12
//	})
//
// Then use it in a struct:
//
//	type RegisterRequest struct {
//	    Password string `json:"password" validate:"required,strong_password"`
//	}
func RegisterRule(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

// RegisterRuleWithMessage adds a custom validation rule and registers a
// human-readable message that appears in FieldError.Message.
//
//	validator.RegisterRuleWithMessage(
//	    "strong_password",
//	    func(fl validator.FieldLevel) bool { return len(fl.Field().String()) >= 12 },
//	    "{0} must be at least 12 characters",
//	)
func RegisterRuleWithMessage(tag string, fn validator.Func, message string) error {
	if err := validate.RegisterValidation(tag, fn); err != nil {
		return err
	}
	return validate.RegisterTranslation(
		tag,
		trans,
		func(t ut.Translator) error {
			return t.Add(tag, message, true)
		},
		func(t ut.Translator, fe validator.FieldError) string {
			msg, _ := t.T(tag, fe.Field())
			return msg
		},
	)
}
