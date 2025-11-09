package migrations

import (
	"context"
	"fmt"
	"sort"

	"github.com/flow/internal/infrastructure/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Runner struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewRunner(db *pgxpool.Pool, logger *logging.Logger) *Runner {
	return &Runner{
		db:     db,
		logger: logger,
	}
}

func (r *Runner) Up(ctx context.Context) error {

	if err := CreateMigrationsTable(ctx, r.db); err != nil {
		return fmt.Errorf("Failed to create migrations table: %w", err)
	}

	applied, err := GetAppliedMigrations(ctx, r.db)
	if err != nil {
		return fmt.Errorf("Failed to get applied migrations: %w", err)
	}

	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("Failed to load migrations: %w", err)
	}

	migrations := make([]*Migration, 0, len(sqlMigrations))
	for _, sm := range sqlMigrations {
		migrations = append(migrations, sm.ToMigration())
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	pendingCount := 0
	for _, m := range migrations {
		if !applied[m.Version] {
			pendingCount++
		}
	}

	if pendingCount == 0 {
		r.logger.Info("No pending migrations", map[string]interface{}{
			"total_migrations": len(migrations),
		})
		return nil
	}

	r.logger.Info(fmt.Sprintf("Found %d pending migrations", pendingCount), map[string]interface{}{
		"pending_count": pendingCount,
		"total_count":   len(migrations),
	})

	for _, m := range migrations {
		if applied[m.Version] {
			r.logger.Debug(fmt.Sprintf("Migration %d already applied, skipping", m.Version), map[string]interface{}{
				"version": m.Version,
				"name":    m.Name,
			})
			continue
		}

		r.logger.Info(fmt.Sprintf("Running migration %d: %s", m.Version, m.Name), map[string]interface{}{
			"version":     m.Version,
			"name":        m.Name,
			"description": m.Description,
		})

		if err := m.Up(ctx, r.db); err != nil {
			r.logger.Error(fmt.Sprintf("Migration %d failed", m.Version), err, map[string]interface{}{
				"version": m.Version,
				"name":    m.Name,
			})
			return fmt.Errorf("Failed to run migration %d (%s): %w", m.Version, m.Name, err)
		}

		if err := RecordMigration(ctx, r.db, m.Version, m.Name); err != nil {
			r.logger.Error("Failed to record migration", err, map[string]interface{}{
				"version": m.Version,
			})
			return fmt.Errorf("Failed to record migration %d: %w", m.Version, err)
		}

		r.logger.Info(fmt.Sprintf("Migration %d applied successfully", m.Version), map[string]interface{}{
			"version": m.Version,
			"name":    m.Name,
		})
	}

	r.logger.Info("All migrations completed successfully", map[string]interface{}{
		"applied_count": pendingCount,
	})
	return nil
}

func (r *Runner) Down(ctx context.Context) error {

	latestVersion, err := GetLatestVersion(ctx, r.db)
	if err != nil {
		return fmt.Errorf("Failed to get latest version: %w", err)
	}

	if latestVersion == 0 {
		r.logger.Info("No migrations to rollback", nil)
		return nil
	}

	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("Failed to load migrations: %w", err)
	}

	var migration *Migration
	for _, sm := range sqlMigrations {
		if sm.Version == latestVersion {
			migration = sm.ToMigration()
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %d not found in registry", latestVersion)
	}

	if migration.Down == nil {
		return fmt.Errorf("migration %d (%s) does not have a down function", migration.Version, migration.Name)
	}

	r.logger.Info(fmt.Sprintf("Rolling back migration %d: %s", migration.Version, migration.Name), map[string]interface{}{
		"version": migration.Version,
		"name":    migration.Name,
	})

	if err := migration.Down(ctx, r.db); err != nil {
		r.logger.Error("Failed to rollback migration", err, map[string]interface{}{
			"version": migration.Version,
		})
		return fmt.Errorf("Failed to rollback migration %d: %w", migration.Version, err)
	}

	if err := RemoveMigration(ctx, r.db, migration.Version); err != nil {
		return fmt.Errorf("Failed to remove migration record: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Migration %d rolled back successfully", migration.Version), map[string]interface{}{
		"version": migration.Version,
		"name":    migration.Name,
	})
	return nil
}

func (r *Runner) Status(ctx context.Context) error {

	if err := CreateMigrationsTable(ctx, r.db); err != nil {
		return fmt.Errorf("Failed to create migrations table: %w", err)
	}

	applied, err := GetAppliedMigrations(ctx, r.db)
	if err != nil {
		return fmt.Errorf("Failed to get applied migrations: %w", err)
	}

	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("Failed to load migrations: %w", err)
	}

	migrations := make([]*Migration, 0, len(sqlMigrations))
	for _, sm := range sqlMigrations {
		migrations = append(migrations, sm.ToMigration())
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	r.logger.Info("Migration Status:", nil)
	r.logger.Info("================", nil)
	for _, m := range migrations {
		status := "pending"
		if applied[m.Version] {
			status = "applied"
		}
		r.logger.Info(fmt.Sprintf("[%s] Version %d: %s", status, m.Version, m.Name), map[string]interface{}{
			"version": m.Version,
			"name":    m.Name,
			"status":  status,
		})
	}

	return nil
}
