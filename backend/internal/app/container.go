package app

import (
	"context"
	"fmt"
	"time"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/config"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/handlers"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"
	postgresRepo "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config *config.Config

	DB              *pgxpool.Pool
	RealtimeService *realtime.Service
	ProgressTracker *progress.Tracker
	Logger          *logging.Logger

	UserRepo          repositories.UserRepository
	ProjectRepo       repositories.ProjectRepository
	DatabaseRepo      repositories.DatabaseRepository
	ProjectMemberRepo repositories.ProjectMemberRepository

	ConnService             *database.ConnectionService
	PgManagementService     *database.PostgreSQLManagementService
	DataManipulationService *database.DataManipulationService
	SchemaManagementService *database.SchemaManagementService
	ProjectConnResolver     *database.ProjectConnectionResolver
	MigrationRunner         *database.ProjectMigrationRunner

	UserService     *services.UserService
	ProjectService  *services.ProjectService
	TableService    *services.TableService
	WorkflowService *services.WorkflowService
	DatabaseService *services.DatabaseService

	UserHandler                     *handlers.UserHandler
	ProjectHandler                  *handlers.ProjectHandler
	ProjectMemberHandler            *handlers.ProjectMemberHandler
	TableHandler                    *handlers.TableHandler
	DbMutationProgressHandler       *handlers.DatabaseMutationProgressHandler
	TableSchemaMutationHandler      *handlers.TableSchemaMutationHandler
	WorkflowHandler                 *handlers.WorkflowHandler
	LogHandler                      *handlers.LogHandler
}

func NewContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	c := &Container{Config: cfg}

	if err := c.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := c.initInfrastructure(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	c.initRepositories()
	c.initServices()
	c.initHandlers()

	return c, nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	dbConfig := database.Config{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		return err
	}
	c.DB = db

	logStorage := logging.NewPostgresLogStorage(db)
	migrationLogger := logging.NewDevelopmentLogger(logStorage, nil)

	if err := database.AutoMigrate(ctx, db, migrationLogger); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (c *Container) initInfrastructure(ctx context.Context) error {
	dbConfig := database.Config{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	dbURL := dbConfig.ConnectionString()

	realtimeService, err := realtime.NewService(dbURL, realtime.Config{
		DatabaseURL: dbURL,
		JWTSecret:   c.Config.JWT.Secret,
	})
	if err != nil {
		return fmt.Errorf("failed to create realtime service: %w", err)
	}

	if err := realtimeService.Start(ctx); err != nil {
		return fmt.Errorf("failed to start realtime service: %w", err)
	}
	c.RealtimeService = realtimeService

	wsHub := realtimeService.GetHub()
	c.ProgressTracker = progress.NewTracker(wsHub)
	c.ProgressTracker.StartCleanupRoutine(ctx, 5*time.Minute, 1*time.Hour)

	logStorage := logging.NewPostgresLogStorage(c.DB)
	logStreamer := logging.NewWebSocketLogStreamer(wsHub)
	c.Logger = logging.NewDevelopmentLogger(logStorage, logStreamer)

	c.ConnService = database.NewConnectionService(c.DB)
	c.PgManagementService = database.NewPostgreSQLManagementService(c.ConnService)
	c.DataManipulationService = database.NewDataManipulationService(c.ConnService)
	c.SchemaManagementService = database.NewSchemaManagementService(c.ConnService)

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

	wsHub := c.RealtimeService.GetHub()

	c.DatabaseService = services.NewDatabaseService(
		c.DatabaseRepo,
		c.ProjectRepo,
		c.PgManagementService,
		c.ConnService,
		c.ProgressTracker,
		wsHub,
	)

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

	c.TableService = services.NewTableService(repoFactory)

	c.WorkflowService = services.NewWorkflowService(
		repoFactory,
		c.ProjectRepo,
		c.DatabaseRepo,
		c.DataManipulationService,
		c.Logger,
	)
}

func (c *Container) initHandlers() {
	wsHub := c.RealtimeService.GetHub()
	logStorage := logging.NewPostgresLogStorage(c.DB)
	logStreamer := logging.NewWebSocketLogStreamer(wsHub)

	c.UserHandler = handlers.NewUserHandler(c.UserService)
	c.ProjectHandler = handlers.NewProjectHandler(c.ProjectService)
	c.ProjectMemberHandler = handlers.NewProjectMemberHandler(c.ProjectMemberRepo, c.UserRepo, c.ProjectRepo)
	c.TableHandler = handlers.NewTableHandler(c.TableService)

	c.DbMutationProgressHandler = handlers.NewDatabaseMutationProgressHandler(
		c.DatabaseService,
		c.ProgressTracker,
		wsHub,
		c.Config.JWT.Secret,
	)

	c.TableSchemaMutationHandler = handlers.NewTableSchemaMutationHandler(
		c.DatabaseRepo,
		c.ProjectRepo,
		c.SchemaManagementService,
	)

	c.WorkflowHandler = handlers.NewWorkflowHandler(c.WorkflowService)
	c.LogHandler = handlers.NewLogHandler(logStorage, logStreamer)
}

func (c *Container) Close() {
	if c.ProjectConnResolver != nil {
		c.ProjectConnResolver.Close()
	}
	if c.RealtimeService != nil {
		c.RealtimeService.Stop()
	}
	if c.DB != nil {
		c.DB.Close()
	}
}
