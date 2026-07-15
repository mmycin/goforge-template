# GoForge Core

Use this skill when a task touches bootstrap, lifecycle, or shared wiring.

Read first:
- `boot/kernel.go`
- `internal/services/kernel.go`
- `app/main.go`

Directory map:
- `boot/` — framework bootstrap (users never edit): server lifecycle, middleware, client helpers
- `core/` — framework packages users import but never edit: config, database, errors, validator, console
- `internal/` — user-owned: services, proto definitions, migrations, queries, tests

Focus:
- App construction
- Dependency injection
- Service registration
- Shared startup and shutdown behavior

Guardrails:
- Keep bootstrap code in `boot/` and `core/`.
- Do not move user-editable logic into `core/` or `boot/`.
- Preserve existing initialization flow unless the task requires a structural change.
