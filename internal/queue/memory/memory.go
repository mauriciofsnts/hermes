package memory

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type MemoryQueue[T any] struct {
	email chan types.Email
}

func (m *MemoryQueue[T]) Read(ctx context.Context) {
	slog.Info("Reading emails from memory")

	for {
		slog.Info("Waiting for emails")
		select {
		case <-ctx.Done():
			// TODO! graceful shutdown
			slog.Info("Context done, stopping read emails from memory")
			return
		case email := <-m.email:
			err := smtp.SendEmail(&email)

			// TODO! error handling?
			if err != nil {
				slog.Error("Error sending email", err)
			}
		}
	}

}

func (m *MemoryQueue[T]) Write(email types.Email) error {
	m.email <- email
	return nil
}

func (m *MemoryQueue[T]) Ping() (string, error) {
	return "Memory queue is up", nil
}

func NewMemoryQueue() types.Queue[types.Email] {
	return &MemoryQueue[types.Email]{
		email: make(chan types.Email, 10),
	}
}
