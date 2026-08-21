# Open questions

Only **true forks** — things we cannot assume. Everything else is decided. If it is in the assumption list, implementers treat it as law.

## True forks

1. **~~User API tokens.~~ Locked.** Not a fork. Mint / list / revoke only on your own user: `GET|POST|DELETE /api/v1/users/:id/tokens`. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. Cookie vs Bearer is the only auth difference. Same `/api/v1` handlers. No org-level mint path. See Decisions locked.
2. **~~Dynamic windows.~~ Locked.** Not a fork. Iterations are not a domain. Planning is a calculation. There is no `iterations` table, no story `iteration_id`, and no `iteration_length_weeks`. Stored setting: **length in days** (`iteration_length_days`) on the project (default 7). UI draws Current + future bands from `pack`. See [velocity-and-planning.md](../velocity-planning/velocity-and-planning.md).
3. **Organisation slug and URL.** We assume a human-facing organisation name plus whatever URL the Technical Lead needs. Fork: public slug uniqueness vs per-user uniqueness.
4. **Project timezone source.** Window end is midnight in the project timezone. Fork: store an explicit project timezone (recommended) vs infer from the creating owner and never show a setting in MVP.
5. **Icebox order vs backlog order.** Tracker Icebox is its own ordered list, not interleaved with Backlog. We assume that. Fork: one global rank including Icebox (filter by state) vs two ranks.
6. **Who creates projects after slice 0.** We assume only organisation owners. Fork: any Member of the org.
7. **CSV export by Members.** We assume **Owners only** for import and export. Fork: Members may export.
8. **Cycle-time clock on reject.** We assume first `started` → `accepted`, and reject does not reset the clock. Fork: clock from the *last* Restart.
9. **`planned` as a story state vs a panel flag (slice 30).** Tracker adds `planned` when Current is manually planned. Manual planning is a later slice. Fork: introduce `planned` as a story state vs keep `unstarted` and a panel flag only. `planned` (if introduced) is not an assignment to a window record. There is no `iterations` table and stories have no `iteration_id`.
10. **Workspaces shared or personal.** We assume personal (per user) in one organisation. Fork: shared team workspaces.

If Dan does not answer, the “we assume” side ships.

## Assumptions already baked (correct these if wrong)

### Tracker copy (locked)

- Four types, different machines, including Reject + Restart → started, stays in Current.
- Icebox is unscheduled, in the product from MVP, not a sequential stage.
- One ordered list; unstarted cannot be dragged above started.
- Auto-plan leaves Current short rather than overfill; Start may overflow Current.
- Initial velocity **10** is bootstrap only: pack uses `initial_velocity` as estimate-points that fit in a full window until at least one corpus Feature exists in a **completed** window (or `time == 0`). Stories accepted in the **open** window are not in the lookback.
- Lookback is the last number of completed windows set by `velocity_strategy` (default 3, setting 1–4).
- Velocity is live duration: `velocity = work / time` from completed Features (`started_at` → `accepted_at`). Incomplete stories pack by `predicted_duration(estimate)` from completed stories of the same estimate. Team strength % later.
- Any Member can accept. No accept ACL. Requester should. My Work surfaces Delivered.
- Owners max 5. Start assigns the clicker. Requester + owners auto-follow and cannot unfollow.
- Tasks unowned, unpointed; warn on Accept.
- Blockers: free-text + optional **story** link; auto-resolve on accept/delete of the linked story; warn on Accept.
- Epics = one purple label + independent epic order + progress. Not a parent ticket.
- Search language: `type:`, `state:`, `estimate:-1`, `owner:`, `label:`, `is:blocked`, `includedone:`.
- Release marker at the **end** of its stories; optional target date; blue/red vs the **computed window start** that contains the marker.
- Manual planning later, Current only. Future **bands** stay auto-planned.
- Default Linear 0,1,2,3. Fibonacci 0,1,2,3,5,8 and Powers of 2 (0,1,2,4,8) later. Custom later and revertible.
- Bugs/chores points toggle later and **reversible**.
- History is undo (latest state-changing activity can be undone).

### Auth and API tokens (locked)

- Humans: email + password **and** magic link in MVP. SSO later. Session cookie for the app.
- **User API token** (Bearer): mint / list / revoke only on your own user: `GET|POST|DELETE /api/v1/users/:id/tokens`. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. Role at or below their own on selected projects. Same Owner / Member / Viewer permissions as that role. No extra scopes. Activity is attributed to the **user** the token belongs to. No org-level mint path.
- Same REST API (`/api/v1/...`) for cookie and Bearer. Same handlers. Webhooks are a product feature for any client. GraphQL later (product-wide).
- Mentions: in-app + email in Phase 1. Slack later.

### Planning (locked)

- Length is **days** (default 7), stored on the project as `iteration_length_days`. There is no `iteration_length_weeks`. Not 1–4 weeks. Not a stored list of windows.
- `pack(ordered_stories, velocity, predicted_duration, iteration_length_days, now, timezone) → bands`. `predicted_duration` is a size → predicted-duration map; unused on cold start. Stories have no `iteration_id`.
- There is no `iterations` table. Stories have no `iteration_id`. Nothing is stored for velocity except story timestamps, estimate, and project settings. Planning is live duration + similar-size predicted duration.
- Accepted stays in Current until the current window ends (midnight at `starts_on + iteration_length_days`, project timezone). Done is a flat aged-accepted list.

### Scope

- CSV import/export Phase 3; always-create, all-or-nothing. Pivotal Tracker import later if at all.
- No custom fields ever.
- No Gantt, no workflow engine, no resource management.
- Icebox in MVP (own slice). Workspaces Phase 3.

### Constraints

- UK / AU / NZ spelling (`organisations`).
- Look locked to `docs/reference/frontend-design-guide.md`.
- Ports: API 8180, frontend 4273, Postgres 5433/5437.
- Core tables: `users`, `projects`, `project_memberships`, `stories`, `labels`, `story_labels`, `activities`. This is the schema. Planning is a calculation — no `iterations` table. Stories have no `iteration_id`. Projects store `iteration_length_days` only. There is no `iteration_length_weeks`. Activities use `user_id`.
- `stories.rank` is a string. Title max 500.
- Business rules in `api/internal/domain/<domain>`, not the database. Frontend is split by domain. See [domain-model.md](./domain-model.md).

### Tenancy and roles

- Signup creates a user + an Organisation + a first Project.
- Roles: Owner / Member / Viewer only, on the project.
- Viewers are read-only (no comments, no follows). A Viewer may create a token on their own user; it can only read.
- Organisation owners can enter every project in the organisation.

## Decisions locked 20–21 Aug 2026

Keyboard and signup (UI Designer proposals, Product Owner accepted):

- **A** = Add Story (Tracker). Not Accept.
- **Accept** = Enter on a focused `delivered` row (or the Accept button).
- **R** = Reject (opens reason, then confirm).
- **Restart** = Enter on a focused `rejected` row.
- **Username** is inferred from the email local-part at signup. No username field in slice 0.
- **Shift+H** is reserved. No project history panel in MVP.
- Estimate keys in MVP are **0 / 1 / 2 / 3** only. Fibonacci / custom keys wait for slice 28.
- **Finished** has no status colour in the design guide. White row + Deliver verb. Do not invent yellow.
- **O** opens the story sheet. Enter stays the primary verb (Start / Finish / Deliver / Accept / Restart). Esc closes. Click the title also opens.
- Epic visual: no new purple hex. Epic pill is Bloom-bordered. Product language “purple label” means epic-marked, not a fifth brand colour.
- Slice 30 manual Current: when auto-plan is off, drag Backlog → Current is legal, and **C** (unshifted) moves the focused unstarted Backlog story into Current. Icebox → Current is still illegal (Pull to Backlog or Start).
- Start from Icebox: estimate-then-Start from Icebox is allowed (Feature must be estimated). It is schedule + start in one verb. The story lands `started` in Current. Shared machine rule for any Member.
- **API tokens (locked).** Mint / list / revoke only on your own user: `GET|POST|DELETE /api/v1/users/:id/tokens`. No mint-for-another-user. Member cannot mint Owner. A Viewer may mint a token on their own user; it can only read. Cookie vs Bearer is the only auth difference. Same `/api/v1` handlers. Role at or below your own. No org-level mint path. No extra scopes.
- **Dynamic windows (locked).** No Iteration aggregate. There is no `iterations` table. Stories have no `iteration_id`. Length in days (`iteration_length_days` only; default 7). There is no `iteration_length_weeks`. Current is the head of the ranked list that fits remaining time at `velocity` (or `initial_velocity` points on cold start). Future bands are visual. Recompute when stories, estimates, start/accept, rank, or settings change. Planning is live duration + similar-size predicted duration.
