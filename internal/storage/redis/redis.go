package redis

import (
	"context"
	"fmt"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Hermes.Redis.Host, config.Hermes.Redis.Port),
		Password: config.Hermes.Redis.Password,
		DB:       0,
	})
}

type RedisStorage[T any] struct{}

func (r *RedisStorage[T]) Read(ctx context.Context) {
	logger.Info("Starting Redis consumer...")

	var emailConsumer *Consumer[types.Email]

	client := NewRedisClient()

	emailConsumer = NewConsumer[types.Email](client, config.Hermes.Redis.Topic)

	readCh := make(chan ReadData[types.Email])

	go emailConsumer.Read(readCh)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping redis consumer...")
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

func (r *RedisStorage[T]) Write(email types.Email) error {

	client := NewRedisClient()

	producer := NewProducer[types.Email](
		*client,
		config.Hermes.Redis.Topic,
	)

	err := producer.Produce(email)

	if err != nil {
		logger.Error("Failed to produce email", err)
		return err
	}

	// TODO: Send email
	logger.Info("Email produced", email)
	return nil
}

func NewRedisStorage() types.Storage[types.Email] {
	return &RedisStorage[types.Email]{}
}
