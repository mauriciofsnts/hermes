package notify

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/helper"
	"github.com/mauriciofsnts/hermes/internal/server/validator"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailControllerInterface interface {
	SendPlainTextEmail(w http.ResponseWriter, r *http.Request)
	SendTemplateEmail(w http.ResponseWriter, r *http.Request)
}

type EmailControler struct {
	EmailControllerInterface
	templateProvider template.TemplateProvider
	queue            types.Queue[types.Mail]
}

func NewEmailController() *EmailControler {
	return &EmailControler{
		templateProvider: template.NewTemplateService(),
		queue:            queue.Queue,
	}
}

func (e *EmailControler) SendPlainTextEmail(w http.ResponseWriter, r *http.Request) {
	queue := e.queue

	if queue == nil {
		slog.Error("Queue not found")
		helper.Err(w, helper.BadRequestErr, "Queue not found")
		return
	}

	body, validationErr := validator.MustGetBody[types.PlainTextEmail](r)

	if validationErr != nil {
		helper.DetailedError(w, helper.BadRequestErr, validationErr.Details)
		return
	}

	mail := types.Mail{
		To:      []string{body.To},
		Subject: body.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    body.Body,
		Type:    types.TEXT,
	}

	queue.Write(mail)

	helper.Created(w, "Email sent successfully")
}

func (e *EmailControler) SendTemplateEmail(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "slug")

	if templateName == "" {
		helper.Err(w, helper.BadRequestErr, "Invalid template name")
		return
	}

	queue := e.queue

	if queue == nil {
		slog.Error("Queue not found")
		helper.Err(w, helper.BadRequestErr, "Queue not found")
		return
	}

	body, validationErr := validator.MustGetBody[types.TemplateEmail](r)

	if validationErr != nil {
		helper.Err(w, helper.BadRequestErr, "Invalid body")
		return
	}

	template, err := e.templateProvider.ParseTemplate(templateName, body.Data)

	if err != nil {
		helper.Err(w, helper.BadRequestErr, "Error parsing template")
		return
	}

	mail := types.Mail{
		To:      []string{body.To},
		Subject: body.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    template.String(),
		Type:    types.HTML,
	}

	queue.Write(mail)

	helper.Created(w, "Email sent successfully")
}
