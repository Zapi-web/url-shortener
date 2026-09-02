DOMAIN=$1
EXPECTED_URL="https://google.com/"

ANSWER=$(curl -s -X POST "$DOMAIN/api/v1/" \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$EXPECTED_URL\", \"user_id\": 1}" | jq -r '.short_url')

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$DOMAIN/$ANSWER")
REDIRECT_URL=$(curl -s -o /dev/null -w "%{redirect_url}" "$DOMAIN/$ANSWER")

if [ "$HTTP_CODE" -ne 302 ] && [ "$HTTP_CODE" -ne 301 ]; then
  echo "Error: not expected status $HTTP_CODE"
  exit 1
fi

if [ "$REDIRECT_URL" != "$EXPECTED_URL" ]; then
  echo "Error: expected $EXPECTED_URL, got $REDIRECT_URL"
  exit 1
fi

echo "test completed: status $HTTP_CODE, link $REDIRECT_URL"