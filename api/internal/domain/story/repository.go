package story

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) List(ctx context.Context, projectID string) ([]Story, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, project_id, title, state, story_type, estimate, rank
		FROM stories WHERE project_id = $1
		ORDER BY rank ASC
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("story: list: %w", err)
	}
	defer rows.Close()
	out := []Story{}
	for rows.Next() {
		var st Story
		var estimate sql.NullInt64
		if err := rows.Scan(&st.ID, &st.ProjectID, &st.Title, &st.State, &st.StoryType, &estimate, &st.Rank); err != nil {
			return nil, fmt.Errorf("story: scan: %w", err)
		}
		if estimate.Valid {
			v := int(estimate.Int64)
			st.Estimate = &v
		}
		out = append(out, st)
	}
	return out, rows.Err()
}
