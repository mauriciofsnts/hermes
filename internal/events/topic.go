package events

import (
	"net"

	"github.com/Pauloo27/logger"
	"github.com/segmentio/kafka-go"
)

const (
	EmailTopic = "kafka-email-topic"
)

var defaultTopics = []string{
	EmailTopic,
}

func CreateTopic() error {
	connection, err := kafka.Dial("tcp", net.JoinHostPort("localhost", "9092"))

	if err != nil {
		logger.Error("Failed to connect to Kafka", err)
		return err
	}

	topics := make([]kafka.TopicConfig, len(defaultTopics))

	for i, topic := range defaultTopics {
		topics[i] = NewTopic(topic)
	}

	err = connection.CreateTopics(topics...)

	if err != nil {
		logger.Error("Failed to create topic", err)
		return err
	}

	return nil
}

func NewTopic(topicName string) kafka.TopicConfig {
	return kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
