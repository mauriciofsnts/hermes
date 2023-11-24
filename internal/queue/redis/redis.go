package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Hermes.Redis.Host, config.Hermes.Redis.Port),
		Password: config.Hermes.Redis.Password,
	})
}

type RedisQueue[T any] struct {
	client   *redis.Client
	consumer *Consumer[types.Email]
}

func (r *RedisQueue[T]) Read(ctx context.Context) {
	slog.Info("Starting Redis consumer...")

	r.consumer = NewConsumer[types.Email](r.client, config.Hermes.Redis.Topic)

	readCh := make(chan types.ReadData[types.Email])

	go r.consumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping redis consumer...")
			_ = r.consumer.Close()
			return
		case data := <-readCh:
			if data.Err != nil {
				slog.Error("Failed to read email", data.Err)
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

func (r *RedisQueue[T]) Write(email types.Email) error {
	producer := NewProducer[types.Email](
		*r.client,
		config.Hermes.Redis.Topic,
	)

	err := producer.Produce(email)

	if err != nil {
		slog.Error("Failed to produce email", err)
		return err
	}

	return nil
}

func (r *RedisQueue[T]) Ping() (string, error) {
	_, err := r.client.Ping(ctx).Result()

	if err != nil {
		slog.Error("Failed to ping redis", err)
		return "", err
	}

	return "Redis is up", nil
}

func NewRedisQueue() types.Queue[types.Email] {
	client := NewRedisClient()

	return &RedisQueue[types.Email]{
		client:   client,
		consumer: NewConsumer[types.Email](client, config.Hermes.Redis.Topic),
	}
}
