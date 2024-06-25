package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers"
	"github.com/mauriciofsnts/hermes/internal/server/middleware"
	"github.com/mauriciofsnts/hermes/internal/server/router"
)

func StartServer(providers *providers.Providers) error {
	r := chi.NewRouter()

	for _, origin := range config.Hermes.Http.AllowedOrigins {
		slog.Info("Allowed origin: %s", origin, nil)
	}

	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.LoggerMiddleware)

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
