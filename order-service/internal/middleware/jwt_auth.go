package middleware

import (
	"context"
	"net/http"

	"order-service/internal/service"
)

type contextKey string

const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

func JWTAuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, `{"error" : "не авторизован"}`, http.StatusUnauthorized)
				return
			}
			claims, err := authService.ValidateToken(cookie.Value)
			if err != nil {
				http.Error(w, `{"error" : "токен недействителен или истек"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
