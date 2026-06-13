# Edge Gateway

> Cổng kết nối tại chỗ - triển khai tại mỗi nhà máy (HN, HCM, Đà Nẵng)

## 🎯 Chức Năng

- Kết nối trực tiếp với PLC/máy móc qua các giao thức công nghiệp (Modbus RTU/TCP, OPC-UA, Profinet, Siemens S7)
- Chuyển đổi protocol sang MQTT để giao tiếp với cloud
- Local cache & buffer — đảm bảo không mất dữ liệu khi mất kết nối internet
- Thực thi lệnh điều khiển local với latency < 10ms
- Rule engine cơ bản chạy tại edge — alarm offline khi mất kết nối
- Hỗ trợ store-and-forward: tự động đồng bộ data khi kết nối được khôi phục

## 📁 Cấu Trúc

```
edge-gateway/
├── protocol-adapter/        # Adapter cho các giao thức công nghiệp
│   ├── modbus/              #   Modbus RTU & Modbus TCP
│   ├── opcua/               #   OPC-UA client
│   ├── s7/                  #   Siemens S7 protocol
│   └── mqtt-client/         #   MQTT client (publish/subscribe)
├── local-cache/             # Buffer dữ liệu khi mất kết nối
│   ├── sqlite-store/        #   SQLite làm storage local
│   └── store-forward/       #   Cơ chế store-and-forward
├── command-executor/        # Thực thi lệnh điều khiển local
│   ├── dispatcher/          #   Nhận lệnh từ cloud, phân phối
│   └── executor/            #   Ghi lệnh xuống PLC
└── edge-rules-engine/       # Rule engine chạy tại edge
    └── rules/               #   Các rule cơ bản (ngưỡng, watchdog)
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go (nhẹ, binary single, chạy tốt trên ARM)
- **Protocol libraries**: `gopcua`, `go-modbus`, `paho.mqtt.golang`
- **Local DB**: SQLite / BoltDB
- **Deploy target**: Raspberry Pi 4, Linux edge server

## 🚀 Chạy

```bash
# Development
go run cmd/edge-gateway/main.go --config configs/dev.yaml

# Build binary
make build

# Build Docker image cho ARM
make docker-arm
```

## 📡 Giao Tiếp

- **Edge → Cloud**: MQTT over TLS (port 8883)
- **Edge → Machine**: Modbus TCP (port 502), OPC-UA (port 4840)
- **Health check**: MQTT heartbeat topic mỗi 10s
