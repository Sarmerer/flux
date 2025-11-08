package mocks

import (
	"context"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockProjectConnectionResolver struct {
	mock.Mock
}

func (m *MockProjectConnectionResolver) GetProjectDB(ctx context.Context, projectID uuid.UUID) (*pgxpool.Pool, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgxpool.Pool), args.Error(1)
}

func (m *MockProjectConnectionResolver) GetDatabaseEntity(ctx context.Context, projectID uuid.UUID) (*entities.Database, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Database), args.Error(1)
}

func (m *MockProjectConnectionResolver) InvalidateCache(projectID uuid.UUID) {
	m.Called(projectID)
}

func (m *MockProjectConnectionResolver) Close() {
	m.Called()
}

type MockSchemaManagementService struct {
	mock.Mock
}

func (m *MockSchemaManagementService) CreateTable(ctx context.Context, db *entities.Database, table *entities.Table, schema interface{}) error {
	args := m.Called(ctx, db, table, schema)
	return args.Error(0)
}

func (m *MockSchemaManagementService) DropTable(ctx context.Context, db *entities.Database, table *entities.Table) error {
	args := m.Called(ctx, db, table)
	return args.Error(0)
}

func (m *MockSchemaManagementService) AddColumn(ctx context.Context, db *entities.Database, table *entities.Table, column database.ColumnDefinition) error {
	args := m.Called(ctx, db, table, column)
	return args.Error(0)
}

func (m *MockSchemaManagementService) RemoveColumn(ctx context.Context, db *entities.Database, table *entities.Table, columnName string) error {
	args := m.Called(ctx, db, table, columnName)
	return args.Error(0)
}

func (m *MockSchemaManagementService) ModifyColumn(ctx context.Context, db *entities.Database, table *entities.Table, oldColumnName string, newColumn database.ColumnDefinition) error {
	args := m.Called(ctx, db, table, oldColumnName, newColumn)
	return args.Error(0)
}

func (m *MockSchemaManagementService) AddForeignKey(ctx context.Context, db *entities.Database, table *entities.Table, fk database.ForeignKeyDefinition) error {
	args := m.Called(ctx, db, table, fk)
	return args.Error(0)
}

func (m *MockSchemaManagementService) RemoveForeignKey(ctx context.Context, db *entities.Database, table *entities.Table, fkName string) error {
	args := m.Called(ctx, db, table, fkName)
	return args.Error(0)
}

func (m *MockSchemaManagementService) UpdateTableSchema(ctx context.Context, db *entities.Database, table *entities.Table, schema database.TableSchema) error {
	args := m.Called(ctx, db, table, schema)
	return args.Error(0)
}

type MockSchemaInspectionService struct {
	mock.Mock
}

func (m *MockSchemaInspectionService) GetTableSchema(ctx context.Context, db *pgxpool.Pool, tableName string) (*entities.TableSchemaInfo, error) {
	args := m.Called(ctx, db, tableName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TableSchemaInfo), args.Error(1)
}

type MockPostgreSQLManagementService struct {
	mock.Mock
}

func (m *MockPostgreSQLManagementService) CreateDatabase(ctx context.Context, db *entities.Database) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

func (m *MockPostgreSQLManagementService) DropDatabase(ctx context.Context, db *entities.Database) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

func (m *MockPostgreSQLManagementService) TestConnection(ctx context.Context, db *entities.Database) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

type MockProjectMigrationRunner struct {
	mock.Mock
}

func (m *MockProjectMigrationRunner) InitializeProjectDatabase(ctx context.Context, db *entities.Database) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}
