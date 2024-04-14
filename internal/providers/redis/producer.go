package redis

import (
	"context"
	"encoding/json"

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
		return err
	}

	pubsub := p.Client.Subscribe(ctx, p.Topic)

	_, err = pubsub.Receive(ctx)

	if err != nil {
		return err
	}

	err = p.Client.Publish(ctx, p.Topic, data).Err()

	if err != nil {
		return err
	}

	return nil
}
