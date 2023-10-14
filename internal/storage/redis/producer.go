package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Producer[T any] struct {
	Client redis.Client
	Ctx    context.Context
	Topic  string
}

func NewProducer[T any](client redis.Client, ctx context.Context, topic string) *Producer[T] {
	return &Producer[T]{
		Client: client,
		Ctx:    ctx,
		Topic:  topic,
	}
}

func (p *Producer[T]) Produce(value T) error {

	pubsub := p.Client.Subscribe(p.Ctx, Topic)

	_, err := pubsub.Receive(p.Ctx)

	if err != nil {
		return err
	}

	// Publish a message.
	err = p.Client.Publish(p.Ctx, p.Topic, value).Err()

	if err != nil {
		return err
	}

	return nil
}
