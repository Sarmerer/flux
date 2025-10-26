package services

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"
	postgresRepos "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/google/uuid"
)

type ProjectRepositoryFactory struct {
	resolver *database.ProjectConnectionResolver
}

func NewProjectRepositoryFactory(resolver *database.ProjectConnectionResolver) *ProjectRepositoryFactory {
	return &ProjectRepositoryFactory{
		resolver: resolver,
	}
}

func (f *ProjectRepositoryFactory) GetTableRepository(ctx context.Context, projectID uuid.UUID) (repositories.TableRepository, error) {
	pool, err := f.resolver.GetProjectDB(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project database connection: %w", err)
	}
	return postgresRepos.NewTableRepository(pool), nil
}

func (f *ProjectRepositoryFactory) GetWorkflowRepository(ctx context.Context, projectID uuid.UUID) (repositories.WorkflowRepository, error) {
	pool, err := f.resolver.GetProjectDB(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project database connection: %w", err)
	}
	return postgresRepos.NewWorkflowRepository(pool), nil
}
