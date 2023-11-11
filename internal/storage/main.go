package storage

import (
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/storage/kafka"
	"github.com/mauriciofsnts/hermes/internal/storage/memory"
	"github.com/mauriciofsnts/hermes/internal/storage/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var storage types.Storage[types.Email]

func NewStorage() types.Storage[types.Email] {
	kafkaEnabled := config.Hermes.Kafka.Enabled
	redisEnabled := config.Hermes.Redis.Enabled

	if kafkaEnabled {
		err := kafka.CreateTopic()

		if err != nil {
			panic(err)
		}

		storage = kafka.NewKafkaStorage()
		return storage
	} else if redisEnabled {
		storage = redis.NewRedisStorage()
		return storage
	} else {
		storage = memory.NewMemoryStorage()
		return storage
	}
}
