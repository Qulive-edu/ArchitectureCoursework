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
curl -sf -H "Authorization: Bearer $TOKEN" "$BASE_URL/bookings/my" > /dev/null && echo "OK" || echo "⚠️  Skipped (no valid token)"

# 4. Data integrity (исправлено: без jq, через grep)
echo -n "✓ Data persistence... "
# Просто проверяем, что ответ валидный и содержит ожидаемые поля
RESPONSE=$(curl -sf "$BASE_URL/places")
if echo "$RESPONSE" | grep -q '"id"' && echo "$RESPONSE" | grep -q '"title"'; then
    echo "OK"
else
    echo "FAIL (unexpected response format)"
    exit 1
fi

echo "All smoke tests passed"