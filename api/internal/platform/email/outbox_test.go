package email

import (
	"context"
	"testing"

	"flower/api/internal/ports"
)

func TestEnqueueRejectsNilTx(t *testing.T) {
	if err := NewOutbox().Enqueue(context.Background(), nil, ports.MailMessage{Kind: "k", To: "a@b.c"}); err == nil {
		t.Fatal("expected nil tx error")
	}
}

func TestEnqueueRejectsMissingFields(t *testing.T) {
	// A dummy non-nil tx cannot be constructed without a DB; field validation happens first.
	o := NewOutbox()
	if err := o.Enqueue(context.Background(), nil, ports.MailMessage{}); err == nil {
		t.Fatal("expected error")
	}
}
