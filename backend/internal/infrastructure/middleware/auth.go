package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"flow-backend/internal/domain/entities"
	"flow-backend/internal/domain/repositories"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserKey     contextKey = "user"
	UserRoleKey contextKey = "user_role"
)

// Claims represents JWT claims with role and permissions
type Claims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				respondWithError(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				respondWithError(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respondWithError(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Get user ID from claims
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

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(ctx context.Context) (*entities.User, bool) {
	user, ok := ctx.Value(UserKey).(*entities.User)
	return user, ok
}

// RequirePermission middleware checks if user has required permission
func RequirePermission(userRepo repositories.UserRepository, permission entities.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				respondWithError(w, "User not found", http.StatusUnauthorized)
				return
			}

			if !user.HasPermission(permission) {
				respondWithError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user to context for handlers to use
			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAnyPermission middleware checks if user has any of the required permissions
func RequireAnyPermission(userRepo repositories.UserRepository, permissions ...entities.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				respondWithError(w, "User not found", http.StatusUnauthorized)
				return
			}

			if !user.HasAnyPermission(permissions...) {
				respondWithError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware checks if user has one of the required roles
func RequireRole(userRepo repositories.UserRepository, roles ...entities.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				respondWithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				respondWithError(w, "User not found", http.StatusUnauthorized)
				return
			}

			if !user.HasRole(roles...) {
				respondWithError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin middleware checks if user is an admin
func RequireAdmin(userRepo repositories.UserRepository) func(http.Handler) http.Handler {
	return RequireRole(userRepo, entities.RoleAdmin)
}

// respondWithError sends a JSON error response
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
