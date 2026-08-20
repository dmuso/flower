# Database migration usage guide

## Overview

Flower uses the [go-migrate](https://github.com/golang-migrate/migrate) library. Migrations are applied by the Go API and can be driven from Make or from `go run ./cmd/server`.

## Initial files

- `api/internal/migrations/000001_create_core_schema.up.sql`
- `api/internal/migrations/000001_create_core_schema.down.sql`

## Usage examples

### Make commands

```bash
nix-shell --run "make migrate"
nix-shell --run "make migrate-down"
nix-shell --run "make migrate-force VERSION=1"
```

### Direct Go commands

```bash
cd api
go run ./cmd/server migrate
go run ./cmd/server rollback
go run ./cmd/server force 1
```

The server runs pending migrations automatically on start, so `nix-shell --run "make dev"` is enough for local schema updates.
