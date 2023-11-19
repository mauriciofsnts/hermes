package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api/router"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue"
)

func Start() {
	err := config.LoadConfig()

	if err != nil {
		slog.Error("Failed to load envs: " + err.Error())
	}

	SetupLog()

	q := queue.NewQueue()

	go queue.StartWorker(q)
	app := router.CreateFiberInstance(q)

	go onShutdown(app)

	err = router.Listen(app)

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
	}
}

func onShutdown(app *fiber.App) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	queue.StopWorker()
	app.Shutdown()
	os.Exit(0)
}
