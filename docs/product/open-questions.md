# Open questions

Only **true forks** — things we cannot assume. Everything else is decided. If it is in the assumption list, implementers treat it as law until Dan strikes it.

## True forks

1. **Agent accept / reject.** Humans: any Member or Owner may accept or reject a Feature/Bug. Agents: this spec forbids Feature/Bug `accept` and `reject` (they deliver; a human accepts). Chore finish→accepted and Release finish→accepted **are** allowed for an agent with `stories:transition`, because those verbs *are* Finish. Fork: should a scoped agent also accept Features?
2. **Organisation slug and URL.** We assume a human-facing organisation name plus whatever URL the Technical Lead needs. Fork: public slug uniqueness vs per-user uniqueness.
3. **Username at signup.** `users.username` is `NOT NULL` in the existing schema. Fork: ask for a username in slice 0, or infer from email local-part and let them edit later.
4. **Project timezone source.** Rollover is midnight in the project timezone. Fork: store an explicit project TZ (recommended) vs infer from the creating owner and never show a setting in MVP.
5. **Icebox order vs backlog order.** Tracker Icebox is its own ordered list, not interleaved with Backlog. We assume that. Fork: one global rank including Icebox (filter by state) vs two ranks.
6. **Who creates projects after slice 0.** We assume only organisation owners. Fork: any Member of the org.
7. **CSV export by Members.** We assume **Owners only** for import and export. Fork: Members may export.
8. **Cycle-time clock on reject.** We assume first `started` → `accepted`, and reject does not reset the clock. Fork: clock from the *last* Restart.
9. **`planned` state.** Tracker adds `planned` when Current is manually planned. Manual planning is a later slice. Fork: introduce `planned` in that slice (copy Tracker) vs keep `unstarted` and a panel flag only.
10. **Workspaces shared or personal.** We assume personal (per account) in one organisation. Fork: shared team workspaces.

If Dan does not answer, the “we assume” side ships.

## Assumptions already baked (correct these if wrong)

### Tracker copy (locked)

- Four types, different machines, including Reject + Restart → started, stays in Current.
- Icebox is unscheduled, in the product from MVP, not a sequential stage.
- One ordered list; unstarted cannot be dragged above started.
- Auto-plan leaves Current short rather than overfill; Start may overflow Current.
- Initial velocity **10** until one iteration completes.
- Velocity strategy last 3 completed iterations, setting 1–4.
- Velocity formula in MVP: accepted **Feature** points only. Team strength % later.
- Any Member can accept. No accept ACL. Requester should. My Work surfaces Delivered.
- Owners max 5. Start assigns the clicker. Requester + owners auto-follow and cannot unfollow.
- Tasks unowned, unpointed; warn on Accept.
- Blockers: free-text + optional **story** link; auto-resolve on accept/delete of the linked story; warn on Accept.
- Epics = one purple label + independent epic order + progress. Not a parent ticket.
- Search language: `type:`, `state:`, `estimate:-1`, `owner:`, `label:`, `is:blocked`, `includedone:`.
- Release marker at the **end** of its stories; optional target date; blue/red vs iteration **start**.
- Manual planning later, Current only.
- Default Linear 0,1,2,3. Fibonacci 0,1,2,3,5,8 and Powers of 2 (0,1,2,4,8) later. Custom later and revertible.
- Bugs/chores points toggle later and **reversible**.
- History is undo (latest state-changing activity can be undone).

### Auth and agents

- Humans: email + password **and** magic link in MVP. SSO later.
- Agents: scoped API tokens, first-class actors, no impersonation.
- Agent slice: REST + webhooks. GraphQL later.
- Mentions: in-app + email in Phase 1. Slack later.

### Scope

- CSV import/export Phase 3; always-create, all-or-nothing. Pivotal Tracker import later if at all.
- No custom fields ever.
- No Gantt, no workflow engine, no resource management.
- Icebox in MVP (own slice). Workspaces Phase 3.

### Existing Flower ground

- UK / AU / NZ spelling (`organisations`).
- Look locked to `docs/reference/frontend-design-guide.md`.
- Ports: API 8180, frontend 4273, Postgres 5433/5437.
- Eight core tables stay; we add, we do not redesign.
- `stories.rank` stays a string. Title max 500.
- Business rules in the Go API, not the database.

### Tenancy and roles

- Signup creates an Account + an Organisation + a first Project.
- Roles: Owner / Member / Viewer only, on the project.
- Viewers are read-only (no comments, no follows, no tokens).
- Organisation owners can enter every project in the organisation.

## Decisions locked 20 Aug 2026

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
- Start from Icebox: estimate-then-Start from Icebox is allowed (Feature must be estimated). It is schedule + start in one verb. The story lands `started` in Current.

