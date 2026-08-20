package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"flower/api/internal/platform/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	connStr := buildConnectionString(dbConfig)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	const (
		maxAttempts = 30
		delay       = time.Second
	)
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if err := db.Ping(); err == nil {
			return db, nil
		} else {
			lastErr = err
			time.Sleep(delay)
		}
	}

	if closeErr := db.Close(); closeErr != nil {
		return nil, fmt.Errorf("failed to ping database after %d attempts: %w (also failed to close handle: %v)", maxAttempts, lastErr, closeErr)
	}
	return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxAttempts, lastErr)
}

func WithTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	if db == nil {
		return fmt.Errorf("db: nil database handle")
	}
	if fn == nil {
		return fmt.Errorf("db: nil transaction function")
	}
	if ctx == nil {
		return fmt.Errorf("db: nil context")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db: begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("db: rollback failed after error (%v): %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db: commit transaction: %w", err)
	}

	return nil
}

func buildConnectionString(dbConfig config.DatabaseConfig) string {
	sslMode := strings.TrimSpace(dbConfig.SSLMode)
	if sslMode == "" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, sslMode)
}
