package routes

import (
	"net/http"

	"github.com/flow/internal/infrastructure/handlers"
	authMiddleware "github.com/flow/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	tableHandler *handlers.TableHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
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
			r.Use(authMiddleware.AuthMiddleware) // Add authentication middleware

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", userHandler.GetProfile)
			})

			// Project routes
			r.Route("/projects", func(r chi.Router) {
				r.Post("/", projectHandler.CreateProject)
				r.Get("/", projectHandler.GetProjects)
				r.Get("/{id}", projectHandler.GetProject)
				r.Put("/{id}", projectHandler.UpdateProject)
				r.Delete("/{id}", projectHandler.DeleteProject)

				// Table routes under projects
				r.Route("/{projectId}/tables", func(r chi.Router) {
					r.Post("/", tableHandler.CreateTable)
					r.Get("/", tableHandler.GetTables)
					r.Get("/{id}", tableHandler.GetTable)
					r.Put("/{id}", tableHandler.UpdateTable)
					r.Delete("/{id}", tableHandler.DeleteTable)
				})
			})
		})
	})

	return r
}
