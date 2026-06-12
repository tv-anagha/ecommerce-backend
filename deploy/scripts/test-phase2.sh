#!/bin/sh
# Phase 2 end-to-end check: one checkout → three independent consumer reactions.
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
DEPLOY_DIR="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"
GATEWAY="${GATEWAY:-http://localhost:8080}"
USER_ID="${USER_ID:-1}"
PRODUCT_ID="${PRODUCT_ID:-1}"

cd "$DEPLOY_DIR"

echo "=== Phase 2 Kafka fan-out test ==="
echo "Compose file: $COMPOSE_FILE"
echo ""

for svc in kafka notification-service fulfillment-service analytics-service order-service; do
  if ! docker compose -f "$COMPOSE_FILE" ps "$svc" 2>/dev/null | grep -qE 'Up|running'; then
    echo "ERROR: $svc is not running. Start with:"
    echo "  cd deploy && docker compose up -d"
    exit 1
  fi
done

echo "Waiting for consumers to connect to Kafka..."
sleep 15

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
echo "Waiting for all three consumers to process order #$ORDER_ID..."
sleep 10

PASS=0
FAIL=0

check_log() {
  svc="$1"
  pattern="$2"
  label="$3"
  if docker compose -f "$COMPOSE_FILE" logs "$svc" --tail 40 2>/dev/null | grep -q "$pattern"; then
    echo "PASS: $label"
    PASS=$((PASS + 1))
  else
    echo "FAIL: $label (pattern: $pattern)"
    FAIL=$((FAIL + 1))
  fi
}

check_log notification-service "Thank you for your order #$ORDER_ID" "notification-service thank-you"
check_log fulfillment-service "shipment queued for order #$ORDER_ID" "fulfillment-service shipment queued"
check_log analytics-service "recorded sale:" "analytics-service revenue recorded"

echo ""
if [ "$FAIL" -eq 0 ]; then
  echo "=== Phase 2 PASS: all three consumers reacted to order #$ORDER_ID ==="
  echo ""
  echo "--- notification-service ---"
  docker compose -f "$COMPOSE_FILE" logs notification-service --tail 5
  echo ""
  echo "--- fulfillment-service ---"
  docker compose -f "$COMPOSE_FILE" logs fulfillment-service --tail 5
  echo ""
  echo "--- analytics-service ---"
  docker compose -f "$COMPOSE_FILE" logs analytics-service --tail 5
  exit 0
fi

echo "=== Phase 2 FAIL: $FAIL consumer(s) missing expected logs ==="
docker compose -f "$COMPOSE_FILE" logs order-service --tail 10 2>/dev/null || true
exit 1
