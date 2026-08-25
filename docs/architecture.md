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
  -> built React assets under `frontend/dist/`
```

The server serves both API routes and frontend assets. Docker Compose is the intended local deployment: one `app` container and one PostgreSQL 16 container.

## 3. Repository Layout

- `cmd/server/main.go`: application composition, database connection/retry, migrations, route registration, static file serving, HTTP startup.
- `internal/config`: environment-based configuration.
- `internal/models`: GORM entities and relationships.
- `internal/repositories`: database access abstractions for users, reports, and authority requests.
- `internal/services`: authentication, report, authority-request, and Kijani provider logic.
- `internal/handlers`: Gin request handlers for auth, reports, comments, likes, status, and assessment.
- `internal/middleware`: CORS, JWT authentication, and role authorization.
- `frontend`: React/Vite application with a centralized API module.
- `docs/spec.md`: product context and Lake Victoria scope.
- `setup-github-issues.sh`: planned milestones, labels, and issue backlog.
- `Dockerfile`, `docker-compose.yml`: container build and local orchestration.
- `migrations/`: versioned initial up/down SQL migrations.

## 4. Layering and Conventions

The intended dependency direction is:

```text
handlers -> services -> repositories -> models/database
```

Authentication, reports, and authority requests follow this rule. The Kijani assessment handler still accesses GORM directly to load its report before invoking the provider.

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

`report_id`, `user_id`, `body`, `is_authority_comment`, and timestamps. The flag is populated from the commenter role. Comment creation updates report activity, and the first authority comment transitions an open report to investigating.

### Like

`report_id`, `user_id`, timestamps, with a composite unique index. The current endpoint toggles a like, but there is no separate DELETE route or explicit error handling around database writes.

### AuthorityRequest

`user_id`, organization name, justification, status, reviewer, and timestamps. Public users can apply, and administrators can list and approve or reject pending requests. Approval changes the applicant's role to authority.

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
- `POST /api/authority-requests`

### Authority/admin restricted

- `PATCH /api/reports/:id/status`
- `PATCH /api/reports/:id/resolve`
- `GET /api/authority-requests`
- `PATCH /api/authority-requests/:id`

The frontend currently consumes auth, report list/detail, report creation, comment, and like operations. It stores the JWT/session in `localStorage` and centralizes requests in `frontend/src/api.js`; shared state and reusable component layers do not exist yet.

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
| Go scaffold and dependency setup | Partial | Builds and has a router package, but no Makefile. |
| Docker/Compose | Partial | Files exist and Compose config parses; runtime requires Docker daemon and image availability. |
| Database schema | Partial | GORM `AutoMigrate` works and an initial up/down migration exists; migration execution is not integrated. |
| Auth/register/login/JWT | Implemented baseline | Core path exists; lacks comprehensive tests and production secret handling. |
| Reports | Partial | List/detail/create and a service layer exist; no photo multipart handling, edit/delete, or category validation. |
| Comments | Partial | Create, authority flagging, activity updates, and the open-to-investigating transition exist; no separate GET endpoint. |
| Likes | Partial | Toggle endpoint exists; explicit unlike contract and robust DB error handling are absent. |
| Authority workflow | Backend implemented | Apply, pending-list, and admin review routes exist; frontend screens are absent. |
| Status automation | Implemented baseline | Authority comments transition open reports to investigating, the hourly job marks reports stale after seven inactive days, and admins can resolve reports manually. |
| Kijani | Mock only | Assessment interface and mock exist; live contract/client and satellite/map endpoints do not. |
| Admin | Not implemented | No admin handlers, analytics, moderation, or frontend. |
| Frontend | Early functional prototype | Auth, feed, detail, comment, and report form pages exist; no shared modules, filters UI, map, admin, or authority application. |
| Tests/CI | Partial | Backend service tests exist; handler, repository, frontend, and CI coverage are absent. |

Overall, the repository is a runnable proof-of-concept for authentication, reports, authority approval, and status automation, not yet the complete environmental reporting platform described by the backlog. The next major backend gaps are broader test coverage and map/admin features.

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
