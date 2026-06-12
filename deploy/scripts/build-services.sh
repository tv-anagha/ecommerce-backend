#!/bin/sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
DEPLOY_DIR="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"

cd "$DEPLOY_DIR"

if [ "$(id -u)" -eq 0 ]; then
  sh "$SCRIPT_DIR/ensure-swap.sh"
else
  sudo sh "$SCRIPT_DIR/ensure-swap.sh"
fi

if [ "$COMPOSE_FILE" = "docker-compose.minimal.yml" ]; then
  SERVICES="product-service user-service cart-service order-service api-gateway"
else
  SERVICES="product-service user-service cart-service order-service notification-service fulfillment-service analytics-service api-gateway"
fi

for svc in $SERVICES; do
  echo ""
  echo "=== Building $svc (one at a time to avoid OOM) ==="
  docker compose -f "$COMPOSE_FILE" build "$svc"
done

echo ""
echo "=== All service images built successfully ==="
