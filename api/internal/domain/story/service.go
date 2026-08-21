package story

import (
	"context"
	"database/sql"
	"fmt"

	"flower/api/internal/domain/planning"
	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type Service struct {
	db     *sql.DB
	access ports.Access
	clock  ports.Clock
	load   func(ctx context.Context, projectID string) (timezone string, createdAt sql.NullTime, lengthDays, weekday, initialVelocity int, err error)
}

func NewService(database *sql.DB, access ports.Access, clock ports.Clock) *Service {
	s := &Service{db: database, access: access, clock: clock}
	s.load = s.loadProject
	return s
}

func (s *Service) loadProject(ctx context.Context, projectID string) (string, sql.NullTime, int, int, int, error) {
	var timezone string
	var created sql.NullTime
	var lengthDays, weekday, velocity int
	err := s.db.QueryRowContext(ctx, `
		SELECT timezone, created_at, iteration_length_days, iteration_start_weekday, initial_velocity
		FROM projects WHERE id = $1
	`, projectID).Scan(&timezone, &created, &lengthDays, &weekday, &velocity)
	if err != nil {
		return "", sql.NullTime{}, 0, 0, 0, fmt.Errorf("story: load project: %w", err)
	}
	return timezone, created, lengthDays, weekday, velocity, nil
}

func (s *Service) List(ctx context.Context, userID, projectID string) (*List, error) {
	if _, _, err := s.access.OpenProject(ctx, userID, projectID); err != nil {
		return nil, types.ErrNotFound
	}
	timezone, created, lengthDays, weekday, velocity, err := s.load(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !created.Valid {
		return nil, fmt.Errorf("story: project created_at is required")
	}
	win, err := planning.CurrentWindow(s.clock.Now(), timezone, created.Time, weekday, lengthDays)
	if err != nil {
		return nil, err
	}
	pack := planning.ColdStartPack(win, velocity)
	stories, err := s.listStories(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &List{
		Stories: stories,
		Pack: Pack{
			CurrentPoints:       pack.CurrentPoints,
			Denominator:         pack.Denominator,
			VelocitySource:      pack.VelocitySource,
			CurrentWindowEndsAt: pack.CurrentWindowEndsAt.UTC(),
		},
	}, nil
}

func (s *Service) listStories(ctx context.Context, projectID string) ([]Story, error) {
	rows, err := s.db.QueryContext(ctx, `
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
