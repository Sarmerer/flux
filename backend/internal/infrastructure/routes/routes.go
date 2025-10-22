package routes

import (
	"net/http"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
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
	userRepo repositories.UserRepository,
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
				// Everyone can view projects
				r.Get("/", projectHandler.GetProjects)
				r.Get("/{id}", projectHandler.GetProject)

				// Require permissions for create/update/delete
				r.With(authMiddleware.RequirePermission(userRepo, entities.PermProjectsCreate)).
					Post("/", projectHandler.CreateProject)

				r.With(authMiddleware.RequirePermission(userRepo, entities.PermProjectsEdit)).
					Put("/{id}", projectHandler.UpdateProject)

				r.With(authMiddleware.RequirePermission(userRepo, entities.PermProjectsDelete)).
					Delete("/{id}", projectHandler.DeleteProject)

				// Project-specific log routes
				r.Get("/{projectId}/logs", logHandler.GetProjectLogs)

				// Table routes under projects
				r.Route("/{projectId}/tables", func(r chi.Router) {
					// Everyone can view tables
					r.Get("/", tableHandler.GetTables)
					r.Get("/{id}", tableHandler.GetTable)

					// Require permissions for create/update/delete
					r.With(authMiddleware.RequirePermission(userRepo, entities.PermTablesCreate)).
						Post("/", tableHandler.CreateTable)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermTablesCreate)).
						Post("/with-progress", dbMutationProgressHandler.CreateTableWithProgress)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermTablesEdit)).
						Put("/{id}", tableHandler.UpdateTable)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermTablesDelete)).
						Delete("/{id}", tableHandler.DeleteTable)

					// Table schema mutation routes - require edit permission
					r.Route("/{tableName}/schema", func(r chi.Router) {
						r.Use(authMiddleware.RequirePermission(userRepo, entities.PermTablesEdit))

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
					// Everyone can view workflows
					r.Get("/", workflowHandler.GetWorkflows)
					r.Get("/{id}", workflowHandler.GetWorkflow)

					// Require permissions for create/update/delete
					r.With(authMiddleware.RequirePermission(userRepo, entities.PermWorkflowsCreate)).
						Post("/", workflowHandler.CreateWorkflow)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermWorkflowsEdit)).
						Put("/{id}", workflowHandler.UpdateWorkflow)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermWorkflowsEdit)).
						Post("/{id}/toggle", workflowHandler.ToggleWorkflowActive)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermWorkflowsEdit)).
						Post("/{id}/execute", workflowHandler.ExecuteWorkflow)

					r.With(authMiddleware.RequirePermission(userRepo, entities.PermWorkflowsDelete)).
						Delete("/{id}", workflowHandler.DeleteWorkflow)
				})

				// Public WebSocket endpoint (auth via WS handshake)
				r.Get("/ws", dbMutationProgressHandler.WebSocketHandler)
			})
		})
	})

	return r
}
