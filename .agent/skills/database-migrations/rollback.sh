#!/usr/bin/env bash
# Remove the latest migration file and revert the Atlas hash.
# Run from the project root.
set -euo pipefail

echo "→ Rolling back latest migration..."
goforge rem:migration

echo "✓ Latest migration rolled back."
echo "  Re-run 'goforge gen:migration <name>' when ready to regenerate."
