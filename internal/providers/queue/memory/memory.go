package memory

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
)

const WorkerPoolSize = 5

type MemoryQueue[T any] struct {
	email chan types.Mail
}

func (m *MemoryQueue[T]) Read(ctx context.Context) {
	slog.Debug("Starting memory queue workers", "worker_count", WorkerPoolSize)

	for i := 0; i < WorkerPoolSize; i++ {
		go m.worker(ctx, i)
	}

	<-ctx.Done()
	slog.Debug("Context done, stopping memory queue workers")
}

func (m *MemoryQueue[T]) worker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Worker stopping", "worker_id", workerID)
			return
		case email := <-m.email:
			slog.Debug("Worker processing email", "worker_id", workerID, "to", email.To)
			err := smtp.SendEmail(&email)

			if err != nil {
				slog.Error("Error sending email", "worker_id", workerID, "to", email.To, "error", err)
			} else {
				slog.Info("Email sent successfully", "worker_id", workerID, "to", email.To)
			}
		}
	}
}

func (m *MemoryQueue[T]) Write(email types.Mail) error {
	slog.Debug("Writing email to memory queue", "to", email.To)
	m.email <- email
	return nil
}

func (m *MemoryQueue[T]) Ping() (string, error) {
	return "Memory queue is up", nil
}

func NewMemoryProvider() worker.Queue[types.Mail] {
	return &MemoryQueue[types.Mail]{
		email: make(chan types.Mail, WorkerPoolSize*2),
	}
}
