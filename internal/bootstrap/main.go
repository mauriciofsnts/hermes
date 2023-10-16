package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/http"
	"github.com/mauriciofsnts/hermes/internal/worker"
	"github.com/pauloo27/logger"
)

func Start() {

	logger.Debug("Starting Hermes...")
	logger.HandleFatal(config.LoadConfig(), "Failed to load config")

	go worker.StartWorker()

	app := http.CreateFiberInstance()

	go onShutdown(app)
	logger.HandleFatal(http.Listen(app), "Failed to start HTTP server")

}

func onShutdown(app *fiber.App) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	// worker.StopWorker()
	app.Shutdown()
	os.Exit(0)
}
