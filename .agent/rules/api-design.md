# API Design Rule

When changing HTTP or gRPC behavior:
- Keep transport thin; business logic belongs in services.
- Validate request input at the handler edge using `core/validator`.
- Return consistent status codes and error shapes using `core/errors` + `c.Error()`.
- The `ErrorHandler` middleware (registered in `boot/server/http.go`) writes all error responses; handlers must NOT call `c.JSON` for errors directly.
- Keep route registration local to the owning service package.
- For gRPC, maintain graceful shutdown behavior and avoid long-running uncontrolled goroutines.
- Prefer additive API changes unless the user explicitly requests a breaking change.
- Imports must use `github.com/mmycin/GoForge/core/...` or `github.com/mmycin/GoForge/boot/...`
