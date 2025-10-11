package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
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
			// Return secret key (in production, this should come from config)
			return []byte("your-secret-key-change-this-in-production"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get user ID from claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
