package mocks

import (
	"context"

	"github.com/flow/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTableService struct {
	mock.Mock
}

func (m *MockTableService) CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	args := m.Called(ctx, req, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) GetTableByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.TableResponse, error) {
	args := m.Called(ctx, id, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	args := m.Called(ctx, id, req, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) DeleteTable(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	args := m.Called(ctx, id, projectID)
	return args.Error(0)
}

func (m *MockTableService) AddColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, column interface{}) error {
	args := m.Called(ctx, id, projectID, column)
	return args.Error(0)
}

func (m *MockTableService) RemoveColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, columnName string) error {
	args := m.Called(ctx, id, projectID, columnName)
	return args.Error(0)
}

func (m *MockTableService) ModifyColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, oldColumnName string, newColumn interface{}) error {
	args := m.Called(ctx, id, projectID, oldColumnName, newColumn)
	return args.Error(0)
}

func (m *MockTableService) AddForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, fk interface{}) error {
	args := m.Called(ctx, id, projectID, fk)
	return args.Error(0)
}

func (m *MockTableService) RemoveForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, foreignKeyName string) error {
	args := m.Called(ctx, id, projectID, foreignKeyName)
	return args.Error(0)
}

func (m *MockTableService) UpdateTableSchema(ctx context.Context, id uuid.UUID, projectID uuid.UUID, schema interface{}) error {
	args := m.Called(ctx, id, projectID, schema)
	return args.Error(0)
}

