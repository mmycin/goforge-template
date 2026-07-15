# HTTP and gRPC

Use this skill when changing transports, middleware, routing, or protocol behavior.

Read first:
- `boot/server/http.go`
- `boot/server/grpc.go`
- `boot/server/middleware/`
- the owning service's `handler.go`, `routes.go`, or `grpc.go`

Checks:
- Are handlers delegating to services cleanly?
- Are middleware expectations still correct?
- Does the gRPC server still shut down gracefully?
- Do route registrations remain local and readable?

Guardrails:
- Keep HTTP and gRPC transport concerns out of core domain logic.
- Preserve existing startup and shutdown semantics.
