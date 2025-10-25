package services

import (
	"context"

	"github.com/flow/internal/domain/entities"

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
	GetProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.ProjectResponse, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type TableServiceInterface interface {
	CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error)
	GetTableByID(ctx context.Context, id uuid.UUID) (*entities.TableResponse, error)
	GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error)
	UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest) (*entities.TableResponse, error)
	DeleteTable(ctx context.Context, id uuid.UUID) error
}

type WorkflowServiceInterface interface {
	CreateWorkflow(ctx context.Context, req *entities.WorkflowCreateRequest, projectID uuid.UUID) (*entities.WorkflowResponse, error)
	GetWorkflowByID(ctx context.Context, id uuid.UUID) (*entities.WorkflowResponse, error)
	GetWorkflowsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.WorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, id uuid.UUID, req *entities.WorkflowUpdateRequest) (*entities.WorkflowResponse, error)
	DeleteWorkflow(ctx context.Context, id uuid.UUID) error
	ToggleWorkflowActive(ctx context.Context, id uuid.UUID, isActive bool) (*entities.WorkflowResponse, error)
	ExecuteWorkflow(ctx context.Context, id uuid.UUID) error
}
