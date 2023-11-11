package memory

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

type ReadData[T any] struct {
	Data *T
	Err  error
}

type MemoryStorage[T any] struct {
	email chan types.Email
}

func (m *MemoryStorage[T]) Read(ctx context.Context) {
	logger.Info("Reading emails from memory")

	for {
		logger.Info("Waiting for emails", m.email)
		select {
		case <-ctx.Done():
			// TODO! graceful shutdown
			logger.Info("Context done, stopping read emails from memory")
			return
		case email := <-m.email:
			err := smtp.SendEmail(&email)

			// TODO! error handling?
			if err != nil {
				logger.Error("Error sending email", err)
			}
		}
	}

}

func (m *MemoryStorage[T]) Write(email types.Email) error {
	m.email <- email
	return nil
}

func NewMemoryStorage() types.Storage[types.Email] {
	return &MemoryStorage[types.Email]{
		email: make(chan types.Email, 10),
	}
}
