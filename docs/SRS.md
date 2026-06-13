# Software Requirements Specification (SRS)

## Industrial IoT Platform — Hệ Thống Giám Sát & Điều Khiển Nhà Máy Công Nghiệp Đa Vùng

| Thuộc tính | Giá trị |
|---|---|
| **Phiên bản** | 1.0 |
| **Ngày phát hành** | 13/06/2026 |
| **Tác giả** | Đội ngũ Industrial IoT Platform |
| **Trạng thái** | Approved |
| **Phân loại** | Internal — Restricted |

---

## Mục Lục

1. [Giới Thiệu](#1-giới-thiệu)
2. [Tổng Quan Hệ Thống](#2-tổng-quan-hệ-thống)
3. [Kiến Trúc Hệ Thống](#3-kiến-trúc-hệ-thống)
4. [Yêu Cầu Chức Năng Chi Tiết](#4-yêu-cầu-chức-năng-chi-tiết)
   - [4.1 Edge Gateway](#41-edge-gateway)
   - [4.2 Gateway Service](#42-gateway-service)
   - [4.3 Device Service](#43-device-service)
   - [4.4 Telemetry Service](#44-telemetry-service)
   - [4.5 Command Service](#45-command-service)
   - [4.6 Alarm Service](#46-alarm-service)
   - [4.7 Notification Service](#47-notification-service)
   - [4.8 Analytics Service](#48-analytics-service)
   - [4.9 Auth Service](#49-auth-service)
   - [4.10 API Gateway](#410-api-gateway)
   - [4.11 Web Portal](#411-web-portal)
   - [4.12 Mobile App](#412-mobile-app)
5. [Giao Diện Hệ Thống](#5-giao-diện-hệ-thống)
6. [Mô Hình Dữ Liệu](#6-mô-hình-dữ-liệu)
7. [Yêu Cầu Phi Chức Năng](#7-yêu-cầu-phi-chức-năng)
8. [Yêu Cầu Bảo Mật](#8-yêu-cầu-bảo-mật)
9. [Yêu Cầu Vận Hành](#9-yêu-cầu-vận-hành)
10. [Ràng Buộc Thiết Kế](#10-ràng-buộc-thiết-kế)
11. [Phụ Lục](#11-phụ-lục)

---

## 1. Giới Thiệu

### 1.1 Mục Đích

Tài liệu này đặc tả các yêu cầu phần mềm chi tiết cho **Industrial IoT Platform** — nền tảng giám sát và điều khiển công nghiệp đa vùng. Tài liệu phục vụ cho:

- Đội ngũ phát triển: hiểu rõ yêu cầu kỹ thuật để triển khai
- Đội ngũ kiểm thử: xây dựng test case dựa trên đặc tả
- Đội ngũ vận hành: hiểu rõ hành vi hệ thống để vận hành và troubleshooting

### 1.2 Phạm Vi

Tài liệu bao phủ toàn bộ 12 component của hệ thống:

1. Edge Gateway (edge-gateway)
2. Gateway Service (gateway-service)
3. Device Service (device-service)
4. Telemetry Service (telemetry-service)
5. Command Service (command-service)
6. Alarm Service (alarm-service)
7. Notification Service (notification-service)
8. Analytics Service (analytics-service)
9. Auth Service (auth-service)
10. API Gateway (api-gateway)
11. Web Portal (web-portal)
12. Mobile App (mobile-app)

### 1.3 Tài Liệu Tham Chiếu

| Ref | Tài liệu |
|-----|---------|
| [REF-01] | IEEE 830-1998 — Recommended Practice for SRS |
| [REF-02] | [Business Requirements Document (BRD)](BRD.md) |
| [REF-03] | [User Requirements Document (URD)](URD.md) |
| [REF-04] | [README.md](../README.md) — Kiến trúc tổng thể |
| [REF-05] | Service READMEs trong repository |

---

## 2. Tổng Quan Hệ Thống

### 2.1 Mô Tả Sản Phẩm

Industrial IoT Platform là nền tảng microservices được triển khai trên Kubernetes, phục vụ giám sát và điều khiển nhà máy công nghiệp. Hệ thống kết nối với máy móc/PLC tại các nhà máy thông qua Edge Gateway, thu thập dữ liệu real-time qua MQTT, xử lý qua Apache Kafka, lưu trữ trong TimescaleDB/PostgreSQL, và hiển thị trên Web Portal và Mobile App.

### 2.2 Kiến Trúc Tổng Thể

Xem chi tiết tại mục [3. Kiến Trúc Hệ Thống](#3-kiến-trúc-hệ-thống) và [README.md](../README.md).

### 2.3 Công Nghệ Chính

| Layer | Công Nghệ | Phiên bản |
|-------|-----------|-----------|
| **Edge Runtime** | Go | 1.22+ |
| **MQTT Broker** | EMQX | 5.x |
| **Message Bus** | Apache Kafka | 3.7+ |
| **API Gateway** | Kong | 3.x |
| **Microservice Framework** | Go (net/http, Chi/Gin) / Java Spring Boot | Go 1.22 / Java 21 |
| **Time-Series DB** | TimescaleDB (PostgreSQL) | 2.15+ / 16 |
| **Relational DB** | PostgreSQL | 16 |
| **Cache** | Redis | 7.x |
| **Object Storage** | MinIO / S3-compatible | — |
| **Identity Provider** | Keycloak | 24+ |
| **Web Frontend** | React + TypeScript + MUI | 19 / 5.x / 6.x |
| **Mobile** | Flutter + Dart | 3.x |
| **Observability** | Prometheus + Grafana + Loki | — |
| **Orchestration** | Kubernetes | 1.29+ |
| **IaC** | Terraform + Ansible | — |

---

## 3. Kiến Trúc Hệ Thống

### 3.1 Kiến Trúc Vật Lý — Multi-Region

```
Region: Hà Nội                  Region: Hồ Chí Minh              Region: Đà Nẵng
┌──────────────────────┐       ┌──────────────────────┐        ┌──────────────────────┐
│  Edge Gateway        │       │  Edge Gateway        │        │  Edge Gateway        │
│  ┌────────────────┐  │       │  ┌────────────────┐  │        │  ┌────────────────┐  │
│  │ Protocol       │  │       │  │ Protocol       │  │        │  │ Protocol       │  │
│  │ Adapter        │  │       │  │ Adapter        │  │        │  │ Adapter        │  │
│  │ • Modbus RTU   │  │       │  │ • Modbus TCP   │  │        │  │ • OPC-UA       │  │
│  │ • Siemens S7   │  │       │  │ • OPC-UA       │  │        │  │ • Modbus TCP   │  │
│  ├────────────────┤  │       │  ├────────────────┤  │        │  ├────────────────┤  │
│  │ Local Cache    │  │       │  │ Local Cache    │  │        │  │ Local Cache    │  │
│  │ (SQLite)       │  │       │  │ (SQLite)       │  │        │  │ (SQLite)       │  │
│  ├────────────────┤  │       │  ├────────────────┤  │        │  ├────────────────┤  │
│  │ Command        │  │       │  │ Command        │  │        │  │ Command        │  │
│  │ Executor       │  │       │  │ Executor       │  │        │  │ Executor       │  │
│  ├────────────────┤  │       │  ├────────────────┤  │        │  ├────────────────┤  │
│  │ Edge Rules     │  │       │  │ Edge Rules     │  │        │  │ Edge Rules     │  │
│  │ Engine         │  │       │  │ Engine         │  │        │  │ Engine         │  │
│  └───────┬────────┘  │       │  └───────┬────────┘  │        │  └───────┬────────┘  │
│          │ MQTT/TLS  │       │          │ MQTT/TLS  │        │          │ MQTT/TLS  │
└──────────┼───────────┘       └──────────┼───────────┘        └──────────┼───────────┘
           │                              │                               │
           └──────────────────────────────┼───────────────────────────────┘
                                          │
                         ┌────────────────┴────────────────┐
                         │        CLOUD (Kubernetes)        │
                         │                                  │
                         │  ┌──────────────────────────┐   │
                         │  │   EMQX Cluster (MQTT)     │   │
                         │  │   • X.509 Auth            │   │
                         │  │   • TLS Termination       │   │
                         │  │   • Topic Routing         │   │
                         │  └──────────┬───────────────┘   │
                         │             │                    │
                         │  ┌──────────▼───────────────┐   │
                         │  │   Gateway Service         │   │
                         │  │   • MQTT→Kafka Bridge     │   │
                         │  │   • Message Validation    │   │
                         │  └──────────┬───────────────┘   │
                         │             │                    │
                         │  ┌──────────▼───────────────┐   │
                         │  │     Apache Kafka          │   │
                         │  │  ┌──────┐ ┌──────┐       │   │
                         │  │  │Tele- │ │Cmd   │ ...   │   │
                         │  │  │metry │ │Events│       │   │
                         │  │  └──┬───┘ └──┬───┘       │   │
                         │  └─────┼────────┼───────────┘   │
                         │        │        │                │
                         │  ┌─────▼──┐ ┌───▼────┐          │
                         │  │Telemetry│ │Command │  ...    │
                         │  │Service  │ │Service │          │
                         │  └────┬────┘ └───┬────┘          │
                         │       │          │                │
                         │  ┌────▼──────────▼─────┐         │
                         │  │    Data Stores       │         │
                         │  │  • TimescaleDB       │         │
                         │  │  • PostgreSQL        │         │
                         │  │  • Redis             │         │
                         │  │  • MinIO             │         │
                         │  └──────────────────────┘         │
                         │                                  │
                         │  ┌──────────────────────────┐   │
                         │  │   API Gateway (Kong)      │   │
                         │  │   • Rate Limiting         │   │
                         │  │   • Auth Plugin           │   │
                         │  │   • Web BFF / Mobile BFF  │   │
                         │  └──────────┬───────────────┘   │
                         │             │                    │
                         │  ┌──────────▼───────────────┐   │
                         │  │   Web Portal (React)      │   │
                         │  │   Mobile App (Flutter)    │   │
                         │  └──────────────────────────┘   │
                         │                                  │
                         │  ┌──────────────────────────┐   │
                         │  │   Observability Stack      │   │
                         │  │   Prometheus + Grafana     │   │
                         │  │   + Loki (logs)            │   │
                         │  └──────────────────────────┘   │
                         └──────────────────────────────────┘
```

### 3.2 Kiến Trúc Logic — Microservices

```
┌─────────────┐  ┌──────────┐  ┌─────────┐  ┌───────────┐
│   Edge      │  │ Gateway  │  │  Device  │  │Telemetry  │
│  Gateway    │  │ Service  │  │ Service  │  │ Service   │
│             │  │          │  │          │  │           │
│ Go          │  │ Go       │  │ Go       │  │ Go        │
│ Linux/ARM   │  │ K8s      │  │ K8s      │  │ K8s       │
│ SQLite      │  │          │  │ PG       │  │ TimescaleDB│
└──────┬──────┘  └────┬─────┘  └────┬─────┘  └─────┬─────┘
       │              │             │               │
       │    MQTT      │             │               │
       └──────────────┘             │               │
              │                     │               │
       ┌──────┴─────────────────────┴───────────────┴──────┐
       │                  Apache Kafka                     │
       │  Topics: telemetry, commands, alarms, events, ... │
       └──────┬────────────────────────────────────────────┘
              │
┌─────────────┼──────────────┬──────────────┬──────────────┐
│             │              │              │              │
│  ┌──────────┴──┐  ┌────────┴───┐  ┌───────┴───┐  ┌──────┴─────┐
│  │  Command    │  │   Alarm    │  │Notifi-    │  │ Analytics  │
│  │  Service    │  │  Service   │  │cation Svc │  │  Service   │
│  │             │  │            │  │           │  │            │
│  │ Go          │  │ Go/Java    │  │ Go        │  │ Java/Python│
│  │ PG + Redis  │  │ PG + Redis │  │ PG        │  │ Flink/Spark│
│  └──────┬──────┘  └─────┬──────┘  └─────┬─────┘  └──────┬─────┘
│         │               │               │               │
│  ┌──────┴───────────────┴───────────────┴───────────────┴─────┐
│  │                      Data Stores                           │
│  │  TimescaleDB │ PostgreSQL │ Redis Cluster │ MinIO/S3       │
│  └─────────────────────────────────────────────────────────────┘
│                               │
│  ┌────────────────────────────┴───────────────────────────────┐
│  │              API Gateway (Kong) + BFF                      │
│  │          Web BFF (Node.js/Go) │ Mobile BFF (Node.js/Go)    │
│  └────────┬────────────────────────────────────┬──────────────┘
│           │                                    │
│  ┌────────┴────────┐                  ┌────────┴────────┐
│  │  Web Portal     │                  │  Mobile App     │
│  │  React + TS     │                  │  Flutter        │
│  │  Port 5173      │                  │  iOS + Android  │
│  └─────────────────┘                  └─────────────────┘
```

### 3.3 Luồng Dữ Liệu Chính

#### 3.3.1 Telemetry Flow (Máy → Cloud → UI)

```
┌─────────┐    Modbus/     ┌──────────┐    MQTT/TLS     ┌──────────┐
│  PLC /  │──────────────→│  Edge    │────────────────→│  EMQX    │
│  Máy móc│   OPC-UA/S7   │ Gateway  │   JSON payload  │ Cluster  │
└─────────┘               └──────────┘                 └────┬─────┘
                                                            │
                                                  ┌─────────▼───────┐
                                                  │ Gateway Service │
                                                  │ • Parse topic   │
                                                  │ • Validate msg  │
                                                  │ • Enrich tags   │
                                                  └─────────┬───────┘
                                                            │  Kafka: "telemetry" topic
                                                  ┌─────────▼───────┐
                                                  │    Kafka        │
                                                  └────┬───────┬────┘
                                                       │       │
                                          ┌────────────▼─┐ ┌───▼──────────┐
                                          │ Telemetry    │ │ Alarm        │
                                          │ Service      │ │ Service      │
                                          │ • Insert DB  │ │ • Evaluate   │
                                          │ • Cache Redis│ │   rules      │
                                          └──────┬───────┘ └──────────────┘
                                                 │
                                          ┌──────▼───────┐
                                          │ TimescaleDB  │
                                          │ + Redis Cache│
                                          └──────┬───────┘
                                                 │  WebSocket / REST / GraphQL
                                          ┌──────▼───────┐
                                          │  API Gateway │
                                          │  (Web BFF)   │
                                          └──────┬───────┘
                                                 │
                                          ┌──────▼───────┐
                                          │  Web Portal  │
                                          │  Dashboard   │
                                          └──────────────┘
```

#### 3.3.2 Command Flow (UI → Cloud → Máy)

```
┌──────────┐   REST/WS    ┌──────────┐   gRPC/REST   ┌──────────┐
│  Web     │─────────────→│ API      │──────────────→│ Command  │
│  Portal  │              │ Gateway  │               │ Service  │
└──────────┘              └──────────┘               └────┬─────┘
                                                          │
                                            ┌─────────────▼──────┐
                                            │ Command Dispatcher │
                                            │ • Validate quyền   │
                                            │ • Tạo command (DB) │
                                            │ • Gửi MQTT         │
                                            └────────┬───────────┘
                                                     │ MQTT topic: devices/{id}/command/req
                                            ┌────────▼───────────┐
                                            │ Gateway Service     │
                                            │ MQTT → Edge         │
                                            └────────┬───────────┘
                                                     │ MQTT
                                            ┌────────▼───────────┐
                                            │ Edge Gateway        │
                                            │ • Nhận lệnh         │
                                            │ • Gửi xuống PLC     │
                                            │ • Phản hồi kết quả  │
                                            └────────┬───────────┘
                                                     │ Modbus/OPC-UA
                                            ┌────────▼───────────┐
                                            │ PLC / Máy móc       │
                                            │ Thực thi lệnh       │
                                            └────────────────────┘

Phản hồi: Máy → Edge → MQTT(res) → Gateway → Kafka → Command Svc → UI update
```

#### 3.3.3 Alarm Flow

```
Telemetry Data (Kafka)
        │
┌───────▼──────────┐
│ Alarm Service    │
│ Rule Engine      │
│ • Threshold      │──→ Condition met?
│ • Rate of Change │
│ • Anomaly        │
│ • State Change   │
└───────┬──────────┘
        │ Alarm raised
┌───────▼──────────┐
│ Alarm State      │
│ Machine          │
│ RAISED → ACKED   │
│   → RESOLVED     │
└───────┬──────────┘
        │
┌───────▼──────────┐     ┌──────────────┐
│ Notification     │────→│ Email (SMTP) │
│ Service          │────→│ SMS (Twilio) │
│                  │────→│ Push (FCM)   │
│                  │────→│ Webhook      │
└───────┬──────────┘     └──────────────┘
        │
┌───────▼──────────┐
│ UI Update        │
│ (WebSocket)      │
└──────────────────┘
```

---

## 4. Yêu Cầu Chức Năng Chi Tiết

### 4.1 Edge Gateway

**Định danh**: `SRV-EDGE` | **Ngôn ngữ**: Go | **Mục tiêu**: ARM/Linux, nhị phân đơn

#### 4.1.1 Protocol Adapter

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| EG-001 | Modbus RTU Support | Kết nối serial RS-485, hỗ trợ: read coils, read holding registers, read input registers. Cấu hình: baud rate, parity, stop bits, slave ID | P0 — MVP |
| EG-002 | Modbus TCP Support | Kết nối TCP port 502, hỗ trợ tất cả function codes phổ biến. Cấu hình: IP:Port, Unit ID | P0 — MVP |
| EG-003 | OPC-UA Support | OPC-UA client kết nối đến OPC-UA server. Hỗ trợ: browse nodes, subscribe data changes, read/write. Security: None, Sign, Sign&Encrypt | P1 — Core |
| EG-004 | Siemens S7 Support | Kết nối Siemens S7-300/400/1200/1500. Đọc/ghi: Data Blocks, Inputs, Outputs, Merkers, Timers, Counters | P2 — Advanced |
| EG-005 | MQTT Publish | Định dạng topic: `factory/{factory_id}/{area}/{machine_id}/{metric}`. Payload: JSON. QoS: 1 (at least once). Giữ kết nối MQTT persistent | P0 — MVP |
| EG-006 | Polling Configuration | Cấu hình polling interval cho từng metric (mặc định 1s, tối thiểu 100ms, tối đa 3600s) | P0 — MVP |

#### 4.1.2 Local Cache & Store-Forward

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| EG-010 | SQLite Local Buffer | Lưu telemetry data vào SQLite khi publish MQTT thất bại. Schema: time, topic, payload, retry_count, created_at | P0 — MVP |
| EG-011 | Store-and-Forward | Tự động gửi lại data từ buffer khi MQTT reconnect. Gửi theo thứ tự thời gian (FIFO). Xóa record sau khi gửi thành công | P0 — MVP |
| EG-012 | Buffer Size Limit | Giới hạn buffer tối đa: 10GB hoặc 7 ngày (configurable). Khi vượt → xóa record cũ nhất | P1 — Core |
| EG-013 | Disk Usage Protection | Dừng buffer khi disk sử dụng > 90%. Cảnh báo khi > 80% | P1 — Core |

#### 4.1.3 Command Executor

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| EG-020 | Nhận lệnh từ Cloud | Subscribe MQTT topic: `devices/{device_id}/command/req`. Parse lệnh, validate tham số | P1 — Core |
| EG-021 | Thực thi lệnh | Gửi lệnh xuống PLC qua protocol tương ứng (Modbus write, OPC-UA write). Timeout: cấu hình được, mặc định 10s | P1 — Core |
| EG-022 | Phản hồi kết quả | Publish kết quả lên topic: `devices/{device_id}/command/res`. Payload: {command_id, status, result, error, timestamp} | P1 — Core |
| EG-023 | Emergency Stop Local | Nhận emergency stop từ cloud HOẶC từ physical button. Gửi lệnh STOP khẩn cấp, không delay | P1 — Core |

#### 4.1.4 Edge Rules Engine

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| EG-030 | Rule cơ bản | Hỗ trợ rule dạng threshold: `if metric > X for Y seconds → alarm`. Tối đa 100 rule/edge | P2 — Advanced |
| EG-031 | Hành động Local | Khi rule kích hoạt: publish alarm MQTT, ghi log, kích hoạt relay output (nếu cấu hình) | P2 — Advanced |
| EG-032 | Hoạt động Offline | Rule engine vẫn hoạt động khi mất kết nối cloud. Sync alarm khi reconnect | P2 — Advanced |

#### 4.1.5 Health & Monitoring

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| EG-040 | Heartbeat | Gửi heartbeat MQTT mỗi 10s lên topic: `devices/{device_id}/status`. Payload: {status: "online", cpu, mem, disk, uptime, version} |
| EG-041 | Health Endpoint | HTTP endpoint GET /health (local network only): trả về trạng thái tất cả connections |
| EG-042 | Auto-Reconnect | MQTT auto-reconnect với exponential backoff: 1s, 2s, 4s, 8s, 16s, 30s (max) |

---

### 4.2 Gateway Service

**Định danh**: `SRV-GATEWAY` | **Ngôn ngữ**: Go | **Mục tiêu**: K8s Deployment

#### 4.2.1 MQTT Broker Connector

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| GW-001 | EMQX Connection | Kết nối đến EMQX cluster. Hỗ trợ multiple brokers (cluster). Auto-reconnect | P0 — MVP |
| GW-002 | Topic Routing | Subscribe MQTT topics theo pattern: `factory/+/+/+/+`. Route message đến Kafka topic tương ứng | P0 — MVP |
| GW-003 | MQTT→Kafka Bridge | Parse MQTT topic thành metadata (factory_id, area, machine_id, metric). Đóng gói thành Kafka message với key = machine_id | P0 — MVP |
| GW-004 | Message Validation | Validate MQTT payload JSON schema. Từ chối message không hợp lệ, ghi log, publish error topic | P1 — Core |
| GW-005 | QoS Management | MQTT QoS 1 → Kafka acks=1. Đảm bảo ít nhất once delivery | P0 — MVP |

#### 4.2.2 Device Authentication

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| GW-010 | X.509 Certificate Auth | Xác thực thiết bị qua X.509 client certificates. Validate: certificate chain, expiration, revocation (CRL/OCSP) | P1 — Core |
| GW-011 | Token Auth | Hỗ trợ JWT token auth cho thiết bị không có X.509. Token được cấp bởi Auth Service | P1 — Core |
| GW-012 | ACL Engine | Access control list cho MQTT topics: device chỉ được publish/subscribe topic của chính nó | P1 — Core |
| GW-013 | Rate Limiting | Rate limit theo device: tối đa N msg/s (cấu hình được, mặc định 100 msg/s) | P2 — Advanced |

#### 4.2.3 MQTT Topic Structure

```
devices/{device_id}/telemetry       # Dữ liệu đo lường (publish)
devices/{device_id}/status          # Trạng thái kết nối (publish)
devices/{device_id}/command/req     # Lệnh từ cloud → device (subscribe)
devices/{device_id}/command/res     # Kết quả lệnh (publish)
devices/{device_id}/alarm           # Cảnh báo từ device (publish)
devices/{device_id}/config          # Cấu hình OTA (subscribe)
devices/{device_id}/firmware        # Firmware update (subscribe)

factory/{factory_id}/{area}/{machine_id}/{metric}  # Telemetry (publish)
```

#### 4.2.4 Kafka Topics Produced

| Kafka Topic | Nguồn MQTT | Key | Payload |
|-------------|-----------|-----|---------|
| `telemetry.raw` | `factory/+/+/+/+` hoặc `devices/+/telemetry` | `machine_id` | TelemetryMessage JSON |
| `device.status` | `devices/+/status` | `device_id` | StatusMessage JSON |
| `device.alarms` | `devices/+/alarm` | `device_id` | AlarmMessage JSON |
| `command.responses` | `devices/+/command/res` | `command_id` | CommandResult JSON |

---

### 4.3 Device Service

**Định danh**: `SRV-DEVICE` | **Ngôn ngữ**: Go | **Database**: PostgreSQL + Redis

#### 4.3.1 Device Registry

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| DV-001 | CRUD Thiết bị | REST API: POST/GET/PUT/DELETE /api/v1/devices. Validate serial number unique. Soft delete | P1 — Core |
| DV-002 | Tìm kiếm & Lọc | GET /api/v1/devices?factory_id=X&area=Y&status=Z&search=keyword. Phân trang, sắp xếp | P1 — Core |
| DV-003 | Device Grouping | Phân cấp: Factory → Area → Line → Device. API lấy cây thiết bị | P1 — Core |
| DV-004 | Device Metadata | Schema mở rộng JSONB cho metadata tùy chỉnh (nhà sản xuất, năm sản xuất, thông số kỹ thuật) | P2 — Advanced |
| DV-005 | Device Status Tracking | Cập nhật status từ Kafka `device.status` topic: online/offline/maintenance/error. Last seen timestamp | P1 — Core |

#### 4.3.2 Device Shadow (Digital Twin)

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| DV-010 | Desired State | API GET/PUT /api/v1/devices/{id}/shadow/desired. User/cloud cập nhật trạng thái mong muốn | P1 — Core |
| DV-011 | Reported State | Nhận từ Kafka (MQTT topic devices/{id}/telemetry). Cập nhật reported state trong Redis + DB | P1 — Core |
| DV-012 | Delta Detection | So sánh desired vs reported. Nếu khác biệt > ngưỡng → tạo sự kiện delta. Hiển thị trong UI | P1 — Core |
| DV-013 | Shadow Sync | Khi device reconnect → tự động đồng bộ desired state xuống device | P1 — Core |

#### 4.3.3 Firmware Management

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| DV-020 | Firmware Upload | Upload firmware binary (tối đa 500MB) lên MinIO. Metadata: version, checksum, release notes | P2 — Advanced |
| DV-021 | OTA Update Dispatch | Chọn firmware + device(s) → tạo update job. Gửi firmware URL qua MQTT | P2 — Advanced |
| DV-022 | Update Status Tracking | Theo dõi trạng thái: pending → downloading → installing → success/failed | P2 — Advanced |
| DV-023 | Automatic Rollback | Nếu cập nhật fail → tự động gửi firmware cũ. Yêu cầu device hỗ trợ dual partition | P3 — Scale |

---

### 4.4 Telemetry Service

**Định danh**: `SRV-TELEMETRY` | **Ngôn ngữ**: Go | **Database**: TimescaleDB + Redis

#### 4.4.1 Data Ingestion

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| TM-001 | Kafka Consumer | Consume từ topic `telemetry.raw` với consumer group. Hỗ trợ multiple consumers (horizontal scaling) | P0 — MVP |
| TM-002 | Schema Validation | Validate TelemetryMessage JSON schema: required fields (factory_id, machine_id, metric, timestamp), types | P0 — MVP |
| TM-003 | Data Enrichment | Enrich message với device metadata từ cache: device name, model, location | P1 — Core |
| TM-004 | Data Transform | Chuẩn hóa unit (metric). Convert timestamp timezone → UTC. Xử lý duplicate (dedup theo time+device+metric) | P1 — Core |
| TM-005 | Dead Letter Queue | Message không parse được → gửi vào Kafka topic `telemetry.dlq` để debug sau | P2 — Advanced |

#### 4.4.2 Time-Series Storage

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| TM-010 | Hypertable Insert | Insert vào TimescaleDB hypertable `telemetry`. Batch insert (100-1000 records/batch) để tối ưu | P0 — MVP |
| TM-011 | Latest Value Cache | Mỗi lần insert → upsert Redis cache key `device:{id}:latest:{metric}` với TTL 24h | P0 — MVP |
| TM-012 | Compression Policy | Tự động compress data > 7 ngày: chuyển sang compressed chunks, giảm 90% storage | P1 — Core |
| TM-013 | Retention Policy | Tự động xóa data > 2 năm. Có thể cấu hình cho từng loại metric | P1 — Core |
| TM-014 | Down-sampling | Tự động tạo aggregate: 1m → 5m, 1h cho data > 30 ngày | P2 — Advanced |

#### 4.4.3 Data Query API

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| TM-020 | REST Query API | `GET /api/v1/telemetry?device_id=X&metric=Y&from=Z&to=W`. Phân trang (cursor-based). Tối đa 10K data points/request | P0 — MVP |
| TM-021 | Latest Value API | `GET /api/v1/telemetry/latest?factory_id=X`. Trả về latest value của tất cả metrics (từ Redis cache) | P0 — MVP |
| TM-022 | Aggregate API | `GET /api/v1/telemetry/aggregate?device_id=X&metric=Y&from=Z&to=W&interval=1h&fn=avg,max,min`. Sử dụng TimescaleDB time_bucket | P1 — Core |
| TM-023 | GraphQL API | GraphQL endpoint cho phép query linh hoạt: chọn fields, nested resources (device info), multiple metrics | P2 — Advanced |
| TM-024 | Real-time Stream | WebSocket/SSE endpoint: `GET /api/v1/telemetry/stream?device_ids=X,Y,Z`. Push real-time data từ Kafka đến client | P1 — Core |
| TM-025 | Export API | `POST /api/v1/telemetry/export`. Export CSV/JSON với filter, time range. Giới hạn 1M records | P2 — Advanced |

#### 4.4.4 TimescaleDB Schema

```sql
-- Hypertable chính
CREATE TABLE telemetry (
    time         TIMESTAMPTZ NOT NULL,
    device_id    UUID NOT NULL,
    factory_id   VARCHAR(10) NOT NULL,
    area         VARCHAR(50),
    metric_name  VARCHAR(100) NOT NULL,
    value        DOUBLE PRECISION,
    unit         VARCHAR(20),
    quality      SMALLINT DEFAULT 0,
    tags         JSONB,
    raw_payload  JSONB
);

SELECT create_hypertable('telemetry', 'time');
CREATE INDEX idx_telemetry_device_time ON telemetry (device_id, time DESC);
CREATE INDEX idx_telemetry_factory_time ON telemetry (factory_id, time DESC);
CREATE INDEX idx_telemetry_metric_name ON telemetry (metric_name, time DESC);

-- Latest value table (cho real-time dashboard)
CREATE TABLE device_latest_value (
    device_id    UUID NOT NULL,
    metric_name  VARCHAR(100) NOT NULL,
    factory_id   VARCHAR(10),
    area         VARCHAR(50),
    value        DOUBLE PRECISION,
    unit         VARCHAR(20),
    quality      SMALLINT DEFAULT 0,
    tags         JSONB,
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (device_id, metric_name)
);

-- Compression
SELECT add_compression_policy('telemetry', INTERVAL '7 days');

-- Retention
SELECT add_retention_policy('telemetry', INTERVAL '2 years');
```

---

### 4.5 Command Service

**Định danh**: `SRV-COMMAND` | **Ngôn ngữ**: Go | **Database**: PostgreSQL + Redis

#### 4.5.1 Command Dispatcher

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| CM-001 | Tạo Command | REST API: `POST /api/v1/commands`. Body: {device_id, command_type, parameters, priority}. Trả về command_id | P1 — Core |
| CM-002 | Command Validation | Kiểm tra: user có quyền điều khiển máy? Máy online? Tham số hợp lệ? Command supported? | P1 — Core |
| CM-003 | MQTT Dispatch | Gửi command qua Kafka → Gateway Service → MQTT topic `devices/{id}/command/req`. QoS 1 | P1 — Core |
| CM-004 | Command Priority | Hỗ trợ: LOW, NORMAL, HIGH, EMERGENCY. EMERGENCY được dispatch trước tất cả | P1 — Core |
| CM-005 | Batch Commands | `POST /api/v1/commands/batch`. Gửi một lệnh đến nhiều máy. Trả về danh sách command_id | P2 — Advanced |

#### 4.5.2 Command State Machine

```
Trạng thái: PENDING → SENT → EXECUTING → DONE / FAILED / TIMEOUT

PENDING:   Command được tạo, chưa gửi
SENT:      Đã gửi MQTT, chờ edge nhận
EXECUTING: Edge báo đang thực thi
DONE:      Thực thi thành công
FAILED:    Thực thi thất bại (lỗi từ edge)
TIMEOUT:   Không nhận được phản hồi sau X giây
```

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| CM-010 | State Transitions | Cập nhật trạng thái command dựa trên Kafka `command.responses`. Valid transitions: PENDING→SENT, SENT→EXECUTING, EXECUTING→DONE/FAILED, ANY→TIMEOUT | P1 — Core |
| CM-011 | Timeout Watchdog | Goroutine định kỳ kiểm tra command quá hạn: SENT > 30s, EXECUTING > 60s → TIMEOUT | P1 — Core |
| CM-012 | Retry Manager | Khi TIMEOUT: tự động retry tối đa 3 lần với exponential backoff (1s, 2s, 4s) | P1 — Core |
| CM-013 | Real-time Status | WebSocket push trạng thái command đến UI. Client subscribe: `ws://.../commands?command_id=X` | P2 — Advanced |

#### 4.5.3 Command History & Audit

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| CM-020 | Command History | `GET /api/v1/commands?device_id=X&from=Y&to=Z`. Phân trang, lọc theo status, type, user | P1 — Core |
| CM-021 | Audit Log | Mọi command đều được ghi: who, what, when, target device, result. Immutable — không thể sửa/xóa | P1 — Core |
| CM-022 | Compliance Export | Export audit log theo filter. Định dạng CSV. Hỗ trợ compliance: ISO 27001, Nghị định 53 | P2 — Advanced |

#### 4.5.4 API

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/commands | Tạo command mới |
| POST | /api/v1/commands/batch | Tạo batch command |
| GET | /api/v1/commands | Lấy danh sách command |
| GET | /api/v1/commands/{id} | Lấy command detail |
| GET | /api/v1/commands/{id}/status | Lấy status real-time |
| POST | /api/v1/commands/{id}/cancel | Hủy command (nếu còn PENDING) |
| GET | /api/v1/commands/audit | Export audit log |

---

### 4.6 Alarm Service

**Định danh**: `SRV-ALARM` | **Ngôn ngữ**: Go/Java | **Database**: PostgreSQL + Redis

#### 4.6.1 Rule Engine

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AL-001 | Threshold Rule | `if metric > X (or < X) for Y seconds → alarm`. X, Y, severity configurable | P1 — Core |
| AL-002 | Rate of Change Rule | `if delta(metric) > X over Y seconds → alarm`. Cho áp suất, nhiệt độ thay đổi đột ngột | P2 — Advanced |
| AL-003 | State Change Rule | `if machine.status: RUNNING → STOPPED (unexpected) → alarm` | P1 — Core |
| AL-004 | Anomaly Rule | `if metric outside 3-sigma band → alarm`. Training window 7d, auto-retrain định kỳ | P2 — Advanced |
| AL-005 | Window Aggregation | Hỗ trợ aggregate: AVG, MIN, MAX, SUM, COUNT trong sliding window | P2 — Advanced |
| AL-006 | Rule CRUD | REST API quản lý rule: POST/GET/PUT/DELETE /api/v1/alarms/rules. Rule lưu dạng JSON trong PostgreSQL | P1 — Core |
| AL-007 | Rule Hot Reload | Rule được cập nhật mà không cần restart service. Sử dụng Redis pub/sub để sync giữa các instance | P1 — Core |

#### 4.6.2 Alarm State Machine

```
            ┌──────────┐
            │  RAISED  │ ← Điều kiện rule = true
            └────┬─────┘
                 │ User acknowledge
            ┌────▼─────────┐
            │ ACKNOWLEDGED │
            └────┬─────────┘
                 │ Điều kiện hết (metric về bình thường)
            ┌────▼──────┐
            │ RESOLVED  │
            └───────────┘

     RAISED ──(timeout X phút)──► ESCALATED
     ACKNOWLEDGED ──(timeout Y phút)──► ESCALATED

     RAISED/ACKED ──(user snooze)──► SNOOZED ──(timeout)──► RAISED
     RAISED/ACKED ──(user suppress)──► SUPPRESSED ──(end schedule)──► resolved (nếu hết) / RAISED
```

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AL-010 | Alarm Creation | Khi rule kích hoạt: tạo alarm record, gán severity, publish Kafka `alarm.events` | P1 — Core |
| AL-011 | Alarm Dedup | Nếu alarm cùng máy + metric đã RAISED → không tạo mới, chỉ cập nhật timestamp | P1 — Core |
| AL-012 | Acknowledge | REST API: `POST /api/v1/alarms/{id}/ack`. Body: {user_id, comment}. Chuyển RAISED→ACKED | P1 — Core |
| AL-013 | Auto Resolve | Khi metric về bình thường trong X giây → tự động RESOLVED. Ghi thời gian resolved | P1 — Core |
| AL-014 | Snooze | REST API: `POST /api/v1/alarms/{id}/snooze`. Body: {duration_minutes}. Tạm ẩn alarm | P2 — Advanced |
| AL-015 | Suppression Schedule | CRUD schedule: chọn máy/mertric, chọn khoảng thời gian (vd: 8:00-10:00 thứ 7). Alarm không kích hoạt | P2 — Advanced |

#### 4.6.3 Escalation Policy

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AL-020 | Escalation Levels | Cấu hình level: Level 1 (Quản đốc) sau 5 phút, Level 2 (GĐ SX) sau 15 phút, Level 3 (Ban GĐ) sau 30 phút | P1 — Core |
| AL-021 | Escalation Actions | Mỗi level có hành động khác nhau: Push, Email, SMS, Voice call. Cấu hình được | P1 — Core |
| AL-022 | On-call Schedule | Tích hợp lịch trực: chọn người nhận escalation theo ngày/giờ | P2 — Advanced |

#### 4.6.4 Alarm Severity Levels

| Severity | Mô tả | Màu | Hành động mặc định |
|----------|-------|-----|-------------------|
| **EMERGENCY** | Nguy hiểm tức thời, cần dừng máy ngay | 🔴 Đỏ | Push + SMS + Email cho tất cả |
| **CRITICAL** | Sắp hỏng, cần xử lý gấp | 🟠 Cam | Push + Email cho Operator + Supervisor |
| **WARNING** | Vượt ngưỡng, cần chú ý | 🟡 Vàng | Push cho Operator |
| **INFO** | Thông báo thông thường | 🔵 Xanh | Không push, chỉ log |

#### 4.6.5 API

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /api/v1/alarms | Danh sách alarm (filter: status, severity, factory, time range) |
| GET | /api/v1/alarms/{id} | Alarm detail + timeline |
| POST | /api/v1/alarms/{id}/ack | Acknowledge alarm |
| POST | /api/v1/alarms/{id}/snooze | Snooze alarm |
| POST | /api/v1/alarms/{id}/suppress | Suppress alarm |
| GET | /api/v1/alarms/rules | Danh sách rule |
| POST | /api/v1/alarms/rules | Tạo rule mới |
| PUT | /api/v1/alarms/rules/{id} | Cập nhật rule |
| DELETE | /api/v1/alarms/rules/{id} | Xóa rule |
| POST | /api/v1/alarms/rules/{id}/test | Test rule với data mẫu |

---

### 4.7 Notification Service

**Định danh**: `SRV-NOTIFICATION` | **Ngôn ngữ**: Go | **Database**: PostgreSQL

#### 4.7.1 Multi-Channel Notifier

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| NF-001 | Email Notification | Gửi email qua SMTP. Hỗ trợ HTML template (Go template). File đính kèm (PDF report) | P1 — Core |
| NF-002 | SMS Notification | Gửi SMS qua Twilio API. Hỗ trợ VNPT/Viettel gateway. Unicode SMS | P2 — Advanced |
| NF-003 | Push Notification | Gửi push qua Firebase Cloud Messaging (Android) và APNs (iOS). Support: title, body, image, action buttons, badge | P2 — Advanced |
| NF-004 | Webhook Notification | HTTP POST đến URL cấu hình. Retry với exponential backoff (3 lần). Signature verification (HMAC) | P2 — Advanced |

#### 4.7.2 Template Engine

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| NF-010 | Template CRUD | Quản lý template cho từng kênh. API: GET/POST/PUT/DELETE /api/v1/notifications/templates | P1 — Core |
| NF-011 | Variable Substitution | Hỗ trợ biến: ${device_name}, ${metric_name}, ${value}, ${threshold}, ${severity}, ${timestamp}, ${factory_name} | P1 — Core |
| NF-012 | Conditional Rendering | Template hiển thị khác nhau theo severity: CRITICAL có thêm nút "Xem ngay", WARNING chỉ thông báo | P2 — Advanced |

#### 4.7.3 User Preferences

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| NF-020 | Channel Preferences | User chọn kênh nhận thông báo cho từng loại event. API: GET/PUT /api/v1/users/{id}/notification-prefs | P2 — Advanced |
| NF-021 | Quiet Hours | User cấu hình giờ yên lặng (vd: 22:00-06:00). Không gửi push/SMS trong thời gian này, chỉ email | P2 — Advanced |
| NF-022 | Delivery Status | Theo dõi trạng thái: PENDING → SENT → DELIVERED → READ / FAILED. API: GET /api/v1/notifications/{id}/status | P2 — Advanced |

#### 4.7.4 Rate Limiting

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NF-030 | Per-User Limit | Tối đa 10 push/user/giờ, 3 SMS/user/giờ, 20 email/user/ngày |
| NF-031 | Alarm Storm Protection | Khi > 50 alarm trong 1 phút → gộp thành 1 notification digest |

---

### 4.8 Analytics Service

**Định danh**: `SRV-ANALYTICS` | **Ngôn ngữ**: Java/Python | **Database**: TimescaleDB + PostgreSQL + MinIO

#### 4.8.1 Real-time Analytics

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AN-001 | OEE Calculator | Tính OEE real-time: Availability × Performance × Quality. Input: telemetry data, output: OEE value mỗi phút | P2 — Advanced |
| AN-002 | Throughput Monitor | Đếm sản phẩm/mỗi phút từ counter metric. So sánh với target | P2 — Advanced |
| AN-003 | Energy Monitor | Tính kWh từ power metric × thời gian. So sánh với baseline, cảnh báo nếu vượt | P2 — Advanced |
| AN-004 | KPI Aggregator | Aggregate KPI real-time: OEE, throughput, yield, energy. Push kết quả lên Kafka `analytics.kpi` để dashboard tiêu thụ | P2 — Advanced |

#### 4.8.2 Batch Analytics

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AN-010 | Daily Report | Job chạy 00:00 mỗi ngày: tính OEE, sản lượng, downtime, alarm count của ngày hôm trước. Lưu vào DB | P2 — Advanced |
| AN-011 | Trend Analysis | Phân tích xu hướng dài hạn (7d, 30d, 90d). Phát hiện xu hướng giảm OEE, tăng downtime | P2 — Advanced |
| AN-012 | Comparative Report | So sánh KPI giữa các nhà máy, ca sản xuất. Xếp hạng hiệu suất | P3 — Scale |

#### 4.8.3 ML Pipeline

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AN-020 | Model Training | Pipeline huấn luyện model: thu thập training data từ TimescaleDB, train (scikit-learn/TensorFlow), evaluate, save model artifact vào MinIO | P2 — Advanced |
| AN-021 | Model Serving | REST API: `POST /api/v1/ml/predict`. Input: device_id, metric, time_range. Output: prediction + confidence | P2 — Advanced |
| AN-022 | Anomaly Detection | ML model phát hiện bất thường: Isolation Forest, Autoencoder. Training window 7d, retrain mỗi tuần | P2 — Advanced |
| AN-023 | Predictive Maintenance | Dự đoán thời gian đến khi hỏng (Remaining Useful Life — RUL). Input: vibration, temperature, pressure trends | P2 — Advanced |
| AN-024 | MLflow Integration | Theo dõi experiment, model version, metrics. MLflow server chung | P3 — Scale |

#### 4.8.4 Export Engine

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AN-030 | CSV Export | Export dữ liệu CSV. Hỗ trợ filter, time range. Giới hạn 1M rows | P2 — Advanced |
| AN-031 | Excel Export | Export Excel với định dạng: style, chart embedded, multiple sheets. Sử dụng thư viện ExcelWriter | P2 — Advanced |
| AN-032 | PDF Report | Generate PDF từ HTML template. Bao gồm: biểu đồ (embedded PNG), bảng số liệu, header/footer | P2 — Advanced |

---

### 4.9 Auth Service

**Định danh**: `SRV-AUTH` | **Core**: Keycloak | **Database**: PostgreSQL

#### 4.9.1 Identity Provider

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AU-001 | Keycloak Integration | Realm: `industrial-iot`. Clients: web-portal, mobile-app, device-gateway. OAuth2 Authorization Code Flow + PKCE (web), Device Flow (edge) | P1 — Core |
| AU-002 | User Federation | Đồng bộ user từ LDAP/Active Directory. Tự động sync định kỳ 1 giờ | P2 — Advanced |
| AU-003 | SSO | Single Sign-On giữa Web Portal và các internal services. Session timeout: 8 giờ (1 ca làm việc) | P1 — Core |

#### 4.9.2 RBAC

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AU-010 | Role Definitions | 4 roles: Admin, Supervisor, Operator, Viewer. Xem chi tiết trong BRD | P1 — Core |
| AU-011 | Scope/Tenant | Mỗi user có scope: Global (tất cả nhà máy), Factory (một nhà máy), Area (một khu vực). Middleware kiểm tra scope trong mọi API request | P1 — Core |
| AU-012 | Policy Engine | Tích hợp Open Policy Agent (OPA). Policy viết bằng Rego. Middleware gọi OPA để evaluate permission | P2 — Advanced |
| AU-013 | Permission CRUD | API quản lý role và permission. Gán role cho user. Không yêu cầu restart | P1 — Core |

#### 4.9.3 Audit Log

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AU-020 | Event Collection | Thu thập mọi hành động: login, logout, CRUD operations, command dispatch, alarm ack, config change | P1 — Core |
| AU-021 | Audit Store | Lưu trong PostgreSQL table `audit_log`. Partition by month. Retention 2 năm | P1 — Core |
| AU-022 | Audit Viewer | Web UI: lọc theo user, action, resource, time. Không thể sửa/xóa. Export CSV | P1 — Core |
| AU-023 | Compliance Report | Báo cáo định kỳ: danh sách user, role, hoạt động đáng ngờ (login fail > 5 lần, command không được phép) | P2 — Advanced |

---

### 4.10 API Gateway

**Định danh**: `SRV-GATEWAY-API` | **Core**: Kong | **BFF**: Node.js/Go

#### 4.10.1 Kong API Gateway

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AP-001 | Route Management | Định nghĩa routes → upstream services. Service discovery qua Kubernetes DNS | P0 — MVP |
| AP-002 | Auth Plugin | Kong plugin xác thực JWT token từ Keycloak. Validate: signature, expiration, issuer. Inject user info vào header | P0 — MVP |
| AP-003 | Rate Limiting | Rate limit theo: IP, user, API key. Cấu hình: N requests/second. Trả về 429 khi vượt | P1 — Core |
| AP-004 | Circuit Breaker | Circuit breaker pattern: sau N lỗi liên tiếp → mở circuit (fail fast), sau timeout → half-open → thử lại | P2 — Advanced |
| AP-005 | Request Transform | Transform request/response: thêm/sửa headers, body, method. Dùng cho versioning API | P2 — Advanced |

#### 4.10.2 Web BFF

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AP-010 | Dashboard API | Aggregate data từ Telemetry + Device + Alarm service → trả về một response duy nhất cho dashboard | P0 — MVP |
| AP-011 | Device API | Proxy các device operation: CRUD, control, shadow. Thêm business logic validation | P1 — Core |
| AP-012 | Alarm API | Proxy alarm operations: list, ack, snooze. Kết hợp thông tin từ Device Service để enrich | P1 — Core |
| AP-013 | Report API | Proxy analytics/report requests. Hỗ trợ export | P2 — Advanced |
| AP-014 | Admin API | Proxy admin operations: user management, system config. Kiểm tra admin role | P1 — Core |

#### 4.10.3 Mobile BFF

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| AP-020 | Monitoring API | Tối ưu payload cho mobile (chỉ trả về fields cần thiết, giảm nesting) | P3 — Scale |
| AP-021 | Control API | Quick control endpoints. Đơn giản hóa request/response cho mobile | P3 — Scale |
| AP-022 | Sync API | Offline data sync: `GET /api/mobile/sync?since=<timestamp>`. Trả về tất cả thay đổi từ timestamp đó | P3 — Scale |

---

### 4.11 Web Portal

**Định danh**: `SRV-WEB` | **Framework**: React 19 + TypeScript + MUI 6

#### 4.11.1 Dashboard

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| WB-001 | Factory Overview | Tổng quan toàn nhà máy: layout, màu trạng thái khu vực, KPI cards (OEE, throughput, alarms) | P0 — MVP |
| WB-002 | Line Detail | Danh sách máy dạng card grid. Mỗi card: tên máy, trạng thái (màu), metrics chính (gauge mini) | P0 — MVP |
| WB-003 | Machine Detail | Biểu đồ line chart real-time cho từng metric. Device shadow (desired vs reported). Alarm history | P0 — MVP |
| WB-004 | KPI Widgets | Widgets: OEE gauge, throughput counter, energy consumption bar, active alarms counter | P1 — Core |
| WB-005 | Multi-Factory Selector | Chuyển đổi giữa các nhà máy (HN, HCM, ĐN). Ghi nhớ lựa chọn cuối cùng | P1 — Core |

#### 4.11.2 Device Control

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| WB-010 | Device Tree | Cây phân cấp: Factory → Area → Line → Device. Expand/collapse. Search | P1 — Core |
| WB-011 | Control Panel | Bảng điều khiển: Start/Stop button, Speed slider, Setpoint input. Confirmation dialog trước khi gửi | P1 — Core |
| WB-012 | Command History | Bảng lịch sử lệnh: thời gian, user, lệnh, trạng thái, kết quả. Phân trang | P1 — Core |
| WB-013 | Emergency Stop | Nút Emergency Stop lớn, màu đỏ, yêu cầu xác nhận 2 bước. Hiển thị trạng thái sau khi gửi | P1 — Core |

#### 4.11.3 Alarm Center

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| WB-020 | Active Alarms | Bảng alarm đang active: filter (severity, factory, area), sort (time, severity). Acknowledge button | P1 — Core |
| WB-021 | Alarm Detail | Click vào alarm: chi tiết thông số, biểu đồ metric, timeline, comment thread | P1 — Core |
| WB-022 | Alarm History | Lịch sử alarm đã resolved. Filter theo time range. Export CSV | P1 — Core |
| WB-023 | Alarm Configuration | UI tạo/chỉnh sửa rule: chọn loại, metric, điều kiện, severity. Form-based, không cần code | P1 — Core |

#### 4.11.4 Reports

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| WB-030 | Production Report | Báo cáo sản xuất: sản lượng, OEE, downtime, alarm count. Chọn ca/ngày/tuần/tháng | P2 — Advanced |
| WB-031 | OEE Report | OEE breakdown: Availability, Performance, Quality. Biểu đồ trend. So sánh các ca | P2 — Advanced |
| WB-032 | Energy Report | Tiêu thụ năng lượng theo máy/khu vực. So sánh baseline. Biểu đồ bar chart | P2 — Advanced |
| WB-033 | Custom Report Builder | Chọn metrics, filter, time range, group by → preview → export PDF/Excel/CSV | P2 — Advanced |

#### 4.11.5 Admin Panel

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| WB-040 | User Management | CRUD user: username, email, role, scope. Gán role (Admin/Supervisor/Operator/Viewer) | P1 — Core |
| WB-041 | Device Management | Bảng danh sách thiết bị: search, filter, CRUD. Xem device shadow, firmware version | P1 — Core |
| WB-042 | System Configuration | Cấu hình tham số hệ thống: data retention, rate limits, notification settings | P1 — Core |
| WB-043 | Audit Viewer | Xem audit log: filter user, action, resource, time. Read-only. Export CSV | P1 — Core |

#### 4.11.6 Technical Requirements

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| WB-050 | Real-time Updates | Sử dụng WebSocket/SSE cho real-time data. Tự động reconnect |
| WB-051 | Responsive Design | Hỗ trợ desktop (1920×1080) và tablet (1024×768) |
| WB-052 | State Management | Zustand/Jotai cho global state. TanStack Query cho server state |
| WB-053 | Chart Library | Recharts hoặc ECharts. Hỗ trợ: line, bar, gauge, scatter, heatmap |
| WB-054 | Testing | Vitest cho unit test (≥ 80% coverage). Playwright cho E2E test |
| WB-055 | Build | Vite, production build < 2MB gzipped |

---

### 4.12 Mobile App

**Định danh**: `SRV-MOBILE` | **Framework**: Flutter 3.x + Dart

#### 4.12.1 Monitoring

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| MB-001 | Mobile Dashboard | KPI cards (cuộn ngang), danh sách máy với status dot màu, pull-to-refresh | P3 — Scale |
| MB-002 | Machine View | Chi tiết máy: gauge chart, metrics list, alarm gần đây. Landscape mode cho biểu đồ full-width | P3 — Scale |
| MB-003 | Multi-Factory | Chuyển đổi nhà máy qua dropdown/picker. Lưu lựa chọn | P3 — Scale |

#### 4.12.2 Alarm Push

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| MB-010 | Push Handling | Nhận FCM/APNs push. Hiển thị notification với action buttons: "Acknowledge", "Xem chi tiết" | P3 — Scale |
| MB-011 | Alarm Inbox | Danh sách alarm: filter theo severity, status. Swipe to acknowledge | P3 — Scale |
| MB-012 | Quick Actions | Từ notification: nhấn "Acknowledge" không cần mở app. Vuốt để Emergency Stop | P3 — Scale |
| MB-013 | Badge Count | App icon badge = số alarm active. Cập nhật real-time | P3 — Scale |

#### 4.12.3 Quick Control

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| MB-020 | Emergency Stop | Nút Emergency Stop lớn, đỏ. Vuốt để xác nhận. Rung phản hồi | P3 — Scale |
| MB-021 | Quick Commands | Danh sách command thường dùng (Start, Stop, Set Speed). Tối ưu cho thao tác nhanh | P3 — Scale |
| MB-022 | QR Scanner | Quét QR/Barcode trên máy → hiển thị thông tin máy + metrics real-time | P3 — Scale |

#### 4.12.4 Offline Mode

| ID | Yêu cầu | Mô tả | Độ ưu tiên |
|----|---------|-------|------------|
| MB-030 | Local Cache | Lưu dữ liệu cuối cùng vào SQLite/Hive. Hiển thị khi offline | P3 — Scale |
| MB-031 | Offline Indicator | Hiển thị banner "Offline" + thời gian mất kết nối. Ẩn các chức năng cần online | P3 — Scale |
| MB-032 | Auto Sync | Khi có mạng trở lại → tự động sync data mới. Background sync | P3 — Scale |

#### 4.12.5 Technical Requirements

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| MB-050 | Platform Support | iOS 15+ (iPhone, iPad), Android 8+ |
| MB-051 | Biometric Auth | FaceID / TouchID (iOS), Fingerprint (Android) cho login |
| MB-052 | State Management | Riverpod hoặc Bloc pattern |
| MB-053 | Push Provider | FCM (Android) + APNs (iOS) |
| MB-054 | Testing | Widget test + Integration test |

---

## 5. Giao Diện Hệ Thống

### 5.1 Giao Diện Ngoài

| Interface | Protocol | Direction | Mô tả |
|-----------|----------|-----------|-------|
| Edge → Cloud | MQTT over TLS 1.2+ | Outbound | Telemetry data, command responses, status |
| Cloud → Edge | MQTT over TLS 1.2+ | Inbound | Commands, config, firmware |
| Edge → PLC/Machine | Modbus TCP/RTU, OPC-UA, S7 | Local | Đọc/ghi dữ liệu máy móc |
| Web Browser → Cloud | HTTPS (TLS 1.3) | Inbound | Web Portal truy cập API |
| Mobile App → Cloud | HTTPS (TLS 1.3) | Inbound | Mobile App truy cập API |
| External → Notification | SMTP, HTTP | Outbound | Email, SMS (Twilio), Push (FCM/APNs) |
| Admin → Infrastructure | SSH, kubectl, HTTPS | Inbound | Quản trị hệ thống |

### 5.2 Giao Diện Trong (Inter-Service)

| From | To | Protocol | Mô tả |
|------|----|----------|-------|
| Gateway Service | Kafka | Kafka Producer API | Publish telemetry, events |
| Kafka | Telemetry Service | Kafka Consumer API | Consume telemetry |
| Kafka | Alarm Service | Kafka Consumer API | Consume telemetry để evaluate rules |
| Kafka | Device Service | Kafka Consumer API | Consume device status |
| Alarm Service | Kafka | Kafka Producer API | Publish alarm events |
| Kafka | Notification Service | Kafka Consumer API | Consume alarm events để gửi notification |
| Command Service | Kafka → Gateway Svc | Kafka → MQTT | Gửi command xuống edge |
| Kafka | Command Service | Kafka Consumer API | Nhận command response |
| API Gateway (BFF) | All Services | REST/gRPC | Aggregate data cho UI |
| All Services | Keycloak | OAuth2/OIDC | Xác thực token |
| All Services | Redis | Redis Protocol | Cache, pub/sub |
| All Services | PostgreSQL/TimescaleDB | PostgreSQL Protocol | Database operations |
| All Services | Prometheus | HTTP/metrics | Export metrics |

### 5.3 Kafka Topics

| Topic | Partitions | Replication | Retention | Mô tả |
|-------|-----------|-------------|-----------|-------|
| `telemetry.raw` | 12 | 3 | 7 days | Dữ liệu telemetry thô từ gateway |
| `telemetry.dlq` | 3 | 3 | 30 days | Dead letter queue |
| `device.status` | 6 | 3 | 7 days | Trạng thái online/offline thiết bị |
| `device.alarms` | 6 | 3 | 7 days | Alarm từ edge device |
| `command.requests` | 6 | 3 | 7 days | Lệnh điều khiển |
| `command.responses` | 6 | 3 | 7 days | Kết quả lệnh |
| `alarm.events` | 6 | 3 | 7 days | Alarm events từ Alarm Service |
| `notification.requests` | 3 | 3 | 7 days | Yêu cầu gửi notification |
| `analytics.kpi` | 3 | 3 | 7 days | KPI tính toán real-time |

### 5.4 MQTT Topics

| Topic Pattern | QoS | Direction | Mô tả |
|---------------|-----|-----------|-------|
| `factory/{factory_id}/{area}/{machine_id}/{metric}` | 1 | Edge → Cloud | Telemetry data |
| `devices/{device_id}/telemetry` | 1 | Edge → Cloud | Telemetry data (alt) |
| `devices/{device_id}/status` | 1 | Edge → Cloud | Online/offline/heartbeat |
| `devices/{device_id}/alarm` | 1 | Edge → Cloud | Alarm từ edge |
| `devices/{device_id}/command/req` | 1 | Cloud → Edge | Lệnh điều khiển |
| `devices/{device_id}/command/res` | 1 | Edge → Cloud | Kết quả lệnh |
| `devices/{device_id}/config` | 1 | Cloud → Edge | Cấu hình OTA |
| `devices/{device_id}/firmware` | 1 | Cloud → Edge | Firmware update URL |

---

## 6. Mô Hình Dữ Liệu

### 6.1 PostgreSQL Schema (Core)

```sql
-- Factories
CREATE TABLE factories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(10) UNIQUE NOT NULL,  -- HN, HCM, DN
    name        VARCHAR(255) NOT NULL,
    address     TEXT,
    timezone    VARCHAR(50) DEFAULT 'Asia/Ho_Chi_Minh',
    status      VARCHAR(20) DEFAULT 'active',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Devices
CREATE TABLE devices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_number   VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    model           VARCHAR(100),
    factory_id      UUID REFERENCES factories(id),
    area            VARCHAR(100),
    line            VARCHAR(100),
    protocol        VARCHAR(50),           -- modbus_rtu, modbus_tcp, opcua, s7, mqtt
    connection_info JSONB,                 -- {ip, port, slave_id, ...}
    capabilities    JSONB,                 -- {metrics: [...], commands: [...]}
    status          VARCHAR(20) DEFAULT 'offline',  -- online, offline, maintenance, error
    metadata        JSONB,
    last_seen       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_devices_factory ON devices(factory_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_area ON devices(factory_id, area);

-- Commands
CREATE TABLE commands (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id       UUID REFERENCES devices(id),
    user_id         VARCHAR(255) NOT NULL,
    command_type    VARCHAR(50) NOT NULL,   -- start, stop, set_speed, set_temp, ...
    parameters      JSONB,
    priority        VARCHAR(20) DEFAULT 'NORMAL',  -- LOW, NORMAL, HIGH, EMERGENCY
    status          VARCHAR(20) DEFAULT 'PENDING',  -- PENDING, SENT, EXECUTING, DONE, FAILED, TIMEOUT
    result          JSONB,
    error_message   TEXT,
    sent_at         TIMESTAMPTZ,
    executed_at     TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_commands_device ON commands(device_id);
CREATE INDEX idx_commands_status ON commands(status);
CREATE INDEX idx_commands_created ON commands(created_at DESC);

-- Alarms
CREATE TABLE alarms (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id       UUID REFERENCES devices(id),
    rule_id         UUID REFERENCES alarm_rules(id),
    metric_name     VARCHAR(100) NOT NULL,
    severity        VARCHAR(20) NOT NULL,   -- EMERGENCY, CRITICAL, WARNING, INFO
    status          VARCHAR(20) DEFAULT 'RAISED',  -- RAISED, ACKNOWLEDGED, RESOLVED, SNOOZED, SUPPRESSED
    trigger_value   DOUBLE PRECISION,
    threshold_value DOUBLE PRECISION,
    acknowledged_by VARCHAR(255),
    acknowledged_at TIMESTAMPTZ,
    resolved_at     TIMESTAMPTZ,
    comment         TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_alarms_device ON alarms(device_id);
CREATE INDEX idx_alarms_status ON alarms(status);
CREATE INDEX idx_alarms_severity ON alarms(severity);
CREATE INDEX idx_alarms_created ON alarms(created_at DESC);

-- Alarm Rules
CREATE TABLE alarm_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    rule_type       VARCHAR(50) NOT NULL,   -- threshold, rate_of_change, state_change, anomaly
    metric_name     VARCHAR(100) NOT NULL,
    condition       JSONB NOT NULL,         -- {operator: ">", value: 800, duration_seconds: 30}
    severity        VARCHAR(20) NOT NULL,
    scope           JSONB,                  -- {factory_id: "...", area: "...", device_ids: [...]}
    enabled         BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Users (đồng bộ từ Keycloak)
CREATE TABLE users (
    id              VARCHAR(255) PRIMARY KEY,  -- Keycloak user ID
    username        VARCHAR(100) UNIQUE NOT NULL,
    email           VARCHAR(255),
    full_name       VARCHAR(255),
    role            VARCHAR(50) NOT NULL,      -- admin, supervisor, operator, viewer
    scope_type      VARCHAR(20) DEFAULT 'global',  -- global, factory, area
    scope_value     VARCHAR(255),              -- factory_id hoặc area
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Audit Log
CREATE TABLE audit_log (
    id              BIGSERIAL,
    user_id         VARCHAR(255),
    username        VARCHAR(100),
    action          VARCHAR(100) NOT NULL,     -- login, logout, create_device, send_command, ack_alarm, ...
    resource_type   VARCHAR(50),               -- device, alarm, user, config
    resource_id     VARCHAR(255),
    details         JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_audit_user ON audit_log(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_log(action, created_at DESC);
```

### 6.2 TimescaleDB Schema (Time-Series)

Xem chi tiết tại [4.4.4 TimescaleDB Schema](#444-timescaledb-schema).

### 6.3 Redis Data Structures

| Key Pattern | Type | TTL | Mô tả |
|-------------|------|-----|-------|
| `device:{id}:latest:{metric}` | String (JSON) | 24h | Latest value của metric |
| `device:{id}:shadow:desired` | Hash | None | Desired state |
| `device:{id}:shadow:reported` | Hash | None | Reported state |
| `device:{id}:status` | String | 30s | online/offline + last seen |
| `session:{token_hash}` | String | 8h | User session cache |
| `rate_limit:{user_id}:{endpoint}` | Sorted Set | 1m | Rate limit counter |
| `alarm:active:{factory_id}` | Set | None | Set of active alarm IDs |
| `rule_engine:rules` | String (JSON) | None | Cached alarm rules |

---

## 7. Yêu Cầu Phi Chức Năng

### 7.1 Hiệu Năng (Performance)

| ID | Yêu cầu | Chỉ số mục tiêu | Cách đo |
|----|---------|----------------|---------|
| NFR-P01 | MQTT Message Throughput | 100,000 msg/s | EMQX metrics |
| NFR-P02 | Kafka Consumer Throughput | 100,000 msg/s | Kafka consumer lag |
| NFR-P03 | Database Write Throughput | 100,000 inserts/s | TimescaleDB metrics |
| NFR-P04 | Telemetry End-to-End Latency | < 2 giây (P99) | MQTT publish → DB insert + cache |
| NFR-P05 | Command Round-Trip Latency | < 500ms (P99) | UI click → MQTT publish |
| NFR-P06 | Dashboard Load Time | < 2 giây (P95) | Lighthouse / Web Vitals |
| NFR-P07 | API Response Time | < 200ms (P95) cho simple queries | APM (Prometheus histogram) |
| NFR-P08 | Alarm Detection Latency | < 5 giây (P99) | Rule trigger → notification sent |
| NFR-P09 | WebSocket Message Latency | < 1 giây (P99) | Kafka → WebSocket push |

### 7.2 Khả Năng Mở Rộng (Scalability)

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NFR-S01 | Horizontal Scaling | Tất cả microservices có thể scale ngang bằng cách tăng replicas |
| NFR-S02 | Device Connections | Hỗ trợ 10,000+ thiết bị kết nối đồng thời |
| NFR-S03 | Multi-Region | Hỗ trợ triển khai edge tại 10+ nhà máy |
| NFR-S04 | Data Volume | Lưu trữ 10TB+ telemetry data/năm |
| NFR-S05 | Kafka Partitions | Thiết kế partition key để consumer scale không giới hạn |
| NFR-S06 | Database Scaling | TimescaleDB: partitioning + compression. PostgreSQL: read replicas |

### 7.3 Độ Sẵn Sàng (Availability)

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NFR-A01 | Cloud Uptime | 99.9% (downtime < 8.76 giờ/năm) |
| NFR-A02 | Edge Resilience | Edge gateway hoạt động độc lập khi mất kết nối cloud |
| NFR-A03 | Database HA | PostgreSQL primary-standby replication. TimescaleDB tương tự |
| NFR-A04 | Kafka HA | 3 brokers, replication factor = 3, min ISR = 2 |
| NFR-A05 | EMQX HA | Cluster 3+ nodes |
| NFR-A06 | Graceful Degradation | Nếu Kafka chết → data ở edge buffer. Nếu DB chết → API đọc từ cache |
| NFR-A07 | Disaster Recovery | RPO < 1 giờ, RTO < 4 giờ |

### 7.4 Độ Tin Cậy (Reliability)

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NFR-R01 | Message Delivery | At-least-once delivery. Telemetry data loss rate < 0.01% |
| NFR-R02 | Edge Store-and-Forward | Không mất dữ liệu khi mất kết nối < 7 ngày |
| NFR-R03 | Kafka Durability | `acks=all`, `min.insync.replicas=2` cho topic quan trọng |
| NFR-R04 | Retry Logic | Retry với exponential backoff cho tất cả external calls |
| NFR-R05 | Circuit Breaker | Ngăn lỗi lan truyền giữa services |

### 7.5 Bảo Mật (Security)

Xem chi tiết tại [8. Yêu Cầu Bảo Mật](#8-yêu-cầu-bảo-mật).

### 7.6 Khả Năng Bảo Trì (Maintainability)

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NFR-M01 | Code Coverage | Unit test coverage ≥ 80% |
| NFR-M02 | API Documentation | OpenAPI 3.0 cho tất cả REST APIs |
| NFR-M03 | Logging | Structured logging (JSON format). Correlation ID xuyên suốt request |
| NFR-M04 | Health Checks | Tất cả services có `/health` và `/ready` endpoints |
| NFR-M05 | Graceful Shutdown | Xử lý SIGTERM, drain connections trước khi exit |
| NFR-M06 | Configuration | Externalized config (YAML files + env vars). Không hard-code |

### 7.7 Khả Năng Quan Sát (Observability)

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| NFR-O01 | Metrics | Prometheus metrics: request rate, latency, error rate, Kafka lag, DB connections |
| NFR-O02 | Tracing | Distributed tracing với Jaeger/Zipkin. Trace ID xuyên suốt Kafka |
| NFR-O03 | Logging | Centralized logging với Loki. Log level configurable per service |
| NFR-O04 | Alerting | Prometheus AlertManager rules: service down, high latency, Kafka lag > 1000, disk > 80% |
| NFR-O05 | Dashboards | Grafana dashboards: service overview, Kafka overview, DB overview, business KPIs |

---

## 8. Yêu Cầu Bảo Mật

### 8.1 Xác Thực & Phân Quyền

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| SEC-01 | User Authentication | OAuth2/OIDC qua Keycloak. Authorization Code Flow + PKCE cho web |
| SEC-02 | Device Authentication | X.509 client certificates cho MQTT devices. Validate chain, expiration, revocation |
| SEC-03 | API Authentication | JWT Bearer token cho REST APIs. Validate: signature (RS256), expiration, issuer, audience |
| SEC-04 | Service-to-Service | mTLS hoặc API key cho inter-service communication |
| SEC-05 | RBAC | 4 roles: Admin, Supervisor, Operator, Viewer. Scope-based access (global/factory/area) |
| SEC-06 | Session Management | Access token TTL: 15 phút. Refresh token TTL: 8 giờ. Session timeout sau 8 giờ idle |

### 8.2 Mã Hóa

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| SEC-10 | Data in Transit | TLS 1.2+ cho mọi kết nối: MQTT, HTTPS, Kafka, DB connections |
| SEC-11 | Data at Rest | Mã hóa database storage (cloud provider KMS). Mã hóa object storage (MinIO SSE) |
| SEC-12 | Secrets Management | Sử dụng Kubernetes Secrets hoặc HashiCorp Vault. Không hard-code secrets |

### 8.3 Audit & Compliance

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| SEC-20 | Audit Trail | Ghi log mọi hành động: user, action, resource, timestamp, IP. Immutable |
| SEC-21 | Command Audit | Mọi lệnh điều khiển máy móc đều được ghi audit log đặc biệt |
| SEC-22 | Compliance | Tuân thủ Nghị định 53/2022/NĐ-CP về bảo vệ dữ liệu cá nhân |
| SEC-23 | Data Retention | Audit log giữ 2 năm. Telemetry data giữ 2 năm |

### 8.4 Network Security

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| SEC-30 | Edge-Cloud VPN | IPsec/Site-to-Site VPN giữa nhà máy và cloud (tùy chọn) |
| SEC-31 | Network Policies | Kubernetes Network Policies: hạn chế east-west traffic giữa services |
| SEC-32 | Firewall | WAF (Web Application Firewall) trước API Gateway |
| SEC-33 | DDoS Protection | Cloud provider DDoS protection + Kong rate limiting |

### 8.5 Device Security

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| SEC-40 | Secure Boot | Edge gateway hỗ trợ secure boot (UEFI) |
| SEC-41 | Firmware Signing | Firmware được ký với private key. Edge verify chữ ký trước khi cài |
| SEC-42 | Minimal Attack Surface | Edge gateway chỉ mở port MQTT outbound. Không mở inbound ports |
| SEC-43 | Physical Security | Cảnh báo nếu edge gateway bị mở case (tamper detection, nếu hardware hỗ trợ) |

---

## 9. Yêu Cầu Vận Hành

### 9.1 Deployment

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| OPS-01 | Containerization | Tất cả services được đóng gói Docker image |
| OPS-02 | Orchestration | Kubernetes deployment với Helm charts |
| OPS-03 | Edge Deployment | Ansible playbooks cho edge gateway. Cài đặt < 30 phút |
| OPS-04 | CI/CD | GitHub Actions pipeline: build → test → scan → push image → deploy |
| OPS-05 | Environment Promotion | Dev → Staging → Production. Canary deployment cho production |
| OPS-06 | Rollback | Tự động rollback nếu health check fail sau deploy |

### 9.2 Backup & Restore

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| OPS-10 | Database Backup | PostgreSQL: pg_dump hàng ngày + WAL archiving liên tục. TimescaleDB: backup policy |
| OPS-11 | Object Storage Backup | MinIO bucket replication hoặc sync sang secondary storage |
| OPS-12 | Backup Retention | Daily backup giữ 30 ngày, weekly giữ 3 tháng, monthly giữ 1 năm |
| OPS-13 | Restore Drill | Diễn tập restore mỗi quý |

### 9.3 Monitoring & Alerting

| ID | Yêu cầu | Mô tả |
|----|---------|-------|
| OPS-20 | Infrastructure Monitoring | CPU, RAM, Disk, Network của tất cả nodes (Prometheus Node Exporter) |
| OPS-21 | Application Monitoring | Request rate, latency, error rate, throughput (Prometheus + Grafana) |
| OPS-22 | Database Monitoring | Connections, query performance, replication lag, disk usage |
| OPS-23 | Kafka Monitoring | Consumer lag, throughput, partition status |
| OPS-24 | MQTT Monitoring | Connections, messages/sec, bytes/sec, session count |
| OPS-25 | Edge Monitoring | Edge health: CPU, RAM, Disk, Network, MQTT connection status |

---

## 10. Ràng Buộc Thiết Kế

### 10.1 Ràng Buộc Kỹ Thuật

| ID | Ràng buộc | Mô tả |
|----|-----------|-------|
| DC-01 | Go Version | Go 1.22+ cho tất cả backend services |
| DC-02 | Protocol Buffers | gRPC services sử dụng proto3 |
| DC-03 | REST API Versioning | URL-based: `/api/v1/`, `/api/v2/` |
| DC-04 | Error Handling | Standard error response format: `{"error": {"code": "...", "message": "...", "details": []}}` |
| DC-05 | Pagination | Cursor-based pagination cho large datasets. Response: `{"data": [], "next_cursor": "..."}` |
| DC-06 | Timezone | Tất cả timestamps là UTC. UI chuyển đổi sang local timezone |
| DC-07 | Character Encoding | UTF-8 throughout |

### 10.2 Ràng Buộc Vận Hành

| ID | Ràng buộc | Mô tả |
|----|-----------|-------|
| DC-10 | Edge OS | Linux (Ubuntu Server 22.04 LTS hoặc Raspberry Pi OS) |
| DC-11 | Edge Resources | Tối thiểu: 1GB RAM, 16GB storage, ARM Cortex-A72 hoặc x86 |
| DC-12 | Cloud Provider | Cloud-agnostic (ưu tiên AWS, hỗ trợ GCP/Azure) |
| DC-13 | Browser Support | Chrome, Firefox, Edge, Safari — 2 versions gần nhất |

---

## 11. Phụ Lục

### A. Service Port Map

| Service | Port | Protocol |
|---------|------|----------|
| EMQX MQTT | 8883 | MQTTS |
| EMQX Dashboard | 18083 | HTTP |
| Kafka Bootstrap | 9092 | TCP |
| Kafka (TLS) | 9093 | TLS |
| TimescaleDB | 5432 | TCP |
| PostgreSQL | 5433 | TCP |
| Redis | 6379 | TCP |
| MinIO API | 9000 | HTTP |
| MinIO Console | 9001 | HTTP |
| Keycloak | 8080 | HTTP |
| Kong Proxy | 8443 | HTTPS |
| Kong Admin | 8001 | HTTP |
| Web Portal (Dev) | 5173 | HTTP |
| Prometheus | 9090 | HTTP |
| Grafana | 3000 | HTTP |
| Gateway Service | 8081 | HTTP |
| Telemetry Service | 8082 | HTTP |
| Device Service | 8083 | HTTP |
| Command Service | 8084 | HTTP |
| Alarm Service | 8085 | HTTP |
| Notification Service | 8086 | HTTP |
| Analytics Service | 8087 | HTTP |
| Web BFF | 8090 | HTTP |
| Mobile BFF | 8091 | HTTP |

### B. Mã Lỗi Chuẩn

| HTTP Status | Error Code | Mô tả |
|-------------|-----------|-------|
| 400 | `BAD_REQUEST` | Request không hợp lệ |
| 401 | `UNAUTHENTICATED` | Chưa đăng nhập hoặc token hết hạn |
| 403 | `PERMISSION_DENIED` | Không có quyền truy cập |
| 404 | `NOT_FOUND` | Resource không tồn tại |
| 409 | `CONFLICT` | Xung đột dữ liệu |
| 422 | `VALIDATION_ERROR` | Dữ liệu không hợp lệ |
| 429 | `RATE_LIMIT_EXCEEDED` | Vượt quá rate limit |
| 500 | `INTERNAL_ERROR` | Lỗi hệ thống |
| 503 | `SERVICE_UNAVAILABLE` | Service đang bảo trì/quá tải |

### C. Ký Hiệu & Viết Tắt

Xem BRD [Phụ Lục A](BRD.md#a-thuật-ngữ--viết-tắt).

### D. Quy Ước Đặt Tên

- **Service names**: lowercase, hyphenated (vd: `gateway-service`, `device-service`)
- **Kafka topics**: lowercase, dot-separated (vd: `telemetry.raw`, `device.status`)
- **MQTT topics**: lowercase, slash-separated (vd: `factory/HN/area_a/lathe_01/temperature`)
- **API endpoints**: `/api/v{n}/{resource}` (vd: `/api/v1/devices`)
- **Environment variables**: UPPERCASE, underscore-separated (vd: `KAFKA_BOOTSTRAP_SERVERS`)
- **Database tables**: lowercase, underscore-separated (vd: `device_latest_value`)
- **Redis keys**: lowercase, colon-separated (vd: `device:123:latest:temperature`)

### E. Tài Liệu Tham Khảo

- [Business Requirements Document (BRD)](BRD.md)
- [User Requirements Document (URD)](URD.md)
- [README.md](../README.md) — Kiến trúc tổng thể
- README của từng service
- [IEEE 830-1998 — Recommended Practice for Software Requirements Specifications](https://standards.ieee.org/standard/830-1998.html)

---

> **Tài liệu này là tài sản trí tuệ của Industrial IoT Platform. Không được sao chép hoặc phân phối khi chưa có sự đồng ý.**
