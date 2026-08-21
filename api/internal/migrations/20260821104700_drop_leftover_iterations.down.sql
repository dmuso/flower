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
