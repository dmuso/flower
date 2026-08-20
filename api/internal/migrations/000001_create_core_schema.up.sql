CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    point_scale VARCHAR(50) NOT NULL,
    iteration_length_weeks INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_projects_slug ON projects (slug);
CREATE INDEX idx_projects_name ON projects (name);

CREATE TABLE project_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_project_memberships_project_user ON project_memberships (project_id, user_id);
CREATE INDEX idx_project_memberships_user_id ON project_memberships (user_id);

ALTER TABLE project_memberships
    ADD CONSTRAINT fk_project_memberships_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

ALTER TABLE project_memberships
    ADD CONSTRAINT fk_project_memberships_user_id
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE TABLE iterations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    number INTEGER NOT NULL,
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_iterations_project_number ON iterations (project_id, number);
CREATE INDEX idx_iterations_project_dates ON iterations (project_id, starts_on, ends_on);

ALTER TABLE iterations
    ADD CONSTRAINT fk_iterations_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

CREATE TABLE stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    iteration_id UUID,
    requester_id UUID NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    story_type VARCHAR(50) NOT NULL,
    state VARCHAR(50) NOT NULL,
    estimate INTEGER,
    rank VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP
);

CREATE INDEX idx_stories_project_id ON stories (project_id);
CREATE INDEX idx_stories_iteration_id ON stories (iteration_id);
CREATE INDEX idx_stories_requester_id ON stories (requester_id);
CREATE INDEX idx_stories_project_state ON stories (project_id, state);
CREATE UNIQUE INDEX idx_stories_project_rank ON stories (project_id, rank);

ALTER TABLE stories
    ADD CONSTRAINT fk_stories_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

ALTER TABLE stories
    ADD CONSTRAINT fk_stories_iteration_id
    FOREIGN KEY (iteration_id) REFERENCES iterations (id) ON DELETE SET NULL;

ALTER TABLE stories
    ADD CONSTRAINT fk_stories_requester_id
    FOREIGN KEY (requester_id) REFERENCES users (id) ON DELETE RESTRICT;

CREATE TABLE labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_labels_project_name ON labels (project_id, name);

ALTER TABLE labels
    ADD CONSTRAINT fk_labels_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

CREATE TABLE story_labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    story_id UUID NOT NULL,
    label_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_story_labels_story_label ON story_labels (story_id, label_id);
CREATE INDEX idx_story_labels_label_id ON story_labels (label_id);

ALTER TABLE story_labels
    ADD CONSTRAINT fk_story_labels_story_id
    FOREIGN KEY (story_id) REFERENCES stories (id) ON DELETE CASCADE;

ALTER TABLE story_labels
    ADD CONSTRAINT fk_story_labels_label_id
    FOREIGN KEY (label_id) REFERENCES labels (id) ON DELETE CASCADE;

CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    story_id UUID,
    actor_id UUID NOT NULL,
    kind VARCHAR(100) NOT NULL,
    summary TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activities_project_id ON activities (project_id, created_at);
CREATE INDEX idx_activities_story_id ON activities (story_id, created_at);
CREATE INDEX idx_activities_actor_id ON activities (actor_id);

ALTER TABLE activities
    ADD CONSTRAINT fk_activities_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;

ALTER TABLE activities
    ADD CONSTRAINT fk_activities_story_id
    FOREIGN KEY (story_id) REFERENCES stories (id) ON DELETE SET NULL;

ALTER TABLE activities
    ADD CONSTRAINT fk_activities_actor_id
    FOREIGN KEY (actor_id) REFERENCES users (id) ON DELETE RESTRICT;
