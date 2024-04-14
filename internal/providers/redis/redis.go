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
		Addr:     fmt.Sprintf("%s:%d", config.Envs.Redis.Host, config.Envs.Redis.Port),
		Password: config.Envs.Redis.Password,
	})
}

type RedisQueue[T any] struct {
	client   *redis.Client
	consumer *Consumer[types.Mail]
}

func (r *RedisQueue[T]) Read(ctx context.Context) {
	slog.Info("Starting Redis consumer...")

	r.consumer = NewConsumer[types.Mail](r.client, config.Envs.Redis.Topic)

	readCh := make(chan types.ReadData[types.Mail])

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

func (r *RedisQueue[T]) Write(email types.Mail) error {
	producer := NewProducer[types.Mail](
		*r.client,
		config.Envs.Redis.Topic,
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

func NewRedisProvider() types.Queue[types.Mail] {
	client := NewRedisClient()

	return &RedisQueue[types.Mail]{
		client:   client,
		consumer: NewConsumer[types.Mail](client, config.Envs.Redis.Topic),
	}
}
