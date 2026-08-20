# Flower product overview

Flower is a task management workflow tool in the spirit of [Pivotal Tracker](https://www.pivotaltracker.com/).

The unit of work is a **story**. Teams plan by ranking stories, estimating them, and letting velocity fill Current and future visual bands. There is no milestone-date picker that overrides the plan. There is no persisted iteration record as the plan.

This page is the one-page intent. Exhaustive rules live in the spec set (see [How to read the spec set](spec-set.md)). If this page and the spec set disagree, the spec set wins.

## Core ideas

- **Organisations** are the tenant. They contain projects and membership.
- **Projects** contain the stories, labels, membership, and planning settings (length in days, timezone, velocity) for one team.
- **Stories** are the work. A story has a type (feature, bug, chore, release), an estimate, a requester, owners, labels, and a state.
- **Icebox** is an **unscheduled holding pen**. You pull from it into the ranked list. It is not a sequential stage on the way to Done, and it is not in the velocity math.
- **Backlog** is the ranked queue: later **computed bands** of one ordered list. Position is priority. Bands are visual, not stored rows.
- **Current** is this window: the head of the ranked list that fits this window’s velocity, plus in-progress work and stories accepted *this* window.
- **Done** is accepted work that has **aged past the current window**. Flat list, newest accepted first. Accepted stories stay in Current until the window ends at midnight in the project timezone.
- **Windows** are computed from length in days (default 7), start weekday, now, and timezone. Velocity is observed from accepted **Feature** points (initial velocity 10 until one window completes). Auto-plan leaves Current short rather than overfilling with the next story. Starting a story may overflow Current. Stories are not assigned to a window row.
- **Labels** group stories. An epic is one purple label plus an independent epic order, not a parent ticket.
- **Activity** records what changed, **who** (the user) changed it, and when. History is undo.

## Story types and state machines

Types have **different** machines. The API owns these values as strings; they are not database enums.

### Feature and Bug

`unscheduled` → `unstarted` → `started` → `finished` → `delivered` → `accepted` or `rejected`.

Reject, then **Restart**, returns the story to `started`. It stays in Current. `rejected` is not a terminal peer of `accepted`.

A Feature cannot be started without an estimate (`0` is a valid estimate). Bugs are unestimated by default and may start without points. Bug points do not count toward velocity.

### Chore

`unscheduled` → `unstarted` → `started` → `accepted`. Finish *is* accept. No delivered, no reject. Unestimated. Does not count toward velocity.

### Release

A marker, not work. Auto-started when created or dragged into Backlog. Finish → accepted. Optional target date. Place the marker at the **end** of that milestone’s stories. Blue if the **computed window** containing the marker starts on or before the target; red if that window starts after the target.

## Who can accept

Roles are Owner, Member, Viewer. **Any Member or Owner can accept.** The requester *should*. Viewers are read-only. There is no accept ACL in MVP. History undoes a mistake.

The HTTP API is the same for the app and for tokens. Humans use a session cookie or a Bearer **user API token** at `/api/v1/users/:id/tokens` (role at or below their own; Member cannot mint Owner). Permissions are Owner / Member / Viewer. Tokens are not org-level.

## What the core schema covers

The first migration creates:

- `users`
- `projects`
- `project_memberships`
- `iterations` (leftover — not the planning model)
- `stories`
- `labels`
- `story_labels`
- `activities`

Business rules (allowed state transitions, who can accept a story, how ranking and `pack` work) belong in `api/internal/domain/<domain>`, not in the database. See [domain-model.md](../core-workflow/domain-model.md). Existing `iterations` / `stories.iteration_id` are leftover; stop using them as the plan. Length is **days** on the project.

## Spec set

- [How to read the spec set](spec-set.md)
- [Tracker: copy exactly vs modernise](tracker-brief.md)
- [Product spec (slices + acceptance criteria)](../core-workflow/product-spec.md)
- [Domain model](../core-workflow/domain-model.md)
- [UI](../core-workflow/ui.md)
- [Technical approach](../core-workflow/technical-approach.md)
- [Velocity and planning](../velocity-planning/velocity-and-planning.md)
- [Multitenancy](../multitenancy/multitenancy.md)
- [Open questions](open-questions.md)
