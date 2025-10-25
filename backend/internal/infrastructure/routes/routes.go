package routes

import (
	"net/http"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/handlers"
	authMiddleware "github.com/flow/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	projectMemberHandler *handlers.ProjectMemberHandler,
	tableHandler *handlers.TableHandler,
	dbMutationProgressHandler *handlers.DatabaseMutationProgressHandler,
	tableSchemaMutationHandler *handlers.TableSchemaMutationHandler,
	workflowHandler *handlers.WorkflowHandler,
	logHandler *handlers.LogHandler,
	projectMemberRepo repositories.ProjectMemberRepository,
	jwtSecret string,
	corsAllowedOrigins []string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(errors.ErrorHandler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		r.Route("/", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware(jwtSecret))

			r.Get("/me", userHandler.GetCurrentUser)

			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", userHandler.GetProfile)
			})

			r.Route("/operations", func(r chi.Router) {
				r.Get("/", dbMutationProgressHandler.GetUserOperations)
				r.Get("/{operationId}", dbMutationProgressHandler.GetOperationStatus)
				r.Delete("/{operationId}", dbMutationProgressHandler.CancelOperation)
			})

			r.Route("/logs", func(r chi.Router) {
				r.Get("/", logHandler.GetLogs)
				r.Delete("/cleanup", logHandler.DeleteOldLogs)
			})

			r.Route("/projects", func(r chi.Router) {

				r.Get("/", projectHandler.GetProjects)

				r.Post("/", projectHandler.CreateProject)

				r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
					Get("/{id}", projectHandler.GetProject)

				r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermProjectsEdit)).
					Put("/{id}", projectHandler.UpdateProject)

				r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermProjectsDelete)).
					Delete("/{id}", projectHandler.DeleteProject)

				r.Route("/{projectId}/members", func(r chi.Router) {
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", projectMemberHandler.GetProjectMembers)

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/me", projectMemberHandler.GetMyProjectRole)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUsersManage)).
						Post("/", projectMemberHandler.AddProjectMember)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUsersManage)).
						Put("/{memberId}", projectMemberHandler.UpdateProjectMember)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUsersManage)).
						Delete("/{memberId}", projectMemberHandler.RemoveProjectMember)
				})

				r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
					Get("/{projectId}/logs", logHandler.GetProjectLogs)

				r.Route("/{projectId}/tables", func(r chi.Router) {

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", tableHandler.GetTables)
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/{id}", tableHandler.GetTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTablesCreate)).
						Post("/", tableHandler.CreateTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTablesCreate)).
						Post("/with-progress", dbMutationProgressHandler.CreateTableWithProgress)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTablesEdit)).
						Put("/{id}", tableHandler.UpdateTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTablesDelete)).
						Delete("/{id}", tableHandler.DeleteTable)

					r.Route("/{tableName}/schema", func(r chi.Router) {
						r.Use(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTablesEdit))

						r.Post("/create", tableSchemaMutationHandler.CreateTableInDatabase)
						r.Delete("/drop", tableSchemaMutationHandler.DropTableFromDatabase)
						r.Post("/columns/add", tableSchemaMutationHandler.AddColumnToTable)
						r.Delete("/columns/remove", tableSchemaMutationHandler.RemoveColumnFromTable)
						r.Put("/columns/modify", tableSchemaMutationHandler.ModifyColumnInTable)
						r.Post("/foreign-keys/add", tableSchemaMutationHandler.AddForeignKeyToTable)
						r.Delete("/foreign-keys/remove", tableSchemaMutationHandler.RemoveForeignKeyFromTable)
					})
				})

				r.Route("/{projectId}/workflows", func(r chi.Router) {

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", workflowHandler.GetWorkflows)
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/{id}", workflowHandler.GetWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowsCreate)).
						Post("/", workflowHandler.CreateWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowsEdit)).
						Put("/{id}", workflowHandler.UpdateWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowsEdit)).
						Post("/{id}/toggle", workflowHandler.ToggleWorkflowActive)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowsEdit)).
						Post("/{id}/execute", workflowHandler.ExecuteWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowsDelete)).
						Delete("/{id}", workflowHandler.DeleteWorkflow)
				})

				r.Get("/ws", dbMutationProgressHandler.WebSocketHandler)
			})
		})
	})

	return r
}
