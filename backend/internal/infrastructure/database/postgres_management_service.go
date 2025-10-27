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

type PostgreSQLManagementService struct {
	connService *ConnectionService
}

func NewPostgreSQLManagementService(connService *ConnectionService) *PostgreSQLManagementService {
	return &PostgreSQLManagementService{
		connService: connService,
	}
}

func (s *PostgreSQLManagementService) CreateDatabase(ctx context.Context, database *entities.Database) error {
	return s.connService.ExecuteWithPostgreSQLServer(ctx, database, func(pool *pgxpool.Pool) error {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", database.Database)
		if _, err := pool.Exec(ctx, createQuery); err != nil {
			return errors.NewDatabaseError(fmt.Sprintf("Failed to create database %s", database.Database), err)
		}
		log.Printf("Successfully created PostgreSQL database: %s", database.Database)
		return nil
	})
}

func (s *PostgreSQLManagementService) DropDatabase(ctx context.Context, database *entities.Database) error {
	return s.connService.ExecuteWithPostgreSQLServer(ctx, database, func(pool *pgxpool.Pool) error {

		terminateQuery := fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, database.Database)

		if _, err := pool.Exec(ctx, terminateQuery); err != nil {

			log.Printf("Warning: Failed to terminate connections for database %s: %v", database.Database, err)
		}

		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", database.Database)
		if _, err := pool.Exec(ctx, dropQuery); err != nil {

			if strings.Contains(err.Error(), "does not exist") {
				log.Printf("Database %s does not exist, skipping drop", database.Database)
				return nil
			}
			return errors.NewDatabaseError(fmt.Sprintf("Failed to drop database %s", database.Database), err)
		}

		log.Printf("Successfully dropped PostgreSQL database: %s", database.Database)
		return nil
	})
}

func (s *PostgreSQLManagementService) TestConnection(ctx context.Context, database *entities.Database) error {
	return s.connService.TestConnection(ctx, database)
}

func (s *PostgreSQLManagementService) CreateTable(ctx context.Context, database *entities.Database, tableName string, schema string) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, schema); err != nil {
			return errors.NewDatabaseError(fmt.Sprintf("Failed to create table %s", tableName), err)
		}
		log.Printf("Successfully created table: %s", tableName)
		return nil
	})
}

func (s *PostgreSQLManagementService) TableExists(ctx context.Context, database *entities.Database, tableName string) (bool, error) {
	var exists bool
	err := s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		checkQuery := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
		if err := pool.QueryRow(ctx, checkQuery, tableName).Scan(&exists); err != nil {
			return errors.NewDatabaseError(fmt.Sprintf("Failed to check if table %s exists", tableName), err)
		}
		return nil
	})

	return exists, err
}

func (s *PostgreSQLManagementService) DropTable(ctx context.Context, database *entities.Database, tableName string) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)
		if _, err := pool.Exec(ctx, dropQuery); err != nil {
			return errors.NewDatabaseError(fmt.Sprintf("Failed to drop table %s", tableName), err)
		}
		log.Printf("Successfully dropped table: %s", tableName)
		return nil
	})
}

func (s *PostgreSQLManagementService) ExecuteQuery(ctx context.Context, database *entities.Database, query string, args ...interface{}) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			return errors.NewDatabaseError("Failed to execute query", err)
		}
		return nil
	})
}
