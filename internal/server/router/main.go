package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/server/api/health"
	"github.com/mauriciofsnts/hermes/internal/server/api/notification"

	hermesMiddleware "github.com/mauriciofsnts/hermes/internal/server/middleware"
)

func RouteApp(root *chi.Mux) {
	root.Route("/api/v1", routeAPI())
}

func routeAPI() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/health", asChiRouter(routeHealth()))

		r.Route("/app", func(r chi.Router) {
			r.Use(hermesMiddleware.AuthMiddleware)

			r.Route("/notify", asChiRouter(routeNotify()))
			r.Route("/templates", asChiRouter(routeTemplate()))
		})
	}
}

func routeHealth() func(api.Router) {
	return func(r api.Router) {
		controller := health.NewHealthController()
		controller.Route(r)
	}
}

func routeNotify() func(api.Router) {
	return func(r api.Router) {
		controller := notification.NewEmailController()
		controller.Route(r)
	}
}

func routeTemplate() func(api.Router) {
	return func(r api.Router) {
		controller := notification.NewEmailController()
		controller.Route(r)
	}
}

func asChiRouter(fn func(api.Router)) func(chi.Router) {
	return func(r chi.Router) {
		fn(&ChiRouter{r})
	}
}
