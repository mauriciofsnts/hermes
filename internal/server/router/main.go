package router

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/mauriciofsnts/hermes/internal/providers"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/server/api/health"
	"github.com/mauriciofsnts/hermes/internal/server/api/notification"
	"github.com/mauriciofsnts/hermes/internal/server/api/template"
)

func RouteApp(root *chi.Mux, provider *providers.Providers) {
	root.Route("/api/v1", routeAPI(provider))
}

func routeAPI(providers *providers.Providers) func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/health", asChiRouter(routeHealth(providers)))

		r.Route("/app", func(r chi.Router) {
			slog.Info("app route")
			// r.Use(hermesMiddleware.AuthMiddleware)
			r.Route("/notify", asChiRouter(routeNotify(providers)))
			r.Route("/templates", asChiRouter(routeTemplate(providers)))
		})
	}
}

func routeHealth(providers *providers.Providers) func(api.Router) {
	return func(r api.Router) {
		controller := health.NewHealthController(providers.Queue)
		controller.Route(r)
	}
}

func routeNotify(providers *providers.Providers) func(api.Router) {
	return func(r api.Router) {
		controller := notification.NewEmailController(providers.Storage, providers.Queue)
		controller.Route(r)
	}
}

func routeTemplate(providers *providers.Providers) func(api.Router) {
	return func(r api.Router) {
		controller := template.NewTemplateController(providers.Storage)
		controller.Route(r)
	}
}

func asChiRouter(fn func(api.Router)) func(chi.Router) {
	return func(r chi.Router) {
		fn(&ChiRouter{r})
	}
}
