package notification

import (
	"encoding/json"
	"log/slog"
	"net/http"

	disgo "github.com/disgoorg/disgo/discord"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/discord"
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
			found := e.Provider.Exists(body.TemplateID)

			if !found {
				return api.Err(api.BadRequestErr, "Template not found")
			}

			template, err := e.Provider.ParseHtmlTemplate(body.TemplateID, recipient.Data)

			if err != nil {
				slog.Error("Failed to parse template", err)
				return api.Err(api.BadRequestErr, "Error parsing template")
			}

			if recipient.Data["to"] == nil {
				return api.Err(api.BadRequestErr, "[to] field is required")
			}

			to, ok := recipient.Data["to"].(string)

			if !ok {
				return api.Err(api.BadRequestErr, "Recipient email is not a string")
			}

			notifications = append(notifications, types.Mail{
				To:      []string{to},
				Subject: body.Subject,
				Sender:  config.Hermes.SMTP.Sender,
				Body:    template.String(),
				Type:    types.HTML,
			})
		case types.DISCORD:
			apiKey := r.Header.Get("X-API-KEY")
			client, err := discord.Connect(apiKey)

			if err != nil {
				return api.Err(api.BadRequestErr, "Failed to connect to Discord")
			}

			embed := disgo.NewEmbedBuilder().SetTitle(body.Subject)

			for k, v := range recipient.Data {
				embed.AddField(k, v.(string), false)
			}

			discord.SendWebhook(client, embed.Build())
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
