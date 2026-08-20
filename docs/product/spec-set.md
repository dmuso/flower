# Flower

Change name: `flower`

This is the Product Owner spec set. It is not an implementation. It does not design the architecture. It does not invent a visual identity.

Eventual repo: [github.com/dmuso/flower](https://github.com/dmuso/flower). How each file lands in that repo is in [LANDING.md](./LANDING.md). **Do not** overwrite `docs/reference/*` or `docs/migrations*`.

## How to read this spec set

1. **[tracker-brief.md](./tracker-brief.md)** — what we copy from classic Pivotal Tracker and what we modernise. Read this first if a rule feels surprising.
2. **This file** — principles, constraints, index.
3. **[product-spec.md](./product-spec.md)** — problem, personas, out of scope, ordered vertical slices, acceptance criteria, required agents.
4. **[domain-model.md](./domain-model.md)** — what each domain owns. Technical Lead owns the tree (`api/internal/domain/<domain>`; frontend split by domain).
5. **[velocity-and-planning.md](./velocity-and-planning.md)** — the planning model (`pack` as a pure function). If a slice disagrees with it, this file wins.
6. **[multitenancy.md](./multitenancy.md)** — users, organisations, projects, roles, isolation.
7. **[open-questions.md](./open-questions.md)** — true forks only, plus the assumption list.
8. **[LANDING.md](./LANDING.md)** — repo path map and the required `docs/product/overview.md` correction.

Then, in the delivery workflow (do not skip):

- **Reviewer** clears this spec set before anyone implements.
- **UI Designer** writes UI notes for new UI. Look is locked: [docs/reference/frontend-design-guide.md](../flower-existing-docs/docs/reference/frontend-design-guide.md) (bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, four-column board). Do not invent a new identity. Shortcut table lives in their `ui.md`.
- **Technical Lead** writes `technical-approach.md` and owns [domain-model.md](./domain-model.md). Stack, ports, and the core schema are constraints.
- **Developer** implements one slice at a time.
- **QA** tests each slice from its acceptance criteria alone.

Companion docs own planning and tenancy. Slices must still be independently demoable.

## Executive summary

Flower is a Pivotal-Tracker-like tracker for small product teams, including the scripts and CI that call the same API the app uses.

There is one ordered list. Position is priority. Icebox is the **unscheduled holding pen**, not a stage you pass through on the way to Done. Backlog is later **computed bands** of the same list. Current is this window: the head of the ranked list that fits this window’s velocity, plus in-progress work and stories accepted *this* window. Done is accepted work that has **aged past the current window** (flat list). Accepted stories stay in Current until the window ends at midnight in the project timezone.

Four story types, **different state machines** (classic Tracker):

- **Feature / Bug:** unscheduled → unstarted → started → finished → delivered → accepted or rejected. Reject, then Restart, returns the story to started. It stays in Current.
- **Chore:** unscheduled → unstarted → started → accepted. No finished, delivered, or reject.
- **Release:** a marker, not work. Auto-started when created or dragged into Backlog. Finish → accepted. Optional target date. Place the marker at the **end** of that milestone’s stories. Blue if on track versus the date; red if the **computed window** that contains the marker **starts** after the target.

Humans estimate Features (0 is valid). A Feature cannot be started without an estimate. Bugs and chores are unestimated by default and do not count toward velocity. Velocity is accepted **Feature** points only. A new project uses **initial velocity 10** until one window completes, then a rolling average of the last 3 completed windows (setting: 1–4). Auto-plan fills Current up to velocity and **leaves Current short** rather than overfilling with the next story. Starting a Backlog or Icebox story jumps it to Current and may overflow.

Length is **days** (default 7), stored on the project. Flower does not persist Tracker-style iteration records as the plan. `pack` is a pure function. Stories are not assigned to a window row.

Any **Member** or Owner can accept. The requester *should*. My Work surfaces their Delivered stories. There is no accept ACL lock in MVP. History is undo. Viewers are read-only.

The HTTP API is the same for the app and for tokens. Humans use a session cookie or a Bearer **user API token** (`/api/v1/users/:id/tokens`). Role at or below their own (Member cannot mint Owner). Owner / Member / Viewer is the whole model. Tokens are not org-level.

## Product principles

Steal the power of classic Pivotal Tracker. Refuse Jira, Monday, and late-era issue-tracker chrome. Keep Flower’s already-locked look.

1. **One ordered list.** Priority is position. Icebox is unscheduled, not in the plan.
2. **The computer does the math.** Humans order and estimate. Flower packs Current and future bands from velocity. Re-plan is live.
3. **Velocity is observed Feature points.** Initial 10, then the rolling average. We do not type a date to look prepared. Team strength % is later, not MVP.
4. **Small scale.** Default Linear 0 / 1 / 2 / 3. `0` is estimated. Unestimated is different and cannot Start a Feature.
5. **Accept is a team verb.** Any Member may accept. Requester should. History undoes a mistake.
6. **The board is the meeting.** Icebox / Backlog / Current / Done plus the open story.
7. **Keyboard is first-class.** Daily actions have a chord. The UI Designer owns the table.
8. **One API.** Cookie or Bearer; same handlers; Owner / Member / Viewer. Tokens live on the user.
9. **Say no.** No Gantt, no custom fields, no workflow engine, no resource management.
10. **Smallest valuable slice.** Thin, end-to-end, user-visible. Not “API first.”

## Constraints (not a design)

Treat as given. Do not redesign in this spec.

| Kind | Locked value |
| --- | --- |
| Database | PostgreSQL 17 |
| API | Go + Gin, host port **8180** |
| Frontend | SolidJS + Bun + TypeScript + Tailwind v4, host port **4273** |
| Local Postgres | **5433** (dev), **5437** (test) |
| Shape | Monorepo, Nix Shell + Docker Compose + Make |
| Tenancy | Multitenant from day one (organisations) |
| Look | bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, column board — `docs/reference/frontend-design-guide.md` |
| Core schema | `users`, `projects`, `project_memberships`, `stories`, `labels`, `story_labels`, `activities`. `projects.iteration_length_days`. `activities.user_id`. No `iterations` table. |
| Added by slices | organisations, story owners, comments, tasks, attachments, epics, `velocity_samples`, `api_tokens`, notifications |
| Domain | `api/internal/domain/<domain>`; frontend split by domain — [domain-model.md](./domain-model.md) |
| Spelling | UK / AU / NZ (`organisations`, not `organizations`) |
| Migrations | Schema-only. No DB enums, triggers, or functions. Business rules in the Go domain packages. |

Schema must match this model. Planning is a calculation: no `iterations` table. Accepted-points history lives in `velocity_samples`. New slices add tables (owners, comments, tasks, organisations, tokens) without inventing a second planning store.

Architecture lives in `technical-approach.md` after the Reviewer clears this set.

## File index

| File | Owns |
| --- | --- |
| [README.md](./README.md) | How to read, summary, principles, constraints |
| [tracker-brief.md](./tracker-brief.md) | Copy-exactly vs modernise |
| [product-spec.md](./product-spec.md) | Problem, personas, slices, AC |
| [domain-model.md](./domain-model.md) | Domain map (Technical Lead) |
| [velocity-and-planning.md](./velocity-and-planning.md) | Window clock, `pack`, releases, charts, examples |
| [multitenancy.md](./multitenancy.md) | Organisations, roles, isolation, workspaces |
| [open-questions.md](./open-questions.md) | True forks + baked assumptions |
| [LANDING.md](./LANDING.md) | Eventual repo paths; overview.md correction |

## Keyboard shortcuts (stub for UI Designer)

Product rules: [product-spec.md](./product-spec.md) slice 16. Full chord table belongs in `ui.md` (board / story / dialog, browser conflicts, one-handed start / finish / deliver / accept / estimate / reorder / search). Implementers must not ship a hidden private keymap.

## What “done” means for this spec

A Reviewer can clear it. A UI Designer can write `ui.md` from it **and** the locked frontend guide, without inventing product rules or a new look. A Technical Lead can write `technical-approach.md` without inventing planning or tenancy. QA can fail a slice from its AC alone.
