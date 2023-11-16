package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api/router"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue"
	"github.com/pauloo27/logger"
)

func Start() {
	logger.Debug("Starting Hermes...")
	logger.HandleFatal(config.LoadConfig(), "Failed to load config")

	q := queue.NewQueue()

	go queue.StartWorker(q)
	app := router.CreateFiberInstance(q)

	go onShutdown(app)
	logger.HandleFatal(router.Listen(app), "Failed to start HTTP server")

}

func onShutdown(app *fiber.App) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	queue.StopWorker()
	app.Shutdown()
	os.Exit(0)
}
