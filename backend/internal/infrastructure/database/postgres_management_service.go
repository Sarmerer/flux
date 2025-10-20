package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQLManagementService handles physical PostgreSQL database operations
// Responsible for CREATE DATABASE, DROP DATABASE, table creation, etc.
type PostgreSQLManagementService struct {
	connService *ConnectionService
}

// NewPostgreSQLManagementService creates a new PostgreSQL management service
func NewPostgreSQLManagementService(connService *ConnectionService) *PostgreSQLManagementService {
	return &PostgreSQLManagementService{
		connService: connService,
	}
}

// CreateDatabase creates a new PostgreSQL database
func (s *PostgreSQLManagementService) CreateDatabase(ctx context.Context, database *entities.Database) error {
	return s.connService.ExecuteWithPostgreSQLServer(ctx, database, func(pool *pgxpool.Pool) error {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", database.Database)
		if _, err := pool.Exec(ctx, createQuery); err != nil {
			return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to create database %s", database.Database))
		}
		log.Printf("Successfully created PostgreSQL database: %s", database.Database)
		return nil
	})
}

// DropDatabase drops a PostgreSQL database
func (s *PostgreSQLManagementService) DropDatabase(ctx context.Context, database *entities.Database) error {
	return s.connService.ExecuteWithPostgreSQLServer(ctx, database, func(pool *pgxpool.Pool) error {
		// Terminate existing connections to the database
		terminateQuery := fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, database.Database)

		if _, err := pool.Exec(ctx, terminateQuery); err != nil {
			// Log but don't fail - database might not exist or connections might already be closed
			log.Printf("Warning: Failed to terminate connections for database %s: %v", database.Database, err)
		}

		// Drop database
		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", database.Database)
		if _, err := pool.Exec(ctx, dropQuery); err != nil {
			// Check if error is due to database not existing
			if strings.Contains(err.Error(), "does not exist") {
				log.Printf("Database %s does not exist, skipping drop", database.Database)
				return nil
			}
			return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to drop database %s", database.Database))
		}

		log.Printf("Successfully dropped PostgreSQL database: %s", database.Database)
		return nil
	})
}

// TestConnection tests if a database connection works
func (s *PostgreSQLManagementService) TestConnection(ctx context.Context, database *entities.Database) error {
	return s.connService.TestConnection(ctx, database)
}

// CreateTable creates a table in a database
func (s *PostgreSQLManagementService) CreateTable(ctx context.Context, database *entities.Database, tableName string, schema string) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, schema); err != nil {
			return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to create table %s", tableName))
		}
		log.Printf("Successfully created table: %s", tableName)
		return nil
	})
}

// TableExists checks if a table exists in a database
func (s *PostgreSQLManagementService) TableExists(ctx context.Context, database *entities.Database, tableName string) (bool, error) {
	var exists bool
	err := s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		checkQuery := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
		if err := pool.QueryRow(ctx, checkQuery, tableName).Scan(&exists); err != nil {
			return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to check if table %s exists", tableName))
		}
		return nil
	})

	return exists, err
}

// DropTable drops a table from a database
func (s *PostgreSQLManagementService) DropTable(ctx context.Context, database *entities.Database, tableName string) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)
		if _, err := pool.Exec(ctx, dropQuery); err != nil {
			return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to drop table %s", tableName))
		}
		log.Printf("Successfully dropped table: %s", tableName)
		return nil
	})
}

// ExecuteQuery executes a raw SQL query (use with caution)
func (s *PostgreSQLManagementService) ExecuteQuery(ctx context.Context, database *entities.Database, query string, args ...interface{}) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			return errors.NewDatabaseError(err).WithDetails("failed to execute query")
		}
		return nil
	})
}
