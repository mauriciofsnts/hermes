package notification

import (
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailController struct {
	Provider template.TemplateProvider
	Queue    types.Queue[types.Mail]
}

func NewEmailController(template template.TemplateProvider, queue types.Queue[types.Mail]) *EmailController {
	return &EmailController{
		Provider: template,
		Queue:    queue,
	}
}

var _ api.Controller = &EmailController{}

func (e *EmailController) Route(r api.Router) {
	r.Post("/notification", e.Notify)
}
