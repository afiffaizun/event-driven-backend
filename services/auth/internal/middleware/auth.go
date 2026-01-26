package middleware

import (
	"net/http"
	"strings"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(auth, " ")
			if len(parts) != 2 {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, err := security.ValidateToken(parts[1], secret)
			if err != nil || claims["type"] != "access" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
