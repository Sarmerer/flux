package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migration represents a database migration
type Migration struct {
	Version     int64
	Name        string
	Up          func(ctx context.Context, db *pgxpool.Pool) error
	Down        func(ctx context.Context, db *pgxpool.Pool) error
	Description string
}

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID        int64     `db:"id"`
	Version   int64     `db:"version"`
	Name      string    `db:"name"`
	AppliedAt time.Time `db:"applied_at"`
}

// Registry holds all registered migrations
type Registry struct {
	migrations []*Migration
}

// NewRegistry creates a new migration registry
func NewRegistry() *Registry {
	return &Registry{
		migrations: make([]*Migration, 0),
	}
}

// Register adds a migration to the registry
func (r *Registry) Register(m *Migration) {
	r.migrations = append(r.migrations, m)
}

// GetMigrations returns all registered migrations sorted by version
func (r *Registry) GetMigrations() []*Migration {
	return r.migrations
}

// CreateMigrationsTable creates the schema_migrations table if it doesn't exist
func CreateMigrationsTable(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version BIGINT NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_schema_migrations_version ON schema_migrations(version);
	`
	_, err := db.Exec(ctx, query)
	return err
}

// GetAppliedMigrations returns all applied migrations
func GetAppliedMigrations(ctx context.Context, db *pgxpool.Pool) (map[int64]bool, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// RecordMigration records a migration as applied
func RecordMigration(ctx context.Context, db *pgxpool.Pool, version int64, name string) error {
	query := `
		INSERT INTO schema_migrations (version, name, applied_at)
		VALUES ($1, $2, NOW())
	`
	_, err := db.Exec(ctx, query, version, name)
	return err
}

// RemoveMigration removes a migration record
func RemoveMigration(ctx context.Context, db *pgxpool.Pool, version int64) error {
	query := "DELETE FROM schema_migrations WHERE version = $1"
	_, err := db.Exec(ctx, query, version)
	return err
}

// GetLatestVersion returns the latest applied migration version
func GetLatestVersion(ctx context.Context, db *pgxpool.Pool) (int64, error) {
	query := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	var version int64
	err := db.QueryRow(ctx, query).Scan(&version)
	return version, err
}
