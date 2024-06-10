package middleware

import (
	"github.com/go-chi/cors"
	"github.com/mauriciofsnts/hermes/internal/config"
)

var CorsConfig = cors.Options{
	AllowedOrigins: config.Hermes.Http.AllowedOrigins,
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	MaxAge:         300, // Maximum value not ignored by any of major browsers
	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
}
