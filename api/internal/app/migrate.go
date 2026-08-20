package app

import (
	"fmt"
	"os"

	"flower/api/internal/platform/config"
	"flower/api/internal/platform/db"

	"go.uber.org/zap"
)

func RunMigrate() {
	cfg, err := config.Load()
	if err != nil {
		zap.L().Error("failed to load configuration",
			zap.String("component", "app"),
			zap.String("operation", "run-migrate"),
			zap.Error(err),
		)
		os.Exit(1)
	}

	if err := db.Migrate(cfg.Database); err != nil {
		zap.L().Error("failed to run migrations",
			zap.String("component", "app"),
			zap.String("operation", "run-migrate"),
			zap.Error(err),
		)
		os.Exit(1)
	}
}

func RunRollback() {
	cfg, err := config.Load()
	if err != nil {
		zap.L().Error("failed to load configuration",
			zap.String("component", "app"),
			zap.String("operation", "run-rollback"),
			zap.Error(err),
		)
		os.Exit(1)
	}

	if err := db.Rollback(cfg.Database); err != nil {
		zap.L().Error("failed to rollback migrations",
			zap.String("component", "app"),
			zap.String("operation", "run-rollback"),
			zap.Error(err),
		)
		os.Exit(1)
	}
}

func RunForce(version string) {
	cfg, err := config.Load()
	if err != nil {
		zap.L().Error("failed to load configuration",
			zap.String("component", "app"),
			zap.String("operation", "run-force"),
			zap.Error(err),
		)
		os.Exit(1)
	}

	var versionInt int
	if _, err := fmt.Sscanf(version, "%d", &versionInt); err != nil {
		zap.L().Error("invalid migration version number",
			zap.String("component", "app"),
			zap.String("operation", "run-force"),
			zap.String("version", version),
			zap.Error(err),
		)
		os.Exit(1)
	}

	if err := db.Force(cfg.Database, versionInt); err != nil {
		zap.L().Error("failed to force migration version",
			zap.String("component", "app"),
			zap.String("operation", "run-force"),
			zap.Error(err),
		)
		os.Exit(1)
	}
}
