#!/bin/sh
# Phase 1 end-to-end check: checkout publishes order.placed; notification-service logs thank-you.
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
DEPLOY_DIR="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"
GATEWAY="${GATEWAY:-http://localhost:8080}"
USER_ID="${USER_ID:-1}"
PRODUCT_ID="${PRODUCT_ID:-1}"

cd "$DEPLOY_DIR"

echo "=== Phase 1 Kafka loop test ==="
echo "Compose file: $COMPOSE_FILE"
echo ""

if ! docker compose -f "$COMPOSE_FILE" ps kafka notification-service order-service 2>/dev/null | grep -q kafka; then
  echo "ERROR: Kafka stack not running. Start with:"
  echo "  cd deploy && docker compose up -d"
  exit 1
fi

echo "Waiting for Kafka to be healthy..."
for i in 1 2 3 4 5 6 7 8 9 10 11 12; do
  if docker compose -f "$COMPOSE_FILE" ps kafka 2>/dev/null | grep -q healthy; then
    break
  fi
  sleep 5
done

echo "Adding item to cart for user $USER_ID..."
curl -sf -X POST "$GATEWAY/api/carts/$USER_ID/items" \
  -H "Content-Type: application/json" \
  -d "{\"productId\":$PRODUCT_ID,\"quantity\":1}" >/dev/null

echo "Placing order..."
ORDER_JSON=$(curl -sf -X POST "$GATEWAY/api/orders" \
  -H "Content-Type: application/json" \
  -d "{\"userId\":$USER_ID}")
echo "$ORDER_JSON"

ORDER_ID=$(echo "$ORDER_JSON" | sed -n 's/^{"id":[[:space:]]*\([0-9][0-9]*\).*/\1/p')
if [ -z "$ORDER_ID" ]; then
  echo "ERROR: could not parse order id from response"
  exit 1
fi

echo ""
echo "Waiting for notification-service to consume order.placed..."
sleep 5

if docker compose -f "$COMPOSE_FILE" logs notification-service --tail 30 2>/dev/null | grep -q "Thank you for your order #$ORDER_ID"; then
  echo ""
  echo "PASS: notification-service logged thank-you for order #$ORDER_ID"
  docker compose -f "$COMPOSE_FILE" logs notification-service --tail 8
  exit 0
fi

echo ""
echo "FAIL: thank-you message not found in notification-service logs."
echo "Check order-service for publish errors:"
docker compose -f "$COMPOSE_FILE" logs order-service --tail 15 2>/dev/null || true
echo ""
echo "notification-service logs:"
docker compose -f "$COMPOSE_FILE" logs notification-service --tail 20 2>/dev/null || true
exit 1
