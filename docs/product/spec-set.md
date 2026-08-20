# Flower

Change name: `flower`

This is the Product Owner spec set. It is not an implementation. It does not design the architecture. It does not invent a visual identity.

Eventual repo: [github.com/dmuso/flower](https://github.com/dmuso/flower). How each file lands in that repo is in [LANDING.md](./LANDING.md). **Do not** overwrite `docs/reference/*` or `docs/migrations*`.

## How to read this spec set

1. **[tracker-brief.md](./tracker-brief.md)** — what we copy from classic Pivotal Tracker and what we modernise. Read this first if a rule feels surprising.
2. **This file** — principles, constraints, index.
3. **[product-spec.md](./product-spec.md)** — problem, personas, out of scope, ordered vertical slices, acceptance criteria, required agents.
4. **[velocity-and-planning.md](./velocity-and-planning.md)** — the planning model. If a slice disagrees with it, this file wins until Dan changes it.
5. **[multitenancy.md](./multitenancy.md)** — accounts, organisations, projects, roles, isolation.
6. **[open-questions.md](./open-questions.md)** — true forks only, plus the assumption list.
7. **[LANDING.md](./LANDING.md)** — repo path map and the required `docs/product/overview.md` correction.

Then, in the delivery workflow (do not skip):

- **Reviewer** clears this spec set before anyone implements.
- **UI Designer** writes UI notes for new UI. Look is locked: [docs/reference/frontend-design-guide.md](../flower-existing-docs/docs/reference/frontend-design-guide.md) (bloom `#C43B6E`, stem `#2F7D4A`, paper `#FBF7F2`, Fraunces + Inter, Lucide, four-column board). Do not invent a new identity. Shortcut table lives in their `ui.md`.
- **Technical Lead** writes `technical-approach.md`. Stack, ports, and the eight core tables are constraints.
- **Developer** implements one slice at a time.
- **QA** tests each slice from its acceptance criteria alone.

Companion docs own planning and tenancy. Slices must still be independently demoable.

## Executive summary

Flower is a Pivotal-Tracker-like tracker for small product teams, including the scripts and CI that call the same API the app uses.

There is one ordered list. Position is priority. Icebox is the **unscheduled holding pen**, not a stage you pass through on the way to Done. Backlog is future iterations of the same list. Current is this iteration: in-progress work, velocity-filled unstarted stories, and stories accepted *this* iteration. Done is accepted work from **completed** iterations only. Accepted stories stay in Current until rollover at midnight in the project timezone.

Four story types, **different state machines** (classic Tracker):

- **Feature / Bug:** unscheduled → unstarted → started → finished → delivered → accepted or rejected. Reject, then Restart, returns the story to started. It stays in Current.
- **Chore:** unscheduled → unstarted → started → accepted. No finished, delivered, or reject.
- **Release:** a marker, not work. Auto-started when created or dragged into Backlog. Finish → accepted. Optional target date. Place the marker at the **end** of that milestone’s stories. Blue if on track versus the date; red if the iteration that contains the marker **starts** after the target.

Humans estimate Features (0 is valid). A Feature cannot be started without an estimate. Bugs and chores are unestimated by default and do not count toward velocity. Velocity is accepted **Feature** points only. A new project uses **initial velocity 10** until one iteration completes, then a rolling average of the last 3 completed iterations (setting: 1–4). Auto-plan fills Current up to velocity and **leaves Current short** rather than overfilling with the next story. Starting a Backlog or Icebox story jumps it to Current and may overflow.

Any **Member** or Owner can accept. The requester *should*. My Work surfaces their Delivered stories. There is no accept ACL lock in MVP. History is undo. Viewers are read-only.

The HTTP API is the same for the app and for tokens. Humans use a session cookie or a Bearer **API token**. A token is minted for a user/membership with a role at or below the minter (Member cannot mint Owner). Owner / Member / Viewer is the whole model. There is no agent product and no `/agents` endpoint.

## Product principles

Steal the power of classic Pivotal Tracker. Refuse Jira, Monday, and late-era issue-tracker chrome. Keep Flower’s already-locked look.

1. **One ordered list.** Priority is position. Icebox is unscheduled, not in the plan.
2. **The computer does the math.** Humans order and estimate. Flower fills iterations from velocity. Re-plan is live.
3. **Velocity is observed Feature points.** Initial 10, then the rolling average. We do not type a date to look prepared. Team strength % is later, not MVP.
4. **Small scale.** Default Linear 0 / 1 / 2 / 3. `0` is estimated. Unestimated is different and cannot Start a Feature.
5. **Accept is a team verb.** Any Member may accept. Requester should. History undoes a mistake.
6. **The board is the meeting.** Icebox / Backlog / Current / Done plus the open story.
7. **Keyboard is first-class.** Daily actions have a chord. The UI Designer owns the table.
8. **One API.** Cookie or Bearer; same handlers; Owner / Member / Viewer. No agent product.
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
| Schema already present | `users`, `projects`, `project_memberships`, `iterations`, `stories`, `labels`, `story_labels`, `activities` |
| Not present yet | story owners, comments, tasks, attachments, epics, organisations, notifications |
| Spelling | UK / AU / NZ (`organisations`, not `organizations`) |
| Migrations | Schema-only. No DB enums, triggers, or functions. Business rules in the Go API. |

Existing tables are **ground**, not a redesign brief. New slices add tables (owners, comments, tasks, …) and add `organisations` without rewriting the eight.

Architecture lives in `technical-approach.md` after the Reviewer clears this set.

## File index

| File | Owns |
| --- | --- |
| [README.md](./README.md) | How to read, summary, principles, constraints |
| [tracker-brief.md](./tracker-brief.md) | Copy-exactly vs modernise |
| [product-spec.md](./product-spec.md) | Problem, personas, slices, AC |
| [velocity-and-planning.md](./velocity-and-planning.md) | Iteration math, packing, releases, charts, examples |
| [multitenancy.md](./multitenancy.md) | Organisations, roles, isolation, workspaces |
| [open-questions.md](./open-questions.md) | True forks + baked assumptions |
| [LANDING.md](./LANDING.md) | Eventual repo paths; overview.md correction |

## Keyboard shortcuts (stub for UI Designer)

Product rules: [product-spec.md](./product-spec.md) slice 16. Full chord table belongs in `ui.md` (board / story / dialog, browser conflicts, one-handed start / finish / deliver / accept / estimate / reorder / search). Implementers must not ship a hidden private keymap.

## What “done” means for this spec

A Reviewer can clear it. A UI Designer can write `ui.md` from it **and** the locked frontend guide, without inventing product rules or a new look. A Technical Lead can write `technical-approach.md` without inventing planning or tenancy. QA can fail a slice from its AC alone.
