package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectConnectionResolver struct {
	connectionService *ConnectionService
	dbRepository      repositories.DatabaseRepository
	cache             map[uuid.UUID]*pgxpool.Pool
	mutex             sync.RWMutex
}

func NewProjectConnectionResolver(
	connectionService *ConnectionService,
	dbRepository repositories.DatabaseRepository,
) *ProjectConnectionResolver {
	return &ProjectConnectionResolver{
		connectionService: connectionService,
		dbRepository:      dbRepository,
		cache:             make(map[uuid.UUID]*pgxpool.Pool),
	}
}

func (r *ProjectConnectionResolver) GetProjectDB(ctx context.Context, projectID uuid.UUID) (*pgxpool.Pool, error) {
	r.mutex.RLock()
	pool, exists := r.cache[projectID]
	r.mutex.RUnlock()

	if exists && pool != nil {
		return pool, nil
	}

	databases, err := r.dbRepository.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database config for project %s: %w", projectID, err)
	}

	if len(databases) == 0 {
		return nil, fmt.Errorf("no database configuration found for project %s", projectID)
	}

	database := databases[0]

	connectionString := database.GetConnectionString()
	newPool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool for project %s: %w", projectID, err)
	}

	if err := newPool.Ping(ctx); err != nil {
		newPool.Close()
		return nil, fmt.Errorf("failed to ping project database for project %s: %w", projectID, err)
	}

	r.mutex.Lock()
	r.cache[projectID] = newPool
	r.mutex.Unlock()

	return newPool, nil
}

func (r *ProjectConnectionResolver) GetDatabaseEntity(ctx context.Context, projectID uuid.UUID) (*entities.Database, error) {
	databases, err := r.dbRepository.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database config for project %s: %w", projectID, err)
	}

	if len(databases) == 0 {
		return nil, fmt.Errorf("no database configuration found for project %s", projectID)
	}

	return databases[0], nil
}

func (r *ProjectConnectionResolver) ExecuteWithProjectDB(ctx context.Context, projectID uuid.UUID, fn func(*pgxpool.Pool) error) error {
	pool, err := r.GetProjectDB(ctx, projectID)
	if err != nil {
		return err
	}
	return fn(pool)
}

func (r *ProjectConnectionResolver) InvalidateCache(projectID uuid.UUID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if pool, exists := r.cache[projectID]; exists {
		pool.Close()
		delete(r.cache, projectID)
	}
}

func (r *ProjectConnectionResolver) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for projectID, pool := range r.cache {
		pool.Close()
		delete(r.cache, projectID)
	}
}
