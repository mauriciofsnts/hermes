package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers"

	q "github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server"
)

func Start(cfg *config.Config) {
	setupLog(cfg)

	queue, err := q.NewQueue(cfg)

	if err != nil {
		slog.Error("Failed to create queue: " + err.Error())
		os.Exit(1)
	}

	slog.Info("Connecting to SMTP server...")
	err = smtp.Ping()

	for i := 0; i < 2 && err != nil; i++ {
		slog.Warn("Failed to connect to SMTP server, retrying", "error", err)
		err = smtp.Ping()

		if i == 1 && err != nil {
			slog.Error("Failed to connect to SMTP server", "error", err)
			os.Exit(1)
		}
	}

	queueManager := q.NewQueueManager(queue)
	go onShutdown(queueManager)

	providers := &providers.Providers{
		// DB:      database.SetupConnection(),
		Queue:   queue,
		Storage: template.NewTemplateService(),
	}

	server.StartServer(providers)

	if err != nil {
		slog.Error("Failed to start HTTP server: " + err.Error())
		os.Exit(0)
	}
}

func onShutdown(queueManager *q.QueueManager) {
	stop := make(chan os.Signal, 1)

	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("Shutdown signal received, draining queue...")

	// Gracefully drain queue with 30 second timeout
	queueManager.DrainAndStop(30 * time.Second)

	slog.Info("Graceful shutdown complete")
	os.Exit(0)
}
