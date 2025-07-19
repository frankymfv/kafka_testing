package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kafka_test/config"
	"kafka_test/models"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

// ChangeDataConsumer represents a Sarama consumer group consumer for ChangeDataMessage
type ChangeDataConsumer struct {
	groupID string
	topic   string
}

// NewChangeDataConsumer creates a new ChangeDataMessage consumer
func NewChangeDataConsumer() (*ChangeDataConsumer, error) {
	return &ChangeDataConsumer{
		groupID: "change-data-consumer-group",
		topic:   config.GetTopicName(),
	}, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *ChangeDataConsumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("ChangeDataConsumer setup completed")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ChangeDataConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("ChangeDataConsumer cleanup completed")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *ChangeDataConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := c.processChangeDataMessage(session, message); err != nil {
				log.Printf("Error processing ChangeDataMessage: %v", err)
				// Continue processing other messages even if one fails
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

// processChangeDataMessage processes a single ChangeDataMessage
func (c *ChangeDataConsumer) processChangeDataMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	// Log message details
	log.Printf("Received message - Topic: %s, Partition: %d, Offset: %d, Key: %s",
		message.Topic, message.Partition, message.Offset, string(message.Key))

	// Print headers if available
	if len(message.Headers) > 0 {
		log.Printf("Message headers:")
		for _, header := range message.Headers {
			log.Printf("  %s: %s", string(header.Key), string(header.Value))
		}
	}

	// Check if this is a ChangeDataMessage
	contentType := ""
	for _, header := range message.Headers {
		if string(header.Key) == "content-type" {
			contentType = string(header.Value)
			break
		}
	}

	if contentType == "application/x-protobuf" {
		// Try to parse as ChangeDataMessage
		var changeDataMsg models.ChangeDataMessage
		if err := proto.Unmarshal(message.Value, &changeDataMsg); err != nil {
			log.Printf("Failed to parse ChangeDataMessage: %v", err)
			log.Printf("Raw message value: %s", string(message.Value))
		} else {
			c.logChangeDataMessage(&changeDataMsg)
		}
	} else {
		// Try to parse as JSON message
		var msg models.Message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			// If JSON parsing fails, treat as plain text message
			log.Printf("Message is plain text: %s", string(message.Value))
		} else {
			log.Printf("Parsed JSON message: %s", msg.String())
		}
	}

	// Simulate some processing time
	time.Sleep(500 * time.Millisecond)

	// Mark message as processed
	session.MarkMessage(message, "")
	log.Printf("Message processed - Partition: %d, Offset: %d \n", message.Partition, message.Offset)

	// Commit immediately for consistency
	session.Commit()
	return nil
}

// logChangeDataMessage logs the details of a ChangeDataMessage
func (c *ChangeDataConsumer) logChangeDataMessage(msg *models.ChangeDataMessage) {
	log.Printf("=== ChangeDataMessage ===")
	log.Printf("Request ID: %s", msg.RequestId)
	log.Printf("Message Number: %d", msg.MessageNumber)
	log.Printf("Total Message Count: %d", msg.TotalMessageCount)
	log.Printf("Tenant UID: %d", msg.TenantUid)
	log.Printf("Records Count: %d", len(msg.Records))

	for i, record := range msg.Records {
		log.Printf("--- Record %d ---", i+1)
		log.Printf("Operation Type: %s", record.OperationType.String())

		switch data := record.Data.(type) {
		case *models.Record_Department:
			dept := data.Department
			log.Printf("Department:")
			log.Printf("  ID: %s", dept.Id)
			log.Printf("  BIID: %s", dept.Biid)
			log.Printf("  Code: %s", dept.Code)
			log.Printf("  Name: %s", dept.Name)
			log.Printf("  Short Name: %s", dept.ShortName)
			log.Printf("  Type: %s", dept.Type)
			log.Printf("  Version: %d", dept.Version)

		case *models.Record_AggregationPattern:
			agg := data.AggregationPattern
			log.Printf("Aggregation Pattern:")
			log.Printf("  ID: %s", agg.Id)
			log.Printf("  BIID: %s", agg.Biid)
			log.Printf("  Name: %s", agg.Name)
			log.Printf("  Version: %d", agg.Version)

		case *models.Record_AggregationPatternDepartment:
			aggDept := data.AggregationPatternDepartment
			log.Printf("Aggregation Pattern Department:")
			log.Printf("  Aggregation Pattern ID: %s", aggDept.AggregationPatternId)
			log.Printf("  Department BIID: %s", aggDept.DepartmentBiid)
			log.Printf("  Path: %s", aggDept.Path)

		default:
			log.Printf("Unknown record type")
		}
	}
	log.Printf("========================")
}

// Start starts the consumer
func (c *ChangeDataConsumer) Start(ctx context.Context) error {
	kafkaConfig := config.GetConsumerConfig()

	// Create consumer group
	group, err := sarama.NewConsumerGroup([]string{config.GetBroker()}, c.groupID, kafkaConfig)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer group.Close()

	log.Printf("ChangeDataConsumer started. Consuming messages from topic: %s with group: %s", c.topic, c.groupID)

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
	log.Println("ChangeDataConsumer context cancelled, stopping...")
	return nil
}

func main() {
	// Create consumer
	c, err := NewChangeDataConsumer()
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer in a goroutine
	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down change data consumer...")
}
