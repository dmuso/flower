# Multitenancy

Change name: `flower`

Flower is multitenant from day one. There is no single-tenant mode. There is no “dev board” that skips organisation isolation.

This file is the tenancy product. The Technical Lead designs how it is enforced. They may not loosen these rules.

Spelling: **organisation**, not organization. Use that in new tables, fields, and docs.

## Schema

`users`, `projects`, `project_memberships`, and `organisations`. Every project has `organisation_id`. Slice 0 ships organisations and attaches projects to them.

`project_memberships.role` is a string. Product values: `owner`, `member`, `viewer` only.

## Hierarchy

```
User (human login)
  └── User API token (credential of that user; list/create/revoke on the user)
  └── Membership in one or more Organisations
        └── Project
              └── Project membership (owner | member | viewer)
        └── Workspace (Phase 3; a view, not a tenant)
```

### User

A person. Email unique. Authenticates with email + password and/or magic link in MVP. SSO later.

Username is inferred from the email local-part. No username field in slice 0. One email, one user. A user can belong to many organisations. A user is not a tenant. Data is always organisation → project → story.

Deleting a user (no slice yet) must not delete an organisation that still has another owner.

### Organisation

The **tenant**. Isolation boundary. Billing boundary when billing exists. All projects, stories, files, webhooks, and audit live under an organisation.

Created on first-run signup (slice 0): the new user becomes the first **organisation owner** and must create a first project in the same flow. No organisation-without-a-project on first run. After someone deletes the last project, empty-organisation state is “create a project.”

URLs must not leak another organisation’s ids. Cross-tenant fetch of a guessed UUID returns **404**, not 403.

### Project

One ranked list, one Icebox, one velocity, one board. Lives in exactly one organisation. Cannot move organisations in this spec.

Creating a project: **organisation owner** only in MVP. Creator is a project owner.

Project columns include `name`, `slug`, `description`, `point_scale`, **`iteration_length_days`** (default 7), and `organisation_id`. Timezone, velocity strategy, initial velocity, and bugs-and-chores-estimable land with their slices. `slug` is unique **per organisation**. Do not invent a new slug scheme in this spec.

### Project membership and roles

Exactly one role per membership:

| Role | Meaning |
| --- | --- |
| **owner** | Full project control. Invite, change roles, settings (length in days, scales, toggles), delete stories, export/import (Phase 3). Can accept / reject. Manages their own API tokens. |
| **member** | Build. Create/edit/estimate/reorder/icebox, all legal verbs including **accept and reject**, tasks, labels, comments, attachments, own tokens, saved searches, My Work. Cannot invite. Cannot change roles. Cannot change length / scale settings. Cannot import/export. |
| **viewer** | Read. Board, stories, activity, charts, search, Icebox, images. Nothing that mutates. Cannot accept. May create a token on their own user; it can only read. |

**Any Member can accept.** Requester *should*. My Work surfaces the requester’s Delivered stories. There is no accept ACL among Members and Owners. History is undo.

**Organisation owner** can do anything a project owner can on every project in the organisation, even without a membership row.

A person may be `viewer` on project A and `member` on project B in the same organisation.

- An organisation always has **at least one** organisation owner.
- A project is always reachable by an organisation owner. UI must not strand a project with zero humans able to administer it.

### Requester vs project owner vs story owner

- **Organisation owner** / **project owner**: tenancy roles.
- **Requester**: default creator. Member or Owner. *Should* accept; is not the only person who can.
- **Story owner**: up to **5** Members or Owners doing the work. Start assigns the clicker.

## Isolation rules

1. **Tenant key is organisation.** Every story, comment, file, webhook, search, and membership is reachable only through an organisation the caller belongs to.
2. **No cross-organisation read.** Session for org A must not receive org B’s stories by id, search, webhook, or export. Cross-tenant id → **404**.
3. **No cross-organisation write.** The project’s organisation must match the session or the **user** the token belongs to.
4. **Tokens are user credentials**, not an organisation mint path. List / create / revoke at `/api/v1/users/:id/tokens`. A token may name projects the user can already open. A request is still organisation-scoped (cross-tenant → 404).
5. **Workspaces do not span organisations.**
6. **Attachments are not public without auth.** Guessed file URL serves nothing to the wrong tenant or a signed-out client. Viewers of **that** project can see the image.
7. **Invites are project-scoped.** Accepting joins that project (and lists the organisation) without access to other projects.
8. **Search, My Work, charts, webhooks, activity** are project-scoped unless a workspace *displays* projects the user can already open. Aggregation is not a new permission.
9. **Emails** may mention a project name and a link. They must not include another tenant’s titles in the same mail.
10. **No product superuser.** Support break-glass is out of scope.

Enumeration: project lists return only projects the user can open. Organisation lists return only organisations they belong to.

## What viewers cannot do

If a mutation exists, they cannot call it.

They **cannot**:

- Create, edit, delete, estimate, reorder, change type
- Start, finish, deliver, accept, reject, restart, undo, icebox, schedule
- Tasks, blockers, labels (they can filter), owners, requester, follow
- Comment, edit description, upload or delete attachments
- Invite, revoke, change roles, settings
- Register webhooks
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
- Create a token on their own user (`/api/v1/users/:id/tokens`); that token can only read

Fail QA if a viewer can `POST` a story via the API.

## What members cannot do

- Invite or change roles
- Delete the project or the organisation
- Change length in days, timezone, start weekday, point scale, velocity strategy, initial velocity, bugs-and-chores toggle, auto-plan flag
- Import or export CSV
- Remove an Owner
- Create a project (organisation owner only, MVP)
- Mint a token for another user, or for a project they do not belong to as member/owner

They **can** accept, reject, restart, undo, and create API tokens on **their own user** (`/api/v1/users/:id/tokens`) with a role at or below their own (Member cannot mint Owner). Same `/api/v1` handlers as a session cookie.

## Workspaces (Phase 3)

A workspace is a **named, personal (per user) set of projects inside one organisation**.

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

Removing someone from their last project in the organisation removes the organisation from their list. It does not delete the user.

## Invites

- Only project owners and organisation owners invite to a project.
- Role on the invite: owner / member / viewer.
- Email already on the project: error, no duplicate.
- 14-day expiry; resend invalidates the old link; revoke kills it now.
- New email: they create a user on accept (password or magic link).
- Invite does not grant other projects.

## Settings

| Setting | Who | Scope |
| --- | --- | --- |
| Organisation name | Organisation owner | Organisation |
| Project name | Project owner | Project |
| Length in days, start weekday, timezone | Project owner | Project |
| Velocity strategy (1–4), initial velocity | Project owner | Project |
| Point scale, bugs-and-chores points, auto-plan Current | Project owner | Project (later slices) |
| Members and invites | Project owner | Project |
| API tokens | The user (own tokens). Role at or below their own on selected projects | User (`/api/v1/users/:id/tokens`) |
| Webhooks | Member or owner | Project |

No per-project custom roles. No per-field permissions.

## Isolation tests (ride with slice 0, 1, and the API-token slice)

- Same story title in org A and org B. Session A search returns only A.
- Session A fetches story id from B → 404.
- Token belonging to a user with no membership in B, used on a project in B → unauthorized / 404.
- Attachment URL from A, signed-out → no file.
- Viewer session, `POST /stories` → forbidden.
- Workspace builder cannot pick org B’s project.
