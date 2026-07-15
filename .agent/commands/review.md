# Review

Use this command for code review or pre-merge checks.

Review for:
- Behavior regressions first, style second.
- Broken dependency injection or accidental globals.
- Missing updates in `internal/services/kernel.go`.
- Config keys added in code but missing from `.env.example`.
- gRPC shutdown and timeout handling.
- Middleware or route registration drift.
- Database model changes without migration impact considered.
- Tests missing for new behavior.

Prefer findings in this order:
1. Bugs
2. Regressions
3. Security or config risks
4. Missing tests
5. Minor maintainability issues
