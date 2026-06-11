# Ecommerce Backend (Microservices)

Go monorepo with API gateway, product catalog, cart, orders, user auth, and **Phase 1 Kafka** (`order.placed` → notification-service).

## Architecture

| Service | Port | Responsibility |
|---------|------|----------------|
| api-gateway | 8080 | Public entry, CORS, reverse proxy |
| product-service | 8081 | Product catalog |
| user-service | 8082 | Register, login, JWT for checkout |
| cart-service | 8083 | Shopping carts |
| order-service | 8084 | Checkout; publishes `order.placed` to Kafka |
| notification-service | 8085 | Consumes `order.placed`, logs thank-you |
| kafka | 9092 | Event broker (Docker internal: `kafka:29092`) |

See [USER_SERVICE_FLOW.md](USER_SERVICE_FLOW.md), [ORDER_SERVICE_FLOW.md](ORDER_SERVICE_FLOW.md), and [KAFKA_IMPLEMENTATION.md](KAFKA_IMPLEMENTATION.md).

## Quick start (Docker — full stack with Kafka)

```bash
cd deploy
docker compose up --build -d
./scripts/test-phase1.sh   # verify order.placed → thank-you loop
```

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/products
```

**Small VMs:** use `docker-compose.minimal.yml` (no Kafka). Phase 1 requires full `docker-compose.yml`.

## Local development (without Docker)

1. Create databases `product_db` and `user_db` (see `deploy/init/`).
2. Copy `deploy/.env.example` to repo root `.env` and set credentials.
3. Start each service in a separate terminal:

```bash
cd product-service && PORT=8081 DB_NAME=product_db go run ./cmd/server
cd user-service && PORT=8082 DB_NAME=user_db JWT_SECRET=dev-secret go run ./cmd/server
cd api-gateway && PORT=8080 go run ./cmd/server
```

## API examples (via gateway)

```bash
# Products
curl http://localhost:8080/api/products
curl http://localhost:8080/api/products/1

# Users (checkout auth)
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret1234"}'

curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret1234"}'

curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <token-from-login>"
```

## EC2

1. Open security group port **8080** (api-gateway).
2. Run `docker compose up -d` from `deploy/`.
3. Set frontend `VITE_API_BASE_URL=http://<EC2_PUBLIC_IP>:8080`.
4. Set a strong `JWT_SECRET` for user-service in production.

## Frontend

```bash
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```
