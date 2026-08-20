# Domain model

Eventual repo path: `docs/core-workflow/domain-model.md`

Local draft: `/workspace/flower-spec/domain-model.md`

Change name: `flower`

Implementer map: what each domain owns and what it must not own. Slice order stays in `product-spec.md`. Packing math stays in `velocity-and-planning.md`. Directory shape stays in `STRUCTURE-FROM-LAYOUT.md` and `technical-approach.md`.

One domain package = one bounded context = one vertical slice of HTTP + service + persistence. Entity types live in that domain (`types.go`). The repository interface lives next to its postgres implementation. Domains do not import each other; shared contracts go in `ports/`. `platform/` never contains story, planning, or tenancy rules. `app/` wires lifecycle. `handlers/` registers routes. There is no `internal/wire/` package.

Frontend: HTTP in `lib/api/<resource>.ts`. Thick-slice rules in `features/<slice>/domain/`. Screens in `pages/`.

Always **user**. `activities.user_id` is the user (session or the user a token authenticates as). A token is a user credential, not an identity type. Planning / velocity is a **calculation**, not a time-box aggregate. There is no iteration table and no story-to-window foreign key.

---

## user

**Owns:** the `users` row, signup / login / magic link / email verify, server-side sessions, and API tokens (`POST` / `GET` / revoke `/api/v1/users/:id/tokens`).

**Invariants:** a token authenticates as `api_tokens.user_id`. It is not a login identity. Role is that user’s project memberships (or an optional cap ≤ minter). Member cannot mint Owner. Viewer cannot mint. No organisation-collection token API.

**Must not own:** organisations, stories, planning, webhooks.

---

## organisation

**Owns:** `organisations`, organisation memberships (MVP: owner), tenant boundary, 404-vs-403 when the organisation is missing.

**Invariants:** there is no single-tenant mode. Cross-tenant ids are `not_found`.

**Must not own:** project settings, stories, tokens.

---

## project

**Owns:** `projects`, `project_memberships`, invites, project settings (`timezone`, `velocity_strategy`, `initial_velocity`, `iteration_start_weekday`, **`iteration_length_days`**, point scale, slug unique per organisation).

**Invariants:** a project belongs to one organisation. Effective role = org owner or membership row. Owner / Member / Viewer is the whole permission model. The only stored window length is a positive integer of **days** (default 7).

**Must not own:** story state machine, rank algorithm, velocity formula, a time-box entity.

---


## story

**Owns:** `stories`, labels / story_labels, owners, followers, type-specific **machine**, **rank** (`rank` + `rank_list`), estimate, revision, `started_at`, `accepted_at`. HTTP: create, patch fields, transitions, reorder.

**Invariants:** Feature start requires estimate (`0` allowed). Icebox Start is `unscheduled` → `started`. No `unstart`. `PATCH` does not write `state`. Member/Owner may accept/reject. Rank is `VARCHAR(64)` fractional. Icebox and ranked lists are independent. Unique `(project_id, rank_list, rank)`. Stories do not store which window or band they are in.

**Must not own:** velocity math, organisation tenancy, tokens. Must not persist a story → window link. Accepting a Feature sets `accepted_at`; it does not store points.

---

## planning

**Owns:** two pure functions — `V(...)` and `pack(ordered stories, V, L, now, TZ)` — and the computed window dates. Nothing persisted.

**Invariants:** planning / velocity is a **calculation**. There is no Iteration entity. The only stored settings are project `iteration_length_days`, start weekday, timezone, `velocity_strategy` (K), and `initial_velocity`. Windows and bands are computed and drawn in the UI. Recompute whenever stories, estimates, accepts, or settings change.

**V** is duration-based. It is calculated from completed Features’ **start/end datetimes** (`started_at` → `accepted_at`), not from summing estimates in a window. Exact rollup (including how K is applied) lives in `velocity-and-planning.md`.

**Predicted duration** for an incomplete story: from completed Features of the **same estimate**. What is projected to finish in the current window is V plus those predicted durations, packed in rank order. We accept stories, not points.

**Cold start:** until one Feature has both `started_at` and `accepted_at`, `pack` uses `initial_velocity` as estimate-points that fit in a window. After that, the duration model takes over.

**Window end** is a clock crossing, not a write.

**Must not own:** story machine, rank, HTTP for transitions, token mint. Must not persist V, window totals, or accepted points.

## activity

**Owns:** `activities` (`kind`, `summary`, `user_id`, `story_id`). Undo reads the latest state-changing row (slice 15).

**Must not own:** permission checks (those are project effective role).

---

## tenancy (organisation + project)

**Owns:** effective role, 404-vs-403, “does this user belong to this organisation / project?”

**Hard boundaries:** cross-tenant id → `not_found`. Same-tenant insufficient role → `forbidden`. Cookie vs Bearer is the only auth difference after the user is resolved.

---

## Later domains (do not pre-create empty packages)

| Domain | Owns | Boundary |
| --- | --- | --- |
| comment | Markdown comments | No raw HTML |
| attachment | Files + auth GET | Not a public URL |
| webhook | Outbound project hooks | Delivery is not a command |
| search | Operators | Does not change rank |
| chart | Completed bars from the same live V sums | Read-only |

---

## What is not a domain

| Not a domain | Why |
| --- | --- |
| Iteration / time-box table | Planning is a calculation. Settings + `pack` + live V. |
| Board | A **frontend projection**. The API returns stories plus pack fields (V, bands, dates). The SPA draws Icebox / Backlog / Current / Done. No `board` domain package. No `/board` API. |
| API token as an identity type | Token is a user credential on `user`. |
| `app/` use-case packages | `app/` is lifecycle + wiring only. |
| `api/internal/wire/` | Does not exist. Do not invent it. |
| A single frontend client file | HTTP is `lib/api/<resource>.ts`. |

---

## Dependency direction

```
domain/ ──────► ports/ ◄────── platform/
    │                              │
    └──────────► types/ ◄──────────┘
```

`app/` wires implementations into ports and starts the process. `handlers/` mounts each domain’s `handler.go` on `/api/v1`.
