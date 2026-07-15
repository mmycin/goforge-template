# GoForge Core

Use this skill when a task touches bootstrap, lifecycle, or shared wiring.

Read first:
- `boot/kernel.go`
- `internal/services/kernel.go`
- `cmd/main.go`

Focus:
- App construction
- Dependency injection
- Service registration
- Shared startup and shutdown behavior

Guardrails:
- Keep bootstrap code in `boot/`.
- Do not move app-specific logic into framework internals without reason.
- Preserve existing initialization flow unless the task requires a structural change.
