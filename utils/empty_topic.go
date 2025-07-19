package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka_test/config"

	"github.com/Shopify/sarama"
)

var (
	topicName = flag.String("topic", config.GetTopicName(), "Topic to empty")
	broker    = flag.String("broker", config.GetBroker(), "Kafka broker address")
	groupID   = flag.String("group", "topic-emptier", "Consumer group ID")
	confirm   = flag.Bool("confirm", false, "Confirm that you want to empty the topic")
)

// TopicEmptier represents a consumer that empties a topic
type TopicEmptier struct {
	topic   string
	groupID string
}

// NewTopicEmptier creates a new topic emptier
func NewTopicEmptier(topic, groupID string) *TopicEmptier {
	return &TopicEmptier{
		topic:   topic,
		groupID: groupID,
	}
}

// Setup is run at the beginning of a new session
func (t *TopicEmptier) Setup(sarama.ConsumerGroupSession) error {
	log.Printf("Topic emptier setup completed for topic: %s", t.topic)
	return nil
}

// Cleanup is run at the end of a session
func (t *TopicEmptier) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("Topic emptier cleanup completed for topic: %s", t.topic)
	return nil
}

// ConsumeClaim consumes messages to empty the topic
func (t *TopicEmptier) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messageCount := 0

	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				log.Printf("No more messages in partition %d, offset %d", message.Partition, message.Offset)
				return nil
			}

			messageCount++
			log.Printf("Consuming message #%d - Partition: %d, Offset: %d, Key: %s, Value: %s",
				messageCount, message.Partition, message.Offset, string(message.Key), string(message.Value))

			// Mark message as consumed
			session.MarkMessage(message, "")
			session.Commit()

		case <-session.Context().Done():
			return nil
		}
	}
}

// EmptyTopic empties the specified topic
func (t *TopicEmptier) EmptyTopic(ctx context.Context) error {
	// Configure consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Version = sarama.V2_8_0_0

	// Create consumer group
	group, err := sarama.NewConsumerGroup([]string{*broker}, t.groupID, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer group.Close()

	log.Printf("Starting to empty topic: %s", t.topic)

	// Start consuming
	topics := []string{t.topic}
	err = group.Consume(ctx, topics, t)
	if err != nil {
		return fmt.Errorf("error from consumer: %w", err)
	}

	return nil
}

func main() {
	flag.Parse()

	if !*confirm {
		log.Printf("WARNING: This will consume and delete ALL messages from topic '%s'", *topicName)
		log.Printf("To confirm, run with -confirm flag")
		log.Printf("Example: go run utils/empty_topic.go -topic %s -confirm", *topicName)
		os.Exit(1)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topic emptier
	emptier := NewTopicEmptier(*topicName, *groupID)

	// Start emptying in a goroutine
	go func() {
		if err := emptier.EmptyTopic(ctx); err != nil {
			log.Printf("Error emptying topic: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down topic emptier...")
}
