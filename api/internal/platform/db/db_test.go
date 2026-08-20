package db

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"flower/api/internal/platform/config"
)

func TestBuildConnectionStringIncludesSSLModeWhenSet(t *testing.T) {
	got := buildConnectionString(config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "flower",
		Password: "secret",
		Name:     "flower_dev",
		SSLMode:  "disable",
	})

	want := "host=localhost port=5432 user=flower password=secret dbname=flower_dev sslmode=disable"
	if got != want {
		t.Fatalf("connection string: got %q, want %q", got, want)
	}
}

func TestBuildConnectionStringOmitsSSLModeWhenBlank(t *testing.T) {
	got := buildConnectionString(config.DatabaseConfig{
		Host:     "db",
		Port:     "5432",
		User:     "flower_dev",
		Password: "flower",
		Name:     "flower_dev",
		SSLMode:  "   ",
	})

	if strings.Contains(got, "sslmode=") {
		t.Fatalf("expected blank SSLMode to omit sslmode, got %q", got)
	}
	if !strings.Contains(got, "host=db") || !strings.Contains(got, "dbname=flower_dev") {
		t.Fatalf("expected host and dbname in connection string, got %q", got)
	}
}

func TestWithTxRejectsNilDatabase(t *testing.T) {
	err := WithTx(context.Background(), nil, func(*sql.Tx) error { return nil })
	if err == nil {
		t.Fatal("expected nil database handle to fail")
	}
	if !strings.Contains(err.Error(), "nil database handle") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithTxRejectsNilFunction(t *testing.T) {
	handle, err := sql.Open("pgx", "host=localhost port=5432 user=flower dbname=flower sslmode=disable")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := handle.Close(); closeErr != nil {
			t.Errorf("close handle: %v", closeErr)
		}
	})

	err = WithTx(context.Background(), handle, nil)
	if err == nil {
		t.Fatal("expected nil transaction function to fail")
	}
	if !strings.Contains(err.Error(), "nil transaction function") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithTxRejectsNilContext(t *testing.T) {
	handle, err := sql.Open("pgx", "host=localhost port=5432 user=flower dbname=flower sslmode=disable")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := handle.Close(); closeErr != nil {
			t.Errorf("close handle: %v", closeErr)
		}
	})

	var ctx context.Context
	err = WithTx(ctx, handle, func(*sql.Tx) error { return nil })
	if err == nil {
		t.Fatal("expected nil context to fail")
	}
	if !strings.Contains(err.Error(), "nil context") {
		t.Fatalf("unexpected error: %v", err)
	}
}
