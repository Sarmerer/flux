package database

import (
	"context"

	"github.com/flow/internal/infrastructure/database/migrations"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AutoMigrate(ctx context.Context, pool *pgxpool.Pool, logger *logging.Logger) error {
	logger.Info("Starting database migrations...", map[string]interface{}{
		"operation": "auto_migrate",
	})

	runner := migrations.NewRunner(pool, logger)

	if err := runner.Up(ctx); err != nil {
		logger.Error("Failed to run migrations", err, map[string]interface{}{
			"operation": "auto_migrate",
		})
		return err
	}

	return nil
}

func ShowMigrationStatus(ctx context.Context, pool *pgxpool.Pool, logger *logging.Logger) error {
	runner := migrations.NewRunner(pool, logger)
	return runner.Status(ctx)
}
