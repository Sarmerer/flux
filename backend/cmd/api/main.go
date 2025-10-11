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
	postgres "github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/handlers"
	postgresRepo "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/flow/internal/infrastructure/routes"
	"github.com/flow/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create database connection
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	ctx := context.Background()
	db, err := postgres.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := postgres.Migrate(ctx, db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)
	projectRepo := postgresRepo.NewProjectRepository(db)
	databaseRepo := postgresRepo.NewDatabaseRepository(db)
	tableRepo := postgresRepo.NewTableRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret)
	projectService := services.NewProjectService(projectRepo, databaseRepo)
	tableService := services.NewTableService(tableRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	tableHandler := handlers.NewTableHandler(tableService)

	// Setup routes
	router := routes.SetupRoutes(userHandler, projectHandler, tableHandler)

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
