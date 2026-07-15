# Next

Use this command to choose the next safe step.

Decision order:
1. Identify the user-visible outcome.
2. Locate the owning package.
3. Prefer editing existing files over creating new ones.
4. If the change affects a service, inspect:
   - `internal/services/<name>/`
   - `internal/services/kernel.go`
   - `internal/database/`
5. If the change affects startup or transport, inspect:
   - `boot/kernel.go`
   - `boot/server/http.go`
   - `boot/server/grpc.go`

Before implementing:
- Check for matching tests.
- Check whether `.env.example` needs to change.
- Check whether graceful shutdown, validation, or middleware behavior is affected.
