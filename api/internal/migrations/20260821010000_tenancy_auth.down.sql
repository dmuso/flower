ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS iteration_length_weeks INTEGER NOT NULL DEFAULT 1;

CREATE TABLE IF NOT EXISTS iterations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    number INTEGER NOT NULL,
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_iterations_project_number ON iterations (project_id, number);
CREATE INDEX IF NOT EXISTS idx_iterations_project_dates ON iterations (project_id, starts_on, ends_on);

ALTER TABLE iterations DROP CONSTRAINT IF EXISTS fk_iterations_project_id;
ALTER TABLE iterations
    ADD CONSTRAINT fk_iterations_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

ALTER TABLE stories ADD COLUMN IF NOT EXISTS iteration_id UUID;
CREATE INDEX IF NOT EXISTS idx_stories_iteration_id ON stories (iteration_id);
ALTER TABLE stories DROP CONSTRAINT IF EXISTS fk_stories_iteration_id;
ALTER TABLE stories
    ADD CONSTRAINT fk_stories_iteration_id
    FOREIGN KEY (iteration_id) REFERENCES iterations (id) ON DELETE SET NULL;

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
