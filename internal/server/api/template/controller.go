package template

import (
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

type TemplateController struct {
	provider template.TemplateProvider
}

func NewTemplateController(storage template.TemplateProvider) *TemplateController {
	return &TemplateController{
		provider: storage,
	}
}

var _ api.Controller = &TemplateController{}

func (c *TemplateController) Route(r api.Router) {
	r.Post("/template", c.CreateTemplate)
	r.Get("/template/{slug}", c.GetRaw)
}
