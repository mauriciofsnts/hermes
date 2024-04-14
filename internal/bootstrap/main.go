package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api/router"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers"
)

func Start() {
	err := config.LoadConfig()

	if err != nil {
		slog.Error("Failed to load envs: " + err.Error())
	}

	SetupLog()

	queue := providers.NewQueue()

	go providers.StartWorker(queue)
	app := router.CreateFiberInstance(queue)

	go onShutdown(app)

	err = router.Listen(app)

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
	}
}

func onShutdown(app *fiber.App) {
	stop := make(chan os.Signal, 1)

	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	providers.StopWorker()
	_ = app.Shutdown()
	os.Exit(0)
}
