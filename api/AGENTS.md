# API development guidelines

## Coding practices

* Always use idiomatic Go conventions
* Use domain driven design concepts, keep related domain concepts together
* Always use TDD, write failing tests first, then code to make a test pass
* Prefer API contract testing over unit tests, include all API endpoints and code paths in tests
* Test against real code and a real server wherever possible, avoid mocking and stubbing
* Always explicitly check error conditions, do not suppress errors
* When checking errors always log with where the error came from and what was attempted
* Be careful on the over-use of dependency injection. If you do not have tests that leverage it, then don't use it. It's additional noise for calling code to have to setup and manage the dependencies to inject. It's better to have the detail of the implementation abstracted from the caller. It's much clearer for individual files to import the modules that they require as dependencies.
* Make logged error messages unique so they can be debugged easily
* This is a production ready application, always implement functions and capability end to end, never use placeholders or TODOs and never skip tests

## Package structure

The codebase follows a domain-driven vertical slices architecture with clear separation between business domains, infrastructure, and HTTP.

```
api/internal/
├── domain/           # Business domains (vertical slices)
├── platform/         # Cross-cutting infrastructure
├── app/              # Application wiring and lifecycle
├── handlers/         # HTTP route registration
└── migrations/       # Schema migrations
```

### domain/ - Business domains

Each domain is self-contained with handler/service/repository pattern.

| Package | Purpose |
|---------|---------|
| `user/` | User accounts |
| `project/` | Projects, memberships, iterations |
| `story/` | Stories, labels, ranking, state transitions |

### platform/ - Infrastructure

| Package | Purpose |
|---------|---------|
| `config/` | Configuration loading and validation |
| `db/` | Database connection and migrations |
| `middleware/` | HTTP middleware (CORS, logging) |

### app/ - Application wiring

Application bootstrap and lifecycle (`Start`, `Shutdown`), plus migrate/rollback/force commands.

### handlers/ - HTTP routes

Route registration for HTTP endpoints.

## Architectural constraints

- Domain packages should not import each other directly
- Keep files under 500 lines where practical
- Run `make test` and `make lint` after changes
