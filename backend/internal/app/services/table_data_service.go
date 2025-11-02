package services

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/google/uuid"
)

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

func (s *TableDataService) getDataRepo(ctx context.Context, projectID uuid.UUID) (repositories.TableDataRepository, error) {
	repo, err := s.repoFactory.GetTableDataRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to get data repository: %w", err))
	}
	return repo, nil
}

type tableDataContext struct {
	repo  repositories.TableDataRepository
	table *entities.TableResponse
}

func (s *TableDataService) prepareTableDataOp(ctx context.Context, projectID, tableID uuid.UUID) (*tableDataContext, error) {
	table, err := s.tableService.GetTableByID(ctx, tableID, projectID)
	if err != nil {
		return nil, err
	}

	repo, err := s.getDataRepo(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &tableDataContext{repo, table}, nil
}

func (s *TableDataService) GetTableData(ctx context.Context, projectID, tableID uuid.UUID, page, limit int) (*TableDataResponse, error) {
	dataCtx, err := s.prepareTableDataOp(ctx, projectID, tableID)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	data, err := dataCtx.repo.Query(ctx, dataCtx.table.Name, limit, offset)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to query table data: %w", err))
	}

	total, err := dataCtx.repo.Count(ctx, dataCtx.table.Name)
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
	dataCtx, err := s.prepareTableDataOp(ctx, projectID, tableID)
	if err != nil {
		return nil, err
	}

	row, err := dataCtx.repo.GetByID(ctx, dataCtx.table.Name, rowID)
	if err != nil {
		return nil, errors.NewNotFoundError("Row not found")
	}

	return row, nil
}

func (s *TableDataService) InsertRow(ctx context.Context, projectID, tableID uuid.UUID, data map[string]interface{}) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.NewValidationError("No data provided")
	}

	dataCtx, err := s.prepareTableDataOp(ctx, projectID, tableID)
	if err != nil {
		return nil, err
	}

	row, err := dataCtx.repo.Insert(ctx, dataCtx.table.Name, data)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to insert row: %w", err))
	}

	return row, nil
}

func (s *TableDataService) UpdateRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}, data map[string]interface{}) error {
	if len(data) == 0 {
		return errors.NewValidationError("No data provided")
	}

	dataCtx, err := s.prepareTableDataOp(ctx, projectID, tableID)
	if err != nil {
		return err
	}

	if err := dataCtx.repo.Update(ctx, dataCtx.table.Name, rowID, data); err != nil {
		if err.Error() == "row not found" {
			return errors.NewNotFoundError("Row not found")
		}
		return errors.NewInternalError(fmt.Errorf("failed to update row: %w", err))
	}

	return nil
}

func (s *TableDataService) DeleteRow(ctx context.Context, projectID, tableID uuid.UUID, rowID interface{}) error {
	dataCtx, err := s.prepareTableDataOp(ctx, projectID, tableID)
	if err != nil {
		return err
	}

	if err := dataCtx.repo.Delete(ctx, dataCtx.table.Name, rowID); err != nil {
		if err.Error() == "row not found" {
			return errors.NewNotFoundError("Row not found")
		}
		return errors.NewInternalError(fmt.Errorf("failed to delete row: %w", err))
	}

	return nil
}
