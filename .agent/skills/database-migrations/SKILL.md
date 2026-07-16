# Database Migrations

Use this skill when models, queries, or schema-affecting behavior change.

Read first:
- `core/database/database.go`         ← framework DB wiring (do not edit)
- `internal/database/migrations/`     ← user-owned: Atlas migration files
- `internal/database/queries/`        ← user-owned: raw SQL for sqlc
- `core/database/atlas.hcl`
- `core/database/sqlc.yaml`

Workflow:
1. Change the model in `internal/services/<name>/model.go`.
2. Run `bash .agent/skills/database-migrations/migrate.sh <migration-name>` to diff and apply.
3. If using sqlc: update queries in `internal/database/queries/`, run `bash .agent/skills/database-migrations/sqlc.sh`.
4. To undo the last migration: `bash .agent/skills/database-migrations/rollback.sh`.
5. Keep `.env.example` aligned if new database settings are required.
6. Verify the service code still handles nil and error cases correctly.

Guardrails:
- `core/database/database.go` is framework-owned — do not edit it.
- User migration files live in `internal/database/migrations/` — always user-editable.
- Do not treat schema changes as incidental.
- Keep migration-related changes small and explicit.
