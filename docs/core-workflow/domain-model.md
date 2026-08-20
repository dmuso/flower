# Domain model

Eventual repo path: `docs/core-workflow/domain-model.md`

Local draft: `/workspace/flower-spec/domain-model.md`

Change name: `flower`

Implementer map: what each domain owns and what it must not own. Slice order stays in `product-spec.md`. Packing **math** stays in `velocity-and-planning.md`. Directory shape stays in `STRUCTURE-FROM-LAYOUT.md` and `technical-approach.md`.

One domain package = one bounded context = one vertical slice of HTTP + service + persistence. Entity types live in that domain (`types.go`). The repository interface lives next to its postgres implementation. Domains do not import each other; shared contracts go in `ports/`. `platform/` never contains story, planning, or tenancy rules. `app/` wires lifecycle. `handlers/` registers routes. There is no `internal/wire/` package.

Frontend: HTTP in `lib/api/<resource>.ts`. Thick-slice rules in `features/<slice>/domain/`. Screens in `pages/`.

Always **user**. Never “actor” as a type. Leftover column `activities.actor_id` stores a user id.

---

## user

**Owns:** the `users` row, signup / login / magic link / email verify, server-side sessions, **API tokens** (`POST/GET/revoke /api/v1/users/:id/tokens`).

**Invariants:** a token is a credential of a user, not an identity type. Bearer authenticates as that user. Role comes from project memberships (or a grant ≤ minter). Member cannot mint Owner. Viewer cannot mint. Session login is password or magic link; a token is not a login identity.

**Must not own:** organisations, stories, planning, webhooks.

---

## organisation

**Owns:** `organisations`, organisation memberships (MVP: owner), tenant boundary, 404-vs-403 when the organisation is missing.

**Invariants:** there is no single-tenant mode. Cross-tenant ids are `not_found`.

**Must not own:** project settings, stories, tokens.

---

## project

**Owns:** `projects`, `project_memberships`, invites, project settings (`timezone`, `velocity_strategy`, `initial_velocity`, `iteration_start_weekday`, **`iteration_length_days`**, point scale, slug unique per organisation).

**Invariants:** a project belongs to one organisation. Effective role = org owner or membership row. Owner / Member / Viewer is the whole permission model.

**Must not own:** story state machine, rank algorithm, velocity formula, iteration rows (there are none).

---

## board

**Owns:** the board read model (Icebox / Backlog / Current / Done payload). Calls planning to attach bands and V. Calls story for ordered lists.

**Invariants:** the SPA renders this payload. The client does not pack.

**Must not own:** persistence of “which iteration a story is in.” That field does not exist in the model.

---

## story

**Owns:** `stories` (except leftover `iteration_id`), labels / story_labels, owners, followers, type-specific **machine**, **rank** (`rank` + `rank_list`), estimate, revision, `started_at`, `accepted_at`. HTTP: create, patch fields, transitions, reorder.

**Invariants:** Feature start requires estimate (`0` allowed). Icebox Start is `unscheduled` → `started`. No `unstart`. `PATCH` does not write `state`. Member/Owner may accept/reject. Rank is `VARCHAR(64)` fractional. Icebox and ranked lists are independent.

**Must not own:** velocity math, window close, organisation tenancy, tokens. **Must not** persist a story → iteration link.

Leftover: `stories.iteration_id` in 000001. Stop writing it. Do not read it as planning truth.

---

## planning

**Owns:** the pack function, velocity formula, window-close writer for `velocity_samples`, board dates.

**Invariants:** planning / velocity is a **calculation, not an Iteration aggregate**. There is no Iteration entity and no `iterations` table in the model. The only length setting is project `iteration_length_days`. Bands (`current` / `next` / `later`) are recomputed whenever velocity, order, or estimates change. Stories do not point at samples.

`velocity_samples` (starts_on, ends_on, accepted_feature_points) is an input to V, written at window close. It is not an Iteration. It is not a home for stories.

**Must not own:** story machine, rank, HTTP for transitions, token mint.

Leftover: 000001 `iterations` table. Stop inserting. Migrate away; do not treat it as this domain.

---

## activity

**Owns:** `activities` rows (kind, summary, leftover `actor_id` = user, story_id). Undo reads the latest state-changing row (slice 15).

**Must not own:** permission checks (those are project effective role).

---

## Later domains (do not pre-create empty packages)

| Domain | Owns | Boundary |
| --- | --- | --- |
| comment | Markdown comments | No raw HTML |
| attachment | Files + auth GET | Not a public URL |
| webhook | Outbound project hooks | Delivery is not a command |
| search | Operators | Does not change rank |
| chart | Completed bars | Read-only |

---

## Hard boundaries (repeat)

1. **Iterations are not a domain and not a table.** Packing does not write iteration rows or `stories.iteration_id`.
2. **A token is not a user type.** It belongs to a user (`/api/v1/users/:id/tokens`).
3. **Organisation is the tenant.** Project belongs to one organisation.
4. **Story does not store which window it is in.** The board calculates current / next / later.
5. **Domains do not import each other.** Cross-domain contracts live in `ports/`.
6. **HTTP for a domain lives in that domain’s `handler.go`.** `handlers/` only registers routes.
