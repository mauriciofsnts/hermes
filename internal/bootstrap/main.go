package bootstrap

import (
	"github.com/Pauloo27/logger"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/http"
)

func Start() {
	logger.Debug("Starting Hermes...")
	logger.HandleFatal(config.LoadConfig(), "Failed to load config")
	logger.HandleFatal(http.Listen(), "Failed to start HTTP server")
}
