package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers"
	"github.com/mauriciofsnts/hermes/internal/server/middleware"
	"github.com/mauriciofsnts/hermes/internal/server/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/mauriciofsnts/hermes/docs" // Swagger docs
)

func StartServer(providers *providers.Providers) error {
	r := chi.NewRouter()

	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.MetricsMiddleware)

	options := cors.Options{
		AllowedOrigins: config.Hermes.Http.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         300,
	}

	r.Use(cors.Handler(options))

	// Register Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Register Swagger UI endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

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
