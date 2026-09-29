## res_nam

Local Lake Victoria pollution reporting application. The KijaniSpace integration is represented by a provider interface and runs in `mock` mode until the live API contract and credentials are available.

### Run locally

```sh
docker compose up --build
```

Open `http://localhost:8080`. Create an account, submit a report using coordinates in the pilot area (latitude -1.6 to 0.7, longitude 31.5 to 35.0), and view it in the feed.

For local (non-Docker) development the server expects PostgreSQL and applies the SQL
migrations in `migrations/` at startup. A standalone migration command is also available:

```sh
make migrate
```

### API

- `POST /api/auth/register`, `POST /api/auth/login`
- `GET /api/reports` (supports `?category=`, `?status=`, `?page=`, `?limit=`), `GET /api/reports/:id`
- `GET /api/reports/:id/comments`
- Authenticated: `POST /api/reports` (JSON or multipart with a `photo` field), `PATCH /api/reports/:id`, `DELETE /api/reports/:id`, comments, likes (`POST`/`DELETE /api/reports/:id/like`), and `GET /api/reports/:id/assessment`
- Authority/admin: `PATCH /api/reports/:id/status`, `GET /api/map/reports`
- Admin: `PATCH /api/reports/:id/resolve`, `GET /api/admin/users`, `PATCH /api/admin/users/:id`, `GET /api/admin/analytics`, and the authority-request review routes

Report categories are `water`, `air`, `waste`, `deforestation`, and `other`. Report statuses are `open`, `investigating`, `stale`, and `resolved`.

Set `KIJANI_MODE=live`, `KIJANI_BASE_URL`, and `KIJANI_API_KEY` only after KijaniSpace provides the endpoint contract. The live provider deliberately returns a clear configuration error until that contract is implemented.
