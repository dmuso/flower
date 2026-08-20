package db

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"flower/api/internal/platform/config"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func getMigrationsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	if !ok {
		baseDir = "."
	}
	joined := filepath.Join(baseDir, "..", "..", "migrations")
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return filepath.Clean(joined)
	}
	return absPath
}

func newMigrator(dbConfig config.DatabaseConfig) (*migrate.Migrate, func(), error) {
	migrationsDir := getMigrationsDir()
	migrationConn, err := Connect(dbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect for migrations: %w", err)
	}
	driver, err := pgxmigrate.WithInstance(migrationConn, &pgxmigrate.Config{})
	if err != nil {
		if closeErr := migrationConn.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf("failed to create pgx migration driver: %w (also failed to close connection: %v)", err, closeErr)
		}
		return nil, nil, fmt.Errorf("failed to create pgx migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, dbConfig.Name, driver)
	if err != nil {
		if closeErr := migrationConn.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf("failed to initialise migrator: %w (also failed to close connection: %v)", err, closeErr)
		}
		return nil, nil, fmt.Errorf("failed to initialise migrator: %w", err)
	}
	cleanup := func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil && !errors.Is(sourceErr, migrate.ErrNilVersion) {
			zap.L().Warn("failed to close migration source cleanly",
				zap.String("component", "platform.db"),
				zap.String("operation", "migration-cleanup"),
				zap.Error(sourceErr),
			)
		}
		if databaseErr != nil && !errors.Is(databaseErr, migrate.ErrNilVersion) {
			zap.L().Warn("failed to close migration database cleanly",
				zap.String("component", "platform.db"),
				zap.String("operation", "migration-cleanup"),
				zap.Error(databaseErr),
			)
		}
		if err := migrationConn.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
			zap.L().Warn("failed to close migration connection",
				zap.String("component", "platform.db"),
				zap.String("operation", "migration-cleanup"),
				zap.Error(err),
			)
		}
	}
	return m, cleanup, nil
}

func Migrate(dbConfig config.DatabaseConfig) error {
	m, cleanup, err := newMigrator(dbConfig)
	if err != nil {
		return fmt.Errorf("new migrator error: %w", err)
	}
	defer cleanup()

	zap.L().Info("starting database migrations",
		zap.String("component", "platform.db"),
		zap.String("operation", "migrate"),
	)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	zap.L().Info("database migrations completed successfully",
		zap.String("component", "platform.db"),
		zap.String("operation", "migrate"),
	)
	return nil
}

func Rollback(dbConfig config.DatabaseConfig) error {
	m, cleanup, err := newMigrator(dbConfig)
	if err != nil {
		return fmt.Errorf("new migrator error: %w", err)
	}
	defer cleanup()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	zap.L().Info("database rollback completed successfully",
		zap.String("component", "platform.db"),
		zap.String("operation", "rollback"),
	)
	return nil
}

func Force(dbConfig config.DatabaseConfig, version int) error {
	m, cleanup, err := newMigrator(dbConfig)
	if err != nil {
		return fmt.Errorf("new migrator error: %w", err)
	}
	defer cleanup()

	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	zap.L().Info("database migration forced",
		zap.String("component", "platform.db"),
		zap.String("operation", "force"),
		zap.Int("version", version),
	)
	return nil
}

func Version(dbConfig config.DatabaseConfig) (string, error) {
	m, cleanup, err := newMigrator(dbConfig)
	if err != nil {
		return "", fmt.Errorf("new migrator error: %w", err)
	}
	defer cleanup()

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return "0", nil
		}
		return "", fmt.Errorf("failed to get migration version: %w", err)
	}
	status := fmt.Sprintf("%d", version)
	if dirty {
		status += " (dirty)"
	}
	return strings.TrimSpace(status), nil
}
