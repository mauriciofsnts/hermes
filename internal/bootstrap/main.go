package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/ctx"

	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	q, err := queue.NewQueue(cfg)

	if err != nil {
		slog.Error("Failed to create queue: " + err.Error())
		os.Exit(1)
	}

	slog.Debug("Connecting to SMTP server...")
	err = smtp.Ping()

	for i := 0; i < 2 && err != nil; i++ {
		slog.Warn("Failed to connect to SMTP server, retrying... Error: ", err)
		err = smtp.Ping()

		if i == 1 && err != nil {
			slog.Error("Failed to connect to SMTP server: %v", err)
			os.Exit(0)
		}
	}

	providers := &ctx.Providers{
		Config: cfg,
		Queue:  q,
	}

	app := server.CreateFiberInstance(providers)

	go queue.StartWorker(q)
	go onShutdown(app)

	err = server.Listen(app)

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
		os.Exit(0)
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
