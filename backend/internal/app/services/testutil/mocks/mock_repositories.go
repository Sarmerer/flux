package mocks

import (
	"context"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *entities.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) CreateTx(ctx context.Context, tx types.Executor, project *entities.Project) error {
	args := m.Called(ctx, tx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByMemberID(ctx context.Context, memberID uuid.UUID) ([]*entities.Project, error) {
	args := m.Called(ctx, memberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *entities.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateTx(ctx context.Context, tx types.Executor, project *entities.Project) error {
	args := m.Called(ctx, tx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, limit, offset int) ([]*entities.Project, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Project), args.Error(1)
}

type MockDatabaseRepository struct {
	mock.Mock
}

func (m *MockDatabaseRepository) Create(ctx context.Context, database *entities.Database) error {
	args := m.Called(ctx, database)
	return args.Error(0)
}

func (m *MockDatabaseRepository) CreateTx(ctx context.Context, tx types.Executor, database *entities.Database) error {
	args := m.Called(ctx, tx, database)
	return args.Error(0)
}

func (m *MockDatabaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Database, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Database), args.Error(1)
}

func (m *MockDatabaseRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Database, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Database), args.Error(1)
}

func (m *MockDatabaseRepository) Update(ctx context.Context, database *entities.Database) error {
	args := m.Called(ctx, database)
	return args.Error(0)
}

func (m *MockDatabaseRepository) UpdateTx(ctx context.Context, tx types.Executor, database *entities.Database) error {
	args := m.Called(ctx, tx, database)
	return args.Error(0)
}

func (m *MockDatabaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDatabaseRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockDatabaseRepository) List(ctx context.Context, limit, offset int) ([]*entities.Database, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Database), args.Error(1)
}

type MockTableRepository struct {
	mock.Mock
}

func (m *MockTableRepository) Create(ctx context.Context, table *entities.Table) error {
	args := m.Called(ctx, table)
	return args.Error(0)
}

func (m *MockTableRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Table, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Table), args.Error(1)
}

func (m *MockTableRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Table, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Table), args.Error(1)
}

func (m *MockTableRepository) Update(ctx context.Context, table *entities.Table) error {
	args := m.Called(ctx, table)
	return args.Error(0)
}

func (m *MockTableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTableRepository) List(ctx context.Context, limit, offset int) ([]*entities.Table, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Table), args.Error(1)
}

type MockWorkflowRepository struct {
	mock.Mock
}

func (m *MockWorkflowRepository) Create(ctx context.Context, workflow *entities.Workflow) error {
	args := m.Called(ctx, workflow)
	return args.Error(0)
}

func (m *MockWorkflowRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Workflow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Workflow, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) Update(ctx context.Context, workflow *entities.Workflow) error {
	args := m.Called(ctx, workflow)
	return args.Error(0)
}

func (m *MockWorkflowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWorkflowRepository) List(ctx context.Context, limit, offset int) ([]*entities.Workflow, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	args := m.Called(ctx, id, isActive)
	return args.Error(0)
}

type MockProjectMemberRepository struct {
	mock.Mock
}

func (m *MockProjectMemberRepository) Create(ctx context.Context, member *entities.ProjectMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) CreateTx(ctx context.Context, tx types.Executor, member *entities.ProjectMember) error {
	args := m.Called(ctx, tx, member)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectMember, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) GetByProjectAndUser(ctx context.Context, projectID, userID uuid.UUID) (*entities.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.ProjectMember, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ProjectMember, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) Update(ctx context.Context, member *entities.ProjectMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) UpdateTx(ctx context.Context, tx types.Executor, member *entities.ProjectMember) error {
	args := m.Called(ctx, tx, member)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

type MockTableDataRepository struct {
	mock.Mock
}

func (m *MockTableDataRepository) Query(ctx context.Context, tableName string, limit, offset int) ([]map[string]interface{}, error) {
	args := m.Called(ctx, tableName, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockTableDataRepository) Count(ctx context.Context, tableName string) (int64, error) {
	args := m.Called(ctx, tableName)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTableDataRepository) GetByID(ctx context.Context, tableName string, id interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, tableName, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTableDataRepository) Insert(ctx context.Context, tableName string, data map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, tableName, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTableDataRepository) Update(ctx context.Context, tableName string, id interface{}, data map[string]interface{}) error {
	args := m.Called(ctx, tableName, id, data)
	return args.Error(0)
}

func (m *MockTableDataRepository) Delete(ctx context.Context, tableName string, id interface{}) error {
	args := m.Called(ctx, tableName, id)
	return args.Error(0)
}
