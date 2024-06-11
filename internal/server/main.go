package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers"
	"github.com/mauriciofsnts/hermes/internal/server/middleware"
	"github.com/mauriciofsnts/hermes/internal/server/router"
)

func StartServer(providers *providers.Providers) error {
	r := chi.NewRouter()

	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.LoggerMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: config.Hermes.Http.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	router.RouteApp(r, providers)

	bindAddr := fmt.Sprintf(":%d", config.Hermes.Http.Port)
	slog.Info("Starting server on %s", bindAddr, nil)

	server := &http.Server{
		Addr:         bindAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
