# Error Handling

Use this skill when a task touches error responses, error wrapping, or HTTP status codes.

Read first:
- `core/errors/errors.go`
- `boot/server/middleware/error.go`

Key concepts:
- All API errors use `*apperrors.AppError` with fields: `code`, `message`, `details`.
- Handlers attach errors to the Gin context via `c.Error(appErr)` and return early.
- `ErrorHandler` middleware (in `boot/server/middleware/error.go`) intercepts these and writes the JSON response.
- Never write raw `c.JSON(...)` for errors inside handlers; use `c.Error(...)`.

Constructor helpers:
- `apperrors.ErrBadRequest(msg)`       → 400
- `apperrors.ErrUnauthorized(msg)`     → 401
- `apperrors.ErrForbidden(msg)`        → 403
- `apperrors.ErrNotFound(msg)`         → 404
- `apperrors.ErrConflict(msg)`         → 409
- `apperrors.ErrUnprocessable(msg)`    → 422
- `apperrors.ErrInternalServer(msg)`   → 500

Wrap underlying errors for logging context:
```go
apperrors.ErrNotFound("todo not found").Wrap(err)
```

Guardrails:
- Keep the 5xx response message generic ("an unexpected error occurred") to avoid leaking internals to clients.
- Always wrap the original `error` so it remains available for logging.
- Do not add new HTTP status code handling without adding a constructor helper.
