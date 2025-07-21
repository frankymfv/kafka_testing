# Multiple Kafka Brokers Setup

This document explains the changes made to support multiple Kafka brokers in your setup.

## Overview

The configuration has been updated to support multiple Kafka brokers for high availability and load balancing. The system now supports:

- **Single broker**: `KAFKA_BROKER` environment variable (backward compatibility)
- **Multiple brokers**: `KAFKA_BROKERS` environment variable (comma-separated list)

## Broker Configuration

### New Brokers
- `localhost:9097`
- `localhost:9098` 
- `localhost:9099`

### Environment Variables

| Variable | Value | Description |
|----------|-------|-------------|
| `KAFKA_BROKER` | `localhost:9092` | Single broker (legacy) |
| `KAFKA_BROKERS` | `localhost:9097,localhost:9098,localhost:9099` | Multiple brokers |

## Configuration Changes

### 1. Config Package (`config/kafka.go`)

Added new function to support multiple brokers:

```go
// GetBrokers returns the Kafka broker addresses from environment variable or default
func GetBrokers() []string {
    if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
        // Split by comma and trim whitespace
        brokerList := strings.Split(brokers, ",")
        for i, broker := range brokerList {
            brokerList[i] = strings.TrimSpace(broker)
        }
        return brokerList
    }
    // Fallback to single broker for backward compatibility
    if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
        return []string{broker}
    }
    return []string{"localhost:9092"}
}
```

### 2. Updated Components

#### Consumers
- `consumer/internal/consumer.go` - Uses `config.GetBrokers()`
- `change_data_consumer/main.go` - Uses `config.GetBrokers()`
- `utils/empty_topic.go` - Uses `config.GetBrokers()` with fallback

#### Producers
- `producer/internal/producer.go` - Uses `config.GetBrokers()`
- `producer/internal/change_data_producer.go` - Uses `config.GetBrokers()`
- `change_data_producer/main.go` - Uses `config.GetBrokers()`

### 3. Environment Configuration

#### Updated Files
- `env.example` - Added `KAFKA_BROKERS` configuration
- `README.md` - Updated documentation
- `KAFKA_CLI_GUIDE.md` - Updated CLI guide
- `kafka-cli.sh` - Updated broker addresses
- `kafka_client.properties` - Added comment about multiple brokers

## Usage

### Environment Setup

Create a `.env` file with:

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

### Priority Order

1. **`KAFKA_BROKERS`** - Multiple brokers (comma-separated)
2. **`KAFKA_BROKER`** - Single broker (fallback)
3. **Default** - `localhost:9092` (if neither is set)

### Examples

#### Single Broker
```bash
export KAFKA_BROKER=localhost:9092
```

#### Multiple Brokers
```bash
export KAFKA_BROKERS=localhost:9097,localhost:9098,localhost:9099
```

#### Mixed Configuration
```bash
export KAFKA_BROKER=localhost:9092
export KAFKA_BROKERS=localhost:9097,localhost:9098,localhost:9099
# KAFKA_BROKERS takes precedence
```

## Benefits

### 1. High Availability
- If one broker fails, others continue to serve requests
- Automatic failover between brokers

### 2. Load Balancing
- Producers and consumers distribute load across multiple brokers
- Better performance under high load

### 3. Scalability
- Easy to add more brokers without code changes
- Horizontal scaling capability

### 4. Fault Tolerance
- Improved resilience against broker failures
- Better message delivery guarantees

## Backward Compatibility

The system maintains full backward compatibility:

- Existing `KAFKA_BROKER` configurations continue to work
- Single broker setups are unaffected
- Gradual migration to multiple brokers is supported

## Testing

### Verify Configuration

```bash
# Check which brokers are being used
go run consumer/main.go
# Look for log: "Consumer started. Consuming messages from topic: ..."

# Test producer
go run producer/main.go
# Look for log: "Producer started. Sending messages to topic: ..."
```

### Kafka CLI Testing

```bash
# List topics using multiple brokers
./kafka-cli.sh topics

# Create topic
./kafka-cli.sh create-topic test-multi-broker

# Send messages
./kafka-cli.sh producer test-multi-broker
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure all brokers are running on specified ports
   - Check firewall settings for all broker ports

2. **Authentication Errors**
   - Verify SASL configuration for all brokers
   - Check credentials in `kafka_client.properties`

3. **Partial Connectivity**
   - Some brokers may be down
   - Check broker health individually

### Debug Commands

```bash
# Test individual broker connectivity
telnet localhost 9097
telnet localhost 9098
telnet localhost 9099

# Check broker status
docker-compose ps

# View logs
docker-compose logs kafka
```

## Next Steps

1. **Update Docker Compose** - Configure multiple Kafka brokers in docker-compose.yml
2. **Load Testing** - Test performance with multiple brokers
3. **Monitoring** - Set up monitoring for all brokers
4. **Production** - Deploy with proper broker configuration 