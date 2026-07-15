# Configuration

Use this skill when adding, renaming, or validating settings.

Read first:
- `core/config/`
- `.env.example`

Workflow:
1. Add config shape in the relevant config type under `core/config/`.
2. Load and validate the value (via `viper` in the relevant `load*Config()` func).
3. Update `.env.example` with the new key and a comment.
4. Check all call sites that depend on the setting.

Guardrails:
- `core/config/` is framework-owned — users should not edit it.
- Security-sensitive settings must be configurable via env vars.
- Avoid silent fallback behavior for critical settings.
- Watch carefully for env key typos and mismatched names.
