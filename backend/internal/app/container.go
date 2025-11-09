package app

import (
	"context"
	"fmt"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/config"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/handlers"
	"github.com/flow/internal/infrastructure/logging"
	postgresRepo "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config *config.Config

	DB     *pgxpool.Pool
	Logger *logging.Logger

	UserRepo          repositories.UserRepository
	ProjectRepo       repositories.ProjectRepository
	DatabaseRepo      repositories.DatabaseRepository
	ProjectMemberRepo repositories.ProjectMemberRepository

	ConnService             *database.ConnectionService
	PgManagementService     *database.PostgreSQLManagementService
	SchemaManagementService *database.SchemaManagementService
	SchemaInspectionService *database.SchemaInspectionService
	ProjectConnResolver     *database.ProjectConnectionResolver
	MigrationRunner         *database.ProjectMigrationRunner

	UserService      *services.UserService
	ProjectService   *services.ProjectService
	TableService     *services.TableService
	TableDataService *services.TableDataService
	WorkflowService  *services.WorkflowService

	UserHandler          *handlers.UserHandler
	ProjectHandler       *handlers.ProjectHandler
	ProjectMemberHandler *handlers.ProjectMemberHandler
	TableHandler         *handlers.TableHandler
	TableDataHandler     *handlers.TableDataHandler
	WorkflowHandler      *handlers.WorkflowHandler
	LogHandler           *handlers.LogHandler
}

func NewContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	c := &Container{Config: cfg}

	if err := c.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("Failed to initialize database: %w", err)
	}

	if err := c.initInfrastructure(ctx); err != nil {
		return nil, fmt.Errorf("Failed to initialize infrastructure: %w", err)
	}

	c.initRepositories()
	c.initServices()
	c.initHandlers()

	return c, nil
}

func (c *Container) buildDatabaseConfig() database.Config {
	return database.Config{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}
}

func (c *Container) initDatabase(ctx context.Context) error {
	dbConfig := c.buildDatabaseConfig()

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		return err
	}
	c.DB = db

	migrationLogger := logging.NewDevelopmentLogger(nil, logging.VerbosityNormal)

	if err := database.AutoMigrate(ctx, db, migrationLogger); err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	return nil
}

func (c *Container) initInfrastructure(ctx context.Context) error {
	c.Logger = logging.NewDevelopmentLogger(nil, logging.VerbosityNormal)

	c.ConnService = database.NewConnectionService(c.DB)
	c.PgManagementService = database.NewPostgreSQLManagementService(c.ConnService)
	c.SchemaManagementService = database.NewSchemaManagementService(c.ConnService)
	c.SchemaInspectionService = database.NewSchemaInspectionService()

	return nil
}

func (c *Container) initRepositories() {
	c.UserRepo = postgresRepo.NewUserRepository(c.DB)
	c.ProjectRepo = postgresRepo.NewProjectRepository(c.DB)
	c.DatabaseRepo = postgresRepo.NewDatabaseRepository(c.DB)
	c.ProjectMemberRepo = postgresRepo.NewProjectMemberRepository(c.DB)

	c.ProjectConnResolver = database.NewProjectConnectionResolver(c.ConnService, c.DatabaseRepo)
	c.MigrationRunner = database.NewProjectMigrationRunner(c.ConnService, c.Logger)
}

func (c *Container) initServices() {
	repoFactory := services.NewProjectRepositoryFactory(c.ProjectConnResolver)

	c.UserService = services.NewUserService(c.UserRepo, c.Config.JWT.Secret)

	c.ProjectService = services.NewProjectService(
		c.ProjectRepo,
		c.DatabaseRepo,
		c.ProjectMemberRepo,
		c.PgManagementService,
		c.MigrationRunner,
		c.DB,
		c.Logger,
		services.ProjectDatabaseConfig{
			Host:     c.Config.ProjectDatabase.Host,
			Port:     c.Config.ProjectDatabase.Port,
			User:     c.Config.ProjectDatabase.User,
			Password: c.Config.ProjectDatabase.Password,
			SSLMode:  c.Config.ProjectDatabase.SSLMode,
		},
	)

	c.TableService = services.NewTableService(
		repoFactory,
		c.ProjectConnResolver,
		c.SchemaManagementService,
		c.SchemaInspectionService,
	)

	c.TableDataService = services.NewTableDataService(
		repoFactory,
		c.TableService,
	)

	c.WorkflowService = services.NewWorkflowService(
		repoFactory,
		c.ProjectRepo,
		c.DatabaseRepo,
		c.ProjectConnResolver,
		c.Logger,
	)
}

func (c *Container) initHandlers() {
	c.UserHandler = handlers.NewUserHandler(c.UserService)
	c.ProjectHandler = handlers.NewProjectHandler(c.ProjectService)
	c.ProjectMemberHandler = handlers.NewProjectMemberHandler(c.ProjectMemberRepo, c.UserRepo, c.ProjectRepo)
	c.TableHandler = handlers.NewTableHandler(c.TableService)
	c.TableDataHandler = handlers.NewTableDataHandler(c.TableDataService)
	c.WorkflowHandler = handlers.NewWorkflowHandler(c.WorkflowService)
	c.LogHandler = handlers.NewLogHandler(nil)
}

func (c *Container) Close() {
	if c.ProjectConnResolver != nil {
		c.ProjectConnResolver.Close()
	}
	if c.DB != nil {
		c.DB.Close()
	}
}
