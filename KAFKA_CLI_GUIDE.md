# Kafka CLI Guide cho Local Setup

Hướng dẫn sử dụng Kafka CLI tools với setup local của bạn.

## Cài đặt Kafka CLI Tools

### macOS
```bash
# Sử dụng Homebrew (đã cài đặt)
brew install kafka

# Hoặc Confluent CLI
brew install confluentinc/tap/confluent
```

### Linux
```bash
# Tải Kafka binary
wget https://downloads.apache.org/kafka/3.6.1/kafka_2.13-3.6.1.tgz
tar -xzf kafka_2.13-3.6.1.tgz
cd kafka_2.13-3.6.1
```

## Cấu hình

Setup của bạn sử dụng SASL/PLAIN authentication:
- **Brokers**: `localhost:9097,localhost:9098,localhost:9099`
- **Username**: `admin`
- **Password**: `admin-secret`

File cấu hình `kafka_client.properties` đã được tạo sẵn.

### Consumer Group IDs

Các consumer group IDs có thể được cấu hình qua environment variables:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `KAFKA_GROUP_ID` | `test-consumer-group` | Main consumer group ID |
| `KAFKA_CHANGE_DATA_GROUP_ID` | `change-data-consumer-group` | ChangeData consumer group ID |
| `KAFKA_TOPIC_EMPTIER_GROUP_ID` | `topic-emptier` | Topic emptier consumer group ID |

## Sử dụng Script Helper

### Liệt kê topics
```bash
./kafka-cli.sh topics
```

### Tạo topic mới
```bash
./kafka-cli.sh create-topic my-new-topic
```

### Xem thông tin topic
```bash
./kafka-cli.sh describe-topic test-topic
```

### Gửi message (Producer)
```bash
./kafka-cli.sh producer test-topic
# Sau đó gõ message và nhấn Enter
```

### Nhận message (Consumer)
```bash
./kafka-cli.sh consumer test-topic
# Hiển thị tất cả messages từ đầu
```

### Liệt kê consumer groups
```bash
./kafka-cli.sh groups
```

### Xem thông tin consumer group
```bash
./kafka-cli.sh describe-group test-consumer-group
```

## Sử dụng Makefile

### Liệt kê topics
```bash
make kafka-topics
```

### Tạo topic mới
```bash
make kafka-create-topic TOPIC=my-topic
```

### Gửi message
```bash
make kafka-producer TOPIC=test-topic
```

### Nhận message
```bash
make kafka-consumer TOPIC=test-topic
```

### Liệt kê consumer groups
```bash
make kafka-groups
```

## Lệnh Kafka CLI trực tiếp

### Topics
```bash
# Liệt kê topics
kafka-topics --bootstrap-server localhost:9092 --command-config kafka_client.properties --list

# Tạo topic
kafka-topics --bootstrap-server localhost:9092 --command-config kafka_client.properties \
  --create --topic my-topic --partitions 3 --replication-factor 1

# Xem thông tin topic
kafka-topics --bootstrap-server localhost:9092 --command-config kafka_client.properties \
  --describe --topic test-topic

# Xóa topic
kafka-topics --bootstrap-server localhost:9092 --command-config kafka_client.properties \
  --delete --topic my-topic
```

### Producer
```bash
# Console producer
kafka-console-producer --bootstrap-server localhost:9092 --producer.config kafka_client.properties \
  --topic test-topic

# Producer với key
kafka-console-producer --bootstrap-server localhost:9092 --producer.config kafka_client.properties \
  --topic test-topic --property "parse.key=true" --property "key.separator=:"
```

### Consumer
```bash
# Console consumer
kafka-console-consumer --bootstrap-server localhost:9092 --consumer.config kafka_client.properties \
  --topic test-topic --from-beginning

# Consumer với group
kafka-console-consumer --bootstrap-server localhost:9092 --consumer.config kafka_client.properties \
  --topic test-topic --group my-group --from-beginning

# Consumer với key
kafka-console-consumer --bootstrap-server localhost:9092 --consumer.config kafka_client.properties \
  --topic test-topic --property "print.key=true" --property "key.separator=:" --from-beginning
```

### Consumer Groups
```bash
# Liệt kê groups
kafka-consumer-groups --bootstrap-server localhost:9092 --command-config kafka_client.properties --list

# Xem thông tin group
kafka-consumer-groups --bootstrap-server localhost:9092 --command-config kafka_client.properties \
  --describe --group test-consumer-group

# Reset offset
kafka-consumer-groups --bootstrap-server localhost:9092 --command-config kafka_client.properties \
  --group test-consumer-group --topic test-topic --reset-offsets --to-earliest --execute
```

## Ví dụ thực tế

### 1. Tạo topic và gửi/nhận message
```bash
# Tạo topic
./kafka-cli.sh create-topic demo-topic

# Terminal 1: Gửi message
./kafka-cli.sh producer demo-topic
# Gõ: Hello World
# Gõ: Test Message
# Ctrl+C để thoát

# Terminal 2: Nhận message
./kafka-cli.sh consumer demo-topic
# Sẽ thấy: Hello World, Test Message
```

### 2. Kiểm tra consumer group
```bash
# Chạy consumer với group
kafka-console-consumer --bootstrap-server localhost:9092 --consumer.config kafka_client.properties \
  --topic test-topic --group demo-group --from-beginning

# Kiểm tra group
./kafka-cli.sh groups
./kafka-cli.sh describe-group demo-group
```

### 3. Xem thông tin chi tiết topic
```bash
./kafka-cli.sh describe-topic test-topic
```

## Troubleshooting

### Lỗi kết nối
```bash
# Kiểm tra Kafka có đang chạy không
docker-compose ps

# Kiểm tra logs
docker-compose logs kafka
```

### Lỗi authentication
```bash
# Kiểm tra file cấu hình
cat kafka_client.properties

# Test kết nối đơn giản
kafka-topics --bootstrap-server localhost:9092 --command-config kafka_client.properties --list
```

### Lỗi topic không tồn tại
```bash
# Tạo topic trước
./kafka-cli.sh create-topic my-topic

# Hoặc để auto-create (đã bật trong docker-compose)
```

## Tips

1. **Sử dụng script helper**: `./kafka-cli.sh` để dễ dàng hơn
2. **Sử dụng Makefile**: `make kafka-*` để tích hợp với workflow
3. **Kafka UI**: Truy cập http://localhost:8080 để xem trực quan
4. **Auto-create topics**: Topics sẽ được tạo tự động khi gửi message đầu tiên
5. **Consumer groups**: Sử dụng để load balancing và fault tolerance 