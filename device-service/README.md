# Device Service

> Quản lý vòng đời thiết bị - Device Registry, Digital Twin & Firmware Management

## 🎯 Chức Năng

- **Device Registry**: CRUD thiết bị, metadata (model, location, capabilities)
- **Device Shadow (Digital Twin)**: Lưu trạng thái mong muốn (desired) và thực tế (reported)
- **Firmware Management**: Over-the-air (OTA) firmware update
- **Device Grouping**: Quản lý thiết bị theo nhà máy, khu vực, dây chuyền
- **Device Provisioning**: Tự động provisioning qua MQTT

## 📁 Cấu Trúc

```
device-service/
├── device-registry/         # Đăng ký & quản lý thiết bị
│   ├── api/                 #   REST/gRPC API
│   ├── model/               #   Device data model
│   ├── repository/          #   Database access layer
│   └── service/             #   Business logic
├── device-shadow/           # Digital Twin
│   ├── desired-state/       #   Trạng thái mong muốn (từ cloud)
│   ├── reported-state/      #   Trạng thái thực tế (từ thiết bị)
│   └── delta-detector/      #   Phát hiện sai khác desired vs reported
└── firmware-management/     # OTA Firmware Update
    ├── firmware-store/      #   Lưu trữ firmware binaries
    ├── update-scheduler/    #   Lập lịch cập nhật
    └── rollback-manager/    #   Rollback khi update thất bại
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go
- **Database**: PostgreSQL 16
- **Cache**: Redis (device shadow real-time state)
- **Object Storage**: MinIO (firmware binaries)

## 📊 Data Model

```sql
-- Device
CREATE TABLE devices (
    id          UUID PRIMARY KEY,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    model       VARCHAR(100),
    factory_id  UUID REFERENCES factories(id),
    area        VARCHAR(100),
    protocol    VARCHAR(50),     -- modbus, opcua, mqtt
    status      VARCHAR(20),     -- online, offline, maintenance
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ
);
```
