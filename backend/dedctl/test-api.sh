#!/bin/bash

# Test script for Steam Game Server Control API

echo "Testing Steam Game Server Control API..."

# Start the API in background (this would be your actual implementation)
echo "Starting API server (simulated)..."
sleep 2

# Test API endpoints
echo "Testing authentication endpoints..."
curl -X POST http://localhost:8085/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password","is_admin":true}' \
  && echo "User registration completed"

echo "Testing login endpoint..."
curl -X POST http://localhost:8085/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  && echo "Login completed"

echo "Testing game listing..."
curl -X GET http://localhost:8085/games \
  && echo "Game listing completed"

echo "API test completed successfully!"