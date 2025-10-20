package routes

import (
	"net/http"

	"github.com/flow/internal/infrastructure/handlers"
	authMiddleware "github.com/flow/internal/infrastructure/middleware"
	"github.com/flow/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	tableHandler *handlers.TableHandler,
	dbMutationProgressHandler *handlers.DatabaseMutationProgressHandler,
	tableSchemaMutationHandler *handlers.TableSchemaMutationHandler,
	workflowHandler *handlers.WorkflowHandler,
	logHandler *handlers.LogHandler,
	jwtSecret string,
	corsAllowedOrigins []string,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(errors.ErrorHandler) // Custom error handler
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS configuration - use configured origins instead of wildcard
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		// Protected routes (authentication required)
		r.Route("/", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware(jwtSecret)) // Add authentication middleware

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", userHandler.GetProfile)
			})

			// Operation management routes
			r.Route("/operations", func(r chi.Router) {
				r.Get("/", dbMutationProgressHandler.GetUserOperations)
				r.Get("/{operationId}", dbMutationProgressHandler.GetOperationStatus)
				r.Delete("/{operationId}", dbMutationProgressHandler.CancelOperation)
			})

			// Log routes
			r.Route("/logs", func(r chi.Router) {
				r.Get("/", logHandler.GetLogs)
				r.Delete("/cleanup", logHandler.DeleteOldLogs)
			})

			// Project routes
			r.Route("/projects", func(r chi.Router) {
				r.Post("/", projectHandler.CreateProject)
				r.Get("/", projectHandler.GetProjects)
				r.Get("/{id}", projectHandler.GetProject)
				r.Put("/{id}", projectHandler.UpdateProject)
				r.Delete("/{id}", projectHandler.DeleteProject)

				// Project-specific log routes
				r.Get("/{projectId}/logs", logHandler.GetProjectLogs)

				// Table routes under projects
				r.Route("/{projectId}/tables", func(r chi.Router) {
					r.Post("/", tableHandler.CreateTable)
					r.Get("/", tableHandler.GetTables)
					r.Get("/{id}", tableHandler.GetTable)
					r.Put("/{id}", tableHandler.UpdateTable)
					r.Delete("/{id}", tableHandler.DeleteTable)

					// Table routes with progress tracking
					r.Post("/with-progress", dbMutationProgressHandler.CreateTableWithProgress)

					// Table schema mutation routes
					r.Route("/{tableName}/schema", func(r chi.Router) {
						r.Post("/create", tableSchemaMutationHandler.CreateTableInDatabase)
						r.Delete("/drop", tableSchemaMutationHandler.DropTableFromDatabase)
						r.Post("/columns/add", tableSchemaMutationHandler.AddColumnToTable)
						r.Delete("/columns/remove", tableSchemaMutationHandler.RemoveColumnFromTable)
						r.Put("/columns/modify", tableSchemaMutationHandler.ModifyColumnInTable)
						r.Post("/foreign-keys/add", tableSchemaMutationHandler.AddForeignKeyToTable)
						r.Delete("/foreign-keys/remove", tableSchemaMutationHandler.RemoveForeignKeyFromTable)
					})
				})

				// Workflow routes under projects
				r.Route("/{projectId}/workflows", func(r chi.Router) {
					r.Post("/", workflowHandler.CreateWorkflow)
					r.Get("/", workflowHandler.GetWorkflows)
					r.Get("/{id}", workflowHandler.GetWorkflow)
					r.Put("/{id}", workflowHandler.UpdateWorkflow)
					r.Delete("/{id}", workflowHandler.DeleteWorkflow)
					r.Post("/{id}/toggle", workflowHandler.ToggleWorkflowActive)
					r.Post("/{id}/execute", workflowHandler.ExecuteWorkflow)
				})

				// Public WebSocket endpoint (auth via WS handshake)
				r.Get("/ws", dbMutationProgressHandler.WebSocketHandler)
			})
		})
	})

	return r
}
