package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mauriciofsnts/hermes/internal/config"

	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	err := queue.NewQueue(cfg)

	if err != nil {
		slog.Error("Failed to create queue: " + err.Error())
		os.Exit(1)
	}

	slog.Info("Connecting to SMTP server...")
	err = smtp.Ping()

	for i := 0; i < 2 && err != nil; i++ {
		slog.Warn("Failed to connect to SMTP server, retrying... Error: ", err)
		err = smtp.Ping()

		if i == 1 && err != nil {
			slog.Error("Failed to connect to SMTP server: %v", err)
			os.Exit(0)
		}
	}

	go queue.StartWorker()
	go onShutdown()

	server.StartServer()

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
		os.Exit(0)
	}
}

func onShutdown() {
	stop := make(chan os.Signal, 1)

	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	queue.StopWorker()
	os.Exit(0)
}
