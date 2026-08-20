# Database guidelines

- Do not use DB migrations for data related changes, schema only. All data changes and updates should use the backfill pattern when the server first boots up
- Business logic belongs in the Go API
- Do not model enums in DB constraints, this should be in the Go API
- Do not use DB triggers, all logic and behaviour must be in the Go API
- Do not use DB functions, all logic and behaviour must be in the Go API
- All table names should be plural
- Primary keys should use UUIDs and be named "id" for all tables
- All table, column and index naming should be Australian/New Zealand English

Migration numbering uses timestamps after the first numbered migrations to avoid clashes between engineers. Use `date +%Y%m%d%H%M%S` when naming files.
