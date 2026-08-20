# Flower product overview

Flower is a task management workflow tool in the spirit of [Pivotal Tracker](https://www.pivotaltracker.com/).

The unit of work is a **story**. Teams plan by ranking stories, estimating them, and letting velocity fill iterations. There is no milestone-date picker that overrides the plan.

This page is the one-page intent. Exhaustive rules live in the spec set (see [How to read the spec set](spec-set.md)). If this page and the spec set disagree, the spec set wins.

## Core ideas

- **Organisations** are the tenant. They contain projects and membership.
- **Projects** contain the stories, iterations, labels, and membership for one team.
- **Stories** are the work. A story has a type (feature, bug, chore, release), an estimate, a requester, owners, labels, and a state.
- **Icebox** is an **unscheduled holding pen**. You pull from it into the ranked list. It is not a sequential stage on the way to Done, and it is not in the velocity math.
- **Backlog** is the ranked queue: future iterations of one ordered list. Position is priority.
- **Current** is this iteration: in-progress work, velocity-filled unstarted stories, and stories accepted *this* iteration.
- **Done** is accepted work from **completed** iterations only. Accepted stories stay in Current until rollover at midnight in the project timezone.
- **Iterations** are time-boxed measuring windows. Velocity is observed from accepted **Feature** points (initial velocity 10 until one iteration completes). Auto-plan leaves Current short rather than overfilling with the next story. Starting a story may overflow Current.
- **Labels** group stories. An epic is one purple label plus an independent epic order, not a parent ticket.
- **Activity** records what changed, who changed it, and when. History is undo.

## Story types and state machines

Types have **different** machines. The API owns these values as strings; they are not database enums.

### Feature and Bug

`unscheduled` → `unstarted` → `started` → `finished` → `delivered` → `accepted` or `rejected`.

Reject, then **Restart**, returns the story to `started`. It stays in Current. `rejected` is not a terminal peer of `accepted`.

A Feature cannot be started without an estimate (`0` is a valid estimate). Bugs are unestimated by default and may start without points. Bug points do not count toward velocity.

### Chore

`unscheduled` → `unstarted` → `started` → `accepted`. Finish *is* accept. No delivered, no reject. Unestimated. Does not count toward velocity.

### Release

A marker, not work. Auto-started when created or dragged into Backlog. Finish → accepted. Optional target date. Place the marker at the **end** of that milestone’s stories. Blue if the iteration containing the marker starts on or before the target; red if that iteration starts after the target.

## Who can accept

Roles are Owner, Member, Viewer. **Any Member or Owner can accept.** The requester *should*. Viewers are read-only. There is no accept ACL in MVP. History undoes a mistake.

Agents authenticate as named actors with a bearer token bound to a project membership, and use the same API as humans. Permissions are the membership role (Owner / Member / Viewer). A Member agent can accept and reject.

## What the core schema covers

The first migration creates:

- `users`
- `projects`
- `project_memberships`
- `iterations`
- `stories`
- `labels`
- `story_labels`
- `activities`

Business rules (allowed state transitions, who can accept a story, how ranking works) belong in the Go API, not in the database.

## Spec set

- [How to read the spec set](spec-set.md)
- [Tracker: copy exactly vs modernise](tracker-brief.md)
- [Product spec (slices + acceptance criteria)](../core-workflow/product-spec.md)
- [UI](../core-workflow/ui.md)
- [Technical approach](../core-workflow/technical-approach.md)
- [Velocity and planning](../velocity-planning/velocity-and-planning.md)
- [Multitenancy](../multitenancy/multitenancy.md)
- [Agent API](../agent-api/agent-api.md)
- [Open questions](open-questions.md)
