# Ecommerce Backend (Microservices)

Go monorepo with an API gateway and four domain services.

## Architecture

| Service | Port | Responsibility |
|---------|------|----------------|
| api-gateway | 8080 | Public entry, CORS, reverse proxy |
| product-service | 8081 | Product catalog |
| user-service | 8082 | Registration, login, profile |
| cart-service | 8083 | Shopping carts |
| order-service | 8084 | Orders from carts |

Frontend calls the gateway only, e.g. `GET http://localhost:8080/api/products`.

## Quick start (Docker)

```bash
cd deploy
docker compose up --build
```

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/products
```

## Local development (without Docker)

1. Create databases: `product_db`, `user_db`, `cart_db`, `order_db` (see `deploy/init/`).
2. Copy `deploy/.env.example` to repo root `.env` and set credentials.
3. For existing local `ecommerce` DB with products, run product-service with `DB_NAME=ecommerce`.
4. Start each service in a separate terminal:

```bash
cd product-service && PORT=8081 DB_NAME=product_db go run ./cmd/server
cd user-service && PORT=8082 DB_NAME=user_db go run ./cmd/server
cd cart-service && PORT=8083 DB_NAME=cart_db go run ./cmd/server
cd order-service && PORT=8084 DB_NAME=order_db go run ./cmd/server
cd api-gateway && PORT=8080 go run ./cmd/server
```

## API examples (via gateway)

```bash
# Products
curl http://localhost:8080/api/products
curl http://localhost:8080/api/products/1

# Users
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret12"}'

curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret12"}'

# Cart (userId=1)
curl http://localhost:8080/api/carts/1
curl -X POST http://localhost:8080/api/carts/1/items \
  -H "Content-Type: application/json" \
  -d '{"productId":1,"quantity":2}'

# Orders
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"userId":1}'
```

## EC2

1. Open security group port **8080** (gateway).
2. Run `docker compose up -d` from `deploy/`.
3. Set frontend `VITE_API_BASE_URL=http://<EC2_PUBLIC_IP>:8080`.
4. Add EC2 frontend origin to `CORS_ALLOWED_ORIGINS` on the gateway.

## Frontend

In `ecommerce-frontend/product-mf`:

```bash
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

Products are loaded from `/api/products` on the gateway.
