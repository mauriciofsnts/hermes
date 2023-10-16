package kafka

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
	kafkaGo "github.com/segmentio/kafka-go"
)

type KafkaStorage[T any] struct{}

func (k *KafkaStorage[T]) Read() {
	logger.Info("Starting Kafka consumer...")

	var emailConsumer *Consumer[types.Email]

	dialer := &kafkaGo.Dialer{Timeout: 10 * time.Second, DualStack: true}

	emailConsumer = NewConsumer[types.Email](dialer, EmailTopic)

	emailConsumer.Read(func(email *types.Email, err error) {
		if err != nil {
			logger.Error("Failed to read email", err)
			return
		}

		// TODO: Send email
		logger.Info("Email received", email)
	})
}

func (k *KafkaStorage[T]) Write(email types.Email) error {
	var producer = NewProducer[types.Email]()

	err := producer.Produce(uuid.New().String(), email, EmailTopic)

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
