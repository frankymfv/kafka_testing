# Kafka Authentication Setup Guide

This guide shows you how to set up SASL/PLAIN authentication for your Kafka cluster.

## Current Status

Your Kafka setup is currently running **without authentication** for simplicity. The authentication code is already implemented in your Go application and will work when you enable it.

## How to Enable Authentication

### Step 1: Update docker-compose.yml

Replace your current Kafka service configuration with this authenticated version:

```yaml
kafka:
  image: confluentinc/cp-kafka:7.4.0
  hostname: kafka
  container_name: kafka
  depends_on:
    - zookeeper
  ports:
    - "9092:9092"
    - "9101:9101"
  environment:
    KAFKA_BROKER_ID: 1
    KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT,SASL_PLAINTEXT:SASL_PLAINTEXT
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092,SASL_PLAINTEXT://kafka:29093
    KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:29092,PLAINTEXT_HOST://0.0.0.0:9092,SASL_PLAINTEXT://0.0.0.0:29093
    KAFKA_SASL_ENABLED_MECHANISMS: PLAIN
    KAFKA_SASL_MECHANISM_INTER_BROKER_PROTOCOL: PLAIN
    KAFKA_OPTS: "-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf"
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
    KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
    KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
    KAFKA_JMX_PORT: 9101
    KAFKA_JMX_HOSTNAME: localhost
    KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
    KAFKA_DELETE_TOPIC_ENABLE: 'true'
    KAFKA_NUM_PARTITIONS: 3
    KAFKA_DEFAULT_REPLICATION_FACTOR: 1
    KAFKA_MIN_INSYNC_REPLICAS: 1
  volumes:
    - kafka-data:/var/lib/kafka/data
    - ./kafka_jaas.conf:/etc/kafka/kafka_jaas.conf
```

### Step 2: Update Kafka UI Configuration

Update the kafka-ui service:

```yaml
kafka-ui:
  image: provectuslabs/kafka-ui:latest
  container_name: kafka-ui
  depends_on:
    - kafka
  ports:
    - "8080:8080"
  environment:
    KAFKA_CLUSTERS_0_NAME: local
    KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29093
    KAFKA_CLUSTERS_0_ZOOKEEPER: zookeeper:2181
    KAFKA_CLUSTERS_0_PROPERTIES_SECURITY_PROTOCOL: SASL_PLAINTEXT
    KAFKA_CLUSTERS_0_PROPERTIES_SASL_MECHANISM: PLAIN
    KAFKA_CLUSTERS_0_PROPERTIES_SASL_JAAS_CONFIG: "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\";"
```

### Step 3: Update Your .env File

Add authentication credentials to your `.env` file:

```bash
# Kafka Authentication
KAFKA_USERNAME=admin
KAFKA_PASSWORD=admin-secret
```

### Step 4: Update Broker Configuration

Update your `config/kafka.go` to use the SASL port when authentication is enabled:

```go
// GetBroker returns the Kafka broker address from environment variable or default
func GetBroker() string {
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		return broker
	}
	// Use SASL port if authentication is enabled
	if os.Getenv("KAFKA_USERNAME") != "" {
		return "localhost:29093"
	}
	return "localhost:9092"
}
```

### Step 5: Restart Services

```bash
docker-compose down -v
docker-compose up -d
```

## Available Users

The JAAS configuration includes these users:

- **admin** / **admin-secret** - Full access (recommended for development)
- **producer** / **producer-secret** - Producer access only
- **consumer** / **consumer-secret** - Consumer access only

## Testing Authentication

### Test with Kafka CLI

```bash
# Create a topic with authentication
docker exec kafka kafka-topics --bootstrap-server localhost:29093 \
  --create --topic test-auth-topic --partitions 3 --replication-factor 1

# List topics with authentication
docker exec kafka kafka-topics --bootstrap-server localhost:29093 --list
```

### Test with Your Go Application

1. Add credentials to your `.env` file
2. Run your producer/consumer applications
3. They will automatically use authentication

## Security Notes

- **SASL/PLAIN** sends credentials in plain text
- For production environments, consider:
  - **SASL/SCRAM** for better security
  - **SSL/TLS** encryption
  - **Kerberos** for enterprise environments

## Troubleshooting

### Common Issues

1. **Zookeeper Authentication Errors**: The setup disables SASL for Zookeeper connections
2. **Port Conflicts**: Make sure ports 29093 and 9092 are available
3. **JAAS Configuration**: Ensure the `kafka_jaas.conf` file is properly mounted

### Verification Commands

```bash
# Check if Kafka is running with authentication
docker-compose ps

# Check Kafka logs
docker-compose logs kafka

# Test connection from inside container
docker exec kafka kafka-topics --bootstrap-server localhost:29093 --list
```

## Current Working Setup

Your current setup works without authentication and includes:

- ✅ Environment variable configuration
- ✅ `.env` file support
- ✅ Authentication code in Go applications
- ✅ Partition configuration (3 partitions)
- ✅ Kafka UI access

To enable authentication, follow the steps above. 