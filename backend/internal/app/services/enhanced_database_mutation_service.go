package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/pkg/progress"
	"github.com/flow/internal/pkg/websocket"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EnhancedDatabaseMutationService extends the basic service with progress tracking and WebSocket support
type EnhancedDatabaseMutationService struct {
	*DatabaseMutationService
	progressTracker *progress.Tracker
	hub             *websocket.Hub
}

// NewEnhancedDatabaseMutationService creates a new enhanced database mutation service
func NewEnhancedDatabaseMutationService(
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	coreDB *pgxpool.Pool,
	progressTracker *progress.Tracker,
	hub *websocket.Hub,
) *EnhancedDatabaseMutationService {
	baseService := NewDatabaseMutationService(dbRepo, projectRepo, coreDB)
	return &EnhancedDatabaseMutationService{
		DatabaseMutationService: baseService,
		progressTracker:         progressTracker,
		hub:                     hub,
	}
}

// CreateProjectDatabaseWithProgress creates a new database with progress tracking
func (s *EnhancedDatabaseMutationService) CreateProjectDatabaseWithProgress(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest, userID uuid.UUID) (*entities.DatabaseResponse, error) {
	// Define operation steps
	steps := []string{
		"Validating request",
		"Verifying project exists",
		"Creating database configuration",
		"Creating PostgreSQL database",
		"Testing connection",
		"Finalizing setup",
	}

	// Start progress tracking
	opCtx := s.progressTracker.NewOperationContext(userID, "create_database", "Creating project database", steps)

	// Step 1: Validate request
	opCtx.UpdateProgress("Validating request", "Validating database creation request", 0.1)
	if err := ValidateDatabaseCreateRequest(req); err != nil {
		opCtx.Fail(fmt.Errorf("validation failed: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Validating request")

	// Step 2: Verify project exists
	opCtx.UpdateProgress("Verifying project exists", "Checking if project exists", 0.2)
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		opCtx.Fail(fmt.Errorf("project not found: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Verifying project exists")

	// Step 3: Create database configuration
	opCtx.UpdateProgress("Creating database configuration", "Saving database configuration to core database", 0.3)
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

	if err := s.dbRepo.Create(ctx, database); err != nil {
		opCtx.Fail(fmt.Errorf("failed to save database configuration: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Creating database configuration")

	// Step 4: Create PostgreSQL database
	opCtx.UpdateProgress("Creating PostgreSQL database", "Creating the actual PostgreSQL database", 0.5)
	if err := s.createPostgreSQLDatabaseWithProgress(ctx, database, opCtx); err != nil {
		// Rollback: remove from core database if PostgreSQL creation fails
		s.dbRepo.Delete(ctx, database.ID)
		opCtx.Fail(fmt.Errorf("failed to create PostgreSQL database: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Creating PostgreSQL database")

	// Step 5: Test connection
	opCtx.UpdateProgress("Testing connection", "Verifying database connection works", 0.8)
	if err := s.testDatabaseConnectionWithProgress(ctx, database, opCtx); err != nil {
		opCtx.Fail(fmt.Errorf("database connection test failed: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Testing connection")

	// Step 6: Finalize setup
	opCtx.UpdateProgress("Finalizing setup", "Completing database setup", 0.9)
	response := database.ToResponse()
	opCtx.Complete(response)

	return &response, nil
}

// createPostgreSQLDatabaseWithProgress creates a PostgreSQL database with progress updates
func (s *EnhancedDatabaseMutationService) createPostgreSQLDatabaseWithProgress(ctx context.Context, database *entities.Database, opCtx *progress.OperationContext) error {
	// Connect to PostgreSQL server (not specific database)
	serverDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.SSLMode)

	opCtx.UpdateProgress("Creating PostgreSQL database", "Connecting to PostgreSQL server", 0.6)
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
	opCtx.UpdateProgress("Creating PostgreSQL database", "Executing CREATE DATABASE command", 0.7)
	createQuery := fmt.Sprintf("CREATE DATABASE %s", database.Database)
	if _, err := pool.Exec(ctx, createQuery); err != nil {
		return fmt.Errorf("failed to create database %s: %w", database.Database, err)
	}

	return nil
}

// testDatabaseConnectionWithProgress tests database connection with progress updates
func (s *EnhancedDatabaseMutationService) testDatabaseConnectionWithProgress(ctx context.Context, database *entities.Database, opCtx *progress.OperationContext) error {
	opCtx.UpdateProgress("Testing connection", "Creating test connection", 0.85)
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
	opCtx.UpdateProgress("Testing connection", "Pinging database", 0.9)
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	return nil
}

// CreateTableWithProgress creates a table with progress tracking
func (s *EnhancedDatabaseMutationService) CreateTableWithProgress(ctx context.Context, projectID uuid.UUID, tableReq *entities.TableCreateRequest, userID uuid.UUID) (*entities.TableResponse, error) {
	// Define operation steps
	steps := []string{
		"Validating table request",
		"Getting project database",
		"Creating table schema",
		"Executing CREATE TABLE",
		"Verifying table creation",
		"Updating metadata",
	}

	// Start progress tracking
	opCtx := s.progressTracker.NewOperationContext(userID, "create_table", "Creating table", steps)

	// Step 1: Validate table request
	opCtx.UpdateProgress("Validating table request", "Validating table creation request", 0.1)
	if err := ValidateTableCreateRequest(tableReq); err != nil {
		opCtx.Fail(fmt.Errorf("validation failed: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Validating table request")

	// Step 2: Get project database
	opCtx.UpdateProgress("Getting project database", "Retrieving project database configuration", 0.2)
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		opCtx.Fail(fmt.Errorf("failed to get project databases: %w", err))
		return nil, err
	}

	if len(databases) == 0 {
		opCtx.Fail(fmt.Errorf("no database found for project"))
		return nil, err
	}

	database := databases[0] // Use first database
	opCtx.CompleteStep("Getting project database")

	// Step 3: Create table schema
	opCtx.UpdateProgress("Creating table schema", "Preparing table creation SQL", 0.3)
	createTableSQL, err := s.generateCreateTableSQL(tableReq)
	if err != nil {
		opCtx.Fail(fmt.Errorf("failed to generate table SQL: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Creating table schema")

	// Step 4: Execute CREATE TABLE
	opCtx.UpdateProgress("Executing CREATE TABLE", "Creating table in database", 0.5)
	if err := s.executeCreateTable(ctx, database, createTableSQL, opCtx); err != nil {
		opCtx.Fail(fmt.Errorf("failed to create table: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Executing CREATE TABLE")

	// Step 5: Verify table creation
	opCtx.UpdateProgress("Verifying table creation", "Checking if table was created successfully", 0.8)
	if err := s.verifyTableCreation(ctx, database, tableReq.Name, opCtx); err != nil {
		opCtx.Fail(fmt.Errorf("table creation verification failed: %w", err))
		return nil, err
	}
	opCtx.CompleteStep("Verifying table creation")

	// Step 6: Update metadata
	opCtx.UpdateProgress("Updating metadata", "Saving table metadata", 0.9)

	// Convert schema to JSON string
	schemaJSON, err := json.Marshal(tableReq.Schema)
	if err != nil {
		opCtx.Fail(fmt.Errorf("failed to marshal schema: %w", err))
		return nil, err
	}

	tableResponse := &entities.TableResponse{
		ID:        uuid.New(),
		Name:      tableReq.Name,
		Schema:    string(schemaJSON),
		ProjectID: projectID,
	}
	opCtx.Complete(tableResponse)

	return tableResponse, nil
}

// generateCreateTableSQL generates SQL for creating a table
func (s *EnhancedDatabaseMutationService) generateCreateTableSQL(tableReq *entities.TableCreateRequest) (string, error) {
	// This is a simplified implementation
	// In a real implementation, you would parse the schema and generate proper SQL
	return fmt.Sprintf("CREATE TABLE %s (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), created_at TIMESTAMP DEFAULT NOW())", tableReq.Name), nil
}

// executeCreateTable executes the CREATE TABLE statement
func (s *EnhancedDatabaseMutationService) executeCreateTable(ctx context.Context, database *entities.Database, sql string, opCtx *progress.OperationContext) error {
	opCtx.UpdateProgress("Executing CREATE TABLE", "Connecting to database", 0.6)
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

	opCtx.UpdateProgress("Executing CREATE TABLE", "Running CREATE TABLE statement", 0.7)
	if _, err := pool.Exec(ctx, sql); err != nil {
		return fmt.Errorf("failed to execute CREATE TABLE: %w", err)
	}

	return nil
}

// verifyTableCreation verifies that the table was created successfully
func (s *EnhancedDatabaseMutationService) verifyTableCreation(ctx context.Context, database *entities.Database, tableName string, opCtx *progress.OperationContext) error {
	opCtx.UpdateProgress("Verifying table creation", "Checking table exists", 0.85)
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

	// Check if table exists
	checkQuery := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
	var exists bool
	if err := pool.QueryRow(ctx, checkQuery, tableName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to verify table creation: %w", err)
	}

	if !exists {
		return fmt.Errorf("table %s was not created", tableName)
	}

	return nil
}

// ValidateTableCreateRequest validates a table creation request
func ValidateTableCreateRequest(req *entities.TableCreateRequest) error {
	if req.Name == "" {
		return fmt.Errorf("table name is required")
	}
	if req.Schema == nil {
		return fmt.Errorf("table schema is required")
	}
	return nil
}
