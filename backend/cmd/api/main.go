package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/config"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/handlers"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"
	postgresRepo "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/flow/internal/infrastructure/routes"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	ctx := context.Background()
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logStorage := logging.NewPostgresLogStorage(db)

	migrationLogger := logging.NewDevelopmentLogger(logStorage, nil)

	log.Println("Running database migrations...")
	if err := database.AutoMigrate(ctx, db, migrationLogger); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := postgresRepo.NewUserRepository(db)
	projectRepo := postgresRepo.NewProjectRepository(db)
	databaseRepo := postgresRepo.NewDatabaseRepository(db)
	tableRepo := postgresRepo.NewTableRepository(db)
	workflowRepo := postgresRepo.NewWorkflowRepository(db)
	projectMemberRepo := postgresRepo.NewProjectMemberRepository(db)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	realtimeService, err := realtime.NewService(dbURL, realtime.Config{
		DatabaseURL: dbURL,
		JWTSecret:   cfg.JWT.Secret,
	})
	if err != nil {
		log.Fatalf("Failed to create realtime service: %v", err)
	}

	if err := realtimeService.Start(ctx); err != nil {
		log.Fatalf("Failed to start realtime service: %v", err)
	}
	defer realtimeService.Stop()

	wsHub := realtimeService.GetHub()
	progressTracker := progress.NewTracker(wsHub)

	progressTracker.StartCleanupRoutine(ctx, 5*time.Minute, 1*time.Hour)

	logStreamer := logging.NewWebSocketLogStreamer(wsHub)
	appLogger := logging.NewDevelopmentLogger(logStorage, logStreamer)

	log.Println("Structured logging system initialized")

	appLogger.Info("Application starting", map[string]interface{}{
		"server_host": cfg.Server.Host,
		"server_port": cfg.Server.Port,
	})

	connService := database.NewConnectionService(db)
	pgManagementService := database.NewPostgreSQLManagementService(connService)
	dataManipulationService := database.NewDataManipulationService(connService)
	schemaManagementService := database.NewSchemaManagementService(connService)

	dbService := services.NewDatabaseService(databaseRepo, projectRepo, pgManagementService, connService, progressTracker, wsHub)
	tableSchemaMutationService := services.NewTableSchemaMutationService(tableRepo, databaseRepo, projectRepo, schemaManagementService)

	userService := services.NewUserService(userRepo, cfg.JWT.Secret)
	projectService := services.NewProjectService(projectRepo, databaseRepo, projectMemberRepo, pgManagementService, db, appLogger)
	tableService := services.NewTableService(tableRepo)
	workflowService := services.NewWorkflowService(workflowRepo, projectRepo, databaseRepo, dataManipulationService, appLogger)

	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberRepo, userRepo, projectRepo)
	tableHandler := handlers.NewTableHandler(tableService)
	dbMutationProgressHandler := handlers.NewDatabaseMutationProgressHandler(dbService, progressTracker, wsHub)
	dbMutationProgressHandler.SetJWTSecret(cfg.JWT.Secret)
	tableSchemaMutationHandler := handlers.NewTableSchemaMutationHandler(tableSchemaMutationService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	logHandler := handlers.NewLogHandler(logStorage, logStreamer)

	router := routes.SetupRoutes(userHandler, projectHandler, projectMemberHandler, tableHandler, dbMutationProgressHandler, tableSchemaMutationHandler, workflowHandler, logHandler, projectMemberRepo, cfg.JWT.Secret, cfg.CORS.AllowedOrigins)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
