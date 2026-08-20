# Velocity and planning

Change name: `flower`

This file is the planning product. If a slice, a mock, or an implementation disagrees with it, this file wins.

Flower copies Tracker's *idea* of velocity windows: one ranked list, pack toward velocity, leave a band short rather than overfill. Flower does **not** persist Tracker-style iteration records as the plan. There is no Iteration aggregate. There is no `iterations` table. Stories are not assigned to a window. The only stored window size is `iteration_length_days`.

Velocity is **live**. It is calculated from previous stories' start and end date/times (`started_at` → `accepted_at`). It is not written to persistence. We accept **stories**, not points.

Stories projected to complete in this window = that live velocity **plus** a prediction of how long each incomplete story will take, from previous **completed** stories of **similar estimated sizing** (the same estimate value).

Deltas versus Tracker are only those listed in [tracker-brief.md](./tracker-brief.md) (initial-velocity bootstrap is Tracker-shaped; team strength is *deferred*).

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
| **Current** | The open window: in-progress (`started` / `finished` / `delivered` / `rejected`) + duration-packed `unstarted` + **accepted this window**. Head of the ranked list that fits remaining time in this window. |
| **Backlog** | Later computed bands of the same ranked list (stories the pack function did not put in Current). Visual only. |
| **Done** | Accepted stories whose `accepted_at` has **aged past the current window**. Flat list, newest accepted first. Not grouped by a window row. |
| **Corpus** | Completed Features (and, if the bugs-and-chores-points toggle is on, estimated bugs/chores) that feed history. See Corpus. |
| **Velocity (`V_rate`)** | Live rate: corpus estimate-points per unit time. Undefined when the corpus is empty or `time == 0`. |
| **Predicted duration `pred(E)`** | How long an incomplete story of estimate `E` is expected to take, from completed corpus stories of the same estimate. |
| **Projected delivery date** | End date of the computed window the story packed into. Not editable. Not a promise. Not a stored assignment. |

## Window clock (chosen rule)

Pick is **calendar windows that end at midnight**, project timezone. The clock defines windows. Packing the **open** window uses time remaining from `now` to `ends_at`.

- Stored: **iteration length in days** `L`. Default **7**. Owner may set any positive integer (typical 7, 14, 21, 28).
- Stored: start weekday (default **Monday**) and project timezone (default `Australia/Melbourne` — see [open-questions.md](./open-questions.md) for the timezone-source fork).
- First window **start date** = the configured start weekday **on or before** the project-created date, in the project timezone.
- Window `i` (i = 0, 1, 2, …) is the half-open range:
  - `starts_on(i) = first_start + i × L days` at `00:00` project TZ
  - `ends_at(i) = starts_on(i) + L days` at `00:00` project TZ (midnight at the start of the next window)
  - Last calendar day of the window = `ends_at(i) − 1 day`. Display as `Ends {that date}`.
- **Current window** = the unique `i` where `starts_on(i) ≤ now < ends_at(i)` in project TZ.
- When `now` crosses `ends_at` of the open window, that window is **completed**. Accepted work whose `accepted_at` fell in it ages into Done. `V_rate` and `pred` recompute from stories. Remaining unstarted stories are re-packed. In-flight stay in Current. Unaccepted work is not failed and not iceboxed.
- Changing length / timezone / start weekday replans immediately. Completed windows stay historically addressable as date ranges from the same function. They are not iteration rows and they do not assign stories.

## Corpus

A story feeds history only when **all** of these hold:

- type is **Feature** (bugs/chores/releases are out unless the later bugs-and-chores-points toggle is on — then bugs/chores **with an estimate** join the corpus; releases never do)
- state is `accepted`
- both `started_at` and `accepted_at` are set

Duration:

```
D = accepted_at − started_at
```

Reject / Restart does **not** reset `started_at`. The clock is first start.

Size key: the story's **current** estimate (live). Edits to estimate move the story's `D` into the new size bucket on the next recompute.

Lookback: corpus stories whose `accepted_at` falls in the last `K` **completed** computed windows. `K` = velocity strategy, default **3**, allowed **1–4** (project setting). If fewer completed windows exist, use what exists. Stories accepted in the **open** window are not in the lookback.

Nothing is stored except story timestamps, the story's estimate, and project settings (`iteration_length_days`, start weekday, timezone, `velocity_strategy`, `initial_velocity`, and later the bugs-and-chores-points toggle).

Unestimated Features are not in the corpus. They do not contribute to velocity.

Icebox is never in the corpus (a story cannot be accepted from Icebox).

## Velocity (rate)

```
work = sum(estimate) of corpus stories
time = sum(D) of corpus stories

if time == 0 or corpus is empty:
    V_rate is undefined          # cold start
else:
    V_rate = work / time         # estimate-points per unit time
```

Compute `D` from timestamps in a consistent unit (fractional days). `V_rate` is then estimate-points per day. `L` and time remaining use the same unit.

Window capacity in estimate-points:

```
capacity           = V_rate × L                         # a full window of L days
capacity_remaining = V_rate × time_remaining_in_window  # open window: now → ends_at
```

Do not kick out in-flight or accepted-this-window stories, even when they make Current exceed capacity.

`V_rate` is a calculation. It is not persisted. Deleting a corpus story, changing its estimate, or editing timestamps (via undo) recomputes `V_rate` and `pred` from the stories that remain.

Team strength % is **not MVP**. Do not implement a hidden 100% strength.

## Predicted duration by size

For estimate `E`:

```
if the corpus has at least one story whose estimate is E:
    pred(E) = mean(D) of those stories
else if V_rate is defined:
    pred(E) = E / V_rate
else:
    # cold start — pack does not use pred; see Cold start
```

Unestimated Features: cost / duration **0** for packing fill rules. They still **cannot Start**. They do not contribute to velocity.

Bugs, chores, and releases: duration **0** for packing unless the bugs-and-chores-points toggle is on and the bug/chore has an estimate — then `pred(estimate)` applies as above.

## Cold start

Until the corpus has at least one Feature with `started_at` and `accepted_at` (or `time == 0`):

- `V_rate` is not used.
- Pack uses `initial_velocity` (default **10**, Owner-editable) as **estimate-points that fit in a full window**. Leave the band **short** rather than overfill.
- Future Backlog bands also pack at those 10 points, so dates exist from day one if Features are estimated and ranked.
- This is bootstrap only.

After the first completion, the duration model is the only model. `V_rate` and `pred` come from the corpus. If a later lookback is empty or `time == 0`, `V_rate` is again undefined and the same bootstrap packing applies — one rule, not a second formula.

Copy on an empty Current with an empty ranked list: “Nothing in this window yet. Pull a story from Icebox or create one.”

Do not show `0` velocity on a new project unless the Owner set initial velocity to 0.

## Packing is a pure function

No Iteration entity. Stories stay on one ranked list. The UI draws bands from the function's output. The board is a view.

```
pack(ordered_stories, pred, V_rate, L, now, TZ) → bands
```

- `ordered_stories` — ranked list, Icebox excluded; highest priority first
- `pred` — the size → predicted-duration map (unused during cold start)
- `V_rate` — live rate, or undefined
- `L`, `now`, `TZ` — window clock
- band 0 = Current (open window); later bands = Backlog chrome

Each band carries a computed `starts_on` / `ends_at` from the window clock, then the stories that packed into it.

Re-run whenever stories, estimates, start / accept (or other verbs), rank, type, length, timezone, start weekday, strategy, initial velocity, or the bugs-and-chores toggle changes. No Recalculate button.

### Duration (or bootstrap cost)

**Duration model** (`V_rate` defined):

```
duration(story) =
  pred(estimate)  if the story has an estimate and is a Feature
                  (or a bug/chore with the points toggle on)
  0               otherwise   # unestimated, bugs, chores, releases
```

**Cold start** (`V_rate` undefined):

```
cost(story) =
  estimate  if type == feature and estimate is not null
  estimate  if bugs_and_chores_may_be_given_points
            and type in {bug, chore} and estimate is not null
  0         otherwise
```

Unestimated Features still occupy a place. They cost 0. They **cannot be Started**. They can be packed into the band being filled only while that band still has remaining budget (points during cold start; predicted-duration time otherwise).

Normative for zero-cost rows: they may sit in the band being filled **only while that band's remaining budget is unused** (cold start: packed Feature-points `< initial_velocity`; duration model: remaining time `> 0` and the band is not already over via Start overflow). Once the band is full or over, remaining zero-cost unstarted stories overflow with the rest. In-flight zero-cost stories (a started Bug) always stay in Current.

### In-flight and accepted-this-window

Always in Current, regardless of cost or predicted duration:

- `started`, `finished`, `delivered`, `rejected`
- `accepted` with `accepted_at` inside the current window (`starts_on ≤ accepted_at < ends_at`)

These can make Current **over** capacity. Auto-plan never kicks them out. They stay in Current until the current window ends. They do not consume the remaining-time budget used to pack unstarted stories.

### Auto-plan algorithm (default)

**Duration model**

1. Place every in-flight and accepted-this-window story into Current (band 0), preserving relative order.
2. Remaining capacity in Current = time from `now` to `ends_at` (the remaining predicted-duration budget). Equivalent point view: `capacity_remaining = V_rate × time_remaining_in_window`.
3. Walk the remaining `unstarted` / Release-`started` stories in rank order:
   - A story's predicted duration is `pred(estimate)` (or 0 if unestimated / non-pointed).
   - If `duration > remaining time`, **do not place it in this band**. Leave the band **short**. Open the next band with a full window of `L` days and try again.
   - If `duration ≤ remaining time`, place it here and subtract.
   - Zero-cost story: place here only if this band still has remaining time and is not already over; else overflow.
   - Release markers walk with the list (cost 0) and sit where they fall — which, if the team placed them at the **end** of the milestone, is after that milestone's last story.
4. Repeat until the ranked list is exhausted.

**Cold start**

Same walk, but the budget is `initial_velocity` estimate-points in each full window (including Current). Leave short rather than overfill. Do not use `pred` or `V_rate`.

Future bands are **visual**. They are not rows. A story in a future band is still just a ranked `unstarted` (or a Release marker). Changing rank or estimate moves it by re-running `pack`.

### Fit rules (normative)

- **Never split a story.**
- **Leave Current (and each future band) short** rather than overfill with the next estimated Feature.
- **Start is the overflow valve.** Start on a Backlog or Icebox story jumps it to Current and **may** push Current over capacity. Icebox Start also unscheduled → started (Icebox start; Feature still needs an estimate).
- **Oversized story (`pred >` a full window of `L` days;** cold start: `cost > initial_velocity`):** it never auto-fills into a band that already has packed estimated work. It becomes the first (and only estimated) occupant of the next empty band — and that band is still **short of a second story** but **over** capacity after it is placed. This is the one auto-plan exception to “leave short rather than overfill,” because the alternative is a story that never packs. Current: if Current already has packed estimated work or in-flight / accepted-this-window Features, the oversized story does **not** enter Current via auto-plan (Start can still pull it). If Current has no estimated work yet, the oversized story **does** enter Current and Current is over capacity.

- **Unstarted cannot sit above started** in the ranked list (drag rejected). The pack function does not violate this: in-flight stories were already started and stay higher if the user left them there; auto-fill only appends unstarted *below* in-flight in Current when packing from Backlog. If the user has unstarted above in-flight, that is already illegal in the UI; the server rejects that rank write.

- **Accepted this window** stay in Current until the window ends, even if that makes Current look “done.”

### Manual planning (Phase 3)

Off by default. When off for Current: Current contains only in-flight, accepted-this-window, and stories **explicitly moved** there. No velocity fill into Current. Backlog future bands still auto-plan (duration model, or bootstrap points if `V_rate` is undefined). Restoring auto-plan reruns `pack`.

## Re-plan is live

Recompute `pack` on: create, delete, icebox, schedule, reorder, estimate change, type change, start (and other verbs), accept, reject, restart, undo, window end, length / TZ / strategy / initial-velocity / bugs-and-chores toggle change.

No Recalculate button. The board after a mutation is already right.

## Window end

Clock crossing only. At `ends_at` of the open window (midnight, project TZ):

1. That window is completed. There is no insert and no sample row.
2. Every story accepted in that window **ages into Done** as a **flat** list (newest accepted first). No group-by-window chrome.
3. Those accepted Features (with `started_at` and `accepted_at`) now sit in a completed window, so they join the lookback corpus if they fall in the last `K` completed windows.
4. `V_rate` and `pred` recompute from stories.
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
- Place at the **end** of the milestone's stories (work **above**, marker **below**).
- Auto-`started` when created in Backlog or dragged from Icebox into the ranked list.
- Finish → `accepted`. The marker may sit accepted in Current until the window ends, then Done, like other accepted stories — or, if Finish happens when it is in a future band, it accepts in place. Prefer: Finish is allowed when the team says the milestone shipped; the marker becomes `accepted` and, if that is in the current window, stays in Current until the window ends.
- Optional **target date**. This is the only date a human types. It does **not** move stories and does not override the plan.
- Colour (when a target date is set):
  - **Blue (on track):** `starts_on` of the **computed window that contains the marker** ≤ target date.
  - **Red (late):** that computed window's `starts_on` **>** target date.
- If the marker is in Icebox: no colour, or muted “unscheduled.”
- No stories above the marker: still a marker; date/colour use the computed window it packed into (often Current with an empty pack). Copy may say “No stories above this release.”
- Two releases: each marker owns the work **above it and below the previous marker**. Colour/date still use the computed window that **contains that marker**.

## Charts

### Velocity (Phase 2)

- X: completed windows (`ends_at`).
- Y: derived, not stored. For each completed window, compute live from stories whose `accepted_at` falls in that window: `work / time` (corpus rules) or total completed estimate / window length.
- Bar per completed window. Line at current `V_rate` converted to a full-window capacity (`V_rate × L`), or at `initial_velocity` while the corpus is empty.
- Current is not a completed bar. Optional faint “accepted so far this window.”
- Empty: “No completed stories yet.”

### Burn-up (Phase 2)

- Accepted line: cumulative estimate of Features accepted by each window end, derived from stories' `accepted_at` and current estimates.
- Scope line: sum of Feature estimates on the ranked list (unestimated = 0, Icebox excluded), live “now.” Historical scope is the same live calculation at each past `ends_at` from stories that exist.
- Bugs/chores/releases add 0 unless the Phase 3 toggle is on.
- Empty: “No scope yet.”
- Cycle time (Phase 3) is the same clock as `D`: first `started_at` → `accepted_at`.

### Cycle time (Phase 3)

- First `started_at` → `accepted_at`. Reject / Restart does not reset the clock.
- p50 / p75 / p95. Filter by type. No per-person chart.

## Panel membership

| Panel | Contains |
| --- | --- |
| Icebox | `unscheduled` only. Own order. No band headers. |
| Current | In-flight + accepted-this-window + auto-filled `unstarted` (and Releases that packed here). Header: packed estimate-points / window capacity, computed `Ends {date}`, over-capacity badge if Current exceeds capacity. Cold start capacity is `initial_velocity`. Duration-model capacity for the open window is `V_rate × time_remaining` (full-window chrome may show `V_rate × L`). |
| Backlog | The rest of the ranked list, under future **band** headers when they have members. No “iteration 3”. |
| Done | Accepted stories aged past the current window. Flat list, newest accepted first. |

## Worked example 1 — new project, cold start, then first completion

Project created Thursday 6 Aug 2026. TZ `Australia/Melbourne`. Length **7 days**, starts Monday. First start = Monday 3 Aug. Current window: Mon 3 Aug 00:00 – Mon 10 Aug 00:00. Displayed end: **9 Aug**. Corpus is empty. `V_rate` is undefined. Pack uses `initial_velocity` **10** estimate-points per full window.

`pack(ordered Features, pred=∅, V_rate=undefined, L=7, now=Thu 6 Aug, TZ=Australia/Melbourne)`:

| Pos | Story | Est |
| --- | --- | --- |
| 1 | Login | 2 |
| 2 | Reset | 1 |
| 3 | Search | 3 |
| 4 | Billing | 5 |
| 5 | Admin | 3 |

Walk Current, remaining **10** points (bootstrap):

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

Luis Starts Billing from Backlog. Billing jumps to Current. Current Feature-points = 6 + 5 = **11 / 10**, over-capacity badge. Admin stays in the 16 Aug band.

They Start Login Thursday 6 Aug 09:00 and Accept it Friday 7 Aug 09:00. `started_at` = Thu 09:00, `accepted_at` = Fri 09:00, `D` = **1 day**, estimate **2**. Login stays in **Current** (accepted-this-window). Still not Done. The open window is not in the lookback, so `V_rate` is still undefined and pack stays on bootstrap.

Monday 10 Aug 00:00 Melbourne: current window ends.

- Login → Done (flat list; newest accepted first). No row is written.
- Corpus now has Login (`accepted_at` in the completed 3–10 Aug window). `work = 2`, `time = 1 day`, **`V_rate = 2` points/day**.
- `pred(2) = 1 day` (mean of Login). No other size in the corpus, so `pred(1) = 1/2 day`, `pred(3) = 1.5 days`, `pred(5) = 2.5 days`.
- In-flight: Billing started. Stays in Current.
- Remaining time in the new window (10–17 Aug) = **7 days**.
- Reset `pred(1) = 0.5d` → Current (6.5d left). Search `pred(3) = 1.5d` → Current (5d). Admin `pred(3) = 1.5d` → Current (3.5d).
- Current: Billing (started) + Reset + Search + Admin. Packed from `pred` against the 7 days remaining, at `V_rate = 2` points/day.

## Worked example 2 — similar-size pred, leave-short, release colour vs computed window

Last three completed windows (3–10, 10–17, 17–24 Aug 2026). `K = 3`. Corpus:

| Story | Est | `started_at` | `accepted_at` | `D` |
| --- | --- | --- | --- | --- |
| Login | 2 | 3 Aug 09:00 | 5 Aug 09:00 | 2d |
| Reset | 1 | 4 Aug 09:00 | 5 Aug 09:00 | 1d |
| Search | 3 | 10 Aug 09:00 | 13 Aug 09:00 | 3d |
| Billing | 5 | 11 Aug 09:00 | 16 Aug 09:00 | 5d |
| Nav polish | 1 | 17 Aug 09:00 | 18 Aug 09:00 | 1d |
| Filters | 2 | 18 Aug 09:00 | 20 Aug 09:00 | 2d |
| Export | 3 | 19 Aug 09:00 | 22 Aug 09:00 | 3d |

`work = 17`, `time = 17 days`, **`V_rate = 1` point/day**. Full-window capacity = `1 × 7` = **7** points.

```
pred(1) = mean(1d, 1d) = 1 day
pred(2) = mean(2d, 2d) = 2 days
pred(3) = mean(3d, 3d) = 3 days
pred(5) = mean(5d)     = 5 days
```

No corpus story of estimate 8, so `pred(8) = 8 / V_rate` = 8 days (oversized versus a 7-day window).

Current window starts **24 Aug 2026**, ends at midnight **31 Aug** (displayed end **30 Aug**). `L` = 7. `now` = **Wednesday 26 Aug 00:00**. Time remaining in Current = **5 days**.

| Pos | Story | Type | Est | State | `pred` |
| --- | --- | --- | --- | --- | --- |
| 1 | Redesign nav | Feature | 3 | started | — (in-flight) |
| 2 | Copy tweaks | Feature | 1 | unstarted | 1d |
| 3 | Checkout | Feature | 5 | unstarted | 5d |
| 4 | Admin | Feature | 3 | unstarted | 3d |
| 5 | v1.0 | Release | — | started (marker) | 0 |
| 6 | SSO | Feature | 5 | unstarted | 5d |

Target date on v1.0: **31 Aug 2026**.

Walk:

- Current seed: Redesign (started). Remaining time **5d**.
- Copy (1d) ≤ 5 → Current. Remaining **4d**.
- Checkout (5d) > 4 → **leave Current short**. Checkout does not enter Current.
- Next band starts **31 Aug**, ends **7 Sep**, remaining **7d**.
- Checkout 5d → that band. Remaining **2d**.
- Admin 3d > 2 → **leave that band short**.
- Next band starts **7 Sep**, remaining **7d**: Admin 3d, v1.0 (0), remaining **4d**.
- SSO 5d > 4 → band starting **14 Sep**.

v1.0 lives in the **computed window that starts 7 Sep**. Target 31 Aug. `starts_on` 7 Sep > 31 Aug → **red**.

If Maya moves v1.0 to sit directly under Checkout, the marker packs into the window starting **31 Aug**. `starts_on` ≤ target → **blue**. Live. Colour is versus the computed window that contains the marker, not a stored row.

SSO is below the marker: not in v1.0.

## Worked example 3 — reject stays in Current; first start to accept

Checkout is Delivered in Current after Luis Started it (overflow). Maya Rejects: “tax is wrong.” State `rejected`, still Current. Restart → `started`, still Current. `started_at` is still the first Start. `V_rate` and `pred` are unchanged (Checkout is not accepted).

They Accept Checkout on **2 Sep** (next window). `D` = first `started_at` → `accepted_at` on 2 Sep. Checkout stays in Current until that window ends, then Done. It joins the corpus when its `accepted_at` sits in a completed lookback window. It does not rewrite an earlier window.

## What we will not do

- Persist an Iteration as the plan. There is no `iterations` table. Stories are not assigned to a window.
- Persist `V_rate`, `pred`, window totals, or accepted-point history. Those are calculated from stories.
- Typed velocity override (“pretend we do 20”) except the **initial velocity** setting, which is bootstrap when `V_rate` is undefined.
- Drag-to-pin a story to a date. Start or reorder; do not staple to 6 Sep.
- Split stories.
- Count bugs/chores as Features (unless the later points toggle is on).
- Show Icebox stories inside band headers.
- Per-person velocity.
- Team strength % in MVP.
- A date picker that overrides the plan. Release target date is a comparison only.
- Store length as 1–4 weeks. Length is **days**.

## QA short script (slices 8 + 20)

1. New project. Capacity displays 10. Create five estimated Features in Backlog totalling more than 10. Confirm Current is **short** of 10 rather than over.
2. Start a Backlog Feature that did not fit. Confirm it is in Current and Current may exceed 10.
3. Accept one Feature. Confirm it is still in Current, not Done. Confirm `started_at` and `accepted_at` are set.
4. Advance the test clock past midnight at the window end. Confirm that Feature is in Done (flat list) and `V_rate` is `estimate / D` from that story (`started_at` → `accepted_at`), not a window point sum.
5. Place a Release at the end of a group. Set a target date. Confirm blue/red vs the **computed window start** that contains the marker.
6. Drag unstarted above started. Must fail.
7. Confirm Icebox stories never appear in Current via auto-plan.
8. Confirm no story is assigned to an iteration row. Changing length in days replans live.
9. After a second Feature of the **same estimate** is accepted, confirm `pred(E)` is the mean of the two durations and that pack uses that duration, not the estimate against an integer V.
