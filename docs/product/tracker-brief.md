# Tracker: copy exactly vs modernise

Classic Pivotal Tracker is the product we are rebuilding. This note is short on purpose. Planning detail lives in [velocity-and-planning.md](../velocity-planning/velocity-and-planning.md). Types and slices live in [product-spec.md](../core-workflow/product-spec.md).

If a later doc invents a third behaviour, this list plus the two companion files win.

## Copy exactly

These are Tracker. Do not “improve” them in MVP.

- **Four types, four machines.** Feature and Bug share one. Chore is shorter. Release is a marker.
  - Feature / Bug: `unscheduled` → `unstarted` → `started` → `finished` → `delivered` → `accepted` or `rejected`. Reject, then **Restart**, → `started`. Stays in Current.
  - Chore: `unscheduled` → `unstarted` → `started` → `accepted`. No finished, delivered, or reject. Finish *is* accept.
  - Release: not work. Auto-`started` when created or dragged into Backlog. Finish → `accepted`. Optional **target date**. Marker sits at the **end** of that milestone’s stories. Blue if the iteration containing the marker starts on or before the target; red if that iteration **starts after** the target.
- **Icebox = unscheduled holding pen.** You pull *from* it into the ranked list. You do not walk through Icebox on the way to Done. Icebox is not in velocity math.
- **One ordered list.** Position is priority. Backlog = future iterations of that list. Current = this iteration. Done = accepted in **completed** iterations only. Accepted stays in Current until rollover at midnight, project timezone.
- **Unstarted cannot be dragged above started.**
- **Estimate to Start a Feature.** `0` is a valid estimate. Bugs and chores are unestimated by default and do not count toward velocity.
- **Auto-plan.** Pack Current toward velocity; **leave Current short** rather than overfill with the next story. Start on a Backlog or Icebox story still jumps it to Current (may overflow). Re-plan live on estimate, order, accept, velocity.
- **Velocity.** Strategy = last 3 completed iterations (project setting 1–4). **Initial velocity 10** until one iteration completes. Count **accepted Feature points only**.
- **Any Member can accept.** Owner / Member / Viewer only. Requester *should* accept; My Work surfaces their Delivered. No accept ACL in MVP. History is undo.
- **Owners: maximum 5.** Start assigns the clicker as an owner. Requester and owners auto-follow and cannot unfollow.
- **Tasks** are unowned, unpointed checklists. Incomplete tasks **warn** on Accept; they do not hard-block.
- **Blockers** are free-text plus an optional **story** link (not a random URL). Auto-resolve when the linked story is accepted or deleted. Warn on Accept; do not hard-block.
- **Epics** are one purple label plus an independent epic order and a progress view. Not a parent ticket.
- **Search** uses Tracker-like filters: `type:`, `state:`, `estimate:-1` (unestimated), `owner:`, `label:`, `is:blocked`, `includedone:`.
- **Default point scale** Linear `0,1,2,3`.
- **Manual planning** is an escape hatch for **Current only**, and it is later. Default is automatic. Future iterations stay auto-planned.

## Modernise (deliberate deltas)

Tracker got these wrong, or the world changed. We are explicit.

| Tracker | Flower |
| --- | --- |
| API token acts as the human who minted it | Agents are **first-class actors** with their own name, scoped token, and attribution. No impersonation. |
| “Bugs and chores may be estimated” could not be turned off | Same toggle is later, and it is **reversible**. |
| Custom point scale could not be reverted | Custom scale is later and **must be revertible**. Fibonacci (`0,1,2,3,5,8`) and Powers of 2 (`0,1,2,4,8`) are later project settings. |
| No first-class organisation tenant in the way we need | **Organisations** (UK spelling) are the tenant. Multitenant from day one. |
| Email/password era only | Email + password **and** magic link in MVP. SSO later. |
| REST only, webhooks as the user | REST + webhooks in the agent slice. GraphQL later unless cheap. |
| In-app + email; later Slack | Mentions: in-app + email in Phase 1. Slack later. |
| CSV / PT import as a pile | CSV import/export Phase 3. Pivotal Tracker import later **if at all**. |
| Team strength % in the velocity formula | **Not MVP.** Average accepted Feature points only. |
| Look-and-feel of Tracker | Flower’s locked look (bloom / stem / paper, Fraunces + Inter, Lucide, columns). See `docs/reference/frontend-design-guide.md`. |

## Never

- Custom fields.
- Gantt, dependency graphs, resource levelling, per-person velocity.
- Custom workflow engines or extra statuses (`in review`, `qa`).
- Treating Icebox as a sequential stage equivalent to Backlog.
- Treating `rejected` as a terminal peer of `accepted`.
- Inventing a milestone date picker that overrides the plan (a Release **target date** is a comparison, not a plan override).
- A new visual identity.

## Existing Flower docs this brief supersedes

`docs/product/overview.md` today lists `rejected` as a peer end-state and is thin on type-specific machines. Root `README.md` frames Icebox as a sequential stage (“stories move through icebox, backlog, current, done”). Both are wrong versus Tracker. [LANDING.md](./LANDING.md) requires the overview to be updated in the same PR as these specs.
