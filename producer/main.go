package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kafka_test/producer/internal"
)

func main() {
	// Create producer
	p, err := producer.NewProducer()
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start producer in a goroutine
	go func() {
		if err := p.Start(ctx, 2*time.Second); err != nil {
			log.Printf("Producer error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down producer...")
}
