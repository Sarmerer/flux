package migrations

import (
	"context"
	"embed"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var migrationFiles embed.FS

// SQLMigration represents a SQL-based migration
type SQLMigration struct {
	Version int64
	Name    string
	UpSQL   string
	DownSQL string
}

// LoadMigrationsFromSQL loads all migrations from embedded SQL files
func LoadMigrationsFromSQL() ([]*SQLMigration, error) {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Parse migration files
	migrationMap := make(map[int64]*SQLMigration)
	filePattern := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		matches := filePattern.FindStringSubmatch(filename)
		if matches == nil {
			continue // Skip non-migration files
		}

		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid version in file %s: %w", filename, err)
		}

		name := matches[2]
		direction := matches[3] // "up" or "down"

		// Read SQL content
		content, err := migrationFiles.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Get or create migration entry
		migration, exists := migrationMap[version]
		if !exists {
			migration = &SQLMigration{
				Version: version,
				Name:    name,
			}
			migrationMap[version] = migration
		}

		// Set SQL based on direction
		if direction == "up" {
			migration.UpSQL = string(content)
		} else {
			migration.DownSQL = string(content)
		}
	}

	// Convert map to sorted slice
	migrations := make([]*SQLMigration, 0, len(migrationMap))
	for _, m := range migrationMap {
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// ToMigration converts SQLMigration to Migration with executable functions
func (sm *SQLMigration) ToMigration() *Migration {
	return &Migration{
		Version:     sm.Version,
		Name:        sm.Name,
		Description: fmt.Sprintf("Migration from %s", sm.Name),
		Up: func(ctx context.Context, db *pgxpool.Pool) error {
			if sm.UpSQL == "" {
				return fmt.Errorf("no up SQL for migration %d", sm.Version)
			}
			_, err := db.Exec(ctx, sm.UpSQL)
			return err
		},
		Down: func(ctx context.Context, db *pgxpool.Pool) error {
			if sm.DownSQL == "" {
				return fmt.Errorf("no down SQL for migration %d", sm.Version)
			}
			_, err := db.Exec(ctx, sm.DownSQL)
			return err
		},
	}
}
