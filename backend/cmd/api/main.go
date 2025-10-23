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
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create database connection
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

	// Initialize basic logger for migrations (no WebSocket streaming yet)
	logStorage := logging.NewPostgresLogStorage(db)
	// Use nil streamer during migrations since hub isn't ready yet
	migrationLogger := logging.NewDevelopmentLogger(logStorage, nil)

	// Run migrations automatically on startup using SQL files and structured logging
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(ctx, db, migrationLogger); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)
	projectRepo := postgresRepo.NewProjectRepository(db)
	databaseRepo := postgresRepo.NewDatabaseRepository(db)
	tableRepo := postgresRepo.NewTableRepository(db)
	workflowRepo := postgresRepo.NewWorkflowRepository(db)

	// Initialize Realtime Service
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

	// Start realtime service
	if err := realtimeService.Start(ctx); err != nil {
		log.Fatalf("Failed to start realtime service: %v", err)
	}
	defer realtimeService.Stop()

	// Get hub and create progress tracker
	wsHub := realtimeService.GetHub()
	progressTracker := progress.NewTracker(wsHub)

	// Start progress tracker cleanup routine
	progressTracker.StartCleanupRoutine(ctx, 5*time.Minute, 1*time.Hour)

	// Update log streamer with WebSocket hub (was initialized earlier without hub)
	logStreamer := logging.NewWebSocketLogStreamer(wsHub)
	appLogger := logging.NewDevelopmentLogger(logStorage, logStreamer)

	log.Println("Structured logging system initialized")

	// Log application startup
	appLogger.Info("Application starting", map[string]interface{}{
		"server_host": cfg.Server.Host,
		"server_port": cfg.Server.Port,
	})

	// Initialize infrastructure services (connection management, PostgreSQL operations)
	connService := database.NewConnectionService(db)
	pgManagementService := database.NewPostgreSQLManagementService(connService)

	// Initialize mutation services
	dbService := services.NewDatabaseService(databaseRepo, projectRepo, pgManagementService, connService, progressTracker, wsHub)
	tableSchemaMutationService := services.NewTableSchemaMutationService(tableRepo, databaseRepo, projectRepo)

	// Initialize services (after database service is available)
	userService := services.NewUserService(userRepo, cfg.JWT.Secret)
	projectService := services.NewProjectService(projectRepo, databaseRepo, pgManagementService, db, appLogger)
	tableService := services.NewTableService(tableRepo)
	workflowService := services.NewWorkflowService(workflowRepo, projectRepo, databaseRepo, pgManagementService, appLogger)

	// Services wired with database handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	tableHandler := handlers.NewTableHandler(tableService)
	dbMutationProgressHandler := handlers.NewDatabaseMutationProgressHandler(dbService, progressTracker, wsHub)
	dbMutationProgressHandler.SetJWTSecret(cfg.JWT.Secret)
	tableSchemaMutationHandler := handlers.NewTableSchemaMutationHandler(tableSchemaMutationService)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	// Initialize log handler
	logHandler := handlers.NewLogHandler(logStorage, logStreamer)

	// Setup routes
	router := routes.SetupRoutes(userHandler, projectHandler, tableHandler, dbMutationProgressHandler, tableSchemaMutationHandler, workflowHandler, logHandler, userRepo, cfg.JWT.Secret, cfg.CORS.AllowedOrigins)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
