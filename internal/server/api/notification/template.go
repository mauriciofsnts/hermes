package notification

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/server/validator"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type CreateTemplateNotificationBody struct {
	To      string         `json:"to" validate:"required"`
	Subject string         `json:"subject" validate:"required"`
	Data    map[string]any `json:"data"`
}

func (e *EmailController) HtmlTemplateNotification(r *http.Request) api.Response {
	templateName := chi.URLParam(r, "slug")

	if templateName == "" {
		return api.Err(api.BadRequestErr, "Invalid template name")
	}

	queue := e.queue

	if queue == nil {
		return api.Err(api.BadRequestErr, "Queue not found")
	}

	body, validationErr := validator.MustGetBody[CreateTemplateNotificationBody](r)

	if validationErr != nil {
		return api.DetailedError(validationErr.Error, validationErr.Details)
	}

	template, err := e.provider.ParseTemplate(templateName, body.Data)

	if err != nil {
		return api.Err(api.BadRequestErr, "Error parsing template")
	}

	notification := types.Mail{
		To:      []string{body.To},
		Subject: body.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    template.String(),
		Type:    types.HTML,
	}

	err = queue.Write(notification)

	if err != nil {
		return api.Err(api.InternalServerErr, "Failed to send notification")
	}

	return api.Created("Notification registered successfully")
}
