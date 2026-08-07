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
	email chan T
}

func (m *MemoryQueue[T]) Read(ctx context.Context) {
	slog.Debug("starting memory queue workers", "worker_count", WorkerPoolSize)

	for i := 0; i < WorkerPoolSize; i++ {
		go m.worker(ctx, i)
	}

	<-ctx.Done()
	slog.Debug("context done, stopping memory queue workers")
}

func (m *MemoryQueue[T]) worker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("worker stopping", "worker_id", workerID)
			return
		case email := <-m.email:
			mail, ok := any(email).(types.Mail)
			if !ok {
				slog.Error("invalid item type in memory queue", "worker_id", workerID)
				continue
			}

			slog.Debug("worker processing email", "worker_id", workerID, "to", mail.To)
			if err := smtp.SendEmail(&mail); err != nil {
				slog.Error("error sending email", "worker_id", workerID, "to", mail.To, "error", err)
			} else {
				slog.Info("email sent successfully", "worker_id", workerID, "to", mail.To)
			}
		}
	}
}

func (m *MemoryQueue[T]) Write(email types.Mail) error {
	slog.Debug("writing email to memory queue", "to", email.To)
	m.email <- any(email).(T)
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
