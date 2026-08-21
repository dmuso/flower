package project

import "time"

const (
	DefaultTimezone              = "Australia/Melbourne"
	DefaultPointScale            = "linear"
	DefaultIterationLengthDays   = 7
	DefaultInitialVelocity       = 10
	DefaultVelocityStrategy      = 3
	DefaultIterationStartWeekday = 1
)

type Project struct {
	ID                     string    `json:"id"`
	OrganisationID         string    `json:"organisation_id"`
	Name                   string    `json:"name"`
	Slug                   string    `json:"slug"`
	Description            *string   `json:"description,omitempty"`
	PointScale             string    `json:"point_scale"`
	Timezone               string    `json:"timezone"`
	VelocityStrategy       int       `json:"velocity_strategy"`
	InitialVelocity        int       `json:"initial_velocity"`
	IterationStartWeekday  int       `json:"iteration_start_weekday"`
	IterationLengthDays    int       `json:"iteration_length_days"`
	CreatedAt              time.Time `json:"created_at"`
}
