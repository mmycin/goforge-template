#!/usr/bin/env bash
# Generate a named Atlas migration from current GORM models and apply it.
# Usage: bash .agent/skills/database-migrations/migrate.sh <migration-name>
# Run from the project root.
set -euo pipefail

MIGRATION_NAME="${1:?Usage: $0 <migration-name>}"

echo "→ Generating migration: $MIGRATION_NAME"
goforge gen:migration "$MIGRATION_NAME"

echo "→ Applying migration..."
goforge migrate

echo "✓ Migration '$MIGRATION_NAME' applied."
