package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection creates a new database connection
func NewConnection(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return pool, nil
}

// Migrate runs database migrations
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		createUsersTable,
		createProjectsTable,
		createDatabasesTable,
		createTablesTable,
	}

	for i, migration := range migrations {
		log.Printf("Running migration %d/%d", i+1, len(migrations))
		if _, err := pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`

const createProjectsTable = `
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    database_url TEXT,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_owner_id ON projects(owner_id);
CREATE INDEX IF NOT EXISTS idx_projects_api_key ON projects(api_key);
`

const createDatabasesTable = `
CREATE TABLE IF NOT EXISTS databases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    database VARCHAR(255) NOT NULL,
    ssl_mode VARCHAR(50) DEFAULT 'disable',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_databases_project_id ON databases(project_id);
`

const createTablesTable = `
CREATE TABLE IF NOT EXISTS tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    schema JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tables_project_id ON tables(project_id);
CREATE INDEX IF NOT EXISTS idx_tables_name ON tables(name);
`
