package ports

import (
	"context"
	"database/sql"
)

type MailMessage struct {
	Kind    string
	To      string
	Subject string
	Body    string
}

type Mailer interface {
	Enqueue(ctx context.Context, tx *sql.Tx, msg MailMessage) error
}
