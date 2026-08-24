#!/usr/bin/env bash
#
# setup-github-issues.sh
#
# Creates labels, milestones, and issues for the environmental reporting
# platform, matching project-spec.md. Requires the GitHub CLI (`gh`) to be
# installed and authenticated (`gh auth login`).
#
# Usage:
#   ./setup-github-issues.sh <owner>/<repo>
#
# Example:
#   ./setup-github-issues.sh jtonny/enviro-platform

set -euo pipefail

REPO="${1:-}"
if [[ -z "$REPO" ]]; then
  echo "Usage: $0 <owner>/<repo>"
  exit 1
fi

echo "==> Target repo: $REPO"
gh repo view "$REPO" >/dev/null || { echo "Cannot access $REPO — check the name and your gh auth."; exit 1; }

# ---------------------------------------------------------------------------
# 1. LABELS
# ---------------------------------------------------------------------------
# format: "name|color|description"
LABELS=(
  "type:feature|1D76DB|New functionality"
  "type:bug|D73A4A|Something isn't working"
  "type:chore|CFD3D7|Tooling, config, non-feature work"
  "type:docs|0075CA|Documentation"
  "type:test|C2E0C6|Test coverage"
  "area:backend|5319E7|Go/Gin/GORM backend"
  "area:frontend|FBCA04|HTML/CSS/JS frontend"
  "area:auth|B60205|Auth & JWT"
  "area:reports|0E8A16|Reports (pollution posts)"
  "area:comments|BFD4F2|Comments"
  "area:likes|F9D0C4|Likes"
  "area:authority-flow|D4C5F9|Authority self-register/approval flow"
  "area:map|006B75|Kijani map integration"
  "area:admin|E99695|Admin dashboard & moderation"
  "area:devops|000000|Docker, CI, deployment"
  "priority:high|B60205|Blocking or core-path work"
  "priority:medium|FBCA04|Standard priority"
  "priority:low|C5DEF5|Nice-to-have / polish"
)

echo "==> Creating labels..."
for entry in "${LABELS[@]}"; do
  IFS="|" read -r name color desc <<< "$entry"
  gh label create "$name" --repo "$REPO" --color "$color" --description "$desc" --force
done

# ---------------------------------------------------------------------------
# 2. MILESTONES
# ---------------------------------------------------------------------------
# format: "title|description"
MILESTONES=(
  "1. Project Setup & Infrastructure|Repo scaffold, Docker, Makefile, config, migrations tooling"
  "2. Auth & Role Management|JWT auth, role middleware, register/login"
  "3. Reports Core|Report CRUD, photo upload, category, severity, location"
  "4. Comments & Likes|Comment creation, authority-comment flagging, like/unlike"
  "5. Authority Approval Flow|AuthorityRequest submission + admin approval, role flip"
  "6. Status Automation|Auto open->investigating->stale scheduled job"
  "7. Map & Kijani Integration|Satellite layer proxy + report pin overlay, role-gated"
  "8. Admin Dashboard & Moderation|User management, moderation, analytics"
  "9. Frontend Pages & Components|All pages, design tokens, api.js/state.js"
  "10. Testing & Deployment|Backend/frontend tests, docker-compose, docs"
)

echo "==> Creating milestones..."
for entry in "${MILESTONES[@]}"; do
  IFS="|" read -r title desc <<< "$entry"
  gh api "repos/$REPO/milestones" -f title="$title" -f description="$desc" -f state="open" \
    --silent || echo "   (milestone '$title' may already exist, skipping)"
done

# ---------------------------------------------------------------------------
# 3. ISSUES
# ---------------------------------------------------------------------------
# Each issue: title|milestone|labels(comma-sep)|body
# Body includes Description, Acceptance Criteria, Affected Files.

create_issue() {
  local title="$1" milestone="$2" labels="$3" body="$4"
  gh issue create --repo "$REPO" --title "$title" --milestone "$milestone" \
    --label "$labels" --body "$body"
}

echo "==> Creating issues..."

# --- Milestone 1: Project Setup & Infrastructure ---
create_issue \
  "Scaffold Go backend clean-architecture skeleton" \
  "1. Project Setup & Infrastructure" \
  "type:chore,area:backend,priority:high" \
"## Description
Set up the initial Go project structure following clean architecture (handlers -> services -> repositories).

## Acceptance Criteria
- [ ] cmd/server/main.go boots Gin with config loaded from env
- [ ] internal/{handlers,services,repositories,models,middleware,config,router} folders exist with placeholder files
- [ ] go.mod initialized

## Affected Files
- cmd/server/main.go
- internal/config/*
- internal/router/*
- go.mod, go.sum"

create_issue \
  "Set up Docker + docker-compose for local dev" \
  "1. Project Setup & Infrastructure" \
  "type:chore,area:devops,priority:high" \
"## Description
Multi-stage Dockerfile for the Go app plus docker-compose with Postgres for local development.

## Acceptance Criteria
- [ ] Dockerfile builds a minimal alpine image
- [ ] docker-compose.yml runs app + postgres
- [ ] .env.example includes KIJANI_API_KEY, JWT_SECRET, DB creds

## Affected Files
- Dockerfile
- docker-compose.yml
- .env.example"

create_issue \
  "Add Makefile with build/run/test/migrate/docker targets" \
  "1. Project Setup & Infrastructure" \
  "type:chore,area:devops,priority:medium" \
"## Description
Standard Makefile targets matching the team workflow.

## Acceptance Criteria
- [ ] make build, run, test, migrate, docker-build, docker-up, docker-down all work

## Affected Files
- Makefile"

create_issue \
  "Choose and wire up migration tool" \
  "1. Project Setup & Infrastructure" \
  "type:chore,area:backend,priority:medium" \
"## Description
Pick golang-migrate or goose, add initial migration for all core tables (users, reports, comments, likes, authority_requests).

## Acceptance Criteria
- [ ] migrations/ contains up/down SQL for all 5 core tables
- [ ] make migrate applies them

## Affected Files
- migrations/*.sql
- Makefile"

# --- Milestone 2: Auth & Role Management ---
create_issue \
  "Implement User model and migration" \
  "2. Auth & Role Management" \
  "type:feature,area:backend,area:auth,priority:high" \
"## Description
GORM User model: id, name, email, password_hash, role, status, created_at.

## Acceptance Criteria
- [ ] Role enum: public, authority, admin
- [ ] Status enum: active, pending, banned
- [ ] Unique constraint on email

## Affected Files
- internal/models/user.go
- migrations/xxxx_create_users.sql"

create_issue \
  "POST /api/auth/register" \
  "2. Auth & Role Management" \
  "type:feature,area:backend,area:auth,priority:high" \
"## Description
Register endpoint, defaults new users to role=public, status=active.

## Acceptance Criteria
- [ ] Hashes password (bcrypt)
- [ ] Validates unique email
- [ ] Returns created user (no password hash in response)

## Affected Files
- internal/handlers/auth_handler.go
- internal/services/auth_service.go
- internal/repositories/user_repository.go"

create_issue \
  "POST /api/auth/login + JWT issuance" \
  "2. Auth & Role Management" \
  "type:feature,area:backend,area:auth,priority:high" \
"## Description
Login endpoint verifying credentials and issuing a JWT with user_id and role claims.

## Acceptance Criteria
- [ ] JWT includes user_id, role, exp
- [ ] Wrong credentials return 401

## Affected Files
- internal/handlers/auth_handler.go
- internal/services/auth_service.go
- pkg/jwt.go (or internal/config)"

create_issue \
  "JWT + role-based auth middleware" \
  "2. Auth & Role Management" \
  "type:feature,area:backend,area:auth,priority:high" \
"## Description
Middleware validating JWT on protected routes and enforcing role checks per route group (public/authority/admin).

## Acceptance Criteria
- [ ] Invalid/missing token -> 401
- [ ] Insufficient role -> 403
- [ ] Route groups in router/ apply the right guard

## Affected Files
- internal/middleware/auth.go
- internal/middleware/require_role.go
- internal/router/*"

# --- Milestone 3: Reports Core ---
create_issue \
  "Report model, migration, and repository" \
  "3. Reports Core" \
  "type:feature,area:backend,area:reports,priority:high" \
"## Description
Report entity: title, description, photo_url, category, severity (1-5), lat/lng, status, last_activity_at.

## Acceptance Criteria
- [ ] Category enum: water, air, waste, deforestation, other
- [ ] Status enum: open, investigating, stale, resolved
- [ ] Severity constrained 1-5 at the DB level

## Affected Files
- internal/models/report.go
- migrations/xxxx_create_reports.sql
- internal/repositories/report_repository.go"

create_issue \
  "POST /api/reports (create report + photo upload)" \
  "3. Reports Core" \
  "type:feature,area:backend,area:reports,priority:high" \
"## Description
Create a report with photo upload handling and location capture.

## Acceptance Criteria
- [ ] Accepts multipart form with photo
- [ ] Validates category/severity/location present
- [ ] Sets status=open, last_activity_at=now

## Affected Files
- internal/handlers/report_handler.go
- internal/services/report_service.go"

create_issue \
  "GET /api/reports and GET /api/reports/:id" \
  "3. Reports Core" \
  "type:feature,area:backend,area:reports,priority:high" \
"## Description
List reports (with category/status filters) and fetch a single report's detail.

## Acceptance Criteria
- [ ] Supports ?category= and ?status= query params
- [ ] Pagination on list endpoint

## Affected Files
- internal/handlers/report_handler.go
- internal/services/report_service.go
- internal/repositories/report_repository.go"

create_issue \
  "PATCH/DELETE /api/reports/:id (owner or admin)" \
  "3. Reports Core" \
  "type:feature,area:backend,area:reports,priority:medium" \
"## Description
Edit/delete a report, restricted to the report's owner or an admin.

## Acceptance Criteria
- [ ] Non-owner, non-admin gets 403
- [ ] Admin can edit/delete any report (moderation path)

## Affected Files
- internal/handlers/report_handler.go
- internal/services/report_service.go"

# --- Milestone 4: Comments & Likes ---
create_issue \
  "Comment model + POST/GET /api/reports/:id/comments" \
  "4. Comments & Likes" \
  "type:feature,area:backend,area:comments,priority:high" \
"## Description
Comment entity with is_authority_comment flag; adding a comment updates the parent report's last_activity_at.

## Acceptance Criteria
- [ ] is_authority_comment set true when commenter role=authority
- [ ] First authority comment triggers report status -> investigating (see Milestone 6)

## Affected Files
- internal/models/comment.go
- migrations/xxxx_create_comments.sql
- internal/handlers/comment_handler.go
- internal/services/comment_service.go"

create_issue \
  "Like model + POST/DELETE /api/reports/:id/likes" \
  "4. Comments & Likes" \
  "type:feature,area:backend,area:likes,priority:medium" \
"## Description
One like per user per report, enforced with a unique DB constraint on (report_id, user_id).

## Acceptance Criteria
- [ ] Duplicate like attempt returns 409 or is idempotent (decide and document)
- [ ] Unlike removes the row

## Affected Files
- internal/models/like.go
- migrations/xxxx_create_likes.sql
- internal/handlers/like_handler.go"

# --- Milestone 5: Authority Approval Flow ---
create_issue \
  "AuthorityRequest model + POST /api/authority-requests" \
  "5. Authority Approval Flow" \
  "type:feature,area:backend,area:authority-flow,priority:high" \
"## Description
Public users submit an AuthorityRequest (organization_name, justification) to apply for authority status.

## Acceptance Criteria
- [ ] status defaults to pending
- [ ] One pending request per user enforced

## Affected Files
- internal/models/authority_request.go
- migrations/xxxx_create_authority_requests.sql
- internal/handlers/authority_request_handler.go"

create_issue \
  "GET/PATCH /api/authority-requests (admin review)" \
  "5. Authority Approval Flow" \
  "type:feature,area:backend,area:authority-flow,priority:high" \
"## Description
Admin lists pending requests and approves/rejects. Approval flips the requesting user's role to authority.

## Acceptance Criteria
- [ ] Approval updates User.role=authority and AuthorityRequest.status=approved
- [ ] Rejection leaves role=public, status=rejected
- [ ] reviewed_by set to the acting admin

## Affected Files
- internal/handlers/authority_request_handler.go
- internal/services/authority_request_service.go"

# --- Milestone 6: Status Automation ---
create_issue \
  "Scheduled job: auto-flag stale reports after 7 days" \
  "6. Status Automation" \
  "type:feature,area:backend,priority:high" \
"## Description
Background job (ticker/cron in the service layer) scanning reports and setting status=stale when last_activity_at is older than 7 days and status is not already resolved.

## Acceptance Criteria
- [ ] Runs on an interval (e.g. hourly) without blocking request handling
- [ ] Skips resolved reports
- [ ] Covered by a table-driven test with mocked time

## Affected Files
- internal/services/report_status_job.go
- cmd/server/main.go (job registration)"

create_issue \
  "Auto-transition open -> investigating on first authority comment" \
  "6. Status Automation" \
  "type:feature,area:backend,priority:high" \
"## Description
When a comment with is_authority_comment=true is created on a report with status=open, flip the report to investigating.

## Acceptance Criteria
- [ ] Only triggers on the transition into investigating, not on every subsequent authority comment
- [ ] Unit tested in comment_service

## Affected Files
- internal/services/comment_service.go
- internal/services/report_service.go"

create_issue \
  "PATCH /api/reports/:id/resolve (admin-only manual resolution)" \
  "6. Status Automation" \
  "type:feature,area:backend,area:admin,priority:medium" \
"## Description
Since no role auto-resolves reports, admin gets a manual endpoint to mark a report resolved.

## Acceptance Criteria
- [ ] Admin-only (403 otherwise)
- [ ] Sets status=resolved

## Affected Files
- internal/handlers/report_handler.go
- internal/services/report_service.go"

# --- Milestone 7: Map & Kijani Integration ---
create_issue \
  "Kijani API client + GET /api/map/satellite proxy" \
  "7. Map & Kijani Integration" \
  "type:feature,area:backend,area:map,priority:high" \
"## Description
Backend client wrapping the Kijani API, proxied through our own endpoint so the API key never reaches the frontend. Authority/admin only.

## Acceptance Criteria
- [ ] KIJANI_API_KEY read from env, never exposed to client
- [ ] Endpoint gated to authority/admin roles

## Affected Files
- internal/services/kijani_client.go
- internal/handlers/map_handler.go
- internal/config/*"

create_issue \
  "GET /api/map/reports (pins for overlay)" \
  "7. Map & Kijani Integration" \
  "type:feature,area:backend,area:map,priority:medium" \
"## Description
Lightweight endpoint returning report lat/lng + minimal metadata for map pins, authority/admin only.

## Acceptance Criteria
- [ ] Returns only fields needed for pins (id, category, severity, status, lat, lng)

## Affected Files
- internal/handlers/map_handler.go
- internal/repositories/report_repository.go"

create_issue \
  "Frontend: Map page with satellite layer + report pins overlaid" \
  "7. Map & Kijani Integration" \
  "type:feature,area:frontend,area:map,priority:high" \
"## Description
Map page combining the Kijani satellite layer and report pins in a single view, visible only when state.js reports role=authority or admin.

## Acceptance Criteria
- [ ] Non-authorized roles never see this page/nav link
- [ ] Pin click shows report summary + link to detail page

## Affected Files
- frontend/pages/map/index.html
- frontend/pages/map/map.js
- frontend/js/api.js
- frontend/js/state.js"

# --- Milestone 8: Admin Dashboard & Moderation ---
create_issue \
  "GET/PATCH /api/admin/users (role & ban management)" \
  "8. Admin Dashboard & Moderation" \
  "type:feature,area:backend,area:admin,priority:high" \
"## Description
Admin endpoints to list users and update role/status (e.g. ban).

## Affected Files
- internal/handlers/admin_handler.go
- internal/services/admin_service.go"

create_issue \
  "GET /api/admin/analytics" \
  "8. Admin Dashboard & Moderation" \
  "type:feature,area:backend,area:admin,priority:medium" \
"## Description
Aggregate stats: report counts by category/status, active users, trends over time.

## Affected Files
- internal/handlers/admin_handler.go
- internal/services/admin_service.go
- internal/repositories/report_repository.go"

create_issue \
  "Frontend: Admin dashboard, users, authority-requests, moderation pages" \
  "8. Admin Dashboard & Moderation" \
  "type:feature,area:frontend,area:admin,priority:high" \
"## Description
The four admin-only pages: analytics dashboard, user management, authority-request review, and content moderation (edit/delete any report/comment).

## Affected Files
- frontend/pages/admin/dashboard/*
- frontend/pages/admin/users/*
- frontend/pages/admin/authority-requests/*
- frontend/pages/admin/moderation/*
- frontend/js/api.js"

# --- Milestone 9: Frontend Pages & Components ---
create_issue \
  "Design tokens + base styles" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,priority:high" \
"## Description
tokens.css defining the full palette/spacing/type scale for the project; base.css for resets.

## Affected Files
- frontend/css/tokens.css
- frontend/css/base.css"

create_issue \
  "api.js: centralized fetch wrappers for all endpoints" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,priority:high" \
"## Description
Single module wrapping every backend endpoint from the spec (auth, reports, comments, likes, authority-requests, map, admin). No other file makes raw fetch calls.

## Affected Files
- frontend/js/api.js"

create_issue \
  "state.js: auth/role state + UI gating helpers" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,priority:high" \
"## Description
Holds current user/JWT/role and exposes helpers like canAccessMap(), isAdmin() used across pages to gate UI. No DOM manipulation in this file.

## Affected Files
- frontend/js/state.js"

create_issue \
  "Login / Register pages" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,priority:high" \
"## Affected Files
- frontend/pages/login/*
- frontend/pages/register/*"

create_issue \
  "Feed + Report Detail + New Report pages" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,area:reports,priority:high" \
"## Description
Feed with category/status filters, report detail with comments+likes, and the new-report form with photo upload and location picker.

## Affected Files
- frontend/pages/feed/*
- frontend/pages/report-detail/*
- frontend/pages/new-report/*"

create_issue \
  "Apply for Authority page" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,area:authority-flow,priority:medium" \
"## Affected Files
- frontend/pages/apply-authority/*"

create_issue \
  "Reusable components: report-card, comment-thread, severity-badge, category-tag, role-aware nav" \
  "9. Frontend Pages & Components" \
  "type:feature,area:frontend,priority:medium" \
"## Affected Files
- frontend/css/components/*
- frontend/js/components/*"

# --- Milestone 10: Testing & Deployment ---
create_issue \
  "Backend unit tests (services + repositories)" \
  "10. Testing & Deployment" \
  "type:test,area:backend,priority:high" \
"## Description
Table-driven tests for services and repositories, including the status-automation job.

## Affected Files
- internal/services/*_test.go
- internal/repositories/*_test.go"

create_issue \
  "Backend handler tests (httptest + mocked services)" \
  "10. Testing & Deployment" \
  "type:test,area:backend,priority:medium" \
"## Affected Files
- internal/handlers/*_test.go"

create_issue \
  "Frontend module tests" \
  "10. Testing & Deployment" \
  "type:test,area:frontend,priority:low" \
"## Description
Plain assertion-based tests for api.js/state.js/utils.js.

## Affected Files
- frontend/js/*.test.js"

create_issue \
  "Write architecture.md" \
  "10. Testing & Deployment" \
  "type:docs,priority:medium" \
"## Description
Generate the architecture.md doc covering project structure, layering rules, data model, auth flow, business rules, API contract, and conventions.

## Affected Files
- architecture.md"

echo "==> Done. Labels, milestones, and issues created on $REPO."
