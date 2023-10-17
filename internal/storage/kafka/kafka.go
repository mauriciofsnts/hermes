package kafka

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
	kafkaGo "github.com/segmentio/kafka-go"
)

type KafkaStorage[T any] struct{}

func (k *KafkaStorage[T]) Read(ctx context.Context) {
	logger.Info("Starting Kafka consumer...")

	var emailConsumer *Consumer[types.Email]

	dialer := &kafkaGo.Dialer{Timeout: 10 * time.Second, DualStack: true}

	emailConsumer = NewConsumer[types.Email](dialer, config.Hermes.Kafka.Topic)

	readCh := make(chan ReadData[types.Email])

	go emailConsumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping Kafka consumer...")
			emailConsumer.Close()
			return
		case data := <-readCh:
			if data.Err != nil {
				logger.Error("Failed to read email", data.Err)
				continue
			}

			smtp.SendEmail(data.Data)
		}
	}

}

func (k *KafkaStorage[T]) Write(email types.Email) error {
	var producer = NewProducer[types.Email]()

	err := producer.Produce(uuid.New().String(), email, config.Hermes.Kafka.Topic)

	if err != nil {
		logger.Error("Failed to produce email", err)
		return err
	}

	logger.Info("Email produced", email)
	return nil
}

func NewKafkaStorage() types.Storage[types.Email] {
	return &KafkaStorage[types.Email]{}
}
