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

		// TODO: avoid clone?
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}

		if res.StatusCode == 0 {
			slog.Error("Missing response status code", "path", r.URL.Path)
			res.StatusCode = http.StatusInternalServerError
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			_ = json.NewEncoder(w).Encode(res.Body)
		}
	}
}
