package database

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionService struct {
	coreDB *pgxpool.Pool
}

func NewConnectionService(coreDB *pgxpool.Pool) *ConnectionService {
	return &ConnectionService{
		coreDB: coreDB,
	}
}

func (s *ConnectionService) GetCoreDB() *pgxpool.Pool {
	return s.coreDB
}

func (s *ConnectionService) ExecuteWithProjectDB(ctx context.Context, database *entities.Database, fn func(*pgxpool.Pool) error) error {
	pool, err := s.createPool(ctx, database.GetConnectionString())
	if err != nil {
		return errors.NewDatabaseError(err).WithDetails(fmt.Sprintf("failed to connect to database %s", database.Name))
	}
	defer pool.Close()

	if err := fn(pool); err != nil {
		return err
	}

	return nil
}

func (s *ConnectionService) ExecuteWithPostgreSQLServer(ctx context.Context, database *entities.Database, fn func(*pgxpool.Pool) error) error {
	serverDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.SSLMode)

	pool, err := s.createPool(ctx, serverDSN)
	if err != nil {
		return errors.NewDatabaseError(err).WithDetails("failed to connect to PostgreSQL server")
	}
	defer pool.Close()

	if err := fn(pool); err != nil {
		return err
	}

	return nil
}

func (s *ConnectionService) TestConnection(ctx context.Context, database *entities.Database) error {
	return s.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if err := pool.Ping(ctx); err != nil {
			return errors.NewAPIError(errors.ErrCodeConnectionError, "database connection test failed").WithDetails(err.Error())
		}
		return nil
	})
}

func (s *ConnectionService) CreatePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return s.createPool(ctx, dsn)
}

func (s *ConnectionService) createPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
