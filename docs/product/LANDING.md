# Landing this spec set in dmuso/flower

Local drafts live under `/workspace/flower-spec/`. They are **not** the repo paths.

The repo already has `docs/product/` and `docs/reference/`. New planning docs go in **feature-named folders**. Product intent stays in `docs/product/`. Do **not** overwrite `docs/reference/*` or `docs/migrations*`.

## Same-PR overview

Update **`docs/product/overview.md`** in the same PR that adds these specs so it matches this spec set. Keep overview as the one-page intent.

The overview must say, in short form:

- Icebox = `unscheduled` holding pen. You pull from it. It is not a sequential stage on the way to Done.
- Feature / Bug / Chore / Release have **different** state machines (see [tracker-brief.md](./tracker-brief.md)).
- `rejected` is not a terminal peer of `accepted`. Reject, then Restart, returns the story to `started` in Current.
- Any Member can accept. Requester should. Viewers cannot.
- Velocity is live from completed Feature start/end times (`started_at` → `accepted_at`) plus similar-size predicted duration. Stories accepted in the open window are not in the lookback. Until at least one corpus Feature exists in a completed window (or `time == 0`), `velocity` is undefined and pack uses initial velocity 10. Auto-plan leaves Current short rather than overfill.
- Windows are **computed** (length in days). Flower does not persist Tracker-style iteration records as the plan.
- Mint / list / revoke only on your own user: `GET|POST|DELETE /api/v1/users/:id/tokens`. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. No org-level mint path.
- Point the reader at the feature folders below for the exhaustive rules.

Do not blindly overwrite overview with this whole spec. Keep overview as the one-page intent; make it match the rules above.

Root `README.md` should describe Icebox as a holding pen, not a pipeline stage. That is a one-line product fix, not a reference overwrite.

## Path map

| Local draft | Eventual repo path | Notes |
| --- | --- | --- |
| `README.md` (this spec set) | `docs/product/spec-set.md` | How to read the set. Do not replace root `README.md`. |
| `tracker-brief.md` | `docs/product/tracker-brief.md` | Copy-exactly vs modernise. |
| `product-spec.md` | `docs/core-workflow/product-spec.md` | Whole-product slices + AC. |
| `domain-model.md` | `docs/core-workflow/domain-model.md` | Domain map. Technical Lead owns the tree. |
| `velocity-and-planning.md` | `docs/velocity-planning/velocity-and-planning.md` | Planning model (live `velocity` + `predicted_duration` + `pack`). |
| `multitenancy.md` | `docs/multitenancy/multitenancy.md` | Organisations, roles, isolation. |
| `open-questions.md` | `docs/product/open-questions.md` | Forks + assumptions. |
| `LANDING.md` | `docs/product/LANDING.md` | Keep until the files have actually moved; then delete or fold into `docs/README.md`. |
| — | `docs/product/overview.md` | **Update in place** (overview content above). |
| `ui.md` | `docs/core-workflow/ui.md` | In this PR. UI Designer file. Must cite `docs/reference/frontend-design-guide.md`. No new palette. |
| `technical-approach.md` | `docs/core-workflow/technical-approach.md` | In this PR. Technical Lead file. Plus per-feature approach notes if a slice needs them. |

Add links from `docs/README.md` to the new product and feature folders. Do not move reference docs.

## What must not be touched

- `docs/reference/frontend-design-guide.md` — look is locked (bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, column board).
- `docs/reference/technology-choices.md` and other `docs/reference/*`.
- `docs/migrations.md`, `docs/migration-usage.md`, `api/internal/migrations/*` (except when a later implementation slice adds a **new** schema-only migration).

## Schema (must match this model)

Core tables: `users`, `projects`, `project_memberships`, `stories`, `labels`, `story_labels`, `activities`.

Planning is a calculation. There is no `iterations` table. Velocity is live from stories; it is not persisted. Project length is `iteration_length_days` (default 7).

Slices add tables they need: organisations, story owners, comments, tasks, attachments, epics, notifications, blockers, followers, API tokens, webhooks.

Notable choices:

- `stories.rank VARCHAR(64)` — fractional / lexicographic rank. Do not switch to integer priority.
- `stories.title VARCHAR(500)` — product max is 500, not a tighter invented limit.
- `story_type` and `state` are strings, not DB enums.
- `activities.user_id` attributes each change to a `users` row. Product language is **user**. An API token authenticates as that user.
- `projects` has `point_scale` and `iteration_length_days` (default 7). Organisation, timezone, velocity strategy, and bugs-and-chores-estimable land with their slices.

## Ports (locked)

API `8180`, frontend `4273`, Postgres `5433` / test `5437`.

## Spelling

UK / AU / NZ in docs and in new identifiers: `organisations`, not `organizations`.

## Repo landing (for later push)

When these drafts land in the repo:

- Index `docs/core-workflow/domain-model.md` from `docs/README.md` and `docs/product/spec-set.md`.
- Mint / list / revoke only on your own user: `GET|POST|DELETE /api/v1/users/:id/tokens`. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. No org-level mint path.
