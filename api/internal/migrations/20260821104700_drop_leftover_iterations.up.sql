-- Idempotent leftover drop for DBs that already applied slice0 before DROPs landed in that file.
DROP INDEX IF EXISTS idx_stories_iteration_id;
ALTER TABLE stories DROP CONSTRAINT IF EXISTS fk_stories_iteration_id;
ALTER TABLE stories DROP COLUMN IF EXISTS iteration_id;
DROP TABLE IF EXISTS iterations;
ALTER TABLE projects DROP COLUMN IF EXISTS iteration_length_weeks;
