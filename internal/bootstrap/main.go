package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	q := queue.NewQueue(cfg)

	go queue.StartWorker(q)

	app := server.CreateFiberInstance()

	go onShutdown(app)

	err := server.Listen(app)

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
		os.Exit(1)
	}
}

func onShutdown(app *fiber.App) {
	stop := make(chan os.Signal, 1)

	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	queue.StopWorker()
	_ = app.Shutdown()
	os.Exit(0)
}
