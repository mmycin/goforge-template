#!/usr/bin/env bash
# Run tests at different scopes.
# Usage: bash .agent/skills/testing/run.sh [scope]
#   scope: http | grpc | integration | unit | all (default: all)
# Run from the project root.
set -euo pipefail

SCOPE="${1:-all}"

run_scope() {
  local label="$1"
  local pkg="$2"
  echo "── $label ──"
  go test -v -count=1 "$pkg" 2>&1 | tail -20
  echo ""
}

case "$SCOPE" in
  http)
    run_scope "HTTP tests" "./internal/tests/http/..."
    ;;
  grpc)
    run_scope "gRPC tests" "./internal/tests/grpc/..."
    ;;
  integration)
    run_scope "Integration tests" "./internal/tests/integration/..."
    ;;
  unit)
    run_scope "Unit tests" "./internal/tests/unit/..."
    ;;
  all)
    run_scope "HTTP tests"        "./internal/tests/http/..."
    run_scope "gRPC tests"        "./internal/tests/grpc/..."
    run_scope "Integration tests" "./internal/tests/integration/..."
    run_scope "Unit tests"        "./internal/tests/unit/..."
    ;;
  *)
    echo "Unknown scope '$SCOPE'. Valid: http | grpc | integration | unit | all"
    exit 1
    ;;
esac

echo "✓ Done."
