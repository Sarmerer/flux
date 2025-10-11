package services

import (
	"context"

	"github.com/flow/internal/domain/entities"

	"github.com/google/uuid"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	Register(ctx context.Context, req *entities.UserCreateRequest) (*entities.UserResponse, error)
	Login(ctx context.Context, req *entities.UserLoginRequest) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserResponse, error)
}

// ProjectServiceInterface defines the interface for project service operations
type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, req *entities.ProjectCreateRequest, ownerID uuid.UUID) (*entities.ProjectResponse, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*entities.ProjectResponse, error)
	GetProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.ProjectResponse, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

// TableServiceInterface defines the interface for table service operations
type TableServiceInterface interface {
	CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error)
	GetTableByID(ctx context.Context, id uuid.UUID) (*entities.TableResponse, error)
	GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error)
	UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest) (*entities.TableResponse, error)
	DeleteTable(ctx context.Context, id uuid.UUID) error
}
