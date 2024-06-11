package redis

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
}

type RedisQueue[T any] struct {
	client   *redis.Client
	consumer *Consumer[types.Mail]
}

func (r *RedisQueue[T]) Read(ctx context.Context) {
	slog.Debug("Starting Redis consumer...")

	r.consumer = NewConsumer[types.Mail](r.client, "hermes")

	readCh := make(chan ReadData[types.Mail])

	go r.consumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
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
			continue
		}
	}

}

func (r *RedisQueue[T]) Write(email types.Mail) error {
	producer := NewProducer[types.Mail](
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

func NewRedisProvider() (worker.Queue[types.Mail], error) {
	client := NewRedisClient(config.Hermes.Redis.Address, config.Hermes.Redis.Password)

	_, err := client.Ping(ctx).Result()

	if err != nil {
		slog.Error("Failed to connect to redis", err)
		return nil, err
	}

	return &RedisQueue[types.Mail]{
		client:   client,
		consumer: NewConsumer[types.Mail](client, config.Hermes.Redis.Topic),
	}, nil
}
