package template

import (
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

type TemplateController struct {
	provider template.TemplateProvider
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		provider: template.NewTemplateService(),
	}
}

var _ api.Controller = &TemplateController{}

func (c *TemplateController) Route(r api.Router) {
	r.Post("/template", c.CreateTemplate)
	r.Get("/template/{slug}", c.GetRaw)
}
