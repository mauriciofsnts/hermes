package discord

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
)

var (
	discordClients = make(map[string]webhook.Client)
	clientsMu      sync.RWMutex
)

func Connect(key string) (webhook.Client, error) {
	app, ok := config.Hermes.AppsByAPIKey[key]
	if !ok || app == nil || app.Discord == nil {
		return nil, errors.New("client has no discord configuration")
	}

	webhookID := app.Discord.ID
	webhookToken := app.Discord.Token

	if webhookID == "" || webhookToken == "" {
		slog.Error("discord webhook ID or token not found")
		return nil, errors.New("discord webhook id or token not found")
	}

	clientsMu.RLock()
	if client, ok := discordClients[key]; ok {
		clientsMu.RUnlock()
		return client, nil
	}
	clientsMu.RUnlock()

	id, err := snowflake.Parse(webhookID)
	if err != nil {
		slog.Error("error parsing snowflake ID", "error", err)
		return nil, err
	}

	client := webhook.New(snowflake.ID(id), webhookToken)

	clientsMu.Lock()
	discordClients[key] = client
	clientsMu.Unlock()

	return client, nil
}

func SendWebhook(client webhook.Client, embed discord.Embed) error {
	message, err := client.CreateEmbeds([]discord.Embed{embed})
	if err != nil {
		slog.Error("failed to send webhook", "error", err)
		return err
	}

	slog.Debug("webhook sent", "message", message)
	return nil
}
