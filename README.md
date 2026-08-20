# Flower

A monorepo for Flower, a task management workflow tool in the spirit of Pivotal Tracker.

Stories live in Icebox until pulled into a ranked backlog; Current is this iteration; Done is accepted work from completed iterations.. The product is built as a Go API, a SolidJS frontend, and PostgreSQL.

## Project Structure

```
flower/
├── api/                # Go backend API
├── frontend/           # SolidJS + Bun client (Tailwind v4)
├── docs/               # Product, architecture, and planning docs
├── Makefile            # Root-level Make orchestrations
├── docker-compose.yml  # Postgres and optional containerised services
├── shell.nix           # Go, Bun, Air, Overmind, golangci-lint
└── AGENTS.md           # Workflow reference for AI agents
```

## Getting Started

### Prerequisites

- `nix-shell`
- Docker, Docker Compose or Orbstack

### Development Setup

1. Clone the repository
2. Copy `.env.example` to `.env`; treat `.env.example` as the source of truth for required variables.
3. Run `nix-shell --run "make setup"` to initialise the development environment
4. Run `nix-shell --run "make start"` to spin up the entire Docker stack when you want everything containerised
5. Run `nix-shell --run "make dev"` for day-to-day hot reload: it ensures Postgres is running in Docker, then launches the Go API via `air` (port 8180) and the frontend via `bun --hot` (port 4273) on your host. Use Ctrl+C to stop the local processes; the Postgres container keeps running until you `nix-shell --run "make stop"`.

### Available Make Commands

- `nix-shell --run "make setup"` – Initialise development environment
- `nix-shell --run "make start"` / `nix-shell --run "make stop"` – Start or stop every container defined in `docker-compose.yml`
- `nix-shell --run "make logs"` / `nix-shell --run "make clean"` – Tail aggregated logs or destroy containers + volumes
- `nix-shell --run "make test"` – Run API + frontend tests
- `nix-shell --run "make lint"` – API lint + frontend type-check
- `nix-shell --run "make dev"` – Start Postgres in Docker and run the API + frontend locally with hot reload (Air + Bun)
- `nix-shell --run "make frontend-dev"` / `nix-shell --run "make frontend-build"` / `nix-shell --run "make frontend-serve"` – Day-to-day SolidJS workflows
- `nix-shell --run "make test-frontend"` / `nix-shell --run "make lint-frontend"` – Bun unit tests + TS checks
- `nix-shell --run "make migrate"`, `nix-shell --run "make migrate-down"`, `nix-shell --run "make migrate-force VERSION=<n>"` – Database maintenance helpers

## Architecture

- **Backend API** (`api`) – Go binary exposing `/health`, `/ready`, `/api/version`, and `/api/v1/*` endpoints on host port `8180`
- **Frontend** (`frontend`) – Bun-powered SolidJS SPA served on host port `4273`
- **Database** (`db`, `db_test`) – PostgreSQL 17 for local development plus a dedicated test instance
- **Orchestration** – Docker Compose runs Postgres (and optionally the API and frontend). `nix-shell` provides Go and Bun. `nix-shell --run "make dev"` is the default local workflow.

## Frontend (SolidJS + Bun)

- Located under `frontend/` and scaffolded without Vite; Bun handles dependency management, bundling, tests, and the HTTP dev server.
- Tailwind v4 runs through the standalone `@tailwindcss/cli`, and CSS is generated via `bun run tailwind:watch` inside the orchestration scripts.
- Key scripts:
  - `make frontend-dev` → launches Bun watch mode, Tailwind watch, and the static dev server on port 4273.
  - `make frontend-build` → compiles Tailwind, bundles `src/index.tsx` into `dist/`, and copies `public/` assets.
  - `make frontend-serve` → serves the compiled assets (helpful for smoke testing the `dist/` output).
  - `make test-frontend` / `make lint-frontend` → run Bun unit tests and strict TS checks.
- API requests target `FRONTEND_API_URL`; it must be set when building or running the frontend (e.g. `FRONTEND_API_URL=http://localhost:8180 make frontend-build` on the host).

## Reference

- Product overview: `docs/product/overview.md`
- Technology choices: `docs/reference/technology-choices.md`
- Development workflow: `docs/reference/development.md`
- Frontend design: `docs/reference/frontend-design-guide.md`
- Database migrations: `docs/migrations.md`
