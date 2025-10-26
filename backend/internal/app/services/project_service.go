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
	"github.com/flow/internal/utils"
	"github.com/flow/internal/validation"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectService struct {
	projectRepo       repositories.ProjectRepository
	dbRepo            repositories.DatabaseRepository
	projectMemberRepo repositories.ProjectMemberRepository
	pgService         *database.PostgreSQLManagementService
	migrationRunner   *database.ProjectMigrationRunner
	coreDB            *pgxpool.Pool
	logger            *logging.Logger
}

func NewProjectService(
	projectRepo repositories.ProjectRepository,
	dbRepo repositories.DatabaseRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	pgService *database.PostgreSQLManagementService,
	migrationRunner *database.ProjectMigrationRunner,
	coreDB *pgxpool.Pool,
	logger *logging.Logger,
) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		dbRepo:            dbRepo,
		projectMemberRepo: projectMemberRepo,
		pgService:         pgService,
		migrationRunner:   migrationRunner,
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
		APIKey:      utils.GenerateAPIKey(),
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

	member := &entities.ProjectMember{
		ID:          uuid.New(),
		ProjectID:   project.ID,
		UserID:      ownerID,
		Role:        entities.RoleAdmin,
		Permissions: allPermissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var database *entities.Database
	if s.pgService != nil {
		database = &entities.Database{
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
	}

	err := utils.WithTransaction(ctx, s.coreDB, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.projectRepo.CreateTx(ctx, tx, project); err != nil {
			logger.Error("Failed to create project", err)
			return errors.NewDatabaseError(err).WithDetails("failed to create project")
		}

		if err := s.projectMemberRepo.CreateTx(ctx, tx, member); err != nil {
			logger.Error("Failed to assign creator as admin", err)
			return errors.NewDatabaseError(err).WithDetails("failed to assign project role")
		}

		if database != nil {
			if err := s.dbRepo.CreateTx(ctx, tx, database); err != nil {
				logger.Error("Failed to insert database configuration", err)
				return errors.NewDatabaseError(err).WithDetails("failed to create database configuration")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if database != nil {
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

		if err := s.migrationRunner.InitializeProjectDatabase(ctx, database); err != nil {
			dbLogger.Error("Failed to initialize project database schema", err, map[string]interface{}{
				"database_name": database.Database,
			})
			_ = s.pgService.DropDatabase(ctx, database)
			return nil, errors.NewDatabaseError(err).WithDetails("failed to initialize project database schema")
		}

		dbLogger.Info("Successfully initialized project database schema", map[string]interface{}{
			"database_name": database.Database,
		})
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
