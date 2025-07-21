package main

import (
	"fmt"
	"log"
	"strings"

	"kafka_test/config"

	"github.com/Shopify/sarama"
)

func main() {
	// Get broker configuration
	brokers := config.GetBrokers()
	log.Printf("Connecting to brokers: %v", brokers)

	// Create admin client
	admin, err := sarama.NewClusterAdmin(brokers, config.GetProducerConfig())
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer admin.Close()

	// List all topics
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}

	if len(topics) == 0 {
		fmt.Println("No topics found.")
		return
	}

	fmt.Printf("\nFound %d topic(s):\n", len(topics))
	fmt.Println(strings.Repeat("=", 50))

	for topicName, topicDetail := range topics {
		fmt.Printf("Topic: %s\n", topicName)
		fmt.Printf("  Partitions: %d\n", topicDetail.NumPartitions)
		fmt.Printf("  Replication Factor: %d\n", topicDetail.ReplicationFactor)
		fmt.Printf("  Config: %v\n", topicDetail.ConfigEntries)
		fmt.Println(strings.Repeat("-", 30))
	}
}
