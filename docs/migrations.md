# Database migrations

This project uses [go-migrate](https://github.com/golang-migrate/migrate) for database schema migrations.

## Overview

The migration system allows you to:

- Create and manage database schema changes
- Apply migrations to the database
- Rollback migrations when needed
- Force specific migration versions
- Track migration history

## Migration files

Migration files are stored in `api/internal/migrations/` with the following naming convention:

- Format: `XXX_description.up.sql` and `XXX_description.down.sql`
- XXX: sequential number with leading zeros for the first migrations (for example `000001`), then timestamps (`YYYYMMDDHHMMSS`) to avoid clashes between engineers
- Description: brief description of the migration

### Current migrations

#### 000001_create_core_schema

Creates the initial Flower schema:

- `users`
- `projects`
- `project_memberships`
- `iterations`
- `stories`
- `labels`
- `story_labels`
- `activities`

Constraints:

- Schema only. Data changes belong in Go boot-time backfills.
- No database enums, triggers, or functions. Business rules live in the Go API.
- Table names are plural. Primary keys are UUIDs named `id`.
- Australian/New Zealand English for names.

## Usage

### Using Make commands (recommended)

```bash
# Run all pending migrations
nix-shell --run "make migrate"

# Rollback the last migration
nix-shell --run "make migrate-down"

# Force to a specific version (use with caution)
nix-shell --run "make migrate-force VERSION=1"
```

### Using direct Go commands

```bash
cd api

# Run migrations
go run ./cmd/server migrate

# Rollback last migration
go run ./cmd/server rollback

# Force to specific version
go run ./cmd/server force 1
```

The API also runs pending migrations on start.

## Creating new migrations

1. Create the up migration file under `api/internal/migrations/`.
2. Create the matching down migration file.
3. Use `date +%Y%m%d%H%M%S` for the numeric prefix after the initial numbered migrations.
