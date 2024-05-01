package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/api/health"
	"github.com/mauriciofsnts/hermes/internal/server/api/notify"
	"github.com/mauriciofsnts/hermes/internal/server/api/template"
)

func StartServer() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	Router(r)

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

func Router(root *chi.Mux) {
	root.Route("/api/v1", func(r chi.Router) {
		hc := health.NewHealthController()
		r.Get("/health", hc.Health)

		r.Route("/notify", func(r chi.Router) {
			nc := notify.NewEmailController()
			r.Post("/", nc.SendPlainTextEmail)
			r.Post("/{slug}", nc.SendTemplateEmail)
		})

		r.Route("/templates", func(r chi.Router) {
			tc := template.NewTemplateController()
			r.Get("/{slug}", tc.GetRaw)
			r.Post("/", tc.Create)
		})
	})
}
