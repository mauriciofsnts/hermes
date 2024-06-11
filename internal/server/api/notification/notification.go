package notification

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

func (e *EmailController) Notify(r *http.Request) api.Response {

	queue := e.Queue

	if queue == nil {
		slog.Error("Queue is not running or not found")
		return api.Err(api.InternalServerErr, "Failed to send notification, contact administrator")
	}

	var body types.NotificationRequest

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return api.Err(api.BadRequestErr, "Failed to parse request body")
	}

	notifications := make([]types.Mail, 0)

	for _, recipient := range body.Recipients {
		switch recipient.Type {
		case types.MAIL:
			notification, err := e.ValidateEmailNotification(body.TemplateID, recipient.Data, body.Subject)
			if err != nil {
				return api.Err(api.BadRequestErr, err.Error())
			} else {
				notifications = append(notifications, *notification)
			}
		case types.DISCORD:
			apiKey := r.Header.Get("X-API-KEY")
			err := e.ValidateDiscordNotification(apiKey, recipient.Data, body.Subject)
			if err != nil {
				return api.Err(api.BadRequestErr, err.Error())
			}
		default:
			return api.Err(api.BadRequestErr, "Recipient type not found")
		}

	}

	for _, notification := range notifications {
		err = queue.Write(notification)
		// save on db to try again later
		if err != nil {
			slog.Error("Failed to write notification to queue", err)
		}
	}

	return api.Created("Notification registered successfully")
}
