package services

import (
	"context"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/infrastructure/database"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *entities.UserCreateRequest) (*entities.UserResponse, error)
	Login(ctx context.Context, req *entities.UserLoginRequest) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserResponse, error)
}

type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, req *entities.ProjectCreateRequest, ownerID uuid.UUID) (*entities.ProjectResponse, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*entities.ProjectResponse, error)
	GetProjectsByMemberID(ctx context.Context, memberID uuid.UUID) ([]*entities.ProjectResponse, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type TableServiceInterface interface {
	CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error)
	GetTableByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.TableResponse, error)
	GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error)
	UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest, projectID uuid.UUID) (*entities.TableResponse, error)
	DeleteTable(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error
	AddColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, column database.ColumnDefinition) error
	RemoveColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, columnName string) error
	ModifyColumn(ctx context.Context, id uuid.UUID, projectID uuid.UUID, oldColumnName string, newColumn database.ColumnDefinition) error
	AddForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, fk database.ForeignKeyDefinition) error
	RemoveForeignKey(ctx context.Context, id uuid.UUID, projectID uuid.UUID, foreignKeyName string) error
	UpdateTableSchema(ctx context.Context, id uuid.UUID, projectID uuid.UUID, schema database.TableSchema) error
}

type WorkflowServiceInterface interface {
	CreateWorkflow(ctx context.Context, req *entities.WorkflowCreateRequest, projectID uuid.UUID) (*entities.WorkflowResponse, error)
	GetWorkflowByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.WorkflowResponse, error)
	GetWorkflowsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.WorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID, req *entities.WorkflowUpdateRequest) (*entities.WorkflowResponse, error)
	DeleteWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error
	ToggleWorkflowActive(ctx context.Context, id uuid.UUID, projectID uuid.UUID, isActive bool) (*entities.WorkflowResponse, error)
	ExecuteWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error
}
