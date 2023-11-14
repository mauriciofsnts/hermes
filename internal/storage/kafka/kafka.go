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

type KafkaStorage[T any] struct {
	producer *Producer[T]
	dialer   *kafkaGo.Dialer
	consumer *Consumer[types.Email]
}

func (k *KafkaStorage[T]) Read(ctx context.Context) {
	logger.Info("Starting Kafka consumer...")

	readCh := make(chan ReadData[types.Email])

	go k.consumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping Kafka consumer...")
			k.consumer.Close()
			return
		case data := <-readCh:
			if data.Err != nil {
				logger.Error("Failed to read content", data.Err)
				continue
			}

			smtp.SendEmail(data.Data)
		}
	}

}

func (k *KafkaStorage[T]) Write(content T) error {
	err := k.producer.Produce(uuid.New().String(), content, config.Hermes.Kafka.Topic)

	if err != nil {
		logger.Error("Failed to produce content", err)
		return err
	}

	return nil
}

func (k *KafkaStorage[T]) Ping() (string, error) {
	conn, err := k.dialer.DialLeader(context.Background(), "tcp", config.Hermes.Kafka.Host, config.Hermes.Kafka.Topic, 0)

	if err != nil {
		logger.Error("Failed to connect to Kafka", err)
		return "", err
	}

	conn.Close()
	return "Kafka is up", nil
}

func NewKafkaStorage() types.Storage[types.Email] {
	dialer := &kafkaGo.Dialer{Timeout: 10 * time.Second, DualStack: true}

	return &KafkaStorage[types.Email]{
		producer: NewProducer[types.Email](),
		dialer:   dialer,
		consumer: NewConsumer[types.Email](dialer, config.Hermes.Kafka.Topic),
	}
}
