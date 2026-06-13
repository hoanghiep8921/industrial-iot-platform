# Telemetry Service

> Thu thập, lưu trữ & truy vấn dữ liệu đo lường từ máy móc công nghiệp

## 🎯 Chức Năng

- **Data Ingestion**: Nhận telemetry data từ Kafka, validate schema, enrich metadata
- **Time-Series Storage**: Lưu trữ hiệu quả trong TimescaleDB với hypertable
- **Data Query API**: RESTful & GraphQL API cho dashboard và analytics
- **Data Retention**: Tự động down-sampling và xóa data cũ theo policy
- **Real-time Stream**: WebSocket/SSE stream cho real-time dashboard

## 📁 Cấu Trúc

```
telemetry-service/
├── data-ingestion/          # Nhận & xử lý data từ Kafka
│   ├── consumer/            #   Kafka consumer groups
│   ├── validator/           #   Schema validation (JSON Schema)
│   ├── enricher/            #   Enrich với device metadata
│   └── transformer/         #   Transform data format
├── time-series-storage/     # Lưu trữ TimescaleDB
│   ├── migrations/          #   Database migrations
│   ├── hypertable-config/   #   Hypertable partitioning
│   ├── retention-policy/    #   Data retention & compression
│   └── repository/          #   Data access layer
└── data-query-api/          # Query API
    ├── rest/                #   REST endpoints
    ├── graphql/             #   GraphQL schema & resolvers
    └── websocket/           #   Real-time data stream
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go
- **Database**: TimescaleDB 2.15+ (PostgreSQL 16)
- **Message Queue**: Kafka consumer
- **Cache**: Redis (latest value cache)
- **API**: REST + GraphQL + WebSocket

## 📊 TimescaleDB Schema

```sql
CREATE TABLE telemetry (
    time        TIMESTAMPTZ NOT NULL,
    device_id   UUID NOT NULL,
    metric_name VARCHAR(100) NOT NULL,  -- temperature, pressure, rpm, ...
    value       DOUBLE PRECISION,
    quality     SMALLINT DEFAULT 0,      -- 0=good, 1=uncertain, 2=bad
    metadata    JSONB
);

SELECT create_hypertable('telemetry', 'time');

-- Compression policy (compress data older than 7 days)
SELECT add_compression_policy('telemetry', INTERVAL '7 days');

-- Retention policy (delete data older than 2 years)
SELECT add_retention_policy('telemetry', INTERVAL '2 years');
```

## 🚀 Chạy

```bash
go run cmd/telemetry/main.go --config configs/dev.yaml
```
