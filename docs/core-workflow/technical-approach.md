# Technical approach

Eventual repo path: `docs/core-workflow/technical-approach.md`

Local draft: `/workspace/flower-spec/technical-approach.md`

Change name: `flower`

This is the Technical Lead approach for implementers. It does not replace the product companions. It does not invent product rules.

**Source of truth (product):**

| Topic | Wins |
| --- | --- |
| Packing, velocity, rollover, releases, dates | `velocity-and-planning.md` |
| Agent verbs, error codes, scopes, webhooks | `agent-api.md` |
| Organisations, roles, isolation | `multitenancy.md` |
| Slice order and acceptance criteria | `product-spec.md` |
| Copy-exactly vs modernise | `tracker-brief.md` |
| True forks + baked assumptions | `open-questions.md` |
| Existing eight tables, rank type, string states | `000001_create_core_schema.up.sql` |
| Quality bar, no fallbacks, TDD, 90% | repo `AGENTS.md` + `api/internal/migrations/AGENTS.md` |

If a slice, mock, or this file disagrees with the velocity doc on packing, the velocity doc wins. If this file disagrees with the agent API on verbs or errors, the agent API wins. If something is unspecified, it is marked **unspecified** below; an implementation assumption is listed only where a column or locked stack forces a choice.

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
| Prophet | Owns 8080 / 4173 — do not bind those |

Existing process commands stay: `nix-shell --run 'make dev'`, `make test`, `make lint`, `make migrate`. `.env.example` remains the source of truth for required variables.

### Where slices live

The API already documents a domain-driven layout (`domain/`, `platform/`, `app/`, `handlers/`). Keep it. Map work by vertical slice, not by “build all tables then all handlers.”

```
flower/
├── api/
│   ├── cmd/                  # process: HTTP, migrate, boot backfills
│   └── internal/
│       ├── domain/           # machines, planner, rank, tenancy rules (no HTTP)
│       ├── platform/         # postgres, clock, email, bus, object store
│       ├── app/              # use-cases: one package per outcome
│       ├── handlers/         # Gin, /api/v1/*
│       └── migrations/       # schema-only; 000001 is ground
├── frontend/src/
│   ├── features/auth/
│   ├── features/board/       # Icebox / Backlog / Current / Done
│   ├── features/story/
│   └── lib/api.ts            # session + agent client
└── docs/                     # this spec set lands per LANDING.md
```

Phase 0 use-case packages (add as the slice lands, not all on day one):

| Slice | `api/internal/app/` | `frontend/src/features/` |
| --- | --- | --- |
| 0 | `auth`, `organisation`, `project`, `board` (empty) | `auth`, `board` (empty columns) |
| 1 | `invite` | project members / invite |
| 2–3 | `story` (create, schedule, icebox) | Icebox + pull |
| 4–6 | `story` (estimate, transition) | estimate + verbs |
| 7 | `rank` | reorder |
| 8 | `planning` | iteration headers, V, Done after rollover |

Shared domain (written with the first slice that needs it, tested in isolation):

- `domain/story/machine` — type-specific verbs. Feature machine is used in Phase 0. Bug / Chore / Release tables are coded in the **same** package when the Feature machine lands (product: machines are specified now so later slices do not invent a second workflow). UI and create-story only expose Feature until slice 18.
- `domain/story/rank` — fractional `VARCHAR(64)` generate / compare / illegal-rank check.
- `domain/planning` — velocity formula + pack algorithm. Velocity doc is normative.
- `domain/tenancy` — organisation scope, effective role, 404-vs-403.

Frontend does **not** recompute velocity or panel membership. It renders what the board payload says. UI Designer owns chords, empty copy, and layout in `ui.md`. Implementers must not invent a palette or a private keymap.

### HTTP mount

Existing repo README exposes `/health`, `/ready`, `/api/version`, and **`/api/v1/*`** on 8180.

`agent-api.md` writes paths as `/v1/...`. That document wins on **verbs and error codes**, not on the mount prefix.

**Decision:** one mount, `/api/v1`. Agents call `http://localhost:8180/api/v1/stories/:id/transitions`. Humans use the same resource paths with a session cookie. Do not run a second `/v1` tree.

SPA origin is `:4273`, API is `:8180` (cross-origin, same-site on localhost). CORS: explicit `FRONTEND_ORIGIN` (never `*`), credentials allowed. Production: put both behind one origin and drop credentialed CORS if possible.

### What must not be disturbed

- The eight core tables’ meaning: `users`, `projects`, `project_memberships`, `iterations`, `stories`, `labels`, `story_labels`, `activities`.
- `stories.rank VARCHAR(64)` fractional / lexicographic. Do not switch to integer priority.
- `stories.title VARCHAR(500)` — product max is 500.
- `story_type` and `state` as strings. No DB enums, checks-as-enums, triggers, or functions.
- `activities.actor_id → users`. Agents get a `users` row (see systems design). Do not make `actor_id` polymorphic in Phase 0.
- `projects.point_scale`, `projects.iteration_length_weeks` — already present; add columns, do not rename.
- Unique `(project_id, user_id)` on memberships; unique `(project_id, number)` on iterations.
- `docs/reference/*`, existing `000001`, Prophet ports, locked look.
- UK / AU / NZ identifiers: `organisations`, not `organizations`.

New slices **add** tables and columns. They do not redesign the eight.

---

## Systems design

### Ground schema (000001) — keep

Exact columns that later slices must treat as given:

- `users`: `id`, `username` UNIQUE NOT NULL, `email` UNIQUE NOT NULL, `password_hash` (NOT NULL today — see slice 0 change), `display_name` NOT NULL, timestamps as `TIMESTAMP` (not `TIMESTAMPTZ`).
- `projects`: `id`, `name`, `slug` (UNIQUE globally today — see slice 0 change), `description`, `point_scale`, `iteration_length_weeks`.
- `project_memberships`: `role VARCHAR(50)` — product values `owner` \| `member` \| `viewer` only, enforced in Go.
- `iterations`: `number`, `starts_on DATE`, `ends_on DATE`.
- `stories`: `iteration_id` nullable SET NULL, `requester_id` RESTRICT, `estimate INTEGER` nullable, `rank VARCHAR(64)`, `accepted_at` nullable, unique `(project_id, rank)` today.
- `labels.name VARCHAR(100)`, unique `(project_id, name)`.
- `activities`: `kind`, `summary`, `actor_id` RESTRICT, `story_id` nullable SET NULL.

`NULL` estimate = unestimated. `0` is an estimate. Do not store `-1`. Search `estimate:-1` (slice 22) means `estimate IS NULL`.

Existing timestamps are `TIMESTAMP`. **Do not** rewrite 000001 columns to `TIMESTAMPTZ` as a drive-by. Store UTC in `TIMESTAMP`. Convert with the project timezone for display and rollover. New tables follow the same convention until a dedicated later migration (not Phase 0) says otherwise.

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

**Slug / URL (open question 2 — unspecified product):** API and SPA routes use UUIDs: `/api/v1/organisations/:organisation_id/projects/:project_id`. No public organisation slug in Phase 0. `projects.slug` stays as the existing column, unique **per organisation**, generated from the project name plus a short suffix on collision. Do not invent a new slug scheme.

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
- `users.password_hash` becomes **nullable** (schema-only). Magic-link-only accounts have `NULL`. Password login on `NULL` is `unauthorized` with a message to use the magic link — that is a real branch, not a fallback hash.
- `users.actor_kind VARCHAR(20) NOT NULL` default `human` (Go: `human` \| `agent`). Session and magic-link login reject `actor_kind != human`.
- `users.last_project_id` optional; session `last_project_id` is enough for “land on my last project” if updated on board load.

**Username (locked 20 Aug 2026):** infer from the email local-part. No username field in slice 0. `users.username` is `NOT NULL UNIQUE`. Slice 0 must persist one: derive from the email local-part (`[a-z0-9_]+`, trim to 100); if taken, append `-` + 4 random unambiguous chars. Allow edit later when a profile slice exists. Do not write an empty string (that is a fallback).

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

1. If the account is an organisation owner of `project.organisation_id` → treat as project `owner` even with no `project_memberships` row.
2. Else the `project_memberships.role` for `(project_id, user_id)`.
3. Else: if the project’s organisation is not one they belong to → **404** on id fetch. If they belong to the org but not the project → **404** (enumeration: project lists only return projects they can open).
4. Same-tenant, known project, insufficient role on a mutation → `forbidden` (403). Viewers can GET the board; `POST /stories` is `forbidden`, not 404.

Cross-tenant id (story, project, attachment, token) → **404** / `not_found`. Never 403 for “exists in another organisation.”

Organisation owner can do anything a project owner can. Any **Member or Owner** can accept / reject Features. Viewers cannot. Agents cannot accept / reject Features (`human_judgment_required`). Do not add an accept ACL.

**Email:** product does not name a vendor. `platform/email` is an interface. Production: SMTP from env. Dev/test: a real **outbox** table written in the same transaction as the token row; a worker or test helper sends/reads it. Tests assert outbox rows. Do not “succeed” without a stored token. Do not log raw links at info level in production (debug in test env only).

### State machines

Business rules live in Go. No DB enum, trigger, or function.

One table-driven package: `domain/story/machine`. Input: type, from-state, action, actor kind, estimate, scopes/role. Output: to-state or a coded error (`invalid_transition`, `unestimated`, `human_judgment_required`, `forbidden`, `validation_failed`).

Normative tables: `product-spec.md` (humans) and `agent-api.md` (agents). Agent API **wins** on verb names and error codes.

**Feature / Bug**

| action | from | to | Agent |
| --- | --- | --- | --- |
| schedule | unscheduled | unstarted | `stories:write` |
| icebox | unstarted only | unscheduled | `stories:write` |
| start | unstarted | started | if `stories:transition`; Feature also needs estimate |
| start | unscheduled (Icebox) | started | if `stories:transition`; this verb is schedule+start. Story lands started in Current (may overflow velocity). A typical CI token with stories:transition and not stories:write may still do this. Feature must already be estimated (`0` allowed). |
| finish | started | finished | `stories:transition` |
| deliver | finished | delivered | `stories:transition` |
| accept | delivered | accepted | **human only** in MVP |
| reject | delivered | rejected | **human only**; non-empty reason |
| restart | rejected | started | `stories:transition` |
| undo | last state-changing activity | previous | human UI (slice 15); agent may undo **own** last state change, not a human Accept |

**Chore:** unscheduled → unstarted → started → accepted. `finish` is the verb; result is `accepted`. No delivered / reject. Agent may finish.

**Release:** unscheduled in Icebox; created in Backlog or scheduled → auto-`started`. `finish` → `accepted`. No estimate.

`PATCH` must not write `state`. There is no `unstart` verb. Icebox is `icebox`. Undo is `undo` (slice 15 / agent contract).

**Start Feature** without estimate → `unestimated`, no mutation. `0` is allowed. Start assigns the clicker (human or agent) as a story owner if owners < 5; if already 5 and clicker is not among them, Start still happens and they are **not** added (product). Auto-follow the clicker (requester and owners cannot unfollow — slice 4 AC requires follow to exist).

**Reject reason (Phase 0):** product calls it a comment; the comments table is slice 13. Phase 0 stores the reason on `activities` (`kind = story.rejected`, `summary` = reason) and returns it on the story payload as `reject_reason` from the latest reject activity. Empty reason → `validation_failed`, no state change. Do not invent a comments table in Phase 0. Slice 13 may also write a comment; **unspecified** whether both exist — do not do both until specified.

**Optimistic concurrency:** add `stories.revision INTEGER NOT NULL` default 1, incremented on every story mutation (fields, rank, state, estimate). Agents send `revision` or `If-Match`. Stale → `conflict`. Humans: last-write-wins on the board is unacceptable for rank/state; the SPA sends `revision` as well.

**`started_at`:** add `stories.started_at TIMESTAMP NULL`. Set on the **first** successful `start` only. Reject / restart do not clear it (cycle-time assumption, slice 27). Do not reset the clock.

**Undo (slice 15, design now):** undo applies only to the latest **state-changing** activity on that story. Reorders do not write activity. Undo itself writes an activity. Implement in the machine as `action=undo` with the previous state from that activity row (store `from_state` / `to_state` in `summary` as structured text or add nullable `activities.payload JSONB` — JSONB is schema, not a function; acceptable). Prefer additive `from_state` / `to_state` VARCHAR columns on `activities` when slice 15 lands rather than parsing `summary`.

Illegal verbs do not change state.

### Velocity recompute

Normative algorithm: `velocity-and-planning.md`. Do not invent a third planner. Do not show a Recalculate button.

**Add to `iterations` (slice 8, not a redesign):**

- `accepted_feature_points INTEGER NULL` — frozen at rollover. Required because deleted Features that were in a completed total stay in that total.
- `completed_at TIMESTAMP NULL` — set at rollover; Current is the iteration with `completed_at IS NULL` and `starts_on <= today <= ends_on` in project TZ.

**Velocity (MVP), in Go, integer:**

```
N = count(completed iterations)
K = project.velocity_strategy   # 1–4, default 3
if N == 0:
    V = project.initial_velocity    # default 10
else:
    window = last min(K, N) completed, by number desc
    if every window row has accepted_feature_points == 0:
        V = project.initial_velocity    # Tracker all-zero revert
    else:
        V = floor(mean(accepted_feature_points of window))
```

Accepted Feature points only. Estimate at **accept** time is what is recorded into the running total; after rollover the frozen integer is the source. Edits after accept do not rewrite a completed iteration. Bugs / chores / releases add 0 until the Phase 3 toggle (do not add the toggle column in Phase 0). Team strength % is not MVP — do not hide a 100% factor.

**Panel membership** is a function of state + `iteration_id` + current window, not a panel enum:

| Panel | Rule |
| --- | --- |
| Icebox | `state = unscheduled`. Own order (`rank_list = icebox`). `iteration_id` NULL. |
| Current | In-flight (`started`/`finished`/`delivered`/`rejected`) **or** `accepted` with `accepted_at` inside the current window **or** `unstarted` (and packed Releases) with `iteration_id = current`. |
| Backlog | Ranked-list stories whose `iteration_id` is a future iteration. |
| Done | `accepted` and `iteration_id` (or `accepted_at` window) is a **completed** iteration. Newest completed iteration first. |

Accepted-this-iteration stays in Current until rollover. Done is empty until the first rollover.

**Planner** (`domain/planning.Pack`) is a pure function: `(now, project, V, ranked stories, iterations) → iteration assignments`. It does **not** change `rank`. It writes `stories.iteration_id` and creates missing future `iterations` rows (number, starts_on, ends_on). First iteration: Monday on or before project-created date in project TZ (or the configured start weekday), length = `iteration_length_weeks`, `ends_on` = last calendar day of the box.

Call Pack **inside the same transaction** as any mutation that the velocity doc lists: create, delete, icebox, schedule, reorder, estimate, type change, start and other verbs, accept, reject, restart, undo, rollover, V / length / TZ / strategy / initial-velocity / (later) toggles. After commit, publish the bus event.

Start on a Backlog or Icebox story jumps it to the ranked list + Current and **may** overflow V. Drag Backlog → Current is **not** Start (slice 7). Auto-plan membership is not drop-to-start.

Leave-short, never-split, oversized-Feature exception, zero-cost packing: copy the velocity doc. Do not rephrase a different fit rule.

**Rollover** at midnight at `ends_on + 1 day 00:00` project TZ:

1. Freeze `accepted_feature_points` for that iteration (sum of Feature estimates accepted in the window).
2. Set `completed_at`.
3. Recompute V.
4. Those accepted stories are now Done (they already have `accepted_at` / `iteration_id`; the board query flips them because the iteration is completed).
5. In-flight stay; leftover unstarted are re-packed. Unaccepted work is not failed and not iceboxed.

Implementation: (a) a process ticker every 30s that selects projects whose current iteration is due, and (b) the same function at the start of every board read/write. Guard with `WHERE completed_at IS NULL` so two paths cannot double-freeze. Inject `platform.Clock` for tests.

**Test clock:** QA must advance past midnight. Production builds do not expose a clock endpoint. Test / `APP_ENV=test` (test compose, Postgres 5437) may expose `POST /api/v1/test/clock` (`{"now": "ISO-8601"}`) compiled behind a build tag or env that **production refuses to start with**. That is a test harness, not a fallback.

Changing length / timezone / start weekday replans immediately. Completed iterations keep historical `starts_on` / `ends_on` and frozen points.

`GET /api/v1/projects/:id/board` returns panels, current V, `initial` vs calculated, points/V, `ends_on`, over-velocity badge if Current Feature-points > V. Frontend does not pack.

### Ranking

Ground: `stories.rank VARCHAR(64)`, unique per project today.

Icebox is its own ordered list (assumption; open question 5). Unique `(project_id, rank)` cannot hold two independent sequences without collisions (both lists would generate `"n"`).

**Add (slice 2), do not change the type of `rank`:**

- `stories.rank_list VARCHAR(16) NOT NULL` — Go values `icebox` \| `ranked`.
- Drop `UNIQUE (project_id, rank)`.
- `UNIQUE (project_id, rank_list, rank)`.

Icebox order is independent. New Icebox story goes to the **top** of Icebox (rank before current first, or a start-rank if empty). Schedule (slice 3) moves `unscheduled` → `unstarted`, `rank_list` `icebox` → `ranked`, rank = **bottom** of the ranked list (after current last). Icebox of `unstarted` only: reverse. Cannot icebox started / finished / delivered / rejected / accepted.

**Algorithm:** lexicographic midpoint on charset `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`. Generate between `before` and `after`; empty list → a mid string (e.g. `U`). If the midpoint would exceed 64 characters, **rebalance** that `(project_id, rank_list)` in the same transaction (rewrite ranks evenly). Rebalance is an implementation event, not an activity.

**Illegal rank (server, not only UI):** unstarted cannot sit above started / finished / delivered / rejected in the **ranked** list. Reject the write, `illegal_rank`, row snaps back. Accepted-this-iteration do not return to Backlog by drag. Viewers cannot reorder.

API (humans and agents with `stories:write`):

```
POST /api/v1/projects/:project_id/stories/reorder
{ "story_id": "...", "before_id": "..." | "after_id": "...", "revision": 3 }
```

Agents must send explicit neighbours. Do not default to “top.” Unique-violation or stale revision → `conflict`. Retry once server-side on unique violation after re-read; then fail.

Reorders do **not** write `activities` (slice 15).

### Real-time (slice 17)

Product: two sessions, same project, mutations appear within **2 seconds** without refresh; focus and comment drafts survive unless that story/comment was deleted; no presence / live cursors; dropped socket shows stale + reconnect; refresh heals; viewers get the same reads.

**Decision:** SSE, not WebSockets (presence is out of scope). Not polling as the primary path (easy to miss the 2s bar under load).

```
GET /api/v1/projects/:project_id/events
Accept: text/event-stream
```

Session or agent `stories:read`. After each committed mutation, `platform/bus` publishes `{project_id, event, story_id, revision}`. The handler fans out to subscribers of that project.

MVP bus: **in-process**. Document single-API-replica for live updates. Interface it (`Subscribe(projectID)`) so a Redis/NATS adapter can land later without changing handlers. Do not add Redis in Phase 0.

Event names align with webhook events where they overlap (`story.created`, `story.updated`, `story.reordered`, `story.started`, …). Payload is enough to patch or to trigger a targeted GET. Do not replace the board with a full refetch on every keystroke if that blows away a comment draft — patch by `story_id`, preserve local draft state.

Client: `EventSource` (or fetch-stream) with cookie. On error: banner “stale”, reconnect with backoff. Manual refresh heals.

Slice 17 is Phase 1. Phase 0 may ship without SSE; the board remains correct on reload. When slice 17 starts, publish from the same post-commit hook used for webhooks so we do not grow a second event list.

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

**Storage (unspecified vendor):** `platform/storage` interface (`Put`, `Get`, `Delete`). Dev: local directory under a data path, keys `{organisation_id}/{project_id}/{story_id}/{attachment_id}`. Production: S3-compatible (env: bucket, region, endpoint). Same interface, no “if S3 fails write to /tmp” fallback.

Validate type by sniffed bytes, not only extension. Over 10 MB or 21st image → `validation_failed`.

### Agent tokens and webhooks (slice 23)

Agents are first-class. They do not impersonate the minting human. Activity **name** is the agent’s name.

**Users row for the agent** (ground: `activities.actor_id → users`):

- Insert `users` with `actor_kind = agent`, `password_hash NULL`, unique `username` = agent name slug, unique `email` = `agent-<uuid>@invalid` (not a login identity; session/magic-link reject this kind).
- `agents` table holds the grant metadata; `users.id` is the actor id.

```
agents
  id UUID PK
  user_id → users RESTRICT          -- actor
  organisation_id → organisations RESTRICT
  name VARCHAR(255) NOT NULL
  created_by_user_id → users RESTRICT
  created_at TIMESTAMP NOT NULL

agent_tokens
  id UUID PK
  agent_id → agents CASCADE
  token_prefix VARCHAR(16) NOT NULL -- listed after create
  token_hash VARCHAR(64) NOT NULL UNIQUE
  scopes TEXT NOT NULL              -- comma-separated; parse in Go
  revoked_at TIMESTAMP NULL
  created_at TIMESTAMP NOT NULL

agent_token_projects
  agent_token_id, project_id
  UNIQUE (agent_token_id, project_id)
```

Secret format: `flr_<random>`. Shown **once**. Store hash only. `Authorization: Bearer flr_...`. Revoke is immediate. Viewer cannot create / list secrets / revoke. Owner or Member, for projects they can write. Bound to **one** organisation; cannot widen. Typical CI: `stories:read` + `stories:transition` + `comments:write`. `stories:accept` is **not granted** in MVP.

`GET /api/v1/me` for an agent: id, name, organisation, projects, scopes — never the minting human’s email as the actor.

Transitions: `POST /api/v1/stories/:id/transitions` with `{ "action": ... }`. Feature/Bug `accept` / `reject` → `human_judgment_required` (and `forbidden` if you also want the HTTP class; agent-api lists both — use **`human_judgment_required`** as the code, HTTP 403). Start unestimated Feature → `unestimated`. Illegal → `invalid_transition` with `from` and `action`. Cross-tenant → `not_found`. Revoked → `unauthorized`.

Create story (agent): `story_type` required (do not default Feature), `requester_id` required and a human Member/Owner, `panel` default `icebox`. Humans in the UI may default Feature — that is UI, not the agent contract.

Bulk: max 50, all-or-nothing, `Idempotency-Key`. Same key + body → original result. Same key + different body → `conflict`.

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

**Webhooks:**

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
  event_id UUID NOT NULL            -- idempotency for the receiver
  event VARCHAR(100) NOT NULL
  payload JSONB NOT NULL
  attempt INTEGER NOT NULL
  next_attempt_at TIMESTAMP NULL
  delivered_at TIMESTAMP NULL
  last_status INTEGER NULL
```

Signature (Technical Lead scheme): header `X-Flower-Signature: t=<unix>,v1=<hex>` where `v1` is HMAC-SHA256 of `{t}.{raw_body}` with the raw webhook secret. Receivers must reject if `|now - t| > 300s`. At-least-once. Retry non-2xx. Timeout 10s. Delivery is not a command; ignore reply body. Worker in the API process.

Events: those listed in `agent-api.md`. Envelope as specified (`event_id`, `organisation_id`, `project_id`, `actor.kind` `human` \| `agent`).

**Rate limits (product: a well-behaved 1 rps single-story transition must not 429):**

| Actor | Limit |
| --- | --- |
| Agent transitions | 60 / minute / token, burst 10 |
| Agent reads | 600 / minute / token |
| Human session | 120 / minute |

In-memory per process in MVP (same single-replica assumption as SSE). `rate_limited` + 429 + `Retry-After`.

Error envelope for **all** 4xx/409 (humans and agents), as agent-api:

```
{ "error": { "code": "invalid_transition", "message": "...", "from": "unstarted", "action": "finish" } }
```

No 200 with a partial surprise. No coerce (start that silently estimates 1).

### Story owners and follow (needed in slice 4, not only 12)

000001 has `requester_id` only. Do not overload it.

```
story_owners
  story_id → stories CASCADE
  user_id → users CASCADE          -- human or agent user
  UNIQUE (story_id, user_id)
  -- max 5 enforced in Go; sixth → owners_full

story_followers
  story_id, user_id
  locked BOOLEAN NOT NULL          -- true for requester + owners
  UNIQUE (story_id, user_id)
```

Viewers do not follow. Requester is a human Member/Owner, never an agent.

### Humans: Phase 0 HTTP (minimum)

Auth and tenancy:

- `POST /api/v1/auth/signup` (email, password, …)
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
- `POST /api/v1/invites/:token/accept` (token in path or body; hash lookup)

Stories:

- `POST /api/v1/projects/:id/stories` (human may omit `story_type` → Feature in the **UI**; API for humans may default Feature; **agents must not**)
- `GET  /api/v1/stories/:id`
- `PATCH /api/v1/stories/:id` (fields, not state)
- `POST /api/v1/stories/:id/transitions`
- `POST /api/v1/projects/:id/stories/reorder`

Exact JSON field names for the human board payload are internals; they must include id, title, type, state, estimate, rank, revision, iteration ends_on, reject_reason if rejected, owners. Do not invent extra product fields (no due date, no custom fields).

---

## Refactor opportunities

Only what a slice cannot ship without. No drive-by.

| Must take | Why | Slice |
| --- | --- | --- |
| `projects.slug` unique per organisation, not globally | Tenancy; existing index is global | 0 |
| `users.password_hash` nullable | Magic-link-only accounts; `NOT NULL` plus a dummy hash would be a fallback | 0 |
| `UNIQUE (project_id, rank)` → `UNIQUE (project_id, rank_list, rank)` | Two lists, one `VARCHAR(64)` column | 2 |
| Add `organisation_id` / timezone / velocity settings on `projects` | Missing vs product; do not recreate `projects` | 0 / 8 |
| Add freeze columns on `iterations` | History-stable velocity | 8 |
| Add `revision`, `rank_list`, `started_at` on `stories` | Concurrency, two lists, cycle-time later | 2 / 4 |
| Add `email_verified_at`, `actor_kind` on `users` | Slice 0 AC + agents | 0 |

**Do not take in these slices:**

- Rewrite the eight tables or rename existing columns.
- Switch `rank` to integer or `bigint` priority.
- `TIMESTAMP` → `TIMESTAMPTZ` across 000001.
- DB enums / triggers / functions for state, role, or type.
- GraphQL.
- A second workflow engine or `planned` state (open question 9; only if the manual-planning slice needs it).
- Extracting a microservice, adding Redis, or adding a queue product for Phase 0.
- “Fixing” Icebox into a sequential stage to match the old root README.
- Changing ports or the locked look.

---

## Technical risks

| Risk | Why it is real | Mitigation |
| --- | --- | --- |
| Planner drift from the velocity doc | Easy to “improve” leave-short or oversized Features | Pure `Pack` + table tests for **all three** worked examples in the velocity doc; QA script in that file; velocity doc wins reviews |
| Double rollover / missed rollover | Ticker + request path; TZ midnight | `completed_at IS NULL` guard; injected `Clock`; test env clock endpoint; assert freeze is idempotent |
| Rank unique collisions / 64-char overflow | Two lists; midpoint can grow | `rank_list` in the unique key; rebalance in-transaction; tests at the 64-char edge |
| Cross-tenant leaks (id guess, search, files, tokens) | Multitenant from day one; 000001 has no org | Every query joins organisation; miss → 404; attachment GET checks org+project+auth; isolation tests ride slices 0, 1, 23 |
| `activities.actor_id` cannot point at a non-user | Agents must be attributed as themselves | Agent `users` row + `actor_kind`; login rejects agents; never attribute the minting Member |
| Session cookie on :4273 → :8180 | Cross-origin | Explicit CORS origin + credentials; prod same-origin proxy |
| SSE lost across API replicas | In-process bus | Single replica in MVP; `Bus` interface; do not pretend it is multi-node |
| Webhook retry storms | At-least-once | Cap attempts (implementation: 8, exponential backoff); `event_id` for receiver idempotency; do not POST back as a user |
| Magic link / invite token leak | Email and logs | Store hashes only; single-use; 14-day invite; do not log raw tokens |
| Clock skew vs project TZ | Rollover definition is midnight **project** TZ | Store TZ; default `Australia/Melbourne`; all comparisons via `Clock` + that TZ |
| Overview.md / root README still describe Icebox as a pipeline and `rejected` as a terminal peer | Agents will implement the old overview | LANDING correction in the implementation PR; this approach and tracker-brief supersede; QA fails those old readings |
| Dummy `password_hash` for magic-link users | Hidden branch, forbidden by AGENTS.md | Nullable column; password login refuses NULL |
| Existing projects without `organisation_id` if 000001 already ran | NOT NULL add | Schema: add nullable, boot **backfill** in Go (create an organisation only for orphan rows if any exist — this is data, not a SQL migration), then schema NOT NULL. Empty greenfield: backfill is a no-op |
| Unique rank violation under concurrent drag | Two Members reorder | `revision` + one retry; then `conflict`; UI snaps back |
| File store credentials / public buckets | Tenant files | Auth-only GET; key prefixed by organisation_id; no public ACL |
| Rate limiter per process | Two replicas double the limit | Accept in MVP; document; do not 429 a 1 rps agent |

---

## Test plan

TDD: failing test first. Tests use the same production logic (no fake planner, no fake machine). Inject `Clock`, email outbox, storage, bus.

QA tests each slice from its acceptance criteria **alone**, plus the companion rule docs named below. Fail without a meeting.

### Slice 0 — signup, organisation, empty board

**API / domain**

- Unverified password signup: org create is rejected; no `organisations` row.
- Verified password signup: organisation + first project; creator is organisation owner and project owner; `point_scale = linear`, `iteration_length_weeks = 1`, `initial_velocity = 10`, `timezone = Australia/Melbourne` (until fork).
- Magic link to a new email: account created, `email_verified_at` set, same org+project flow.
- Login password and login magic link both land on last project.
- Board GET: four panels, all empty, V displays 10, no fake stories, no fake dates.
- Reload: same organisation / project ids.
- Cross-tenant: session A `GET` project/story id from org B → 404.
- `actor_kind = agent` cannot create a session.

**QA**

- Stranger signs up, names organisation `Acme` and project `Trail`, sees Icebox / Backlog / Current / Done empty (Current may show dates + velocity 10).
- Sign out, sign in, same empty board. Fail if another tenant’s name appears.
- Fail if a new palette appears (frontend guide).
- Fail if unverified password user can create an organisation.

### Slice 1 — invite

**API**

- Owner invites `alex@example.com` as `member`: outbox has one email; pending invite listed; token hashed.
- New email accept: signup, lands as Member.
- Existing account accept: project in their list.
- Invite as `viewer`: board GET works; `POST /stories`, invite, settings → `forbidden`.
- Email already a member → visible error, one membership row.
- Revoke → consume fails. Expiry 14 days. Resend invalidates old hash.
- Member or Viewer `POST` invite → `forbidden`.
- Isolation: invite accept does not grant a second project.

**QA**

- Two humans, one empty project. Viewer cannot create a story in UI **or** via API.

### Slice 2 — Feature in Icebox

**API**

- Member/Owner create: `story_type=feature`, `state=unscheduled`, `estimate IS NULL`, `requester_id=me`, `iteration_id IS NULL`, `rank_list=icebox`, rank at top.
- Empty title → `validation_failed`. Title 501 chars → rejected (column 500).
- Viewer create → `forbidden`.
- Second create sits above the first in Icebox.
- Current / Backlog / Done unchanged. Planner ignores Icebox (no `iteration_id`).

**QA**

- Default add lands in Icebox. Fail if it appears in Current or affects V.

### Slice 3 — schedule / icebox

**API**

- Schedule: `unscheduled` → `unstarted`, leaves Icebox, `rank_list=ranked`, rank at **bottom**.
- Icebox: `unstarted` → `unscheduled`, `iteration_id` NULL, no projected date.
- Icebox of started / finished / delivered / rejected / accepted → `invalid_transition` / `illegal_rank` as appropriate; no change.
- Viewer move → `forbidden`.

**QA**

- Pull one Feature to Backlog; Icebox empty state returns. Fail if schedule Starts the story.

### Slice 4 — estimate and start

**API**

- Estimate 0 / 1 / 2 / 3; stays `unstarted`. Other values → `validation_failed`.
- Start unestimated Feature → `unestimated`, still `unstarted`.
- Start estimated: `started`, `started_at` set once, `iteration_id=current` (overflow allowed), clicker in `story_owners`, locked follower.
- Clear estimate only while `unstarted`. Started Feature cannot become NULL.
- Start from Icebox (after estimate): `started` in Current, not `unstarted` in Backlog.
- Viewer estimate/start → `forbidden`.
- Sixth distinct owner on Start: Start succeeds, clicker not added (`owners_full` only when they **try** to add a sixth via assign — Start does not error).

**QA**

- Fail the slice if Start on unestimated succeeds. Fail if 0 is treated as unestimated.

### Slice 5 — finish, deliver, accept

**API**

- started → finish → `finished`, still Current.
- finished → deliver → `delivered`, still Current.
- Any Member or Owner accept → `accepted`, `accepted_at` set, still Current, Done empty.
- Viewer finish/deliver/accept → `forbidden`.
- Agent accept Feature → `human_judgment_required`, no change.
- finish on unstarted, accept on finished → `invalid_transition`, no change.
- Tasks/blockers do not exist yet; do not block accept.

**QA**

- Accepted story still in Current. Fail if it jumps to Done before rollover.

### Slice 6 — reject and restart

**API**

- delivered + reject + reason → `rejected`, still Current; reason on activity.
- Empty reason → no change, `validation_failed`.
- Viewer / agent reject Feature → forbidden / `human_judgment_required`.
- Restart → `started`, still Current; reason remains in activity.
- Finish + deliver again; new accept required.
- Fail if reject jumps to `started` with no Restart, or if `rejected` is treated as terminal like `accepted`.

**QA**

- Rejected work is not in Backlog or Icebox.

### Slice 7 — reorder

**API**

- Member/Owner reorder persists (reload).
- Unstarted above started / finished / delivered / rejected → `illegal_rank`, snap back.
- Accepted-this-iteration cannot be dragged to Backlog.
- Ranked-list drop into Current does **not** Start and does not write `state`.
- Icebox reorder independent; Icebox → Backlog is schedule (slice 3), not Start.
- Viewer reorder → `forbidden`.
- Keyboard reorder path exists (chord from UI Designer; fail if none).

**QA**

- Drag unstarted above started must fail. Fail if drag-to-Current starts the story.

### Slice 8 — auto-plan

Normative tests = velocity doc worked examples 1–3 **and** its QA short script.

**API / domain (table-driven)**

- New project V = 10. Five estimated Features totalling > 10: Current **short**, not over. Next Feature that would exceed 10 stays in the next Backlog iteration.
- Never split.
- Start a Backlog Feature that did not fit → Current, points may exceed 10, over-velocity badge.
- Accept one Feature → still Current, not Done.
- Advance test clock past `ends_on + 1 day 00:00` project TZ → that Feature in Done; V = that iteration’s accepted Feature points; `accepted_feature_points` frozen.
- Reorder / estimate / accept / start / icebox / length change → board already recomputed.
- Owner sets length 2–4 weeks → dates and pack change; completed iterations unchanged.
- Icebox never auto-fills into Current.
- Oversized Feature (`cost > V`): next empty iteration, over-V exception (velocity doc).
- All-zero last K completed → V returns to initial.
- Bugs/chores/releases not required; when present later they follow the velocity doc.

**QA**

- Run the seven-step script in `velocity-and-planning.md` (steps 5–6 that need Release / illegal drag may wait for slices 18 / 7 respectively; step 7 is this slice).
- Fail if a Recalculate button is required. Fail if V is typed in as a target.

### Isolation tests (ride 0, 1, 23)

From `multitenancy.md`: same title in org A and B, search/list only A; A fetches B’s id → 404; token from A on B → 404/`unauthorized`; attachment from A, signed-out → no file; viewer `POST /stories` → `forbidden`.

### Later phases (notes, not Phase 0 work)

| Slices | Test notes |
| --- | --- |
| 9 tasks | Toggle persists; complete-all does not Finish; accept **warns**; viewer read-only |
| 10 blockers | Free-text + optional same-project story; auto-resolve on accept/delete; warn on accept |
| 11 labels | Existing tables; `[a-z0-9-]+`; column max 100; do not exceed 100; filter does not change rank |
| 12 owners / mentions | Max 5; requester human; in-app + email; viewer mentionable, cannot follow |
| 13 Markdown / comments | No raw HTML; `javascript:` rejected; tombstone delete |
| 14 attachments | Type/size/count; auth GET; paste; missing-image after delete |
| 15 activity / undo | No reorder spam; undo only latest state change; undo is an activity |
| 16 keyboard | Every listed action has a chord from `ui.md` |
| 17 SSE | Two sessions < 2s; draft preserved; stale + reconnect |
| 18 types | Same machine package; Bug start without points; Chore finish=accept; Release colour vs **starts_on** |
| 19 epics | Purple label; independent order; progress Feature points only |
| 20 dates | `iterations.ends_on`; Icebox “Not scheduled”; no plan-overriding picker |
| 21 charts | Completed bars only; empty states; V line “initial 10” while N=0 |
| 22 search | Operators as specified; Done excluded unless `includedone:` |
| 23 agents | Agent-api QA short script (9 steps) — that file wins |
| 24 My Work / saved search | Owner/requester rules as specified |
| 25 workspaces | One organisation; not a permission boundary |
| 26 CSV | Owner only (assumption); create-only, all-or-nothing |
| 27 cycle time | First `started_at` → `accepted_at`; reject does not reset |
| 28 scales / bugs-chores toggle | No history conversion; toggle reversible |
| 29 split panes | Both live-update |
| 30 manual Current | Future iterations still auto-plan |

---

## Quality bar

From repo `AGENTS.md` and `api/internal/migrations/AGENTS.md`. Non-negotiable.

**Testing**

- TDD: failing test, then code. `nix-shell --run 'make test'` and `make lint` after every code change.
- 90%+ line coverage on API and frontend. Do not lower the gate.
- No placeholder, fake, or unused production branches. Tests exercise the same machine, planner, and auth as runtime.
- **No fallbacks.** Dummy password hashes, silent estimate-on-start, “if planner fails return last board”, catch-all `state = body.State`, or defaulting agent `story_type` to Feature are firable. Fail closed with a coded error.
- Do not ignore test, script, or lint errors.
- Documentation-only changes: review the diff; do not run the code gates unless code changed.

**Types and API**

- Go: wrap with `fmt.Errorf("%w")`. Structured Zap fields: `component`, `operation`, `organisation_id`, `project_id`, `actor_id`. Never log passwords, raw session/magic/agent/webhook secrets, or Authorization headers.
- Frontend: TypeScript strict; oxlint; `tsc --noEmit`. Types for board/story match the API.
- One error envelope. Agent codes are the codes.

**Migrations**

- Schema only. Data backfills run in Go on boot.
- No DB enums, triggers, or functions. Plural table names. UUID `id` PKs.
- Australian / New Zealand / UK names (`organisations`).
- After 000001, name files with `date +%Y%m%d%H%M%S`.
- Additive columns and new tables. Do not rewrite 000001.
- `point_scale` / `role` / `state` / `story_type` stay strings.

**Logging and config**

- Env from `.env.example`. No secrets in the repo.
- Clock, email, storage, bus are injected. Production has one real clock: `time.Now`.

**Feature flags and compatibility**

- No flags to hide Icebox, velocity, or accept-by-any-Member.
- Greenfield compatibility: old overview is wrong; do not keep `rejected` as a terminal peer to “match” it.
- Human SPA and agents share machines and rank rules.

**Authz checklist (every mutating handler)**

1. Authenticate (session or `flr_` bearer).
2. Resolve organisation from the project, not from a client-supplied org header alone.
3. Effective role or token scope.
4. Machine / rank / estimate rules.
5. Transaction: mutate, activity, pack, revision++.
6. Commit then bus.

**Spelling:** organisation, behaviour, colour (docs). Code identifiers for new tables: `organisations`, `organisation_id`.

---

## Implementation assumptions (forced by schema or stack)

These are not new product rules. Product remains unspecified where marked.

| Item | Status | Assumption if we must ship |
| --- | --- | --- |
| Username at signup | Locked 20 Aug 2026 | Infer from email local-part; uniquify; no username field in slice 0; editable later |
| Organisation public slug | Unspecified (fork 2) | UUID routes; no public org slug in Phase 0 |
| `display_name` | Unspecified | Copy username at create |
| Project TZ | Fork 3; velocity doc already says store + default Melbourne | `projects.timezone` default `Australia/Melbourne` |
| Icebox vs one rank | Assumed two lists (fork 4) | `rank_list` + unique `(project_id, rank_list, rank)` |
| Who creates projects after slice 0 | Assumed organisation owners | Enforce in Go |
| Reject reason vs comments table | Comments are slice 13 | Phase 0: activity only |
| HTTP prefix `/v1` vs `/api/v1` | Conflict: agent-api vs existing README | **`/api/v1`** mount; agent verbs/errors unchanged |
| Human session mechanism | Unspecified | Server-side cookie sessions |
| Live transport | Unspecified | SSE + in-process bus, slice 17 |
| Attachment / email vendors | Unspecified | Storage and email interfaces; local dir + SMTP/outbox |
| Rate-limit numbers | “TL picks” | 60 transitions/min/token, burst 10 |
| Webhook HMAC scheme | “TL specifies” | `t=<unix>,v1=<hmac-sha256>` |
| Agent user email | Unspecified | `agent-<uuid>@invalid`, `actor_kind=agent` |
| Magic/verify link TTL | Unspecified | 30 minutes |
| Password hash | Unspecified | bcrypt, cost 12, nullable column |
| Iteration start weekday storage | Implied Monday | ISO `1` = Monday on `projects.iteration_start_weekday` |
| `planned` state | Fork 9, later | Do not add in Phase 0 |

If Dan answers a fork, update this file in the same PR as the code change. Do not leave the doc lying.

---

## Required reading (implementers)

- This file
- `product-spec.md` (the slice you are in, plus machines)
- `velocity-and-planning.md` if the slice touches the plan (3, 4, 7, 8, 18, 20, 21, 30)
- `multitenancy.md` if the slice touches authz (0, 1, 23, 25)
- `agent-api.md` if the slice touches agents (23, and any transition error code)
- `open-questions.md` assumptions
- `000001_create_core_schema.up.sql`
- repo `AGENTS.md`, `api/internal/migrations/AGENTS.md`
- `docs/reference/frontend-design-guide.md` for UI work (do not invent a look)
- `ui.md` once the UI Designer has written it
