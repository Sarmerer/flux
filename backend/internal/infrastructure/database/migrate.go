package database

import (
	"context"

	"github.com/flow/internal/infrastructure/database/migrations"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AutoMigrate runs database migrations automatically on startup
// This is the production-ready way to run migrations
// Uses SQL files and project logger
func AutoMigrate(ctx context.Context, pool *pgxpool.Pool, logger *logging.Logger) error {
	logger.Info("Starting database migrations...", map[string]interface{}{
		"operation": "auto_migrate",
	})

	// Create runner - it loads migrations from SQL files automatically
	runner := migrations.NewRunner(pool, logger)

	// Run migrations
	if err := runner.Up(ctx); err != nil {
		logger.Error("Failed to run migrations", err, map[string]interface{}{
			"operation": "auto_migrate",
		})
		return err
	}

	return nil
}

// ShowMigrationStatus displays migration status using project logger
func ShowMigrationStatus(ctx context.Context, pool *pgxpool.Pool, logger *logging.Logger) error {
	runner := migrations.NewRunner(pool, logger)
	return runner.Status(ctx)
}
