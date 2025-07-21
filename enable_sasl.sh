`#!/bin/bash

echo "🔐 Enabling SASL Security for Kafka..."
echo "======================================"

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: docker-compose.yml not found. Please run this script from the project root."
    exit 1
fi

# Backup current configuration
echo "📋 Backing up current configuration..."
cp docker-compose.yml docker-compose.yml.backup

# Update docker-compose.yml with SASL configuration
echo "🔧 Updating docker-compose.yml with SASL configuration..."

# Update Kafka service configuration
sed -i '' '/KAFKA_LISTENER_SECURITY_PROTOCOL_MAP:/d' docker-compose.yml
sed -i '' '/KAFKA_ADVERTISED_LISTENERS:/d' docker-compose.yml
sed -i '' '/KAFKA_SASL_ENABLED_MECHANISMS:/d' docker-compose.yml
sed -i '' '/KAFKA_SASL_MECHANISM_INTER_BROKER_PROTOCOL:/d' docker-compose.yml
sed -i '' '/KAFKA_OPTS:/d' docker-compose.yml

# Add SASL configuration after KAFKA_ZOOKEEPER_CONNECT
sed -i '' '/KAFKA_ZOOKEEPER_CONNECT:/a\
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT,SASL_PLAINTEXT:SASL_PLAINTEXT\
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092,SASL_PLAINTEXT://kafka:29093\
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:29092,PLAINTEXT_HOST://0.0.0.0:9092,SASL_PLAINTEXT://0.0.0.0:29093\
      KAFKA_SASL_ENABLED_MECHANISMS: PLAIN\
      KAFKA_SASL_MECHANISM_INTER_BROKER_PROTOCOL: PLAIN\
      KAFKA_OPTS: "-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf -Dzookeeper.sasl.client=false"' docker-compose.yml

# Add JAAS configuration volume
sed -i '' '/kafka-data:\/var\/lib\/kafka\/data/a\
      - .\/kafka_jaas.conf:\/etc\/kafka\/kafka_jaas.conf' docker-compose.yml

# Update Kafka UI configuration
sed -i '' '/KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092/a\
      KAFKA_CLUSTERS_0_PROPERTIES_SECURITY_PROTOCOL: SASL_PLAINTEXT\
      KAFKA_CLUSTERS_0_PROPERTIES_SASL_MECHANISM: PLAIN\
      KAFKA_CLUSTERS_0_PROPERTIES_SASL_JAAS_CONFIG: "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\";"' docker-compose.yml

# Update Kafka UI bootstrap servers
sed -i '' 's/KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092/KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29093/' docker-compose.yml

# Add authentication credentials to .env file
echo "🔑 Adding authentication credentials to .env file..."
echo "" >> .env
echo "# Kafka Authentication" >> .env
echo "KAFKA_USERNAME=admin" >> .env
echo "KAFKA_PASSWORD=admin-secret" >> .env

# Update Go configuration to use SASL port
echo "🔧 Updating Go configuration..."
sed -i '' '/\/\/ Use SASL port if authentication is enabled/a\
	// Use SASL port if authentication is enabled\
	if os.Getenv("KAFKA_USERNAME") != "" {\
		return "localhost:29093"\
	}' config/kafka.go

echo "✅ SASL configuration updated!"
echo ""
echo "🔄 Restarting services..."
docker-compose down -v
docker-compose up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 20

echo ""
echo "🔍 Checking service status..."
docker-compose ps

echo ""
echo "🎉 SASL Security is now enabled!"
echo ""
echo "📋 Available credentials:"
echo "   Username: admin"
echo "   Password: admin-secret"
echo ""
echo "🌐 Access points:"
echo "   - Kafka UI: http://localhost:8080"
echo "   - Kafka SASL Port: localhost:29093"
echo ""
echo "📖 For more information, see AUTHENTICATION_SETUP.md" `