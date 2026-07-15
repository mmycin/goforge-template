# Validation

Use this skill when a task changes request validation or adds new input rules.

Read first:
- `core/validator/validator.go`

Key concepts:
- The `validator` package wraps `go-playground/validator/v10` with English translations and JSON field names.
- Annotate request structs with `validate:` tags.
- Call `validator.Validate(req)` after `ShouldBindJSON`; it returns a `*apperrors.AppError` (422) with per-field details, or nil.

Example:
```go
type CreateUserRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

if appErr := validator.Validate(req); appErr != nil {
    _ = c.Error(appErr)
    return
}
```

Custom rules:
```go
// Register once at startup (e.g., in init() or boot phase)
validator.RegisterRule("strong_password", func(fl validator.FieldLevel) bool {
    return len(fl.Field().String()) >= 12
})
```

Guardrails:
- Always validate input at the handler edge, not inside service methods.
- Use JSON tag names on structs so error messages match what the client sends.
- Do not bypass the validator package with raw `validate.Struct()` calls.
