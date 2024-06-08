package notification

import (
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailController struct {
	provider template.TemplateProvider
	queue    types.Queue[types.Mail]
}

func NewEmailController(template template.TemplateProvider, queue types.Queue[types.Mail]) *EmailController {
	return &EmailController{
		provider: template,
		queue:    queue,
	}
}

var _ api.Controller = &EmailController{}

func (e *EmailController) Route(r api.Router) {
	r.Post("/notification", e.Notify)
}
