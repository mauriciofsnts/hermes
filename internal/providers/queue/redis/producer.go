package redis

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Producer[T any] struct {
	Client redis.Client
	Topic  string
}

func NewProducer[T any](client redis.Client, topic string) *Producer[T] {
	return &Producer[T]{
		Client: client,
		Topic:  topic,
	}
}

func (p *Producer[T]) Produce(value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		slog.Error("Failed to marshal value for Redis", "topic", p.Topic, "error", err)
		return err
	}

	result := p.Client.Publish(ctx, p.Topic, data)
	err = result.Err()
	if err != nil {
		slog.Error("Failed to publish to Redis", "topic", p.Topic, "error", err)
		return err
	}

	subscribers := result.Val()
	if subscribers > 0 {
		slog.Debug("Message published to Redis", "topic", p.Topic, "subscribers", subscribers)
	} else {
		slog.Warn("Message published but no subscribers listening", "topic", p.Topic)
	}

	return nil
}
