# Agent API

Agents are first-class teammates. They are not a human’s personal token with a different User-Agent. This file is the contract. The Technical Lead designs storage and HTTP internals. They may not loosen the rules.

REST + webhooks in this slice. GraphQL later unless a named slice proves it cheaper — default remains REST.

## Why agents are first-class

Coding agents and CI already move work. In Tracker they had to impersonate a person. That breaks audit (“who delivered this?”), permissions (“the token can do everything Dan can”), and the state machine (“PATCH state=accepted because I guessed”).

Flower:

- An **agent** has a name, an organisation, a scope list, and one or more projects.
- Actions are attributed to the agent. The human who minted the token is recorded on the grant, not as the actor.
- The state machines in [product-spec.md](../core-workflow/product-spec.md) are the contract. Verbs, not free-form state writes.
- Agents must not guess. If a rule would require judgement (Feature/Bug accept or reject), the API refuses.

## Auth and scopes

`Authorization: Bearer flr_<secret>`

- Created by a project **Owner** or **Member** for projects they can write.
- Bound to **one organisation**. Cannot be widened.
- Bound to an explicit project list (one or more in that organisation).
- Secret shown **once**. Prefix may be listed later for revocation.
- Revoke is immediate. Next call → `unauthorized`.
- Viewers cannot create, list secrets, or revoke.
- Rate limits: Technical Lead picks numbers. Product: a single-story transition at ≤ 1 rps must not 429 a well-behaved agent.

Scopes (explicit, additive):

| Scope | Can |
| --- | --- |
| `stories:read` | GET stories, iterations, search, activity |
| `stories:write` | Create, edit title/description/labels/estimate/type (legal), schedule/icebox, reorder with explicit `before_id` / `after_id` |
| `stories:transition` | `start`, `finish`, `deliver`, `restart`, chore-finish (accept), release-finish (accept) |
| `stories:accept` | **Not granted by default.** Reserved. MVP agents **do not** receive this. Feature/Bug accept/reject stay human. See open questions. |
| `comments:write` | Post comments, including `@` mentions |
| `tasks:write` | Add / toggle tasks |
| `webhooks:manage` | Register / delete webhook endpoints for the bound projects |
| `project:read` | Project settings, membership names, labels, epics (read) |

A token without `stories:read` cannot poll after it writes. Typical CI token: `stories:read` + `stories:transition` + `comments:write`.

`GET /v1/me` returns the agent (id, name, organisation, projects, scopes), never the minting human’s email as the actor.

Existing `activities.actor_id` → `users`. Implementation ground: an agent needs a `users` row or a later actor model. Product: the activity **name** is the agent’s name. Do not attribute the minting Member.

## State machine as the contract

Do **not** `PATCH { "state": "started" }`.

```
POST /v1/stories/:id/transitions
{ "action": "start" | "finish" | "deliver" | "reject" | "accept" | "restart" | "schedule" | "icebox" }
```

`reject` requires `{ "action": "reject", "reason": "..." }`. Empty reason → `validation_failed`.

Legal actions depend on **type and from-state** (copy of the product machines):

**Feature / Bug**

| action | from | to | Agent? |
| --- | --- | --- | --- |
| schedule | unscheduled | unstarted | if `stories:write` |
| icebox | unstarted | unscheduled | if `stories:write` |
| start | unstarted | started | if `stories:transition`; Feature also needs estimate |
| finish | started | finished | if `stories:transition` |
| deliver | finished | delivered | if `stories:transition` |
| accept | delivered | accepted | **human only** in MVP |
| reject | delivered | rejected | **human only** in MVP |
| restart | rejected | started | if `stories:transition` |

**Chore**

| action | from | to | Agent? |
| --- | --- | --- | --- |
| start | unstarted | started | yes |
| finish | started | accepted | yes (`finish` is the verb; result is `accepted`) |
| accept / reject / deliver | — | — | `invalid_transition` |

**Release**

| action | from | to | Agent? |
| --- | --- | --- | --- |
| schedule / create-in-backlog | unscheduled | started | yes (auto-start) |
| finish | started | accepted | yes |
| start / deliver / reject | — | — | `invalid_transition` |

Undo is a human UI on activity, also available as:

```
POST /v1/stories/:id/transitions
{ "action": "undo" }
```

Agents **may** undo their own last state change on that story (automation error). They may not undo a human’s Accept.

Start on a Feature without an estimate → **`unestimated`**, no mutation. Start assigns the agent as a story owner if owners < 5.

Optimistic concurrency: send `revision` (or `If-Match`) as the Technical Lead specifies. Stale → `conflict`.

## Key resources (REST, v1)

All paths are organisation-aware. Prefer `/v1/projects/:project_id/...` after the token is bound. Cross-tenant ids → 404.

```
GET    /v1/me
GET    /v1/projects/:project_id
GET    /v1/projects/:project_id/stories
POST   /v1/projects/:project_id/stories
GET    /v1/stories/:id
PATCH  /v1/stories/:id          # fields, not state
DELETE /v1/stories/:id
POST   /v1/stories/:id/transitions
POST   /v1/stories/:id/comments
GET    /v1/stories/:id/activity
POST   /v1/projects/:project_id/stories/reorder
POST   /v1/projects/:project_id/stories/bulk-transitions
GET    /v1/projects/:project_id/iterations
GET    /v1/projects/:project_id/search?q=
POST   /v1/projects/:project_id/webhooks
GET    /v1/projects/:project_id/webhooks
DELETE /v1/webhooks/:id
```

Create story body (agent must be explicit):

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

- `story_type` required. Do not default to feature if omitted — `validation_failed`. (Humans default to Feature in the UI; agents must say.)
- `requester_id` required and must be a human Member or Owner on the project. Agents cannot be requester.
- `panel` default `icebox` (`unscheduled`). `backlog` → `unstarted` (Release → auto-`started`).
- Estimate on create is allowed only if the type/scale/toggle allows it.

`PATCH` may change title, description, estimate (legal values), labels, owners (max 5), requester (human). It may **not** change `state`. Type change only from `unscheduled` / `unstarted`.

Reorder:

```
{ "story_id": "...", "before_id": "..." | "after_id": "..." }
```

Server rejects a rank that would put `unstarted` above `started` (`illegal_rank`).

Search `q` accepts the Tracker-like language in product-spec slice 22.

## Bulk ops

```
POST /v1/projects/:project_id/stories/bulk-transitions
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

The agent or a human with `webhooks:manage` registers a URL per project.

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

- Signed (`X-Flower-Signature` HMAC, scheme the Technical Lead specifies). Agents **must** verify.
- At-least-once. Retries on non-2xx. Timeout ~10s.
- Receiver treats `event_id` as idempotent. Delivery is not a command. Flower does not listen to the agent’s reply body.
- No “act on this webhook by POSTing back as the user.”

## Deterministic errors

Every 4xx/409 error:

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
| `unauthorized` | Missing, revoked, or foreign token |
| `forbidden` | Scope or role missing |
| `human_judgment_required` | Agent called Feature/Bug `accept` or `reject` |
| `unestimated` | `start` on a Feature with no estimate |
| `invalid_transition` | Verb illegal for type + from-state |
| `illegal_rank` | Unstarted above started, or Icebox/rank mismatch |
| `validation_failed` | Missing title, empty reject reason, bad type, no requester, omitted `story_type` |
| `not_found` | Unknown id **or** other-tenant id |
| `conflict` | Revision / idempotency mismatch |
| `rate_limited` | 429 |
| `owners_full` | Sixth owner |

No 200 with a partial surprise. No coerce (`start` that silently estimates 1).

## What an agent must never guess

- **Estimate.** Do not invent points. Only set estimate when the caller (prompt, human, or a prior explicit field) supplied it. Do not overwrite a human estimate unless the request body includes the new value **and** `stories:write`.
- **Type.** Always send `story_type`. Never infer Feature vs Bug from the title.
- **State.** Always send a verb. Never write `state` directly. Never treat webhook delivery as “so it must be accepted.”
- **Accept / reject on Feature or Bug.** Deliver, then stop. A human accepts.
- **Position.** Reorder only with explicit `before_id` / `after_id`. Do not “put it at the top to be helpful.”
- **Organisation / project.** Use ids from the token or a previous list call. Do not infer from a name.
- **Requester.** Always a human id. Never itself.
- **Retries.** Non-idempotent transitions retry only with the same `Idempotency-Key`.
- **Secrets.** Never log the raw token. Never store a human password.
- **Impersonation.** Never send a header claiming to be Maya.

If the agent does not know, it must fail closed (`validation_failed` / skip), not guess.

## Human vs agent (quick)

| Action | Member (human) | Agent with typical CI scopes |
| --- | --- | --- |
| Start estimated Feature | yes | yes |
| Start unestimated Feature | no | no (`unestimated`) |
| Finish / Deliver Feature | yes | yes |
| Accept / Reject Feature | yes | no |
| Restart rejected Feature | yes | yes |
| Finish Chore / Release | yes (accepts) | yes |
| Reorder | yes | only with explicit neighbours |
| Invite, settings, scale | Owner | no |

## QA short script (slice 23)

1. Member mints a token with `stories:read` + `stories:transition`. Secret shown once; a second view does not reveal it.
2. Agent starts an estimated Feature → 200, actor is the agent name on the story activity.
3. Agent starts an unestimated Feature → `unestimated`, still `unstarted`.
4. Agent accepts a delivered Feature → `human_judgment_required`.
5. Agent finishes a started Chore → `accepted`.
6. Agent `PATCH` `{ "state": "delivered" }` → rejected (no such field write).
7. Revoke token → next GET is `unauthorized`.
8. Token from org A on org B story id → 404.
9. Webhook fires on deliver; signature verifies; duplicate `event_id` is safe.
