# Multitenancy

Change name: `flower`

Flower is multitenant from day one. There is no single-tenant mode. There is no “dev board” that skips organisation isolation.

This file is the tenancy product. The Technical Lead designs how it is enforced. They may not loosen these rules.

Spelling: **organisation**, not organization. Use that in new tables, fields, and docs.

## Existing ground

`000001` has `users`, `projects`, `project_memberships` — **no organisations table**, no `projects.organisation_id`. That is a gap, not a licence to skip tenancy. Slice 0 **adds** organisations and attaches projects to them. Do not redesign the eight core tables to get there.

`project_memberships.role` is a string. Product values: `owner`, `member`, `viewer` only.

## Hierarchy

```
Account (human login)
  └── Membership in one or more Organisations
        └── Project
              └── Project membership (owner | member | viewer)
        └── Agent (named actor with a token)
        └── Workspace (Phase 3; a view, not a tenant)
```

### Account

A person. Email unique. Authenticates with email + password and/or magic link in MVP. SSO later.

`users.username` is `NOT NULL` today — see [open-questions.md](./open-questions.md). One email, one account. An account can belong to many organisations. An account is not a tenant. Data is always organisation → project → story.

Deleting an account (no slice yet) must not delete an organisation that still has another owner.

### Organisation

The **tenant**. Isolation boundary. Billing boundary when billing exists. All projects, stories, files, tokens, webhooks, and audit live under an organisation.

Created on first-run signup (slice 0): the new account becomes the first **organisation owner** and must create a first project in the same flow. No organisation-without-a-project on first run. After someone deletes the last project, empty-organisation state is “create a project.”

URLs must not leak another organisation’s ids. Cross-tenant fetch of a guessed UUID returns **404**, not 403.

### Project

One ranked list, one Icebox, one velocity, one board. Lives in exactly one organisation. Cannot move organisations in this spec.

Creating a project: **organisation owner** only in MVP. Creator is a project owner.

Existing columns stay: `name`, `slug`, `description`, `point_scale`, `iteration_length_weeks`. Add (without rewriting): `organisation_id`, and later timezone, velocity strategy, initial velocity, bugs-and-chores-estimable. `slug` is unique today **globally**; under tenancy it should be unique **per organisation**. Technical Lead migrates that carefully; do not invent a new slug scheme in this spec.

### Project membership and roles

Exactly one role per membership:

| Role | Meaning |
| --- | --- |
| **owner** | Full project control. Invite, change roles, settings (iteration clock, scales, toggles), delete stories, create agent tokens, export/import (Phase 3). Can accept / reject. |
| **member** | Build. Create/edit/estimate/reorder/icebox, all legal verbs including **accept and reject**, tasks, labels, comments, attachments, tokens for projects they can write, saved searches, My Work. Cannot invite. Cannot change roles. Cannot change iteration / scale settings. Cannot import/export. |
| **viewer** | Read. Board, stories, activity, charts, search, Icebox, images. Nothing that mutates. Cannot accept. |

**Any Member can accept.** Requester *should*. My Work surfaces the requester’s Delivered stories. There is no accept ACL among Members and Owners. History is undo. This supersedes `docs/product/overview.md` (“the requester accepted the work”) as a hard rule.

**Organisation owner** can do anything a project owner can on every project in the organisation, even without a membership row.

A person may be `viewer` on project A and `member` on project B in the same organisation.

- An organisation always has **at least one** organisation owner.
- A project is always reachable by an organisation owner. UI must not strand a project with zero humans able to administer it.

### Requester vs project owner vs story owner

- **Organisation owner** / **project owner**: tenancy roles.
- **Requester**: default creator. Human Member or Owner. Never an agent. *Should* accept; is not the only person who can.
- **Story owner**: up to **5** Members, Owners, or agents doing the work. Start assigns the clicker.

## Isolation rules

1. **Tenant key is organisation.** Every story, comment, file, token, webhook, iteration, search, and membership is reachable only through an organisation the caller belongs to.
2. **No cross-organisation read.** Session for org A must not receive org B’s stories by id, search, webhook, or export. Cross-tenant id → **404**.
3. **No cross-organisation write.** The project’s organisation must match the session or token scope.
4. **Tokens do not span organisations.** Minted in one organisation. A second organisation needs a new token.
5. **Workspaces do not span organisations.**
6. **Attachments are not public without auth.** Guessed file URL serves nothing to the wrong tenant or a signed-out client. Viewers of **that** project can see the image.
7. **Invites are project-scoped.** Accepting joins that project (and lists the organisation) without access to other projects.
8. **Search, My Work, charts, webhooks, activity** are project-scoped unless a workspace *displays* projects the user can already open. Aggregation is not a new permission.
9. **Emails** may mention a project name and a link. They must not include another tenant’s titles in the same mail.
10. **No product superuser.** Support break-glass is out of scope.

Enumeration: project lists return only projects the account can open. Organisation lists return only organisations they belong to.

## What viewers cannot do

If a mutation exists, they cannot call it.

They **cannot**:

- Create, edit, delete, estimate, reorder, change type
- Start, finish, deliver, accept, reject, restart, undo, icebox, schedule
- Tasks, blockers, labels (they can filter), owners, requester, follow
- Comment, edit description, upload or delete attachments
- Invite, revoke, change roles, settings
- Create or revoke agent tokens, register webhooks
- Export or import
- Create or delete epics (they can view and filter)
- Change point scale or planning mode
- See other members’ emails if display names suffice

They **can**:

- Sign in
- Open projects they were invited to
- Read the board (Icebox / Backlog / Current / Done), story detail, Markdown, images, activity, tasks, blockers, labels, epics, charts, search
- Change their own password / request a magic link
- Make a **personal** workspace of projects they can view (Phase 3)

Fail QA if a viewer can `POST` a story via the API.

## What members cannot do

- Invite or change roles
- Delete the project or the organisation
- Change iteration length, timezone, start weekday, point scale, velocity strategy, initial velocity, bugs-and-chores toggle, auto-plan flag
- Import or export CSV
- Remove an Owner
- Create a project (organisation owner only, MVP)
- Mint a token for a project they do not belong to as member/owner

They **can** accept, reject, restart, undo, and create agent tokens for projects they can write. Tokens are agents, not impersonations. See [agent-api.md](./agent-api.md).

## Workspaces (Phase 3)

A workspace is a **named, personal (per account) set of projects inside one organisation**.

- Not a tenant. Not a permission boundary. Not shared in this slice.
- Default: “All my projects” in that organisation.
- Cannot contain two organisations.
- Deleting the workspace deletes the view, not the projects.

Until Phase 3: organisation → project switcher only.

## Organisation membership vs project membership

On invite accept:

- The organisation appears in their switcher. Organisation roles in the UI are **organisation owner** or **not**. Everyone else is “in the organisation because they have at least one project membership.”
- They receive the invited **project** role.

Organisation owners are appointed by existing organisation owners.

Removing someone from their last project in the organisation removes the organisation from their list. It does not delete their account.

## Invites

- Only project owners and organisation owners invite to a project.
- Role on the invite: owner / member / viewer.
- Email already on the project: error, no duplicate.
- 14-day expiry; resend invalidates the old link; revoke kills it now.
- New email: they create an account on accept (password or magic link).
- Invite does not grant other projects.

## Settings

| Setting | Who | Scope |
| --- | --- | --- |
| Organisation name | Organisation owner | Organisation |
| Project name | Project owner | Project |
| Iteration length, start weekday, timezone | Project owner | Project |
| Velocity strategy (1–4), initial velocity | Project owner | Project |
| Point scale, bugs-and-chores points, auto-plan Current | Project owner | Project (later slices) |
| Members and invites | Project owner | Project |
| Agent tokens | Project owner or member (their grants) | Organisation + listed projects |
| Webhooks | Member or owner | Project |

No per-project custom roles. No per-field permissions.

## Isolation tests (ride with slice 0, 1, and the agent slice)

- Same story title in org A and org B. Session A search returns only A.
- Session A fetches story id from B → 404.
- Token from A on a project in B → unauthorized / 404.
- Attachment URL from A, signed-out → no file.
- Viewer session, `POST /stories` → forbidden.
- Workspace builder cannot pick org B’s project.
