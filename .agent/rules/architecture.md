# Architecture Rule

Treat this repository as a generated application built on GoForge.

Structure:
- `boot/` — framework bootstrap and server infrastructure (users never edit this).
- `core/` — framework-owned packages that users consume but do not edit:
  - `core/config/`    — env/config loading
  - `core/database/`  — DB connection and GORM wrapper
  - `core/errors/`    — standardized AppError types and constructors
  - `core/validator/` — request validation wrapper (go-playground/validator)
  - `core/console/`   — CLI command runner and logger
- `internal/` — user-owned, everything here is meant to be edited:
  - `internal/services/<name>/` — service handler, routes, model, docs, grpc
  - `internal/services/kernel.go` — auto-generated composition root (do not hand-edit)
  - `internal/database/migrations/` — Atlas/GORM migration files
  - `internal/database/queries/`    — raw SQL queries (used by sqlc)
  - `internal/proto/<name>/`        — proto definitions and generated gRPC code
  - `internal/tests/`               — HTTP, gRPC, integration, and unit tests

Rules:
- Do not introduce package-level mutable state in `internal/` or `core/`.
- Prefer constructor injection with `New*` functions.
- Keep framework bootstrap concerns in `boot/` and `core/`, not in service packages.
- Keep service-specific behavior inside its own `internal/services/<name>/` package.
- If a change belongs in the generator rather than the generated app, stop and confirm whether `goforge-cli` should be updated too.
