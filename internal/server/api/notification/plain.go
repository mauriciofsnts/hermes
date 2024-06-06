package notification

import (
	"log/slog"
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/server/validator"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type CreatePlainNotificationBody struct {
	To      string `json:"to" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body"`
}

func (e *EmailController) PlainTextNotification(r *http.Request) api.Response {
	queue := e.queue

	if queue == nil {
		slog.Error("Queue is not running or not found")
		return api.Err(api.InternalServerErr, "Failed to send notification, contact administrator")
	}

	body, validationErr := validator.MustGetBody[CreatePlainNotificationBody](r)

	if validationErr != nil {
		return api.DetailedError(validationErr.Error, validationErr.Details)
	}

	notification := types.Mail{
		To:      []string{body.To},
		Subject: body.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    body.Body,
		Type:    types.TEXT,
	}

	err := queue.Write(notification)

	if err != nil {
		slog.Error("Failed to send notification", err)
		return api.Err(api.InternalServerErr, "Failed to register notification")
	}

	return api.Created("Notification registered successfully")
}
