package routes

import (
	"net/http"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/handlers"
	"github.com/flow/internal/infrastructure/logging"
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
	tableDataHandler *handlers.TableDataHandler,
	workflowHandler *handlers.WorkflowHandler,
	logHandler *handlers.LogHandler,
	projectMemberRepo repositories.ProjectMemberRepository,
	appLogger *logging.Logger,
	jwtSecret string,
	corsAllowedOrigins []string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(authMiddleware.LoggingMiddleware(appLogger))
	r.Use(errors.ErrorHandler)
	r.Use(middleware.Recoverer)

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

			r.Route("/logs", func(r chi.Router) {
				r.Get("/", logHandler.GetLogs)
				r.Delete("/cleanup", logHandler.DeleteOldLogs)
			})

			r.Route("/projects", func(r chi.Router) {

				r.Get("/", projectHandler.GetProjects)

				r.Post("/", projectHandler.CreateProject)

				r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
					Get("/{id}", projectHandler.GetProject)

				r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermProjectEdit)).
					Put("/{id}", projectHandler.UpdateProject)

				r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermProjectDelete)).
					Delete("/{id}", projectHandler.DeleteProject)

				r.Route("/{projectId}/members", func(r chi.Router) {
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", projectMemberHandler.GetProjectMembers)

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/me", projectMemberHandler.GetMyProjectRole)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUserManage)).
						Post("/", projectMemberHandler.AddProjectMember)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUserManage)).
						Put("/{memberId}", projectMemberHandler.UpdateProjectMember)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermUserManage)).
						Delete("/{memberId}", projectMemberHandler.RemoveProjectMember)
				})

				r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
					Get("/{projectId}/logs", logHandler.GetProjectLogs)

				r.Route("/{projectId}/tables", func(r chi.Router) {

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", tableHandler.GetTables)
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/{id}", tableHandler.GetTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableCreate)).
						Post("/", tableHandler.CreateTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableEdit)).
						Put("/{id}", tableHandler.UpdateTable)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableDelete)).
						Delete("/{id}", tableHandler.DeleteTable)

					r.Route("/{id}/schema", func(r chi.Router) {
						r.Use(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableEdit))

						r.Put("/update", tableHandler.UpdateTableSchema)
						r.Post("/columns", tableHandler.AddColumn)
						r.Delete("/columns", tableHandler.RemoveColumn)
						r.Put("/columns", tableHandler.ModifyColumn)
						r.Post("/foreign-keys", tableHandler.AddForeignKey)
						r.Delete("/foreign-keys", tableHandler.RemoveForeignKey)
					})

					r.Route("/{id}/data", func(r chi.Router) {
						r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
							Get("/", tableDataHandler.GetTableData)

						r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableEdit)).
							Post("/", tableDataHandler.InsertRow)

						r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableEdit)).
							Put("/{rowId}", tableDataHandler.UpdateRow)

						r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermTableEdit)).
							Delete("/{rowId}", tableDataHandler.DeleteRow)
					})
				})

				r.Route("/{projectId}/workflows", func(r chi.Router) {

					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/", workflowHandler.GetWorkflows)
					r.With(authMiddleware.RequireProjectMembership(projectMemberRepo)).
						Get("/{id}", workflowHandler.GetWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowCreate)).
						Post("/", workflowHandler.CreateWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowEdit)).
						Put("/{id}", workflowHandler.UpdateWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowEdit)).
						Post("/{id}/toggle", workflowHandler.ToggleWorkflowActive)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowEdit)).
						Post("/{id}/execute", workflowHandler.ExecuteWorkflow)

					r.With(authMiddleware.RequireProjectPermission(projectMemberRepo, entities.PermWorkflowDelete)).
						Delete("/{id}", workflowHandler.DeleteWorkflow)
				})
			})
		})
	})

	return r
}
