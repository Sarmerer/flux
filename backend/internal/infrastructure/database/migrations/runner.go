package migrations

import (
	"context"
	"fmt"
	"sort"

	"github.com/flow/internal/infrastructure/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Runner handles migration execution with structured logging
type Runner struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRunner creates a new migration runner
func NewRunner(db *pgxpool.Pool, logger *logging.Logger) *Runner {
	return &Runner{
		db:     db,
		logger: logger,
	}
}

// Up runs all pending migrations
func (r *Runner) Up(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	if err := CreateMigrationsTable(ctx, r.db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := GetAppliedMigrations(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Load migrations from SQL files
	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Convert to Migration objects
	migrations := make([]*Migration, 0, len(sqlMigrations))
	for _, sm := range sqlMigrations {
		migrations = append(migrations, sm.ToMigration())
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Count pending migrations
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

	// Run pending migrations
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

		// Begin transaction
		tx, err := r.db.Begin(ctx)
		if err != nil {
			r.logger.Error("Failed to begin transaction", err, map[string]interface{}{
				"version": m.Version,
			})
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		// Run the migration
		if err := m.Up(ctx, r.db); err != nil {
			tx.Rollback(ctx)
			r.logger.Error(fmt.Sprintf("Migration %d failed", m.Version), err, map[string]interface{}{
				"version": m.Version,
				"name":    m.Name,
			})
			return fmt.Errorf("failed to run migration %d (%s): %w", m.Version, m.Name, err)
		}

		// Record the migration
		if err := RecordMigration(ctx, r.db, m.Version, m.Name); err != nil {
			tx.Rollback(ctx)
			r.logger.Error("Failed to record migration", err, map[string]interface{}{
				"version": m.Version,
			})
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			r.logger.Error("Failed to commit migration", err, map[string]interface{}{
				"version": m.Version,
			})
			return fmt.Errorf("failed to commit migration %d: %w", m.Version, err)
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

// Down rolls back the last migration
func (r *Runner) Down(ctx context.Context) error {
	// Get latest applied migration
	latestVersion, err := GetLatestVersion(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	if latestVersion == 0 {
		r.logger.Info("No migrations to rollback", nil)
		return nil
	}

	// Load migrations from SQL files
	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Find the migration
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

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Run the down migration
	if err := migration.Down(ctx, r.db); err != nil {
		tx.Rollback(ctx)
		r.logger.Error("Failed to rollback migration", err, map[string]interface{}{
			"version": migration.Version,
		})
		return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
	}

	// Remove migration record
	if err := RemoveMigration(ctx, r.db, migration.Version); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Migration %d rolled back successfully", migration.Version), map[string]interface{}{
		"version": migration.Version,
		"name":    migration.Name,
	})
	return nil
}

// Status shows migration status
func (r *Runner) Status(ctx context.Context) error {
	// Create migrations table if needed
	if err := CreateMigrationsTable(ctx, r.db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := GetAppliedMigrations(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Load migrations from SQL files
	sqlMigrations, err := LoadMigrationsFromSQL()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Convert to Migration objects
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
