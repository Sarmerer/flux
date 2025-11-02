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

	"github.com/flow/internal/app"
	"github.com/flow/internal/config"
	"github.com/flow/internal/infrastructure/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()

	log.Println("Initializing application container...")
	container, err := app.NewContainer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	log.Println("Application initialized successfully")

	container.Logger.Info("Application starting", map[string]interface{}{
		"server_host": cfg.Server.Host,
		"server_port": cfg.Server.Port,
	})

	router := routes.SetupRoutes(
		container.UserHandler,
		container.ProjectHandler,
		container.ProjectMemberHandler,
		container.TableHandler,
		container.TableDataHandler,
		container.DbMutationProgressHandler,
		container.TableSchemaMutationHandler,
		container.WorkflowHandler,
		container.LogHandler,
		container.ProjectMemberRepo,
		container.Logger,
		cfg.JWT.Secret,
		cfg.CORS.AllowedOrigins,
	)

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
