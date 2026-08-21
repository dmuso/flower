package story

import "time"

type Story struct {
	ID        string     `json:"id"`
	ProjectID string     `json:"project_id"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	StoryType string     `json:"story_type"`
	Estimate  *int       `json:"estimate"`
	Rank      string     `json:"rank"`
}

type Pack struct {
	CurrentPoints       int       `json:"current_points"`
	Denominator         int       `json:"denominator"`
	VelocitySource      string    `json:"velocity_source"`
	CurrentWindowEndsAt time.Time `json:"current_window_ends_at"`
}

type List struct {
	Stories []Story `json:"stories"`
	Pack    Pack    `json:"pack"`
}
