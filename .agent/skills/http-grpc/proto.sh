#!/usr/bin/env bash
# Compile proto definitions and regenerate gRPC stubs.
# Usage:
#   bash .agent/skills/http-grpc/proto.sh              # compile all services
#   bash .agent/skills/http-grpc/proto.sh <name>       # compile one service
#   bash .agent/skills/http-grpc/proto.sh --clean      # remove all generated .pb.go files
# Run from the project root.
set -euo pipefail

ARG="${1:-}"

case "$ARG" in
  --clean)
    echo "→ Removing generated proto files..."
    goforge rem:proto
    echo "✓ Generated proto files removed."
    ;;
  "")
    echo "→ Compiling all proto definitions in internal/proto/..."
    goforge gen:proto
    echo "✓ All proto stubs regenerated."
    ;;
  *)
    echo "→ Compiling proto for service: $ARG"
    goforge gen:proto "$ARG"
    echo "✓ Proto stubs regenerated for '$ARG'."
    echo "  Output: internal/proto/$ARG/gen/"
    ;;
esac
