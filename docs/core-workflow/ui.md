# Flower UI

Change name: `flower`

Eventual path: `docs/core-workflow/ui.md` (see `docs/product/LANDING.md`).

This is the implementation UI spec. Future-state only. Look is locked in [`docs/reference/frontend-design-guide.md`](../reference/frontend-design-guide.md). Product rules are locked in [`docs/core-workflow/product-spec.md`](./product-spec.md), [`docs/product/tracker-brief.md`](../product/tracker-brief.md), and [`docs/velocity-planning/velocity-and-planning.md`](../velocity-planning/velocity-and-planning.md).

Spelling: UK / AU / NZ (`organisation`).

Do not implement from this file until the Reviewer has cleared the spec set.

---

## 1. What this file is

A developer should be able to build every screen in slices 0–30 without choosing a colour, a typeface, a column name, or a verb label.

If a control is not in the frontend design guide, do not invent a second look for it. Reuse **Primary**, **Secondary**, **Destructive**, **Story card**, **Board column**, **Column header**, **Text input**, **Lucide** (stroke 2px). If a product rule is not in the product spec / tracker brief / velocity doc, do not invent it — add an open question.

---

## 2. Locked look

Cite: `docs/reference/frontend-design-guide.md`. Do not add a fifth brand colour, a second type pair, or a Kanban redesign.

### Colour (exact)

| Token | Hex | Use |
| --- | --- | --- |
| Bloom (primary) | `#C43B6E` | Primary buttons, current-window highlight, key links, selected story ring |
| Bloom hover | `#d24d7d` | Primary hover (`hover:bg-bloom/90` on the button) |
| Bloom active | `#a3325a` | Primary active |
| Stem (secondary) | `#2F7D4A` | Accepted / done, Secondary-as-success |
| Pollen (accent) | `#E8A317` | Started / in-progress, **estimate chips** |
| Ink | `#1C1917` | Text. Never pure black |
| Paper | `#FBF7F2` | Page and column background |
| Surface | `#FFFFFF` | Story cards, dialogs, settings cards |
| Border on paper | `#E8DFD6` | Column rules, header rules |
| Border on white | `#D9D0C7` | Card edges (`border-paper-300`) |

Status colours from the same guide — no extras:

| Status | Hex | Use |
| --- | --- | --- |
| Accepted / success | `#2F7D4A` | Accepted row meta, Done, success toast |
| Rejected / destructive | `#B42318` | Rejected meta, reject reason, late Release |
| Started / in progress | `#E8A317` | Started indicator (not body text — contrast) |
| Delivered | `#3A6EA5` | Delivered meta, on-track Release |

**Type is not a colour.** The guide does not give Feature / Bug / Chore / Release swatches. Do not invent them. Type = Lucide icon + the word if the row is open. State uses the table above.

Pollen `#E8A317` is a chip and a 6–8px indicator, never small body text on paper.

### Type (exact)

- **Headings:** Fraunces, semi-bold. Column titles, dialog titles, settings page titles, empty-state titles.
- **Body:** Inter Regular (400). Story titles, panel copy, descriptions, comments.
- **Meta / labels:** Inter Medium (500), often uppercase, muted (`text-ink-700`). Column headers, band headers, estimate chip text, owner initials.

### Components (exact classes)

**Primary button** — the one next action when it is a button:

`rounded-full bg-bloom px-4 py-2 text-sm font-medium text-white hover:bg-bloom/90`

**Secondary button** — visible, quieter:

`rounded-full border border-paper-300 px-4 py-2 text-sm font-medium text-ink hover:bg-paper-50`

**Destructive button** — Reject, Delete, Revoke:

`rounded-full bg-red-700 px-4 py-2 text-sm font-medium text-white hover:bg-red-800`

**Disabled:** always `disabled:cursor-not-allowed disabled:opacity-50`.

**Story card** (`StoryRow` collapsed, `StorySheet` open):

`rounded-lg border border-paper-300 bg-white` + `p-3`

Hover `hover:border-ink/20`. Selected `border-bloom ring-2 ring-bloom/30`.

**Board column** (`BoardColumn`):

`flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper`

**Column header** (`ColumnHeader`):

`border-b border-paper-300 px-4 py-3 text-xs font-medium uppercase tracking-wider text-ink-700`

**Text input:**

`rounded-lg border border-paper-300 px-3 py-2 text-sm focus:border-bloom focus:outline-none focus:ring-1 focus:ring-bloom`

**Field label:** `text-sm font-medium text-ink-800`

**Settings card:** same surface as a story, padding `p-6`.

**Spacing:** story `p-3`; settings `p-6`; section `gap-6`.

**Motion:** `200ms`–`300ms`, `ease-out`. Expand / collapse, column toggle, illegal-drop snap-back. No bounce. No hover theatre. If a click hesitates, the design already failed.

**Icons:** Lucide, stroke 2px, colour of the nearby text.

| Meaning | Lucide |
| --- | --- |
| Feature | `star` |
| Bug | `bug` |
| Chore | `wrench` |
| Release | `flag` |
| Blocked | `octagon-alert` |
| Search | `search` |
| Close | `x` |
| Overflow / later panel actions | `ellipsis` — do **not** hide the next action here |

### Four columns, not Kanban

Left → right, always these names, this order:

**Icebox · Backlog · Current · Done**

They are four views of planning, not four stages of a card. A story does not “walk” Icebox → Backlog → Current → Done.

The **board** in this file is that projection (stories into Icebox / Backlog / Current / Done). It is a view, not a domain. There is no board API.

| Column | What it is | What it is not |
| --- | --- | --- |
| **Icebox** | `unscheduled` holding pen. Own order. No band headers. No points / velocity. | Not step 1 of a pipeline. Not in the plan. |
| **Backlog** | Later **computed bands** of the ranked list (velocity + estimates, live). Date-band headers once stories exist. | Not “ready”. Not a todo column. |
| **Current** | This band: in-flight + velocity-filled unstarted + accepted *in this window*. Bloom highlight on the header. | Not “in progress only”. Not Done. |
| **Done** | Accepted work that has **aged past the current window**. Flat list, newest accepted first. | Not a drop target. Not where Accept sends a story. Not grouped by window. |

Current header: **points / denominator** · `Ends {computed end}` · badge `Over capacity` when Feature-points exceed the denominator.

- **Cold start** (no completed window yet): denominator is initial velocity (default 10). Example `0 / 10`.
- **After the first completed window:** denominator is **capacity**. For a full band: velocity times length in days. For the open Current window: remaining days in the window times velocity.

Do not show a formula. Do not name capacity with a single letter. Say **capacity** in chrome if you need a word; the numbers are enough.

Current column header uses Bloom as the live-window highlight (guide). The other three headers stay ink-700 uppercase. That is how you know which band is live without memorising the UI.

The only setting is **window length in days** (default 7). Bands recompute live. You never assign a story to a numbered window.


---

## 3. Board chrome

### Top bar (`AppBar`)

One row, paper, bottom border `#E8DFD6`. No fat nav.

Left: **organisation name** (switcher) · **project name** (switcher) · optional workspace name (slice 25).

Centre: nothing. The board is the meeting.

Right: search field (placeholder `Search this project` — `/` focuses it) · `?` · user menu (Sign out, settings).

Viewers see the same bar. They do not see Invite or Settings that mutate.

You always know where you are: organisation + project are visible, Current’s header is Bloom, the selected `StoryRow` has the Bloom ring.

### Story row anatomy (`StoryRow`)

A ticket on paper, not a table row. One line when collapsed. `p-3`. Inter 400 title.

Left → right:

1. **Type icon** (Lucide, 16px).
2. **EstimateChip** — Features only in Phase 0–1. Pollen fill, ink text, the point or `·` if unestimated. Bugs / chores: no chip unless the Phase 3 toggle is on. Releases: no chip.
3. **Title**.
4. **Label chips** (meta, Inter 500). Epic label is Bloom-bordered — see slice 19. Ordinary labels are ink on paper-border pills. Do not invent a rainbow.
5. **Blocked** badge if any open blocker (`octagon-alert` + `Blocked`).
6. **Tasks** `done/total` once slice 9 exists.
7. **Owner initials** (up to 5). If they used an API token, initials are still the **user the token belongs to**.
8. **Projected date** (slice 20): computed band end date, or `Not scheduled` in Icebox. Quiet meta. Not a field.
9. **StateButton** — the one next verb for this type + state. Primary (Bloom) or Stem when the verb is Accept-as-Finish on a Chore/Release. Always visible. Not hover-only.

Selected: Bloom ring. Open (`StorySheet`): same card, expanded in the column (Tracker). Slice 29 can pin the sheet in a second pane.

**Do not** hide Start / Finish / Deliver / Accept / Reject / Restart behind `ellipsis`. Secondary verbs that are legal right now sit on the row quieter (Secondary class, or a text button in Inter 500). Illegal verbs are absent, not disabled-for-mystery.

### StateButton — the next verb

The button label **is** the verb. Not “Update status”.

#### Feature and Bug

| State | Column it may sit in | StateButton (one next action) | Quieter, still visible |
| --- | --- | --- | --- |
| `unscheduled` | Icebox | **Pull to Backlog** | Estimate chips. **Start** only if estimated (Start jumps to Current; may overflow). |
| `unstarted` | Backlog or Current | **Start** if estimated; **Estimate** (chips) if not | **Icebox** (unstarted only). Estimate chips if already estimated (can change / clear). |
| `started` | Current | **Finish** | — |
| `finished` | Current | **Deliver** | — |
| `delivered` | Current | **Accept** (Primary / Bloom) | **Reject** (Destructive, not hidden) |
| `rejected` | Current | **Restart** | Reason stays visible as meta. Not a delete. |
| `accepted` | Current until it ages past the current window, then Done | None | Accepted meta in Stem. No drag back to Backlog. |

Start on a Feature with no estimate: button disabled, `title` / accessible name: `Needs an estimate first`. Clicking or `S` does not change state; see slice 4 error copy.

Bug: same machine. No estimate required. No EstimateChip unless Phase 3 toggle.

#### Chore

| State | StateButton | Notes |
| --- | --- | --- |
| `unscheduled` | **Pull to Backlog** | Unestimated. No points. |
| `unstarted` | **Start** | |
| `started` | **Finish** | Finish **is** accept. Label stays **Finish**. Do not show Accept / Reject / Deliver. |
| `accepted` | None | Stem meta. |

#### Release

| State | StateButton | Notes |
| --- | --- | --- |
| `unscheduled` | **Pull to Backlog** | Pull auto-starts the marker. |
| `started` | **Finish** | Finish **is** accept. Optional **Target date** on the open sheet. |
| `accepted` | None | Colour: none if no target; Delivered blue `#3A6EA5` if the **computed band start** ≤ target; Rejected `#B42318` if the band start > target. |

Place the marker at the **end** of that milestone’s stories (work above, marker below). Open-sheet hint: `Sit this at the end of the stories it ships.`

### Legal drags (this is not Kanban)

| Drag | Result |
| --- | --- |
| Reorder inside Icebox | Icebox order. Persists. |
| Icebox → Backlog | **Schedule.** `unscheduled` → `unstarted`. Lands at the **bottom** of the ranked list. Release auto-`started`. |
| Backlog → Icebox | **Icebox.** `unstarted` only → `unscheduled`. Drops out of the plan. |
| Reorder inside the ranked list | Priority. Unstarted **cannot** sit above started / finished / delivered / rejected. Illegal drop snaps back. |
| Backlog → Current | **Rejected** while auto-plan is on. Does not Start. Snap back. **Slice 30 only (manual Current):** this drag is legal. `C` on a focused unstarted Backlog story does the same. Icebox → Current stays illegal. |
| Icebox → Current | **Rejected.** Does not Start. Use **Start** (after estimate on a Feature). |
| Anything → Done | **Rejected.** Done is aged-out accepted work only. |
| Accepted-this-window → Backlog | **Rejected.** |
| `started` / `finished` / `delivered` / `rejected` / `accepted` → Icebox | **Rejected.** |

Copy for the snap-back lives in slice 7.

Keyboard reorder is first-class (Tab, Space, arrows, Space). See §7.

### Band chrome (computed, not stored)

- **Current `ColumnHeader` extra line** (Inter 500, not uppercase): `{points} / {denominator}` · `Ends {computed end}` · badge `Over capacity` when points > denominator. Cold start: `0 / 10` (or the Owner’s initial velocity). After the first completed window: denominator is capacity (full band: velocity times length in days; open window: remaining days times velocity).
- **Backlog:** future **band** subheaders: `Ends {computed end}` · packed points / capacity for that band. Empty board: no fake headers (slice 0).
- **Done:** flat accepted list that has aged past the current window, newest accepted first. No window groups.
- **Icebox:** title + count only. Never a band header. Never a date.

The current window is the length-in-days band that contains now (default 7 days). Start of the first band is computed from project created; later bands follow every length-in-days. Changing length replans immediately.

No Recalculate button. The board after a mutation is already right.

---

## 4. Manifesto, applied

1. **One next action.** On the board it is the focused story’s `StateButton`, or **Add a story** on an empty Icebox. Never two Bloom buttons.
2. **Secondary stays visible.** Icebox, Reject, Estimate chips, Pull to Backlog — quieter, not inside More.
3. **Confirm vs undo.** Confirm: hits someone else, or cannot be undone (Reject + reason, Delete story, Revoke invite, Revoke token, Import CSV, destroy project / organisation, turn off auto-plan). Undo: local and reversible (Start, Finish, Deliver, Accept, Restart, reorder, estimate, schedule / icebox). Accept is **not** a confirm — History undoes it. Incomplete tasks / open blockers: **warn**, still accept.
4. **Copy.** Short. Tracker words are the user’s words. Strings are written below, not “the empty state should be friendly”.
5. **Empty / loading / error still have a next action.**
6. **Occasional user.** Power is chords and `?`, not extra chrome. No density-mode picker in MVP.
7. **Keyboard. Fast.** Optimistic verbs. Snap back on error. Jank is a bug.
8. **Infer.** See §6.
9. **Opinionated.** Default add is Icebox. Default type is Feature. Default scale is Linear 0/1/2/3. Default velocity 10. No Typeform onboarding.

Also: if it is not keyboardable and does not hold contrast, it is not done. Focus is visible (Bloom ring / input ring). You can name the organisation, the project, the column, and the selected story without hovering.

---

## 5. Patterns

### First: in-repo guide

Paper columns, white tickets, Bloom current, Stem accepted, Pollen estimates, Fraunces headings, Inter body, Lucide, pill buttons, 200–300ms ease-out. Four full-height columns.

### Then: classic Tracker, old Stripe, old Trello

| Steal | How Flower uses it |
| --- | --- |
| Tracker board | Columns are panels of one list. Expand the story in the column. Type icon + estimate + title + owners + state verb on one row. Tab / Space / arrows to reorder. Shift+letter toggles panels. `A` adds a story. `/` search. `?` the map. |
| Tracker verbs | Buttons named Start, Finish, Deliver, Accept, Reject, Restart. Not a status dropdown. |
| Tracker Icebox | Holding pen on its own order. Pull to plan; do not walk through it. |
| Tracker My Work | Owner’s in-flight + requester’s Delivered. `Shift+W`. |
| Tracker search | `type:`, `state:`, `estimate:-1`, `owner:`, `label:`, `is:blocked`, `includedone:`. |
| Old Stripe | One obvious submit. Quiet secondary. Calm sentences. Fast settings. |
| Old Trello | A card is a card. Drag is immediate. The thing you clicked is the thing that moved. |
| Basecamp (philosophy only) | Opinionated defaults. No cockpit. |

### Refuse

- New Stripe (gradient marketing chrome, hover reveal, motion for its own sake)
- GitHub PRs (latency, hover theatre, everything behind checks)
- Jira (custom fields, workflow designer, density as punishment)
- Qualtrics / Typeform (one question per screen, fake progress)
- Monday.com (rainbow boards, widget chrome)
- Linear’s navigation (command-palette-as-the-app, hidden verbs)
- Kanban stage columns
- A status `<select>`
- Hover-only Start / Accept
- Recalculate
- A date picker that moves the plan (Release **target date** is a comparison only)
- Presence avatars, live cursors (slice 17 forbids them)


---

## 6. Infer vs ask

Flower fills what it can. The human types what only a human knows.

| Thing | Infer (Flower) | Ask (human) |
| --- | --- | --- |
| Email | — | Email |
| Password or magic link | — | Their choice. Do not force both. |
| Username | Infer from email local-part; editable later. Locked 20 Aug 2026. | Do not ask at signup. |
| Email verification | Send the mail | They click the link |
| Organisation name | — | Name (slice 0, required) |
| Organisation slug / URL | From the name | Do not ask |
| First project name | — | Name (slice 0, required) |
| Project slug | From the name, unique per organisation | Do not ask |
| Creator roles | Organisation owner + project owner | — |
| Last project | Land here on sign-in | — |
| Story type | Feature | Change only if they mean Bug / Chore / Release (slice 18) |
| Story title | — | Title (required, max 500) |
| Description | Empty | Only if they have a brief |
| State on create | `unscheduled` in Icebox (default add) | — |
| Requester | Creating user | Change later to another Member/Owner |
| Owners | Start assigns the clicker | Add / remove others, max 5 |
| Follow | Requester + owners, cannot unfollow | Others may follow |
| Estimate | Unestimated. `0` is a value only if they pick it | 0 / 1 / 2 / 3 on a Feature before Start |
| Band / projected date | Planner assigns from velocity + rank + length in days. Live. | Never a due-date picker. |
| Current membership | Auto-plan up to velocity, leave short | Start (may overflow). No drag-to-Current |
| Velocity | Live from completed-story durations; initial 10 is bootstrap | Owner may set initial velocity and `velocity_strategy` 1–4 in settings |
| Window length | **Days** (default 7) | Owner may set length in days. |
| Timezone | Store; default `Australia/Melbourne` | Owner setting, not a signup field |
| Icebox vs Backlog on create | Icebox | Pull to Backlog, or add-in-Backlog as a quieter column action |
| New Icebox position | Top (most recently captured) | Reorder if they care |
| New scheduled position | Bottom of the ranked list | Reorder if they care |
| Who accepts | Any Member or Owner — no ACL picker | — |
| Reject reason | — | Required |
| Blocker | — | Free text; optional **story** picker (not a URL) |
| Linked-blocker resolve | Auto when the linked story is accepted or deleted | — |
| Mentions (slice 12b) | Resolve `@` against project Members/Owners/Viewers | The `@name` |
| Images | Clipboard paste → attachment + embed | File picker as quieter option |
| Activity | User, time, from → to | — |
| Reorders in activity | Do not log | — |
| API token | Same API and defaults as the frontend. Minted on the **user** (`/users/:id/tokens`). Bound to a project role. | Name, role, projects |
| Webhook (slice 23b) | Signing secret; `t=<unix>,v1=<hmac-sha256>` on each POST | URL |
| Token secret | Show once | They copy it now |
| Search operators | Parse Tracker-like language | The query |
| Saved search name | — | A name |
| My Work set | Owners in started/finished/delivered/rejected + requester’s delivered | — |
| CSV | Export columns as specified | File to import (creates only) |
| Point scale | Linear 0/1/2/3 | Owner changes later (slice 28) |

---

## 7. Keyboard

Scope: chords work on the **board** when focus is not in an input, textarea, or contenteditable. Inside those: `Esc`, `Ctrl/Cmd+S`, `Ctrl/Cmd+Enter` only (plus native typing). `?` opens this map. Viewers: navigation works; mutating chords announce `Viewers can look, not change.` (`aria-live`).

### Tracker chords Flower copies

From classic Tracker help. Same letters. Same jobs, mapped onto Flower panels.

| Key | Tracker | Flower |
| --- | --- | --- |
| `Shift+A` | Open skiplink navigation | **Skiplink.** Jump: Icebox · Backlog · Current · Done · open story · My Work (once it exists) · search. This is how you know where you are without the mouse. |
| `Shift+B` | Toggle Backlog | Toggle **Backlog** column |
| `Shift+C` | Toggle Current | Toggle **Current** column |
| `Shift+D` | Toggle Done | Toggle **Done** column. Not “story is done”. |
| `Shift+E` | Toggle Epics panel | Toggle **Epics** (slice 19). Until then the chord is reserved — do not reuse. |
| `Shift+H` | Toggle Project History | Reserved. No project history panel. Do not bind it. |
| `Shift+I` | Toggle Icebox | Toggle **Icebox** column |
| `Shift+L` | Toggle Labels | Focus / toggle the **label filter** (slice 11). Reserved until then. |
| `Shift+W` | Toggle My Work | Toggle **My Work** (slice 24). Reserved until then. |
| `A` | Add Story | **Add a story** in Icebox (`unscheduled`, Feature). Title focused. Not Accept. |
| `E` | Add Epic | **Add an epic** (slice 19). Reserved until then. |
| `/` | Search | Focus search / filter |
| `Tab` / `Shift+Tab` | Select story | Move story focus along the current column, then the next visible column. In dialogs, composer, and `StorySheet` fields: **native** Tab. |
| `Space` | Start / commit move | First `Space` picks up the focused story; arrows move one slot; second `Space` commits. `Esc` cancels. Illegal slot snaps back (same copy as drag). |
| `↑` `↓` | Move story (while picked up) | Same. When not picked up: move focus one story (with Tab). |
| `←` `→` | (not in Tracker table) | Move focus to the adjacent **column**, same vertical neighbourhood. |
| `Esc` | Collapse open story (no auto-save in Tracker) | Close `StorySheet` or dialog; return focus to the row. Flower **does** keep what the field already saved (title blurs save). Unsaved comment draft stays in the field if they reopen. |
| `?` | This help | Shortcut overlay. One next action: `Esc` / `Close`. |
| `Ctrl/Cmd+S` | Save open story | Save title + description on the open sheet. `preventDefault` so the browser does not save HTML. |
| `Ctrl/Cmd+Enter` | Save comment | Post the comment. |

Default add is **Icebox** (product spec slice 2). `A` does not create in the focused column.

### Flower additions Tracker never shipped

Slice 16 requires a keyboard path for estimate 0–3, start, finish, deliver, accept, reject (reason then confirm), restart, icebox / schedule.

| Key | Flower |
| --- | --- |
| `0` `1` `2` `3` | Set estimate on the focused Feature (legal states: any except accepted). `0` is estimated. Illegal on Bug/Chore/Release in MVP: announce `This type is not estimated.` Clear (unstarted Feature only): `-`. |
| `S` | **Start** the focused story if the machine allows. Feature + unestimated: no state change; announce the slice 4 copy. Assigns the clicker as owner when owners < 5. |
| `F` | **Finish.** Feature/Bug: `started` → `finished`. Chore/Release: Finish **is** accept. |
| `D` | **Deliver.** Unshifted. `Shift+D` still toggles Done. Illegal if not `finished`. |
| `Enter` | If the row is collapsed: activate the primary `StateButton` (Start / Finish / Deliver / Accept / Restart). If the sheet is already open: native. **Locked 20 Aug 2026:** `A` is Add Story. **Enter on a focused `delivered` row Accepts** (any Member/Owner). Click the title or `O` to open the sheet without activating the verb. |
| `O` | **Open** the focused story’s `StorySheet`. `Esc` closes. Click the title does the same. Locked 20 Aug 2026. Enter stays the verb, not open. |
| `R` | **Reject** when focused story is `delivered`. Opens the reason dialog. Confirm is **Reject**. Empty reason does nothing. **Locked 20 Aug 2026:** `R` is Reject only. Restart is **Enter** on a focused `rejected` row. |
| `I` | **Icebox** the focused `unstarted` story. `Shift+I` still toggles the column. Illegal states announce `Only unstarted stories can go to Icebox.` |
| `B` | **Pull to Backlog** (schedule) from Icebox. `Shift+B` still toggles the column. |
| `C` | **Slice 30 only (manual Current):** move the focused unstarted Backlog story into Current. `Shift+C` still toggles the column. Illegal while auto-plan is on: announce `Start the story to bring it into Current.` Icebox → Current stays illegal. |

`Shift+letter` = panel. Unshifted letter = verb or add. Do not invent a leader key.

### Viewer, dialog, browser

| Situation | Behaviour |
| --- | --- |
| Viewer + `S` / `A` / `0`–`3` / drag commit | `aria-live`: `Viewers can look, not change.` |
| Focus in input | Chords off, except `Esc`, `Ctrl/Cmd+S`, `Ctrl/Cmd+Enter` |
| `/` vs Firefox quick find | `preventDefault` on the board |
| `Space` vs page scroll | `preventDefault` when a story is focused or picked up |
| `?` overlay | Trap focus. `Esc` closes. |

---

## 8. Shared states and copy

Use these strings. Do not rewrite them to sound “on brand”.

### Confirm vs undo (product-wide)

| Action | Pattern | Copy |
| --- | --- | --- |
| Start, Finish, Deliver, Accept, Restart | Optimistic + **Undo** on the latest state-changing activity (slice 15). Optional quiet toast. | Toast: `Started.` / `Finished.` / `Delivered.` / `Accepted.` / `Restarted.` Action link: `Undo` |
| Estimate, reorder, schedule, icebox | Optimistic. Snap back on error. No confirm. | — |
| Reject | **Confirm.** Reason required. Hits the owners. | See slice 6 |
| Accept with incomplete tasks or open blockers | **Warn, still accept.** No hard block. | `Accepted, with unfinished tasks.` / `Accepted, with open blockers.` + `Undo` |
| Delete story | **Confirm.** Not covered by state-Undo. | `Delete this story? This cannot be undone.` Buttons: `Delete` / `Keep it` |
| Revoke invite | **Confirm** | `Revoke this invite? The link will stop working.` `Revoke` / `Keep invite` |
| Revoke token | **Confirm** | `Revoke this token? It will stop working now.` `Revoke` / `Keep token` |
| Import CSV | **Confirm** | `Import will create these stories. It cannot update existing ones.` `Import` / `Cancel` |
| Turn off auto-plan (slice 30) | **Confirm** | `Current will no longer fill from velocity. Future windows still will.` `Turn off` / `Keep automatic` |

Never confirm Accept, Start, or a reorder.

### Global loading / error

| State | Copy | Next action |
| --- | --- | --- |
| First board paint | Column skeletons (paper, no fake titles). | Wait. No spinner-on-click for verbs. |
| Mutation in flight | The row already shows the new state (optimistic). | — |
| Mutation failed | `That didn’t save. The story is back as it was.` | Retry is the same verb, still on the row. |
| Lost session | `You’re signed out.` | `Sign in` |
| Forbidden (viewer) | `Viewers can look, not change.` | — |
| Cross-tenant / missing | `We can’t find that.` (404, never “you aren’t allowed to see org B”) | `Back to your projects` |
| Stale socket (slice 17) | `Connection lost. We’ll reconnect.` | Automatic reconnect. Secondary: `Refresh` |


---

## 9. Flows

Each slice: screen, one next action, secondary (visible, quieter), confirm vs undo, empty / loading / error / success.

Phase 0 Features only unless noted. Bug / Chore / Release controls appear when slice 18 has landed — same patterns, different `StateButton` table.

---

### Slice 0 — Sign up, organisation, project, empty board

**Screens:** `SignUp` → (`CheckEmail` if password) → `NameOrganisation` → `NameProject` → `Board` (four empty columns).

**SignUp**

- One next action: **Sign up** (Primary).
- Secondary: **Email me a link instead** (Secondary). Sign in link in the footer, quieter: `Already a user? Sign in`.
- Fields: Email, Password. Username: **do not ask.** Infer from the email local-part. No username field in slice 0.
- Confirm: none. They can sign out.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `Email and a password. That’s enough to start.` | `Sign up` |
| Loading | `Signing you up…` | — |
| Error (taken) | `That email already belongs to a user.` | `Sign in` |
| Error (blank) | `Need an email and a password.` | Fix and `Sign up` |
| Success (password) | `Check your email to verify.` | Open mail; the verify link continues. |

**CheckEmail / verify**

- One next action: the link in the email.
- Unverified password signup **cannot** create an organisation (AC). If they sneak to the app: `Verify your email to continue.` + `Resend the email`.

**Magic link (no user yet)**

- One next action: **Email me a link**.
- Success: `Check your email for a link.` Same organisation + project flow after the click.

**SignIn**

- One next action: **Sign in**.
- Secondary: **Email me a link instead**.
- Success: last project. If none (shouldn’t happen after slice 0): `NameOrganisation`.

| State | Copy | Next |
| --- | --- | --- |
| Error (bad password) | `That email and password don’t match.` | Try again, or `Email me a link instead` |
| Loading | `Signing in…` | — |

**NameOrganisation**

- One next action: **Create organisation** (Primary).
- Field: `Organisation name`. Infer slug.
- Empty: `What’s the organisation called?`
- Error (blank): `Name the organisation.`
- Success: `NameProject`.

**NameProject**

- One next action: **Create project**.
- Field: `Project name`. Infer slug, window length 7 days, initial velocity 10, Linear 0/1/2/3, timezone default.
- Empty: `Name the first project.`
- Error (blank): `Name the project.`
- Success: `Board`.

**Empty Board**

Four `BoardColumn`s. Current header may show the computed band end and `0 / 10`. No fake stories. No fake Backlog band headers.

| Column | Empty copy | One next action |
| --- | --- | --- |
| Icebox | `Nothing waiting. Capture a story.` | `Add a story` (also `A`) |
| Backlog | `Nothing ranked yet. Pull from Icebox when you’re ready.` | Quieter: none until Icebox has a row. Do not put a second Add here as Primary. |
| Current | `Nothing in this window yet. Pull a story from Icebox or create one.` | Same as velocity doc. Cold-start header `0 / 10`. |
| Done | `Nothing in Done yet. Accepted work lands here after it ages out of this window.` | — look at Current |

Secondary on the board: Invite (slice 1, Owners), project switcher, `?`.

Viewer: same empty columns, no `Add a story`.

Reload stays in this organisation / project.

---

### Slice 1 — Invite a member

**Screen:** `InviteDialog` from a quieter `Invite` Secondary in the AppBar (Owners and organisation owners only). `MemberList` on project settings (`p-6` card).

- One next action: **Send invite**.
- Fields: Email, Role (`member` default; `owner` / `member` / `viewer`).
- Secondary: **Resend**, **Revoke** on each pending row (Revoke confirms).
- Members / Viewers do not see `Invite`.

| State | Copy | Next |
| --- | --- | --- |
| Empty list | `No one else is on this project yet.` | `Invite` |
| Loading | `Sending invite…` | — |
| Error (already a member) | `They’re already on this project.` | Close |
| Error (blank email) | `Need an email.` | Fix |
| Success | `Invite sent. It expires in 14 days.` | Pending row appears |
| Resend | `New link sent. The old one no longer works.` | — |
| Revoke confirm | `Revoke this invite? The link will stop working.` | `Revoke` / `Keep invite` |
| Expired | `This invite expired.` | Owner: `Resend` |

**Accept invite (new email):** signup (password or magic link) → land on the project as the invited role.

**Accept invite (existing user):** sign in → project is in their list → board.

**Viewer board:** they can see; they cannot create, invite, or change settings. Mutating controls absent, not merely disabled.

---

### Slice 2 — Create a Feature in Icebox

**Screen:** `Board` + `StoryComposer` at the **top** of Icebox.

- One next action: **Add** (or `Enter` in the title).
- Field: Title only. Infer type Feature, state `unscheduled`, requester = me, no estimate, no band.
- Secondary: `Cancel` (`Esc`) discards an empty composer. Type switch is absent until slice 18.
- Confirm: none.

| State | Copy | Next |
| --- | --- | --- |
| Composer empty | Placeholder `A feature title` | Type, `Enter` |
| Error (empty title) | `Need a title.` | Type one |
| Error (over 500) | `Titles stop at 500 characters.` | Shorten |
| Viewer | Composer absent | — |
| Success | Row appears at the **top** of Icebox. Current / Backlog / Done unchanged. | Estimate or pull (slices 3–4), or `A` again |

Icebox with rows: composer is a quiet `Add a story` at the top, not a modal.

---

### Slice 3 — Pull Icebox → Backlog (and back)

**Screen:** `Board`. Drag, or **Pull to Backlog** on the Icebox row, or `B` with the row focused.

- One next action on an Icebox Feature: **Pull to Backlog**.
- Secondary: reorder Icebox; Estimate; Start if estimated (slice 4 — jumps to Current, not Backlog).
- Result: `unstarted`, **bottom** of the ranked list, leaves Icebox. Velocity may then auto-fill it into Current if it fits — that is slice 8, not a drag-to-Current.
- Icebox an `unstarted` Backlog row: **Icebox** / `I`. Becomes `unscheduled`. Loses the projected date.
- Confirm: none. Undo: drag back, or the opposite verb.

| State | Copy | Next |
| --- | --- | --- |
| Empty Icebox (after the last pull) | `Nothing waiting. Capture a story.` | `Add a story` |
| Illegal icebox (started+) | `Only unstarted stories can go to Icebox.` | — |
| Viewer | No pull handle | — |
| Success (scheduled) | Story in Backlog (or Current if auto-plan takes it). | Estimate / Start / reorder |
| Success (iceboxed) | Story at the top of Icebox. `Not scheduled.` | — |

---

### Slice 4 — Estimate and Start

**Screen:** `StoryRow` EstimateChips `0` `1` `2` `3` (Pollen) + **Start**.

- One next action if unestimated: pick a point. Chips **are** the action.
- One next action if estimated: **Start**.
- Secondary: change estimate; clear with `—` / `-` while `unstarted`; Pull / Icebox as legal.
- Start assigns the clicker as owner and follower. If 5 owners already and the clicker is not among them: Start still happens; announce `Started. Couldn’t add you as owner — this story already has 5.`
- Start from Icebox: estimate first (or estimate then Start). Story becomes `started` in **Current**, not `unstarted` in Backlog. May overflow velocity.
- Viewer: chips visible, not clickable.
- Confirm: none. Undo: slice 15 (`started` → `unstarted`).

| State | Copy | Next |
| --- | --- | --- |
| Unestimated Feature, Start clicked or `S` | `Estimate this feature before you start it. 0 is fine.` | Chips |
| Clear while started | Chips locked. `Started features keep their estimate.` | — |
| Viewer Start | `Viewers can look, not change.` | — |
| Success | Row is `started` in Current. Pollen indicator. You are an owner. | **Finish** (slice 5) |

---

### Slice 5 — Finish, Deliver, Accept

**Screen:** Current `StoryRow`. Accepted stays in **Current**. Done stays empty until that work ages past the current window.

| Now | One next action | Secondary |
| --- | --- | --- |
| `started` | **Finish** | — |
| `finished` | **Deliver** | — |
| `delivered` | **Accept** (Bloom, any Member or Owner) | **Reject** (slice 6) |

Requester *should* accept. The UI does not lock the button. My Work (slice 24) is where the requester finds Delivered; until then, Delivered is obvious in Current (Delivered blue meta).

Viewer: verbs absent. API tokens use the same Accept path as a person in that role. No token-only chrome.

Confirm: none on Finish / Deliver / Accept. Undo: activity.

Incomplete tasks / blockers (later slices): accept anyway; toast `Accepted, with unfinished tasks.` + `Undo`.

| State | Copy | Next |
| --- | --- | --- |
| Illegal verb (Finish on unstarted, Accept on finished) | `That move isn’t available yet.` | The legal StateButton |
| Success Accept | `Accepted.` Story remains in Current, Stem meta. | Next delivered, or keep working |
| Done still empty | Correct until the accept ages past this window | — |

---

### Slice 6 — Reject and Restart

**Screen:** `RejectDialog` from **Reject** or `R` on a `delivered` row.

- One next action in the dialog: **Reject** (Destructive), disabled until the reason has text.
- Secondary: **Keep delivered**.
- Confirm: yes (this *is* the confirm). Reason required. Hits the owners (email + in-app once slice 12b exists).
- Result: `rejected`, **still Current**, reason on the row and in activity.
- Then one next action: **Restart** (or Enter on that focused row). → `started`, still Current. Reason remains in activity. They Finish + Deliver again. A new Accept is required.
- Fail the product if Reject jumps to `started`, or if `rejected` looks terminal like Accepted (no Restart, or it sits in Done).

| State | Copy | Next |
| --- | --- | --- |
| Dialog | Title `Reject this story?` Body label `Why` placeholder `What needs to change?` | Type, `Reject` |
| Empty reason | `Say why, then reject.` | Type |
| Viewer | Control absent | — |
| Success reject | `Rejected.` Reason visible. | `Restart` |
| Success restart | `Restarted.` | `Finish` |

---

### Slice 7 — Reorder

**Screen:** `Board`. Drag or Tab → Space → arrows → Space.

- One next action while picked up: place it, `Space` to commit.
- Secondary: `Esc` cancel.
- Confirm: none. Persist on reload. Viewers cannot.

**Snap-back copy** (also `aria-live`):

| Illegal drop | Copy |
| --- | --- |
| Unstarted above started / finished / delivered / rejected | `Unstarted stories stay below work that’s already started.` |
| Backlog or Icebox → Current | `Start the story to bring it into Current. Dragging doesn’t start it.` |
| Accepted-this-window → Backlog | `Accepted work stays in Current until it ages out of this window.` |
| In-flight / accepted → Icebox | `Only unstarted stories can go to Icebox.` |
| Anything → Done | `Done is aged-out accepted work. You can’t drag here.` |

Success: the row sits where they put it. Auto-plan may move **other** unstarted rows between Current and Backlog. That is live replan, not a fight with their drop.

Icebox reorder is independent.

---

### Slice 8 — Auto-plan from velocity

**Screen:** same `Board`. No new verb. This is the Tracker moment.

- One next action on a new project with estimated Features in the ranked list: **read Current** (already filled, left **short** of the cold-start denominator) or **Start** a Backlog story (may overflow).
- Secondary: **window length in days** in project settings (Owner). Default 7. Changing length replans live. No Recalculate.
- Current header, cold start: `6 / 10` · `Ends 9 Aug`. After the first completed window: `{points} / {capacity}` (capacity as defined in band chrome). Over: points above capacity + badge `Over capacity`.
- Accepted in this window: still in Current.
- When accepted work ages past the current window: those rows sit in Done as a **flat list**. Velocity updates from completed-story durations (`started_at` → `accepted_at`) over the last number of completed windows set by `velocity_strategy` (default 3, setting 1–4).

| State | Copy | Next |
| --- | --- | --- |
| Empty ranked list | Current empty copy from slice 0 | Add or pull |
| Left short | Header shows e.g. `6 / 10` on cold start. No apology chrome. | Start something if they mean to overflow |
| Over (because Start) | Badge `Over capacity` | Finish the work; do not kick rows out |
| Cold start | Header uses initial velocity **10**, never a fake 0 | — |

The window boundary is midnight, project timezone. No close-window button.


---

### Slice 9 — Tasks

**Screen:** `TaskList` on `StorySheet`. Unowned, unpointed.

- One next action: **Add a task** (empty) or toggle the next unchecked box.
- Secondary: reorder tasks (drag / Space pattern).
- Completing all tasks does **not** Finish or Accept.
- Viewer: see, no toggle.
- Accept with incomplete: warn, still accept.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `Break the work into checks if you want. Not required.` | `Add a task` |
| Blank submit | `Need some text.` | Type |
| Success toggle | `done/total` on the row updates | Next box |
| Accept warn | `Accepted, with unfinished tasks.` | `Undo` if they didn’t mean it |

---

### Slice 10 — Blockers

**Screen:** `BlockerForm` on `StorySheet`. Badge on `StoryRow`.

- One next action: **Add a blocker**.
- Fields: required text; optional **story picker** (same project). Not a URL.
- Secondary: **Clear** on a blocker (Member/Owner).
- Does not block Start / Finish / Deliver / Accept. Accept warns.
- Auto-resolve copy in activity: `Blocker cleared because {title} was accepted.` (or deleted).

| State | Copy | Next |
| --- | --- | --- |
| Empty | `What’s in the way?` | Type, optional story, `Add a blocker` |
| Blank | `Say what’s blocking it.` | Type |
| Success | Badge `Blocked` on the row | Keep working, or Accept with warning |
| Accept warn | `Accepted, with open blockers.` | `Undo` |

---

### Slice 11 — Labels

**Screen:** label pills on the row; `LabelInput` on the sheet; filter in the AppBar / `Shift+L`.

- One next action on the sheet: add a label (lowercase `[a-z0-9-]+`, 1–32, max existing `VARCHAR(100)`).
- One next action on the board once labels exist: **filter** (hides non-matches in all columns; rank unchanged).
- Secondary: Owner deletes an **unused** project label in settings.
- Viewer: filter yes, add no.
- Epic-as-purple is slice 19. This slice is ordinary labels.

| State | Copy | Next |
| --- | --- | --- |
| Empty on story | Placeholder `label` | Type, Enter |
| Bad name | `Use lowercase letters, numbers, and hyphens.` | Fix |
| Filter on, zero rows | `Nothing with that label.` | `Clear filter` |
| Success | Pill on the row | Filter or keep adding |

---

### Slice 12 — Owners, requester, follow

**Screen:** `PeopleFields` on `StorySheet`. No notification chrome. Mentions and mail are slice 12b.

- One next action on a story with no owner: **Add owner** (Start already did this for the clicker).
- Secondary: change requester (another Member or Owner on the project); Follow (if allowed).
- Max 5 owners. Sixth: `This story already has 5 owners.`
- Requester and owners: followed, **cannot** unfollow. Control reads `Following` and is disabled, title `Requesters and owners stay subscribed.`
- Viewer: see owners and requester; cannot assign or follow.

| State | Copy | Next |
| --- | --- | --- |
| Empty owners | `No owners yet. Start assigns you.` | `Add owner` or Start |
| Success | Initials on the row | — |

---

### Slice 12b — Mentions and notifications

Product order is 12 → 12b → 13. `@` is allowed in a description or a comment. The description picker can ship in 12b (description already exists). The comment picker attaches to `CommentList` when slice 13 lands. Do not block 12b on 13.

**Screen:** `@` picker in the description (and in comments once slice 13 exists); in-app notice + email. No notification cockpit on the board.

- One next action when typing `@`: pick a project Member, Owner, or Viewer.
- Unresolved `@` stays plain text.
- Mail + in-app: mention, assigned as owner, delivered → requester, rejected → owners.
- Viewer can be mentioned (email) and cannot follow, comment, or assign.
- In-app notice: one next action is **Open** the story. Secondary: dismiss.

| State | Copy | Next |
| --- | --- | --- |
| Empty notices | `Nothing new.` | Back to the board |
| Unresolved @ | Stays as typed | — |
| Success mention | They get in-app + email | `Open` |

---

### Slice 13 — Description and comments

**Screen:** Markdown `Description` + `CommentList` on `StorySheet`.

Allowed: h1–h3, bold, italic, strikethrough, lists, fenced code, inline code, block quotes, links. No raw HTML. `javascript:` rejected.

- One next action: **Post** on a comment (`Ctrl/Cmd+Enter`), or **Add description** when empty.
- Secondary: Edit (shows Edited), Delete comment (tombstone, not silent erase).
- Empty description is legal. Empty comment is not.
- Viewer: read only.
- `@` picker on comments is slice 12b, attached here when both have landed. Do not invent a second mention UI.

| State | Copy | Next |
| --- | --- | --- |
| Empty description | `Add a brief if the title isn’t enough.` | Type, `Ctrl/Cmd+S` |
| Empty comment | Placeholder `Write a comment` | Type, `Post` |
| Blank post | `Need some text.` | Type |
| Bad URL | `That link isn’t allowed.` | Fix |
| Deleted comment | `This comment was deleted.` | — |
| Success | Comment with author + time | — |

---

### Slice 14 — Images

**Screen:** paste into the sheet; quieter **Upload image**.

- One next action: **paste** (infer).
- Rules: png/jpg/jpeg/gif/webp, ≤ 10 MB, max 20 per story.
- Remote hotlinks do not render.
- Delete → embed shows missing-image.
- Viewer: see, no upload.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `Paste a screenshot.` | Paste |
| Wrong type | `Images only — png, jpg, gif, or webp.` | Pick another |
| Too big | `Images need to be 10 MB or smaller.` | — |
| At 20 | `This story already has 20 images.` | Delete one |
| Missing embed | `Image deleted.` | — |
| Success | Image in the description | — |

---

### Slice 15 — Activity and Undo

**Screen:** `ActivityList` on `StorySheet`.

- One next action when the latest event is state-changing: **Undo**.
- Secondary: read the list (created, estimated, scheduled/iceboxed, started, finished, delivered, accepted, rejected + reason, restarted, requester/owners/title, blocker, attachment, comment, label). User, time, from → to.
- Reorders do not appear.
- Undo restores previous state. Undo is itself an activity.
- Cannot undo a comment with this control.
- Viewer: read, no Undo.

| State | Copy | Next |
| --- | --- | --- |
| Empty (new story) | `Created just now.` | — |
| Undo Accept | Story is `delivered` again. `Undid accept.` | — |
| Undo on a comment selected | Control absent | — |
| Success | New activity line | — |

---

### Slice 16 — Keyboard morning

**Screen:** `Board` + `ShortcutOverlay` (`?`).

- One next action: do the work without the mouse — focus, open, create, estimate, verbs, reorder, icebox/schedule, search.
- Secondary: `?` map. Overlay is not the product.
- Viewer: navigate and open; mutating chords announce no permission.
- Fail the slice if any AC verb has no keyboard path. Letter keys plus Enter-on-StateButton and `O` to open are specified in §7. Locked 20 Aug 2026.

| State | Copy | Next |
| --- | --- | --- |
| Overlay | Title `Keyboard`. Table from §7. | `Esc` |
| No permission | `Viewers can look, not change.` | — |
| Success | A morning without the mouse | — |

---

### Slice 17 — Live updates

**Screen:** same `Board` / `StorySheet`. No presence. No live cursors.

- One next action: keep typing. Flower patches the other session within 2 seconds.
- Secondary: `Refresh` only when the socket is dead.
- Do not steal focus. Do not discard a comment draft unless that comment or story was deleted (`This story was deleted.` / `That comment was deleted.`).

| State | Copy | Next |
| --- | --- | --- |
| Live | No chrome | — |
| Dropped socket | `Connection lost. We’ll reconnect.` | Wait; quieter `Refresh` |
| Reconnected | Banner clears. No parade. | — |
| Refresh heal | Full board reload, focus restored if the story still exists | — |

---

### Slice 18 — Bug, Chore, Release

**Screen:** type control on `StoryComposer` / `StorySheet`. Default **Feature**.

- One next action: still **Add** / the type’s StateButton.
- Secondary: type segmented control — Feature · Bug · Chore · Release. Visible, quieter than the title. Only while `unscheduled` or `unstarted`. Change to Release clears estimate.
- Bugs: Start without points; 0 velocity.
- Chores: Finish accepts; no Reject.
- Releases: target date optional (`type=date`). Hint: `A comparison, not a deadline we plan to.` Colour rules in §3. Place at the end of the milestone. Empty above: `No stories above this release.`

| State | Copy | Next |
| --- | --- | --- |
| Type locked (started+) | `Type is locked after the story is started.` | — |
| Release in Icebox | No colour. `Not scheduled.` | Pull to Backlog (auto-start) |
| Success | Row icon matches type | The type’s next verb |

---

### Slice 19 — Epics

**Screen:** `EpicPanel` (`Shift+E`) + epic pill on the row (one epic label). Ordinary labels still allowed.

- One next action on empty: **Add an epic** (`E`).
- Secondary: assign / remove epic on a story; filter; Owner delete unused epic label (stories remain).
- Progress: `accepted Feature points / estimated Feature points`. Unestimated add 0 to the denominator. Independent epic **order**, not story rank.
- Viewer: view + filter.

Tracker epics are purple. The guide has no purple. **Do not invent a purple hex.** Epic pill: Bloom `#C43B6E` border and Inter 500 Bloom text, paper fill — a label, not a selection. Selection remains `ring-2 ring-bloom/30` on the card.

| State | Copy | Next |
| --- | --- | --- |
| Empty panel | `Group a theme. Not a parent ticket.` | `Add an epic` |
| Success | Bloom-bordered pill + progress | Filter or rank epics |

---

### Slice 20 — Projected date

**Screen:** date meta on `StoryRow` / sheet. Not a field.

- One next action: none to type — **reorder or estimate** if the date is wrong. The planner moves.
- Icebox: `Not scheduled.`
- Release date = computed **band end** that contains the marker. Target date is the only date they type. It does not move stories.
- Live on reorder / estimate / accept / velocity change.

| State | Copy | Next |
| --- | --- | --- |
| Scheduled | computed band end in project timezone | — |
| Icebox | `Not scheduled.` | Pull to Backlog |
| Late release | Marker uses `#B42318` | Reorder work above it, or live with red |

---

### Slice 21 — Velocity chart and burn-up

**Screen:** `Charts` (settings or a quiet AppBar link). Viewers can see.

- One next action: **read**, then back to the board (`Esc` / `Back to the board`).
- Velocity: bar per **completed window**, derived from completed-story durations in that window; line at velocity. While the corpus is empty, line label `Initial 10`. Current is not a completed bar. Optional faint “accepted so far”.
- Burn-up: cumulative accepted Feature estimates vs scoped Feature estimates, derived from stories (live now).

| State | Copy | Next |
| --- | --- | --- |
| Empty velocity | `No completed stories yet.` | `Back to the board` |
| Empty burn-up | `No scope yet.` | Add estimated Features |
| Success | Real bars only | — |

No fake history.

---

### Slice 22 — Search

**Screen:** AppBar search (`/`). Results do **not** reorder the board. Overlay list or a results column (slice 29 can pane it).

Operators (combinable): `type:feature|bug|chore|release`, `state:`, `estimate:-1`, `estimate:0`…`3`, `owner:name`, `label:api`, `is:blocked`, `includedone:`. Done excluded unless `includedone:`.

- One next action: type the query.
- Empty result one next action: **Clear**.
- Scoped to this project.

| State | Copy | Next |
| --- | --- | --- |
| Empty query | Placeholder `Search this project` | Type or `/` |
| No hits | `Nothing matches. You can clear the search.` | `Clear` |
| Success | Result rows; click opens / reveals on the board | Open |

Reveal (Tracker steal): a quieter control to scroll the story into its column. You do not drag to prioritize from search (Tracker). Rank on the board.

---

### Slice 23 — API token

Generic API tokens, minted on the **user** (`/users/:id/tokens`). Same API as the frontend. Same Owner / Member / Viewer permissions. A Member token can accept. No special token type. No scope checklist.

**Screen:** `TokenForm` on **user** settings. Each user manages their own tokens.

- One next action: **Create token**.
- Fields: Name, **Role** (`member` default; you can pick a role at or below your own on the selected projects), projects you belong to. The token is that role on those projects.
- You may mint a role at or below your own. A Member cannot mint an Owner token.
- After create: secret once. One next action: **Copy secret**.
- Secondary: **Revoke** (confirm).

| State | Copy | Next |
| --- | --- | --- |
| Empty | `A named token with a project role. Same permissions as you in that role.` | Name + role + `Create token` |
| Secret shown | `Copy this now. We won’t show it again.` | `Copy secret` |
| Copied | `Copied. Store it with the token.` | Done |
| Revoke | `Revoke this token? It will stop working now.` | `Revoke` |

Activity **user** is the person the token belongs to. Token name may appear in the summary as `via {token name}`. Not a distinct identity. Always say **user**.

Bearer only. Webhooks are slice 23b.

---

### Slice 23b — Project webhook

Outbound hook on the **project**. Owner or Member. Viewer cannot register. Not a token. Not a second API.

**Screen:** `WebhookForm` in project settings (`p-6`).

- One next action when none exist: **Add webhook**.
- Field: URL. Infer the signing secret. Flower POSTs with header `t=<unix>,v1=<hmac-sha256>`.
- After create: secret once. One next action: **Copy secret**.
- One next action when a hook exists: **Revoke** (confirm). Then they can add again.
- Viewer: page absent.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `Flower will POST when stories change.` | URL + `Add webhook` |
| Loading | `Saving webhook…` | — |
| Error (blank) | `Need a URL.` | Fix |
| Error (bad URL) | `That URL isn’t allowed.` | Fix |
| Secret shown | `Copy this now. We won’t show it again.` | `Copy secret` |
| Success | Hook row with the URL | — |
| Revoke confirm | `Revoke this webhook? Posts will stop.` | `Revoke` / `Keep webhook` |
| Viewer | Page absent | — |

---

### Slice 24 — Saved search and My Work

**Screens:** `SaveSearchDialog`; `MyWork` panel (`Shift+W`).

**Save search**

- One next action: **Save** with a name (per user per project).
- Secondary: reopen (applies the same filter); **Delete** (deletes the save, not stories).
- Empty name: `Name this search.`

**My Work** (Members and Owners only)

Contains: stories **I own** in `started|finished|delivered|rejected`, plus stories **I request** that are `delivered`.

- One next action: the top row’s `StateButton` (Finish / Deliver / Accept / Restart as the machine says). For a requester, that is usually **Accept**.
- Secondary: open the story; reveal it on the board. **No drag reorder** in My Work (Tracker).
- Viewers: no My Work; they may search.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `Nothing on you.` | `Back to the board` (or `A` to capture work) |
| Loading | Skeleton rows | — |
| Error | `Couldn’t load My Work.` | `Retry` |
| Success | Rows with the same StateButton as the board | Accept / Finish / … |


---

### Slice 25 — Workspace

**Screen:** `WorkspaceSwitcher` + `WorkspaceEditor`. Personal, one organisation. Default: all projects they can open.

Acceptance story: create personal workspace **Core** containing **Trail** and **Checkout** in one org.

- One next action when none exist: **Create workspace**. Name it `Core`, then **Add project**.
- Picker lists Trail, Checkout, and any other project in **this** organisation they can open. Another organisation’s project is not in the picker.
- One next action when Core exists: **Open**.
- Secondary: remove a project from the set; **Remove workspace**. Confirm: `Remove this workspace? Projects stay.` Deleting Core does not delete Trail or Checkout.
- Permissions stay per project. Viewer on Checkout stays Viewer inside Core.

| State | Copy | Next |
| --- | --- | --- |
| Empty custom | `All your projects in {org} are already the default.` | `New workspace` |
| Picker | This organisation only | Add Trail, add Checkout |
| Other-org project | Absent from the picker | — |
| Delete confirm | `Remove this workspace? Projects stay.` | `Remove` / `Keep` |
| Success | **Core** in the AppBar | Open Trail or Checkout (slice 29 can pane both) |

---

### Slice 26 — CSV

**Screen:** Owner settings. Members cannot export (assumption). Viewers neither.

- One next action: **Export CSV** or **Import CSV**.
- Import creates only, all-or-nothing. Confirm before run.
- Attachments and comments are not in the file.

| State | Copy | Next |
| --- | --- | --- |
| Export success | File downloads | — |
| Import confirm | `Import will create these stories. It cannot update existing ones.` | `Import` |
| Import fail | `Nothing was imported. Fix the file and try again.` | Pick file |
| Forbidden | Controls absent for Member/Viewer | — |

---

### Slice 27 — Cycle time

**Screen:** `CycleTime` chart. First `start` → `accepted`. Reject does not reset.

- One next action: read p50 / p75 / p95. Filter by type (quieter).
- No per-person leaderboard. Ever.

| State | Copy | Next |
| --- | --- | --- |
| Empty | `No accepted stories in this range yet.` | `Back to the board` |

---

### Slice 28 — Point scale and bugs/chores points

**Screen:** Owner project settings.

- One next action: pick `linear` (0,1,2,3) / `fibonacci` (0,1,2,3,5,8) / `powers_of_two` (0,1,2,4,8) / `custom`.
- Custom must offer **Revert to linear** without converting history.
- Switching does not convert estimates. Illegal-on-new-scale values stay until edited. Chips on those rows still show the old number.
- Toggle: **Bugs and chores may be given points** (default off). Turning it **off** is allowed.
- Members: read-only.

| State | Copy | Next |
| --- | --- | --- |
| Revert | `Back to 0, 1, 2, 3. Existing estimates stay as they are.` | `Revert to linear` |
| Toggle off | `New bugs and chores won’t take points. Existing points stay until edited.` | Save (Owner) |

Keyboard 0–3 remains. Extra values on Fibonacci / powers / custom appear as extra chips — letter keys for 5/8 are **not** specified (open if that slice needs chords).

---

### Slice 29 — Side-by-side panels

**Screen:** `SplitBoard`. Two panes. Each pane is: board (any column set), open story, My Work, search results, Icebox-only, or another project in the workspace.

- One next action: **work in the focused pane**. The other recedes (no Bloom header unless its Current is showing).
- Secondary: split / unsplit (quieter). Reload restores the last split. Both panes live-update (slice 17).
- `Shift+A` skiplinks include both panes.
- Tracker steal: toggle columns inside a pane with Shift+B/C/D/I/W/E/L. Refuse Linear’s “everything is a floating palette”.

| State | Copy | Next |
| --- | --- | --- |
| First split | Two panes, last-focused story or My Work on the right | Continue the verb they were on |
| Reload | Same split | — |

Default before they split: four columns in one pane (the MVP board).

---

### Slice 30 — Manual Current

**Screen:** Owner settings + a quiet badge on Current: `Manual`.

- One next action when turning off: confirm **Turn off**.
- Current then holds only in-flight, accepted-this-window, and stories **explicitly moved** there (drag Backlog → Current, or `C` on a focused unstarted Backlog story). No velocity fill into Current. Backlog future windows still auto-plan. Icebox → Current stays illegal.
- Restore: **Use automatic planning** — replans Current from velocity. No confirm (reversible by turning off again; the replan is the product doing math, not a delete).
- Default remains automatic. Not MVP.

| State | Copy | Next |
| --- | --- | --- |
| Confirm off | `Current will no longer fill from velocity. Future windows still will.` | `Turn off` / `Keep automatic` |
| Manual Current empty of unstarted | `Nothing pulled into Current. Start a story or move one here.` | Start, or move |
| Restored | Header loses `Manual`. Current packed, left short | — |

---

## 10. Role chrome (do not contradict the spec)

| | Owner | Member | Viewer |
| --- | --- | --- | --- |
| Add / estimate / verbs including Accept & Reject | Yes | Yes | No |
| Invite, roles, window length (days) / scale / auto-plan / timezone | Yes | No | No |
| API tokens (own) | Yes | Yes | Yes (own tokens; Viewer role on the token still cannot mutate) |
| Project webhooks | Yes | Yes | No |
| My Work | Yes | Yes | No |
| Search, charts, read board | Yes | Yes | Yes |
| CSV export/import | Yes | No (assumption) | No |

Organisation owners can do project-owner things on every project in the organisation.

The requester is not a special ACL. Copy may say `You requested this` on My Work Delivered rows. The Accept button is the same Bloom button.

---

## 11. Manifesto notes (tensions)

**Density of four columns vs the occasional user.** The board is dense on purpose (Tracker). Occasional users still get one Bloom verb per row and empty-state sentences. They do not get a wizard. `?` and skiplinks (`Shift+A`) are how they learn location. We do not add a “simple mode”.

**Icebox must not look like a fourth Kanban column.** Same paper column class as the others (guide). Difference is **chrome and words**, not a new colour: no band header, no `points / capacity`, empty copy is “waiting / capture / pull”, the row verb is **Pull to Backlog** not “Move to next”. Current is the only Bloom-highlighted column. Dragging toward Current does not advance a pipeline — it snaps back. Hiding Icebox (`Shift+I`) is allowed so it is not always sitting there as “step 0”.

**Icebox is leftmost.** The product names columns Icebox, Backlog, Current, Done. Tracker often parked Icebox to the right. Leftmost + holding-pen copy is a real tension (it reads as a funnel). Mitigation is the verb language and the snap-back, not a reorder of the columns. Reordering the four columns is **not** in this spec.

**Enter means the verb. O opens.** Tracker used click for Accept and Enter to open. Flower uses Enter to activate the row’s StateButton when collapsed so Accept / Restart have a keyboard path while `A` stays Add Story. Opening the sheet is `O` or click-the-title. Locked 20 Aug 2026.

**Tab selects stories (Tracker) vs native Tab.** On the board, Tab is story focus so a morning works one-handed. In any dialog or field, Tab is native. Contrast and focus ring are mandatory so this is learnable.

**Accept is team-wide, requester should.** The UI does not nag Luis off the Accept button (no Typeform “are you the requester?”). My Work puts Maya’s Delivered at her fingertips. Undo fixes a polite mistake.

**Computed bands.** Length is a number of days. Current / Backlog headers are live bands from velocity + estimates. Done is a flat aged-out list. After the first completed window the Current denominator is capacity, not velocity.

**Optimistic verbs vs “wait is a bug”.** The row updates before the server. Failure copy snaps back. Do not put a spinner on Start.

**Charts vs the board.** Charts are a separate quiet screen (slice 21 / 27). They are not dashboard widgets on the board (Monday refuse).

**Epic colour.** Guide has no purple. Epic uses Bloom-as-label, not a new hex. Selection ring stays the ring, not a fill.

---

## 12. Locked 20 Aug 2026

Do not reopen these. Product Owner accepted them. They live in `docs/product/open-questions.md` as well.

1. **`A` is Add Story.** Not Accept.
2. **Accept** is Enter on a focused `delivered` row (or the Accept button).
3. **`R` is Reject** (reason, then confirm).
4. **Restart** is Enter on a focused `rejected` row.
5. **Username** is inferred from the email local-part at signup. No username field in slice 0.
6. **`Shift+H` is reserved.** No project history panel in MVP. Do not bind it to story activity or a fake feed.
7. Estimate keys in MVP are **0 / 1 / 2 / 3** only. Fibonacci / custom keys wait for slice 28.
8. **Finished** has no status colour. White row + Deliver verb. Do not invent yellow.
9. **`O` opens** the story sheet. Enter stays the primary verb. Esc closes. Click the title also opens.
10. Epic visual: no new purple hex. Epic pill is Bloom-bordered.
11. Slice 30 manual Current: drag Backlog → Current is legal, and **C** (unshifted) moves the focused unstarted Backlog story into Current. Icebox → Current stays illegal.

Also law: Icebox is its own order; organisation owners create projects; CSV is Owner-only; cycle clock is first start → accepted; workspaces are personal; `planned` only if slice 30 needs it; project timezone stored, default `Australia/Melbourne`, not asked at signup; API tokens are generic (same API, same roles; a Member token can accept).

---

## 13. QA glance (UI)

A reviewer can fail a build from this list without a meeting.

- Palette is bloom / stem / paper / pollen / the four status hexes. No new brand colour.
- Fraunces headings, Inter body, Lucide stroke 2.
- Columns named Icebox, Backlog, Current, Done. Current header Bloom-highlighted. Icebox has no band header.
- Empty Icebox / Backlog / Current / Done use the slice 0 sentences.
- `StoryRow` shows Type icon, EstimateChip (Feature), title, StateButton. StateButton is not hover-only and not inside More.
- Unestimated Feature: Start refused with `Estimate this feature before you start it. 0 is fine.`
- Accept stays in Current. Done empty until that work ages past the current window. No window groups.
- Reject requires reason; Restart returns to started in Current.
- Drag unstarted above started snaps back with the slice 7 sentence.
- Drag into Current or Done snaps back. Start is the overflow valve.
- Current can read `6 / 10` on cold start, or `{points} / {capacity}` + `Over capacity` after the first completed window.
- `A` adds to Icebox. `0`–`3` estimate. `S` `F` `D`. `Shift+I/B/C/D` toggle columns. `/` search. `?` map. `Esc` back to the board.
- Viewer never sees a mutating Primary.
- Confirm on Reject / Delete / Revoke; not on Accept / Start.
- Contrast: ink on paper, white on Bloom, ink on Pollen chips. Focus ring visible.
- Copy uses Icebox, Current, window, band, velocity, capacity, points, story, epic. Length is in days.

---

## References

- [`docs/reference/frontend-design-guide.md`](../reference/frontend-design-guide.md) — look, components, type, colour.
- [`docs/core-workflow/product-spec.md`](./product-spec.md) — slices, machines, AC.
- [`docs/product/tracker-brief.md`](../product/tracker-brief.md) — copy-exactly vs modernise.
- [`docs/velocity-planning/velocity-and-planning.md`](../velocity-planning/velocity-and-planning.md) — pack, leave-short, dates, charts, panel membership.
- [`docs/multitenancy/multitenancy.md`](../multitenancy/multitenancy.md) — roles, isolation, settings who.
- [`docs/product/open-questions.md`](../product/open-questions.md) — true forks.
- Classic Tracker help: keyboard shortcuts, story panels, prioritizing with the keyboard, backlog-to-current, My Work.
