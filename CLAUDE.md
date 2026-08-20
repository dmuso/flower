# Flower

* DO NOT use fallbacks, they hide bugs and create unintended application behaviour that is more difficult to identify and correct. IF YOU USE FALLBACKS YOU WILL BE FIRED!
* You must have evidence of any claim that you make.
* Use accurate language and obtain code evidence that is presented to the user.
* This is a complex project and will require you to gather evidence across many areas, database, backend, and frontend.
* Be proactive and eager to obtain learning and information so that your solutions are correct and fit for purpose.
* All work in the project is critical, production ready work, it's more important to get it right than to get it done quickly, so take as long as you need. You don't need to rush work to have it done within one session, prioritise quality over quantity. If your solution is only partly complete, but is of the highest quality, that is fine, just provide the user with an update and a prompt to use in a following session.

## Project Context

Flower is a task management workflow tool in the spirit of Pivotal Tracker, built as a monorepo application.

Use UK spelling for tables, fields, functions, variables, etc. Example: "organisations" instead of "organizations".

## Development Workflow

* Use `nix-shell` to run every `make` command so that you execute within the defined environment.
* Use Test Driven Development when making changes - write a failing test, then write code to make it pass.
* Always run tests and linting after a change to verify your work.
* Do NOT ignore errors from tests, scripts, and linting checks. Errors are the important feedback you need to ensure that your code works as intended. Always action errors and warnings, never ignore them.

### Environment Initialisation

1. **Local Hot Reload**: Run `nix-shell --run 'make dev'` when iterating on the Go API and frontend—this brings up Postgres via Docker, then runs the API (Air) and frontend (Bun dev server) on the host with hot reload.
2. **Testing**: Run `nix-shell --run 'make test'` to ensure everything is working.

### Daily Development

- **Make Commands**: Use Make commands from the project root for common operations
- **Code Changes**: Make changes in the respective service directories
- **Testing**: Always run `nix-shell --run 'make test'` and `nix-shell --run 'make lint'` after making changes.

### Project Structure

```
flower/
├── api/                # Go backend API (cmd/, internal/)
├── docs/               # Project documentation
├── frontend/           # SolidJS + Bun SPA
├── tmp/                # Temporary file location
├── docker-compose.yml  # Postgres and optional containerised services
├── Makefile            # Monorepo orchestration entry point
└── README.md
```

## Service Management

### Backend API (Go)

- **Port**: 8180 (host). Prophet uses 8080; do not bind Flower to 8080.
- **Framework**: Standard Go with Gin, exposing RESTful endpoints
- **Testing**: `nix-shell --run "cd api && make test"`
- **Build**: `nix-shell --run "cd api && make build"`

### Frontend (SolidJS + Bun)

- **Port**: 4273 (host). Prophet uses 4173; do not bind Flower to 4173.
- **Development**: `nix-shell --run "make frontend-dev"` (Bun watch + Tailwind)
- **Build/Test**: `nix-shell --run "make frontend-build | test-frontend | lint-frontend"`

### Supporting Services

1. **db / db-test** – PostgreSQL instances (primary + test)

## Docker Management

### Services

1. **api** – Go backend application (host port 8180)
2. **frontend** – SolidJS SPA (host port 4273)
3. **db / db-test** – PostgreSQL instances (host ports 5433/5437)

### Environment Variables

Copy `.env.example` to `.env` and adjust values for your machine. Treat `.env.example` as the single source of truth for required variables.

## Best Practices

1. **TDD**: Always write tests first
2. **Code Coverage**: Maintain 90%+ line coverage
3. **Environment**: Use environment variables for configuration
4. **Real Implementations**: Never introduce placeholder or "fake" code paths—tests must exercise the same production-grade logic as runtime code with no exceptions.
5. Store AI planning docs in `docs/` directory

- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems
- ❌ Do NOT clutter repo root with planning documents

## References

- [Product overview](docs/product/overview.md)
- [Frontend Design Guide](docs/reference/frontend-design-guide.md)
- [Technology choices](docs/reference/technology-choices.md)
- [Development](docs/reference/development.md)
- [Migrations](docs/migrations.md)
