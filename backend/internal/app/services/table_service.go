package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
)

// TableService handles table-related business logic
type TableService struct {
	tableRepo repositories.TableRepository
}

// NewTableService creates a new TableService
func NewTableService(tableRepo repositories.TableRepository) *TableService {
	return &TableService{
		tableRepo: tableRepo,
	}
}

// CreateTable creates a new table
func (s *TableService) CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	// Convert schema to JSON string
	schemaJSON, err := json.Marshal(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}

	// Create table
	table := &entities.Table{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Schema:    string(schemaJSON),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tableRepo.Create(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

// GetTableByID retrieves a table by ID
func (s *TableService) GetTableByID(ctx context.Context, id uuid.UUID) (*entities.TableResponse, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

// GetTablesByProjectID retrieves all tables for a project
func (s *TableService) GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error) {
	tables, err := s.tableRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}

	var responses []*entities.TableResponse
	for _, table := range tables {
		response := table.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// UpdateTable updates a table
func (s *TableService) UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest) (*entities.TableResponse, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("table not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		table.Name = req.Name
	}
	if req.Schema != nil {
		schemaJSON, err := json.Marshal(req.Schema)
		if err != nil {
			return nil, fmt.Errorf("invalid schema format: %w", err)
		}
		table.Schema = string(schemaJSON)
	}

	table.UpdatedAt = time.Now()

	if err := s.tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	response := table.ToResponse()
	return &response, nil
}

// DeleteTable deletes a table
func (s *TableService) DeleteTable(ctx context.Context, id uuid.UUID) error {
	// Check if table exists
	_, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	if err := s.tableRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	return nil
}
