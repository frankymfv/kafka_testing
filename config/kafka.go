package config

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	"github.com/xdg-go/scram"
)

func init() {
	// Load .env file if it exists
	godotenv.Load()
}

// GetTopicName returns the Kafka topic name from environment variable or default
func GetTopicName() string {
	if topic := os.Getenv("KAFKA_TOPIC_NAME"); topic != "" {
		return topic
	}
	return "test-topic"
}

// GetBroker returns the Kafka broker address from environment variable or default
func GetBroker() string {
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		return broker
	}
	return "localhost:9092"
}

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

// GetGroupID returns the Kafka consumer group ID from environment variable or default
func GetGroupID() string {
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		return groupID
	}
	return "test-consumer-group"
}

// GetChangeDataGroupID returns the ChangeData consumer group ID from environment variable or default
func GetChangeDataGroupID() string {
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		return groupID
	}
	return "change-data-consumer-group"
}

// GetTopicEmptierGroupID returns the topic emptier consumer group ID from environment variable or default
func GetTopicEmptierGroupID() string {
	if groupID := os.Getenv("KAFKA_TOPIC_EMPTIER_GROUP_ID"); groupID != "" {
		return groupID
	}
	return "topic-emptier"
}

// GetProducerConfig returns the Kafka producer configuration
func GetProducerConfig() *sarama.Config {
	// Get security settings from environment variables with defaults
	securityProtocol := os.Getenv("KAFKA_SECURITY_PROTOCOL")
	if securityProtocol == "" {
		securityProtocol = "PLAINTEXT"
	}

	saslMechanism := os.Getenv("KAFKA_SASL_MECHANISM")
	if saslMechanism == "" {
		saslMechanism = "PLAIN"
	}

	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	sslVerify := true
	if sslVerifyStr := os.Getenv("KAFKA_SSL_VERIFY"); sslVerifyStr != "" {
		sslVerify = sslVerifyStr == "true"
	}

	return GetProducerConfigWithSecurity(securityProtocol, saslMechanism, username, password, sslVerify)
}

// GetProducerConfigWithSecurity returns the Kafka producer configuration with custom security settings
func GetProducerConfigWithSecurity(securityProtocol, saslMechanism, username, password string, sslVerify bool) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Version = sarama.V2_8_0_0

	// Configure security protocol
	switch strings.ToUpper(securityProtocol) {
	case "PLAINTEXT":
		config.Net.TLS.Enable = false
		config.Net.SASL.Enable = false
	case "SSL":
		config.Net.TLS.Enable = true
		config.Net.SASL.Enable = false
	case "SASL_PLAINTEXT":
		config.Net.TLS.Enable = false
		config.Net.SASL.Enable = true
	case "SASL_SSL":
		config.Net.TLS.Enable = true
		config.Net.SASL.Enable = true
	default:
		// Default to PLAINTEXT if invalid protocol
		config.Net.TLS.Enable = false
		config.Net.SASL.Enable = false
	}

	// Configure SASL if enabled
	if config.Net.SASL.Enable {
		if username == "" || password == "" {
			// Fallback to environment variables if not provided
			username = os.Getenv("KAFKA_USERNAME")
			password = os.Getenv("KAFKA_PASSWORD")
		}

		if username != "" && password != "" {
			config.Net.SASL.User = username
			config.Net.SASL.Password = password
			config.Net.SASL.Handshake = true

			// Configure SASL mechanism
			switch strings.ToUpper(saslMechanism) {
			case "PLAIN":
				config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			case "SCRAM-SHA-256":
				config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
				config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{HashGeneratorFcn: scram.SHA256}
				}
			case "SCRAM-SHA-512":
				config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
				config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
					return &XDGSCRAMClient{HashGeneratorFcn: scram.SHA512}
				}
			default:
				// Default to PLAIN if invalid mechanism
				config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			}
		}
	}

	// Configure SSL verification
	if config.Net.TLS.Enable && !sslVerify {
		config.Net.TLS.Config = nil
	}

	return config
}

// GetSASLSSLConfig returns a configuration optimized for SASL_SSL with SCRAM-SHA-512
func GetSASLSSLConfig(username, password string) *sarama.Config {
	// Get security settings from environment variables with defaults
	securityProtocol := os.Getenv("KAFKA_SECURITY_PROTOCOL")
	if securityProtocol == "" {
		securityProtocol = "SASL_SSL"
	}

	saslMechanism := os.Getenv("KAFKA_SASL_MECHANISM")
	if saslMechanism == "" {
		saslMechanism = "SCRAM-SHA-512"
	}

	sslVerify := true
	if sslVerifyStr := os.Getenv("KAFKA_SSL_VERIFY"); sslVerifyStr != "" {
		sslVerify = sslVerifyStr == "true"
	}

	return GetProducerConfigWithSecurity(securityProtocol, saslMechanism, username, password, sslVerify)
}

// GetConsumerConfig returns the Kafka consumer configuration
func GetConsumerConfig() *sarama.Config {
	// Get security settings from environment variables with defaults
	securityProtocol := os.Getenv("KAFKA_SECURITY_PROTOCOL")
	if securityProtocol == "" {
		securityProtocol = "PLAINTEXT"
	}

	saslMechanism := os.Getenv("KAFKA_SASL_MECHANISM")
	if saslMechanism == "" {
		saslMechanism = "PLAIN"
	}

	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	sslVerify := true
	if sslVerifyStr := os.Getenv("KAFKA_SSL_VERIFY"); sslVerifyStr != "" {
		sslVerify = sslVerifyStr == "true"
	}

	// Create base config using the security-aware producer config
	config := GetProducerConfigWithSecurity(securityProtocol, saslMechanism, username, password, sslVerify)

	// Override with consumer-specific settings
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false // Disable auto-commit for manual commit

	return config
}

// XDGSCRAMClient implements sarama.SCRAMClient for SCRAM authentication
type XDGSCRAMClient struct {
	HashGeneratorFcn scram.HashGeneratorFcn
	client           *scram.Client
	conversation     *scram.ClientConversation
}

// Begin implements sarama.SCRAMClient.Begin
func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	var err error
	x.client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.conversation = x.client.NewConversation()
	return nil
}

// Step implements sarama.SCRAMClient.Step
func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	return x.conversation.Step(challenge)
}

// Done implements sarama.SCRAMClient.Done
func (x *XDGSCRAMClient) Done() bool {
	return x.conversation.Done()
}
