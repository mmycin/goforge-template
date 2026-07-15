# Database Migrations

Use this skill when models, queries, or schema-affecting behavior change.

Read first:
- `internal/database/database.go`
- `internal/database/migrations/`
- `atlas.hcl`
- `sqlc.yaml`

Workflow:
1. Change the model or query.
2. Decide whether Atlas migration output is needed.
3. Keep `.env.example` aligned if new database settings are required.
4. Verify the service code still handles nil and error cases correctly.

Guardrails:
- Do not treat schema changes as incidental.
- Keep migration-related changes small and explicit.
