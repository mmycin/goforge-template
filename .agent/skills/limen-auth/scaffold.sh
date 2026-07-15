#!/usr/bin/env bash
# Scaffold the auth service, remove Limen-incompatible generated files,
# and generate + apply the initial migration.
# Run from the project root.
set -euo pipefail

echo "→ Generating auth service..."
goforge gen:service auth

echo "→ Removing files Limen does not need..."
rm -f  internal/services/auth/handler.go
rm -f  internal/services/auth/grpc.go
rm -f  internal/services/auth/docs.go
rm -rf internal/proto/auth

echo "→ Installing Limen dependencies..."
bash .agent/skills/limen-auth/install-deps.sh

echo ""
echo "✓ Scaffold complete."
echo ""
echo "Next steps (manual):"
echo "  1. Rewrite internal/services/auth/model.go   — User, Session, Verification, Account"
echo "  2. Rewrite internal/services/auth/service.go — limen.New(...) + var Auth"
echo "  3. Rewrite internal/services/auth/routes.go  — engine.Any(\"/api/auth/*action\", ...)"
echo "  4. Update internal/services/kernel.go        — register 4 models + auth.Initialize()"
echo "  5. Run: bash .agent/skills/limen-auth/migrate.sh"
