#!/usr/bin/env bash
# Install Limen authentication dependencies.
# Run from the project root.
set -euo pipefail

echo "→ Installing Limen core..."
go get github.com/thecodearcher/limen

echo "→ Installing Limen GORM adapter..."
go get github.com/thecodearcher/limen/adapters/gorm

echo "→ Installing credential-password plugin..."
go get "github.com/thecodearcher/limen/plugins/credential-password"

echo "→ Tidying module graph..."
go mod tidy

echo "✓ Limen dependencies installed."
