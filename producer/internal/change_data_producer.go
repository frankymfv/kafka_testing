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
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChangeDataProducer represents a Kafka producer for ChangeDataMessage
type ChangeDataProducer struct {
	producer sarama.SyncProducer
	topic    string
	counter  uint32
}

// NewChangeDataProducer creates a new ChangeDataMessage producer
func NewChangeDataProducer() (*ChangeDataProducer, error) {
	kafkaConfig := config.GetProducerConfig()

	producer, err := sarama.NewSyncProducer([]string{config.GetBroker()}, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &ChangeDataProducer{
		producer: producer,
		topic:    config.GetTopicName(),
		counter:  0,
	}, nil
}

// Close closes the producer
func (p *ChangeDataProducer) Close() error {
	return p.producer.Close()
}

// createSampleDepartment creates a sample department record
func (p *ChangeDataProducer) createSampleDepartment() *models.Department {
	now := time.Now()
	today := &date.Date{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
		Day:   int32(now.Day()),
	}

	return &models.Department{
		Id:                    fmt.Sprintf("dept-%d", p.counter),
		Biid:                  fmt.Sprintf("biid-%d", p.counter),
		Type:                  "department",
		TenantUid:             12345,
		Code:                  fmt.Sprintf("DEPT%03d", p.counter),
		Name:                  fmt.Sprintf("Department %d", p.counter),
		ShortName:             fmt.Sprintf("Dept%d", p.counter),
		ValidFrom:             today,
		ValidTo:               nil, // No end date
		TransactionFrom:       timestamppb.Now(),
		TransactionTo:         nil,
		TransactionFromBy:     1001,
		TransactionToBy:       0,
		TransactionFromByName: "System User",
		TransactionToByName:   "",
		TransactionFromBySrv:  "kafka-producer",
		TransactionToBySrv:    "",
		DispOrder:             p.counter,
		SearchKey:             fmt.Sprintf("dept%d", p.counter),
		Version:               1,
		ParentBiid:            nil,
		ParentCode:            nil,
		ParentName:            nil,
	}
}

// createSampleAggregationPattern creates a sample aggregation pattern
func (p *ChangeDataProducer) createSampleAggregationPattern() *models.AggregationPattern {
	return &models.AggregationPattern{
		Id:                    fmt.Sprintf("agg-%d", p.counter),
		Biid:                  fmt.Sprintf("agg-biid-%d", p.counter),
		TenantUid:             12345,
		Name:                  fmt.Sprintf("Aggregation Pattern %d", p.counter),
		TransactionFrom:       timestamppb.Now(),
		TransactionTo:         nil,
		TransactionFromBy:     1001,
		TransactionToBy:       0,
		TransactionFromByName: "System User",
		TransactionToByName:   "",
		TransactionFromBySrv:  "kafka-producer",
		TransactionToBySrv:    "",
		Version:               1,
	}
}

// SendChangeDataMessage sends a ChangeDataMessage to Kafka
func (p *ChangeDataProducer) SendChangeDataMessage(ctx context.Context) error {
	p.counter++

	// Create sample records
	var records []*models.Record

	// Add a department record
	dept := p.createSampleDepartment()
	deptRecord := &models.Record{
		OperationType: models.Record_CREATE,
		Data: &models.Record_Department{
			Department: dept,
		},
	}
	records = append(records, deptRecord)

	// Add an aggregation pattern record (every 3rd message)
	if p.counter%3 == 0 {
		agg := p.createSampleAggregationPattern()
		aggRecord := &models.Record{
			OperationType: models.Record_UPDATE,
			Data: &models.Record_AggregationPattern{
				AggregationPattern: agg,
			},
		}
		records = append(records, aggRecord)
	}

	// Create ChangeDataMessage
	changeDataMsg := &models.ChangeDataMessage{
		RequestId:         uuid.New().String(),
		MessageNumber:     p.counter,
		TotalMessageCount: p.counter + 1,
		TenantUid:         12345,
		Records:           records,
	}

	// Serialize to protobuf
	protoData, err := proto.Marshal(changeDataMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	// Create Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(changeDataMsg.RequestId),
		Value: sarama.ByteEncoder(protoData),
		Headers: []sarama.RecordHeader{
			{Key: []byte("request-id"), Value: []byte(changeDataMsg.RequestId)},
			{Key: []byte("message-number"), Value: []byte(strconv.FormatUint(uint64(changeDataMsg.MessageNumber), 10))},
			{Key: []byte("tenant-uid"), Value: []byte(strconv.FormatUint(changeDataMsg.TenantUid, 10))},
			{Key: []byte("content-type"), Value: []byte("application/x-protobuf")},
			{Key: []byte("message-type"), Value: []byte("ChangeDataMessage")},
		},
	}

	// Send the message
	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("ChangeDataMessage sent successfully - Partition: %d, Offset: %d, RequestID: %s, Records: %d",
		partition, offset, changeDataMsg.RequestId, len(changeDataMsg.Records))

	return nil
}

// Start starts the producer loop
func (p *ChangeDataProducer) Start(ctx context.Context, interval time.Duration) error {
	log.Printf("ChangeDataProducer started. Sending messages to topic: %s every %v", p.topic, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("ChangeDataProducer context cancelled, stopping...")
			return nil
		case <-ticker.C:
			if err := p.SendChangeDataMessage(ctx); err != nil {
				log.Printf("Error sending ChangeDataMessage: %v", err)
			}
		}
	}
}
