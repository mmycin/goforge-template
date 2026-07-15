#!/usr/bin/env bash
# Smoke-test the Limen auth endpoints against a running local server.
# Requires curl. Run from the project root while `goforge app serve` is running.
set -euo pipefail

BASE="${1:-http://localhost:3000}"
COOKIE_JAR="$(mktemp /tmp/limen-cookies-XXXXXX.txt)"
trap 'rm -f "$COOKIE_JAR"' EXIT

echo "→ Testing against $BASE"
echo ""

# ---------- Register ----------
echo "── Register ──"
curl -sf -X POST "$BASE/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@goforge.dev","password":"GoForge123!","first_name":"Go","last_name":"Forge"}' \
  | jq .
echo ""

# ---------- Login ----------
echo "── Login ──"
curl -sf -X POST "$BASE/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@goforge.dev","password":"GoForge123!"}' \
  -c "$COOKIE_JAR" \
  | jq .
echo ""

# ---------- Me ----------
echo "── Me (authenticated) ──"
curl -sf -X GET "$BASE/api/auth/me" \
  -b "$COOKIE_JAR" \
  | jq .
echo ""

# ---------- Logout ----------
echo "── Logout ──"
curl -sf -X POST "$BASE/api/auth/logout" \
  -b "$COOKIE_JAR" \
  | jq .
echo ""

# ---------- Me after logout (should 401) ----------
echo "── Me (after logout — expect 401) ──"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE/api/auth/me" -b "$COOKIE_JAR")
if [ "$STATUS" = "401" ]; then
  echo "✓ Got 401 as expected."
else
  echo "✗ Expected 401, got $STATUS"
  exit 1
fi

echo ""
echo "✓ All Limen auth smoke tests passed."
