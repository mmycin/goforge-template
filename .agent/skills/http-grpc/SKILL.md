# HTTP and gRPC

Use this skill when changing transports, middleware, routing, or protocol behavior.

Read first:
- `boot/server/http.go`
- `boot/server/grpc.go`
- `boot/server/middleware/`
- the owning service's `internal/services/<name>/handler.go`, `routes.go`, or `grpc.go`

Proto files:
- Definitions live in `internal/proto/<name>/<name>.proto`
- Generated code goes in `internal/proto/<name>/gen/`

Checks:
- Are handlers delegating to services cleanly?
- Are middleware expectations still correct?
- Does the gRPC server still shut down gracefully?
- Do route registrations remain local and readable?

Guardrails:
- `boot/server/` is framework-owned — do not edit it for service-specific logic.
- Keep HTTP and gRPC transport concerns out of core domain logic.
- Preserve existing startup and shutdown semantics.
