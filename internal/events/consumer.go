package events

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Consumer[T any] struct {
	reader *kafka.Reader
	Dialer *kafka.Dialer
	Topic  string
}

func NewConsumer[T any](dialer *kafka.Dialer, topic string) *Consumer[T] {
	return &Consumer[T]{
		Dialer: dialer,
		Topic:  topic,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			Dialer:  dialer,
			GroupID: "hermes",
		}),
	}
}

// TODO: Move to a channel
func (c *Consumer[T]) Read(callback func(*T, error)) {
	for {
		message, err := c.reader.ReadMessage(context.Background())

		if err != nil {
			callback(nil, err)
			continue
		}

		var model T

		err = json.Unmarshal(message.Value, &model)

		if err != nil {
			callback(nil, err)
			continue
		}

		callback(&model, nil)
	}
}

func (c *Consumer[T]) Close() error {
	return c.reader.Close()
}
