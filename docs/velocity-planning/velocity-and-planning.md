# Velocity and planning

Change name: `flower`

This file is the planning product. If a slice, a mock, or an implementation disagrees with it, this file wins.

Flower copies Tracker's *idea* of velocity windows: one ranked list, pack toward velocity, leave a band short rather than overfill. Flower does **not** persist Tracker-style iteration records as the plan. There is no Iteration aggregate. There is no `iterations` table. Stories are not assigned to a window. Project setting: `iteration_length_days`. Accepted-points history lives in `velocity_samples`.

Deltas versus Tracker are only those listed in [tracker-brief.md](./tracker-brief.md) (initial-velocity behaviour is Tracker; team strength is *deferred*, not replaced with a fiction).

## Definitions

| Term | Meaning |
| --- | --- |
| **Story** | A row. Types: feature, bug, chore, release. |
| **Estimate** | A point value from the project scale, or none. `0` is an estimate. Missing is unestimated (`estimate:-1` in search). |
| **Priority** | Position in the **ranked** list (Backlog + Current + accepted-this-window). Icebox has its own order and is not in this list. |
| **Length in days** | The only stored window size on the project. Default **7**. A positive integer the Owner may change. Not weeks. Not a list of rows. |
| **Window (computed)** | A contiguous date range of length L, derived from start weekday + L + now + timezone. Never a stored entity. Never a story assignment. |
| **Band** | A visual group the UI draws by running the pack function: Current is the open window; later bands sit in Backlog. |
| **Icebox** | `unscheduled` holding pen. Not in velocity math. Not a pipeline stage. |
| **Current** | The open window: in-progress (`started` / `finished` / `delivered` / `rejected`) + velocity-filled `unstarted` + **accepted this window**. Head of the ranked list that fits this window's velocity. |
| **Backlog** | Later computed bands of the same ranked list (stories the pack function did not put in Current). Visual only. |
| **Done** | Accepted stories whose `accepted_at` has **aged past the current window**. Flat list, newest accepted first. Not grouped by a window row. |
| **Velocity** | Rolling average of accepted **Feature** points over the last K **completed windows**, or **initial 10** if none have completed. |
| **Projected delivery date** | End date of the computed window the story packed into. Not editable. Not a promise. Not a stored assignment. |

## Window clock (chosen rule)

Pick is **calendar windows that end at midnight**, project timezone. Not “now + remaining days.”

- Stored: **iteration length in days** `L`. Default **7**. Owner may set any positive integer (typical 7, 14, 21, 28).
- Stored: start weekday (default **Monday**) and project timezone (default `Australia/Melbourne` until the timezone fork is resolved — see [open-questions.md](./open-questions.md)).
- First window **start date** = the configured start weekday **on or before** the project-created date, in the project timezone.
- Window `i` (i = 0, 1, 2, …) is the half-open range:
  - `starts_on(i) = first_start + i × L days` at `00:00` project TZ
  - `ends_at(i) = starts_on(i) + L days` at `00:00` project TZ (midnight at the start of the next window)
  - Last calendar day of the window = `ends_at(i) − 1 day`. Display as `Ends {that date}`.
- **Current window** = the unique `i` where `starts_on(i) ≤ now < ends_at(i)` in project TZ.
- When `now` crosses `ends_at` of the open window, that window is **completed**. Accepted work whose `accepted_at` fell in it ages into Done. Velocity recomputes. Remaining unstarted stories are re-packed. In-flight stay in Current. Unaccepted work is not failed and not iceboxed.
- Changing length / timezone / start weekday replans immediately. Completed windows stay historically addressable as date ranges (from the same function, or from `velocity_samples`). They are not iteration rows and they do not assign stories.

## What counts toward velocity

**Accepted Feature points only.**

| Type | Counts toward velocity? |
| --- | --- |
| Feature, estimated, accepted in the window | Yes — estimate at accept time |
| Feature, 0 points, accepted | Yes, adds 0 |
| Bug | No (unless Phase 3 toggle is on *and* it has points) |
| Chore | No (same toggle) |
| Release | Never |
| Unestimated anything | Never |
| Icebox | Never (cannot be accepted from Icebox) |

Further:

- Estimate at **accept** time is what is recorded. Edits after accept do not rewrite a completed window.
- Rejected Features add nothing. If later accepted, they add to the window that contains that accept.
- Accepted-this-window Features sit in Current and count toward **this** window’s running total; they enter the velocity average only when the window **completes**.
- Deleting a Feature after a window completed does not rewrite that window’s total. Persist **`velocity_samples`** (window start + accepted Feature points) so history is stable. Samples are not an Iteration and do not assign stories.
- Team strength % is **not MVP**. Do not implement a hidden 100% strength. The MVP formula is a plain average of accepted Feature points per completed window, not Tracker’s normalised-per-week formula.

## Velocity formula (MVP)

Let `K` = velocity strategy, default **3**, allowed **1–4** (project setting).

Let `N` = number of completed windows.

```
if N == 0:
    V = initial_velocity        # default 10, Owner-editable
else:
    V = floor( mean( accepted_feature_points of the last min(K, N) completed windows ) )
```

- Current (in progress) is not in the average.
- After a streak of K completed windows with **zero** accepted Features, Tracker reverts to initial velocity. **Copy that:** if the last K completed windows are all 0, `V = initial_velocity` again.
- Display V as an integer. `floor` of the mean (Tracker rounds down after its fuller formula; we do not invent half-up).

Initial velocity default **10**. Owner may change it before the first completion. After `N ≥ 1`, the calculated average is what the board uses (until the all-zero revert).

## Cold start

Not empty, not fake.

- V = 10 (unless the Owner set a different initial).
- Auto-plan **does** fill Current from the ranked list up to V, leaving Current **short** rather than overfilling.
- Icebox stays out of the fill.
- Future Backlog bands also pack at V = 10, so dates exist from day one if Features are estimated and ranked. That is Tracker. We do **not** hide dates until the first accept.
- Copy on an empty Current with an empty ranked list: “Nothing in this window yet. Pull a story from Icebox or create one.”
- Do not show `0` velocity on a new project unless the Owner set initial velocity to 0.

## Packing is a pure function

No Iteration entity. Stories stay on one ranked list. The UI draws bands from the function’s output.

```
pack(
  ordered_stories,   # ranked list, Icebox excluded; highest priority first
  velocity,          # V, integer
  length_in_days,    # L
  now,
  timezone
) → bands            # band 0 = Current (open window); later bands = Backlog chrome
```

Each band carries a computed `starts_on` / `ends_at` from the window clock, then the stories that packed into it. Re-run whenever velocity, rank, estimate, accept, type, start (or other verbs), length, timezone, start weekday, strategy, initial velocity, or (later) the bugs-and-chores toggle changes. No Recalculate button.

### Cost

```
cost(story) =
  estimate  if type == feature and estimate is not null
  estimate  if bugs_and_chores_may_be_given_points
            and type in {bug, chore} and estimate is not null
  0         otherwise   # bugs, chores, releases, unestimated
```

Unestimated Features still occupy a place. They cost 0. They **cannot be Started**. They can be packed into Current while Current feature-points have not yet reached V (Tracker: unestimated bugs/chores keep filling Current until points *exceed* / have not yet exceeded V).

Normative for zero-cost rows: they may sit in the band being filled **only while that band’s packed Feature-points `< V`**. Once the band has reached V (equal or over via Start overflow), remaining zero-cost unstarted stories overflow with the rest. In-flight zero-cost stories (a started Bug) always stay in Current.

### In-flight and accepted-this-window

Always in Current, regardless of cost:

- `started`, `finished`, `delivered`, `rejected`
- `accepted` with `accepted_at` inside the current window (`starts_on ≤ accepted_at < ends_at`)

These can make Current **over** velocity. Auto-plan never kicks them out. They stay in Current until the current window ends.

### Auto-plan algorithm (default)

1. Place every in-flight and accepted-this-window story into Current (band 0), preserving relative order.
2. `remaining = V - feature_points(Current)`. If `remaining <= 0`, do not auto-fill more estimated Features into Current. (Zero-cost unstarted also stay out once points ≥ V.)
3. Walk the remaining `unstarted` / Release-`started` stories in rank order:
   - If estimated Feature `cost > remaining`, **do not place it in this band**. Leave the band **short**. Open the next band with `remaining = V` and try again.
   - If `cost <= remaining`, place it here and subtract.
   - Zero-cost story: place here only if this band’s Feature-points `< V`; else overflow.
   - Release markers walk with the list (cost 0) and sit where they fall — which, if the team placed them at the **end** of the milestone, is after that milestone’s last story.
4. Repeat until the ranked list is exhausted.

Future bands are **visual**. They are not rows. A story in a future band is still just a ranked `unstarted` (or a Release marker). Changing rank or estimate moves it by re-running `pack`.

### Fit rules (normative)

- **Never split a story.**
- **Leave Current (and each future band) short** rather than overfill with the next estimated Feature.
- **Start is the overflow valve.** Start on a Backlog or Icebox story jumps it to Current and **may** push Current over V. Icebox Start also unscheduled → started (Icebox start; Feature still needs an estimate).
- **Oversized Feature (`cost > V`):** it never auto-fills into a band that already has Feature-points. It becomes the first (and only estimated) occupant of the next empty band — and that band is still **short of a second story** but **over** V after it is placed.

  Tracker: if a story is larger than velocity it still has to live somewhere. Copy: an estimated Feature with `cost > V` **is** placed into the next band that has **no** estimated Feature yet, and that band is marked **over velocity**. This is the one auto-plan exception to “leave short rather than overfill,” because the alternative is a story that never packs. Current: if Current already has Feature-points > 0, the oversized story does **not** enter Current via auto-plan (Start can still pull it). If Current has 0 Feature-points, the oversized story **does** enter Current and Current is over V.

- **Unstarted cannot sit above started** in the ranked list (drag rejected). The pack function does not violate this: in-flight stories were already started and stay higher if the user left them there; auto-fill only appends unstarted *below* in-flight in Current when packing from Backlog. If the user has unstarted above in-flight, that is already illegal in the UI; the server rejects that rank write.

- **Accepted this window** stay in Current until the window ends, even if that makes Current look “done.”

### Manual planning (Phase 3)

Off by default. When off for Current: Current contains only in-flight, accepted-this-window, and stories **explicitly moved** there. No velocity fill into Current. Backlog future bands still auto-plan at V. Restoring auto-plan reruns `pack`.

## Re-plan is live

Recompute `pack` on: create, delete, icebox, schedule, reorder, estimate change, type change, start (and other verbs), accept, reject, restart, undo, window end, V change, length / TZ / strategy / initial-velocity / bugs-and-chores toggle change.

No Recalculate button. The board after a mutation is already right.

## Window end

At `ends_at` of the open window (midnight, project TZ):

1. That window’s accepted Feature points become a `velocity_samples` row (sum of Feature estimates accepted while it was current).
2. It now participates in V.
3. Recompute V.
4. Every story accepted in that window **ages into Done** as a **flat** list (newest accepted first). No group-by-window chrome.
5. In-flight and remaining unstarted Current stories remain; unstarted are re-packed (they are not sticky to “this window” unless still in-flight).
6. Unaccepted work is not failed and not iceboxed.

There is no button to close a window. There is no stored row to close.

## Projected delivery date

```
projected_date(story) =
  last calendar day of the computed band pack() assigned
  or null if Icebox / not on the ranked list
```

Display the date only, project TZ. Icebox: “Not scheduled.”

Not a field. Not a due date. No “date slipped” email in MVP.

## Releases

Copy Tracker colour *idea*; compare against the **computed** window that contains the marker.

- Marker. Cost 0. Not work.
- Place at the **end** of the milestone’s stories (work **above**, marker **below**).
- Auto-`started` when created in Backlog or dragged from Icebox into the ranked list.
- Finish → `accepted`. The marker may sit accepted in Current until the window ends, then Done, like other accepted stories — or, if Finish happens when it is in a future band, it accepts in place. Prefer: Finish is allowed when the team says the milestone shipped; the marker becomes `accepted` and, if that is in the current window, stays in Current until the window ends.
- Optional **target date**. This is the only date a human types. It does **not** move stories and does not override the plan.
- Colour (when a target date is set):
  - **Blue (on track):** `starts_on` of the **computed window that contains the marker** ≤ target date.
  - **Red (late):** that computed window’s `starts_on` **>** target date.
- If the marker is in Icebox: no colour, or muted “unscheduled.”
- No stories above the marker: still a marker; date/colour use the computed window it packed into (often Current with an empty pack). Copy may say “No stories above this release.”
- Two releases: each marker owns the work **above it and below the previous marker**. Colour/date still use the computed window that **contains that marker**.

## Charts

### Velocity (Phase 2)

- X: completed windows (`ends_at`).
- Y: accepted Feature points.
- Bar per completed window. Line at current V (show “initial 10” while N = 0).
- Current is not a completed bar. Optional faint “accepted so far this window.”
- Empty: “No completed windows yet.” (V may still be 10.)

### Burn-up (Phase 2)

- Accepted line: cumulative accepted Feature points.
- Scope line: snapshot sum of Feature estimates (unestimated = 0, Icebox excluded) at each window end, plus live “now.”
- Bugs/chores/releases add 0 unless the Phase 3 toggle is on.
- Empty: “No scope yet.”

### Cycle time (Phase 3)

- First `started_at` → `accepted_at`. Reject / Restart does not reset the clock.
- p50 / p75 / p95. Filter by type. No per-person chart.

## Panel membership

| Panel | Contains |
| --- | --- |
| Icebox | `unscheduled` only. Own order. No band headers. |
| Current | In-flight + accepted-this-window + auto-filled `unstarted` (and Releases that packed here). Header: points / V, computed `Ends {date}`, over-velocity badge if Feature-points > V. |
| Backlog | The rest of the ranked list, under future **band** headers when they have members. No “iteration 3”. |
| Done | Accepted stories aged past the current window. Flat list, newest accepted first. |

## Worked example 1 — new project, initial velocity 10

Project created Thursday 6 Aug 2026. TZ `Australia/Melbourne`. Length **7 days**, starts Monday. First start = Monday 3 Aug. Current window: Mon 3 Aug 00:00 – Mon 10 Aug 00:00. Displayed end: **9 Aug**. N = 0. V = **10**.

`pack(ordered Features, V=10, L=7, now=Thu 6 Aug, TZ=Australia/Melbourne)`:

| Pos | Story | Est |
| --- | --- | --- |
| 1 | Login | 2 |
| 2 | Reset | 1 |
| 3 | Search | 3 |
| 4 | Billing | 5 |
| 5 | Admin | 3 |

Walk Current, remaining 10:

- Login 2 → Current (8 left)
- Reset 1 → Current (7)
- Search 3 → Current (4)
- Billing 5 > 4 → **leave Current short** (points 6 / 10). Billing does not enter Current.
- Next band (10–16 Aug), remaining 10: Billing 5, Admin 3 → both fit (2 left).

Board:

- Current (ends 9 Aug): Login, Reset, Search — **6 / 10**. Projected dates 9 Aug.
- Backlog band ending 16 Aug: Billing, Admin. Dates 16 Aug.
- Icebox empty.
- Done empty.

Luis Starts Billing from Backlog. Billing jumps to Current. Current Feature-points = 6 + 5 = **11 / 10**, over-velocity badge. Admin stays in the 16 Aug band.

They Accept Login (2) on Friday. Login stays in **Current** (accepted-this-window). Still not Done.

Monday 10 Aug 00:00 Melbourne: current window ends.

- Window 3–10 Aug accepted Feature points = 2 (Login only).
- N = 1. V = **2**.
- Login → Done (flat list; newest accepted first).
- In-flight: Billing started (5). Current already over the new V. No auto-fill. Reset and Search if still unstarted are re-packed: remaining = 2 − 5 < 0, so they go to future bands at V = 2.

## Worked example 2 — steady V, leave-short, release colour vs computed window

Last three completed windows accepted Feature points 8, 7, 9. V = floor((8+7+9)/3) = floor(8) = **8**.

Current window starts **24 Aug 2026**, ends at midnight **31 Aug** (displayed end **30 Aug**). L = 7.

| Pos | Story | Type | Est | State |
| --- | --- | --- | --- | --- |
| 1 | Redesign nav | Feature | 3 | started |
| 2 | Copy tweaks | Feature | 1 | unstarted |
| 3 | Checkout | Feature | 5 | unstarted |
| 4 | Admin | Feature | 3 | unstarted |
| 5 | v1.0 | Release | — | started (marker) |
| 6 | SSO | Feature | 5 | unstarted |

Target date on v1.0: **31 Aug 2026**.

Walk:

- Current seed: Redesign (3). remaining 5.
- Copy (1) ≤ 5 → Current. remaining 4.
- Checkout (5) > 4 → **leave Current short** (4 / 8). Checkout does not enter Current.
- Next band starts **31 Aug**, ends **6 Sep**, remaining 8.
- Checkout 5 → that band. remaining 3.
- Admin 3 ≤ 3 → that band. remaining 0.
- v1.0 (0) → same band (end of those stories).
- SSO 5 > 0 → band starting **7 Sep**.

v1.0 lives in the **computed window that starts 31 Aug**. Target 31 Aug. `starts_on` ≤ target → **blue**.

If Maya reorders Checkout below v1.0, the marker sits earlier. Suppose the marker then packs into Current (24 Aug start) still blue. If instead the work above the marker grows so the marker packs into the window starting **7 Sep**, `starts_on` 7 Sep > 31 Aug → **red**. Live. Colour is versus the computed window that contains the marker, not a stored row.

SSO is below the marker: not in v1.0.

## Worked example 3 — reject stays in Current; accept after window end

Checkout is Delivered in Current after Luis Started it (overflow). Maya Rejects: “tax is wrong.” State `rejected`, still Current. Restart → `started`, still Current. Velocity unchanged.

They Accept Checkout on **2 Sep** (next window). Points land in the computed window that contains 2 Sep, not the window they first started it in.

## What we will not do

- Persist an Iteration as the plan. There is no `iterations` table. Stories are not assigned to a window.
- Typed velocity override (“pretend we do 20”) except the **initial velocity** setting, which stops mattering after N ≥ 1 (and the all-zero revert).
- Drag-to-pin a story to a date. Start or reorder; do not staple to 6 Sep.
- Split stories.
- Count bugs/chores as Features.
- Show Icebox stories inside band headers.
- Per-person velocity.
- Team strength % in MVP.
- A date picker that overrides the plan. Release target date is a comparison only.
- Store length as 1–4 weeks. Length is **days**.

## QA short script (slices 8 + 20)

1. New project. V displays 10. Create five estimated Features in Backlog totalling more than 10. Confirm Current is **short** of 10 rather than over.
2. Start a Backlog Feature that did not fit. Confirm it is in Current and Current may exceed 10.
3. Accept one Feature. Confirm it is still in Current, not Done.
4. Advance the test clock past midnight at the window end. Confirm that Feature is in Done (flat list) and V is that window’s Feature points.
5. Place a Release at the end of a group. Set a target date. Confirm blue/red vs the **computed window start** that contains the marker.
6. Drag unstarted above started. Must fail.
7. Confirm Icebox stories never appear in Current via auto-plan.
8. Confirm no story is assigned to an iteration row. Changing length in days replans live.
