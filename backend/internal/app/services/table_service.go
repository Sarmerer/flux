package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"

	"github.com/google/uuid"
)

type TableService struct {
	repoFactory *ProjectRepositoryFactory
}

func NewTableService(repoFactory *ProjectRepositoryFactory) *TableService {
	return &TableService{
		repoFactory: repoFactory,
	}
}

func (s *TableService) CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table repository: %w", err)
	}

	schemaJSON, err := json.Marshal(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}

	table := &entities.Table{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Schema:      string(schemaJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tableRepo.Create(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) GetTableByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.TableResponse, error) {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table repository: %w", err)
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
		return nil, fmt.Errorf("failed to get table repository: %w", err)
	}

	tables, err := tableRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
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
		return nil, fmt.Errorf("failed to get table repository: %w", err)
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
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

func (s *TableService) DeleteTable(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	tableRepo, err := s.repoFactory.GetTableRepository(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get table repository: %w", err)
	}

	_, err = tableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	if err := tableRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	return nil
}
