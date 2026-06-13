# Gateway Service

> IoT Gateway - Entry point từ Edge Gateways lên Cloud

## 🎯 Chức Năng

- Kết nối và quản lý EMQX MQTT Broker cluster
- Xác thực thiết bị qua X.509 certificates và token
- MQTT topic routing và message validation
- Chuyển đổi MQTT messages sang Kafka events
- Quản lý device sessions và heartbeat monitoring
- Rate limiting cho device connections

## 📁 Cấu Trúc

```
gateway-service/
├── mqtt-broker-connector/   # Kết nối EMQX cluster
│   ├── emqx-client/         #   EMQX client management
│   ├── topic-router/        #   MQTT topic routing rules
│   └── bridge-config/       #   MQTT bridge configuration
├── device-auth/             # Xác thực thiết bị
│   ├── x509-validator/      #   X.509 certificate validation
│   ├── token-auth/          #   Token-based authentication
│   └── acl-engine/          #   Access control list cho MQTT topics
└── protocol-handler/        # Xử lý message protocol
    ├── message-validator/   #   Message schema validation
    ├── mqtt-kafka-bridge/   #   MQTT → Kafka message bridge
    └── topic-manager/       #   Quản lý MQTT topic hierarchy
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go
- **MQTT**: EMQX 5.x client libraries
- **Message Queue**: Kafka producer/consumer (Sarama)
- **Auth**: X.509, JWT tokens

## 📡 MQTT Topic Structure

```
devices/{device_id}/telemetry       # Dữ liệu đo lường từ máy móc
devices/{device_id}/status          # Trạng thái kết nối (online/offline)
devices/{device_id}/command/req     # Lệnh điều khiển từ cloud
devices/{device_id}/command/res     # Phản hồi lệnh từ máy móc
devices/{device_id}/alarm           # Cảnh báo từ thiết bị
devices/{device_id}/config          # Cấu hình thiết bị
```

## 🚀 Chạy

```bash
go run cmd/gateway/main.go --config configs/dev.yaml
```
