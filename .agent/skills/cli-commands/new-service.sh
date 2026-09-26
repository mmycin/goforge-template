#!/usr/bin/env bash
# Scaffold a new service, generate events, and run the full post-generation workflow.
# Usage: bash .agent/skills/cli-commands/new-service.sh <service-name>
# Run from the project root.
set -euo pipefail

NAME="${1:?Usage: $0 <service-name>}"

echo "→ Generating service: $NAME"
goforge gen:service "$NAME"

echo "→ Generating events for: $NAME"
goforge gen:event "$NAME"

echo "→ Generating migration: add_${NAME}_table"
goforge gen:migration "add_${NAME}_table"

echo "→ Applying migration..."
goforge migrate

echo ""
echo "✓ Service '$NAME' is ready."
echo ""
echo "Files created:"
echo "  internal/services/$NAME/model.go"
echo "  internal/services/$NAME/service.go"
echo "  internal/services/$NAME/handler.go"
echo "  internal/services/$NAME/routes.go"
echo "  internal/services/$NAME/docs.go"
echo "  internal/services/$NAME/grpc.go"
echo "  internal/services/$NAME/events.go"
echo "  internal/services/$NAME/listeners.go"
echo "  internal/proto/$NAME/$NAME.proto"
echo "  internal/services/kernel.go  (regenerated)"
echo "  internal/services/kernel_models.go  (regenerated)"
