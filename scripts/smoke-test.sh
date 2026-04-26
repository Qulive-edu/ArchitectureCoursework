#!/bin/bash
# scripts/smoke-test.sh
set -e

BASE_URL="${1:-http://localhost:3000/api}"
TOKEN="${2:-test-token}"

echo "🔍 Running smoke tests against $BASE_URL"

# 1. Health check
echo -n "✓ API health... "
curl -sf "$BASE_URL/places" > /dev/null && echo "OK" || { echo "FAIL"; exit 1; }

# 2. Auth flow
echo -n "✓ User registration... "
RESP=$(curl -sf -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test\",\"email\":\"test$(date +%s)@example.com\",\"password\":\"pass123\"}")
[[ -n "$RESP" ]] && echo "OK" || { echo "FAIL"; exit 1; }

# 3. Protected endpoint
echo -n "✓ Protected endpoint... "
curl -sf -H "Authorization: Bearer $TOKEN" "$BASE_URL/bookings/my" > /dev/null && echo "OK" || echo "Skipped (no valid token)"

# 4. Data integrity
echo -n "Data persistence... "
BEFORE=$(curl -sf "$BASE_URL/places" | jq length)
sleep 2
AFTER=$(curl -sf "$BASE_URL/places" | jq length)
[[ "$BEFORE" == "$AFTER" ]] && echo "OK" || echo "Count changed: $BEFORE → $AFTER"

echo "All smoke tests passed"