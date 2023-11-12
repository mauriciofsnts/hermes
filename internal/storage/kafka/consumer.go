package kafka

import (
	"context"
	"encoding/json"

	"github.com/mauriciofsnts/hermes/internal/config"
	kafkaGo "github.com/segmentio/kafka-go"
)

type Consumer[T any] struct {
	reader *kafkaGo.Reader
	Dialer *kafkaGo.Dialer
	Topic  string
}

func NewConsumer[T any](dialer *kafkaGo.Dialer, topic string) *Consumer[T] {
	return &Consumer[T]{
		Dialer: dialer,
		Topic:  topic,
		reader: kafkaGo.NewReader(kafkaGo.ReaderConfig{
			Brokers: config.Hermes.Kafka.Brokers,
			Topic:   topic,
			Dialer:  dialer,
			GroupID: "hermes",
		}),
	}
}

type ReadData[T any] struct {
	Data *T
	Err  error
}

func (c *Consumer[T]) Read(ch chan<- ReadData[T]) {
	for {
		message, err := c.reader.ReadMessage(context.Background())

		if err != nil {
			ch <- ReadData[T]{nil, err}
			continue
		}

		var model T

		err = json.Unmarshal(message.Value, &model)

		if err != nil {
			ch <- ReadData[T]{nil, err}
			continue
		}

		ch <- ReadData[T]{&model, nil}
	}
}

func (c *Consumer[T]) Close() error {
	return c.reader.Close()
}
