package kafka

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	kafkaGo "github.com/segmentio/kafka-go"
)

type KakfaQueue[T any] struct {
	producer *Producer[T]
	dialer   *kafkaGo.Dialer
	consumer *Consumer[types.Email]
}

func (k *KakfaQueue[T]) Read(ctx context.Context) {
	slog.Info("Starting Kafka consumer...")

	readCh := make(chan types.ReadData[types.Email])

	go k.consumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping Kafka consumer...")

			k.consumer.Close()
			return
		case data := <-readCh:
			if data.Err != nil {
				slog.Error("Failed to read content", data.Err)
				continue
			}

			smtp.SendEmail(data.Data)
		}
	}

}

func (k *KakfaQueue[T]) Write(content T) error {
	err := k.producer.Produce(uuid.New().String(), content, config.Hermes.Kafka.Topic)

	if err != nil {
		slog.Error("Failed to produce content", err)
		return err
	}

	return nil
}

func (k *KakfaQueue[T]) Ping() (string, error) {
	conn, err := k.dialer.DialLeader(context.Background(), "tcp", config.Hermes.Kafka.Host, config.Hermes.Kafka.Topic, 0)

	if err != nil {
		slog.Error("Failed to connect to Kafka", err)
		return "", err
	}

	conn.Close()
	return "Kafka is up", nil
}

func NewKafkaQueue() types.Queue[types.Email] {
	dialer := &kafkaGo.Dialer{Timeout: 10 * time.Second, DualStack: true}

	return &KakfaQueue[types.Email]{
		producer: NewProducer[types.Email](),
		dialer:   dialer,
		consumer: NewConsumer[types.Email](dialer, config.Hermes.Kafka.Topic),
	}
}
