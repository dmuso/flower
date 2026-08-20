# Technical approach

Eventual repo path: `docs/core-workflow/technical-approach.md`

Local draft: `/workspace/flower-spec/technical-approach.md`

Change name: `flower`

This is the Technical Lead approach for implementers. It does not replace the product companions. It does not invent product rules.

**Source of truth (product):**

| Topic | Wins |
| --- | --- |
| Packing **math** (leave-short, never-split, oversized Feature, zero-cost, Start overflow, window clock) | `velocity-and-planning.md` |
| No `iterations` table / no `stories.iteration_id` as planning truth | this file + `domain-model.md` + Dan’s law (velocity doc now agrees) |
| Auth, user API tokens, who may mint | this file + `product-spec.md` roles (Owner / Member / Viewer) |
| Organisations, roles, isolation | `multitenancy.md` |
| Slice order and acceptance criteria | `product-spec.md` |
| Copy-exactly vs modernise | `tracker-brief.md` |
| True forks + baked assumptions | `open-questions.md` |
| Directory shape | `STRUCTURE-FROM-LAYOUT.md` (Flower names) + this file |
| Domain boundaries | `domain-model.md` |
| Existing 000001 columns | `000001_create_core_schema.up.sql` is ground: do not rewrite it. Stop using leftover planning columns. |
| Quality bar, no fallbacks, TDD, 90% | repo `AGENTS.md` + `api/internal/migrations/AGENTS.md` |

Cookie sessions and Bearer API tokens call the same `/api/v1` as the frontend. Project role (Owner / Member / Viewer) wins on what they can do. If something is unspecified, it is marked **unspecified** below; an implementation assumption is listed only where a column or locked stack forces a choice.

Do not implement until the Reviewer has cleared the spec set. Do not open a PR from this document. Do not overwrite `docs/reference/*` or `docs/migrations*`. Same-PR product correction of `docs/product/overview.md` is owned by LANDING, not by this file.

Spelling: UK / AU / NZ (`organisations`). Always **user**, never “actor.”

---

## Architecture

### Locked stack and ports

Do not change these.

| Piece | Locked value |
| --- | --- |
| Database | PostgreSQL 17 |
| API | Go + Gin, host **8180** |
| Frontend | SolidJS + Bun + TypeScript + Tailwind v4, host **4273** |
| Local Postgres | **5433** (dev), **5437** (test) |
| Shape | Monorepo, Nix Shell + Docker Compose + Make |
| Look | `docs/reference/frontend-design-guide.md` (bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, four-column board) |
| Host ports already taken | Do **not** bind **8080** or **4173** |

Existing process commands stay: `nix-shell --run 'make dev'`, `make test`, `make lint`, `make migrate`. `.env.example` remains the source of truth for required variables.

### Where slices live

House pattern: `STRUCTURE-FROM-LAYOUT.md`. Domain-driven **vertical slices**. A domain is handler + service + repository + types. Handler *bodies* live in the owning domain. `app/` is process wiring and lifecycle (`Start`, `Shutdown`), not one use-case package per outcome. There is **no** `api/internal/wire/` package. Wiring lives in `app/` and `handlers/` (`handlers/wire.go` is a file, not a layer).

```
flower/
├── api/
│   ├── cmd/                  # process: HTTP, migrate, boot backfills
│   └── internal/
│       ├── domain/           # vertical slices: handler + service + repository + types
│       ├── platform/         # db, clock, auth middleware, email, bus, object store
│       ├── ports/            # cross-domain interfaces only
│       ├── app/              # lifecycle + wiring (Start / Shutdown), not use-case packages
│       ├── handlers/         # route registration for /api/v1/*
│       ├── types/            # minimal shared types (claims)
│       └── migrations/       # schema-only; 000001 is ground; then stop using leftover tables
├── frontend/src/
│   ├── pages/                # route screens
│   ├── components/           # shared chrome
│   ├── features/<slice>/     # thick slices: domain / state / effects / components
│   ├── lib/api/<resource>.ts # one module per HTTP resource; not one api.ts
│   └── routes/
└── docs/
```

Dependency direction: `domain/` → `ports/` ← `platform/`. Domains do not import each other. SQL lives only in the owning domain’s repository files. Platform never contains story, planning, or tenancy rules.

Typical domain package (add when the first slice needs it — do not pre-create empty packages):

```
api/internal/domain/<name>/
├── types.go
├── service.go
├── repository.go      # interface + postgres impl (this is not ports/)
├── handler.go         # HTTP lives here
└── *_test.go
```

Phase 0 domains (same names as `domain-model.md`):

| Slice need | `api/internal/domain/` | Frontend |
| --- | --- | --- |
| Auth / identity / API tokens | `user` | `pages/` login; `lib/api/auth.ts`, `lib/api/users.ts` |
| Organisations | `organisation` | `lib/api/organisations.ts` |
| Projects + membership + invite | `project` | `lib/api/projects.ts`; invite on project page |
| Board read model | `board` | `features/board/`; `lib/api/board.ts` |
| Stories + machine + rank | `story` (`story/machine`, `story/rank`) | `features/story/`; `lib/api/stories.ts` |
| Planning / velocity | `planning` | `features/planning/` (render only) |
| Effective role | inside `organisation` + `project` | — |

**Membership and invite sit with project** (project-scoped). Organisation membership sits with organisation. API tokens sit with `user` (credential of a user, not an identity type).

`app/` and `handlers/` register the new domain’s routes on `/api/v1`. They do not grow a use-case package per slice.

Frontend: `lib/api/core.ts` owns `request`, auth header, errors. Each resource is its own module. Thick slices (`board`, `story`, `planning`) use `features/<name>/{domain,state,effects,components}`. The SPA **does not recompute** velocity or panel membership.

UI Designer owns chords, empty copy, and layout in `ui.md`. Implementers must not invent a palette or a private keymap.

### HTTP mount

Existing repo README exposes `/health`, `/ready`, `/api/version`, and **`/api/v1/*`** on 8180.

**Decision:** one mount, `/api/v1`. Cookie sessions and Bearer API tokens (`Authorization: Bearer flr_...`) call the same resource paths. Example: `POST /api/v1/stories/:id/transitions` with `{ "action": ... }`. Do not run a second tree. Cookie vs Bearer is the only auth difference. Owner / Member / Viewer is the whole permission model. There is no scope list on the token.

SPA origin is `:4273`, API is `:8180` (cross-origin, same-site on localhost). CORS: explicit `FRONTEND_ORIGIN` (never `*`), credentials allowed. Production: put both behind one origin and drop credentialed CORS if possible.

### What must not be disturbed

- `000001` itself: do not rewrite or rename its columns in Phase 0. Additive migrations and Go boot backfills only.
- `stories.rank VARCHAR(64)` fractional / lexicographic. Do not switch to integer priority.
- `stories.title VARCHAR(500)` — product max is 500.
- `story_type` and `state` as strings. No DB enums, checks-as-enums, triggers, or functions.
- `docs/reference/*`, locked look, host ports 8180 / 4273 / 5433 / 5437.
- UK / AU / NZ identifiers: `organisations`, not `organizations`.
- Unique `(project_id, user_id)` on memberships.

Leftover 000001 (stop using / migrate away — **not** the model):

| Leftover | Treat as |
| --- | --- |
| `iterations` table | Unused. Do not insert. Do not read for the plan. Drop in a later schema migration after Go no longer references it. |
| `stories.iteration_id` | Unused. Do not write as planning truth. Do not read to decide Current / Backlog / Done. |
| `projects.iteration_length_weeks` | Unused after backfill. Add `iteration_length_days`; boot backfill `days = weeks * 7`; then read days only. Drop weeks later. |
| `activities.actor_id` | Leftover **column name**. It maps to the user (`users.id`). Do not add an actor type. Prose says user. Prefer `user_id` in new tables and logs. |

New slices **add** tables and columns. They do not redesign 000001. They do stop treating leftover planning columns as the model.

---

## Systems design

### Ground schema (000001) — keep the file, not leftover planning

Exact columns that later slices must treat as given (do not rewrite 000001):

- `users`: `id`, `username` UNIQUE NOT NULL, `email` UNIQUE NOT NULL, `password_hash` (NOT NULL today — see slice 0), `display_name` NOT NULL, timestamps as `TIMESTAMP` (not `TIMESTAMPTZ`).
- `projects`: `id`, `name`, `slug` (UNIQUE globally today — see slice 0), `description`, `point_scale`, `iteration_length_weeks` (leftover; see days).
- `project_memberships`: `role VARCHAR(50)` — product values `owner` | `member` | `viewer` only, enforced in Go.
- `iterations`: leftover table (`number`, `starts_on`, `ends_on`). Not a domain. Not planning truth.
- `stories`: `iteration_id` leftover nullable, `requester_id` RESTRICT, `estimate INTEGER` nullable, `rank VARCHAR(64)`, `accepted_at` nullable, unique `(project_id, rank)` today (becomes `(project_id, rank_list, rank)`).
- `labels.name VARCHAR(100)`, unique `(project_id, name)`.
- `activities`: `kind`, `summary`, leftover `actor_id` RESTRICT (maps to the user), `story_id` nullable SET NULL.

`NULL` estimate = unestimated. `0` is an estimate. Do not store `-1`. Search `estimate:-1` (slice 22) means `estimate IS NULL`.

Existing timestamps are `TIMESTAMP`. **Do not** rewrite 000001 columns to `TIMESTAMPTZ` as a drive-by. Store UTC in `TIMESTAMP`. Convert with the project timezone for display and window end. New tables follow the same convention until a dedicated later migration (not Phase 0) says otherwise.

### Tenancy (organisations)

There is no single-tenant mode. Slice 0 adds the tenant. Isolation rules in `multitenancy.md` are not optional.

**New tables (slice 0):**

```
organisations
  id UUID PK
  name VARCHAR(255) NOT NULL
  created_at, updated_at TIMESTAMP

organisation_memberships
  id UUID PK
  organisation_id -> organisations CASCADE
  user_id -> users CASCADE
  role VARCHAR(50) NOT NULL          -- Go: 'owner' only in MVP
  UNIQUE (organisation_id, user_id)
```

**Add to `projects` (do not rewrite the table):**

- `organisation_id UUID NOT NULL` -> `organisations` RESTRICT (**unspecified** whether org delete is allowed later; RESTRICT so an organisation with projects cannot disappear).
- Drop `idx_projects_slug`. Create `UNIQUE (organisation_id, slug)`.
- `timezone VARCHAR(64) NOT NULL` — velocity doc: store it; default `Australia/Melbourne` until the fork is resolved.
- `velocity_strategy INTEGER NOT NULL` default 3 (allowed 1-4, Go).
- `initial_velocity INTEGER NOT NULL` default 10.
- `iteration_start_weekday INTEGER NOT NULL` default 1 (ISO 1 = Monday).
- `iteration_length_days INTEGER NOT NULL` — the only stored window-length setting. Default **7**. Owner may set any **positive integer** (velocity doc; typical 7, 14, 21, 28). Do not keep a weeks setting in the model.

Boot backfill (Go, not a 000001 rewrite): for existing rows, `iteration_length_days = iteration_length_weeks * 7`. After backfill, **stop reading** `iteration_length_weeks`. Drop that leftover column in a later schema-only migration.

**Slug / URL (open question 2 — unspecified product):** API and SPA routes use UUIDs: `/api/v1/organisations/:organisation_id/projects/:project_id`. No public organisation slug in Phase 0. `projects.slug` stays as the existing column, unique **per organisation**, generated from the project name plus a short suffix on collision. Do not invent a new slug scheme.

**Authn (humans):**

- Email + password **and** magic link (assumption list). SSO later — do not build it.
- Server-side sessions, not JWT:

```
sessions
  id UUID PK
  user_id -> users CASCADE
  token_hash VARCHAR(64) NOT NULL UNIQUE   -- SHA-256 hex of cookie secret
  expires_at TIMESTAMP NOT NULL
  last_project_id UUID NULL -> projects SET NULL
```

Cookie: `flower_session`, HttpOnly, SameSite=Lax, Path=/, Secure in production. Hash the raw token; never store it.

- `users.email_verified_at TIMESTAMP NULL`. Password signup cannot create an organisation until this is set. A magic-link hit on a new email **is** verification.
- `users.password_hash` becomes **nullable** (schema-only). Magic-link-only accounts have `NULL`. Password login on `NULL` is `unauthorized` with a message to use the magic link — that is a real branch, not a fallback hash.

Session login is for users with passwords or magic links. An API token is not a login identity and does not create a session.

**Username (locked 20 Aug 2026):** infer from the email local-part. No username field in slice 0. `users.username` is `NOT NULL UNIQUE`. Slice 0 must persist one: derive from the email local-part (`[a-z0-9_]+`, trim to 100); if taken, append `-` + 4 random unambiguous chars. Allow edit later when a profile slice exists. Do not write an empty string (that is a fallback).

**Display name (unspecified product):** column is `NOT NULL`. Set it to the username at create. Do not invent a profile editor in Phase 0.

**Auth tokens (verify + magic; not API tokens):**

```
auth_tokens
  id UUID PK
  kind VARCHAR(50) NOT NULL          -- verify_email | magic_link
  email VARCHAR(255) NOT NULL
  token_hash VARCHAR(64) NOT NULL UNIQUE
  expires_at TIMESTAMP NOT NULL
  consumed_at TIMESTAMP NULL
```

Single-use. Hash only. Expiry: magic link and verify **unspecified** in product; implementation: 30 minutes. Invite expiry is specified: 14 days (own table).

**Invites (slice 1) — project domain:**

```
project_invites
  id UUID PK
  project_id -> projects CASCADE
  organisation_id -> organisations RESTRICT
  email VARCHAR(255) NOT NULL
  role VARCHAR(50) NOT NULL          -- owner | member | viewer
  token_hash VARCHAR(64) NOT NULL UNIQUE
  invited_by_user_id -> users RESTRICT
  expires_at TIMESTAMP NOT NULL      -- now + 14 days
  accepted_at TIMESTAMP NULL
  revoked_at TIMESTAMP NULL
  created_at TIMESTAMP NOT NULL
```

Resend: revoke (or consume) the old row, insert a new hash. Email already a member: `validation_failed`, no second membership. Accepting joins **that** project only and lists the organisation. Members and Viewers cannot invite.

**Effective role (one function, used by every handler):**

1. If the account is an organisation owner of `project.organisation_id` -> treat as project `owner` even with no `project_memberships` row.
2. Else the `project_memberships.role` for `(project_id, user_id)`.
3. Else: if the project’s organisation is not one they belong to -> **404** on id fetch. If they belong to the org but not the project -> **404** (enumeration: project lists only return projects they can open).
4. Same-tenant, known project, insufficient role on a mutation -> `forbidden` (403). Viewers can GET the board; `POST /stories` is `forbidden`, not 404.

A Bearer token authenticates as `api_tokens.user_id`. Then this **same** function on that user, except if the token has grant rows the effective role on **that project only** is the grant role (already capped at the minter’s role when minted). Token used on a project the user cannot open (and has no grant) -> `not_found` (404). Do not check a scope list.

Cross-tenant id (story, project, attachment, token) -> **404** / `not_found`. Never 403 for “exists in another organisation.”

Organisation owner can do anything a project owner can. Any **Member or Owner** can accept / reject Features (cookie session or Bearer). Viewers cannot. Cookie vs Bearer is the only auth difference. Insufficient role -> `forbidden`.

**Email:** product does not name a vendor. `ports` email interface (implemented in `platform`). Production: SMTP from env. Dev/test: a real **outbox** table written in the same transaction as the token row; a worker or test helper sends/reads it. Tests assert outbox rows. Do not “succeed” without a stored token. Do not log raw links at info level in production (debug in test env only).

### User API tokens

Token belongs to a user. Bearer authenticates as that user. Not org-level. Not a second identity. Each user manages their own (account settings).

```
POST   /api/v1/users/:id/tokens
GET    /api/v1/users/:id/tokens
DELETE /api/v1/users/:id/tokens/:token_id
```

`POST /api/v1/users/:id/tokens/:token_id/revoke` may be the same revoke handler as `DELETE`. No organisation-level mint path.

```
api_tokens
  id UUID PK
  user_id -> users RESTRICT           -- who the token authenticates as
  created_by_user_id -> users RESTRICT
  name VARCHAR(255) NOT NULL
  token_prefix VARCHAR(16) NOT NULL
  token_hash VARCHAR(64) NOT NULL UNIQUE
  revoked_at TIMESTAMP NULL
  created_at TIMESTAMP NOT NULL
```

**Assumption (inherit-memberships, grants allowed):** if the mint body omits projects, the token inherits the user’s full project memberships (no grant rows). If the body lists projects, store a small grant table so a user can mint a reduced-role token for themselves:

```
api_token_grants
  token_id -> api_tokens CASCADE
  project_id -> projects RESTRICT
  role VARCHAR(20) NOT NULL         -- owner | member | viewer
  UNIQUE (token_id, project_id)
```

Each grant role must be at or below the minter’s effective role on that project. Member cannot mint Owner. Viewer cannot mint a Member or Owner grant. Product-spec slice 23 allows a Viewer to create a token on their **own** user; that Bearer is read-only. Dan’s briefing said Viewer cannot mint — **slice 23 AC wins** for that slice (Viewer self-mint, read-only). Prefer inherit when the client does not send grants.

Secret format: `flr_<random>`. Shown **once**. Store hash + prefix only. Revoke is immediate (next call -> `unauthorized`).

**Who may mint:** `:id` is the authenticated user (each user manages their own). Do not mint on another user’s path in Phase 0. Viewer self-mint is read-only as above.

`GET /api/v1/me` with Bearer: the user the token authenticates as (id, name, organisations, projects, **role per project**). Same shape as a cookie session. Optional token `name`. Never a scope list.

A typical CI token named `"CI"` minted for a Member user can start / finish / deliver **and** accept / reject.

Transitions: same `POST /api/v1/stories/:id/transitions` with `{ "action": ... }`. Member / Owner accept / reject **succeeds** (session or Bearer). Viewer -> `forbidden`. Start unestimated Feature -> `unestimated`. Illegal -> `invalid_transition` with `from` and `action`. Cross-tenant -> `not_found`. Revoked token -> `unauthorized`.

Create story: **one API**. The API may default omitted `story_type` to Feature. `requester_id` is the authenticated user and must be a Member / Owner. Omitted panel -> icebox.

Bulk (shared route; session or Bearer): max 50, all-or-nothing, `Idempotency-Key`. Same key + body -> original result. Same key + different body -> `conflict`.

```
idempotency_keys
  key VARCHAR(255) NOT NULL
  project_id UUID NOT NULL
  request_hash VARCHAR(64) NOT NULL
  status_code INTEGER NOT NULL
  response_body JSONB NOT NULL
  created_at TIMESTAMP NOT NULL
  UNIQUE (project_id, key)
```

Token table and handlers live in `domain/user`. They may land with or before the first Bearer test (slice 5 AC). Slice 23 is the independently acceptable mint + webhook story; it does not invent a second mint path.

### State machines

Business rules live in Go. No DB enum, trigger, or function.

One table-driven package: `domain/story/machine`. Input: type, from-state, action, **effective role**, estimate — not a scope list. Output: to-state or a coded error (`invalid_transition`, `unestimated`, `forbidden`, `validation_failed`).

`product-spec.md` wins on verbs. Who = Owner / Member; Viewers forbidden.

**Feature / Bug**

| action | from | to | Who |
| --- | --- | --- | --- |
| schedule | unscheduled | unstarted | Owner, Member |
| icebox | unstarted only | unscheduled | Owner, Member |
| start | unstarted | started | Owner, Member; Feature must be estimated (`0` allowed) |
| start | unscheduled (Icebox) | started | Owner, Member; schedule+start. Lands started in Current (may overflow V). Feature already estimated (`0` allowed). |
| finish | started | finished | Owner, Member |
| deliver | finished | delivered | Owner, Member |
| accept | delivered | accepted | Owner, Member |
| reject | delivered | rejected | Owner, Member; non-empty reason |
| restart | rejected | started | Owner, Member |
| undo | last state-changing activity | previous | Owner, Member (slice 15). Viewer cannot. |

Viewers cannot call any of the above. Insufficient role -> `forbidden`. There is no `unstart` verb.

**Chore:** unscheduled -> unstarted -> started -> accepted. `finish` is the verb; result is `accepted`. No delivered / reject.

**Release:** unscheduled in Icebox; created in Backlog or scheduled -> auto-`started`. `finish` -> `accepted`. No estimate.

`PATCH` must not write `state`. Icebox verb is `icebox`. Undo is `undo` (slice 15): latest state-changing activity for any Owner/Member — not “own last change only.”

**Start Feature** without estimate -> `unestimated`, no mutation. `0` is allowed. Start assigns the clicker as a story owner if owners < 5; if already 5 and clicker is not among them, Start still happens and they are **not** added. Auto-follow the clicker.

**Reject reason (Phase 0):** product calls it a comment; the comments table is slice 13. Phase 0 stores the reason on `activities` (`kind = story.rejected`, `summary` = reason) and returns it on the story payload as `reject_reason` from the latest reject activity. Empty reason -> `validation_failed`, no state change.

**Optimistic concurrency:** add `stories.revision INTEGER NOT NULL` default 1, incremented on every story mutation (fields, rank, state, estimate). Clients send `revision` or `If-Match`. Stale -> `conflict`.

**`started_at`:** add `stories.started_at TIMESTAMP NULL`. Set on the **first** successful `start` only. Reject / restart do not clear it.

**Undo (slice 15, design now):** latest **state-changing** activity only. Reorders do not write activity. Undo itself writes an activity. Prefer additive `from_state` / `to_state` VARCHAR columns on `activities` when slice 15 lands.

Illegal verbs do not change state.

### Planning (calculation, not an aggregate)

The only stored window setting is `projects.iteration_length_days`. The UI draws **bands** from velocity + estimates. Membership of a story in current / next / later is recalculated whenever velocity, order, or estimates change. Stories do not store which band they are in.

Normative fit rules and window clock: `velocity-and-planning.md`. Do not invent a third planner. Do not show a Recalculate button.

**Window clock (from the velocity doc):** calendar windows that end at midnight, project TZ.

- `L = project.iteration_length_days` (positive integer, default 7).
- First start = configured start weekday on or before project-created date, project TZ.
- Window `i`: `starts_on(i) = first_start + i * L days` at 00:00 TZ; `ends_at(i) = starts_on(i) + L days` at 00:00 TZ (half-open `[starts_on, ends_at)`). Displayed end = last calendar day (`ends_at - 1 day`).
- Current window = the unique `i` where `starts_on <= now < ends_at`.
- Changing length / timezone / start weekday replans immediately.

**Velocity history** (calculation input, not a planning aggregate):

```
velocity_samples
  id UUID PK
  project_id -> projects CASCADE
  starts_on DATE NOT NULL
  ends_on DATE NOT NULL            -- last calendar day of the closed window
  accepted_feature_points INTEGER NOT NULL
  created_at TIMESTAMP NOT NULL
  UNIQUE (project_id, starts_on)
```

Stories do **not** point at this table. A row is written at window end only: freeze accepted Feature points for the window that just closed. Deleted Features that were already in a completed total stay in that frozen integer. This is the “velocity observation” the velocity doc names. It is not an Iteration. It is not a home for stories.

**Velocity (MVP), in Go, integer, in `domain/planning`:**

```
N = count(velocity_samples for the project)
K = project.velocity_strategy   # 1-4, default 3
if N == 0:
    V = project.initial_velocity    # default 10
else:
    window = last min(K, N) samples, by starts_on desc
    if every sample has accepted_feature_points == 0:
        V = project.initial_velocity
    else:
        V = floor(mean(accepted_feature_points of window))
```

Accepted Feature points only. Estimate at **accept** time is what is recorded; after window end the frozen integer is the source. Bugs / chores / releases add 0 until the Phase 3 toggle (do not add the toggle column in Phase 0). Team strength % is not MVP.

**`Pack`** (`domain/planning`) is a pure function (velocity doc signature):

```
pack(ordered_stories, velocity, length_in_days, now, timezone) -> bands
```

Each band is `{ starts_on, ends_at, kind: current | next | later, story_ids }`. It does **not** change `rank`. It does **not** write `stories.iteration_id`. It does **not** insert `iterations` rows. Frontend does not pack.

Call Pack **inside the same transaction** as any mutation the velocity doc lists: create, delete, icebox, schedule, reorder, estimate, type change, start and other verbs, accept, reject, restart, undo, window end, V / length / TZ / strategy / initial-velocity / (later) toggles. After commit, publish the bus event. The board payload is the packed bands; it is not a stored assignment.

Start on a Backlog or Icebox story jumps it to the ranked list + Current and **may** overflow V. Drag Backlog -> Current is **not** Start (slice 7).

**Panel membership** (state + `accepted_at` + Pack bands — not a leftover `iteration_id`):

| Panel | Rule |
| --- | --- |
| Icebox | `state = unscheduled`. Own order (`rank_list = icebox`). Not packed. |
| Current | In-flight (`started` / `finished` / `delivered` / `rejected`) **or** `accepted` with `accepted_at` inside the current window **or** `unstarted` (and packed Releases) that Pack puts in the current band. |
| Backlog | Ranked-list stories Pack puts in next / later bands. Visual headers only. |
| Done | `accepted` and `accepted_at` has **aged past the current window**. **Flat list, newest accepted first.** Not grouped by a sample row. |

Accepted-this-window stays in Current until window end. Done is empty until the first window end. Ranked list is Backlog + Current order (`rank_list = ranked`).

**Window end** at `ends_at` (midnight, project TZ):

1. Freeze `accepted_feature_points` = sum of Feature estimates accepted in that window (`starts_on <= accepted_at < ends_at`). Insert `velocity_samples`. Guard with unique `(project_id, starts_on)`.
2. Recompute V.
3. Those accepted stories age into Done (flat list).
4. In-flight stay; leftover unstarted are re-packed. Unaccepted work is not failed and not iceboxed.

Implementation: (a) a process ticker every 30s that selects projects whose current window is due, and (b) the same function at the start of every board read/write. Inject `platform` clock (via `ports`).

**Test clock:** QA must advance past midnight. Production builds do not expose a clock endpoint. Test / `APP_ENV=test` (test compose, Postgres 5437) may expose `POST /api/v1/test/clock` (`{"now": "ISO-8601"}`) compiled behind a build tag or env that **production refuses to start with**. That is a test harness, not a fallback.

`GET /api/v1/projects/:id/board` returns panels, current V, `initial` vs calculated, points/V, current displayed `ends` date, over-velocity badge if Current Feature-points > V, and computed band headers. Frontend does not pack.

Projected delivery date (slice 20): last calendar day of the band Pack assigned. Icebox: “Not scheduled.” Not a field on the story.

### Ranking

Ground: `stories.rank VARCHAR(64)`, unique per project today.

Icebox is its own ordered list (assumption; open question 4). Unique `(project_id, rank)` cannot hold two independent sequences without collisions.

**Add (slice 2), do not change the type of `rank`:**

- `stories.rank_list VARCHAR(16) NOT NULL` — Go values `icebox` | `ranked`.
- Drop `UNIQUE (project_id, rank)`.
- `UNIQUE (project_id, rank_list, rank)`.

New Icebox story goes to the **top** of Icebox. Schedule (slice 3) moves `unscheduled` -> `unstarted`, `rank_list` `icebox` -> `ranked`, rank = **bottom** of the ranked list. Icebox of `unstarted` only: reverse. Cannot icebox started / finished / delivered / rejected / accepted.

**Algorithm:** lexicographic midpoint on charset `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`. Generate between `before` and `after`; empty list -> a mid string (e.g. `U`). If the midpoint would exceed 64 characters, **rebalance** that `(project_id, rank_list)` in the same transaction. Rebalance is an implementation event, not an activity.

**Illegal rank (server, not only UI):** unstarted cannot sit above started / finished / delivered / rejected in the **ranked** list. Reject the write, `illegal_rank`, row snaps back. Accepted-this-window do not return to Backlog by drag. Viewers cannot reorder.

```
POST /api/v1/projects/:project_id/stories/reorder
{ "story_id": "...", "before_id": "..." | "after_id": "...", "revision": 3 }
```

Member or Owner, session or Bearer. Explicit neighbours. Do not default to “top.” Unique-violation or stale revision -> `conflict`. Retry once server-side on unique violation after re-read; then fail.

Reorders do **not** write `activities` (slice 15).

### Real-time (slice 17)

Product: two sessions, same project, mutations appear within **2 seconds** without refresh; focus and comment drafts survive unless that story/comment was deleted; no presence / live cursors; dropped socket shows stale + reconnect; refresh heals; viewers get the same reads.

**Decision:** SSE, not WebSockets (presence is out of scope). Not polling as the primary path.

```
GET /api/v1/projects/:project_id/events
Accept: text/event-stream
```

Authenticated session or Bearer; effective role >= viewer. After each committed mutation, the bus publishes `{project_id, event, story_id, revision}`. MVP bus: **in-process** (`ports` interface, `platform` impl). Document single-API-replica. Do not add an external broker in Phase 0.

Event names align with webhook events where they overlap (`story.created`, `story.updated`, `story.reordered`, `story.started`, ...). Patch by `story_id`; preserve local draft state.

Slice 17 is Phase 1. Phase 0 may ship without SSE; the board remains correct on reload. When slice 17 starts, publish from the same post-commit hook used for webhooks.

### Attachments (slice 14)

Product: `png`, `jpg`, `jpeg`, `gif`, `webp`, <= 10 MB, max 20 per story; clipboard paste; no remote hotlink render; delete -> missing-image on embeds; viewers see, cannot upload; guessed URL serves nothing to the wrong tenant or a signed-out client.

```
attachments
  id UUID PK
  organisation_id -> organisations RESTRICT
  project_id -> projects CASCADE
  story_id -> stories CASCADE
  uploaded_by_user_id -> users RESTRICT
  content_type VARCHAR(100) NOT NULL
  byte_size INTEGER NOT NULL
  storage_key VARCHAR(512) NOT NULL
  created_at TIMESTAMP NOT NULL
```

Serve only via authenticated `GET /api/v1/attachments/:id` (effective role >= viewer on that project). Check organisation **and** project before streaming. Do not put objects on a public bucket URL.

**Storage (unspecified vendor):** `ports` storage interface (`Put`, `Get`, `Delete`), implemented in `platform`. Dev: local directory, keys `{organisation_id}/{project_id}/{story_id}/{attachment_id}`. Production: S3-compatible (env). Same interface, no fallback path.

Validate type by sniffed bytes. Over 10 MB or 21st image -> `validation_failed`.

### Webhooks (outbound; slice 23)

Webhooks are outbound project hooks (not a command API):

```
POST /api/v1/projects/:project_id/webhooks
GET  /api/v1/projects/:project_id/webhooks
POST /api/v1/webhooks/:id/revoke
```

Owner or Member; Viewer cannot. Secret shown once.

```
webhooks
  id UUID PK
  organisation_id, project_id
  url TEXT NOT NULL
  secret_hash VARCHAR(64) NOT NULL
  secret_prefix VARCHAR(16) NOT NULL
  created_by_user_id
  created_at

webhook_deliveries
  id UUID PK
  webhook_id
  event_id UUID NOT NULL
  event VARCHAR(100) NOT NULL
  payload JSONB NOT NULL
  attempt INTEGER NOT NULL
  next_attempt_at TIMESTAMP NULL
  delivered_at TIMESTAMP NULL
  last_status INTEGER NULL
```

Signature: header `X-Flower-Signature: t=<unix>,v1=<hex>` where `v1` is HMAC-SHA256 of `{t}.{raw_body}` with the raw webhook secret. Receivers must reject if `|now - t| > 300s`. At-least-once. Retry non-2xx. Timeout 10s. Delivery is not a command. Worker in the API process.

Events align with SSE names. Product-spec lists a completed-window event (the time box rolling over — a `velocity_samples` insert, not an `iterations` row). Envelope: `event_id`, `organisation_id`, `project_id`, `user_id`, optional `token_name`.

**Rate limits (product: a well-behaved 1 rps single-story transition must not 429).** Limits apply by credential:

| Credential | Limit |
| --- | --- |
| Bearer token (mutations) | 60 / minute / token, burst 10 |
| Bearer token (reads) | 600 / minute / token |
| Human session | 120 / minute |

In-memory per process in MVP. `rate_limited` + 429 + `Retry-After`.

Error envelope for **all** 4xx/409 (session or Bearer) — **one** envelope:

```
{ "error": { "code": "invalid_transition", "message": "...", "from": "unstarted", "action": "finish" } }
```

No 200 with a partial surprise. No coerce (start that silently estimates 1).

### Story owners and follow (needed in slice 4, not only 12)

000001 has `requester_id` only. Do not overload it.

```
story_owners
  story_id -> stories CASCADE
  user_id -> users CASCADE
  UNIQUE (story_id, user_id)
  -- max 5 enforced in Go; sixth -> owners_full

story_followers
  story_id, user_id
  locked BOOLEAN NOT NULL          -- true for requester + owners
  UNIQUE (story_id, user_id)
```

Viewers do not follow. Requester is a Member or Owner.

### Phase 0 HTTP (session cookie or Bearer; same routes)

Auth and tenancy:

- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/magic-link`
- `POST /api/v1/auth/magic-link/consume`
- `POST /api/v1/auth/verify-email/consume`
- `POST /api/v1/auth/logout`
- `GET  /api/v1/me`
- `POST /api/v1/organisations` (verified only; first-run flow may combine org + project)
- `POST /api/v1/organisations/:id/projects` (organisation owner)
- `GET  /api/v1/organisations/:id/projects`
- `GET  /api/v1/projects/:id`
- `GET  /api/v1/projects/:id/board`
- `POST /api/v1/projects/:id/invites`
- `POST /api/v1/invites/:token/accept`

Tokens (user-scoped):

- `POST   /api/v1/users/:id/tokens`
- `GET    /api/v1/users/:id/tokens`
- `DELETE /api/v1/users/:id/tokens/:token_id`

Stories:

- `POST /api/v1/projects/:id/stories` (omitted `story_type` may default to Feature)
- `GET  /api/v1/stories/:id`
- `PATCH /api/v1/stories/:id` (fields, not state)
- `POST /api/v1/stories/:id/transitions`
- `POST /api/v1/projects/:id/stories/reorder`

Board payload internals must include id, title, type, state, estimate, rank, revision, current-window end date, reject_reason if rejected, owners. Do not invent extra product fields. Do not include a stored iteration id.

---

