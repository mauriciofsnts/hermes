package worker

import (
	"time"

	"github.com/Pauloo27/logger"
	"github.com/mauriciofsnts/hermes/internal/events"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/segmentio/kafka-go"
)

var emailConsumer *events.Consumer[types.Email]

func StartWorker() {
	logger.Debug("Starting worker...")

	dialer := &kafka.Dialer{Timeout: 10 * time.Second, DualStack: true}

	emailConsumer = events.NewConsumer[types.Email](dialer, events.EmailTopic)

	emailConsumer.Read(func(email *types.Email, err error) {
		if err != nil {

			// if errors.Is(err, io.EOF) {
			// 	return
			// }

			logger.Error("Failed to read email", err)
			return
		}

		logger.Info("Email received", email)
	})

}

func StopWorker() {
	logger.Debug("Stopping worker...")

	err := emailConsumer.Close()

	if err != nil {
		logger.Error("Failed to stop worker", err)
	}

}
