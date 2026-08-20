# Agent API

Change name: `flower`

This file is **not** a second HTTP contract. Agents call the same `/api/v1` the SPA uses. The Technical Lead designs storage and HTTP internals in `technical-approach.md`. They may not loosen the rules below.

REST is the product API (humans and agents). Outbound webhooks are a product feature for any Member/Owner client, not an agent privilege. GraphQL later unless a named slice proves it cheaper — product-wide, not an agent surface.

## Why agents are first-class

Coding agents and CI already move work. In Tracker they had to impersonate a person. That breaks audit (“who delivered this?”) and permissions (“the token can do everything Dan can”).

Flower:

- An **agent** is a first-class named actor: a `users` row (`actor_kind=agent`) with a project role (Owner / Member / Viewer), the same as a human on that project.
- Actions are attributed to the agent. The human who minted the token is recorded on the grant, not as the actor.
- The state machines in [product-spec.md](./product-spec.md) are the contract (type + from-state + **role**). Verbs, not free-form state writes. Same verbs for humans and agents.
- What an agent can do is the **project role** of the token’s actor. There is no special agent API, no special verbs, and no agent scope list.

## Same API, two credentials

Mount: **`/api/v1`**. One tree. Do not implement `/v1` beside it.

| Who | How they authenticate |
| --- | --- |
| Human (browser) | Session cookie `flower_session` |
| Agent | `Authorization: Bearer flr_<secret>` |

Both hit the same routes. Example: `POST /api/v1/stories/:id/transitions`.

## Authorization is the project role

The token authenticates as the agent’s `users` row. Then the same effective-role function as a human session ([multitenancy.md](./multitenancy.md), [technical-approach.md](./technical-approach.md)):

1. Organisation owner of the project’s org → treat as project `owner`.
2. Else `project_memberships.role` for `(project_id, user_id)`.
3. Else cross-tenant or unknown project → **404** / `not_found`.
4. Same-tenant, insufficient role on a mutation → **403** / `forbidden`.

Typical CI “Grove” is **Member** on the project: start / finish / deliver **and** accept / reject. Owner agent has the same mutations plus Owner-only resources (invite, settings, delete story, …). Viewer token is read-only — same as a Viewer human.

Do **not** implement `stories:read`, `stories:write`, `stories:transition`, `stories:accept`, `comments:write`, `tasks:write`, `webhooks:manage`, or `project:read` as token scopes. Do not reserve `stories:accept`. Do not return `human_judgment_required`.

Resource permissions already specified still apply (Viewers cannot comment; Members cannot invite; only Owners change scale; requester is never an agent). Those are role rules, not a second scope list.

## Mint and revoke

`Authorization: Bearer flr_<secret>` after mint.

- Created by a project **Owner** or **Member** for projects they can write. No scope picker.
- The minting body assigns the agent **Owner / Member / Viewer** on listed projects via `project_memberships` (assumption: a Member cannot assign `owner`; see `technical-approach.md`).
- Bound to **one organisation**. Cannot be widened. A second organisation needs a new agent.
- Secret shown **once**. Prefix may be listed later for revocation.
- Revoke is immediate. Next call → `unauthorized`.
- Viewers cannot create, list secrets, or revoke.
- Rate limits: Technical Lead picks numbers (by token vs session, not by verb). Product: a single-story transition at ≤ 1 rps must not 429 a well-behaved agent.

```
POST /api/v1/organisations/:organisation_id/agents
{ "name": "Grove", "projects": [ { "project_id": "...", "role": "member" } ] }
```

Returns the agent (id, name, organisation, projects + roles) and the raw secret once.

```
GET  /api/v1/organisations/:organisation_id/agents
POST /api/v1/agent-tokens/:id/revoke
```

`GET /api/v1/me` for an agent returns the agent (id, name, organisation, projects, **role per project**), never the minting human’s email as the actor, never a scope list.

Existing `activities.actor_id` → `users`. The activity **name** is the agent’s name. Do not attribute the minting Member.

## State machine as the contract

Do **not** `PATCH { "state": "started" }`.

```
POST /api/v1/stories/:id/transitions
{ "action": "start" | "finish" | "deliver" | "reject" | "accept" | "restart" | "schedule" | "icebox" | "undo" }
```

`reject` requires `{ "action": "reject", "reason": "..." }`. Empty reason → `validation_failed`.

Legal actions depend on **type and from-state** (copy of the product machines) and **effective role**. Viewer → `forbidden` on every mutating verb.

**Feature / Bug**

| action | from | to | Who |
| --- | --- | --- | --- |
| schedule | unscheduled | unstarted | Owner, Member |
| icebox | unstarted | unscheduled | Owner, Member |
| start | unstarted | started | Owner, Member; Feature also needs estimate |
| start | unscheduled (Icebox) | started | Owner, Member; this verb is schedule+start. Story lands started in Current (may overflow velocity). Feature must already be estimated (`0` allowed). |
| finish | started | finished | Owner, Member |
| deliver | finished | delivered | Owner, Member |
| accept | delivered | accepted | Owner, Member (human or agent) |
| reject | delivered | rejected | Owner, Member (human or agent) |
| restart | rejected | started | Owner, Member |
| undo | last state-changing activity | previous | Owner, Member (same as a human; not “own last change only”) |

**Chore**

| action | from | to | Who |
| --- | --- | --- | --- |
| start | unstarted | started | Owner, Member |
| finish | started | accepted | Owner, Member (`finish` is the verb; result is `accepted`) |
| accept / reject / deliver | — | — | `invalid_transition` |

**Release**

| action | from | to | Who |
| --- | --- | --- | --- |
| schedule / create-in-backlog | unscheduled | started | Owner, Member (auto-start) |
| finish | started | accepted | Owner, Member |
| start / deliver / reject | — | — | `invalid_transition` |

```
POST /api/v1/stories/:id/transitions
{ "action": "undo" }
```

Undo is the latest state-changing activity on that story (slice 15). Same rule for a Member/Owner agent as for a human.

Start on a Feature without an estimate → **`unestimated`**, no mutation. Start assigns the actor as a story owner if owners < 5.

Optimistic concurrency: send `revision` (or `If-Match`) as the Technical Lead specifies. Stale → `conflict`.

## Key resources (REST, `/api/v1`)

All paths are organisation-aware. Prefer `/api/v1/projects/:project_id/...` after the actor has memberships. Cross-tenant ids → 404.

```
GET    /api/v1/me
GET    /api/v1/projects/:project_id
GET    /api/v1/projects/:project_id/stories
POST   /api/v1/projects/:project_id/stories
GET    /api/v1/stories/:id
PATCH  /api/v1/stories/:id          # fields, not state
DELETE /api/v1/stories/:id
POST   /api/v1/stories/:id/transitions
POST   /api/v1/stories/:id/comments
GET    /api/v1/stories/:id/activity
POST   /api/v1/projects/:project_id/stories/reorder
POST   /api/v1/projects/:project_id/stories/bulk-transitions
GET    /api/v1/projects/:project_id/iterations
GET    /api/v1/projects/:project_id/search?q=
POST   /api/v1/projects/:project_id/webhooks
GET    /api/v1/projects/:project_id/webhooks
DELETE /api/v1/webhooks/:id
POST   /api/v1/organisations/:organisation_id/agents
GET    /api/v1/organisations/:organisation_id/agents
POST   /api/v1/agent-tokens/:id/revoke
```

Create story body (shared with the human client):

```
{
  "title": "Fix login timeout",
  "story_type": "bug",
  "requester_id": "<human uuid>",
  "description": "optional markdown",
  "estimate": null,
  "panel": "icebox" | "backlog"
}
```

- One API. Omitted `story_type` **may default to Feature** (the SPA also posts Feature). Not “required, server does not default.” See `technical-approach.md`.
- `requester_id` is required and must be a human Member/Owner for every actor (including humans). Agents still cannot be requester.
- `panel` omitted → `icebox` (`unscheduled`) for everyone. `backlog` → `unstarted` (Release → auto-`started`).
- Estimate on create is allowed only if the type/scale/toggle allows it.

`PATCH` may change title, description, estimate (legal values), labels, owners (max 5), requester (human). It may **not** change `state`. Type change only from `unscheduled` / `unstarted`.

Reorder:

```
{ "story_id": "...", "before_id": "..." | "after_id": "...", "revision": 3 }
```

Callers must send explicit neighbours. Do not default to “top.” Server rejects a rank that would put `unstarted` above `started` (`illegal_rank`).

Search `q` accepts the Tracker-like language in product-spec slice 22.

## Bulk ops

Shared route. Humans may call it.

```
POST /api/v1/projects/:project_id/stories/bulk-transitions
Idempotency-Key: <client uuid>
{
  "items": [ { "story_id": "...", "action": "deliver" }, ... ]
}
```

- Max 50.
- **All-or-nothing.** One illegal item fails the batch; none apply.
- Replay of the same `Idempotency-Key` + body returns the original result.
- Same key, different body → `conflict`.

## Webhooks

Outbound events. Member or Owner on the project registers a URL (same as a human). Viewer cannot.

Events (at least):

- `story.created`, `story.updated`, `story.reordered`
- `story.started`, `story.finished`, `story.delivered`, `story.accepted`, `story.rejected`, `story.restarted`
- `comment.created`
- `iteration.completed`
- `membership.changed`

Envelope:

```
{
  "event_id": "uuid",
  "event": "story.delivered",
  "organisation_id": "...",
  "project_id": "...",
  "occurred_at": "ISO-8601",
  "actor": { "kind": "agent" | "human", "id": "...", "name": "Grove" },
  "story": { "...subset..." }
}
```

- Signed. Header `X-Flower-Signature: t=<unix>,v1=<hex>` where `v1` is HMAC-SHA256 of `{t}.{raw_body}` with the raw webhook secret. Receivers must reject if `|now - t| > 300s`. Any receiver **must** verify — not an agent-only duty.
- At-least-once. Retries on non-2xx. Timeout 10s.
- Receiver treats `event_id` as idempotent. Delivery is not a command. Flower does not listen to the agent’s reply body.
- No “act on this webhook by POSTing back as the user.”

## Deterministic errors

Every 4xx/409 error — **one** envelope for humans and agents:

```
{
  "error": {
    "code": "invalid_transition",
    "message": "Cannot finish a story that is unstarted",
    "from": "unstarted",
    "action": "finish"
  }
}
```

| code | When |
| --- | --- |
| `unauthorized` | Missing, revoked, or foreign token (or missing/invalid session) |
| `forbidden` | Role missing (Viewer mutation; Member calling Owner-only). Same code for a Viewer human and a Viewer agent. |
| `unestimated` | `start` on a Feature with no estimate |
| `invalid_transition` | Verb illegal for type + from-state |
| `illegal_rank` | Unstarted above started, or Icebox/rank mismatch |
| `validation_failed` | Missing title, empty reject reason, bad type, no requester |
| `not_found` | Unknown id **or** other-tenant id |
| `conflict` | Revision / idempotency mismatch |
| `rate_limited` | 429 |
| `owners_full` | Sixth owner |

There is no `human_judgment_required`. Insufficient role is `forbidden`.

No 200 with a partial surprise. No coerce (`start` that silently estimates 1).

## What an agent must never guess

- **Estimate.** Do not invent points. Only set estimate when the caller (prompt, human, or a prior explicit field) supplied it. Do not overwrite a human estimate unless the request body includes the new value **and** the actor is Member or Owner.
- **Type.** Never infer Feature vs Bug from the title. Omitted `story_type` may default to Feature (same as the SPA posting Feature).
- **State.** Always send a verb. Never write `state` directly. Never treat webhook delivery as “so it must be accepted.”
- **Accept / reject.** Do not treat Delivered as Accepted. A Member/Owner agent *may* accept or reject; that is a role, not a guess. A Viewer agent gets `forbidden`.
- **Position.** Reorder only with explicit `before_id` / `after_id`. Do not “put it at the top to be helpful.”
- **Organisation / project.** Use ids from `/me` or a previous list call. Do not infer from a name.
- **Requester.** Always a human id. Never itself.
- **Retries.** Non-idempotent transitions retry only with the same `Idempotency-Key`.
- **Secrets.** Never log the raw token. Never store a human password.
- **Impersonation.** Never send a header claiming to be Maya.

If the agent does not know, it must fail closed (`validation_failed` / skip), not guess.

## Human vs agent (same role, same verbs)

| Action | Member (human) | Member (agent) | Viewer (either) |
| --- | --- | --- | --- |
| Start estimated Feature | yes | yes | `forbidden` |
| Start unestimated Feature | `unestimated` | `unestimated` | `forbidden` |
| Finish / Deliver Feature | yes | yes | `forbidden` |
| Accept / Reject Feature | yes | yes | `forbidden` |
| Restart rejected Feature | yes | yes | `forbidden` |
| Finish Chore / Release | yes (accepts) | yes (accepts) | `forbidden` |
| Icebox start (estimated Feature) | yes | yes | `forbidden` |
| Reorder | yes | yes, explicit neighbours | `forbidden` |
| Invite, settings, scale | Owner only | Owner agent only | `forbidden` |

## QA short script (slice 23)

Same board verbs with Bearer. No scope picker.

1. Member mints an agent named Grove as **Member** on the project. Secret shown once; a second view does not reveal it.
2. Agent starts an estimated Feature → 200, actor is the agent name on the story activity.
3. Agent starts an unestimated Feature → `unestimated`, still `unstarted`.
4. Agent accepts a delivered Feature → 200, `accepted`. Agent rejects a delivered Feature with a reason → 200, `rejected`.
5. Viewer agent (or a Viewer human) accept / reject → `forbidden`, no change.
6. Agent finishes a started Chore → `accepted`.
7. Agent `PATCH` `{ "state": "delivered" }` → rejected (no such field write).
8. Revoke token → next GET is `unauthorized`. Token from org A on org B story id → 404.
9. Webhook fires on deliver; `X-Flower-Signature` verifies (`t=<unix>,v1=<hex>`); duplicate `event_id` is safe.
