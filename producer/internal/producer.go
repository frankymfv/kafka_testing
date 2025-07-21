package producer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"kafka_test/config"
	"kafka_test/models"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	counter  int
}

// NewProducer creates a new Kafka producer
func NewProducer() (*Producer, error) {
	kafkaConfig := config.GetProducerConfig()

	producer, err := sarama.NewSyncProducer(config.GetBrokers(), kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    config.GetTopicName(),
		counter:  0,
	}, nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// SendMessage sends a message to Kafka
func (p *Producer) SendMessage(ctx context.Context) error {
	p.counter++

	// Alternate between plain text and JSON messages
	if p.counter%2 == 0 {
		// Send plain text message
		content := fmt.Sprintf("Plain text message #%d - ID: %s - Time: %s",
			p.counter,
			uuid.New().String(),
			time.Now().Format(time.RFC3339))

		kafkaMsg := &sarama.ProducerMessage{
			Topic: p.topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("key-%d", p.counter)),
			Value: sarama.StringEncoder(content),
			Headers: []sarama.RecordHeader{
				{Key: []byte("source"), Value: []byte("go-producer")},
				{Key: []byte("counter"), Value: []byte(strconv.Itoa(p.counter))},
				{Key: []byte("type"), Value: []byte("plain-text")},
			},
		}

		partition, offset, err := p.producer.SendMessage(kafkaMsg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		log.Printf("Plain text message sent successfully - Partition: %d, Offset: %d, Message: %s",
			partition, offset, content)
	} else {
		// Send JSON message
		msg := models.NewMessage(p.counter, fmt.Sprintf("JSON message content #%d", p.counter))
		msg.AddHeader("message-id", msg.ID)
		msg.AddHeader("counter", strconv.Itoa(msg.Counter))
		msg.AddHeader("timestamp", msg.Timestamp.Format(time.RFC3339))
		msg.AddHeader("type", "json")

		// Convert to JSON
		jsonData, err := msg.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to convert message to JSON: %w", err)
		}

		// Create Kafka message
		kafkaMsg := &sarama.ProducerMessage{
			Topic: p.topic,
			Key:   sarama.StringEncoder(msg.ID),
			Value: sarama.StringEncoder(jsonData),
			Headers: []sarama.RecordHeader{
				{Key: []byte("message-id"), Value: []byte(msg.ID)},
				{Key: []byte("counter"), Value: []byte(strconv.Itoa(msg.Counter))},
				{Key: []byte("timestamp"), Value: []byte(msg.Timestamp.Format(time.RFC3339))},
				{Key: []byte("type"), Value: []byte("json")},
			},
		}

		// Send the message
		partition, offset, err := p.producer.SendMessage(kafkaMsg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		log.Printf("JSON message sent successfully - Partition: %d, Offset: %d, Message: %s",
			partition, offset, msg.String())
	}

	return nil
}

// Start starts the producer loop
func (p *Producer) Start(ctx context.Context, interval time.Duration) error {
	log.Printf("Producer started. Sending messages to topic: %s every %v", p.topic, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Producer context cancelled, stopping...")
			return nil
		case <-ticker.C:
			if err := p.SendMessage(ctx); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
