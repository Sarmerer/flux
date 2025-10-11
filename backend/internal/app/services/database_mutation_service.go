package services

import (
	"context"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseMutationService handles database mutation operations
type DatabaseMutationService struct {
	dbRepo      repositories.DatabaseRepository
	projectRepo repositories.ProjectRepository
	coreDB      *pgxpool.Pool
}

// NewDatabaseMutationService creates a new DatabaseMutationService
func NewDatabaseMutationService(
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	coreDB *pgxpool.Pool,
) *DatabaseMutationService {
	return &DatabaseMutationService{
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
		coreDB:      coreDB,
	}
}

// CreateProjectDatabase creates a new isolated PostgreSQL database for a project
func (s *DatabaseMutationService) CreateProjectDatabase(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	// Validate request
	if err := ValidateDatabaseCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Create database entity
	database := &entities.Database{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Host:      req.Host,
		Port:      req.Port,
		Username:  req.Username,
		Password:  req.Password,
		Database:  req.Database,
		SSLMode:   req.SSLMode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save database configuration to core database
	if err := s.dbRepo.Create(ctx, database); err != nil {
		return nil, fmt.Errorf("failed to save database configuration: %w", err)
	}

	// Create the actual PostgreSQL database
	if err := s.createPostgreSQLDatabase(ctx, database); err != nil {
		// Rollback: remove from core database if PostgreSQL creation fails
		s.dbRepo.Delete(ctx, database.ID)
		return nil, fmt.Errorf("failed to create PostgreSQL database: %w", err)
	}

	response := database.ToResponse()
	return &response, nil
}

// UpdateProjectDatabase updates an existing project database configuration
func (s *DatabaseMutationService) UpdateProjectDatabase(ctx context.Context, databaseID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	// Validate request
	if err := ValidateDatabaseCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing database
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return nil, fmt.Errorf("database not found: %w", err)
	}

	// Update fields
	database.Name = req.Name
	database.Host = req.Host
	database.Port = req.Port
	database.Username = req.Username
	database.Password = req.Password
	database.Database = req.Database
	database.SSLMode = req.SSLMode
	database.UpdatedAt = time.Now()

	// Update in core database
	if err := s.dbRepo.Update(ctx, database); err != nil {
		return nil, fmt.Errorf("failed to update database configuration: %w", err)
	}

	response := database.ToResponse()
	return &response, nil
}

// DeleteProjectDatabase deletes a project database
func (s *DatabaseMutationService) DeleteProjectDatabase(ctx context.Context, databaseID uuid.UUID) error {
	// Get database configuration
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Drop the PostgreSQL database
	if err := s.dropPostgreSQLDatabase(ctx, database); err != nil {
		return fmt.Errorf("failed to drop PostgreSQL database: %w", err)
	}

	// Remove from core database
	if err := s.dbRepo.Delete(ctx, databaseID); err != nil {
		return fmt.Errorf("failed to remove database configuration: %w", err)
	}

	return nil
}

// TestDatabaseConnection tests connection to a project database
func (s *DatabaseMutationService) TestDatabaseConnection(ctx context.Context, databaseID uuid.UUID) error {
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return fmt.Errorf("database not found: %w", err)
	}

	// Create temporary connection to test
	dsn := database.GetConnectionString()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("invalid database configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	return nil
}

// createPostgreSQLDatabase creates a new PostgreSQL database
func (s *DatabaseMutationService) createPostgreSQLDatabase(ctx context.Context, database *entities.Database) error {
	// Connect to PostgreSQL server (not specific database)
	serverDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.SSLMode)

	config, err := pgxpool.ParseConfig(serverDSN)
	if err != nil {
		return fmt.Errorf("invalid server configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer pool.Close()

	// Create database
	createQuery := fmt.Sprintf("CREATE DATABASE %s", database.Database)
	if _, err := pool.Exec(ctx, createQuery); err != nil {
		return fmt.Errorf("failed to create database %s: %w", database.Database, err)
	}

	return nil
}

// dropPostgreSQLDatabase drops a PostgreSQL database
func (s *DatabaseMutationService) dropPostgreSQLDatabase(ctx context.Context, database *entities.Database) error {
	// Connect to PostgreSQL server (not specific database)
	serverDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.SSLMode)

	config, err := pgxpool.ParseConfig(serverDSN)
	if err != nil {
		return fmt.Errorf("invalid server configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer pool.Close()

	// Terminate existing connections to the database
	terminateQuery := fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, database.Database)
	pool.Exec(ctx, terminateQuery)

	// Drop database
	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", database.Database)
	if _, err := pool.Exec(ctx, dropQuery); err != nil {
		return fmt.Errorf("failed to drop database %s: %w", database.Database, err)
	}

	return nil
}

// GetProjectDatabases retrieves all databases for a project
func (s *DatabaseMutationService) GetProjectDatabases(ctx context.Context, projectID uuid.UUID) ([]*entities.DatabaseResponse, error) {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}

	var responses []*entities.DatabaseResponse
	for _, database := range databases {
		response := database.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}
