# RES NAM Architecture

## 1. Purpose and Scope

RES NAM is a Lake Victoria environmental pollution reporting platform. The first release is geographically limited to the pilot bounding box:

- Latitude: `-1.6` to `0.7`
- Longitude: `31.5` to `35.0`

The intended users are the public, including fishermen, maritime authorities, and administrators. Users submit pollution reports; authorities and administrators review and update them; administrators manage users and moderation. KijaniSpace integration is represented by a provider interface and currently runs in mock mode.

This document records the architecture that exists now and the larger target architecture described by `setup-github-issues.sh`. Planned items are explicitly marked so later workflows do not mistake them for implemented functionality.

## 2. Current System Shape

The repository is a Go monolith with a static frontend and PostgreSQL persistence:

```text
Browser
  -> Gin HTTP server (`cmd/server`)
       -> middleware (CORS, JWT, role checks)
       -> handlers (HTTP parsing and responses)
       -> GORM / repositories
       -> PostgreSQL
       -> services (authentication and Kijani assessment)
  -> static files under `frontend/`
```

The server serves both API routes and frontend assets. Docker Compose is the intended local deployment: one `app` container and one PostgreSQL 16 container.

## 3. Repository Layout

- `cmd/server/main.go`: application composition, database connection/retry, migrations, route registration, static file serving, HTTP startup.
- `internal/config`: environment-based configuration.
- `internal/models`: GORM entities and relationships.
- `internal/repositories`: database access abstractions; currently used by authentication and partly independent of report handlers.
- `internal/services`: authentication rules and Kijani provider abstraction/mock implementation.
- `internal/handlers`: Gin request handlers for auth, reports, comments, likes, status, and assessment.
- `internal/middleware`: CORS, JWT authentication, and role authorization.
- `frontend`: plain HTML/CSS/JavaScript pages using browser `fetch` calls.
- `docs/spec.md`: product context and Lake Victoria scope.
- `setup-github-issues.sh`: planned milestones, labels, and issue backlog.
- `Dockerfile`, `docker-compose.yml`: container build and local orchestration.
- `migrations/`: currently empty; schema creation is performed with GORM `AutoMigrate`.

## 4. Layering and Conventions

The intended dependency direction is:

```text
handlers -> services -> repositories -> models/database
```

Current code only partially follows this rule. Authentication uses `AuthService` and `UserRepository`; report handlers access `*gorm.DB` directly, so report business logic and persistence are coupled to HTTP handlers. New work should preserve the direction above and move report operations behind services/repositories before adding substantial features.

Handlers should validate transport input, call a service, and translate domain errors into HTTP responses. Services should own business rules and transaction boundaries. Repositories should own query construction. Models should describe persistence and JSON relationships, not request-specific behavior.

## 5. Implemented Data Model

All models are currently migrated at startup with GORM.

### User

`id`, `name`, unique `email`, `password_hash`, `role`, `status`, timestamps.

Roles currently supported by the code are `public`, `authority`, and `admin`. Status values are intended to be `active`, `pending`, and `banned`, although validation is not centralized.

### Report

`user_id`, `title`, `description`, `photo_url`, `category`, `severity`, `latitude`, `longitude`, `status`, `last_activity_at`, timestamps, plus user/comments/likes relationships.

Creation enforces severity `1..5` and the Lake Victoria pilot bounds. Status defaults to `open`. Category and status are strings rather than database enums.

### Comment

`report_id`, `user_id`, `body`, `is_authority_comment`, and timestamps. The flag is present in the schema but is not currently populated from the commenter role, and comment creation does not update report activity or status.

### Like

`report_id`, `user_id`, timestamps, with a composite unique index. The current endpoint toggles a like, but there is no separate DELETE route or explicit error handling around database writes.

### AuthorityRequest

`user_id`, organization name, justification, status, reviewer, and timestamps. The model exists, but no request/review handlers or routes are implemented.

## 6. Implemented API Surface

### Public

- `GET /health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/reports`
- `GET /api/reports/:id`

### Authenticated

- `POST /api/reports`
- `POST /api/reports/:id/comments`
- `POST /api/reports/:id/like` (toggle behavior)
- `GET /api/reports/:id/assessment`

### Authority/admin restricted

- `PATCH /api/reports/:id/status`

The frontend currently consumes auth, report list/detail, report creation, and comment operations. It stores the JWT/session in `localStorage` and makes raw `fetch` calls per page; centralized `api.js` and `state.js` do not exist yet.

## 7. Authentication and Authorization Flow

Registration validates basic name/email/password requirements, hashes the password with bcrypt, creates a public active user, and returns a 24-hour HS256 JWT. Login verifies the password and active status before issuing the same token shape.

The JWT middleware requires `Authorization: Bearer <token>`, validates HS256 and token validity, then places `user_id` and `role` in Gin context. `RequireRole` allows exact role matches. The JWT secret is read from `JWT_SECRET`; Docker Compose currently supplies the development value `change-me`, which is unsuitable for production.

## 8. Target Workflows From the Issue Script

The script describes the following milestones:

1. Project setup and infrastructure: scaffold, Docker, Makefile, and a real migration tool.
2. Auth and role management: complete model constraints and role middleware.
3. Reports core: CRUD, photo upload, filters, pagination, and owner/admin permissions.
4. Comments and likes: authority flags, activity updates, comment retrieval, like/unlike semantics.
5. Authority approval: public applications and admin review with role changes.
6. Status automation: stale-report job, authority-comment transition, and admin resolution.
7. Map/Kijani: server-side satellite proxy, role-gated map pins, and frontend map.
8. Admin/moderation: user management, analytics, authority requests, and content moderation.
9. Frontend architecture: design tokens, centralized API/state modules, reusable components, and remaining pages.
10. Testing and deployment: unit, handler, frontend tests, and documented deployment workflow.

## 9. Progress Assessment

| Area | State | Assessment |
|---|---|---|
| Go scaffold and dependency setup | Partial | Builds and starts composition, but no router package or Makefile. |
| Docker/Compose | Partial | Files exist and Compose config parses; runtime requires Docker daemon and image availability. |
| Database schema | Partial | GORM `AutoMigrate` works; versioned up/down migrations are absent. |
| Auth/register/login/JWT | Implemented baseline | Core path exists; lacks comprehensive tests and production secret handling. |
| Reports | Partial | List/detail/create exist; no photo multipart handling, edit/delete, category validation, or service layer. |
| Comments | Partial | Create exists; GET, authority flagging, activity updates, and status transition are absent. |
| Likes | Partial | Toggle endpoint exists; explicit unlike contract and robust DB error handling are absent. |
| Authority workflow | Not implemented | Model only. |
| Status automation | Not implemented | No scheduler/job or resolution endpoint. |
| Kijani | Mock only | Assessment interface and mock exist; live contract/client and satellite/map endpoints do not. |
| Admin | Not implemented | No admin handlers, analytics, moderation, or frontend. |
| Frontend | Early functional prototype | Auth, feed, detail, comment, and report form pages exist; no shared modules, filters UI, map, admin, or authority application. |
| Tests/CI | Not implemented | `go test ./...` compiles but reports no tests. |

Overall, the repository is a runnable proof-of-concept for authentication and basic report browsing/submission, not yet the complete environmental reporting platform described by the backlog. The most important architectural gap is the direct GORM access in report handlers, followed by missing migrations, tests, and role workflows.

## 10. Recommended Workflow Order

Before expanding features, establish a reliable baseline:

1. Fix and verify route/static-file composition with a real Docker/PostgreSQL smoke test.
2. Add a Makefile and versioned migrations, or explicitly document why `AutoMigrate` remains temporary.
3. Introduce report/comment/like services and repositories, keeping handlers thin.
4. Add tests for auth, authorization, report validation, ownership, and status transitions.
5. Implement authority approval and status automation before map/admin UI, because those workflows define role and state contracts.
6. Add centralized frontend API/state modules, then build map and admin pages against stable APIs.
7. Replace the Kijani mock only after the external API contract and credential strategy are available.

## 11. Operational Notes and Risks

- Docker Compose uses development credentials and a fixed JWT secret.
- The server retries PostgreSQL connection for 40 seconds, then exits.
- `uploads/` is copied into the image, but there is no implemented upload handler or persistent volume.
- CORS currently allows every origin.
- There are no automated tests, CI workflow, observability, rate limiting, or documented backup strategy.
- The issue setup script is a planning tool; creating its GitHub issues does not indicate that the corresponding code exists.

