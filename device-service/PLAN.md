# Kế Hoạch Triển Khai Device Service

> **Branch**: `feature/device-service`
> **Ngày**: 2026-06-14
> **Mục tiêu**: Quản lý vòng đời thiết bị đầy đủ — Registry, Digital Twin, Firmware OTA

---

## 1. Tổng Quan

Device Service là service nền tảng, cung cấp **3 subsystem chính**:

| Subsystem | Mô tả |
|-----------|-------|
| **Device Registry** | CRUD thiết bị, metadata, grouping theo factory/area/line |
| **Device Shadow** | Digital Twin — desired state vs reported state, delta detection |
| **Firmware Management** | OTA firmware update, upload/download, scheduler, rollback |

### Vị trí trong hệ thống

```
[Máy móc/PLC] → [Edge Gateway] → MQTT → Gateway Service → Kafka
                                                              ↓
                                          ┌───────────────────┴──────────┐
                                          │  Device Service               │
                                          │  ┌─────────────────────────┐  │
                                          │  │ Device Registry (PG)    │  │
                                          │  │ Device Shadow  (PG+Redis)│  │◄── Các service khác gọi API
                                          │  │ Firmware Mgmt  (MinIO)  │  │
                                          │  └─────────────────────────┘  │
                                          └──────────────────────────────┘
```

---

## 2. Kiến Trúc Chi Tiết

### 2.1 Cấu trúc thư mục

```
device-service/
├── cmd/
│   └── device/
│       └── main.go                       # Entry point, wire-up dependencies
├── internal/
│   ├── model/
│   │   ├── device.go                     # Device, DeviceStatus, DeviceFilter
│   │   ├── shadow.go                     # DeviceShadow, DesiredState, ReportedState
│   │   └── firmware.go                   # Firmware, FirmwareUpdate, UpdateStatus
│   ├── repository/
│   │   ├── postgres.go                   # PostgreSQL: devices, shadows, firmware metadata
│   │   └── migrations.go                 # Auto-migration on startup
│   ├── storage/
│   │   └── minio.go                      # MinIO client: upload/download/delete firmware
│   ├── service/
│   │   ├── registry_service.go           # Device CRUD business logic
│   │   ├── shadow_service.go             # Shadow sync, delta calculation
│   │   ├── firmware_service.go           # Firmware upload, campaign, rollback
│   │   └── provisioning.go               # Auto-provisioning via MQTT
│   └── api/
│       ├── server.go                     # HTTP server + router setup
│       ├── registry_handler.go           # Device CRUD handlers
│       ├── shadow_handler.go             # Shadow handlers
│       ├── firmware_handler.go           # Firmware upload/download/update handlers
│       └── middleware.go                 # Logging, CORS, recovery
├── migrations/
│   └── 001_init.sql                      # 5 tables: devices, shadows, connections, firmwares, firmware_updates
├── Dockerfile                            # Docker multi-stage build
├── Makefile                              # Build, run, test
├── go.mod
└── README.md
```

### 2.2 Công nghệ sử dụng

| Layer | Công nghệ | Lý do |
|-------|-----------|-------|
| Ngôn ngữ | Go 1.22 | Cùng stack với gateway, telemetry |
| HTTP Router | `gorilla/mux` | Giống telemetry-service |
| Database | PostgreSQL 16 (`pgx/v5`) | Metadata, shadow, firmware registry |
| Cache | Redis (`go-redis/v9`) | Real-time device state, shadow cache |
| Object Storage | MinIO (`minio-go/v7`) | Firmware binary storage |
| Logging | `go.uber.org/zap` | Giống các service khác |
| Container | Docker multi-stage | Giống pattern hiện có |

---

## 3. Data Model

### 3.1 Bảng `devices`

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE devices (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    model         VARCHAR(100),
    vendor        VARCHAR(100),
    factory_id    VARCHAR(50) NOT NULL,         -- HN, HCM, DN
    area          VARCHAR(100),                 -- area_a, area_b
    line          VARCHAR(100),                 -- line_1, line_2
    protocol      VARCHAR(50) DEFAULT 'mqtt',   -- modbus, opcua, mqtt
    capabilities  JSONB DEFAULT '[]',           -- ["temperature", "vibration", "speed"]
    status        VARCHAR(20) DEFAULT 'offline', -- online, offline, maintenance, error
    firmware_version VARCHAR(50),               -- Current firmware version
    metadata      JSONB DEFAULT '{}',
    last_seen_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_devices_factory   ON devices(factory_id);
CREATE INDEX idx_devices_status    ON devices(status);
CREATE INDEX idx_devices_protocol  ON devices(protocol);
CREATE INDEX idx_devices_area      ON devices(factory_id, area);
```

### 3.2 Bảng `device_shadows`

```sql
CREATE TABLE device_shadows (
    device_id       UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    desired_state   JSONB NOT NULL DEFAULT '{}',
    reported_state  JSONB NOT NULL DEFAULT '{}',
    delta           JSONB NOT NULL DEFAULT '{}',
    version         BIGINT DEFAULT 1,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 3.3 Bảng `device_connections`

```sql
CREATE TABLE device_connections (
    device_id       UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    connected       BOOLEAN DEFAULT FALSE,
    connected_at    TIMESTAMPTZ,
    disconnected_at TIMESTAMPTZ,
    ip_address      VARCHAR(45),
    mqtt_client_id  VARCHAR(255),
    last_heartbeat  TIMESTAMPTZ
);
```

### 3.4 Bảng `firmwares`

```sql
CREATE TABLE firmwares (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version         VARCHAR(50) NOT NULL,
    model           VARCHAR(100),              -- Target device model (NULL = all)
    description     TEXT,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    checksum_sha256 VARCHAR(64) NOT NULL,
    minio_path      VARCHAR(500) NOT NULL,      -- Path in MinIO bucket
    status          VARCHAR(20) DEFAULT 'draft', -- draft, released, deprecated
    created_by      VARCHAR(100),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    released_at     TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_firmwares_version_model ON firmwares(version, model);
CREATE INDEX idx_firmwares_status ON firmwares(status);
```

### 3.5 Bảng `firmware_updates`

```sql
CREATE TABLE firmware_updates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firmware_id     UUID NOT NULL REFERENCES firmwares(id),
    device_id       UUID NOT NULL REFERENCES devices(id),
    campaign_name   VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'pending',  -- pending, downloading, installing, success, failed, rolled_back
    from_version    VARCHAR(50),
    to_version      VARCHAR(50) NOT NULL,
    priority        INT DEFAULT 0,                   -- 0=normal, 1=high, 2=critical
    retry_count     INT DEFAULT 0,
    max_retries     INT DEFAULT 3,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    error_message   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_fw_updates_device  ON firmware_updates(device_id);
CREATE INDEX idx_fw_updates_status  ON firmware_updates(status);
CREATE INDEX idx_fw_updates_campaign ON firmware_updates(campaign_name);
```

### 3.6 Go Models

```go
// ========== Device Registry ==========

type Device struct {
    ID              string         `json:"id"`
    SerialNumber    string         `json:"serial_number"`
    Name            string         `json:"name"`
    Model           string         `json:"model,omitempty"`
    Vendor          string         `json:"vendor,omitempty"`
    FactoryID       string         `json:"factory_id"`
    Area            string         `json:"area,omitempty"`
    Line            string         `json:"line,omitempty"`
    Protocol        string         `json:"protocol"`
    Capabilities    []string       `json:"capabilities"`
    Status          DeviceStatus   `json:"status"`
    FirmwareVersion string         `json:"firmware_version,omitempty"`
    Metadata        map[string]any `json:"metadata,omitempty"`
    LastSeenAt      *time.Time     `json:"last_seen_at,omitempty"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

type DeviceStatus string
const (
    StatusOnline      DeviceStatus = "online"
    StatusOffline     DeviceStatus = "offline"
    StatusMaintenance DeviceStatus = "maintenance"
    StatusError       DeviceStatus = "error"
)

// ========== Device Shadow ==========

type DeviceShadow struct {
    DeviceID      string         `json:"device_id"`
    DesiredState  map[string]any `json:"desired_state"`
    ReportedState map[string]any `json:"reported_state"`
    Delta         map[string]any `json:"delta"`
    Version       int64          `json:"version"`
    UpdatedAt     time.Time      `json:"updated_at"`
}

// ========== Firmware ==========

type Firmware struct {
    ID             string        `json:"id"`
    Version        string        `json:"version"`
    Model          string        `json:"model,omitempty"`
    Description    string        `json:"description,omitempty"`
    FileName       string        `json:"file_name"`
    FileSize       int64         `json:"file_size"`
    ChecksumSHA256 string        `json:"checksum_sha256"`
    Status         FirmwareStatus `json:"status"`
    CreatedBy      string        `json:"created_by,omitempty"`
    CreatedAt      time.Time     `json:"created_at"`
    ReleasedAt     *time.Time    `json:"released_at,omitempty"`
}

type FirmwareStatus string
const (
    FirmwareDraft      FirmwareStatus = "draft"
    FirmwareReleased   FirmwareStatus = "released"
    FirmwareDeprecated FirmwareStatus = "deprecated"
)

type FirmwareUpdate struct {
    ID           string           `json:"id"`
    FirmwareID   string           `json:"firmware_id"`
    DeviceID     string           `json:"device_id"`
    CampaignName string           `json:"campaign_name,omitempty"`
    Status       UpdateStatus     `json:"status"`
    FromVersion  string           `json:"from_version"`
    ToVersion    string           `json:"to_version"`
    Priority     int              `json:"priority"`
    RetryCount   int              `json:"retry_count"`
    MaxRetries   int              `json:"max_retries"`
    StartedAt    *time.Time       `json:"started_at,omitempty"`
    CompletedAt  *time.Time       `json:"completed_at,omitempty"`
    ErrorMessage string           `json:"error_message,omitempty"`
    CreatedAt    time.Time        `json:"created_at"`
}

type UpdateStatus string
const (
    UpdatePending      UpdateStatus = "pending"
    UpdateDownloading  UpdateStatus = "downloading"
    UpdateInstalling   UpdateStatus = "installing"
    UpdateSuccess      UpdateStatus = "success"
    UpdateFailed       UpdateStatus = "failed"
    UpdateRolledBack   UpdateStatus = "rolled_back"
)
```

---

## 4. REST API

### 4.1 Device Registry Endpoints

| Method | Path | Mô tả |
|--------|------|-------|
| `GET`    | `/api/v1/health`                   | Health check |
| `POST`   | `/api/v1/devices`                  | Đăng ký thiết bị mới |
| `GET`    | `/api/v1/devices`                  | Danh sách (filter: `factory_id`, `area`, `status`, `protocol`, `page`, `limit`) |
| `GET`    | `/api/v1/devices/{id}`             | Chi tiết thiết bị |
| `PUT`    | `/api/v1/devices/{id}`             | Cập nhật thiết bị |
| `DELETE` | `/api/v1/devices/{id}`             | Xóa thiết bị |
| `PATCH`  | `/api/v1/devices/{id}/status`      | Cập nhật trạng thái |
| `GET`    | `/api/v1/stats`                    | Thống kê: tổng số, theo factory, theo status |

### 4.2 Device Shadow Endpoints

| Method | Path | Mô tả |
|--------|------|-------|
| `GET`    | `/api/v1/devices/{id}/shadow`             | Xem shadow |
| `PUT`    | `/api/v1/devices/{id}/shadow/desired`     | Cloud → Device |
| `PUT`    | `/api/v1/devices/{id}/shadow/reported`    | Device → Cloud |

### 4.3 Firmware Management Endpoints

| Method | Path | Mô tả |
|--------|------|-------|
| `POST`   | `/api/v1/firmwares`                      | Đăng ký firmware metadata |
| `POST`   | `/api/v1/firmwares/{id}/upload`          | Upload file firmware binary (multipart) |
| `GET`    | `/api/v1/firmwares`                      | Danh sách firmware |
| `GET`    | `/api/v1/firmwares/{id}`                 | Chi tiết firmware |
| `GET`    | `/api/v1/firmwares/{id}/download`        | Tải firmware binary |
| `PUT`    | `/api/v1/firmwares/{id}/release`         | Release firmware (draft → released) |
| `PUT`    | `/api/v1/firmwares/{id}/deprecate`       | Deprecate firmware |
| `DELETE` | `/api/v1/firmwares/{id}`                 | Xóa firmware |
| `POST`   | `/api/v1/firmware-updates`               | Tạo chiến dịch cập nhật (hàng loạt hoặc đơn lẻ) |
| `GET`    | `/api/v1/firmware-updates`               | Danh sách cập nhật (filter: `device_id`, `status`, `campaign_name`) |
| `GET`    | `/api/v1/firmware-updates/{id}`          | Chi tiết cập nhật |
| `POST`   | `/api/v1/firmware-updates/{id}/retry`    | Retry cập nhật thất bại |
| `POST`   | `/api/v1/firmware-updates/{id}/rollback` | Rollback firmware |

### 4.4 Ví dụ Request/Response

**POST /api/v1/firmwares**
```json
// Request
{
  "version": "v2.3.1",
  "model": "CNC-L500",
  "description": "Fix lỗi spindle speed control, cải thiện PID loop"
}
// Response 201
{
  "id": "f1e2d3c4-...",
  "version": "v2.3.1",
  "status": "draft",
  "created_at": "2026-06-14T10:00:00+07:00"
}
```

**POST /api/v1/firmware-updates (chiến dịch hàng loạt)**
```json
// Request
{
  "firmware_id": "f1e2d3c4-...",
  "campaign_name": "CNC-L500-v2.3.1-rollout",
  "device_ids": ["a1b2c3d4-...", "b2c3d4e5-..."],
  "priority": 1
}
// Response 201
{
  "campaign_name": "CNC-L500-v2.3.1-rollout",
  "total_devices": 2,
  "status": "scheduled"
}
```

---

## 5. Luồng Xử Lý Chính

### 5.1 Firmware Upload Flow
```
POST /api/v1/firmwares              → Tạo firmware metadata (status=draft)
POST /api/v1/firmwares/{id}/upload   → Upload file .bin/.hex qua multipart
                                      → Tính SHA256 checksum
                                      → Lưu vào MinIO bucket "firmwares"
                                      → Cập nhật file_name, file_size, checksum
```

### 5.2 OTA Update Flow
```
1. Admin tạo firmware + upload binary → status = "released"
2. Admin tạo firmware-updates (campaign) → chọn devices cần update
3. Device Service:
   → Tạo record firmware_updates cho từng device (status=pending)
   → Publish MQTT: /factory/{id}/device/{id}/fw-update với MinIO download URL
4. Edge Gateway nhận MQTT → thông báo thiết bị
5. Thiết bị download firmware từ MinIO → cài đặt → báo kết quả
6. Device Service nhận kết quả → update firmware_updates.status + devices.firmware_version
```

### 5.3 Rollback Flow
```
POST /api/v1/firmware-updates/{id}/rollback
  → Kiểm tra update hiện tại là "success"
  → Tạo firmware_update mới với to_version = from_version cũ
  → Gửi lệnh rollback qua MQTT
  → Cập nhật status = "rolled_back"
```

### 5.4 Device Shadow Flow
```
PUT /api/v1/devices/{id}/shadow/desired   ← Cloud set desired state
  → Lưu desired_state
  → So sánh với reported_state → tính delta
  → Nếu có delta:
      → Publish MQTT /factory/{id}/device/{id}/desired
      → Redis cache: SET device:{id}:desired <json> (TTL 5m)

PUT /api/v1/devices/{id}/shadow/reported  ← Device reports state
  → Lưu reported_state
  → Redis cache: SET device:{id}:reported <json>
  → Update devices.last_seen_at
  → Nếu có sai khác → cập nhật delta
```

---

## 6. Kế Hoạch Triển Khai (9 bước)

### Bước 1: Khởi tạo project skeleton (30 phút)
- `go.mod` với module `github.com/industrial-iot/device-service`
- `Makefile` (run, build, test, lint, clean, docker-build)
- `Dockerfile` (multi-stage: golang:1.22-alpine → alpine:3.19)
- `migrations/001_init.sql` với 5 bảng
- Cấu trúc thư mục đầy đủ

### Bước 2: Data Models (30 phút)
- `internal/model/device.go` — Device, DeviceStatus, DeviceFilter
- `internal/model/shadow.go` — DeviceShadow, ShadowUpdate
- `internal/model/firmware.go` — Firmware, FirmwareStatus, FirmwareUpdate, UpdateStatus

### Bước 3: Repository Layer (2 giờ)
- `internal/repository/postgres.go`:
  - Kết nối PostgreSQL pool (pgx/v5)
  - Device CRUD: Create, GetByID, List (with filters + pagination), Update, Delete
  - Shadow: GetShadow, UpsertDesired, UpsertReported
  - Connection: UpsertConnection, UpdateHeartbeat
- `internal/repository/firmware_repo.go`:
  - Firmware CRUD: Create, GetByID, List, Update, Delete
  - FirmwareUpdate CRUD: Create, GetByID, List (by device, campaign, status), Update
- `internal/repository/migrations.go` — Chạy 001_init.sql khi start

### Bước 4: MinIO Storage (1 giờ)
- `internal/storage/minio.go`:
  - Kết nối MinIO client (`minio-go/v7`)
  - UploadFirmware(file, fileName) → minioPath
  - DownloadFirmware(minioPath) → io.Reader
  - DeleteFirmware(minioPath)
  - GeneratePresignedURL(minioPath, expiry) → download URL

### Bước 5: Business Logic Layer (2 giờ)
- `internal/service/registry_service.go`:
  - Validate device data (serial_number format, factory_id, etc.)
  - Check unique constraints
  - Auto-create shadow + connection records on register
- `internal/service/shadow_service.go`:
  - ComputeDelta(desired, reported) → delta map
  - SyncDesiredState(deviceID, state) → publish MQTT, cache Redis
  - ProcessReportedState(deviceID, state) → update delta, heartbeat
- `internal/service/firmware_service.go`:
  - Validate firmware upload (checksum, size limit)
  - CreateUpdateCampaign(firmwareID, deviceIDs) → batch create
  - ProcessUpdateResult(updateID, status) → state machine
  - Rollback(updateID) → create reverse update

### Bước 6: REST API — Registry + Shadow (1.5 giờ)
- `internal/api/server.go` — Router setup, middleware chain
- `internal/api/registry_handler.go` — 8 handlers (CRUD + status + stats)
- `internal/api/shadow_handler.go` — 3 handlers
- `internal/api/middleware.go` — Logging, CORS, recovery, requestID

### Bước 7: REST API — Firmware (1.5 giờ)
- `internal/api/firmware_handler.go` — 10 handlers
- Multipart file upload handling
- File download with proper Content-Type + Content-Disposition
- Campaign batch creation

### Bước 8: Wire-up main.go (1 giờ)
- Parse config từ environment variables
- Khởi tạo: PostgreSQL pool → Repository → Services → API Server
- Khởi tạo MinIO client
- Khởi tạo Redis client
- Graceful shutdown với signal handling
- Auto-migration on startup

### Bước 9: Tích hợp Docker (1 giờ)
- Dockerfile (theo pattern gateway/telemetry)
- Thêm device-service + MinIO vào `deployment/docker-compose/dev.yml`
- Test end-to-end với Docker Compose
- Verify tất cả endpoints

---

## 7. Tích hợp Docker Compose

Thêm 2 service mới vào `deployment/docker-compose/dev.yml`:

```yaml
  # Object Storage - Firmware binaries
  minio:
    image: minio/minio:latest
    container_name: iiot-minio
    ports:
      - "9000:9000"   # API
      - "9001:9001"   # Console
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"
    networks:
      - iiot-net
    healthcheck:
      test: ["CMD", "mc", "ready", "local"]
      interval: 10s
      retries: 5

  # Device Service
  device-service:
    build:
      context: ../../device-service
      dockerfile: Dockerfile
    container_name: iiot-device
    ports:
      - "8083:8083"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      minio:
        condition: service_healthy
    environment:
      POSTGRES_HOST: postgres
      POSTGRES_PORT: "5432"
      POSTGRES_DATABASE: industrial_iot
      POSTGRES_USER: iiot
      POSTGRES_PASSWORD: iiot_dev_2024
      REDIS_HOST: redis
      REDIS_PORT: "6379"
      REDIS_PASSWORD: iiot_dev_2024
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
      MINIO_BUCKET: firmwares
      MINIO_USE_SSL: "false"
      API_PORT: "8083"
    networks:
      - iiot-net
    restart: unless-stopped
```

Thêm vào volumes section:
```yaml
  minio-data:
```

---

## 8. Tổng Thời Gian & Độ Phức Tạp

| Bước | Thời gian | Độ phức tạp |
|------|-----------|-------------|
| 1. Project skeleton | 30 phút | ⭐ |
| 2. Data Models | 30 phút | ⭐ |
| 3. Repository Layer | 2 giờ | ⭐⭐ |
| 4. MinIO Storage | 1 giờ | ⭐⭐ |
| 5. Business Logic | 2 giờ | ⭐⭐⭐ |
| 6. API — Registry + Shadow | 1.5 giờ | ⭐⭐⭐ |
| 7. API — Firmware | 1.5 giờ | ⭐⭐⭐ |
| 8. Wire-up main.go | 1 giờ | ⭐⭐ |
| 9. Docker tích hợp | 1 giờ | ⭐⭐ |
| **Tổng** | **~11 giờ** | |

---

## 9. Coding Patterns Thống Nhất

Tuân thủ patterns từ gateway-service và telemetry-service:

- ✅ `go.uber.org/zap` structured logging với `CapitalColorLevelEncoder`
- ✅ `pgx/v5` connection pool cho PostgreSQL
- ✅ `gorilla/mux` cho HTTP routing
- ✅ `minio-go/v7` cho MinIO object storage
- ✅ `go-redis/v9` cho Redis caching
- ✅ Environment variables cho config (`getEnv(key, default)`)
- ✅ Docker multi-stage build (`golang:1.22-alpine` → `alpine:3.19`)
- ✅ Health check endpoint `/api/v1/health`
- ✅ CORS middleware
- ✅ Logging middleware
- ✅ Graceful shutdown với signal handling
- ✅ Container name prefix: `iiot-*`
- ✅ Auto-migration on startup

---

## 10. Goals & Acceptance Criteria

Mỗi goal có test case cụ thể để xác minh hoàn thành. Dùng `curl` hoặc `docker logs` để kiểm tra.

### G1: Device Registry hoạt động đầy đủ

| ID | Test | Expected Result |
|----|------|-----------------|
| G1.1 | `POST /api/v1/devices` với payload hợp lệ | `201 Created` + device JSON có `id`, `status="offline"` |
| G1.2 | `POST /api/v1/devices` serial_number trùng | `409 Conflict` + `{"error": "serial_number already exists"}` |
| G1.3 | `POST /api/v1/devices` thiếu trường bắt buộc | `400 Bad Request` + message chi tiết |
| G1.4 | `GET /api/v1/devices` (không filter) | `200 OK` + `{"count": N, "devices": [...]}` |
| G1.5 | `GET /api/v1/devices?factory_id=HN&status=online` | `200 OK` + chỉ trả về device khớp filter |
| G1.6 | `GET /api/v1/devices?page=1&limit=10` | `200 OK` + phân trang đúng |
| G1.7 | `GET /api/v1/devices/{id}` | `200 OK` + chi tiết thiết bị |
| G1.8 | `GET /api/v1/devices/{non-existent-id}` | `404 Not Found` |
| G1.9 | `PUT /api/v1/devices/{id}` cập nhật name | `200 OK` + device với name mới |
| G1.10 | `DELETE /api/v1/devices/{id}` | `204 No Content`, shadow + connections bị cascade xóa |
| G1.11 | `PATCH /api/v1/devices/{id}/status` `{"status":"online"}` | `200 OK` + `status="online"`, `last_seen_at` được cập nhật |
| G1.12 | `GET /api/v1/stats` | `200 OK` + `{"total": N, "by_factory": {...}, "by_status": {...}}` |

### G2: Device Shadow hoạt động đúng

| ID | Test | Expected Result |
|----|------|-----------------|
| G2.1 | `GET /api/v1/devices/{id}/shadow` (mới tạo) | `200 OK` + `desired={}`, `reported={}`, `delta={}`, `version=1` |
| G2.2 | `PUT .../shadow/desired` `{"temp": 200}` | `200 OK` + `desired={"temp":200}`, `delta={"temp":200}`, version tăng |
| G2.3 | `PUT .../shadow/reported` `{"temp": 195}` | `200 OK` + `reported={"temp":195}`, `delta={"temp":200}` (còn khác) |
| G2.4 | `PUT .../shadow/reported` `{"temp": 200}` | `200 OK` + `delta={}` (không còn khác biệt) |
| G2.5 | Redis cache sau desired update | `GET device:{id}:desired` trong Redis có giá trị đúng |

### G3: Firmware Management hoạt động

| ID | Test | Expected Result |
|----|------|-----------------|
| G3.1 | `POST /api/v1/firmwares` metadata | `201 Created` + `status="draft"` |
| G3.2 | `POST /api/v1/firmwares/{id}/upload` file `.bin` | `200 OK` + `file_size`, `checksum_sha256` được cập nhật, file có trong MinIO |
| G3.3 | `GET /api/v1/firmwares/{id}/download` | `200 OK` + Content-Type `application/octet-stream`, file đúng checksum |
| G3.4 | `PUT /api/v1/firmwares/{id}/release` | `200 OK` + `status="released"`, `released_at` != null |
| G3.5 | `PUT /api/v1/firmwares/{id}/deprecate` | `200 OK` + `status="deprecated"` |
| G3.6 | `DELETE /api/v1/firmwares/{id}` | `204 No Content`, file trong MinIO cũng bị xóa |
| G3.7 | Upload file > 100MB | `413 Request Entity Too Large` |

### G4: OTA Firmware Update hoạt động

| ID | Test | Expected Result |
|----|------|-----------------|
| G4.1 | `POST /api/v1/firmware-updates` (đơn lẻ) | `201 Created` + `status="pending"` |
| G4.2 | `POST /api/v1/firmware-updates` (campaign 5 devices) | `201 Created` + `campaign_name`, 5 bản ghi được tạo |
| G4.3 | Edge báo `status="downloading"` → `"installing"` → `"success"` | Mỗi lần PUT update status → response đúng |
| G4.4 | Device firmware_version được cập nhật sau success | `GET /api/v1/devices/{id}` → `firmware_version="v2.3.1"` |
| G4.5 | `POST .../retry` với update failed | `200 OK` + `retry_count` tăng, `status="pending"` |
| G4.6 | `POST .../rollback` với update success | `200 OK` + `status="rolled_back"`, firmware_version về cũ |
| G4.7 | Retry quá `max_retries` | `status="failed"`, không retry thêm được |

### G5: Tích hợp hệ thống

| ID | Test | Expected Result |
|----|------|-----------------|
| G5.1 | `GET /api/v1/health` | `200 OK` + `{"status":"healthy","service":"device-service"}` |
| G5.2 | `docker compose ... up -d` | Container `iiot-device`, `iiot-minio` chạy healthy |
| G5.3 | Migration tự động khi start | Log: `"Migration 001_init.sql applied successfully"` |
| G5.4 | Gọi API từ host qua `localhost:8083` | Kết nối và response thành công |
| G5.5 | Container restart → data không mất | Device đã tạo trước restart vẫn query được |
| G5.6 | Service log đúng format | Log có timestamp, level, message, fields (zap format) |

### G6: Chất lượng code

| ID | Test | Expected Result |
|----|------|-----------------|
| G6.1 | `go test ./...` | Tất cả test pass, coverage ≥ 70% |
| G6.2 | `go vet ./...` | Không có warning |
| G6.3 | Code pattern khớp gateway/telemetry service | getEnv(), zap logger, pgx pool, graceful shutdown |
| G6.4 | Không hardcode config | Tất cả config qua env vars, có default value |

---

## 11. Tiêu Chí Done Tổng

- [ ] **G1** — Device Registry: 12/12 test pass
- [ ] **G2** — Device Shadow: 5/5 test pass
- [ ] **G3** — Firmware Management: 7/7 test pass
- [ ] **G4** — OTA Firmware Update: 7/7 test pass
- [ ] **G5** — Tích hợp hệ thống: 6/6 test pass
- [ ] **G6** — Chất lượng code: 4/4 test pass
- [ ] **Tổng: 41/41 acceptance criteria hoàn thành**
