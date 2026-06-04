# User service flow (short)

Basic auth for checkout: register → login → use JWT on protected calls.

## Ports and paths

| Where | URL |
|-------|-----|
| Gateway (frontend) | `http://localhost:8080/api/users/...` |
| user-service (internal) | `http://localhost:8082/users/...` |

The gateway rewrites `/api/users` → `/users` and forwards to `USER_SERVICE_URL`.

## Endpoints

| Method | Gateway | user-service | Purpose |
|--------|---------|--------------|---------|
| POST | `/api/users` or `/api/register` | `/users` or `/register` | Register |
| POST | `/api/users/login` or `/api/login` | `/users/login` or `/login` | Login → `token`, `userId`, `email` |
| GET | `/api/users/me` | `/users/me` | Who am I? Header: `Authorization: Bearer <token>` |
| GET | — | `/health` | Health check |

## Request flow

```text
Frontend → api-gateway:8080 → user-service:8082 → user_db (Postgres)
```

1. **Register** — password is hashed with bcrypt; email is stored in `user_db`.
2. **Login** — password is checked; service issues a **JWT** (24h, signed with `JWT_SECRET`).
3. **Checkout** — frontend keeps the token and sends `Authorization: Bearer <token>` to `/api/users/me` (or later to cart/order services) so the backend knows `userId`.

## Layer map

```text
cmd/server/main.go     → routes
internal/handler       → JSON in/out, read Bearer header
internal/service       → register, login, me
internal/repository    → GORM queries
internal/auth          → sign / parse JWT
internal/database      → Postgres + AutoMigrate users table
```

## Try it

```bash
# Register
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"shopper@example.com","password":"secret1234"}'

# Login
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"shopper@example.com","password":"secret1234"}'

# Me (replace TOKEN)
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer TOKEN"
```

Set `JWT_SECRET` in production (see `deploy/.env.example`).
