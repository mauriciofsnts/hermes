package middleware

import (
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/config"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")

		_, ok := config.Hermes.AppsByAPIKey[apiKey]

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
