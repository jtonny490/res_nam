## res_nam

Local Lake Victoria pollution reporting application. The KijaniSpace integration is represented by a provider interface and runs in `mock` mode until the live API contract and credentials are available.

### Structure

```
backend/    Go + Gin + GORM API (Dockerfile: Go build -> Alpine, listens on :8080)
frontend/   React + Vite SPA (Dockerfile: Node build -> Nginx, serves the SPA and
            reverse-proxies /api, /uploads and /health to the backend)
docker-compose.yml   db (Postgres 16) + backend (internal) + frontend (public)
```

The backend is API-only; Nginx serves the built SPA and proxies API calls, so the
frontend and API share one origin.

### Run with Docker (recommended)

```sh
cp .env.example .env      # then edit values
make docker-up            # or: docker-compose up --build
```

Open `http://localhost:8080`. Create an account, submit a report using coordinates in the pilot area (latitude -1.6 to 0.7, longitude 31.5 to 35.0), and view it in the feed.

### Run locally without Docker

The backend applies the SQL migrations in `backend/migrations/` at startup.

```sh
# 1. database
docker run -d --name res_nam_pg -p 5432:5432 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=res_nam \
  postgres:16-alpine

# 2. backend (DB_HOST etc. default to localhost; blank DATABASE_URL to use it)
cd backend && go run ./cmd/server      # http://localhost:8080

# 3. frontend with hot reload (proxies /api to localhost:8080)
cd frontend && npm install && npm run dev   # http://localhost:5173
```

### Configuration

Environment variables (see `.env.example`):

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | Full Postgres URL; takes precedence over `DB_*`. TLS (`sslmode=require`) is added when absent, for managed providers such as Render. |
| `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASSWORD`/`DB_NAME`/`DB_SSLMODE` | Discrete connection settings, used when `DATABASE_URL` is empty. |
| `JWT_SECRET` | HS256 signing secret (rotate for production). |
| `CORS_ORIGIN` | Allowed origin; `*` for local development. |
| `KIJANI_MODE`/`KIJANI_BASE_URL`/`KIJANI_API_KEY` | KijaniSpace provider (mock until the live contract is available). |
| `MIGRATIONS_DIR` | Migrations directory (default `migrations`). |
| `PORT` | Host port the frontend is published on (default `8080`). |

### API

- `POST /api/auth/register`, `POST /api/auth/login`
- `GET /api/reports` (supports `?category=`, `?status=`, `?page=`, `?limit=`), `GET /api/reports/:id`
- `GET /api/reports/:id/comments`
- Authenticated: `POST /api/reports` (JSON or multipart with a `photo` field), `PATCH /api/reports/:id`, `DELETE /api/reports/:id`, comments, likes (`POST`/`DELETE /api/reports/:id/like`), and `GET /api/reports/:id/assessment`
- Authority/admin: `PATCH /api/reports/:id/status`, `GET /api/map/reports`
- Admin: `PATCH /api/reports/:id/resolve`, `GET /api/admin/users`, `PATCH /api/admin/users/:id`, `GET /api/admin/analytics`, and the authority-request review routes
- `GET /health`

Report categories are `water`, `air`, `waste`, `deforestation`, and `other`. Report statuses are `open`, `investigating`, `stale`, and `resolved`.

Set `KIJANI_MODE=live`, `KIJANI_BASE_URL`, and `KIJANI_API_KEY` only after KijaniSpace provides the endpoint contract. The live provider deliberately returns a clear configuration error until that contract is implemented.
