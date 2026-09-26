#!/usr/bin/env bash
# Remove a service, its proto files, and roll back its last migration.
# Usage: bash .agent/skills/cli-commands/remove-service.sh <service-name>
# Run from the project root.
set -euo pipefail

NAME="${1:?Usage: $0 <service-name>}"

echo "→ Rolling back latest migration..."
goforge rem:migration

echo "→ Removing service: $NAME"
goforge rem:service "$NAME"

echo "→ Removing proto files..."
goforge rem:proto

echo "✓ Service '$NAME' removed, proto files removed, and last migration rolled back."
echo ""
echo "Remember to remove the service wiring from internal/services/kernel.go manually."
