package bootstrap

import (
	"github.com/Pauloo27/logger"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/discord"
)

func Start() {
	logger.HandleFatal(config.LoadConfig(), "Error loading config")
	logger.HandleFatal(discord.Start(), "Error starting discord")
}
