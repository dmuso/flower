# Technical approach

Eventual repo path: `docs/core-workflow/technical-approach.md`

Change name: `flower`

This is the Technical Lead approach for implementers. It does not replace the product companions. It does not invent product rules.

**Source of truth (product):**

| Topic | Wins |
| --- | --- |
| Packing, velocity, rollover, releases, dates | `velocity-and-planning.md` |
| Auth, API tokens, attribution | this file + `product-spec.md` roles (Owner / Member / Viewer) |
| Organisations, roles, isolation | `multitenancy.md` |
| Slice order and acceptance criteria | `product-spec.md` |
| Copy-exactly vs modernise | `tracker-brief.md` |
| True forks + baked assumptions | `open-questions.md` |
| Intended schema (users, orgs, projects, stories, tokens) | this file + `domain-model.md` |
| Quality bar, no fallbacks, TDD, 90% | repo `AGENTS.md` + `api/internal/migrations/AGENTS.md` |
| Directory names and vertical-slice ownership | this file + `domain-model.md` (house layout) |

If a slice, mock, or this file disagrees with the velocity doc on packing, the velocity doc wins. Cookie sessions and Bearer API tokens call the same `/api/v1` as the frontend. Project role (Owner / Member / Viewer) wins on what they can do. If something is unspecified, it is marked **unspecified** below; an implementation assumption is listed only where a column or locked stack forces a choice.

Do not implement until the Reviewer has cleared the spec set. Do not open a PR from this document. Do not overwrite `docs/reference/*` or `docs/migrations*`. Same-PR product correction of `docs/product/overview.md` is owned by LANDING, not by this file.

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
| Host ports already taken | Do not bind 8080 or 4173 (taken on this host). |

Existing process commands stay: `nix-shell --run 'make dev'`, `make test`, `make lint`, `make migrate`. `.env.example` remains the source of truth for required variables.

### Where slices live

House layout: domain-driven vertical slices. A slice is one business domain cut through HTTP, service, persistence, and tests — not “all tables, then all handlers.”

`domain/` **includes HTTP**. Handler *bodies* live in `domain/<name>/handler.go`. `app/` is process **wiring and lifecycle** only (`Start`, `Shutdown`, bind ports). `handlers/` registers routes on the existing `/api/v1` tree. There is **no** `api/internal/wire/` package. Do not invent one. `handlers/wire.go` is a file, not a layer.

Core logic lives in `api/internal/domain/<domain>`. Do not dump use-cases into `api/internal/app`.

```
flower/
├── api/
│   ├── cmd/                  # process: HTTP, migrate, tools
│   ├── contracts/            # OpenAPI + contract tests against a real server
│   └── internal/
│       ├── domain/           # vertical slices: handler + service + repository + types
│       ├── platform/         # db, clock, auth, middleware, object store, observability
│       ├── ports/            # cross-domain interfaces only
│       ├── app/              # lifecycle + wiring (not use-case packages)
│       ├── handlers/         # Gin route registration for /api/v1/*
│       ├── types/            # auth claims and other minimal shared types
│       └── migrations/       # schema-only, timestamped
├── frontend/src/
│   ├── pages/                # route-level screens
│   ├── components/           # shared chrome + area UI
│   ├── features/<slice>/     # domain / state / effects / components
│   ├── lib/api/<resource>.ts # one module per HTTP resource; barrel index.ts
│   └── routes/
└── docs/
```

Documented dependency direction:

```
domain/ ──────► ports/ ◄────── platform/
    │                              │
    └──────────► types/ ◄──────────┘
```

Constraints:

- Domain packages must not import each other directly; cross-domain contracts go in `ports/`.
- Keep files under 500 lines where practical.
- Wiring packages (`app/`, `handlers/`, bootstrap helpers) are allowed higher coupling.
- Prefer explicit imports over unused dependency injection.
- Add the domain package when the **first slice that needs it** lands. Do not pre-create empty packages.
- SQL lives only in repository files, and only for tables that domain owns.
- Platform never contains story, planning, or tenancy rules.

### Phase 0 domain packages

Add as the slice lands, under `api/internal/domain/` — **not** under `app/`.

| Slice need | `api/internal/domain/` package | Notes |
| --- | --- | --- |
| Auth / identity + user tokens | `user` (or `auth` if login is its own boundary) | Handler + repository + type; session cookie hashing stays in `platform/auth`. API tokens are a user credential on this domain. |
| Organisations | `organisation` | Handler / service / repository / type; bootstrap seeding as a subpackage if needed. |
| Projects + membership | `project` | Same four files; roles live here or in `permission`. Owns iteration **length in days** and other project settings. |
| Stories | `story` | Create, schedule, icebox, estimate, transition. |
| Rank | `story/rank` | Keep with story unless another domain must own it. |
| Planning / velocity | `planning` | Pure `velocity` + `pack`. A calculation. Persists nothing. |
| Tenancy / 404-vs-403 | `tenancy` or inside `organisation` + `permission` | Effective role. |

Typical domain package:

```
api/internal/domain/<name>/
├── type.go / types.go     # entities and value objects
├── service.go             # orchestration
├── repository.go          # Repository interface + postgres implementation
├── handler.go             # Gin handlers for that domain
├── validation.go          # when needed
├── *_test.go
└── <subpackage>/          # only when the root is too large
```

Larger story domain (do not pre-create until the slice lands):

```
api/internal/domain/story/
├── handler.go
├── service.go
├── repository.go
├── types.go
├── machine/               # type-specific verbs (Feature first)
└── rank/                  # fractional rank generate / compare
```

Shared engines written with the first slice that needs them, tested in isolation:

- `domain/story/machine` — type-specific verbs. Feature machine is used in Phase 0. Bug / Chore / Release tables are coded in the **same** package when the Feature machine lands (product: machines are specified now so later slices do not invent a second workflow). UI and create-story only expose Feature until slice 18.
- `domain/story/rank` — fractional `VARCHAR(64)` generate / compare / illegal-rank check.
- `domain/planning` — velocity formula + pack algorithm. Velocity doc is normative. Pack does **not** persist a time-box assignment on the story.
- `domain/tenancy` — organisation scope, effective role, 404-vs-403.

### Frontend map

```
frontend/src/
├── pages/            # Login, board, story, org home, …
├── components/       # shared chrome + area UI
├── features/         # thick slices only
├── lib/
│   └── api/          # one module per resource; core.ts + barrel index.ts
├── routes/
├── guides/
├── styles/
├── types/            # ambient TS, not domain models
└── assets/
```

House rules:

- **`lib/api/` is split by resource.** `core.ts` owns `request`, auth header, retries, errors. `organisations.ts`, `auth.ts`, `users.ts`, `projects.ts`, `stories.ts`, … own that resource’s functions and types. `index.ts` re-exports. Do not collapse this into a single client file.
- **`pages/`** are the screens the router mounts. A page composes components and feature modules; it does not own HTTP verbs.
- **`components/`** is shared and area UI. Not a dumping ground for domain rules.
- **`features/<name>/`** when a slice is thick enough:

```
frontend/src/features/<slice>/
├── domain/           # pure functions, parsing, serialisation (no HTTP)
├── state/            # store + actions
├── effects/          # bootstrap, autosave, hooks
├── components/       # slice-only UI
├── index.ts
└── types.ts
```

Use `features/` for board, story, and planning. Keep auth/session helpers in `lib/` with screens in `pages/`.

| Slice | `pages/` | `features/` | `lib/api/` |
| --- | --- | --- | --- |
| 0 Auth + empty board | `Login.tsx`, board **page** (projection) | `features/board/` view only | `auth.ts`, `organisations.ts`, `projects.ts`, `stories.ts` |
| 1 Invite | project members | `features/project/` or page-only | members / invite module |
| 2–6 Stories | board + story | `features/story/` | `stories.ts` |
| 7 Rank | — | board/story state | same `stories.ts` (reorder endpoint) |
| 8 Planning | board headers | `features/planning/` (render pack fields) | `stories.ts` + pack fields on the story/project payload |
| 23 Tokens | user settings | page-only is enough | `users.ts` (mint / list / revoke own user) |
| 23b Webhooks | project settings | page-only | `webhooks.ts` |

The SPA **does not recompute** velocity. It draws Icebox / Backlog / Current / Done from stories plus pack fields the API already calculated (velocity, band, dates). There is no `board` domain and no `/board` API. UI Designer owns chords, empty copy, and layout in `ui.md`. Implementers must not invent a palette or a private keymap.

### HTTP mount

Existing repo README exposes `/health`, `/ready`, `/api/version`, and **`/api/v1/*`** on 8180.

**Decision:** one mount, `/api/v1`. Cookie sessions and Bearer API tokens (`Authorization: Bearer flr_...`) call the same resource paths. Example: `POST /api/v1/stories/:id/transitions`. Do not run a second `/v1` tree. Cookie vs Bearer is the only auth difference. Owner / Member / Viewer is the whole permission model.

SPA origin is `:4273`, API is `:8180` (cross-origin, same-site on localhost). CORS: explicit `FRONTEND_ORIGIN` (never `*`), credentials allowed. Production: put both behind one origin and drop credentialed CORS if possible.

### What must not be disturbed

- Intended tables: `users`, `organisations`, `projects`, `project_memberships`, `stories`, `labels`, `story_labels`, `activities`, `api_tokens`. New slices add tables and columns. They do not invent an iteration table, or a story-to-window foreign key.
- `stories.rank VARCHAR(64)` fractional / lexicographic. Do not switch to integer priority.
- `stories.title VARCHAR(500)` — product max is 500.
- `story_type` and `state` as strings. No DB enums, checks-as-enums, triggers, or functions.
- `activities.user_id` → `users`. An API token belongs to a `users` row. Activity is attributed to that user. Do not add a second identity column.
- `projects.point_scale`. The only stored window length is `iteration_length_days` (positive integer, default 7).
- Unique `(project_id, user_id)` on memberships.
- `docs/reference/*`, locked look, locked Flower ports.
- UK / AU / NZ identifiers: `organisations`, not `organizations`.

---

## Systems design

### Intended schema (Phase 0 core)

- `users`: `id`, `username` UNIQUE NOT NULL, `email` UNIQUE NOT NULL, `password_hash` nullable (magic-link-only), `display_name` NOT NULL, `email_verified_at`, timestamps as `TIMESTAMP`.
- `organisations`, `organisation_memberships`.
- `projects`: `id`, `organisation_id`, `name`, `slug` unique per organisation, `description`, `point_scale`, `timezone`, `velocity_strategy`, `initial_velocity`, `iteration_start_weekday`, **`iteration_length_days`**.
- `project_memberships`: `role VARCHAR(50)` — `owner` \| `member` \| `viewer` in Go. Unique `(project_id, user_id)`.
- `stories`: `requester_id` RESTRICT, `estimate INTEGER` nullable, `rank VARCHAR(64)`, `rank_list`, `revision`, `started_at`, `accepted_at`. No window foreign key. Unique `(project_id, rank_list, rank)`.
- `labels.name VARCHAR(100)`, unique `(project_id, name)`.
- `activities`: `kind`, `summary`, **`user_id`** RESTRICT, `story_id` nullable SET NULL.
- `api_tokens`: belongs to `user_id`.

There is no `iterations` table. Stories do not store a window id.

`NULL` estimate = unestimated. `0` is an estimate. Do not store `-1`. Search `estimate:-1` (slice 22) means `estimate IS NULL`.

Timestamps are `TIMESTAMP`. Store UTC. Convert with the project timezone. Do not switch the model to `TIMESTAMPTZ` as a drive-by.

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
  organisation_id → organisations CASCADE
  user_id → users CASCADE
  role VARCHAR(50) NOT NULL          -- Go: 'owner' only in MVP
  UNIQUE (organisation_id, user_id)
```

**Add to `projects` (do not rewrite the table):**

- `organisation_id UUID NOT NULL` → `organisations` RESTRICT (or CASCADE if product later allows org delete — **unspecified**; use RESTRICT so an organisation with projects cannot disappear).
- Drop `idx_projects_slug`. Create `UNIQUE (organisation_id, slug)`.
- `timezone VARCHAR(64) NOT NULL` — velocity doc: store it; default `Australia/Melbourne` until the fork is resolved.
- `velocity_strategy INTEGER NOT NULL` default 3 (allowed 1–4, Go).
- `initial_velocity INTEGER NOT NULL` default 10.
- `iteration_start_weekday INTEGER NOT NULL` default 1 (ISO 1 = Monday). Product: default start weekday Monday.
- `iteration_length_days INTEGER NOT NULL` — the only stored window length. Default 7. Owner may set any positive integer (velocity doc; typical 7, 14, 21, 28).

**Slug / URL:** API and SPA routes use UUIDs: `/api/v1/organisations/:organisation_id/projects/:project_id`. No public organisation slug in Phase 0. `projects.slug` stays as the existing column, unique **per organisation**, generated from the project name plus a short suffix on collision. Do not invent a new slug scheme.


The person is a **user** (`users` table, `/api/v1/users/...`).

**Authn (humans):**

- Email + password **and** magic link (assumption list). SSO later — do not build it.
- Server-side sessions, not JWT:

```
sessions
  id UUID PK
  user_id → users CASCADE
  token_hash VARCHAR(64) NOT NULL UNIQUE   -- SHA-256 hex of cookie secret
  expires_at TIMESTAMP NOT NULL
  last_project_id UUID NULL → projects SET NULL
```

Cookie: `flower_session`, HttpOnly, SameSite=Lax, Path=/, Secure in production. Hash the raw token; never store it.

- `users.email_verified_at TIMESTAMP NULL`. Password signup cannot create an organisation until this is set. A magic-link hit on a new email **is** verification.
- `users.password_hash` becomes **nullable** (schema-only). Magic-link-only users have `NULL`. Password login on `NULL` is `unauthorized` with a message to use the magic link — that is a real branch, not a fallback hash.
- `users.last_project_id` optional; session `last_project_id` is enough for “land on my last project” if updated on board load.

Session login is for users with passwords or magic links. An API token is not a login identity and does not create a session.

**Username:** infer from the email local-part (`[a-z0-9_]+`, trim to 100); if taken, append `-` + 4 random unambiguous chars. `users.username` is unique. No username field in slice 0. Editable later when a profile slice exists.

**Display name (unspecified product):** column is `NOT NULL`. Set it to the username at create. Do not invent a profile editor in Phase 0.

**Auth tokens (verify + magic + invite accept):**

```
auth_tokens
  id UUID PK
  kind VARCHAR(50) NOT NULL          -- verify_email | magic_link
  email VARCHAR(255) NOT NULL
  token_hash VARCHAR(64) NOT NULL UNIQUE
  expires_at TIMESTAMP NOT NULL
  consumed_at TIMESTAMP NULL
```

Single-use. Hash only. Expiry: magic link and verify **unspecified** in product; implementation: 30 minutes for magic/verify. Invite expiry is specified: 14 days (own table).

**Invites (slice 1):**

```
project_invites
  id UUID PK
  project_id → projects CASCADE
  organisation_id → organisations RESTRICT
  email VARCHAR(255) NOT NULL
  role VARCHAR(50) NOT NULL          -- owner | member | viewer
  token_hash VARCHAR(64) NOT NULL UNIQUE
  invited_by_user_id → users RESTRICT
  expires_at TIMESTAMP NOT NULL      -- now + 14 days
  accepted_at TIMESTAMP NULL
  revoked_at TIMESTAMP NULL
  created_at TIMESTAMP NOT NULL
```

Resend: revoke (or consume) the old row, insert a new hash. Email already a member: `validation_failed`, no second membership (unique `(project_id, user_id)` already exists). Accepting joins **that** project only and lists the organisation. Members and Viewers cannot invite.

**Effective role (one function, used by every handler):**

1. If the user is an organisation owner of `project.organisation_id` → treat as project `owner` even with no `project_memberships` row.
2. Else the `project_memberships.role` for `(project_id, user_id)`.
3. Else: if the project’s organisation is not one they belong to → **404** on id fetch. If they belong to the org but not the project → **404** (enumeration: project lists only return projects they can open).
4. Same-tenant, known project, insufficient role on a mutation → `forbidden` (403). Viewers can GET stories; `POST /stories` is `forbidden`, not 404.

A Bearer token authenticates as `api_tokens.user_id`. Then this same function. If the token carries an optional role cap, effective role on a project is the **min** of that cap and the user’s membership (or org-owner elevation). The cap was already required to be at or below the minter when minted. Token used on a project the user cannot open → `not_found` (404), same as cross-tenant. Do not check a scope list.

Cross-tenant id (story, project, attachment, token) → **404** / `not_found`. Never 403 for “exists in another organisation.”

Organisation owner can do anything a project owner can. Any **Member or Owner** can accept / reject Features (cookie session or Bearer). Viewers cannot. Cookie vs Bearer is the only auth difference. Insufficient role → `forbidden`. Do not add an accept ACL or extra error codes. Do not check a scope list.

**Email:** product does not name a vendor. `platform/email` is an interface. Production: SMTP from env. Dev/test: a real **outbox** table written in the same transaction as the token row; a worker or test helper sends/reads it. Tests assert outbox rows. Do not “succeed” without a stored token. Do not log raw links at info level in production (debug in test env only).

### API tokens (user-scoped)

A generic **API token** (Bearer `flr_...`) is a **user credential**, not an identity type. It authenticates as that user. Humans can also use a session cookie. Same `/api/v1` handlers. What the caller can do is Owner / Member / Viewer. The token does not add scopes.

Mint, list, and revoke on **that user**. There is no organisation-collection token API and no org-level grants table.

```
GET    /api/v1/users/:id/tokens
POST   /api/v1/users/:id/tokens
{ "name": "CI", "role": "member" }

DELETE /api/v1/users/:id/tokens/:token_id
```

`:id` is the authenticated user. Mint / list / revoke only on your own user. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. No organisation-collection token API.

```
api_tokens
  id UUID PK
  user_id → users RESTRICT            -- who the token authenticates as
  organisation_id → organisations RESTRICT   -- bound to one organisation
  created_by_user_id → users RESTRICT -- the minter (always the same user)
  name VARCHAR(255) NOT NULL
  role VARCHAR(20) NULL               -- optional cap: owner | member | viewer
  token_prefix VARCHAR(16) NOT NULL
  token_hash VARCHAR(64) NOT NULL UNIQUE
  revoked_at TIMESTAMP NULL
  created_at TIMESTAMP NOT NULL
```

Role:

- If `role` is omitted, the token uses that user’s project memberships (and org-owner elevation) as-is.
- If `role` is set, it is a cap at or below the minter’s effective role. Member cannot assign `owner`. The token cannot exceed that user’s own memberships.
- Effective role on a request = the user’s membership (or org-owner elevation), then the optional cap. No scope list.

Returns the token (id, name, user, optional role) and the raw secret **once**. List shows prefix + name, never the secret. Revoke is immediate (next call → `unauthorized`). Secret format: `flr_<random>`. Store hash only.

Activity is attributed to the user the token authenticates as (`activities.user_id`). Optional token `name` may appear as an activity label.

`GET /api/v1/me` with Bearer: the user the token authenticates as (id, name, organisation, projects, **role per project**). Same shape as a cookie session. Never a scope list.

Typical CI token named `"CI"` with cap **Member** (or a Member user’s memberships) can start / finish / deliver **and** accept / reject. Viewer cap is read-only.

Token handlers live in `domain/user` (slice 23). Webhooks are slice 23b, a separate story. Bearer auth must work as soon as transitions are tested.

Create story: **one API**. Omitted `story_type` may default to Feature. `requester_id` is the authenticated user and must be a Member / Owner. Default icebox if panel omitted.

Bulk (session or Bearer): max 50, all-or-nothing, `Idempotency-Key`. Table `idempotency_keys` unique `(project_id, key)`.

### State machines

Business rules in Go. No DB enum. One table-driven package: `domain/story/machine`. Input: type, from-state, action, effective role, estimate. Output: to-state or `invalid_transition` / `unestimated` / `forbidden` / `validation_failed`.

`product-spec.md` wins on verbs. Who = Owner / Member. Viewers forbidden.

**Feature / Bug**

| action | from | to | Who |
| --- | --- | --- | --- |
| schedule | unscheduled | unstarted | Owner, Member |
| icebox | unstarted only | unscheduled | Owner, Member |
| start | unstarted | started | Owner, Member; Feature needs estimate (`0` allowed) |
| start | unscheduled (Icebox) | started | Owner, Member; schedule+start; lands started in Current (may overflow velocity) |
| finish | started | finished | Owner, Member |
| deliver | finished | delivered | Owner, Member |
| accept | delivered | accepted | Owner, Member |
| reject | delivered | rejected | Owner, Member; non-empty reason |
| restart | rejected | started | Owner, Member |
| undo | last state-changing activity | previous | Owner, Member (slice 15) |

No `unstart`. `PATCH` must not write `state`. Icebox Start of an estimated Feature lands `started` in Current.

**Chore:** unscheduled -> unstarted -> started -> accepted. `finish` is accept. **Release:** schedule -> auto-`started`; `finish` -> `accepted`. No estimate.

Start Feature without estimate -> `unestimated`. Start assigns the clicker as owner if owners < 5; else Start still happens and they are not added. Auto-follow the clicker.

Reject reason in Phase 0: `activities` (`kind = story.rejected`, `summary` = reason). Comments table is slice 13.

Add `stories.revision` (optimistic concurrency) and `stories.started_at` (first successful start only; reject/restart do not clear it).

Undo (slice 15): latest state-changing activity only. Reorders do not write activity. Prefer `from_state` / `to_state` columns on `activities` when that slice lands.

### Planning (calculation, not an aggregate)

The only stored settings are `projects.iteration_length_days`, start weekday, timezone, `velocity_strategy`, and `initial_velocity`. Do not persist velocity, window totals, or accepted points. We accept stories (`accepted_at`), not points.

Windows and bands are computed and drawn in the UI. Recompute whenever stories, estimates, accepts, or settings change. Stories do not store which band they are in.

Normative formula text: `velocity-and-planning.md`. This file states the implementer shape only. No Recalculate button.

**Window clock:** calendar windows that end at midnight, project timezone. Length is `iteration_length_days`. First start = configured weekday on or before project-created date. Each window is half-open `[starts_on, ends_at)` where `ends_at = starts_on + iteration_length_days` at 00:00 in that timezone. Current window = the unique window containing `now`. Changing length / timezone / weekday replans immediately. **Window end** is a clock crossing, not a write.

**Velocity** (live, `domain/planning`): calculated from previous Features’ **start/end datetimes** (`started_at` → `accepted_at`). Not a sum of estimates in a window. How those durations roll up, and how velocity_strategy applies, is the velocity doc.

A story accepted in the **open** window is not in the lookback. Velocity stays undefined (`pack` uses `initial_velocity` 10) until at least one corpus Feature exists in a **completed** window. When `now` crosses `ends_at`, those accepted Features join the lookback if they fall in the last number of completed windows set by `velocity_strategy`; then velocity = work / time and `pack` uses `predicted_duration(estimate)`.

**predicted_duration** is a size → predicted-duration map, built from completed Features of the **same estimate**. Unused on cold start. Pack fills the current window using velocity plus those predicted durations, in rank order:

```
pack(ordered_stories, velocity, predicted_duration, iteration_length_days, now, timezone) → bands
```

`started_at` is set on the first successful Start only; reject / restart do not clear it (same as cycle time).

Call Pack when serving stories and after mutations that change stories, estimates, accepts, or settings. After commit, publish the bus event. Inject platform clock. Test / `APP_ENV=test` (Postgres 5437) may expose `POST /api/v1/test/clock`. Production refuses to start with that harness.

Start on a Backlog or Icebox story jumps it to Current and may overflow the current window. Drag Backlog -> Current is not Start.

**Panels** (derived in the UI from pack + state):

| Panel | Rule |
| --- | --- |
| Icebox | `state = unscheduled`. Own order (`rank_list = icebox`). |
| Current | In-flight or accepted-this-window or unstarted (and packed Releases) whose pack band is current. |
| Backlog | Ranked-list stories in next / later bands. Visual headers only. |
| Done | Accepted stories whose `accepted_at` is before the current window. Flat list, newest accepted first. |

Accepted-this-window stays in Current until the clock crosses `ends_at`. Done is empty until that first crossing.

There is **no** `/board` API. `GET /api/v1/projects/:id/stories` (and project settings) returns stories plus computed pack fields: band, velocity, initial vs calculated, current window dates, over-capacity badge. The SPA draws four columns. One formula in `domain/planning`.

### Ranking

Ground: `stories.rank VARCHAR(64)`, unique per project today.

Icebox is its own ordered list. Unique `(project_id, rank)` cannot hold two independent sequences without collisions.

**Add (slice 2), do not change the type of `rank`:**

- `stories.rank_list VARCHAR(16) NOT NULL` — Go values `icebox` \| `ranked`.
- Drop `UNIQUE (project_id, rank)`.
- `UNIQUE (project_id, rank_list, rank)`.

Icebox order is independent. New Icebox story goes to the **top** of Icebox. Schedule (slice 3) moves `unscheduled` → `unstarted`, `rank_list` `icebox` → `ranked`, rank = **bottom** of the ranked list. Icebox of `unstarted` only: reverse. Cannot icebox started / finished / delivered / rejected / accepted.

**Algorithm:** lexicographic midpoint on charset `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`. Generate between `before` and `after`; empty list → a mid string (e.g. `U`). If the midpoint would exceed 64 characters, **rebalance** that `(project_id, rank_list)` in the same transaction. Rebalance is an implementation event, not an activity.

**Illegal rank (server, not only UI):** unstarted cannot sit above started / finished / delivered / rejected in the **ranked** list. Reject the write, `illegal_rank`, row snaps back. Accepted-this-window do not return to Backlog by drag. Viewers cannot reorder.

API (Member or Owner, session or Bearer):

```
POST /api/v1/projects/:project_id/stories/reorder
{ "story_id": "...", "before_id": "..." | "after_id": "...", "revision": 3 }
```

The shared API requires explicit neighbours. Do not default to “top.” Unique-violation or stale revision → `conflict`. Retry once server-side on unique violation after re-read; then fail.

Reorders do **not** write `activities` (slice 15).

### Real-time (slice 17)

Product: two sessions, same project, mutations appear within **2 seconds** without refresh; focus and comment drafts survive unless that story/comment was deleted; no presence / live cursors; dropped socket shows stale + reconnect; refresh heals; viewers get the same reads.

**Decision:** SSE, not WebSockets (presence is out of scope). Not polling as the primary path.

```
GET /api/v1/projects/:project_id/events
Accept: text/event-stream
```

Authenticated session or Bearer; effective role ≥ viewer. After each committed mutation, `platform/bus` publishes `{project_id, event, story_id, revision}`. The handler fans out to subscribers of that project.

MVP bus: **in-process**. Document single-API-replica for live updates. Interface it (`Subscribe(projectID)`) so a Redis/NATS adapter can land later without changing handlers. Do not add Redis in Phase 0.

Event names align with webhook events where they overlap. Payload is enough to patch or to trigger a targeted GET. Patch by `story_id`; preserve local draft state.

Client: `EventSource` (or fetch-stream) with cookie. On error: banner “stale”, reconnect with backoff. Manual refresh heals.

Slice 17 is Phase 1. Phase 0 may ship without SSE; the board remains correct on reload. When slice 17 starts, publish from the same post-commit hook used for webhooks.

### Attachments (slice 14)

Product: `png`, `jpg`, `jpeg`, `gif`, `webp`, ≤ 10 MB, max 20 per story; clipboard paste; no remote hotlink render; delete → missing-image on embeds; viewers see, cannot upload; guessed URL serves nothing to the wrong tenant or a signed-out client.

```
attachments
  id UUID PK
  organisation_id → organisations RESTRICT
  project_id → projects CASCADE
  story_id → stories CASCADE
  uploaded_by_user_id → users RESTRICT
  content_type VARCHAR(100) NOT NULL
  byte_size INTEGER NOT NULL
  storage_key VARCHAR(512) NOT NULL
  created_at TIMESTAMP NOT NULL
```

Serve only via authenticated `GET /api/v1/attachments/:id` (effective role ≥ viewer on that project). Check organisation **and** project before streaming. Do not put objects on a public bucket URL. Do not render `http(s):` images that are not this attachment id.

**Storage (unspecified vendor):** `platform/storage` interface (`Put`, `Get`, `Delete`). Dev: local directory, keys `{organisation_id}/{project_id}/{story_id}/{attachment_id}`. Production: S3-compatible from env. Same interface, no “if S3 fails write to /tmp” fallback.

Validate type by sniffed bytes, not only extension. Over 10 MB or 21st image → `validation_failed`.

### Outbound webhooks (slice 23b)

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

Signature: header `X-Flower-Signature: t=<unix>,v1=<hex>` where `v1` is HMAC-SHA256 of `{t}.{raw_body}` with the raw webhook secret. Receivers must reject if `|now - t| > 300s`. At-least-once. Retry non-2xx. Timeout 10s. Delivery is not a command; ignore reply body. Worker in the API process.

Events align with SSE names where they overlap (`story.created`, `story.updated`, `story.reordered`, `story.started`, `story.finished`, `story.delivered`, `story.accepted`, `story.rejected`, and later comment/attachment events). Envelope: `event_id`, `organisation_id`, `project_id`, `user_id` (the user), optional `token_name`. A window crossing is a clock event, not a row insert.

**Rate limits (product: a well-behaved 1 rps single-story transition must not 429).**

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

`stories.requester_id` is the requesting user. Do not overload it.

```
story_owners
  story_id → stories CASCADE
  user_id → users CASCADE
  UNIQUE (story_id, user_id)
  -- max 5 enforced in Go; sixth → owners_full

story_followers
  story_id, user_id
  locked BOOLEAN NOT NULL
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
- `POST /api/v1/organisations`
- `POST /api/v1/organisations/:id/projects`
- `GET  /api/v1/organisations/:id/projects`
- `GET  /api/v1/projects/:id`
- `GET  /api/v1/projects/:id/stories` (includes pack fields; SPA draws the board)
- `POST /api/v1/projects/:id/invites`
- `POST /api/v1/invites/:token/accept`

Stories:

- `POST /api/v1/projects/:id/stories`
- `GET  /api/v1/stories/:id`
- `PATCH /api/v1/stories/:id` (fields, not state)
- `POST /api/v1/stories/:id/transitions`
- `POST /api/v1/projects/:id/stories/reorder`

User tokens (slice 23):

- `GET    /api/v1/users/:id/tokens`
- `POST   /api/v1/users/:id/tokens`
- `DELETE /api/v1/users/:id/tokens/:token_id`

Story payload must include id, title, type, state, estimate, rank, revision, pack band, current-window `ends_on`, reject_reason if rejected, owners. Do not invent extra product fields. Do not add a `/board` resource.

## Refactor opportunities

Only what a slice cannot ship without. No drive-by.

| Must take | Why | Slice |
| --- | --- | --- |
| `projects.slug` unique per organisation, not globally | Tenancy; existing index is global | 0 |
| `users.password_hash` nullable | Magic-link-only users; `NOT NULL` plus a dummy hash would be a fallback | 0 |
| `UNIQUE (project_id, rank)` → `UNIQUE (project_id, rank_list, rank)` | Two lists, one `VARCHAR(64)` column | 2 |
| Add `organisation_id` / timezone / velocity settings / `iteration_length_days` on `projects` | Missing vs product; do not recreate `projects` | 0 / 8 |
| Do not add an iteration table or a story window FK | Planning is a calculation | 8 |
| Add `revision`, `rank_list`, `started_at` on `stories` | Concurrency, two lists, cycle-time later | 2 / 4 |
| Add `email_verified_at` on `users` | Slice 0 AC | 0 |

**Do not take in these slices:**

- Redesign the intended tables or invent an iteration entity.
- Switch `rank` to integer or `bigint` priority.
- `TIMESTAMP` → `TIMESTAMPTZ` as a drive-by.
- DB enums / triggers / functions for state, role, or type.
- GraphQL.
- A second workflow engine or `planned` state (only if the manual-planning slice needs it).
- Extracting a microservice, adding Redis, or adding a queue product for Phase 0.
- Changing Flower ports (8180 / 4273 / 5433 / 5437) or the locked look.
- Inventing `api/internal/wire/` or dumping use-cases into `app/`.
- Persisting Current / next / later as rows, or adding a story window foreign key.

---

## Technical risks

| Risk | Why it is real | Mitigation |
| --- | --- | --- |
| Planner drift from the velocity doc | Easy to “improve” leave-short or oversized Features | Pure `Pack` + table tests for **all three** worked examples in the velocity doc; QA script in that file; velocity doc wins reviews |
| Inventing an iteration table | Easy to persist bands or velocity | Pack and velocity are functions. Persist nothing. Tests assert no window FK |
| Clock vs project timezone | Window end is midnight in the project timezone | Injected `Clock`; test env clock endpoint; crossing is a read, not a write |
| Rank unique collisions / 64-char overflow | Two lists; midpoint can grow | `rank_list` in the unique key; rebalance in-transaction; tests at the 64-char edge |
| Cross-tenant leaks (id guess, search, files, tokens) | Multitenant from day one | Every query joins organisation; miss → 404; attachment GET checks org+project+auth; isolation tests ride slices 0, 1, 23 |
| Token treated as a second identity | Activity needs a user; a token is not a login | Token authenticates as `api_tokens.user_id`; `activities.user_id` is that user; optional token `name` in activity payload/summary; session login stays password/magic-link |
| Session cookie on :4273 → :8180 | Cross-origin | Explicit CORS origin + credentials; prod same-origin proxy |
| SSE lost across API replicas | In-process bus | Single replica in MVP; `Bus` interface; do not pretend it is multi-node |
| Webhook retry storms | At-least-once | Cap attempts (implementation: 8, exponential backoff); `event_id` for receiver idempotency; do not POST back as a user |
| Magic link / invite token leak | Email and logs | Store hashes only; single-use; 14-day invite; do not log raw tokens |
| Clock skew vs project timezone | Rollover definition is midnight **project** timezone | Store timezone; default `Australia/Melbourne`; all comparisons via `Clock` + that timezone |
| Dummy `password_hash` for magic-link users | Hidden branch, forbidden by AGENTS.md | Nullable column; password login refuses NULL |
| Greenfield org on project | Slice 0 | `projects.organisation_id` is NOT NULL in the intended schema |
| Unique rank violation under concurrent drag | Two Members reorder | `revision` + one retry; then `conflict`; UI snaps back |
| File store credentials / public buckets | Tenant files | Auth-only GET; key prefixed by organisation_id; no public ACL |
| Rate limiter per process | Two replicas double the limit | Accept in MVP; document; do not 429 a well-behaved 1 rps single-story transition |
| Org-level token API creeps back | Easy to hang tokens off the tenant collection | Mint / list / revoke only on `/api/v1/users/:id/tokens`. Role = memberships or a cap ≤ minter |

---

## Test plan

TDD: failing test first. Tests use the same production logic (no fake planner, no fake machine). Inject `Clock`, email outbox, storage, bus.

QA tests each slice from its acceptance criteria **alone**, plus the companion rule docs named below. Fail without a meeting.

### Slice 0 — signup, organisation, empty board

**API / domain**

- Unverified password signup: org create is rejected; no `organisations` row.
- Verified password signup: organisation + first project; creator is organisation owner and project owner; `point_scale = linear`, `iteration_length_days = 7`, `initial_velocity = 10`, `timezone = Australia/Melbourne` (until fork).
- Magic link to a new email: user created, `email_verified_at` set, same org+project flow.
- Login password and login magic link both land on last project.
- Story list on a new project: empty; SPA shows four empty columns; Current header shows `0 / 10` (points / initial velocity); no fake stories, no fake dates.
- Reload: same organisation / project ids.
- Cross-tenant: session A `GET` project/story id from org B → 404.

**QA**

- Stranger signs up, names organisation `Acme` and project `Trail`, sees Icebox / Backlog / Current / Done empty (Current may show dates + `0 / 10`).
- Sign out, sign in, same empty board. Fail if another tenant’s name appears.
- Fail if a new palette appears (frontend guide).
- Fail if unverified password user can create an organisation.

### Slice 1 — invite

**API**

- Owner invites `alex@example.com` as `member`: outbox has one email; pending invite listed; token hashed.
- New email accept: signup, lands as Member.
- Existing user accept: project in their list.
- Invite as `viewer`: story list works; `POST /stories`, invite, settings → `forbidden`.
- Email already a member → visible error, one membership row.
- Revoke → consume fails. Expiry 14 days. Resend invalidates old hash.
- Member or Viewer `POST` invite → `forbidden`.
- Isolation: invite accept does not grant a second project.

**QA**

- Two humans, one empty project. Viewer cannot create a story in UI **or** via API.

### Slice 2 — Feature in Icebox

**API**

- Member/Owner create: `story_type=feature`, `state=unscheduled`, `estimate IS NULL`, `requester_id=me`, `rank_list=icebox`, rank at top. Not packed. No projected date.
- Empty title → `validation_failed`. Title 501 chars → rejected (column 500).
- Viewer create → `forbidden`.
- Second create sits above the first in Icebox.
- Current / Backlog / Done unchanged. Planner ignores Icebox.

**QA**

- Default add lands in Icebox. Fail if it appears in Current or affects velocity.

### Slice 3 — schedule / icebox

**API**

- Schedule: `unscheduled` → `unstarted`, leaves Icebox, `rank_list=ranked`, rank at **bottom**. Pack may place it in Current or a later band.
- Icebox: `unstarted` → `unscheduled`, drops out of pack, no projected date.
- Icebox of started / finished / delivered / rejected / accepted → `invalid_transition` / `illegal_rank` as appropriate; no change.
- Viewer move → `forbidden`.

**QA**

- Pull one Feature to Backlog; Icebox empty state returns. Fail if schedule Starts the story.

### Slice 4 — estimate and start

**API**

- Estimate 0 / 1 / 2 / 3; stays `unstarted`. Other values → `validation_failed`.
- Start unestimated Feature → `unestimated`, still `unstarted`.
- Start estimated: `started`, `started_at` set once, lands in Current (overflow allowed), clicker in `story_owners`, locked follower.
- Clear estimate only while `unstarted`. Started Feature cannot become NULL.
- Icebox Start (after estimate): `started` in Current, not `unstarted` in Backlog. This is schedule+start. There is no `unstart`.
- Viewer estimate/start → `forbidden`.
- Sixth distinct owner on Start: Start succeeds, clicker not added (`owners_full` only when they **try** to add a sixth via assign — Start does not error).

**QA**

- Fail the slice if Start on unestimated succeeds. Fail if 0 is treated as unestimated. Fail if Icebox Start lands as `unstarted` in Backlog.

### Slice 5 — finish, deliver, accept

**API**

- started → finish → `finished`, still Current.
- finished → deliver → `delivered`, still Current.
- Any Member or Owner accept → `accepted`, `accepted_at` set, still Current, Done empty.
- Viewer finish/deliver/accept → `forbidden` (session or Bearer).
- Member or Owner accept Feature (session or Bearer) → `accepted`. Viewer Bearer accept → `forbidden`, no change.
- finish on unstarted, accept on finished → `invalid_transition`, no change.
- Tasks/blockers do not exist yet; do not block accept.

**QA**

- Accepted story still in Current. Fail if it jumps to Done before window end.

### Slice 6 — reject and restart

**API**

- delivered + reject + reason → `rejected`, still Current; reason on activity.
- Empty reason → no change, `validation_failed`.
- Viewer reject (session or Bearer) → `forbidden`. Member or Owner reject with reason (session or Bearer) → `rejected`.
- Restart → `started`, still Current; reason remains in activity.
- Finish + deliver again; new accept required.
- Fail if reject jumps to `started` with no Restart, or if `rejected` is treated as terminal like `accepted`.

**QA**

- Rejected work is not in Backlog or Icebox.

### Slice 7 — reorder

**API**

- Member/Owner reorder persists (reload).
- Unstarted above started / finished / delivered / rejected → `illegal_rank`, snap back.
- Accepted-this-window cannot be dragged to Backlog.
- Ranked-list drop into Current does **not** Start and does not write `state`.
- Icebox reorder independent; Icebox → Backlog is schedule (slice 3), not Start.
- Viewer reorder → `forbidden`.
- Keyboard reorder path exists (chord from UI Designer; fail if none).

**QA**

- Drag unstarted above started must fail. Fail if drag-to-Current starts the story.

### Slice 8 — auto-plan

Normative tests = velocity doc worked examples 1–3 **and** its QA short script.

**API / domain (table-driven)**

- New project velocity = 10. Five estimated Features totalling > 10: Current **short**, not over. Next Feature that would exceed 10 stays in the next Backlog band.
- Never split.
- Start a Backlog Feature that did not fit → Current, points may exceed 10, `Over capacity` badge.
- Accept one Feature → still Current, not Done.
- Advance test clock past window `ends_at` project timezone → that Feature in Done (flat list, newest accepted first). Accepted Features join the lookback if they fall in the last number of completed windows set by `velocity_strategy`; then velocity = work / time and pack uses `predicted_duration(estimate)`.
- Reorder / estimate / accept / start / icebox / length change → board already recomputed. No window id written on the story.
- Owner sets `iteration_length_days` to 14 → dates and pack change; velocity recomputes from durations.
- Icebox never auto-fills into Current.
- Oversized Feature (predicted_duration(estimate) exceeds the current window): next empty window, over-capacity exception (velocity doc).
- Cold start (no corpus Feature in a completed window) → velocity stays undefined; pack uses `initial_velocity` 10; `predicted_duration` unused. A Feature accepted in the open window does not enter the lookback.
- Bugs/chores/releases not required; when present later they follow the velocity doc.

**QA**

- Run the seven-step script in `velocity-and-planning.md` (steps 5–6 that need Release / illegal drag may wait for slices 18 / 7 respectively; step 7 is this slice).
- Fail if a Recalculate button is required. Fail if velocity is typed in as a target. Fail if the board requires an iteration table to render bands.

### Isolation tests (ride 0, 1, 23)

From `multitenancy.md`: same title in org A and B, search/list only A; A fetches B’s id → 404; Bearer whose user cannot open that project → 404; attachment from A, signed-out → no file; viewer `POST /stories` → `forbidden` (session or Bearer — same code).

### Later phases (notes, not Phase 0 work)

| Slices | Test notes |
| --- | --- |
| 9 tasks | Toggle persists; complete-all does not Finish; accept **warns**; viewer read-only |
| 10 blockers | Free-text + optional same-project story; auto-resolve on accept/delete; warn on accept |
| 11 labels | Existing tables; `[a-z0-9-]+`; column max 100; do not exceed 100; filter does not change rank |
| 12 owners / mentions | Max 5; requester is a Member/Owner; in-app + email; viewer mentionable, cannot follow |
| 13 Markdown / comments | No raw HTML; `javascript:` rejected; tombstone delete |
| 14 attachments | Type/size/count; auth GET; paste; missing-image after delete |
| 15 activity / undo | No reorder spam; undo only latest state change; undo is an activity; `activities.user_id` is the user |
| 16 keyboard | Every listed action has a chord from `ui.md` |
| 17 SSE | Two sessions < 2s; draft preserved; stale + reconnect |
| 18 types | Same machine package; Bug start without points; Chore finish=accept; Release colour vs **starts_on** of the calculated window |
| 19 epics | Purple label; independent order; progress Feature points only |
| 20 dates | Calculated window `ends_on`; Icebox “Not scheduled”; no plan-overriding picker |
| 21 charts | Bars from the same live velocity calculation; empty / cold-start: velocity line “initial 10” |
| 22 search | Operators as specified; Done excluded unless `includedone:` |
| 23 API tokens | Mint / list / revoke only on your own user: `GET` / `POST` / `DELETE` `/api/v1/users/:id/tokens`. No mint-for-another-user. Bearer start/finish/deliver/accept as Member. Viewer may mint a token on their own user; it can only read. Member cannot mint Owner. Token whose user cannot open project B → 404. |
| 23b Webhooks | Outbound project hooks. Not part of the mint slice. |
| 24 My Work / saved search | Owner/requester rules as specified |
| 25 workspaces | One organisation; not a permission boundary |
| 26 CSV | Owner only (assumption); create-only, all-or-nothing |
| 27 cycle time | First `started_at` → `accepted_at`; reject does not reset |
| 28 scales / bugs-chores toggle | No history conversion; toggle reversible |
| 29 split panes | Both live-update |
| 30 manual Current | Later bands still auto-plan |

---

## Quality bar

From repo `AGENTS.md` and `api/internal/migrations/AGENTS.md`. Non-negotiable.

**Testing**

- TDD: failing test, then code. `nix-shell --run 'make test'` and `make lint` after every code change.
- 90%+ line coverage on API and frontend. Do not lower the gate.
- No placeholder, fake, or unused production branches. Tests exercise the same machine, planner, and auth as runtime.
- **No fallbacks.** Dummy password hashes, silent estimate-on-start, “if planner fails return last board”, or catch-all `state = body.State` are firable. Fail closed with a coded error.
- Do not ignore test, script, or lint errors.
- Documentation-only changes: review the diff; do not run the code gates unless code changed.

**Types and API**

- Go: wrap with `fmt.Errorf("%w")`. Structured Zap fields: `component`, `operation`, `organisation_id`, `project_id`, `user_id`. Never log passwords, raw session/magic/token/webhook secrets, or Authorization headers.
- Frontend: TypeScript strict; oxlint; `tsc --noEmit`. Types for board/story match the API.
- One error envelope. The codes are the codes (shared; session or Bearer).

**Migrations**

- Schema only. Data backfills run in Go on boot.
- No DB enums, triggers, or functions. Plural table names. UUID `id` PKs.
- Australian / New Zealand / UK names (`organisations`).
- Name migration files with `date +%Y%m%d%H%M%S`.
- Additive columns and new tables. Do not add an iteration table or a story window FK.
- `point_scale` / `role` / `state` / `story_type` stay strings.

**Logging and config**

- Env from `.env.example`. No secrets in the repo.
- Clock, email, storage, bus are injected. Production has one real clock: `time.Now`.

**Feature flags and compatibility**

- No flags to hide Icebox, velocity, or accept-by-any-Member.
- Do not treat `rejected` as a terminal peer of accepted.
- Cookie session and Bearer share machines and rank rules.

**Authz checklist (every mutating handler)**

1. Authenticate (session or `flr_` bearer).
2. Resolve organisation from the project, not from a client-supplied org header alone.
3. Effective role (Owner / Member / Viewer) — **not** a scope list.
4. Machine / rank / estimate rules.
5. Transaction: mutate, activity, pack, revision++.
6. Commit then bus.

**Spelling:** organisation, behaviour, colour (docs). Code identifiers for new tables: `organisations`, `organisation_id`. Always **user**, never a second identity word, in product language.

---

## Implementation assumptions (forced by schema or stack)

These are not new product rules. Product remains unspecified where marked.

| Item | Status | Assumption if we must ship |
| --- | --- | --- |
| Username at signup | Locked 20 Aug 2026 | Infer from email local-part; uniquify; no username field in slice 0; editable later |
| Organisation public slug | Unspecified (fork 2) | UUID routes; no public org slug in Phase 0 |
| `display_name` | Unspecified | Copy username at create |
| Project timezone | Fork 3; velocity doc already says store + default Melbourne | `projects.timezone` default `Australia/Melbourne` |
| Icebox vs one rank | Assumed two lists (fork 4) | `rank_list` + unique `(project_id, rank_list, rank)` |
| Who creates projects after slice 0 | Assumed organisation owners | Enforce in Go |
| Reject reason vs comments table | Comments are slice 13 | Phase 0: activity only |
| HTTP prefix | Locked | **`/api/v1`** mount; same bodies and errors for session and Bearer; no second `/v1` tree |
| Create-story `story_type` | Shared contract | One API; omitted `story_type` may default to Feature (UI also posts Feature) |
| Create-story `requester_id` | Product: requester is a Member/Owner | `requester_id` is the authenticated user (session or Bearer) and must be a Member / Owner |
| Create-story `panel` | Product: default add is Icebox | Omitted `panel` → `icebox` for every user |
| Token mint path + role-cap | This file | `GET` / `POST` / `DELETE` `/api/v1/users/:id/tokens` only on your own user. No mint-for-another-user. Role = that user's project memberships, or an optional cap on `api_tokens.role` at or below the minter. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. Token is a user credential, not an identity type. No org-collection token routes. No grants table (cap is a column). |
| Planning persistence | This file | None. Settings + live duration velocity + `pack(ordered_stories, velocity, predicted_duration, iteration_length_days, now, timezone) → bands`. Velocity stays undefined until at least one corpus Feature exists in a completed window (`pack` uses `initial_velocity` 10; `predicted_duration` unused). A story accepted in the open window is not in the lookback. |
| Undo | Product-spec: Owner, Member on latest state change | Same for session or Bearer. Not “own last change only”. |
| Human session mechanism | Unspecified | Server-side cookie sessions |
| Live transport | Unspecified | SSE + in-process bus, slice 17 |
| Attachment / email vendors | Unspecified | Storage and email interfaces; local dir + SMTP/outbox |
| Rate-limit numbers | “TL picks” | 60 mutations / minute / Bearer token, burst 10; 600 reads / minute / token; 120 / minute / human session |
| Webhook HMAC scheme | “TL specifies” | `t=<unix>,v1=<hmac-sha256>` |
| Magic/verify link TTL | Unspecified | 30 minutes |
| Password hash | Unspecified | bcrypt, cost 12, nullable column |
| Iteration start weekday storage | Implied Monday | ISO `1` = Monday on `projects.iteration_start_weekday` |
| `planned` state | Fork 9, later | Do not add in Phase 0 |

If Dan answers a fork, update this file in the same PR as the code change. Do not leave the doc lying.

---

## Required reading (implementers)

- This file
- `domain-model.md` (where code goes; bounded contexts)
- `product-spec.md` (the slice you are in, plus machines)
- `velocity-and-planning.md` if the slice touches the plan (3, 4, 7, 8, 18, 20, 21, 30)
- `multitenancy.md` if the slice touches authz (0, 1, 23, 25)
- `open-questions.md` assumptions
- this file’s intended schema + `domain-model.md`
- repo `AGENTS.md`, `api/internal/migrations/AGENTS.md`
- `docs/reference/frontend-design-guide.md` for UI work (do not invent a look)
- `ui.md` once the UI Designer has written it
