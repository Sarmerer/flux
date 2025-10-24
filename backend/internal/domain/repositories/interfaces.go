package repositories

import (
	"context"

	"github.com/flow/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
}

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.Project, error)
	Update(ctx context.Context, project *entities.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Project, error)
}

// DatabaseRepository defines the interface for database data operations
type DatabaseRepository interface {
	Create(ctx context.Context, database *entities.Database) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Database, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Database, error)
	Update(ctx context.Context, database *entities.Database) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Database, error)
}

// TableRepository defines the interface for table data operations
type TableRepository interface {
	Create(ctx context.Context, table *entities.Table) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Table, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Table, error)
	Update(ctx context.Context, table *entities.Table) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Table, error)
}

// WorkflowRepository defines the interface for workflow data operations
type WorkflowRepository interface {
	Create(ctx context.Context, workflow *entities.Workflow) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Workflow, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Workflow, error)
	Update(ctx context.Context, workflow *entities.Workflow) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Workflow, error)
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error
}

type ProjectMemberRepository interface {
	Create(ctx context.Context, member *entities.ProjectMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectMember, error)
	GetByProjectAndUser(ctx context.Context, projectID, userID uuid.UUID) (*entities.ProjectMember, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.ProjectMember, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ProjectMember, error)
	Update(ctx context.Context, member *entities.ProjectMember) error
	Delete(ctx context.Context, id uuid.UUID) error
}
