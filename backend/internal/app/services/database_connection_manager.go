package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConnectionManager manages connections to project databases
type DatabaseConnectionManager struct {
	dbRepo      repositories.DatabaseRepository
	connections map[uuid.UUID]*pgxpool.Pool
	mutex       sync.RWMutex
}

// NewDatabaseConnectionManager creates a new DatabaseConnectionManager
func NewDatabaseConnectionManager(dbRepo repositories.DatabaseRepository) *DatabaseConnectionManager {
	return &DatabaseConnectionManager{
		dbRepo:      dbRepo,
		connections: make(map[uuid.UUID]*pgxpool.Pool),
	}
}

// GetConnection gets or creates a connection to a project database
func (m *DatabaseConnectionManager) GetConnection(ctx context.Context, projectID uuid.UUID) (*pgxpool.Pool, error) {
	m.mutex.RLock()
	if pool, exists := m.connections[projectID]; exists {
		m.mutex.RUnlock()
		// Test connection
		if err := pool.Ping(ctx); err == nil {
			return pool, nil
		}
		// Connection is stale, remove it
		m.mutex.Lock()
		delete(m.connections, projectID)
		m.mutex.Unlock()
	} else {
		m.mutex.RUnlock()
	}

	// Create new connection
	return m.createConnection(ctx, projectID)
}

// createConnection creates a new connection to a project database
func (m *DatabaseConnectionManager) createConnection(ctx context.Context, projectID uuid.UUID) (*pgxpool.Pool, error) {
	// Get project database configuration
	databases, err := m.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return nil, fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0] // Use first database for now

	// Create connection
	dsn := database.GetConnectionString()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Store connection
	m.mutex.Lock()
	m.connections[projectID] = pool
	m.mutex.Unlock()

	return pool, nil
}

// CloseConnection closes a specific project database connection
func (m *DatabaseConnectionManager) CloseConnection(projectID uuid.UUID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if pool, exists := m.connections[projectID]; exists {
		pool.Close()
		delete(m.connections, projectID)
	}

	return nil
}

// CloseAllConnections closes all project database connections
func (m *DatabaseConnectionManager) CloseAllConnections() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for projectID, pool := range m.connections {
		pool.Close()
		delete(m.connections, projectID)
	}
}

// GetConnectionStats returns statistics about active connections
func (m *DatabaseConnectionManager) GetConnectionStats() map[uuid.UUID]ConnectionStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[uuid.UUID]ConnectionStats)
	for projectID, pool := range m.connections {
		stats[projectID] = ConnectionStats{
			TotalConns:        pool.Stat().TotalConns(),
			AcquiredConns:     pool.Stat().AcquiredConns(),
			ConstructingConns: pool.Stat().ConstructingConns(),
			IdleConns:         pool.Stat().IdleConns(),
		}
	}

	return stats
}

// ConnectionStats represents connection pool statistics
type ConnectionStats struct {
	TotalConns        int32 `json:"total_conns"`
	AcquiredConns     int32 `json:"acquired_conns"`
	ConstructingConns int32 `json:"constructing_conns"`
	IdleConns         int32 `json:"idle_conns"`
}

// ExecuteInProjectDatabase executes a function with a project database connection
func (m *DatabaseConnectionManager) ExecuteInProjectDatabase(ctx context.Context, projectID uuid.UUID, fn func(*pgxpool.Pool) error) error {
	pool, err := m.GetConnection(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	return fn(pool)
}

// ExecuteTransactionInProjectDatabase executes a function within a transaction
func (m *DatabaseConnectionManager) ExecuteTransactionInProjectDatabase(ctx context.Context, projectID uuid.UUID, fn func(pgx.Tx) error) error {
	return m.ExecuteInProjectDatabase(ctx, projectID, func(pool *pgxpool.Pool) error {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		if err := fn(tx); err != nil {
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})
}

// HealthCheck checks the health of all project database connections
func (m *DatabaseConnectionManager) HealthCheck(ctx context.Context) map[uuid.UUID]error {
	m.mutex.RLock()
	connections := make(map[uuid.UUID]*pgxpool.Pool)
	for projectID, pool := range m.connections {
		connections[projectID] = pool
	}
	m.mutex.RUnlock()

	results := make(map[uuid.UUID]error)
	for projectID, pool := range connections {
		if err := pool.Ping(ctx); err != nil {
			results[projectID] = err
			// Remove stale connection
			m.CloseConnection(projectID)
		} else {
			results[projectID] = nil
		}
	}

	return results
}
