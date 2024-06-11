package memory

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type MemoryQueue[T any] struct {
	email chan types.Mail
}

func (m *MemoryQueue[T]) Read(ctx context.Context) {
	slog.Debug("Reading emails from memory")

	for {
		select {
		case <-ctx.Done():
			// TODO! graceful shutdown
			slog.Debug("Context done, stopping read emails from memory")
			return
		case email := <-m.email:
			slog.Debug("Sending email..")
			err := smtp.SendEmail(&email)

			// TODO! error handling?
			if err != nil {
				slog.Error("Error sending email", err)
			}
		}
	}

}

func (m *MemoryQueue[T]) Write(email types.Mail) error {
	slog.Debug("Writing email to memory")
	m.email <- email
	return nil
}

func (m *MemoryQueue[T]) Ping() (string, error) {
	return "Memory queue is up", nil
}

func NewMemoryProvider() worker.Queue[types.Mail] {
	return &MemoryQueue[types.Mail]{
		email: make(chan types.Mail, 10),
	}
}
