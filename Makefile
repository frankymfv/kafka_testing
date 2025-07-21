.PHONY: help start stop clean build producer consumer change-data-producer change-data-consumer test kafka-topics kafka-create-topic kafka-producer kafka-consumer kafka-groups list-topics

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
	@echo ""
	@echo "Kafka CLI Commands:"
	@echo "  kafka-topics       - List all Kafka topics"
	@echo "  kafka-create-topic - Create a new topic (TOPIC=name)"
	@echo "  kafka-producer     - Start console producer (TOPIC=name)"
	@echo "  kafka-consumer     - Start console consumer (TOPIC=name)"
	@echo "  kafka-groups       - List consumer groups"
	@echo "  list-topics        - List all Kafka topics"

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
	go build -o list_topics/list_topics ./list_topics

build_consumer_linux:
	GOOS=linux GOARCH=amd64 go build -o my_consumer_linux ./change_data_consumer

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

# Kafka CLI commands
kafka-topics:
	@echo "Listing Kafka topics..."
	@./kafka-cli.sh topics

kafka-create-topic:
	@echo "Creating Kafka topic..."
	@./kafka-cli.sh create-topic $(TOPIC)

kafka-producer:
	@echo "Starting Kafka console producer..."
	@./kafka-cli.sh producer $(TOPIC)

kafka-consumer:
	@echo "Starting Kafka console consumer..."
	@./kafka-cli.sh consumer $(TOPIC)

kafka-groups:
	@echo "Listing consumer groups..."
	@./kafka-cli.sh groups

# List all Kafka topics
list-topics:
	@echo "Listing all Kafka topics..."
	@go run list_topics/main.go 

copy-to-bastion-consumer:
	@echo "Building my_consumer..."
	@GOOS=linux GOARCH=amd64 go build -o my_consumer_linux ./change_data_consumer
	@echo "Copying files to bastion..."
	@scp -i ~/.ssh/mfj_bastion -P 1135 my_consumer_linux bui.cong.hoang@nlb2.bastion.test.musubu.co.in:~/

copy-to-bastion-list-topics:
	@echo "Building list_topics..."
	@GOOS=linux GOARCH=amd64 go build -o list_kafka_topic ./list_topic
	@echo "Copying files to bastion..."
	@scp -i ~/.ssh/mfj_bastion -P 1135 list_kafka_topic bui.cong.hoang@nlb2.bastion.test.musubu.co.in:~/