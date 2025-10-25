package services

import (
	"context"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/flow/internal/validation"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectService struct {
	projectRepo       repositories.ProjectRepository
	dbRepo            repositories.DatabaseRepository
	projectMemberRepo repositories.ProjectMemberRepository
	pgService         *database.PostgreSQLManagementService
	coreDB            *pgxpool.Pool
	logger            *logging.Logger
}

func NewProjectService(
	projectRepo repositories.ProjectRepository,
	dbRepo repositories.DatabaseRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	pgService *database.PostgreSQLManagementService,
	coreDB *pgxpool.Pool,
	logger *logging.Logger,
) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		dbRepo:            dbRepo,
		projectMemberRepo: projectMemberRepo,
		pgService:         pgService,
		coreDB:            coreDB,
		logger:            logger,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *entities.ProjectCreateRequest, ownerID uuid.UUID) (*entities.ProjectResponse, error) {

	if err := validation.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}

	project := &entities.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		APIKey:      generateAPIKey(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	logger := s.logger.WithContext(logging.LogContext{
		UserID:    &ownerID,
		Operation: "create_project",
	})

	logger.Info("Creating new project", map[string]interface{}{
		"project_name": req.Name,
		"owner_id":     ownerID.String(),
	})

	tx, err := s.coreDB.Begin(ctx)
	if err != nil {
		logger.Error("Failed to begin transaction for project creation", err)
		return nil, errors.NewDatabaseError(err).WithDetails("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO projects (id, name, description, owner_id, database_url, api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if _, err := tx.Exec(ctx, query, project.ID, project.Name, project.Description, project.OwnerID, project.DatabaseURL, project.APIKey, project.CreatedAt, project.UpdatedAt); err != nil {
		logger.Error("Failed to insert project", err)
		return nil, errors.NewDatabaseError(err).WithDetails("failed to create project")
	}

	allPermissions := entities.Permissions{
		entities.PermProjectsCreate,
		entities.PermProjectsEdit,
		entities.PermProjectsDelete,
		entities.PermProjectsView,
		entities.PermWorkflowsCreate,
		entities.PermWorkflowsEdit,
		entities.PermWorkflowsDelete,
		entities.PermTablesCreate,
		entities.PermTablesEdit,
		entities.PermTablesDelete,
		entities.PermSettingsManage,
		entities.PermUsersManage,
	}
	memberQuery := `
		INSERT INTO project_members (id, project_id, user_id, role, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	memberID := uuid.New()
	if _, err := tx.Exec(ctx, memberQuery, memberID, project.ID, ownerID, entities.RoleAdmin, allPermissions, time.Now(), time.Now()); err != nil {
		logger.Error("Failed to assign creator as admin", err)
		return nil, errors.NewDatabaseError(err).WithDetails("failed to assign project role")
	}

	if s.pgService != nil {
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

		dbQuery := `
			INSERT INTO databases (id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		if _, err := tx.Exec(ctx, dbQuery, database.ID, database.ProjectID, database.Name, database.Host, database.Port, database.Username, database.Password, database.Database, database.SSLMode, database.CreatedAt, database.UpdatedAt); err != nil {
			logger.Error("Failed to insert database configuration", err)
			return nil, errors.NewDatabaseError(err).WithDetails("failed to create database configuration")
		}

		dbLogger := s.logger.WithContext(logging.LogContext{
			ProjectID:  &project.ID,
			DatabaseID: &database.ID,
			UserID:     &ownerID,
			Operation:  "create_project",
		})

		if err := s.pgService.CreateDatabase(ctx, database); err != nil {
			dbLogger.Error("Failed to create physical database", err, map[string]interface{}{
				"database_name": database.Database,
			})

			return nil, err
		}

		dbLogger.Info("Successfully created physical database", map[string]interface{}{
			"database_name": database.Database,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit transaction for project creation", err)
		return nil, errors.NewDatabaseError(err).WithDetails("failed to commit transaction")
	}

	projectLogger := s.logger.WithContext(logging.LogContext{
		ProjectID: &project.ID,
		UserID:    &ownerID,
		Operation: "create_project",
	})

	projectLogger.Info("Successfully created project with database", map[string]interface{}{
		"project_id":   project.ID.String(),
		"project_name": project.Name,
	})

	response := project.ToResponse()
	return &response, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*entities.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Project")
	}

	response := project.ToResponse()
	return &response, nil
}

func (s *ProjectService) GetProjectsByOwnerID(ctx context.Context, userID uuid.UUID) ([]*entities.ProjectResponse, error) {

	memberships, err := s.projectMemberRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to retrieve project memberships")
	}

	responses := make([]*entities.ProjectResponse, 0)
	for _, membership := range memberships {
		project, err := s.projectRepo.GetByID(ctx, membership.ProjectID)
		if err != nil {

			continue
		}
		response := project.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error) {

	if err := validation.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Project")
	}

	project.Name = req.Name
	project.Description = req.Description
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to update project")
	}

	response := project.ToResponse()
	return &response, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Project")
	}

	logger := s.logger.WithContext(logging.LogContext{
		ProjectID: &id,
		Operation: "delete_project",
	})

	logger.Info("Deleting project", map[string]interface{}{
		"project_name": project.Name,
	})

	databases, err := s.dbRepo.GetByProjectID(ctx, id)
	if err != nil {
		logger.Warn("Failed to retrieve databases for project", map[string]interface{}{
			"error": err.Error(),
		})

	}

	for _, database := range databases {
		dbLogger := s.logger.WithContext(logging.LogContext{
			ProjectID:  &id,
			DatabaseID: &database.ID,
			Operation:  "delete_project",
		})

		if err := s.pgService.DropDatabase(ctx, database); err != nil {

			dbLogger.Warn("Failed to drop physical database", map[string]interface{}{
				"database_name": database.Database,
				"error":         err.Error(),
			})
		} else {
			dbLogger.Info("Successfully dropped physical database", map[string]interface{}{
				"database_name": database.Database,
			})
		}
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete project from database", err)
		return errors.NewDatabaseError(err).WithDetails("failed to delete project")
	}

	logger.Info("Successfully deleted project and all associated resources", map[string]interface{}{
		"project_name": project.Name,
	})
	return nil
}
