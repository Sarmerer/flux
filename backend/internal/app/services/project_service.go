package services

import (
	"context"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
)

// ProjectService handles project-related business logic
type ProjectService struct {
	projectRepo repositories.ProjectRepository
	dbRepo      repositories.DatabaseRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(projectRepo repositories.ProjectRepository, dbRepo repositories.DatabaseRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		dbRepo:      dbRepo,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, req *entities.ProjectCreateRequest, ownerID uuid.UUID) (*entities.ProjectResponse, error) {
	// Create project
	project := &entities.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		APIKey:      generateAPIKey(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Create default database for the project
	database := &entities.Database{
		ID:        uuid.New(),
		ProjectID: project.ID,
		Name:      fmt.Sprintf("%s_db", req.Name),
		Host:      "localhost",
		Port:      5432,
		Username:  "postgres",
		Password:  "password",
		Database:  fmt.Sprintf("project_%s", project.ID.String()[:8]),
		SSLMode:   "disable",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.dbRepo.Create(ctx, database); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	project.DatabaseURL = database.GetConnectionString()
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project with database URL: %w", err)
	}

	response := project.ToResponse()
	return &response, nil
}

// GetProjectByID retrieves a project by ID
func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*entities.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	response := project.ToResponse()
	return &response, nil
}

// GetProjectsByOwnerID retrieves all projects for a user
func (s *ProjectService) GetProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.ProjectResponse, error) {
	projects, err := s.projectRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var responses []*entities.ProjectResponse
	for _, project := range projects {
		response := project.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	project.Name = req.Name
	project.Description = req.Description
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	response := project.ToResponse()
	return &response, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}
