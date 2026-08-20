# Technical reference index

## Code style

### Defensive coding

Follows defensive coding patterns: all error conditions are checked, and all error conditions are logged with a contextual error message. No undefined behaviour.

### Idiomatic Go

The Go API backend follows idiomatic Go conventions.

### 90% code coverage target

Maintain a 90% line coverage unit test target across all codebases.

### Linting

All code is lint-checked. The API uses golangci-lint. The frontend uses oxlint and `tsc --noEmit`.

## Architecture

### Domain-driven vertical slices

The API follows a domain-driven layout with vertical slices (`domain/`, `platform/`, `app/`, `handlers/`).

### Error handling

Use `fmt.Errorf` with `%w` for error wrapping.

### Structured logging

Use Zap with contextual fields (`component`, `operation`, error).

### Database

PostgreSQL is the source of truth. Schema changes go through go-migrate. Data backfills run in the Go API on boot, not in SQL migrations.

## Documents

- [Technology choices](technology-choices.md)
- [Development](development.md)
- [Frontend design guide](frontend-design-guide.md)
- [Bootstrap credentials](bootstrap-credentials.md)
