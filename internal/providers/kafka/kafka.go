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
	consumer *Consumer[types.Mail]
}

func (k *KakfaQueue[T]) Read(ctx context.Context) {
	slog.Info("Starting Kafka consumer...")

	readCh := make(chan types.ReadData[types.Mail])

	go k.consumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping Kafka consumer...")

			_ = k.consumer.Close()
			return
		case data := <-readCh:
			if data.Err != nil {
				slog.Error("Failed to read content", data.Err)
				continue
			}

			err := smtp.SendEmail(data.Data)

			if err != nil {
				slog.Error("Failed to send email", err)
				continue
			}
		}
	}

}

func (k *KakfaQueue[T]) Write(content T) error {
	err := k.producer.Produce(uuid.New().String(), content, config.Envs.Kafka.Topic)

	if err != nil {
		slog.Error("Failed to produce content", err)
		return err
	}

	return nil
}

func (k *KakfaQueue[T]) Ping() (string, error) {
	conn, err := k.dialer.DialLeader(context.Background(), "tcp", config.Envs.Kafka.Host, config.Envs.Kafka.Topic, 0)

	if err != nil {
		slog.Error("Failed to connect to Kafka", err)
		return "", err
	}

	_ = conn.Close()
	return "Kafka is up", nil
}

func NewKafkaProvider() types.Queue[types.Mail] {
	dialer := &kafkaGo.Dialer{Timeout: 10 * time.Second, DualStack: true}

	return &KakfaQueue[types.Mail]{
		producer: NewProducer[types.Mail](),
		dialer:   dialer,
		consumer: NewConsumer[types.Mail](dialer, config.Envs.Kafka.Topic),
	}
}
