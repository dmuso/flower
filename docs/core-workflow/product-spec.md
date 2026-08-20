# Flower product spec

Change name: `flower`

Eventual path: `docs/core-workflow/product-spec.md` (see [LANDING.md](./LANDING.md)).

This spec is the product. Companion files own exhaustive rules for planning and tenancy. Do not implement until the Reviewer clears the spec set.

Spelling: UK / AU / NZ (`organisation`). Look: locked Flower identity, not a new one — UI Designer starts from `docs/reference/frontend-design-guide.md`.

## Problem

Small product teams still cannot see an honest plan.

Who feels it:

- **Product owners** with a backlog in one tool, dates in a spreadsheet, and status in standup. They cannot answer “when does this land” without planning theater.
- **Developers** asked to update tickets that do not match how work moves, and who lose the top of the list to unestimated Features or work that was never accepted.
- **Automation** (scripts, CI, Grove) expected to keep the board current by scraping, impersonating a person, or guessing which status string means done.

What is missing:

- Classic Tracker’s model, implemented faithfully: one ordered list, Icebox as an unscheduled holding pen, type-specific state machines, velocity from accepted Feature points.
- An honest empty/new project: initial velocity 10, auto-fill, no fake milestone picker.
- The same story API the frontend uses, so a script can move a story with an **API token** rather than a person’s password.

What better looks like:

A team signs up, makes an organisation and a project, and sees Icebox / Backlog / Current / Done. They type Features into Icebox, pull them into the ranked list, estimate 0/1/2/3, start only estimated Features, finish, deliver, and accept (any Member may; the requester should). Rejected work Restarts to started and stays in Current. Flower packs the ranked list into Current and future visual bands from velocity. A second person sees the move without refresh. A Member or Owner mints a **user API token** (`/api/v1/users/:id/tokens`) and a client (CI, a script, Grove) calls the same story API as the frontend.

Keep `docs/product/overview.md` as the one-page intent and make it match this spec set (see [LANDING.md](./LANDING.md)).

## Personas

### Product owner — Maya

Maya runs a 6–12 person product. She is the requester on most Features.

She needs: the next story at the top; accept or reject herself (or any Member if she is away); a projected date from velocity; a Release marker at the **end** of a milestone with an optional target date that turns red when the plan slips; Icebox for not-now work.

She does not need: swimlanes, custom fields, a second status, Gantt.

Daily: reorder, pull from Icebox, skim Delivered (My Work), accept or reject, glance at the Release colour.

### Developer — Luis

Luis implements the top of Current. He wants the board to keep up.

He needs: keyboard; estimate then Start (he becomes an owner); tasks; a blocker that can point at another story; labels; Markdown; image paste; deliver and move on. He *may* accept (any Member can) but should not make a habit of accepting his own work when Maya is the requester.

He does not need: time tracking, per-person capacity, an “In Review” column.

Daily: Start the next estimated Feature, check tasks, paste a screenshot, Deliver, next.

### Automation — Grove

Grove is a **client** (CI, a script, or a coding job) that keeps the board current. It authenticates with a **user API token** minted for a Member or Owner at `/api/v1/users/:id/tokens`. It is not a distinct user type.

It needs: a Bearer token (not someone’s password); the **same** `/api/v1/...` story API the app uses; `unestimated` if it tries to Start a Feature without points. Activity is attributed to the **user** the token belongs to.

Permissions are Owner / Member / Viewer — the whole model. A token minted for a Member can do everything Luis can, including accept and reject. A Viewer token can only read. The token does not add scopes.

Same machines as anyone — do not guess an estimate, invent a type or state, reorder without an explicit position, or treat Delivered as Accepted.

Daily: Start when the branch opens, comment the PR URL, Finish when the work is in, Deliver when CI is green. Same verbs as the board.

## Out of scope

- Gantt, roadmap timelines, dependency arrows as planning UI.
- Custom fields. None. Ever. Labels, types, and epics (epic-marked labels) are the vocabulary.
- Workflow engines, custom statuses, column designers.
- Resource management, individual capacity, PTO, “who is overloaded.”
- Time tracking. Per-person velocity.
- SSO / SAML (later). Slack (later). GraphQL (later unless actually cheaper than REST for a named slice).
- Pivotal Tracker import (later if at all). CSV is Phase 3.
- Public projects, anonymous boards.
- Native mobile apps.
- AI inside Flower (auto-estimate, summarise). Clients *use* Flower via the same API.
- A wiki. Images on a story only.
- Multiple backlogs per project, portfolio rollups across organisations, SAFe.
- Team strength % (later). Manual planning of Current (later, Phase 3).
- A new visual identity.

## Assumptions (law; see [open-questions.md](./open-questions.md))

See also [tracker-brief.md](./tracker-brief.md). Short form:

- Email + password and magic link for humans in MVP; SSO later. Session cookie for the app.
- **User API token** (Bearer): list / create / revoke at `/api/v1/users/:id/tokens`. Each user manages their own. Role at or below their own on the selected projects. Member cannot mint Owner. Same Owner / Member / Viewer permissions as that role. No extra scopes. Activity is attributed to the **user** the token belongs to. Not org-level.
- Roles: Owner / Member / Viewer only. **Any Member or Owner can accept.** Viewers read-only. Requester *should* accept; My Work surfaces Delivered. No accept ACL in MVP. History is undo.
- Feature/Bug reject → `rejected`; Restart → `started`; stays in Current.
- Default points Linear 0/1/2/3. Fibonacci 0/1/2/3/5/8 and Powers of 2 later. Custom later and revertible.
- Icebox is **MVP**, own slice. Not Phase 3.
- Same REST API (`/api/v1/...`) for cookie and Bearer. Same handlers. Webhooks are a product feature for any Member/Owner client.
- Mentions: in-app + email in Phase 1.
- CSV Phase 3. No custom fields ever.
- Initial velocity 10. Feature points only.
- Owners max 5. Start assigns the clicker.
- Ports 8180 / 4273 / 5433 / 5437. Schema must match this model.

## Tech stack (constraints)

PostgreSQL 17. Go + Gin API on **8180**. SolidJS on Bun on **4273**. Monorepo. Multitenant from day one. Domain logic lives in `api/internal/domain/<domain>`. Frontend is split by domain (not one `api.ts`). Technical Lead owns the tree — see [domain-model.md](./domain-model.md).

## Types and state machines

MVP ships **Feature** plus Icebox/Backlog/Current/Done. Phase 2 adds Bug, Chore, Release. Machines are specified now so later slices do not invent a second workflow.

### Feature and Bug

`unscheduled` → `unstarted` → `started` → `finished` → `delivered` → `accepted`  
`delivered` → `rejected`  
`rejected` + **Restart** → `started` (stays in Current)

| Verb | From | To | Who |
| --- | --- | --- | --- |
| schedule (Icebox → list) | unscheduled | unstarted | Owner, Member |
| icebox (list → Icebox) | unstarted only | unscheduled | Owner, Member |
| estimate | any except accepted | same; sets points | Owner, Member |
| start | unstarted **or** unscheduled (Icebox) | started | Owner, Member. Feature must already be estimated. Icebox start is schedule + start; story lands in Current. |
| finish | started | finished | Owner, Member |
| deliver | finished | delivered | Owner, Member |
| accept | delivered | accepted | Owner or Member (any). Not Viewer. |
| reject | delivered | rejected | Owner or Member (any). Not Viewer. Reason required (comment). |
| restart | rejected | started | Owner, Member |
| undo | last state-changing activity | previous state | Owner, Member |

**Start a Feature** requires an estimate (`0` allowed). Start a Bug does not. Start assigns the clicker as a story owner (max 5; if already 5 and clicker is not among them, Start still happens and they are **not** added — show that).

### Chore

`unscheduled` → `unstarted` → `started` → `accepted`

Finish **is** accept. No delivered, no reject. Unestimated by default. Does not count toward velocity.

### Release

Marker, not work. `unscheduled` in Icebox. Created in Backlog or dragged there → auto-`started`. Finish → `accepted`.

No estimate. Optional **target date** (a comparison, not a plan override). Place at the **end** of the milestone’s stories (all work for that release sits **above** the marker). Colour: blue if the **computed window** that contains the marker starts on or before the target; red if that window **starts after** the target. No colour if there is no target date.

## Story owners, requester, follow

- **Requester:** creating Member or Owner. May be changed to another Member/Owner. Never a Viewer as requester unless we later allow it — **requester must be able to accept**, so Member or Owner.
- **Owners:** 0–5 project Members or Owners. Start assigns the clicker.
- **Follow:** requester and owners auto-follow and **cannot unfollow**. Others may follow (Member/Owner). Viewers do not follow.

## Schema

Schema must match this model.

Core: `users`, `projects`, `project_memberships`, `stories` (incl. `requester_id`, `estimate`, `rank`, `state`, `story_type`, `accepted_at`), `labels`, `story_labels`, `activities` (`user_id` = the user who did the thing).

Planning is a calculation. There is no `iterations` table. Length is **`iteration_length_days`** on the project (default 7). Accepted-points history lives in `velocity_samples`. Product language is **user**.

Slices add tables they need: organisations, story owners, comments, tasks, attachments, epics, blockers, followers, tokens.

## Ordered vertical slices

Each slice is a thin end-to-end user-visible outcome. Independently acceptable. Do not start slice n+1 until slice n is agreed, unless a human parallelises.

---

# Phase 0 — MVP

A team can sign up, share a project, Icebox a Feature, rank it, estimate, run the Feature machine, reorder, and see an honest plan (initial velocity 10).

### Slice 0 — Owner signs up, creates an organisation and first project, sees an empty board

**Why slice 0:** later slices need an account, a tenant, and a board. Still demoable.

**Why independently acceptable:** a stranger becomes Owner of an organisation that contains one project, signs out, signs in, and sees the same empty four columns.

Acceptance criteria:

- Given I have no account, when I sign up with email + password and verify email, then I name an **organisation** and a first **project** before I see a board. Username is inferred from the email local-part. No username field.
- Given I have no account, when I open a magic link to a new email, then an account is created and I am on the same organisation + project flow.
- Given I have an account, when I sign in with password **or** magic link, then I land on my last project.
- Given I created organisation `Acme` and project `Trail`, when the board loads, then I see four columns — **Icebox, Backlog, Current, Done** — each empty, each with one next action. No fake stories. No fake dates. Current may show the computed window end and **velocity 10** (initial). Backlog may show future **band** headers once stories exist; on an empty board, empty copy is enough. No stored iteration list.
- Columns follow `docs/reference/frontend-design-guide.md` (full-height, paper, bloom current highlight). Fail if a new palette appears.
- Unverified password signup cannot create an organisation.
- Creator is organisation owner and project owner.
- Reload stays in `Acme` / `Trail`, not another tenant.

### Slice 1 — Owner invites a member with a role

**Why independently acceptable:** two humans share one project, even if the board is empty.

Acceptance criteria:

- Given I am a project Owner, when I invite `alex@example.com` as `member`, then they get an email with a single-use link and I see a pending invite.
- Given they have no account, when they accept, then they complete signup and land on the project as Member.
- Given they have an account, when they accept and sign in, then the project is in their list.
- Roles are exactly `owner`, `member`, `viewer`.
- Given I invite as `viewer`, when they open the board, then they can see it and **cannot** create a story, invite, or change settings.
- Email already a member → visible error, no second membership.
- Revoke kills the link. Expiry 14 days. Resend invalidates the old link.
- Members and Viewers cannot invite.

### Slice 2 — Member creates a feature in the Icebox

**Why independently acceptable:** not-now work has a place that does not pollute the plan.

Acceptance criteria:

- Given I am a Member or Owner, when I create a Feature from Icebox (or the default add, which is Icebox), then it appears in Icebox as type `feature`, state `unscheduled`, no estimate, requester = me, not in the ranked list (no band, no projected date).
- Empty title is rejected. Title max **500** (existing column).
- Viewer cannot create (UI and API).
- Icebox is ordered (own order). New story goes to the **top** of Icebox (most recently captured). 
- Current / Backlog / Done unchanged. Velocity math ignores Icebox.

### Slice 3 — Member pulls a story from Icebox into the Backlog

**Why independently acceptable:** scheduling is a deliberate pull, not a pipeline step. This is the Icebox MVP slice.

Acceptance criteria:

- Given an `unscheduled` Feature in Icebox, when I move it to Backlog (drag or action), then it is `unstarted`, it leaves Icebox, and it sits at the **bottom** of the ranked list (does not jump the queue).
- Given an `unstarted` Feature in Backlog, when I move it to Icebox, then it is `unscheduled`, drops out of `pack`, and loses any projected date.
- I cannot icebox `started`, `finished`, `delivered`, `rejected`, or `accepted` stories.
- Viewer cannot move.
- After this slice, an empty Icebox shows the empty state again.

### Slice 4 — Member estimates a feature and starts it

**Why independently acceptable:** work enters Current the honest way.

Acceptance criteria:

- Given an `unstarted` Feature, when I set estimate to `0`, `1`, `2`, or `3`, then it shows that estimate and stays `unstarted` until I Start.
- Given an unestimated Feature, when I Start (UI or API), then Start is rejected, state stays `unstarted`, message says it needs an estimate. Fail the slice if Start succeeds.
- Given an estimated Feature, when I Start it, then it is `started`, it appears in Current (even if that overflows velocity), I am a story owner, and follow is recorded; Follow UI is slice 12.
- Clearing estimate is allowed only while `unstarted`. Started Features cannot become unestimated.
- Viewer cannot estimate or start.
- Starting from Icebox: estimate first (or estimate-and-start from Icebox). The story becomes `started` in Current, not `unstarted` in Backlog.

### Slice 5 — Member finishes and delivers; any member accepts

**Why independently acceptable:** a Feature can go from started to accepted. Accepted stays in **Current** until rollover.

Acceptance criteria:

- Started → Finish → `finished`, stays in Current.
- Finished → Deliver → `delivered`, stays in Current. Requester sees it in My Work once that slice exists; until then, Delivered is visible in Current.
- Given `delivered`, when **any** Member or Owner accepts, then state is `accepted`, `accepted_at` is set, and the story **remains in Current** (not Done). Done is empty until that accept ages past the current window.
- Viewer cannot finish / deliver / accept.
- Illegal verbs fail and do not change state (finish on unstarted, accept on finished, …).
- Incomplete tasks / open blockers do not exist yet; when they do, Accept **warns** and still proceeds.

### Slice 6 — Member rejects a delivered story and restarts it

**Why independently acceptable:** bad work returns to the builder without a hallway chat.

Acceptance criteria:

- Given `delivered`, when a Member or Owner rejects with a reason, then state is `rejected`, the reason is visible, the story stays in Current (not Backlog, not Icebox).
- Empty reason → reject does not happen.
- Viewer cannot reject. A Member can.
- Given `rejected`, when someone clicks **Restart**, then state is `started`, still Current, reason remains in activity.
- After Restart they can Finish and Deliver again. A new Accept is required.
- Fail the slice if reject jumps straight to `started` with no Restart, or if `rejected` is treated as terminal like `accepted`.

### Slice 7 — Member reorders the ranked list

**Why independently acceptable:** priority is real.

Acceptance criteria:

- Drag and keyboard reorder persist after reload for Members and Owners.
- **Unstarted cannot be dragged above started** (nor above finished / delivered / rejected). The drop is rejected and the row snaps back. Copy explains why.
- Accepted-this-window stories live in Current; they do not go back to Backlog by drag.
- Dragging Backlog → Current does **not** Start. Auto-plan membership is not a drop-to-start. Start is a verb. (Manual planning later is the escape hatch.)
- Icebox order is independent and reorderable. Drag Icebox → Backlog is slice 3 (schedule), not a silent Start.
- Viewers cannot reorder.
- UI Designer sets the keyboard chord; a keyboard reorder path must exist.

### Slice 8 — Team sees Current and future bands auto-fill from velocity

**Why independently acceptable:** the Tracker moment, computed live. Rules: [velocity-and-planning.md](./velocity-and-planning.md).

Acceptance criteria:

- Brand-new project: velocity **10** (initial). **Current is the head of the ranked list that fits this window’s velocity.** Auto-fill estimated unstarted Features from the top up to 10 points, **leaving Current short** rather than putting a story that would exceed 10 into Current.
- Future bands in Backlog are **visual** only — the same `pack` function, not stored rows, not story assignments.
- Starting a Backlog or Icebox Feature still jumps it to Current and **may** make Current > 10.
- Window **end** (chosen rule): midnight at `starts_on + L days` in the project timezone, where `starts_on` is the configured start weekday on or before project-created, and `L` is length in days. After the first window completes, velocity becomes accepted **Feature** points of that window (and then the rolling last 3, setting 1–4, default 3).
- A story is never split. If the next Feature is 3 and remaining is 1, it stays in the next Backlog **band**.
- Reorder, estimate, accept, start, icebox, length change, or window end → plan already recomputed. No Recalculate button.
- Accepted this window: still in Current until the current window ends. After the window ends: those accepted stories are in Done as a **flat** list (newest accepted first), not grouped by a window row.
- Length default **7 days**. Owner may set a positive number of days. Changing length replans. Not 1–4 weeks. Not a stored iteration list.
- Bugs/chores/releases are not required for this slice (Features only). When they exist, they follow the velocity doc.
- An estimated Feature with cost > V auto-fills only into a band with 0 Feature-points and marks that band over velocity; it never sits unpacked.
- Fail the slice if the planner persists window rows or assigns a story to a window, or if Current is labelled “iteration 3”.

Phase 0 is not done without slice 8.

---

# Phase 1 — Story depth and team UX

Still Features only unless a later Phase 2 slice has landed.

### Slice 9 — Member adds tasks on a story

**Why independently acceptable:** a story can be checked off without becoming five stories.

Acceptance criteria:

- Tasks are unowned, unpointed checklists. Text required. Toggle persists. Story shows `done/total`.
- Completing all tasks does not Finish or Accept.
- Tasks do not affect velocity.
- Reorder tasks persists.
- Viewer can see, cannot toggle.
- Accept with incomplete tasks: **warn**, do not hard-block (once Accept UI can show a confirm; if Accept is a single click, show a non-blocking warning and still accept).

### Slice 10 — Member marks a story blocked

**Why independently acceptable:** stuck work is visible on the row.

Acceptance criteria:

- Blocker = required free-text + optional **link to another story in the same project** (picker, not a free URL).
- Badge on the row. Multiple blockers allowed.
- When the linked story is **accepted** or **deleted**, that blocker auto-resolves.
- Blockers do not prevent Start / Finish / Deliver / Accept.
- Accept with open blockers: **warn**, do not hard-block.
- Viewer can see, cannot add or clear.

### Slice 11 — Member labels a story

**Why independently acceptable:** filter the board to `api`.

Acceptance criteria:

- Label: lowercase `[a-z0-9-]+`, 1–32 chars (or existing `labels.name` VARCHAR(100) — do not exceed 100). Visible on the row.
- Filter by label hides non-matching stories in all columns. Does not change rank or state.
- Creating a label on a story creates it for the project. Owner can delete an unused project label.
- Viewer can filter, cannot add.
- Epic-as-purple-label is Phase 2; this slice is ordinary labels. `labels` / `story_labels` are the tables.

### Slice 12 — Member assigns owners, requester, and followers

Why independently acceptable: the row shows who owns the story; the requester can be handed off.

Acceptance criteria:
- Given a story, when I add owners, then 0–5 project Members/Owners are accepted; the sixth is rejected with that limit. Start already assigned the clicker (slice 4).
- Requester change: another Member or Owner on the project. Never a Viewer.
- Requester and owners auto-follow and cannot unfollow. Other Members/Owners may follow/unfollow.
- Viewer cannot assign, change requester, or follow.
- This slice adds the owners table. Do not overload requester_id.

### Slice 12b — Member @mentions a teammate and they are notified

Why independently acceptable: an @ reaches a human without Slack.

Acceptance criteria:
- `@` in a comment or description that resolves to a project Member, Owner, or Viewer → in-app + email. Unresolved `@` is plain text. Slack out of scope.
- Email + in-app also for: assigned as owner, delivered (to requester), rejected (to owners).
- Viewer can be mentioned and still cannot follow, comment, or assign.

### Slice 13 — Member writes Markdown description and comments

**Why independently acceptable:** the brief and the conversation live on the story.

Acceptance criteria:

- Description and comments: headings h1–h3, bold, italic, strikethrough, lists, fenced code, inline code, block quotes, links. No raw HTML. `javascript:` URLs rejected.
- Comments have author and timestamp. Empty comments rejected. Description may be empty.
- Edited mark. Deleted comments tombstoned, not silently erased.
- Viewer can read, cannot write.
- Comments table does not exist yet; this slice adds it. `stories.description` already exists.

### Slice 14 — Member attaches images

**Why independently acceptable:** paste a screenshot.

Acceptance criteria:

- `png`, `jpg`, `jpeg`, `gif`, `webp`, ≤ 10 MB. Max 20 images per story. Other types rejected.
- Clipboard paste into the story creates an attachment and embeds it.
- Remote hotlinked images that are not attachments do not render.
- Viewer can see, cannot upload.
- Delete → missing-image state on embeds.

### Slice 15 — Member reads activity and undoes the last state change

**Why independently acceptable:** anyone can see who rejected, and a mistaken Accept is reversible.

Acceptance criteria:

- Activity list on the story: created, estimated, scheduled/iceboxed, started, finished, delivered, accepted, rejected (reason), restarted, requester/owners/title changed, blocker, attachment, comment, label. **User**, time, from → to.
- Reorders do not spam activity.
- **Undo** on the latest state-changing activity restores the previous state (Accept → Delivered, Reject → Delivered, Start → Unstarted, Icebox Start → Unscheduled, Restart → Rejected, Finish → Started, Deliver → Finished). Undo is itself an activity.
- You cannot undo a non-state event (a comment) with this control.
- Viewer can read activity, cannot undo.
- `activities` records `kind`, `summary`, and `user_id`. Product language is **user**.

### Slice 16 — Member runs the board from the keyboard

**Why independently acceptable:** a morning without the mouse.

Acceptance criteria:

- Keyboard can: move focus, **O** open / Esc close, create (**A**), estimate 0–3, start (**S**), finish (**F**), deliver (**D**), accept (**Enter** on delivered), reject (**R**, reason then confirm), restart (**Enter** on rejected), reorder one slot, icebox/schedule, `/` focuses search or filter.
- Escape closes the open story or dialog and returns focus to the board.
- Viewer: navigate and open; mutating chords announce no permission.
- **Stub:** full table is the UI Designer’s in `ui.md`. Fail if any action above has no keyboard path.

### Slice 17 — Second member sees updates in real time

**Why independently acceptable:** grooming without refresh.

Acceptance criteria:

- Two sessions, same project: create, move Icebox↔Backlog, reorder, estimate, state change, comment, task toggle appear on the other board within 2 seconds without manual refresh.
- The other user does not lose focus or a comment draft unless that story/comment was deleted.
- No presence avatars or live cursors in this slice.
- Dropped socket: visible stale state + reconnect. Refresh heals.
- Viewers receive the same reads.

---

# Phase 2 — Types, dates, search, API

### Slice 18 — Member files a bug, chore, or release

**Why independently acceptable:** not everything is a Feature, and the machines differ.

Acceptance criteria:

- Create type Feature (default), Bug, Chore, or Release.
- Feature: estimate required to Start; points count when accepted.
- Bug: same machine as Feature; unestimated by default; may Start without points; does **not** count toward velocity.
- Chore: unscheduled → unstarted → started → accepted; Finish accepts; no reject; unestimated; 0 velocity.
- Release: marker; auto-started in Backlog; Finish accepts; optional target date; sits at the end of its stories; blue/red vs the **computed window** that contains the marker (velocity doc).
- Only `unstarted` or `unscheduled` stories may change type. Changing to Release clears estimate.
- “Bugs and chores may be given points” is **not** in this slice (Phase 3, reversible).

### Slice 19 — Member groups stories into an epic

**Why independently acceptable:** a theme without a parent ticket or a custom field.

Acceptance criteria:

- An epic **is** a project label marked epic (Bloom-bordered pill; no new purple hex). A story has at most one epic label (ordinary labels still allowed).
- Epics have an **independent order** (epic panel / list), not the story rank.
- Progress: accepted Feature points / estimated Feature points of stories with that epic label. Unestimated add 0 to the denominator. Bugs/chores do not add points (unless the later toggle is on).
- Filter by epic. Removing the epic label does not delete the story.
- Deleting the epic label unassigns it; stories remain.
- Viewer can view and filter. Member can assign. Owner can delete the epic label.

### Slice 20 — Member sees a projected delivery date

**Why independently acceptable:** Maya can answer “when.”

Rules: [velocity-and-planning.md](./velocity-and-planning.md).

Acceptance criteria:

- A story that `pack` placed in a computed band shows that band’s **end date** as the projection (last calendar day of the window).
- Icebox and (if ever unpacked) stories with no band: no date — “Not scheduled.”
- Release date = end date of the **computed window that contains the marker** (the marker sits at the end of its stories, so that window is when the last story above it packed). Target date optional; marker is blue or red versus that window’s **start**. **No date picker that sets the plan.** Target date is the only date the user types, and it does not move stories.
- Reorder / estimate / accept / velocity / length change updates dates and colours live.
- Fail the slice if a projection is read from a stored window row or a story-to-window assignment.

### Slice 21 — Member reads the velocity chart and burn-up

**Why independently acceptable:** burn vs added scope is visible.

Acceptance criteria:

- Velocity: bar per **completed window** (accepted Feature points); line at current velocity (initial 10 until N≥1). Current window is not a completed bar.
- Burn-up: cumulative accepted Feature points vs scoped Feature points (snapshot at window end + live “now”). Releases/bugs/chores add 0 to scope unless the later toggle is on.
- Empty charts have empty states, no fake history.
- Viewers can see charts.

### Slice 22 — Member searches with Tracker-like language

**Why independently acceptable:** a known story is findable.

Acceptance criteria:

- Full-text: title, description, comments.
- Operators, combinable, Tracker-like: `type:feature|bug|chore|release`, `state:`, `estimate:-1` (unestimated), `estimate:0`…`estimate:3`, `owner:name`, `label:api`, `is:blocked`, `includedone:` (Done is excluded unless this is present).
- Results do not reorder the board. Empty result: one next action (clear).
- Scoped to the current project in this slice.

### Slice 23 — Member or Owner mints a user API token and calls the same story API

**Why independently acceptable:** Grove (or CI, or a script) is a **client with a user token**. It delivers when CI is green, on the same endpoints Luis uses, with a Bearer token.

Acceptance criteria:

- List / create / revoke on the **user**: `GET|POST|DELETE /api/v1/users/:id/tokens`. Each user manages their own tokens (account settings, not an organisation or project screen). Owner or Member creates a generic **API token** with a role at or below their own on projects they belong to. Member cannot mint Owner. Secret shown once. The token does not add scopes beyond Owner / Member / Viewer.
- Bearer token. Humans can also use a session cookie. Same `/api/v1` handlers, same request bodies, same errors. If the frontend mutates with `POST /api/v1/stories/:id/transitions`, the token uses that same shared machine. There is no org-level mint path. There is no `unstart` verb.
- Viewer-bound token cannot mutate (start / finish / deliver / accept / reject / reorder → `forbidden`). A Viewer may still create a token on their own user; it can only read.
- Member-bound token **can** accept a delivered Feature (and reject with a reason). Same as a Member with a cookie.
- Start on unestimated Feature → `unestimated`, no state change (shared machine rule).
- Illegal transition → `invalid_transition` with `from` and `action`.
- Icebox Start of an estimated Feature lands `started` in Current — shared Member rule.
- Revoked token → `unauthorized`.
- Activity is attributed to the **user** the token belongs to. Token name may appear as `via {token name}`.
- Webhooks, if registered, are a product feature for any Member/Owner client. Events include story created/updated/reordered, story started/finished/delivered/accepted/rejected/restarted, comment created, window completed, membership changed. Signed outbound POST; delivery is not a command.

### Slice 24 — Member saves a search and opens My Work

**Why independently acceptable:** Monday morning has a landing place.

Saved search and My Work are two outcomes on one landing; acceptable together because both are read/filter, not two systems. Do not split unless easy.

Acceptance criteria:

- Save current search (text + operators) with a name, per user per project. Reopen applies the same filter. Delete does not delete stories.
- **My Work:** stories I own in `started|finished|delivered|rejected`, plus stories I request that are `delivered`. Empty: “Nothing on you.”
- Viewers do not get My Work. They may use search.

---

# Phase 3 — Multi-project and power tools

### Slice 25 — Member opens a named workspace of projects in one organisation

Workspace = a personal named set of projects in one organisation. Not a tenant. See [multitenancy.md](./multitenancy.md).
Why independently acceptable: Maya switches Trail and Checkout without leaving Acme, and cannot pull in another organisation.

Acceptance criteria:
- Given I can open Trail and Checkout in Acme, when I create workspace `Core` and add both, then opening `Core` shows those projects. Default is all projects I can open.
- Given a project in org B, when I try to add it to `Core`, then it is not in the picker and is not added.
- Permissions stay per project: Viewer on Trail still cannot mutate Trail from the workspace.
- Given `Core`, when I delete the workspace, then Trail and Checkout still exist.
- Cannot add another organisation’s project. Deleting a workspace does not delete projects.

### Slice 26 — Owner exports and imports CSV

**Why independently acceptable:** Maya can snapshot or seed a project from a spreadsheet without a Pivotal importer.

- Owner export: id, title, type, state, estimate, labels, description, requester email, owner emails, epic name, blocker text.
- Import **creates** only. No upsert. All-or-nothing. Not a Pivotal dump.
- Members and Viewers cannot import. Viewers cannot export. Members cannot export (assumption).
- Attachments and comments are not in the CSV.

### Slice 27 — Member reads cycle time

**Why independently acceptable:** the team can see how long work takes, without a per-person scoreboard.

- First `start` → `accepted`. Reject does not reset the clock.
- p50 / p75 / p95 over an accepted-at range, filterable by type.
- No per-person leaderboard.

### Slice 28 — Owner changes the point scale (including revertible custom)

**Why independently acceptable:** the team can change how points work on one settings page — the bugs/chores points toggle is on that same page, not a second epic.

- Settings: `linear` (0,1,2,3, default), `fibonacci` (0,1,2,3,5,8), `powers_of_two` (0,1,2,4,8), `custom` (Owner-defined list, must include a way to **revert** to linear without converting history).
- Switching does **not** convert existing estimates. Illegal-on-new-scale values remain on old stories until edited.
- Members cannot change the scale.
- This is the slice that also ships the reversible **“bugs and chores may be given points”** project toggle (default off). Turning it off is allowed (Tracker’s irreversibility was hated).

### Slice 29 — Member opens side-by-side panels

**Why independently acceptable:** Luis can keep the board and a story (or another project) open without a second window.

- Two panes: board, story, My Work, search, Icebox, or another project in a workspace.
- Reload restores the last split. Both panes live-update.

### Slice 30 — Owner turns off automatic planning for Current

**Why independently acceptable:** Owner can stop auto-fill for Current without changing how future bands plan.

- Escape hatch, Current **only**. Future bands stay auto-planned.
- Manual Current: only in-progress, accepted-this-window, and stories explicitly moved there (drag Backlog → Current, or **C** on a focused unstarted Backlog story). Icebox → Current is still illegal. No velocity fill into Current.
- Default remains automatic. A control to restore auto-plan replans Current from velocity.
- Not MVP.

---

## Required agents

| Agent | Why they are in |
| --- | --- |
| **Reviewer** | Always. Spec is not ready until they agree it is clear, vertical, and testable. |
| **UI Designer** | Entire product is new UI. Manifesto + **locked** frontend guide. They write `ui.md` (including shortcuts). They do not invent bloom/stem/paper. |
| **Technical Lead** | Stack, ports, and the core schema are constrained; approach is not. They write `technical-approach.md` and own `api/internal/domain/<domain>` plus the frontend domain split. Tenancy, machines, live updates, **user** API tokens (`/api/v1/users/:id/tokens`), `pack` as a pure function. |
| **Developer** | One slice at a time. Documented PR. |
| **QA** | Tests each slice from AC plus the companion rule docs. Can fail a slice without a meeting. |

A human merges. Do not implement before the Reviewer clears the spec.

## References

- This brief and the locked Tracker copy-exactly list ([tracker-brief.md](./tracker-brief.md)).
- Domain map (Technical Lead): [domain-model.md](./domain-model.md).
- Classic Pivotal Tracker help: story states, types, velocity, backlog-to-current, icebox, releases, automatic vs manual planning.
- [github.com/dmuso/flower](https://github.com/dmuso/flower) — `docs/product/overview.md`, `docs/reference/frontend-design-guide.md`. Schema must match this model.
- House rules: vertical slices, write-spec, review-spec, delivery workflow, UI manifesto (Pivotal / old Stripe / old Trello; refuse Jira and Monday).
