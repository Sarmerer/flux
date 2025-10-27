package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"
	"github.com/flow/internal/validation"

	"github.com/google/uuid"
)

type DatabaseService struct {
	dbRepo          repositories.DatabaseRepository
	projectRepo     repositories.ProjectRepository
	pgService       *database.PostgreSQLManagementService
	connService     *database.ConnectionService
	progressTracker *progress.Tracker
	hub             *realtime.Hub
}

func NewDatabaseService(
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	pgService *database.PostgreSQLManagementService,
	connService *database.ConnectionService,
	progressTracker *progress.Tracker,
	hub *realtime.Hub,
) *DatabaseService {
	return &DatabaseService{
		dbRepo:          dbRepo,
		projectRepo:     projectRepo,
		pgService:       pgService,
		connService:     connService,
		progressTracker: progressTracker,
		hub:             hub,
	}
}

func (s *DatabaseService) CreateProjectDatabase(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	if s.progressTracker == nil {
		return s.createDatabase(ctx, projectID, req, nil)
	}
	return s.createDatabase(ctx, projectID, req, nil)
}

func (s *DatabaseService) createDatabase(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest, opCtx *progress.OperationContext) (*entities.DatabaseResponse, error) {
	if opCtx != nil {
		opCtx.UpdateProgress("Validating request", "Validating database creation request", 0.1)
	}
	if err := validation.ValidateDatabaseCreateRequest(req); err != nil {
		if opCtx != nil {
			opCtx.Fail(errors.NewValidationError("validation failed").WithDetails(err.Error()))
		}
		return nil, errors.NewValidationError("validation failed").WithDetails(err.Error())
	}
	if opCtx != nil {
		opCtx.CompleteStep("Validating request")
		opCtx.UpdateProgress("Verifying project exists", "Checking if project exists", 0.2)
	}

	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if opCtx != nil {
			opCtx.Fail(errors.NewNotFoundError("Project not found"))
		}
		return nil, errors.NewNotFoundError("Project not found")
	}
	if opCtx != nil {
		opCtx.CompleteStep("Verifying project exists")
		opCtx.UpdateProgress("Creating database configuration", "Saving database configuration", 0.3)
	}

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
		dbErr := errors.NewDatabaseError("Failed to save database configuration", err)
		if opCtx != nil {
			opCtx.Fail(dbErr)
		}
		return nil, dbErr
	}
	if opCtx != nil {
		opCtx.CompleteStep("Creating database configuration")
		opCtx.UpdateProgress("Creating PostgreSQL database", "Creating the actual PostgreSQL database", 0.5)
	}

	if err := s.pgService.CreateDatabase(ctx, database); err != nil {
		s.dbRepo.Delete(ctx, database.ID)
		if opCtx != nil {
			opCtx.Fail(err)
		}
		return nil, err
	}
	if opCtx != nil {
		opCtx.CompleteStep("Creating PostgreSQL database")
	}

	response := database.ToResponse()
	if opCtx != nil {
		opCtx.Complete(response)
	}
	return &response, nil
}

func (s *DatabaseService) CreateProjectDatabaseWithProgress(ctx context.Context, projectID uuid.UUID, req *entities.DatabaseCreateRequest, userID uuid.UUID) (*entities.DatabaseResponse, error) {
	if s.progressTracker == nil {
		return s.createDatabase(ctx, projectID, req, nil)
	}

	steps := []string{
		"Validating request",
		"Verifying project exists",
		"Creating database configuration",
		"Creating PostgreSQL database",
		"Testing connection",
		"Finalizing setup",
	}

	opCtx := s.progressTracker.NewOperationContext(userID, "create_database", "Creating project database", steps)
	return s.createDatabase(ctx, projectID, req, opCtx)
}

func (s *DatabaseService) UpdateProjectDatabase(ctx context.Context, databaseID uuid.UUID, req *entities.DatabaseCreateRequest) (*entities.DatabaseResponse, error) {
	if err := validation.ValidateDatabaseCreateRequest(req); err != nil {
		return nil, errors.NewValidationError("validation failed").WithDetails(err.Error())
	}

	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return nil, errors.NewNotFoundError("Database not found")
	}

	database.Name = req.Name
	database.Host = req.Host
	database.Port = req.Port
	database.Username = req.Username
	database.Password = req.Password
	database.Database = req.Database
	database.SSLMode = req.SSLMode
	database.UpdatedAt = time.Now()

	if err := s.dbRepo.Update(ctx, database); err != nil {
		return nil, errors.NewDatabaseError("Failed to update database configuration", err)
	}

	response := database.ToResponse()
	return &response, nil
}

func (s *DatabaseService) DeleteProjectDatabase(ctx context.Context, databaseID uuid.UUID) error {

	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return errors.NewNotFoundError("Database not found")
	}

	if err := s.pgService.DropDatabase(ctx, database); err != nil {
		return err
	}

	if err := s.dbRepo.Delete(ctx, databaseID); err != nil {
		return errors.NewDatabaseError("Failed to delete database configuration", err)
	}

	return nil
}

func (s *DatabaseService) TestDatabaseConnection(ctx context.Context, databaseID uuid.UUID) error {
	database, err := s.dbRepo.GetByID(ctx, databaseID)
	if err != nil {
		return errors.NewNotFoundError("Database not found")
	}

	return s.pgService.TestConnection(ctx, database)
}

func (s *DatabaseService) GetProjectDatabases(ctx context.Context, projectID uuid.UUID) ([]*entities.DatabaseResponse, error) {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to retrieve databases", err)
	}

	responses := make([]*entities.DatabaseResponse, 0)
	for _, database := range databases {
		response := database.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *DatabaseService) CreateTableWithProgress(ctx context.Context, projectID uuid.UUID, tableReq *entities.TableCreateRequest, userID uuid.UUID) (*entities.TableResponse, error) {

	if s.progressTracker == nil {
		return nil, errors.NewAPIError(errors.ErrCodeOperationFailed, "progress tracking not available")
	}

	steps := []string{
		"Validating table request",
		"Getting project database",
		"Creating table schema",
		"Executing CREATE TABLE",
		"Verifying table creation",
		"Updating metadata",
	}

	opCtx := s.progressTracker.NewOperationContext(userID, "create_table", "Creating table", steps)

	opCtx.UpdateProgress("Validating table request", "Validating table creation request", 0.1)
	if err := validation.ValidateTableCreateRequest(tableReq); err != nil {
		opCtx.Fail(errors.NewValidationError("validation failed").WithDetails(err.Error()))
		return nil, err
	}
	opCtx.CompleteStep("Validating table request")

	opCtx.UpdateProgress("Getting project database", "Retrieving project database configuration", 0.2)
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		opCtx.Fail(errors.NewDatabaseError("Failed to retrieve project database", err))
		return nil, err
	}

	if len(databases) == 0 {
		err := errors.NewNotFoundError("Project database not found")
		opCtx.Fail(err)
		return nil, err
	}

	database := databases[0]
	opCtx.CompleteStep("Getting project database")

	opCtx.UpdateProgress("Creating table schema", "Preparing table creation SQL", 0.3)
	createTableSQL, err := s.generateCreateTableSQL(tableReq)
	if err != nil {
		opCtx.Fail(errors.NewAPIError(errors.ErrCodeOperationFailed, "Failed to generate table SQL").WithDetails(err.Error()))
		return nil, err
	}
	opCtx.CompleteStep("Creating table schema")

	opCtx.UpdateProgress("Executing CREATE TABLE", "Creating table in database", 0.5)
	if err := s.pgService.CreateTable(ctx, database, tableReq.Name, createTableSQL); err != nil {
		opCtx.Fail(err)
		return nil, err
	}
	opCtx.CompleteStep("Executing CREATE TABLE")

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

	opCtx.UpdateProgress("Updating metadata", "Saving table metadata", 0.9)

	schemaJSON, err := json.Marshal(tableReq.Schema)
	if err != nil {
		opCtx.Fail(errors.NewAPIError(errors.ErrCodeOperationFailed, "Failed to marshal schema").WithDetails(err.Error()))
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

func (s *DatabaseService) generateCreateTableSQL(tableReq *entities.TableCreateRequest) (string, error) {
	var columns []string

	if columnsData, ok := tableReq.Schema["columns"].([]interface{}); ok {
		for _, colData := range columnsData {
			colMap, ok := colData.(map[string]interface{})
			if !ok {
				continue
			}

			name, _ := colMap["name"].(string)
			colType, _ := colMap["type"].(string)
			nullable, _ := colMap["nullable"].(bool)
			defaultValue, _ := colMap["default_value"].(string)
			unique, _ := colMap["unique"].(bool)

			if name == "" || colType == "" {
				continue
			}

			colDef := fmt.Sprintf("%s %s", name, colType)

			if !nullable {
				colDef += " NOT NULL"
			}

			if defaultValue != "" {
				colDef += fmt.Sprintf(" DEFAULT %s", defaultValue)
			}

			if unique {
				colDef += " UNIQUE"
			}

			columns = append(columns, colDef)
		}
	}

	if len(columns) == 0 {
		columns = append(columns, "id UUID PRIMARY KEY DEFAULT gen_random_uuid()")
		columns = append(columns, "created_at TIMESTAMP DEFAULT NOW()")
	}

	if primaryKey, ok := tableReq.Schema["primary_key"].([]interface{}); ok && len(primaryKey) > 0 {
		var pkColumns []string
		for _, pk := range primaryKey {
			if pkStr, ok := pk.(string); ok {
				pkColumns = append(pkColumns, pkStr)
			}
		}
		if len(pkColumns) > 0 {
			columns = append(columns, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", ")))
		}
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", tableReq.Name, strings.Join(columns, ", ")), nil
}
