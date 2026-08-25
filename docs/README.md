## res_nam

Local Lake Victoria pollution reporting application. The KijaniSpace integration is represented by a provider interface and runs in `mock` mode until the live API contract and credentials are available.

### Run locally

```sh
docker compose up --build
```

Open `http://localhost:8080`. Create an account, submit a report using coordinates in the pilot area (latitude -1.6 to 0.7, longitude 31.5 to 35.0), and view it in the feed.

### API

- `POST /api/auth/register`, `POST /api/auth/login`
- `GET /api/reports`, `GET /api/reports/:id`
- Authenticated: `POST /api/reports`, comments, likes, and `GET /api/reports/:id/assessment`
- Authority/admin: `PATCH /api/reports/:id/status`
- Admin: `PATCH /api/reports/:id/resolve`

Set `KIJANI_MODE=live`, `KIJANI_BASE_URL`, and `KIJANI_API_KEY` only after KijaniSpace provides the endpoint contract. The live provider deliberately returns a clear configuration error until that contract is implemented.
