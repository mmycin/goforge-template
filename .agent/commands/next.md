# Next

Use this command to choose the next safe step.

Decision order:
1. Identify the user-visible outcome.
2. Locate the owning package.
3. Prefer editing existing files over creating new ones.
4. If the change affects a service, inspect:
   - `internal/services/<name>/`
   - `internal/services/kernel.go`
   - `internal/database/` (migrations, queries)
   - `internal/proto/<name>/` (proto definitions)
5. If the change affects startup or transport, inspect:
   - `boot/kernel.go`
   - `boot/server/http.go`
   - `boot/server/grpc.go`
6. If the change affects framework config, errors, or validation, inspect:
   - `core/config/`
   - `core/errors/`
   - `core/validator/`

Before implementing:
- Check for matching tests in `internal/tests/`.
- Check whether `.env.example` needs to change.
- Check whether graceful shutdown, validation, or middleware behavior is affected.
