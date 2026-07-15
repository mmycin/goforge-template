# State

Use this command to map the current system before changing it.

Capture:
- Entry point:       `app/main.go`
- App bootstrap:     `boot/kernel.go`
- HTTP server:       `boot/server/http.go`
- gRPC server:       `boot/server/grpc.go`
- Config loading:    `core/config/`
- Database wiring:   `core/database/database.go`
- Errors package:    `core/errors/errors.go`
- Validator package: `core/validator/validator.go`
- Service wiring:    `internal/services/kernel.go`
- Migrations:        `internal/database/migrations/`
- Queries:           `internal/database/queries/`
- Proto definitions: `internal/proto/<name>/`
- Tests:             `internal/tests/`
- Target domain:     `internal/services/<name>/`

Answer these questions:
1. Where is the behavior wired in?
2. What dependencies are injected?
3. What config or env keys influence it?
4. What tests currently cover it?
