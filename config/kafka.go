package config

import (
	"os"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists
	godotenv.Load()
}

// GetTopicName returns the Kafka topic name from environment variable or default
func GetTopicName() string {
	if topic := os.Getenv("KAFKA_TOPIC_NAME"); topic != "" {
		return topic
	}
	return "test-topic"
}

// GetBroker returns the Kafka broker address from environment variable or default
func GetBroker() string {
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		return broker
	}
	return "localhost:9092"
}

// GetGroupID returns the Kafka consumer group ID from environment variable or default
func GetGroupID() string {
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		return groupID
	}
	return "test-consumer-group"
}

// GetProducerConfig returns the Kafka producer configuration
func GetProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Version = sarama.V2_8_0_0

	// Add SASL authentication if credentials are provided
	if username := os.Getenv("KAFKA_USERNAME"); username != "" {
		if password := os.Getenv("KAFKA_PASSWORD"); password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			config.Net.SASL.User = username
			config.Net.SASL.Password = password
			config.Net.SASL.Handshake = true
		}
	}

	return config
}

// GetConsumerConfig returns the Kafka consumer configuration
func GetConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false // Disable auto-commit for manual commit
	config.Version = sarama.V2_8_0_0

	// Add SASL authentication if credentials are provided
	if username := os.Getenv("KAFKA_USERNAME"); username != "" {
		if password := os.Getenv("KAFKA_PASSWORD"); password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			config.Net.SASL.User = username
			config.Net.SASL.Password = password
			config.Net.SASL.Handshake = true
		}
	}

	return config
}
