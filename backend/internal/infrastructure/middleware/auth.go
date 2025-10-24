package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
)

type contextKey string

const (
	UserIDKey        contextKey = "user_id"
	ProjectMemberKey contextKey = "project_member"
	ProjectIDKey     contextKey = "project_id"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				respondWithError(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				respondWithError(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respondWithError(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				respondWithError(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				respondWithError(w, "Invalid user ID format", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetProjectMemberFromContext(ctx context.Context) (*entities.ProjectMember, bool) {
	member, ok := ctx.Value(ProjectMemberKey).(*entities.ProjectMember)
	return member, ok
}

func GetProjectIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	projectID, ok := ctx.Value(ProjectIDKey).(uuid.UUID)
	return projectID, ok
}

func RequireProjectPermission(projectMemberRepo repositories.ProjectMemberRepository, permission entities.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				respondWithError(w, "Project ID required", http.StatusBadRequest)
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				respondWithError(w, "Invalid project ID", http.StatusBadRequest)
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				respondWithError(w, "Not a member of this project", http.StatusForbidden)
				return
			}

			if !member.HasPermission(permission) {
				respondWithError(w, "Insufficient permissions for this project", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireProjectAnyPermission(projectMemberRepo repositories.ProjectMemberRepository, permissions ...entities.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				respondWithError(w, "Project ID required", http.StatusBadRequest)
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				respondWithError(w, "Invalid project ID", http.StatusBadRequest)
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				respondWithError(w, "Not a member of this project", http.StatusForbidden)
				return
			}

			if !member.HasAnyPermission(permissions...) {
				respondWithError(w, "Insufficient permissions for this project", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireProjectRole(projectMemberRepo repositories.ProjectMemberRepository, roles ...entities.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				respondWithError(w, "Project ID required", http.StatusBadRequest)
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				respondWithError(w, "Invalid project ID", http.StatusBadRequest)
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				respondWithError(w, "Not a member of this project", http.StatusForbidden)
				return
			}

			if !member.HasRole(roles...) {
				respondWithError(w, "Insufficient permissions for this project", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireProjectAdmin(projectMemberRepo repositories.ProjectMemberRepository) func(http.Handler) http.Handler {
	return RequireProjectRole(projectMemberRepo, entities.RoleAdmin)
}

func RequireProjectMembership(projectMemberRepo repositories.ProjectMemberRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				respondWithError(w, "Project ID required", http.StatusBadRequest)
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				respondWithError(w, "Invalid project ID", http.StatusBadRequest)
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				respondWithError(w, "Not a member of this project", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(statusCode),
		},
	})
}
