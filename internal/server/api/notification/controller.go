package notification

import (
	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailController struct {
	provider template.TemplateProvider
	queue    types.Queue[types.Mail]
}

func NewEmailController() *EmailController {
	return &EmailController{
		provider: template.NewTemplateService(),
		queue:    queue.Queue,
	}
}

var _ api.Controller = &EmailController{}

func (e *EmailController) Route(r api.Router) {
	r.Post("/email/plain", e.PlainTextNotification)
	r.Post("/email/template", e.HtmlTemplateNotification)
}
