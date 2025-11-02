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
