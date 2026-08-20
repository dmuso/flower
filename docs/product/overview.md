# Flower product overview

Flower is a task management workflow tool in the spirit of [Pivotal Tracker](https://www.pivotaltracker.com/).

The unit of work is a **story**. Teams plan by ranking stories, estimating them, and moving them through a small set of panels that match how software actually gets delivered.

## Core ideas

- **Projects** contain the stories, iterations, labels, and membership for one team.
- **Stories** are the work. A story has a type (feature, bug, chore, release), an estimate, a requester, owners, labels, and a state.
- **Icebox** holds stories that are not yet scheduled.
- **Backlog** is the ranked queue of stories that will be scheduled next.
- **Current** is the iteration currently in flight.
- **Done** is accepted work from recent iterations.
- **Iterations** are time-boxed planning windows. Velocity is derived from accepted points, not entered as a target.
- **Labels** group stories across the board.
- **Activity** records what changed, who changed it, and when.

## Story states

Stories move through a linear workflow. The API owns these values as strings; they are not database enums.

| State | Meaning |
| --- | --- |
| `unscheduled` | In the icebox. |
| `unstarted` | Ranked in the backlog or current iteration, not started. |
| `started` | Someone is working on it. |
| `finished` | Implementation is complete. |
| `delivered` | Ready for acceptance. |
| `accepted` | The requester accepted the work. |
| `rejected` | The requester sent it back. |

## Story types

| Type | Meaning |
| --- | --- |
| `feature` | User-facing work, estimated in points. |
| `bug` | Defect. Unestimated unless the project requires it. |
| `chore` | Supporting work. Unestimated. |
| `release` | A marker in the backlog, not a unit of work. |

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
