# Industrial IoT Platform

> Hệ thống giám sát và điều khiển nhà máy công nghiệp đa vùng (Multi-Region Industrial IoT Monitoring & Control Platform)

## 🏭 Tổng Quan

Industrial IoT Platform là nền tảng giám sát và điều khiển công nghiệp, cho phép:
- **Thu thập dữ liệu** từ máy móc/PLC qua các giao thức công nghiệp (Modbus, OPC-UA, MQTT)
- **Giám sát real-time** trạng thái máy móc, dây chuyền sản xuất
- **Điều khiển từ xa** máy móc từ Web Portal và Mobile App
- **Cảnh báo thông minh** khi có sự cố hoặc vượt ngưỡng
- **Phân tích dữ liệu** để dự đoán bảo trì, tối ưu sản xuất
- **Triển khai đa vùng** hỗ trợ nhà máy ở HN, HCM, Đà Nẵng

## 📐 Kiến Trúc Hệ Thống

```
                          ┌─────────────────────────────────┐
                          │         CLOUD (Kubernetes)       │
                          │                                  │
  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
  │  API     │  │  Auth    │  │ Command  │  │Telemetry  │  │
  │ Gateway  │  │ Service  │  │ Service  │  │ Service   │  │
  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬─────┘  │
       │              │             │               │        │
  ┌────┴──────────────┴─────────────┴───────────────┴────┐  │
  │                    Message Bus (Kafka)                 │  │
  └────┬──────────────┬──────────────────────────────────┘  │
       │              │                                      │
  ┌────┴────┐  ┌──────┴──────┐  ┌──────────┐ ┌──────────┐  │
  │  Alarm  │  │Notification │  │Analytics │ │ Device   │  │
  │ Service │  │  Service    │  │ Service  │ │ Service  │  │
  └─────────┘  └─────────────┘  └──────────┘ └──────────┘  │
       │              │             │               │        │
  ┌────┴──────────────┴─────────────┴───────────────┴────┐  │
  │                  Data Stores                          │  │
  │  [TimescaleDB] [PostgreSQL] [Redis] [MinIO]          │  │
  └──────────────────────────────────────────────────────┘  │
                          └─────────────────────────────────┘
                               ▲          ▲          ▲
                               │  MQTT    │  MQTT    │  MQTT
                               │          │          │
                    ┌──────────┴──┐ ┌─────┴─────┐ ┌─┴──────────┐
                    │ Edge HN     │ │ Edge HCM  │ │ Edge DN     │
                    │ ┌────────┐  │ │ ┌───────┐ │ │ ┌────────┐  │
                    │ │Protocol│  │ │ │Proto- │ │ │ │Protocol│  │
                    │ │Adapter │  │ │ │col    │ │ │ │Adapter │  │
                    │ │        │  │ │ │Adapter│ │ │ │        │  │
                    │ └───┬────┘  │ │ └───┬───┘ │ │ └───┬────┘  │
                    │     │Modbus │ │     │     │ │     │       │
                    └─────┼───────┘ └─────┼─────┘ └─────┼───────┘
                          │OPC-UA        │              │
                     [Máy móc HN]   [Máy móc HCM]  [Máy móc ĐN]
```

## 📁 Cấu Trúc Repository

```
industrial-iot-platform/
│
├── edge-gateway/                    # Edge Gateway - triển khai tại mỗi nhà máy
│   ├── protocol-adapter/            #   Modbus, OPC-UA, MQTT client adapter
│   ├── local-cache/                 #   Buffer data khi mất kết nối internet
│   ├── command-executor/            #   Thực thi lệnh điều khiển local
│   └── edge-rules-engine/           #   Rule engine cơ bản chạy tại edge
│
├── gateway-service/                 # IoT Gateway - entry point edge → cloud
│   ├── mqtt-broker-connector/       #   Kết nối EMQX cluster
│   ├── device-auth/                 #   Xác thực thiết bị (X.509 cert, token)
│   └── protocol-handler/            #   MQTT topic routing, message validation
│
├── device-service/                  # Quản lý vòng đời thiết bị
│   ├── device-registry/             #   CRUD thiết bị, metadata
│   ├── device-shadow/               #   Digital Twin (desired vs reported state)
│   └── firmware-management/         #   OTA firmware update
│
├── telemetry-service/               # Thu thập & lưu trữ dữ liệu đo lường
│   ├── data-ingestion/              #   Nhận & validate data từ Kafka
│   ├── time-series-storage/         #   Lưu trữ TimescaleDB/InfluxDB
│   └── data-query-api/              #   Query API (REST/GraphQL)
│
├── command-service/                 # Điều khiển cloud → máy móc
│   ├── command-dispatcher/          #   Nhận lệnh từ UI, gửi qua MQTT
│   ├── command-tracker/             #   Theo dõi trạng thái lệnh
│   └── command-history/             #   Audit log lệnh điều khiển
│
├── alarm-service/                   # Cảnh báo & Alert
│   ├── rule-engine/                 #   Đánh giá điều kiện alarm
│   ├── alarm-state-machine/         #   Lifecycle: raised → acked → resolved
│   └── escalation-policy/           #   Escalation rule
│
├── notification-service/            # Gửi thông báo đa kênh
│   ├── email-notifier/              #   Email notification
│   ├── sms-notifier/                #   SMS notification
│   ├── push-notifier/               #   Mobile push notification
│   └── webhook-notifier/            #   Webhook tích hợp hệ thống khác
│
├── analytics-service/               # Phân tích dữ liệu
│   ├── real-time-analytics/         #   Stream processing (Kafka Streams)
│   ├── batch-analytics/             #   Batch processing
│   └── ml-pipeline/                 #   Predictive maintenance models
│
├── auth-service/                    # Xác thực & Phân quyền
│   ├── identity-provider/           #   Keycloak integration
│   ├── rbac/                        #   Role-based access control
│   └── audit-log/                   #   User activity audit trail
│
├── api-gateway/                     # API Gateway / BFF
│   ├── web-bff/                     #   Backend-for-Frontend cho Web Portal
│   ├── mobile-bff/                  #   Backend-for-Frontend cho Mobile App
│   └── rate-limiter/                #   Rate limiting & circuit breaker
│
├── web-portal/                      # Web Portal (React)
│   ├── dashboard/                   #   Real-time dashboard
│   ├── device-control/              #   Giao diện điều khiển máy móc
│   ├── alarm-center/                #   Trung tâm cảnh báo
│   ├── reports/                     #   Báo cáo & analytics
│   └── admin/                       #   Quản trị hệ thống
│
├── mobile-app/                      # Mobile App (Flutter)
│   ├── monitoring/                  #   Dashboard trên mobile
│   ├── alarm-push/                  #   Push notification alarm
│   └── quick-control/               #   Điều khiển nhanh
│
├── deployment/                      # Infrastructure as Code
│   ├── kubernetes/                  #   Helm charts & K8s manifests
│   ├── terraform/                   #   Cloud resource provisioning
│   ├── docker-compose/              #   Local dev environment
│   └── edge-deployment/             #   Ansible playbooks cho edge
│
├── Makefile                         # Build, test, deploy commands
├── docker-compose.yml              # Root compose (references deployment/)
└── README.md                        # ← This file
```

## 🚀 Bắt Đầu Nhanh

### Yêu cầu
- Docker & Docker Compose
- Go 1.22+ / Node.js 20+ / Java 21 (tùy service)
- Kubernetes cluster (cho production)

### Khởi động môi trường dev

```bash
# 1. Clone repository
git clone <repo-url> industrial-iot-platform
cd industrial-iot-platform

# 2. Khởi động infrastructure
make dev

# 3. Kiểm tra trạng thái
docker compose -f deployment/docker-compose/dev.yml ps

# 4. Truy cập các service:
#   - EMQX Dashboard:    http://localhost:18083 (admin/admin123)
#   - Grafana:            http://localhost:3000  (admin/admin123)
#   - Keycloak:           http://localhost:8080  (admin/admin123)
#   - Prometheus:         http://localhost:9090
```

### Phát triển từng service

Mỗi service có Makefile và README riêng với hướng dẫn chi tiết:

```bash
cd telemetry-service
make dev    # Chạy service ở chế độ development
make test   # Chạy unit tests
make build  # Build binary/Docker image
```

## 🔄 Luồng Dữ Liệu Chính

### 1. Telemetry Flow (Máy móc → Cloud → UI)
```
Machine → Edge Gateway (MQTT) → EMQX → Kafka → Telemetry Service → TimescaleDB
                                                                    ↓
UI Dashboard ← Web BFF ← Data Query API ←───────────────────────────┘
```

### 2. Command Flow (UI → Cloud → Máy móc)
```
UI → Web BFF → Command Service → Kafka → Gateway Service → MQTT → Edge Gateway → Machine
                    ↓                                                    ↓
              Command History DB                              MQTT Response → UI Update
```

### 3. Alarm Flow
```
Telemetry Data → Rule Engine → Alarm State Machine → Notification Service → User
                              ↓
                         Alarm DB → UI Alarm Center
```

## 🛠 Công Nghệ Sử Dụng

| Layer | Công Nghệ | Mục Đích |
|-------|-----------|----------|
| **Edge** | Go | Gateway nhẹ, chạy trên ARM/Linux |
| **MQTT** | EMQX 5.x | MQTT Broker, 1M+ concurrent connections |
| **Message Bus** | Apache Kafka | Event streaming, inter-service communication |
| **API Gateway** | Kong / Traefik | Rate limiting, auth, routing |
| **Services** | Go / Java Spring Boot | Microservices backend |
| **Time-Series DB** | TimescaleDB (PostgreSQL) | Lưu trữ telemetry data |
| **Main DB** | PostgreSQL 16 | Metadata, user, alarm, device registry |
| **Cache** | Redis Cluster | Session, device state, real-time data |
| **Auth** | Keycloak | OAuth2 / OIDC Identity Provider |
| **Web Portal** | React + TypeScript + MUI | Dashboard & control UI |
| **Mobile App** | Flutter | Cross-platform mobile app |
| **Observability** | Prometheus + Grafana + Loki | Metrics, logs, tracing |
| **CI/CD** | GitHub Actions | Build, test, deploy pipeline |
| **Infra** | Kubernetes + Terraform | Container orchestration & IaC |

## 🔒 Bảo Mật

- **Device Auth**: X.509 certificates + MQTT TLS
- **API Security**: OAuth2/OIDC qua Keycloak
- **Network**: VPN giữa cloud và edge gateways
- **Audit**: Ghi log tất cả command và user action
- **RBAC**: Phân quyền chi tiết theo role (Operator, Supervisor, Admin)

## 📊 Multi-Region Deployment

```
Region: Hà Nội            Region: Hồ Chí Minh         Region: Đà Nẵng
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Edge Gateway x2  │     │ Edge Gateway x2  │     │ Edge Gateway x2  │
│ (Active/Standby) │     │ (Active/Standby) │     │ (Active/Standby) │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │ (MQTT over TLS)
                         ┌───────┴───────┐
                         │ Cloud (K8s)   │
                         │ EMQX Cluster  │
                         │ All Services  │
                         │ DB HA Setup   │
                         └───────────────┘
```

## 📝 Roadmap

- [x] Kiến trúc tổng thể
- [ ] Phase 1: MVP - Edge Gateway + Telemetry + Dashboard cơ bản
- [ ] Phase 2: Core - Device Management + Alarm + Command + Auth
- [ ] Phase 3: Advanced - Analytics + ML + Rule Engine + Notification
- [ ] Phase 4: Multi-Region - HA deployment, DR plan
- [ ] Phase 5: Mobile App + Push notification

## 📄 License

Proprietary - All rights reserved
