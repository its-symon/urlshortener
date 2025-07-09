#!/bin/bash

# CONFIG
URL="http://localhost:8080/shorten"
BODY='{"long_url":"https://google.com"}'
HEADERS=(-H "Content-Type: application/json")

# LOOP 12 TIMES (2 more than limit)
for i in {1..12}; do
  echo "Request $i:"
  curl -s -X POST "${URL}" "${HEADERS[@]}" -d "$BODY"
  echo -e "\n------"
done
