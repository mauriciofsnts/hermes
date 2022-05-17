package bootstrap

import (
	"github.com/Pauloo27/logger"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/discord"
)

func Start() {
	err := config.LoadConfig()
	logger.HandleFatal(err, "Error loading config")

	discord.Start()
}
