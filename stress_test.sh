#!/bin/bash

# CONFIG
URL="http://localhost:8080/shorten"
BODY='{"long_url":"https://google.com"}'
API_KEY="c0f79396c338551e76b660d8995c4824a0a5552970e8427a0dcb22fe0cc0505f"
HEADERS=(-H "Content-Type: application/json" -H "x-api-key: $API_KEY")

# LOOP 12 TIMES
for i in {1..12}; do
  echo "Request $i:"
  curl -s -X POST "${URL}" "${HEADERS[@]}" -d "$BODY"
  echo -e "\n------"
done
