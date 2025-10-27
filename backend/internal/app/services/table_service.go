package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
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
		return nil, fmt.Errorf("invalid table name: %s", req.Name)
	}

	var tableSchema validation.TableSchema
	schemaBytes, err := json.Marshal(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}

	if err := json.Unmarshal(schemaBytes, &tableSchema); err != nil {
		return nil, fmt.Errorf("schema does not match expected structure: %w", err)
	}

	if err := validation.ValidateTableSchema(&tableSchema); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get project database: %w", err)
	}

	if err := s.schemaService.CreateTable(ctx, db, req.Name, tableSchema); err != nil {
		return nil, fmt.Errorf("Failed to create table: %w", err)
	}

	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		if dropErr := s.schemaService.DropTable(ctx, db, req.Name); dropErr != nil {
			return nil, fmt.Errorf("Failed to get table repository (and failed to rollback physical table): %w", err)
		}
		return nil, fmt.Errorf("Failed to get table repository: %w", err)
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

	if err := tableRepo.Create(ctx, table); err != nil {
		if dropErr := s.schemaService.DropTable(ctx, db, req.Name); dropErr != nil {
			return nil, fmt.Errorf("Failed to create table metadata (and failed to rollback physical table): %w", err)
		}
		return nil, fmt.Errorf("Failed to create table metadata: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) GetTableByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get table repository: %w", err)
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get table repository: %w", err)
	}

	tables, err := tableRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tables: %w", err)
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
		return nil, fmt.Errorf("Failed to get table repository: %w", err)
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
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
			return nil, fmt.Errorf("invalid schema format: %w", err)
		}
		table.Schema = string(schemaJSON)
	}

	table.UpdatedAt = time.Now()

	if err := tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("Failed to update table: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) DeleteTable(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to get table repository: %w", err)
	}

	table, err := tableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	db, err := s.resolver.GetDatabaseEntity(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to get project database: %w", err)
	}

	if err := s.schemaService.DropTable(ctx, db, table.Name); err != nil {
		return fmt.Errorf("Failed to drop physical table: %w", err)
	}

	if err := tableRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Failed to delete table metadata: %w", err)
	}

	return nil
}
