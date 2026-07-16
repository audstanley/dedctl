#!/bin/bash

# Test script for dedctl API

echo "Testing dedctl API..."

# 1. Login to get JWT token
echo ""
echo "=== 1. Login ==="
TOKEN=$(curl -s -X POST http://localhost:8085/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)

if [ -z "$TOKEN" ]; then
  echo "ERROR: Login failed. Is the backend running on port 8085?"
  exit 1
fi
echo "Login successful, token: ${TOKEN:0:20}..."

# 2. Test server-info (no auth required)
echo ""
echo "=== 2. Server Info (no auth) ==="
curl -s -X GET http://localhost:8085/server-info | python3 -m json.tool 2>/dev/null

# 3. Test game listing (auth required)
echo ""
echo "=== 3. List Games ==="
curl -s -X GET "http://localhost:8085/games?token=${TOKEN}" | python3 -m json.tool 2>/dev/null

# 4. Test game status (auth required)
echo ""
echo "=== 4. Game Status ==="
curl -s -X GET "http://localhost:8085/games/corekeeper/status?token=${TOKEN}" | python3 -m json.tool 2>/dev/null

# 5. Test logs streaming (auth required) - only first 2 lines
echo ""
echo "=== 5. Logs (first 2 lines) ==="
curl -s -N "http://localhost:8085/games/corekeeper/logs?token=${TOKEN}" --max-time 3 | head -n 4

echo ""
echo "API test completed!"
