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
	producer *Producer[types.Mail]
	dlq      *database.DLQService
}

func (r *RedisQueue[T]) Read(ctx context.Context) {
	slog.Debug("starting Redis consumer...")

	consumer := NewConsumer[types.Mail](r.client, config.Hermes.Redis.Topic)
	readCh := make(chan ReadData[types.Mail])

	go consumer.Read(ctx, readCh)

	for {
		select {
		case <-ctx.Done():
			if err := consumer.Close(); err != nil {
				slog.Error("failed to close redis consumer", "error", err)
			}
			return
		case data := <-readCh:
			if data.Err != nil {
				slog.Error("failed to read email", "error", data.Err)
				continue
			}

			if err := smtp.SendEmail(data.Data); err != nil {
				slog.Error("failed to send email", "error", err)
				r.storeInDLQ(data.Data, err)
				continue
			}
		}
	}
}

func (r *RedisQueue[T]) storeInDLQ(email *types.Mail, sendErr error) {
	if r.dlq == nil {
		return
	}

	emailJSON, jsonErr := json.Marshal(email)
	if jsonErr != nil {
		slog.Error("failed to marshal email for DLQ", "error", jsonErr)
		return
	}

	if dlqErr := r.dlq.Store(string(emailJSON), sendErr.Error(), "unknown"); dlqErr != nil {
		slog.Error("failed to store email in DLQ", "error", dlqErr)
	} else {
		slog.Info("email stored in DLQ for retry")
	}
}

func (r *RedisQueue[T]) Write(email types.Mail) error {
	if err := r.producer.Produce(email); err != nil {
		slog.Error("failed to produce email", "error", err)
		return err
	}
	return nil
}

func (r *RedisQueue[T]) Ping() (string, error) {
	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("failed to ping redis", "error", err)
		return "", err
	}
	return "Redis is up", nil
}

func NewRedisProvider() (worker.Queue[types.Mail], error) {
	client := NewRedisClient(config.Hermes.Redis.Address, config.Hermes.Redis.Password)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return nil, err
	}

	var dlqService *database.DLQService
	dlqService, err = database.NewDLQService("hermes_dlq.db")
	if err != nil {
		slog.Warn("failed to initialize DLQ service, continuing without DLQ", "error", err)
		dlqService = nil
	} else {
		slog.Info("DLQ service initialized successfully")
	}

	return &RedisQueue[types.Mail]{
		client:   client,
		producer: NewProducer[types.Mail](*client, config.Hermes.Redis.Topic),
		dlq:      dlqService,
	}, nil
}
