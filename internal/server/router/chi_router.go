package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

type ChiRouter struct {
	chi chi.Router
}

var _ api.Router = &ChiRouter{}

func (c *ChiRouter) Get(path string, handler api.WrappedHandler) {
	c.chi.Get(path, c.wrap(handler))
}

func (c *ChiRouter) Post(path string, handler api.WrappedHandler) {
	c.chi.Post(path, c.wrap(handler))
}

func (c *ChiRouter) Put(path string, handler api.WrappedHandler) {
	c.chi.Put(path, c.wrap(handler))
}

func (c *ChiRouter) Delete(path string, handler api.WrappedHandler) {
	c.chi.Delete(path, c.wrap(handler))
}

func (c *ChiRouter) wrap(handler api.WrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := handler(r)

		if res.StatusCode == 0 {
			slog.Error("missing response status code", "path", r.URL.Path)
			res = api.Err(api.InternalServerErr, "missing response status code")
		}

		for k, values := range res.Header {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			if err := json.NewEncoder(w).Encode(res.Body); err != nil {
				slog.Error("failed to encode response body", "error", err)
			}
		}
	}
}
