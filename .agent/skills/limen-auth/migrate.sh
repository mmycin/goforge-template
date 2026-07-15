#!/usr/bin/env bash
# Generate and apply the Limen schema migration.
# Run from the project root after rewriting model.go and kernel.go.
set -euo pipefail

echo "→ Generating Limen migration..."
goforge gen:migration add_limen_auth_tables

echo "→ Applying migration..."
goforge migrate

echo "✓ Limen tables created."
