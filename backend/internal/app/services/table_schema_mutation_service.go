package services

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"

	"github.com/google/uuid"
)

type TableSchemaMutationService struct {
	tableRepo   repositories.TableRepository
	dbRepo      repositories.DatabaseRepository
	projectRepo repositories.ProjectRepository
	schemaSvc   *database.SchemaManagementService
}

func NewTableSchemaMutationService(
	tableRepo repositories.TableRepository,
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	schemaSvc *database.SchemaManagementService,
) *TableSchemaMutationService {
	return &TableSchemaMutationService{
		tableRepo:   tableRepo,
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
		schemaSvc:   schemaSvc,
	}
}

type ColumnDefinition = database.ColumnDefinition
type TableSchema = database.TableSchema
type IndexDefinition = database.IndexDefinition
type ForeignKeyDefinition = database.ForeignKeyDefinition

type AddColumnRequest struct {
	Column ColumnDefinition `json:"column" validate:"required"`
}

type RemoveColumnRequest struct {
	ColumnName string `json:"column_name" validate:"required"`
}

type ModifyColumnRequest struct {
	ColumnName string           `json:"column_name" validate:"required"`
	Column     ColumnDefinition `json:"column" validate:"required"`
}

type AddForeignKeyRequest struct {
	ForeignKey ForeignKeyDefinition `json:"foreign_key" validate:"required"`
}

type RemoveForeignKeyRequest struct {
	ForeignKeyName string `json:"foreign_key_name" validate:"required"`
}

func (s *TableSchemaMutationService) CreateTableInDatabase(ctx context.Context, projectID uuid.UUID, tableName string, schema TableSchema) error {
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}
	if !isValidTableName(tableName) {
		return fmt.Errorf("invalid table name")
	}

	if err := ValidateTableSchema(&schema); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.CreateTable(ctx, database, tableName, schema)
}

func (s *TableSchemaMutationService) DropTableFromDatabase(ctx context.Context, projectID uuid.UUID, tableName string) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.DropTable(ctx, database, tableName)
}

func (s *TableSchemaMutationService) AddColumnToTable(ctx context.Context, projectID uuid.UUID, tableName string, req *AddColumnRequest) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.AddColumn(ctx, database, tableName, req.Column)
}

func (s *TableSchemaMutationService) RemoveColumnFromTable(ctx context.Context, projectID uuid.UUID, tableName string, req *RemoveColumnRequest) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.RemoveColumn(ctx, database, tableName, req.ColumnName)
}

func (s *TableSchemaMutationService) ModifyColumnInTable(ctx context.Context, projectID uuid.UUID, tableName string, req *ModifyColumnRequest) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.ModifyColumn(ctx, database, tableName, req.ColumnName, req.Column)
}

func (s *TableSchemaMutationService) AddForeignKeyToTable(ctx context.Context, projectID uuid.UUID, tableName string, req *AddForeignKeyRequest) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.AddForeignKey(ctx, database, tableName, req.ForeignKey)
}

func (s *TableSchemaMutationService) RemoveForeignKeyFromTable(ctx context.Context, projectID uuid.UUID, tableName string, req *RemoveForeignKeyRequest) error {
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	return s.schemaSvc.RemoveForeignKey(ctx, database, tableName, req.ForeignKeyName)
}
