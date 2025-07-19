package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Message represents a Kafka message with metadata
type Message struct {
	ID        string            `json:"id"`
	Counter   int               `json:"counter"`
	Content   string            `json:"content"`
	Timestamp time.Time         `json:"timestamp"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// NewMessage creates a new message with the given counter and content
func NewMessage(counter int, content string) *Message {
	return &Message{
		ID:        uuid.New().String(),
		Counter:   counter,
		Content:   content,
		Timestamp: time.Now(),
		Headers:   make(map[string]string),
	}
}

// AddHeader adds a header to the message
func (m *Message) AddHeader(key, value string) {
	m.Headers[key] = value
}

// ToJSON converts the message to JSON string
func (m *Message) ToJSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}
	return string(data), nil
}

// String returns a formatted string representation of the message
func (m *Message) String() string {
	return fmt.Sprintf("Message #%d - ID: %s - Time: %s - Content: %s",
		m.Counter, m.ID, m.Timestamp.Format(time.RFC3339), m.Content)
}
