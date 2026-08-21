DROP TABLE IF EXISTS email_outbox;
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS sessions;

ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_last_project_id;
ALTER TABLE users DROP COLUMN IF EXISTS last_project_id;

DROP INDEX IF EXISTS idx_projects_organisation_slug;
CREATE UNIQUE INDEX idx_projects_slug ON projects (slug);

ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_projects_organisation_id;
ALTER TABLE projects DROP COLUMN IF EXISTS iteration_length_days;
ALTER TABLE projects DROP COLUMN IF EXISTS iteration_start_weekday;
ALTER TABLE projects DROP COLUMN IF EXISTS initial_velocity;
ALTER TABLE projects DROP COLUMN IF EXISTS velocity_strategy;
ALTER TABLE projects DROP COLUMN IF EXISTS timezone;
ALTER TABLE projects DROP COLUMN IF EXISTS organisation_id;

ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;

DROP TABLE IF EXISTS organisation_memberships;
DROP TABLE IF EXISTS organisations;
