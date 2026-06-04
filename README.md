# Ecommerce Backend (Microservices)

Go monorepo with an API gateway and product service. User, cart, and order services can be added in later commits.

## Architecture

| Service | Port | Responsibility |
|---------|------|----------------|
| api-gateway | 8080 | Public entry, CORS, reverse proxy |
| product-service | 8081 | Product catalog |

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

1. Create database `product_db` (see `deploy/init/`).
2. Copy `deploy/.env.example` to repo root `.env` and set credentials.
3. For an existing local `ecommerce` DB with products, run product-service with `DB_NAME=ecommerce`.
4. Start each service in a separate terminal:

```bash
cd product-service && PORT=8081 DB_NAME=product_db go run ./cmd/server
cd api-gateway && PORT=8080 go run ./cmd/server
```

## API examples (via gateway)

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/products
curl http://localhost:8080/api/products/1
```

## EC2

1. Open security group port **8080** (api-gateway).
2. Run `docker compose up -d` from `deploy/`.
3. Set frontend `VITE_API_BASE_URL=http://<EC2_PUBLIC_IP>:8080`.
4. Add EC2 frontend origin to `CORS_ALLOWED_ORIGINS` on the gateway.

## Frontend

In `ecommerce-frontend/product-mf`:

```bash
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

Products are loaded from `/api/products` on the gateway.
