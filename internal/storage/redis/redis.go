package redis

import (
	"strconv"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Hermes.Redis.Host + ":" + strconv.Itoa(config.Hermes.Redis.Port),
		Password: config.Hermes.Redis.Password,
		DB:       0,
	})
}

type RedisStorage[T any] struct{}

func (r *RedisStorage[T]) Read() {
	logger.Info("Starting Redis consumer...")

	var emailConsumer *Consumer[types.Email]

	client := NewRedisClient()

	emailConsumer = NewConsumer[types.Email](client, Topic)

	emailConsumer.Read(func(email *types.Email, err error) {
		if err != nil {
			logger.Error("Failed to read email", err)
			return
		}

		logger.Info("Email received", email)
	})

}

func (r *RedisStorage[T]) Write(email types.Email) error {

	client := NewRedisClient()

	producer := NewProducer[types.Email](
		*client,
		Topic,
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
