# Kafka Test Project

This project demonstrates a complete Kafka setup using Go with Shopify Sarama library, featuring a producer, consumer with manual commit, and Kafka UI for visualization.

## Project Structure

```
kafka_test/
├── docker-compose.yml    # Kafka, Zookeeper, and Kafka UI services
├── go.mod               # Go module dependencies
├── producer/
│   └── main.go         # Kafka producer application
├── consumer/
│   └── main.go         # Kafka consumer application with manual commit
└── README.md           # This documentation
```

## Architecture Overview

### Components

1. **Kafka Broker**: Apache Kafka message broker running on port 9092
2. **Zookeeper**: Coordination service for Kafka (port 2181)
3. **Kafka UI**: Web interface for monitoring Kafka (port 8080)
4. **Go Producer**: Application that sends messages to Kafka
5. **Go Consumer**: Application that consumes messages with manual commit

### Message Flow

```
Go Producer → Kafka Topic → Go Consumer (Manual Commit)
                ↓
            Kafka UI (Monitoring)
```

## Role of Kafka in This Project

### 1. Message Broker
Kafka acts as a distributed message broker that:
- Receives messages from the producer
- Stores them in topics with persistence
- Delivers messages to consumers
- Handles message ordering and partitioning

### 2. Decoupling
Kafka decouples the producer and consumer:
- Producer can send messages without waiting for consumer processing
- Consumer can process messages at its own pace
- Both applications can be scaled independently

### 3. Reliability
Kafka provides:
- Message persistence (messages survive broker restarts)
- Fault tolerance through replication
- At-least-once delivery semantics
- Message ordering within partitions

### 4. Scalability
Kafka enables:
- Horizontal scaling of both producers and consumers
- Partition-based parallelism
- Consumer group load balancing

## How It Works

### 1. Message Production
The producer (`producer/main.go`):
- Connects to Kafka broker on `localhost:9092`
- Sends messages to topic `test-topic` every 2 seconds
- Includes message headers (message-id, counter, timestamp)
- Uses synchronous producer for guaranteed delivery
- Generates unique UUIDs for each message

### 2. Message Consumption
The consumer (`consumer/main.go`):
- Connects to Kafka as part of consumer group `test-consumer-group`
- Consumes messages from topic `test-topic`
- **Manually commits** each message after processing
- Processes messages with 500ms delay to simulate work
- Logs message details including headers

### 3. Manual Commit Strategy
This project demonstrates manual commit by:
- Disabling auto-commit: `config.Consumer.Offsets.AutoCommit.Enable = false`
- Using `session.MarkMessage(message, "")` to commit each message individually
- This ensures messages are only marked as consumed after successful processing

### 4. Monitoring
Kafka UI provides:
- Real-time topic monitoring
- Message browsing and inspection
- Consumer group status
- Partition information
- Message headers and payload viewing

## Configuration

### Environment Variables

The application now uses environment variables for Kafka configuration. You can set these variables or use the default values:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `KAFKA_TOPIC_NAME` | `test-topic` | Kafka topic name |
| `KAFKA_BROKER` | `localhost:9092` | Single Kafka broker address (host:port) |
| `KAFKA_BROKERS` | (none) | Multiple Kafka broker addresses (comma-separated) |
| `KAFKA_GROUP_ID` | `test-consumer-group` | Main Kafka consumer group ID |
| `KAFKA_CHANGE_DATA_GROUP_ID` | `change-data-consumer-group` | ChangeData consumer group ID |
| `KAFKA_TOPIC_EMPTIER_GROUP_ID` | `topic-emptier` | Topic emptier consumer group ID |
| `KAFKA_USERNAME` | (none) | Kafka username for authentication |
| `KAFKA_PASSWORD` | (none) | Kafka password for authentication |

### Example Environment Setup

The application automatically loads environment variables from a `.env` file in the project root. Create a `.env` file with:

```bash
# Kafka Configuration
KAFKA_TOPIC_NAME=test-topic
KAFKA_BROKER=localhost:9092
KAFKA_BROKERS=localhost:9097,localhost:9098,localhost:9099

# Consumer Group IDs
KAFKA_GROUP_ID=test-consumer-group
KAFKA_CHANGE_DATA_GROUP_ID=change-data-consumer-group
KAFKA_TOPIC_EMPTIER_GROUP_ID=topic-emptier

# Kafka Authentication (optional)
KAFKA_USERNAME=admin
KAFKA_PASSWORD=admin-secret
```

You can also set environment variables directly:

```bash
export KAFKA_TOPIC_NAME=my-topic
export KAFKA_BROKER=kafka.example.com:9092
export KAFKA_BROKERS=kafka1.example.com:9092,kafka2.example.com:9092,kafka3.example.com:9092
export KAFKA_GROUP_ID=my-consumer-group
export KAFKA_CHANGE_DATA_GROUP_ID=my-change-data-group
export KAFKA_TOPIC_EMPTIER_GROUP_ID=my-topic-emptier
export KAFKA_USERNAME=admin
export KAFKA_PASSWORD=admin-secret
```

## Authentication

This setup includes SASL/PLAIN authentication for Kafka. The default credentials are:

- **Username**: `admin`
- **Password**: `admin-secret`

### Available Users

The JAAS configuration includes these users:
- `admin` / `admin-secret` - Full access
- `producer` / `producer-secret` - Producer access
- `consumer` / `consumer-secret` - Consumer access

### Security Notes

- SASL/PLAIN sends credentials in plain text
- For production, consider using SASL/SCRAM or SSL/TLS
- The authentication is optional - if no credentials are provided, the application will connect without authentication

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later

### 1. Start Kafka Infrastructure
```bash
docker-compose up -d
```

This starts:
- Zookeeper on port 2181
- Kafka on port 9092
- Kafka UI on port 8080

### 2. Install Go Dependencies
```bash
go mod tidy
```

### 3. Run the Producer
```bash
go run producer/main.go
```

The producer will start sending messages every 2 seconds.

### 4. Run the Consumer
```bash
go run consumer/main.go
```

The consumer will start processing messages with manual commit.

### 5. Monitor with Kafka UI
Open http://localhost:8080 in your browser to:
- View topics and messages
- Monitor consumer groups
- Inspect message headers and content

## Key Features

### Producer Features
- **Synchronous Production**: Ensures message delivery confirmation
- **Message Headers**: Includes metadata (ID, counter, timestamp)
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals
- **Retry Logic**: Configurable retry attempts for failed sends

### Consumer Features
- **Manual Commit**: Explicit control over message acknowledgment
- **Consumer Groups**: Enables load balancing and fault tolerance
- **Message Headers**: Processes and logs message metadata
- **Graceful Shutdown**: Proper cleanup on termination

### Kafka Configuration
- **Auto Topic Creation**: Topics are created automatically
- **Single Broker**: Development setup with one Kafka broker
- **Persistence**: Messages are stored on disk
- **JMX Monitoring**: Available on port 9101

## Testing the Setup

1. **Start the infrastructure**: `docker-compose up -d`
2. **Run producer**: `go run producer/main.go`
3. **Run consumer**: `go run consumer/main.go`
4. **Monitor**: Visit http://localhost:8080
5. **Verify**: Check logs to see messages being produced and consumed

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure Docker Compose services are running
2. **Topic Not Found**: Topics are auto-created when first message is sent
3. **Consumer Not Receiving**: Check consumer group configuration
4. **Port Conflicts**: Ensure ports 9092, 2181, and 8080 are available

### Logs
- Producer logs show sent messages with partition/offset
- Consumer logs show received messages and commit confirmations
- Docker logs: `docker-compose logs kafka`

## Cleanup

```bash
# Stop applications (Ctrl+C)
# Stop infrastructure
docker-compose down -v
```

## Production Considerations

For production use, consider:
- Multiple Kafka brokers for high availability
- Proper security (SASL/SSL)
- Monitoring and alerting
- Log aggregation
- Resource planning
- Backup strategies 