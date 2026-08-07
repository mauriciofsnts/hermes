package notification

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

// Notify godoc
//
//	@Summary		Send notification
//	@Description	Send a notification through multiple channels (email, Discord, etc.)
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			X-API-Key	header		string						true	"API Key"
//	@Param			request		body		types.NotificationRequest	true	"Notification request"
//	@Success		201			{object}	map[string]interface{}		"Notification registered successfully"
//	@Failure		400			{object}	map[string]interface{}		"Invalid request"
//	@Failure		500			{object}	map[string]interface{}		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/api/v1/app/notify/notification [post]
func (e *EmailController) Notify(r *http.Request) api.Response {
	queue := e.Queue
	if queue == nil {
		slog.Error("queue is not running or not found")
		return api.Err(api.InternalServerErr, "failed to send notification, contact administrator")
	}

	var body types.NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return api.Err(api.BadRequestErr, "failed to parse request body")
	}

	apiKey := r.Header.Get("x-api-key")
	client := config.Hermes.AppsByAPIKey[apiKey]
	if client == nil {
		return api.Err(api.UnauthorizedErr, "invalid API key")
	}

	notifications := make([]types.Mail, 0, len(body.Recipients))

	for _, recipient := range body.Recipients {
		switch recipient.Type {
		case types.MAIL:
			if !slices.Contains(client.EnabledFeatures, config.FeatEmail) {
				return api.Err(api.BadRequestErr, "email feature is not enabled")
			}

			notification, err := e.ValidateEmailNotification(body.TemplateID, recipient.Data, body.Subject)
			if err != nil {
				return api.Err(api.BadRequestErr, err.Error())
			}
			notifications = append(notifications, *notification)

		case types.DISCORD:
			if !slices.Contains(client.EnabledFeatures, config.FeatDiscord) {
				return api.Err(api.BadRequestErr, "discord feature is not enabled")
			}

			if err := e.SendDiscordNotification(apiKey, recipient.Data, body.Subject); err != nil {
				return api.Err(api.BadRequestErr, err.Error())
			}

		default:
			return api.Err(api.BadRequestErr, "recipient type not found")
		}
	}

	for _, notification := range notifications {
		if err := queue.Write(notification); err != nil {
			slog.Error("failed to write notification to queue", "error", err)
			return api.Err(api.InternalServerErr, "failed to enqueue notification")
		}
	}

	return api.Created("Notification registered successfully")
}
