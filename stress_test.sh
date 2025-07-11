#!/bin/bash

# CONFIG
BASE_URL="http://localhost:8080"
REGISTER_URL="$BASE_URL/register"
LOGIN_URL="$BASE_URL/login"
API_KEY_URL="$BASE_URL/generate-api-key"
SHORTEN_URL="$BASE_URL/shorten"

EMAIL="user123@example.com"
PASSWORD="string"

# Register
echo "[+] Registering user..."
curl -s -X POST "$REGISTER_URL" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\", \"password\":\"$PASSWORD\"}"
echo -e "\n------"

# Login
echo "[+] Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\", \"password\":\"$PASSWORD\"}")

# Extract JWT token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [[ "$TOKEN" == "null" || -z "$TOKEN" ]]; then
  echo "[-] Failed to extract token. Exiting."
  exit 1
fi

echo "[+] JWT Token received."
echo "------"

# Generate API Key
echo "[+] Generating API key..."
API_KEY_RESPONSE=$(curl -s -X POST "$API_KEY_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")

API_KEY=$(echo "$API_KEY_RESPONSE" | jq -r '.api_key')

if [[ "$API_KEY" == "null" || -z "$API_KEY" ]]; then
  echo "[-] Failed to extract API key. Exiting."
  exit 1
fi

echo "[+] API Key: $API_KEY"
echo "------"

# Prepare headers
HEADERS=(-H "Content-Type: application/json" -H "x-api-key: $API_KEY")
BODY='{"long_url":"https://google.com"}'

# Loop and send requests
for i in {1..12}; do
  echo "Request $i:"
  curl -s -X POST "$SHORTEN_URL" "${HEADERS[@]}" -d "$BODY"
  echo -e "\n------"
done
