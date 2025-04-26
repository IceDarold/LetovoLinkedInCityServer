package middleware

import (
	"net/http"
	"strings"
)

// Проверка авторизации по заголовку Authorization: Bearer <token>
func AuthMiddleware(validToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != validToken {
				http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
				return
			}

			// токен валиден — передаём дальше
			next.ServeHTTP(w, r)
		})
	}
}
