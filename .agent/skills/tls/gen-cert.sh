#!/usr/bin/env bash
# Generate a self-signed TLS certificate for local development.
# Outputs cert.pem and key.pem in the project root.
# Run from the project root.
set -euo pipefail

DAYS="${1:-365}"
CN="${2:-localhost}"

echo "→ Generating self-signed cert (CN=$CN, valid ${DAYS} days)..."

openssl req -x509 \
  -newkey rsa:4096 \
  -keyout key.pem \
  -out cert.pem \
  -days "$DAYS" \
  -nodes \
  -subj "/CN=$CN" \
  -addext "subjectAltName=DNS:$CN,IP:127.0.0.1"

echo "✓ cert.pem and key.pem created."
echo ""
echo "Set in .env:"
echo "  TLS_ENABLED=true"
echo "  TLS_CERT_FILE=cert.pem"
echo "  TLS_KEY_FILE=key.pem"
echo ""
echo "Both files are gitignored — do not commit them."
