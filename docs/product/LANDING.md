# Landing this spec set in dmuso/flower

Local drafts live under `/workspace/flower-spec/`. They are **not** the repo paths.

The repo already has `docs/product/` and `docs/reference/`. New planning docs go in **feature-named folders**. Product intent stays in `docs/product/`. Do **not** overwrite `docs/reference/*` or `docs/migrations*`.

## Same-PR correction (required)

Update **`docs/product/overview.md`** in the same PR that adds these specs. The current overview is incomplete and slightly wrong versus classic Tracker. This spec set **supersedes** it.

The updated overview must say, in short form:

- Icebox = `unscheduled` holding pen. You pull from it. It is not a sequential stage on the way to Done.
- Feature / Bug / Chore / Release have **different** state machines (see [tracker-brief.md](./tracker-brief.md)).
- `rejected` is not a terminal peer of `accepted`. Reject, then Restart, returns the story to `started` in Current.
- Any Member can accept. Requester should. Viewers cannot.
- Velocity = accepted Feature points. Initial velocity 10. Auto-plan leaves Current short rather than overfill.
- Windows are **computed** (length in days). Flower does not persist Tracker-style iteration records as the plan.
- Tokens are **user-scoped** (`/api/v1/users/:id/tokens`).
- Point the reader at the feature folders below for the exhaustive rules.

Do not blindly overwrite overview with this whole spec. Keep overview as the one-page intent; make it **correct**.

Optional, same PR if the line is still wrong: root `README.md` sentence “Stories move through icebox, backlog, current iteration, and done” should be rephrased so Icebox is a holding pen, not a pipeline stage. That is a one-line product fix, not a reference overwrite.

## Path map

| Local draft | Eventual repo path | Notes |
| --- | --- | --- |
| `README.md` (this spec set) | `docs/product/spec-set.md` | How to read the set. Do not replace root `README.md`. |
| `tracker-brief.md` | `docs/product/tracker-brief.md` | Copy-exactly vs modernise. |
| `product-spec.md` | `docs/core-workflow/product-spec.md` | Whole-product slices + AC. |
| `domain-model.md` | `docs/core-workflow/domain-model.md` | Domain map. Technical Lead owns the tree. |
| `velocity-and-planning.md` | `docs/velocity-planning/velocity-and-planning.md` | Planning model (`pack` as a pure function). |
| `multitenancy.md` | `docs/multitenancy/multitenancy.md` | Organisations, roles, isolation. |
| `open-questions.md` | `docs/product/open-questions.md` | Forks + assumptions. |
| `LANDING.md` | `docs/product/LANDING.md` | Keep until the files have actually moved; then delete or fold into `docs/README.md`. |
| — | `docs/product/overview.md` | **Update in place** (correction above). |
| `ui.md` | `docs/core-workflow/ui.md` | In this PR. UI Designer file. Must cite `docs/reference/frontend-design-guide.md`. No new palette. |
| `technical-approach.md` | `docs/core-workflow/technical-approach.md` | In this PR. Technical Lead file. Plus per-feature approach notes if a slice needs them. |

Add links from `docs/README.md` to the new product and feature folders. Do not move reference docs.

## What must not be touched

- `docs/reference/frontend-design-guide.md` — look is locked (bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, column board).
- `docs/reference/technology-choices.md` and other `docs/reference/*`.
- `docs/migrations.md`, `docs/migration-usage.md`, `api/internal/migrations/*` (except when a later implementation slice adds a **new** schema-only migration).

## Existing schema (ground, not a redesign)

Already in `000001_create_core_schema`:

`users`, `projects`, `project_memberships`, `iterations`, `stories`, `labels`, `story_labels`, `activities`.

`iterations` and `stories.iteration_id` are leftover. Product does not persist iteration records as the plan. Do not treat 000001 as the planning model.

Not present: organisations, story owners, comments, tasks, attachments, epics, notifications, blockers, followers, API tokens, webhooks.

Technical Lead adds tables for slices that need them. They do not redesign the eight. Notable existing choices to keep:

- `stories.rank VARCHAR(64)` — fractional / lexicographic rank. Do not switch to integer priority.
- `stories.title VARCHAR(500)` — product max is 500, not a tighter invented limit.
- `story_type` and `state` are strings, not DB enums.
- `activities` attributes each change to a `users` row (000001 column name is leftover). Product language is **user**. An API token authenticates as that user.
- `projects` has `point_scale` and a leftover weeks-named length column. Product length is **days** (default 7). No organisation_id, no timezone, no velocity strategy, no bugs-and-chores-estimable flag yet.

## Ports (locked)

API `8180`, frontend `4273`, Postgres `5433` / test `5437`.

## Spelling

UK / AU / NZ in docs and in new identifiers: `organisations`, not `organizations`.

## Repo landing (for later push)

When these drafts land in the repo:

- Index `docs/core-workflow/domain-model.md` from `docs/README.md` and `docs/product/spec-set.md`.
- Do not add an org-level token mint path.
