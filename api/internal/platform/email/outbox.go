package email

import (
	"context"
	"database/sql"
	"fmt"

	"flower/api/internal/ports"
)

type Outbox struct{}

func NewOutbox() *Outbox {
	return &Outbox{}
}

func (o *Outbox) Enqueue(ctx context.Context, tx *sql.Tx, msg ports.MailMessage) error {
	if tx == nil {
		return fmt.Errorf("email outbox: nil transaction")
	}
	if msg.To == "" {
		return fmt.Errorf("email outbox: to is required")
	}
	if msg.Kind == "" {
		return fmt.Errorf("email outbox: kind is required")
	}
	_, err := tx.ExecContext(ctx, `
		INSERT INTO email_outbox (kind, to_email, subject, body)
		VALUES ($1, $2, $3, $4)
	`, msg.Kind, msg.To, msg.Subject, msg.Body)
	if err != nil {
		return fmt.Errorf("email outbox: insert: %w", err)
	}
	return nil
}
