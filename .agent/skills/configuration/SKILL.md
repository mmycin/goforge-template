# Configuration

Use this skill when adding, renaming, or validating settings.

Read first:
- `internal/config/`
- `.env.example`

Workflow:
1. Add config shape in the relevant config type.
2. Load and validate the value.
3. Update `.env.example`.
4. Check all call sites that depend on the setting.

Guardrails:
- Security-sensitive settings must be configurable.
- Avoid silent fallback behavior for critical settings.
- Watch carefully for env key typos and mismatched names.
