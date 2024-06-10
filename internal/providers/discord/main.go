package discord

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
)

// map api key +discord.client
var discordClients = make(map[string]webhook.Client)

func Connect(key string) (webhook.Client, error) {

	if config.Hermes.AppsByAPIKey[key].Discord == nil {
		return nil, errors.New("client has no discord configuration")

	}

	webhookId := config.Hermes.AppsByAPIKey[key].Discord.ID
	webhookToken := config.Hermes.AppsByAPIKey[key].Discord.Token

	if webhookId == "" || webhookToken == "" {
		slog.Error("Discord webhook ID or token not found")
		return nil, errors.New("discord webhook id or token not found")
	}

	// if client already exists, return it
	if client, ok := discordClients[key]; ok {
		return client, nil
	}

	id, err := snowflake.Parse(webhookId)

	if err != nil {
		slog.Error("Error parsing snowflake ID", err)
		return nil, err
	}

	client := webhook.New(snowflake.ID(id), webhookToken)
	return client, nil

}

func SendWebhook(client webhook.Client, embed discord.Embed) {
	message, err := client.CreateEmbeds([]discord.Embed{embed})

	if err != nil {
		slog.Error("Failed to send webhook: ", err)
	}

	slog.Info("Webhook sent: ", message)
}
