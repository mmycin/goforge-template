#!/usr/bin/env bash
# Regenerate SQLC type-safe query code from internal/database/queries/*.sql
# Run from the project root.
set -euo pipefail

echo "→ Running sqlc generation..."
goforge gen:sqlc

echo "✓ SQLC code regenerated at internal/database/gen/"
