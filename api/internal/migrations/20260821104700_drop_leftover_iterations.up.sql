-- For DBs that already applied the tenancy/auth migration before the leftover DROPs were added.
DROP INDEX IF EXISTS idx_stories_iteration_id;
ALTER TABLE stories DROP CONSTRAINT IF EXISTS fk_stories_iteration_id;
ALTER TABLE stories DROP COLUMN IF EXISTS iteration_id;
DROP TABLE IF EXISTS iterations;
ALTER TABLE projects DROP COLUMN IF EXISTS iteration_length_weeks;
