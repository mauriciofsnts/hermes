package kafka

import (
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/config"
	kafkaGo "github.com/segmentio/kafka-go"
)

func CreateTopic() error {
	var defaultTopics = []string{}

	connection, err := kafkaGo.Dial("tcp", config.Hermes.Kafka.Address)

	if err != nil {
		slog.Error("Failed to connect to Kafka", err)
		return err
	}

	if config.Hermes.Kafka.Topic != "" {
		defaultTopics = append(defaultTopics, config.Hermes.Kafka.Topic)
	}

	topics := make([]kafkaGo.TopicConfig, len(defaultTopics))

	for i, topic := range defaultTopics {
		topics[i] = NewTopic(topic)
	}

	err = connection.CreateTopics(topics...)

	if err != nil {
		slog.Error("Failed to create topic", err)
		return err
	}

	return nil
}

func NewTopic(topicName string) kafkaGo.TopicConfig {
	return kafkaGo.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
