package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka_test/consumer/internal"
)

func main() {
	// Create consumer
	c, err := internal.NewConsumer()
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
	log.Println("Shutting down consumer...")
}
