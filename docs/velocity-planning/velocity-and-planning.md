# Velocity and planning

This file is the planning product. If a slice, a mock, or an implementation disagrees with it, this file wins until Dan changes it.

Flower plans like classic Pivotal Tracker. We copy the model. We do not invent a third one. Deltas are only those listed in [tracker-brief.md](../product/tracker-brief.md) (initial-velocity behaviour is Tracker; team strength is *deferred*, not replaced with a fiction).

## Definitions

| Term | Meaning |
| --- | --- |
| **Story** | A row. Types: feature, bug, chore, release. |
| **Estimate** | A point value from the project scale, or none. `0` is an estimate. Missing is unestimated (`estimate:-1` in search). |
| **Priority** | Position in the **ranked** list (Backlog + Current + accepted-this-iteration). Icebox has its own order and is not in this list. |
| **Iteration** | A contiguous time box. Default 1 week. Existing table: `iterations(number, starts_on, ends_on)`. |
| **Icebox** | `unscheduled` holding pen. Not in velocity math. Not a pipeline stage. |
| **Current** | This iteration: in-progress (`started` / `finished` / `delivered` / `rejected`) + velocity-filled `unstarted` + **accepted this iteration**. |
| **Backlog** | Future iterations of the same ranked list (stories the planner did not put in Current). |
| **Done** | Accepted stories from **completed** iterations only. |
| **Velocity** | Rolling average of accepted **Feature** points over the last K completed iterations, or **initial 10** if none have completed. |
| **Projected delivery date** | End date of the iteration the planner assigned. Not editable. Not a promise. |

## Iteration clock

- Default length: **1 week**. Owner may set **1, 2, 3, or 4** weeks (`projects.iteration_length_weeks`).
- Rollover at **midnight** at the start of the next iteration, **project timezone**.
- Iterations start at date `starts_on` (existing column). Default start weekday: Monday. First iteration starts on the Monday on or before project-created (or the configured start weekday).
- `ends_on` is the last calendar day of the box. Projected dates display `ends_on`.
- Changing length / timezone / start weekday replans immediately. Completed iterations keep historical `starts_on` / `ends_on` and recorded accepted Feature points.
- Project timezone is a setting we still need (not in `000001`). Until the fork is resolved: store it; default `Australia/Melbourne`.

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

- Estimate at **accept** time is what is recorded. Edits after accept do not rewrite the completed iteration.
- Rejected Features add nothing. If later accepted, they add to the iteration of that accept.
- Accepted-this-iteration Features sit in Current and count toward **this** iteration’s running total; they enter the velocity average only when the iteration **completes**.
- Deleted Features that were already in a completed iteration’s recorded total stay in that total (history is stable).
- Team strength % is **not MVP**. Do not implement a hidden 100% strength. The MVP formula is a plain average of accepted Feature points per completed iteration, not Tracker’s normalised-per-week formula.

## Velocity formula (MVP)

Let `K` = velocity strategy, default **3**, allowed **1–4** (project setting, not in `000001` yet).

Let `N` = number of completed iterations.

```
if N == 0:
    V = initial_velocity        # default 10, Owner-editable
else:
    V = floor( mean( accepted_feature_points of the last min(K, N) completed iterations ) )
```

- Current (in progress) is not in the average.
- After a streak of K completed iterations with **zero** accepted Features, Tracker reverts to initial velocity. **Copy that:** if the last K completed iterations are all 0, `V = initial_velocity` again.
- Display V as an integer. `floor` of the mean (Tracker rounds down after its fuller formula; we do not invent half-up).

Initial velocity default **10**. Owner may change it before the first completion. After `N ≥ 1`, the calculated average is what the board uses (until the all-zero revert).

## Cold start

Not empty, not fake.

- V = 10 (unless the Owner set a different initial).
- Auto-plan **does** fill Current from the ranked list up to V, leaving Current **short** rather than overfilling.
- Icebox stays out of the fill.
- Future Backlog iterations also pack at V = 10, so dates exist from day one if Features are estimated and ranked. That is Tracker. We do **not** hide dates until the first accept.
- Copy on an empty Current with an empty ranked list: “Nothing in this iteration yet. Pull a story from Icebox or create one.”
- Do not show `0` velocity on a new project unless the Owner set initial velocity to 0.

## How the planner fills Current and future iterations

One ranked list L: every non-accepted, non-icebox story, plus accepted-this-iteration stories (they occupy Current but are not “packed” as upcoming work). Highest priority first.

### Cost

```
cost(story) =
  estimate  if type == feature and estimate is not null
  estimate  if bugs_and_chores_may_be_given_points
            and type in {bug, chore} and estimate is not null
  0         otherwise   # bugs, chores, releases, unestimated
```

Unestimated Features still occupy a place. They cost 0. They **cannot be Started**. They can be packed into Current while Current feature-points have not yet reached V (Tracker: unestimated bugs/chores keep filling Current until points *exceed* / have not yet exceeded V).

Normative for zero-cost rows: they may be assigned to the iteration being filled **only while that iteration’s packed Feature-points `< V`**. Once the iteration has reached V (equal or over via Start overflow), remaining zero-cost unstarted stories overflow with the rest. In-flight zero-cost stories (a started Bug) always stay in Current.

### In-flight and accepted-this-iteration

Always in Current, regardless of cost:

- `started`, `finished`, `delivered`, `rejected`
- `accepted` with `accepted_at` inside the current iteration window

These can make Current **over** velocity. Auto-plan never kicks them out.

### Auto-plan algorithm (default)

1. Place every in-flight and accepted-this-iteration story into Current, preserving relative order.
2. `remaining = V - feature_points(Current)`. If `remaining <= 0`, do not auto-fill more estimated Features into Current. (Zero-cost unstarted also stay out once points ≥ V.)
3. Walk the remaining `unstarted` / Release-`started` stories in L:
   - If estimated Feature `cost > remaining`, **do not place it in this iteration**. Leave the iteration **short**. Open the next iteration with `remaining = V` and try again.
   - If `cost <= remaining`, place it here and subtract.
   - Zero-cost story: place here only if this iteration’s Feature-points `< V`; else overflow.
   - Release markers walk with the list (cost 0) and sit where they fall — which, if the team placed them at the **end** of the milestone, is after that milestone’s last story.
4. Repeat until L is exhausted.

### Fit rules (normative)

- **Never split a story.**
- **Leave Current (and each future iteration) short** rather than overfill with the next estimated Feature.
- **Start is the overflow valve.** Start on a Backlog or Icebox story jumps it to Current and **may** push Current over V. Icebox Start also unschedules→starts (Feature still needs an estimate).
- **Oversized Feature (`cost > V`):** it never auto-fills into an iteration that already has Feature-points. It becomes the first (and only estimated) occupant of the next empty iteration — and that iteration is still **short of a second story** but **over** V after it is placed? Conflict.

  Tracker: if a story is larger than velocity it still has to live somewhere. Copy: an estimated Feature with `cost > V` **is** placed into the next iteration that has **no** estimated Feature yet, and that iteration is marked **over velocity**. This is the one auto-plan exception to “leave short rather than overfill,” because the alternative is a story that never packs. Current: if Current already has Feature-points > 0, the oversized story does **not** enter Current via auto-plan (Start can still pull it). If Current has 0 Feature-points, the oversized story **does** enter Current and Current is over V.

- **Unstarted cannot sit above started** in the ranked list (drag rejected). The planner does not violate this: in-flight stories were already started and stay higher if the user left them there; auto-fill only appends unstarted *below* in-flight in Current when packing from Backlog. If the user has unstarted above in-flight, that is already illegal in the UI; the server rejects that rank write.

- **Accepted this iteration** stay in Current until rollover, even if that makes Current look “done.”

### Manual planning (Phase 3)

Off by default. When off for Current: Current contains only in-flight, accepted-this-iteration, and stories **explicitly moved** there. No velocity fill into Current. Backlog future iterations still auto-plan at V. Restoring auto-plan reruns the algorithm.

## Re-plan is live

Recompute on: create, delete, icebox, schedule, reorder, estimate change, type change, start (and other verbs), accept, reject, restart, undo, rollover, V change, length / TZ / strategy / initial-velocity / bugs-and-chores toggle change.

No Recalculate button. The board after a mutation is already right.

## Iteration rollover

At midnight at `ends_on + 1 day 00:00` project TZ:

1. Freeze that iteration’s accepted Feature points.
2. Mark it completed. It now participates in V.
3. Recompute V.
4. Every story accepted in that window **moves to Done**, grouped under that iteration.
5. In-flight and leftover unstarted Current stories remain; unstarted are re-packed (they are not sticky to “this week” unless still in-flight).
6. Unaccepted work is not failed and not iceboxed.

## Projected delivery date

```
projected_date(story) =
  iterations.ends_on of the assigned iteration
  or null if Icebox / no iteration
```

Display the date only, project TZ. Icebox: “Not scheduled.”

Not a field. Not a due date. No “date slipped” email in MVP.

## Releases

Copy Tracker.

- Marker. Cost 0. Not work.
- Place at the **end** of the milestone’s stories (work **above**, marker **below**).
- Auto-`started` when created in Backlog or dragged from Icebox into the ranked list.
- Finish → `accepted`. The marker may sit accepted in Current until rollover, then Done, like other accepted stories — or, if Finish happens when it is in a future iteration, it accepts in place. Prefer: Finish is allowed when the team says the milestone shipped; the marker becomes `accepted` and, if that is in the current window, stays in Current until rollover.
- Optional **target date**. This is the only date a human types. It does **not** move stories and does not override the plan.
- Colour (when a target date is set):
  - **Blue (on track):** `starts_on` of the iteration that **contains the marker** ≤ target date.
  - **Red (late):** that iteration’s `starts_on` **>** target date.
- If the marker has no iteration (Icebox): no colour, or muted “unscheduled.”
- No stories above the marker: still a marker; date/colour use the iteration it sits in (often Current leftover / empty pack). Copy may say “No stories above this release.”
- Two releases: each marker owns the work **above it and below the previous marker**. Colour/date still use the iteration that **contains that marker**.

## Charts

### Velocity (Phase 2)

- X: completed iterations (`ends_on`).
- Y: accepted Feature points.
- Bar per completed iteration. Line at current V (show “initial 10” while N = 0).
- Current is not a completed bar. Optional faint “accepted so far this iteration.”
- Empty: “No completed iterations yet.” (V may still be 10.)

### Burn-up (Phase 2)

- Accepted line: cumulative accepted Feature points.
- Scope line: snapshot sum of Feature estimates (unestimated = 0, Icebox excluded) at each rollover, plus live “now.”
- Bugs/chores/releases add 0 unless the Phase 3 toggle is on.
- Empty: “No scope yet.”

### Cycle time (Phase 3)

- First `started_at` → `accepted_at`. Reject / Restart does not reset the clock.
- p50 / p75 / p95. Filter by type. No per-person chart.

## Panel membership

| Panel | Contains |
| --- | --- |
| Icebox | `unscheduled` only. Own order. No iteration headers. |
| Current | In-flight + accepted-this-iteration + auto-filled `unstarted` (and Releases that packed here). Header: points / V, `ends_on`, over-velocity badge if Feature-points > V. |
| Backlog | The rest of the ranked list, under future iteration headers when they have members. |
| Done | Accepted in completed iterations, newest iteration first. |

## Worked example 1 — new project, initial velocity 10

Project created Thursday 6 Aug 2026. TZ `Australia/Melbourne`. Length 1 week, starts Monday. Current box: Mon 3 Aug – Sun 9 Aug. N = 0. V = **10**.

Ranked list (all Features, all estimated, all `unstarted`):

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
- Next iteration (10–16 Aug), remaining 10: Billing 5, Admin 3 → both fit (2 left).

Board:

- Current (ends 9 Aug): Login, Reset, Search — **6 / 10**. Projected dates 9 Aug.
- Backlog iteration ending 16 Aug: Billing, Admin. Dates 16 Aug.
- Icebox empty.
- Done empty.

Luis Starts Billing from Backlog. Billing jumps to Current. Current Feature-points = 6 + 5 = **11 / 10**, over-velocity badge. Admin stays in the 16 Aug iteration.

They Accept Login (2) on Friday. Login stays in **Current** (accepted-this-iteration). Still not Done.

Monday 10 Aug 00:00 Melbourne: rollover.

- Iteration 1 accepted Feature points = 2 (Login only).
- N = 1. V = **2**.
- Login → Done under “Ended 9 Aug.”
- In-flight: Billing started (5). Current already over the new V. No auto-fill. Reset and Search if still unstarted are re-packed: remaining = 2 − 5 < 0, so they go to future iterations at V = 2.

## Worked example 2 — steady V, leave-short, release colour

Last three completed iterations accepted Feature points 8, 7, 9. V = floor((8+7+9)/3) = floor(8) = **8**.

Current ends **30 Aug 2026** (started **24 Aug**).

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
- Next iteration starts **31 Aug**, ends **6 Sep**, remaining 8.
- Checkout 5 → that iteration. remaining 3.
- Admin 3 ≤ 3 → that iteration. remaining 0.
- v1.0 (0) → same iteration (end of those stories).
- SSO 5 > 0 → iteration starting **7 Sep**.

v1.0 lives in the iteration that **starts 31 Aug**. Target 31 Aug. `starts_on` ≤ target → **blue**.

If Maya reorders Checkout below v1.0, the marker sits earlier. Suppose marker then packs into Current (24 Aug start) still blue. If instead the work above the marker grows so the marker packs into the iteration starting **7 Sep**, `starts_on` 7 Sep > 31 Aug → **red**. Live.

SSO is below the marker: not in v1.0.

## Worked example 3 — reject stays in Current; accept after rollover

Checkout is Delivered in Current after Luis Started it (overflow). Maya Rejects: “tax is wrong.” State `rejected`, still Current. Restart → `started`, still Current. Velocity unchanged.

They Accept Checkout on **2 Sep** (next iteration). Points land in the iteration that contains 2 Sep, not the one they first started it in.

## What we will not do

- Typed velocity override (“pretend we do 20”) except the **initial velocity** setting, which stops mattering after N ≥ 1 (and the all-zero revert).
- Drag-to-pin a story to a date. Start or reorder; do not staple to 6 Sep.
- Split stories.
- Count bugs/chores as Features.
- Show Icebox stories inside iteration headers.
- Per-person velocity.
- Team strength % in MVP.
- A date picker that overrides the plan. Release target date is a comparison only.

## QA short script (slices 8 + 20)

1. New project. V displays 10. Create five estimated Features in Backlog totalling more than 10. Confirm Current is **short** of 10 rather than over.
2. Start a Backlog Feature that did not fit. Confirm it is in Current and Current may exceed 10.
3. Accept one Feature. Confirm it is still in Current, not Done.
4. Advance the test clock past midnight rollover. Confirm that Feature is in Done and V is that iteration’s Feature points.
5. Place a Release at the end of a group. Set a target date. Confirm blue/red vs the iteration **start** that contains the marker.
6. Drag unstarted above started. Must fail.
7. Confirm Icebox stories never appear in Current via auto-plan.
