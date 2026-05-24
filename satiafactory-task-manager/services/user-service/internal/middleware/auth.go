package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dratbo/satisfactory-task-manager/user-service/internal/config"
	"github.com/dratbo/satisfactory-task-manager/user-service/internal/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims, err := jwt.ValidateToken(tokenString, cfg.JWTSecret) // ← изменено
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
