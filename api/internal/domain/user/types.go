package user

import "time"

type User struct {
	ID              string
	Username        string
	Email           string
	PasswordHash    *string
	DisplayName     string
	EmailVerifiedAt *time.Time
	LastProjectID   *string
	CreatedAt       time.Time
}

type Session struct {
	ID            string
	UserID        string
	TokenHash     string
	ExpiresAt     time.Time
	LastProjectID *string
}

type AuthToken struct {
	ID         string
	Kind       string
	Email      string
	TokenHash  string
	ExpiresAt  time.Time
	ConsumedAt *time.Time
}

const (
	KindVerifyEmail = "verify_email"
	KindMagicLink   = "magic_link"
	TokenTTL        = 30 * time.Minute
	SessionTTL      = 30 * 24 * time.Hour
)

type Me struct {
	ID              string         `json:"id"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	DisplayName     string         `json:"display_name"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	Organisations   []Organisation `json:"organisations"`
	LastProject     *Project       `json:"last_project"`
}

type Organisation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Project struct {
	ID             string `json:"id"`
	OrganisationID string `json:"organisation_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}
