package services

import (
	"context"
	"fmt"

	"github.com/flow/internal/errors"
	"github.com/google/uuid"
)

type TableDataServiceInterface interface {
	GetTableData(ctx context.Context, projectID, tableID uuid.UUID, page, limit int) (*TableDataResponse, error)
	GetRowByID(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}) (map[string]interface{}, error)
	InsertRow(ctx context.Context, projectID, tableID uuid.UUID, data map[string]interface{}) (map[string]interface{}, error)
	UpdateRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}, data map[string]interface{}) error
	DeleteRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}) error
}

type TableDataResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
}

type TableDataService struct {
	repoFactory  *ProjectRepositoryFactory
	tableService TableServiceInterface
}

func NewTableDataService(
	repoFactory *ProjectRepositoryFactory,
	tableService TableServiceInterface,
) *TableDataService {
	return &TableDataService{
		repoFactory:  repoFactory,
		tableService: tableService,
	}
}

func (s *TableDataService) GetTableData(ctx context.Context, projectID, tableID uuid.UUID, page, limit int) (*TableDataResponse, error) {
	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return nil, err
	}

	dataRepo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	data, err := dataRepo.Query(ctx, table.Name, limit, offset)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to query table data: %w", err))
	}

	total, err := dataRepo.Count(ctx, table.Name)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to count table rows: %w", err))
	}

	if data == nil {
		data = []map[string]interface{}{}
	}

	return &TableDataResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *TableDataService) GetRowByID(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}) (map[string]interface{}, error) {
	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return nil, err
	}

	dataRepo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}

	row, err := dataRepo.GetByID(ctx, table.Name, rowID)
	if err != nil {
		return nil, errors.NewNotFoundError("Row not found")
	}

	return row, nil
}

func (s *TableDataService) InsertRow(ctx context.Context, projectID, tableID uuid.UUID, data map[string]interface{}) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("No data provided")
	}

	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return nil, err
	}

	dataRepo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}

	row, err := dataRepo.Insert(ctx, table.Name, data)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to insert row: %w", err))
	}

	return row, nil
}

func (s *TableDataService) UpdateRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}, data map[string]interface{}) error {
	if len(data) == 0 {
		return errors.NewValidationError("No data provided")
	}

	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return err
	}

	dataRepo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}

	if err := dataRepo.Update(ctx, table.Name, rowID, data); err != nil {
		if err.Error() == "row not found" {
			return errors.NewNotFoundError("Row not found")
		}
		return errors.NewInternalError(fmt.Errorf("failed to update row: %w", err))
	}

	return nil
}

func (s *TableDataService) DeleteRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}) error {
	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return err
	}

	dataRepo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}

	if err := dataRepo.Delete(ctx, table.Name, rowID); err != nil {
		if err.Error() == "row not found" {
			return errors.NewNotFoundError("Row not found")
		}
		return errors.NewInternalError(fmt.Errorf("failed to delete row: %w", err))
	}

	return nil
}
