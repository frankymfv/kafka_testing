.PHONY: help start stop clean build producer consumer change-data-producer change-data-consumer test

# Default target
help:
	@echo "Available commands:"
	@echo "  start     - Start Kafka infrastructure (Docker Compose)"
	@echo "  stop      - Stop Kafka infrastructure"
	@echo "  clean     - Stop and remove all containers and volumes"
	@echo "  build     - Build Go applications"
	@echo "  producer  - Run the Kafka producer"
	@echo "  consumer  - Run the Kafka consumer"
	@echo "  change-data-producer  - Run the ChangeDataMessage producer"
	@echo "  change-data-consumer  - Run the ChangeDataMessage consumer"
	@echo "  test      - Run a complete test (start infra, producer, consumer)"
	@echo "  logs      - Show Kafka logs"
	@echo "  ui        - Open Kafka UI in browser"

# Start Kafka infrastructure
start:
	@echo "Starting Kafka infrastructure..."
	docker-compose up -d
	@echo "Kafka UI available at: http://localhost:8080"

# Stop Kafka infrastructure
stop:
	@echo "Stopping Kafka infrastructure..."
	docker-compose down

# Clean everything
clean:
	@echo "Cleaning up Kafka infrastructure..."
	docker-compose down -v
	@echo "Removing Go build artifacts..."
	rm -rf producer/producer consumer/consumer

# Build Go applications
build:
	@echo "Building Go applications..."
	go mod tidy
	go build -o producer/producer ./producer
	go build -o consumer/consumer ./consumer
	go build -o change_data_producer/change_data_producer ./change_data_producer
	go build -o change_data_consumer/change_data_consumer ./change_data_consumer

# Run producer
producer:
	@echo "Running Kafka producer..."
	go run producer/main.go

# Run consumer
consumer:
	@echo "Running Kafka consumer..."
	go run consumer/main.go

# Run ChangeDataMessage producer
change-data-producer:
	@echo "Running ChangeDataMessage producer..."
	go run change_data_producer/main.go

# Run ChangeDataMessage consumer
change-data-consumer:
	@echo "Running ChangeDataMessage consumer..."
	go run change_data_consumer/main.go

# Show Kafka logs
logs:
	docker-compose logs -f kafka

# Open Kafka UI
ui:
	@echo "Opening Kafka UI..."
	open http://localhost:8080

# Complete test setup
test: start
	@echo "Waiting for Kafka to be ready..."
	@sleep 10
	@echo "Starting producer in background..."
	@go run producer/main.go &
	@echo "Starting consumer..."
	@go run consumer/main.go 