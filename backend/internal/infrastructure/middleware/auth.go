package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
)

const (
	UserIDKey        string = "user_id"
	ProjectMemberKey string = "project_member"
	ProjectIDKey     string = "project_id"
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
				errors.WriteError(w, errors.NewUnauthorizedError("Authorization header missing"))
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				errors.WriteError(w, errors.NewUnauthorizedError("Invalid authorization header format"))
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
				errors.WriteError(w, errors.NewUnauthorizedError("Invalid or expired token"))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				errors.WriteError(w, errors.NewUnauthorizedError("Invalid token claims"))
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				errors.WriteError(w, errors.NewUnauthorizedError("User ID not found in token"))
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				errors.WriteError(w, errors.NewUnauthorizedError("Invalid user ID in token"))
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

func RequireProjectPermission(projectMemberRepo repositories.ProjectMemberRepository, permission entities.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				errors.WriteError(w, errors.NewUnauthorizedError("User not authenticated"))
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				errors.WriteError(w, errors.NewValidationError("Project ID required").WithField("projectId"))
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				errors.WriteError(w, errors.NewValidationError("Invalid project ID").WithField("projectId"))
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				errors.WriteError(w, errors.NewNotFoundError("Project member not found"))
				return
			}

			if !member.HasPermission(permission) {
				errors.WriteError(w, errors.NewForbiddenError("Insufficient permissions"))
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireProjectMembership(projectMemberRepo repositories.ProjectMemberRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				errors.WriteError(w, errors.NewUnauthorizedError("User not authenticated"))
				return
			}

			projectIDStr := chi.URLParam(r, "projectId")
			if projectIDStr == "" {
				projectIDStr = chi.URLParam(r, "id")
			}

			if projectIDStr == "" {
				errors.WriteError(w, errors.NewValidationError("Project ID required").WithField("projectId"))
				return
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				errors.WriteError(w, errors.NewValidationError("Invalid project ID").WithField("projectId"))
				return
			}

			member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
			if err != nil {
				errors.WriteError(w, errors.NewNotFoundError("Project member not found"))
				return
			}

			ctx := context.WithValue(r.Context(), ProjectMemberKey, member)
			ctx = context.WithValue(ctx, ProjectIDKey, projectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
