# Development

## Tooling split

- **Nix Shell** provides Go, Bun, Air, Overmind, and golangci-lint. Run every `make` command inside `nix-shell`.
- **Docker Compose** runs PostgreSQL. Day-to-day development does not run the API or frontend in Docker.

## Local hot reload

```bash
nix-shell --run "make dev"
```

This:

1. Starts the `db` service in Docker Compose (Postgres on host port 5433, container port 5432)
2. Starts Overmind with `Procfile`
3. Runs the API with Air on port 8180
4. Runs Tailwind watch and the Bun SPA on port 4273

Stop with Ctrl+C. Postgres keeps running until `nix-shell --run "make stop"`.

## Tests

```bash
nix-shell --run "make test"
nix-shell --run "make lint"
```

API tests that need a database use `.env.test` (Postgres on host port 5437 via the `db_test_host` Compose service).

Flower host ports are distinct from Prophet so both stacks can run on one machine:

| Service | Flower host | Prophet host |
| --- | --- | --- |
| API | 8180 | 8080 |
| Frontend | 4273 | 4173 |
| Postgres | 5433 | 5432 |
| Postgres test | 5437 | 5436 |

## Environment files

- `.env.example` is the source of truth for required variables.
- Copy it to `.env` for local development.
- `.env.test` is used by the test Makefile targets.
