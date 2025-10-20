package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"

	"github.com/google/uuid"
)

// DatabaseService handles database configuration operations with optional progress tracking
// This service is focused on database CONFIGURATION, not physical database operations
type DatabaseService struct {
	dbRepo      repositories.DatabaseRepository
	projectRepo repositories.ProjectRepository
	pgService   *database.PostgreSQLManagementService
	connService *database.ConnectionService
	progressTracker *progress.Tracker
	hub         *realtime.Hub
}

// NewDatabaseService creates a new database service
func NewDatabaseService(
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	pgService *database.PostgreSQLManagementService,
	connService *database.ConnectionService,
	progressTracker *progress.Tracker,
	hub *realtime.Hub,
) *DatabaseService {
	return &DatabaseService{
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
		pgService:   pgService,
		connService: connService,
		progressTracker: progressTracker,
		hub:         hub,
	}
}

// CreateProjectDatabase creates a new isolated PostgreSQL database for a project (without progress tracking)
func (s *DatabaseService) CreateProjectDatabase(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	// Validate request
	if err := ValidateDatabaseCreateRequest(req); err != nil {
		return nil, errors.NewValidationError("validation failed").WithDetails(err.Error())
	}

	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewNotFoundError("Project")
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
		return nil, errors.NewDatabaseError(err).WithDetails("failed to save database configuration")
	}

	// Create the actual PostgreSQL database using the management service
	if err := s.pgService.CreateDatabase(ctx, database); err != nil {
		// Rollback: remove from core database if PostgreSQL creation fails
		s.dbRepo.Delete(ctx, database.ID)
		return nil, err
	}

	response := database.ToResponse()
	return &response, nil
}

// CreateProjectDatabaseWithProgress creates a new database with progress tracking
func (s *DatabaseService) CreateProjectDatabaseWithProgress(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest, userID uuid.UUID) (*entities.DatabaseResponse, error) {
	// If progress tracker is not available, fall back to simple version
	if s.progressTracker == nil {
		return s.CreateProjectDatabase(ctx, projectID, req)
	}

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
		opCtx.Fail(errors.NewValidationError("validation failed").WithDetails(err.Error()))
		return nil, err
	}
	opCtx.CompleteStep("Validating request")

	// Step 2: Verify project exists
	opCtx.UpdateProgress("Verifying project exists", "Checking if project exists", 0.2)
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		opCtx.Fail(errors.NewNotFoundError("Project"))
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
		opCtx.Fail(errors.NewDatabaseError(err).WithDetails("failed to save database configuration"))
		return nil, err
	}
	opCtx.CompleteStep("Creating database configuration")

	// Step 4: Create PostgreSQL database
	opCtx.UpdateProgress("Creating PostgreSQL database", "Creating the actual PostgreSQL database", 0.5)
	if err := s.pgService.CreateDatabase(ctx, database); err != nil {
		// Rollback: remove from core database if PostgreSQL creation fails
		s.dbRepo.Delete(ctx, database.ID)
		opCtx.Fail(err)
		return nil, err
	}
	opCtx.CompleteStep("Creating PostgreSQL database")

	// Step 5: Test connection
	opCtx.UpdateProgress("Testing connection", "Verifying database connection works", 0.8)
	if err := s.pgService.TestConnection(ctx, database); err != nil {
		opCtx.Fail(err)
		return nil, err
	}
	opCtx.CompleteStep("Testing connection")

	// Step 6: Finalize setup
	opCtx.UpdateProgress("Finalizing setup", "Completing database setup", 0.9)
	response := database.ToResponse()
	opCtx.Complete(response)

	return &response, nil
}

// UpdateProjectDatabase updates an existing project database configuration
func (s *DatabaseService) UpdateProjectDatabase(ctx context.Context, databaseID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	// Validate request
	if err := ValidateDatabaseCreateRequest(req); err != nil {
		return nil, errors.NewValidationError("validation failed").WithDetails(err.Error())
	}

	// Get existing database
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return nil, errors.NewNotFoundError("Database")
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
		return nil, errors.NewDatabaseError(err).WithDetails("failed to update database configuration")
	}

	response := database.ToResponse()
	return &response, nil
}

// DeleteProjectDatabase deletes a project database
func (s *DatabaseService) DeleteProjectDatabase(ctx context.Context, databaseID uuid.UUID) error {
	// Get database configuration
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return errors.NewNotFoundError("Database")
	}

	// Drop the PostgreSQL database using the management service
	if err := s.pgService.DropDatabase(ctx, database); err != nil {
		return err
	}

	// Remove from core database
	if err := s.dbRepo.Delete(ctx, databaseID); err != nil {
		return errors.NewDatabaseError(err).WithDetails("failed to remove database configuration")
	}

	return nil
}

// TestDatabaseConnection tests connection to a project database
func (s *DatabaseService) TestDatabaseConnection(ctx context.Context, databaseID uuid.UUID) error {
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return errors.NewNotFoundError("Database")
	}

	return s.pgService.TestConnection(ctx, database)
}

// GetProjectDatabases retrieves all databases for a project
func (s *DatabaseService) GetProjectDatabases(ctx context.Context, projectID uuid.UUID) ([]*entities.DatabaseResponse, error) {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to get databases")
	}

	var responses []*entities.DatabaseResponse
	for _, database := range databases {
		response := database.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// CreateTableWithProgress creates a table with progress tracking
func (s *DatabaseService) CreateTableWithProgress(ctx context.Context, projectID uuid.UUID, tableReq *entities.TableCreateRequest, userID uuid.UUID) (*entities.TableResponse, error) {
	// If progress tracker is not available, return error
	if s.progressTracker == nil {
		return nil, errors.NewAPIError(errors.ErrCodeOperationFailed, "progress tracking not available")
	}

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
		opCtx.Fail(errors.NewValidationError("validation failed").WithDetails(err.Error()))
		return nil, err
	}
	opCtx.CompleteStep("Validating table request")

	// Step 2: Get project database
	opCtx.UpdateProgress("Getting project database", "Retrieving project database configuration", 0.2)
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		opCtx.Fail(errors.NewDatabaseError(err).WithDetails("failed to get project databases"))
		return nil, err
	}

	if len(databases) == 0 {
		err := errors.NewNotFoundError("Project database")
		opCtx.Fail(err)
		return nil, err
	}

	database := databases[0] // Use first database
	opCtx.CompleteStep("Getting project database")

	// Step 3: Create table schema
	opCtx.UpdateProgress("Creating table schema", "Preparing table creation SQL", 0.3)
	createTableSQL, err := s.generateCreateTableSQL(tableReq)
	if err != nil {
		opCtx.Fail(errors.NewAPIError(errors.ErrCodeOperationFailed, "failed to generate table SQL").WithDetails(err.Error()))
		return nil, err
	}
	opCtx.CompleteStep("Creating table schema")

	// Step 4: Execute CREATE TABLE
	opCtx.UpdateProgress("Executing CREATE TABLE", "Creating table in database", 0.5)
	if err := s.pgService.CreateTable(ctx, database, tableReq.Name, createTableSQL); err != nil {
		opCtx.Fail(err)
		return nil, err
	}
	opCtx.CompleteStep("Executing CREATE TABLE")

	// Step 5: Verify table creation
	opCtx.UpdateProgress("Verifying table creation", "Checking if table was created successfully", 0.8)
	exists, err := s.pgService.TableExists(ctx, database, tableReq.Name)
	if err != nil {
		opCtx.Fail(err)
		return nil, err
	}
	if !exists {
		err := errors.NewAPIError(errors.ErrCodeOperationFailed, fmt.Sprintf("table %s was not created", tableReq.Name))
		opCtx.Fail(err)
		return nil, err
	}
	opCtx.CompleteStep("Verifying table creation")

	// Step 6: Update metadata
	opCtx.UpdateProgress("Updating metadata", "Saving table metadata", 0.9)

	// Convert schema to JSON string
	schemaJSON, err := json.Marshal(tableReq.Schema)
	if err != nil {
		opCtx.Fail(errors.NewAPIError(errors.ErrCodeOperationFailed, "failed to marshal schema").WithDetails(err.Error()))
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
func (s *DatabaseService) generateCreateTableSQL(tableReq *entities.TableCreateRequest) (string, error) {
	// This is a simplified implementation
	// In a real implementation, you would parse the schema and generate proper SQL
	return fmt.Sprintf("CREATE TABLE %s (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), created_at TIMESTAMP DEFAULT NOW())", tableReq.Name), nil
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
