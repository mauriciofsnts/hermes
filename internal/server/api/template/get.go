package template

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

func (t *TemplateController) GetRaw(r *http.Request) api.Response {
	id := chi.URLParam(r, "slug")

	html, err := t.provider.Get(id)

	if err != nil {
		return api.DetailedError(api.NotFoundErr, "Template not found")
	}

	return api.OkHTML(string(html))
}
