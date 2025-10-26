package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var projectInitialSchema string

type ProjectMigrationRunner struct {
	connService *ConnectionService
	logger      *logging.Logger
}

func NewProjectMigrationRunner(connService *ConnectionService, logger *logging.Logger) *ProjectMigrationRunner {
	return &ProjectMigrationRunner{
		connService: connService,
		logger:      logger,
	}
}

func (r *ProjectMigrationRunner) InitializeProjectDatabase(ctx context.Context, database *entities.Database) error {
	r.logger.Info("Initializing project database schema", map[string]interface{}{
		"database_name": database.Database,
		"project_id":    database.ProjectID,
	})

	err := r.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		_, err := pool.Exec(ctx, projectInitialSchema)
		if err != nil {
			return fmt.Errorf("failed to execute initial schema: %w", err)
		}

		r.logger.Info("Successfully initialized project database schema", map[string]interface{}{
			"database_name": database.Database,
			"project_id":    database.ProjectID,
		})

		return nil
	})

	if err != nil {
		r.logger.Error("Failed to initialize project database schema", err, map[string]interface{}{
			"database_name": database.Database,
			"project_id":    database.ProjectID,
		})
		return err
	}

	return nil
}

func (r *ProjectMigrationRunner) GetProjectSchemaVersion(ctx context.Context, database *entities.Database) (int64, error) {
	var version int64

	err := r.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		query := `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`
		return pool.QueryRow(ctx, query).Scan(&version)
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get schema version: %w", err)
	}

	return version, nil
}

func (r *ProjectMigrationRunner) RunProjectMigration(ctx context.Context, database *entities.Database, version int64, name string, sql string) error {
	r.logger.Info("Running project migration", map[string]interface{}{
		"database_name": database.Database,
		"version":       version,
		"name":          name,
	})

	err := r.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx, sql)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`,
			version, name,
		)
		if err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		if err = tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		r.logger.Error("Failed to run project migration", err, map[string]interface{}{
			"database_name": database.Database,
			"version":       version,
			"name":          name,
		})
		return err
	}

	r.logger.Info("Successfully ran project migration", map[string]interface{}{
		"database_name": database.Database,
		"version":       version,
		"name":          name,
	})

	return nil
}
