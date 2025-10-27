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

var migrationFiles embed.FS

type SQLMigration struct {
	Version int64
	Name    string
	UpSQL   string
	DownSQL string
}

func LoadMigrationsFromSQL() ([]*SQLMigration, error) {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("Failed to read migrations directory: %w", err)
	}

	migrationMap := make(map[int64]*SQLMigration)
	filePattern := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		matches := filePattern.FindStringSubmatch(filename)
		if matches == nil {
			continue
		}

		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid version in file %s: %w", filename, err)
		}

		name := matches[2]
		direction := matches[3]

		content, err := migrationFiles.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read migration file %s: %w", filename, err)
		}

		migration, exists := migrationMap[version]
		if !exists {
			migration = &SQLMigration{
				Version: version,
				Name:    name,
			}
			migrationMap[version] = migration
		}

		if direction == "up" {
			migration.UpSQL = string(content)
		} else {
			migration.DownSQL = string(content)
		}
	}

	migrations := make([]*SQLMigration, 0, len(migrationMap))
	for _, m := range migrationMap {
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

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
