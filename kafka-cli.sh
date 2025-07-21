#!/bin/bash

# Kafka CLI Helper Script
# Usage: ./kafka-cli.sh [command] [options]

KAFKA_BROKER="localhost:9097,localhost:9098,localhost:9099"
CONFIG_FILE="kafka_client.properties"

case "$1" in
    "topics")
        echo "Listing topics..."
        kafka-topics --bootstrap-server $KAFKA_BROKER --command-config $CONFIG_FILE --list
        ;;
    "create-topic")
        if [ -z "$2" ]; then
            echo "Usage: $0 create-topic <topic-name>"
            exit 1
        fi
        echo "Creating topic: $2"
        kafka-topics --bootstrap-server $KAFKA_BROKER --command-config $CONFIG_FILE \
          --create --topic "$2" --partitions 3 --replication-factor 1
        ;;
    "describe-topic")
        if [ -z "$2" ]; then
            echo "Usage: $0 describe-topic <topic-name>"
            exit 1
        fi
        echo "Describing topic: $2"
        kafka-topics --bootstrap-server $KAFKA_BROKER --command-config $CONFIG_FILE \
          --describe --topic "$2"
        ;;
    "producer")
        if [ -z "$2" ]; then
            echo "Usage: $0 producer <topic-name>"
            exit 1
        fi
        echo "Starting producer for topic: $2"
        echo "Type messages and press Enter. Ctrl+C to exit."
        kafka-console-producer --bootstrap-server $KAFKA_BROKER --producer.config $CONFIG_FILE \
          --topic "$2"
        ;;
    "consumer")
        if [ -z "$2" ]; then
            echo "Usage: $0 consumer <topic-name>"
            exit 1
        fi
        echo "Starting consumer for topic: $2"
        echo "Press Ctrl+C to exit."
        kafka-console-consumer --bootstrap-server $KAFKA_BROKER --consumer.config $CONFIG_FILE \
          --topic "$2" --from-beginning
        ;;
    "groups")
        echo "Listing consumer groups..."
        kafka-consumer-groups --bootstrap-server $KAFKA_BROKER --command-config $CONFIG_FILE --list
        ;;
    "describe-group")
        if [ -z "$2" ]; then
            echo "Usage: $0 describe-group <group-name>"
            exit 1
        fi
        echo "Describing consumer group: $2"
        kafka-consumer-groups --bootstrap-server $KAFKA_BROKER --command-config $CONFIG_FILE \
          --describe --group "$2"
        ;;
    *)
        echo "Kafka CLI Helper"
        echo "Usage: $0 [command] [options]"
        echo ""
        echo "Commands:"
        echo "  topics                    - List all topics"
        echo "  create-topic <name>       - Create a new topic"
        echo "  describe-topic <name>     - Describe topic details"
        echo "  producer <topic>          - Start console producer"
        echo "  consumer <topic>          - Start console consumer"
        echo "  groups                    - List consumer groups"
        echo "  describe-group <name>     - Describe consumer group"
        echo ""
        echo "Examples:"
        echo "  $0 topics"
        echo "  $0 create-topic my-topic"
        echo "  $0 producer test-topic"
        echo "  $0 consumer test-topic"
        ;;
esac 