package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"kafka_test/config"
	"kafka_test/models"

	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	groupID string
	topic   string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer() (*Consumer, error) {
	return &Consumer{
		groupID: config.GetGroupID(),
		topic:   config.GetTopicName(),
	}, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer setup completed")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer cleanup completed")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := c.processMessage(session, message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue processing other messages even if one fails
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single Kafka message
func (c *Consumer) processMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	// Log message details
	log.Printf("Received message - Topic: %s, Partition: %d, Offset: %d, Key: %s, Value: %s",
		message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

	// Print headers if available
	if len(message.Headers) > 0 {
		log.Printf("Message headers:")
		for _, header := range message.Headers {
			log.Printf("  %s: %s", string(header.Key), string(header.Value))
		}
	}

	// Try to parse as JSON message first
	var msg models.Message
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		// If JSON parsing fails, treat as plain text message
		log.Printf("Message is plain text: %s", string(message.Value))
	} else {
		log.Printf("Parsed JSON message: %s", msg.String())
	}

	// Simulate some processing time
	time.Sleep(500 * time.Millisecond)

	// Mark message as processed
	session.MarkMessage(message, "")
	log.Printf("Message processed - Partition: %d, Offset: %d \n", message.Partition, message.Offset)

	// Commit immediately for consistency (you could also batch commits for better performance)
	session.Commit()
	return nil
}

// Start starts the consumer
func (c *Consumer) Start(ctx context.Context) error {
	kafkaConfig := config.GetConsumerConfig()

	// Create consumer group
	group, err := sarama.NewConsumerGroup(config.GetBrokers(), c.groupID, kafkaConfig)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer group.Close()

	log.Printf("Consumer started. Consuming messages from topic: %s with group: %s", c.topic, c.groupID)

	// Start consuming in a goroutine
	go func() {
		for {
			topics := []string{c.topic}
			err := group.Consume(ctx, topics, c)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Consumer context cancelled, stopping...")
	return nil
}
