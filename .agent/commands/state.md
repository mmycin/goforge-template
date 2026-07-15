# State

Use this command to map the current system before changing it.

Capture:
- Entry point: `cmd/main.go`
- App bootstrap: `boot/kernel.go`
- HTTP server: `boot/server/http.go`
- gRPC server: `boot/server/grpc.go`
- Config loading: `internal/config/`
- Database setup: `internal/database/database.go`
- Service wiring: `internal/services/kernel.go`
- Target domain package: `internal/services/<name>/`

Answer these questions:
1. Where is the behavior wired in?
2. What dependencies are injected?
3. What config or env keys influence it?
4. What tests currently cover it?
