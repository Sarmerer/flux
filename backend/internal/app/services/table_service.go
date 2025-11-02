package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/validation"

	"github.com/google/uuid"
)

type TableService struct {
	repoFactory   *ProjectRepositoryFactory
	resolver      *database.ProjectConnectionResolver
	schemaService *database.SchemaManagementService
}

func NewTableService(
	repoFactory *ProjectRepositoryFactory,
	resolver *database.ProjectConnectionResolver,
	schemaService *database.SchemaManagementService,
) *TableService {
	return &TableService{
		repoFactory:   repoFactory,
		resolver:      resolver,
		schemaService: schemaService,
	}
}

func (s *TableService) CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	if !validation.IsValidTableName(req.Name) {
		return nil, errors.NewValidationError("Invalid table name").
			WithField("name").
			WithDetails(fmt.Sprintf("Table name '%s' contains invalid characters or format", req.Name))
	}

	var tableSchema validation.TableSchema
	schemaBytes, err := json.Marshal(req.Schema)
	if err != nil {
		return nil, errors.NewValidationError("Invalid schema format").
			WithField("schema").
			WithDetails(err.Error())
	}

	if err := json.Unmarshal(schemaBytes, &tableSchema); err != nil {
		return nil, errors.NewValidationError("Schema does not match expected structure").
			WithField("schema").
			WithDetails(err.Error())
	}

	if err := validation.ValidateTableSchema(&tableSchema); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return nil, apiErr
		}
		return nil, errors.NewValidationError("Schema validation failed").
			WithField("schema").
			WithDetails(err.Error())
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return nil, errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	table := &entities.Table{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Schema:      string(schemaBytes),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.schemaService.CreateTable(ctx, db, table, tableSchema); err != nil {
		return nil, errors.NewDatabaseError("Failed to create table", err)
	}

	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		if dropErr := s.schemaService.DropTable(ctx, db, table); dropErr != nil {
			return nil, errors.NewInternalError(err).
				WithDetails("Failed to get table repository and failed to rollback physical table: " + dropErr.Error())
		}
		return nil, errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	if err := tableRepo.Create(ctx, table); err != nil {
		if dropErr := s.schemaService.DropTable(ctx, db, table); dropErr != nil {
			return nil, errors.NewDatabaseError("Failed to create table metadata", err).
				WithDetails("Additionally failed to rollback physical table: " + dropErr.Error())
		}
		return nil, errors.NewDatabaseError("Failed to create table metadata", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) GetTableByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Table not found")
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	tables, err := tableRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get tables", err)
	}

	responses := make([]*entities.TableResponse, 0)
	for _, table := range tables {
		response := table.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *TableService) UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Table not found")
	}

	if req.Name != "" {
		table.Name = req.Name
	}
	if req.Description != "" {
		table.Description = req.Description
	}
	if req.Schema != nil {
		schemaJSON, err := json.Marshal(req.Schema)
		if err != nil {
			return nil, errors.NewValidationError("Invalid schema format").
				WithField("schema").
				WithDetails(err.Error())
		}
		table.Schema = string(schemaJSON)
	}

	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return nil, errors.NewDatabaseError("Failed to update table", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) DeleteTable(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.DropTable(ctx, db, table); err != nil {
		return errors.NewDatabaseError("Failed to drop physical table", err)
	}

	if err := tableRepo.Delete(ctx, id); err != nil {
		return errors.NewDatabaseError("Failed to delete table metadata", err)
	}

	return nil
}

func (s *TableService) AddColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, column database.ColumnDefinition) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.AddColumn(ctx, db, table, column); err != nil {
		return errors.NewDatabaseError("Failed to add column", err)
	}

	var tableSchema database.TableSchema
	if err := json.Unmarshal([]byte(table.Schema), &tableSchema); err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to parse table schema")
	}

	tableSchema.Columns = append(tableSchema.Columns, column)
	schemaBytes, err := json.Marshal(tableSchema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}

func (s *TableService) RemoveColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, columnName string) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.RemoveColumn(ctx, db, table, columnName); err != nil {
		return errors.NewDatabaseError("Failed to remove column", err)
	}

	var tableSchema database.TableSchema
	if err := json.Unmarshal([]byte(table.Schema), &tableSchema); err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to parse table schema")
	}

	for i, col := range tableSchema.Columns {
		if col.Name == columnName {
			tableSchema.Columns = append(tableSchema.Columns[:i], tableSchema.Columns[i+1:]...)
			break
		}
	}

	schemaBytes, err := json.Marshal(tableSchema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}

func (s *TableService) ModifyColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, oldColumnName string, newColumn database.ColumnDefinition) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.ModifyColumn(ctx, db, table, oldColumnName, newColumn); err != nil {
		return errors.NewDatabaseError("Failed to modify column", err)
	}

	var tableSchema database.TableSchema
	if err := json.Unmarshal([]byte(table.Schema), &tableSchema); err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to parse table schema")
	}

	for i, col := range tableSchema.Columns {
		if col.Name == oldColumnName {
			tableSchema.Columns[i] = newColumn
			break
		}
	}

	schemaBytes, err := json.Marshal(tableSchema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}

func (s *TableService) AddForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, fk database.ForeignKeyDefinition) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.AddForeignKey(ctx, db, table, fk); err != nil {
		return errors.NewDatabaseError("Failed to add foreign key", err)
	}

	var tableSchema database.TableSchema
	if err := json.Unmarshal([]byte(table.Schema), &tableSchema); err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to parse table schema")
	}

	tableSchema.ForeignKeys = append(tableSchema.ForeignKeys, fk)
	schemaBytes, err := json.Marshal(tableSchema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}

func (s *TableService) RemoveForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, foreignKeyName string) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.RemoveForeignKey(ctx, db, table, foreignKeyName); err != nil {
		return errors.NewDatabaseError("Failed to remove foreign key", err)
	}

	var tableSchema database.TableSchema
	if err := json.Unmarshal([]byte(table.Schema), &tableSchema); err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to parse table schema")
	}

	for i, fkDef := range tableSchema.ForeignKeys {
		if fkDef.Name == foreignKeyName {
			tableSchema.ForeignKeys = append(tableSchema.ForeignKeys[:i], tableSchema.ForeignKeys[i+1:]...)
			break
		}
	}

	schemaBytes, err := json.Marshal(tableSchema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}

func (s *TableService) UpdateTableSchema(ctx context.Context, id uuid.UUID, projectID uuid.UUID, schema database.TableSchema) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to get table repository")
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Table not found")
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return errors.NewNotFoundError("Project database not found").WithDetails(err.Error())
	}

	if err := s.schemaService.UpdateTableSchema(ctx, db, table, schema); err != nil {
		return errors.NewDatabaseError("Failed to update table schema", err)
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal table schema")
	}

	table.Schema = string(schemaBytes)
	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return errors.NewDatabaseError("Failed to update table metadata", err)
	}

	return nil
}
