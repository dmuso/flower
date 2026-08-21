-- Organisations, auth sessions/tokens, email outbox, project tenancy.
-- Additive tenancy/auth. Drops leftover iteration storage: planning is a calculation.

CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organisation_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_organisation_memberships_org_user
    ON organisation_memberships (organisation_id, user_id);
CREATE INDEX idx_organisation_memberships_user_id
    ON organisation_memberships (user_id);

ALTER TABLE organisation_memberships
    ADD CONSTRAINT fk_organisation_memberships_organisation_id
    FOREIGN KEY (organisation_id) REFERENCES organisations (id) ON DELETE CASCADE;

ALTER TABLE organisation_memberships
    ADD CONSTRAINT fk_organisation_memberships_user_id
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE users
    ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE users
    ADD COLUMN email_verified_at TIMESTAMP;

ALTER TABLE projects
    ADD COLUMN organisation_id UUID;

ALTER TABLE projects
    ADD COLUMN timezone VARCHAR(64) NOT NULL DEFAULT 'Australia/Melbourne';

ALTER TABLE projects
    ADD COLUMN velocity_strategy INTEGER NOT NULL DEFAULT 3;

ALTER TABLE projects
    ADD COLUMN initial_velocity INTEGER NOT NULL DEFAULT 10;

ALTER TABLE projects
    ADD COLUMN iteration_start_weekday INTEGER NOT NULL DEFAULT 1;

ALTER TABLE projects
    ADD COLUMN iteration_length_days INTEGER NOT NULL DEFAULT 7;

ALTER TABLE projects
    ADD CONSTRAINT fk_projects_organisation_id
    FOREIGN KEY (organisation_id) REFERENCES organisations (id) ON DELETE RESTRICT;

ALTER TABLE projects
    ALTER COLUMN organisation_id SET NOT NULL;

DROP INDEX idx_projects_slug;

CREATE UNIQUE INDEX idx_projects_organisation_slug
    ON projects (organisation_id, slug);

ALTER TABLE users
    ADD COLUMN last_project_id UUID;

ALTER TABLE users
    ADD CONSTRAINT fk_users_last_project_id
    FOREIGN KEY (last_project_id) REFERENCES projects (id) ON DELETE SET NULL;

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    last_project_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_sessions_token_hash ON sessions (token_hash);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);

ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_user_id
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_last_project_id
    FOREIGN KEY (last_project_id) REFERENCES projects (id) ON DELETE SET NULL;

CREATE TABLE auth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    consumed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_auth_tokens_token_hash ON auth_tokens (token_hash);
CREATE INDEX idx_auth_tokens_email ON auth_tokens (email);

CREATE TABLE email_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind VARCHAR(50) NOT NULL,
    to_email VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_outbox_to_email ON email_outbox (to_email, created_at);

DROP INDEX IF EXISTS idx_stories_iteration_id;
ALTER TABLE stories DROP CONSTRAINT IF EXISTS fk_stories_iteration_id;
ALTER TABLE stories DROP COLUMN IF EXISTS iteration_id;
DROP TABLE IF EXISTS iterations;
ALTER TABLE projects DROP COLUMN IF EXISTS iteration_length_weeks;
