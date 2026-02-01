package redis

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/database"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

type RedisQueue[T any] struct {
	client   *redis.Client
	consumer *Consumer[types.Mail]
	dlq      *database.DLQService
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
				slog.Error("Failed to read email", "error", data.Err)
				continue
			}

			err := smtp.SendEmail(data.Data)

			if err != nil {
				slog.Error("Failed to send email", "error", err)

				// Store in DLQ if available
				if r.dlq != nil {
					emailJSON, jsonErr := json.Marshal(data.Data)
					if jsonErr == nil {
						dlqErr := r.dlq.Store(string(emailJSON), err.Error(), "unknown")
						if dlqErr != nil {
							slog.Error("Failed to store email in DLQ", "error", dlqErr)
						} else {
							slog.Info("Email stored in DLQ for retry")
						}
					}
				}
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
		slog.Error("Failed to produce email", "error", err)
		return err
	}

	return nil
}

func (r *RedisQueue[T]) Ping() (string, error) {
	_, err := r.client.Ping(context.Background()).Result()

	if err != nil {
		slog.Error("Failed to ping redis", "error", err)
		return "", err
	}

	return "Redis is up", nil
}

func NewRedisProvider() (worker.Queue[types.Mail], error) {
	client := NewRedisClient(config.Hermes.Redis.Address, config.Hermes.Redis.Password)

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		return nil, err
	}

	// Initialize DLQ service (optional)
	var dlqService *database.DLQService
	dlqService, err = database.NewDLQService("hermes_dlq.db")
	if err != nil {
		slog.Warn("Failed to initialize DLQ service, continuing without DLQ", "error", err)
		dlqService = nil
	} else {
		slog.Info("DLQ service initialized successfully")
	}

	return &RedisQueue[types.Mail]{
		client:   client,
		consumer: NewConsumer[types.Mail](client, config.Hermes.Redis.Topic),
		dlq:      dlqService,
	}, nil
}
