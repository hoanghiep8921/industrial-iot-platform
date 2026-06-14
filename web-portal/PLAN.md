# Web Portal — Kế hoạch triển khai

> **Mục tiêu**: Dashboard React giám sát & quản lý thiết bị IoT
> **Thời gian**: 3-4 ngày | **Trang**: 3 (Dashboard, Devices, Device Detail)

---

## Context

- **telemetry-service** (:8082) và **device-service** (:8085) đã có REST API đầy đủ + CORS
- **web-portal/** trống — chỉ có README spec
- Web portal gọi thẳng backend API qua Vite proxy (không cần BFF/API Gateway ngay)
- Backend trả về 2 kiểu JSON key khác nhau: telemetry dùng PascalCase (`Time`, `DeviceID`), device dùng snake_case (`serial_number`)

---

## Tech Stack

| Lớp | Công nghệ |
|-----|-----------|
| Framework | React 19 + TypeScript + Vite |
| UI | MUI v6 |
| Charts | Recharts |
| State | Zustand |
| Package Manager | pnpm |
| Container | node:20-alpine → nginx:alpine (multi-stage) |

---

## Cấu trúc thư mục

```
web-portal/
├── src/
│   ├── api/
│   │   ├── client.ts              # fetch wrapper + error handling
│   │   ├── telemetryApi.ts        # PascalCase → camelCase
│   │   └── deviceApi.ts           # snake_case → camelCase
│   ├── stores/
│   │   ├── useTelemetryStore.ts   # telemetry data + actions
│   │   └── useDeviceStore.ts      # devices, shadow, firmware
│   ├── pages/
│   │   ├── DashboardPage.tsx      # Overview + chart + device cards
│   │   ├── DeviceListPage.tsx     # Table + filters + add/edit form
│   │   └── DeviceDetailPage.tsx   # Info + chart + shadow + firmware
│   ├── components/
│   │   ├── layout/  (AppLayout, Sidebar, TopBar)
│   │   ├── dashboard/  (StatsOverview, TelemetryChart, DeviceStatusCards)
│   │   ├── devices/  (DeviceTable, DeviceFilters, DeviceFormDialog, DeviceStatusBadge)
│   │   ├── device-detail/  (DeviceInfoCard, TelemetryHistoryChart, ShadowStateViewer, FirmwareInfoCard)
│   │   └── common/  (LoadingSpinner, ErrorAlert, EmptyState)
│   ├── hooks/usePolling.ts       # Polling every 10s
│   ├── types/  (telemetry.ts, device.ts, shadow.ts, firmware.ts)
│   ├── utils/  (transformKeys.ts: pascalToCamel + snakeToCamel)
│   ├── theme.ts, App.tsx, main.tsx
├── vite.config.ts                 # React plugin + proxy config
├── Dockerfile, nginx.conf, Makefile
```

---

## Data Flow

```
Backend API → fetch() → transformKeys (camelCase) → Zustand store → React component
                                                         ↑
                                              usePolling (10s interval)
```

### Vite Proxy (dev mode)

```
/api/v1/telemetry/*   → localhost:8082
/api/v1/devices/*     → localhost:8085
/api/v1/firmwares/*   → localhost:8085
/api/v1/firmware-updates/* → localhost:8085
/api/v1/*             → localhost:8085 (stats, health)
```

### Nginx Proxy (production container)

```
/api/v1/telemetry → telemetry-service:8082
/api/v1/          → device-service:8083
/                 → SPA (index.html fallback)
```

---

## 3 Trang MVP

### 1. Dashboard (`/`)
- **StatsOverview**: 4 Card (total devices, online, offline, errors) từ `GET /api/v1/stats`
- **TelemetryChart**: Recharts LineChart — `GET /api/v1/telemetry/recent?limit=100`, metric selector dropdown, auto-poll 10s
- **DeviceStatusCards**: Grid Card cho từng device — `GET /api/v1/telemetry/latest`, hiển thị metric mới nhất + status badge

### 2. Device List (`/devices`)
- **DeviceFilters**: factory dropdown, status chips, search box, "Add Device" button
- **DeviceTable**: MUI DataGrid — columns: name, serial, factory, area, status, last_seen
  - Row click → navigate `/devices/:id`
- **DeviceFormDialog**: Create/Edit form (POST/PUT)

### 3. Device Detail (`/devices/:id`)
- **DeviceInfoCard**: metadata hiển thị dạng field-value
- **TelemetryHistoryChart**: time-series 24h — `GET /api/v1/telemetry/device/{id}?from=&to=`, time range picker
- **ShadowStateViewer**: JSON diff desired vs reported + delta highlight
- **FirmwareInfoCard**: version hiện tại + danh sách firmware + update status

---

## Triển khai 8 bước

| # | Bước | Thời gian |
|---|------|:--:|
| 1 | Project scaffold: Vite + MUI + Router + theme + Dockerfile | 1h |
| 2 | Types + API layer + key transformers | 1.5h |
| 3 | Zustand stores (telemetry + device) + usePolling hook | 1h |
| 4 | Layout (Sidebar, TopBar, AppLayout) + common components | 1h |
| 5 | Dashboard page (StatsOverview + Chart + DeviceCards) | 2h |
| 6 | Device List page (Filters + DataGrid + FormDialog) | 2h |
| 7 | Device Detail page (Info + Chart + Shadow + Firmware) | 2h |
| 8 | Docker + nginx + docker-compose integration | 1h |
| | **Tổng** | **~11.5h** |

---

## Docker Integration

```yaml
# deployment/docker-compose/dev.yml (thêm)
web-portal:
  build:
    context: ../../web-portal
    dockerfile: Dockerfile
  container_name: iiot-web-portal
  ports:
    - "8090:80"
  depends_on:
    - telemetry-service
    - device-service
  networks:
    - iiot-net
  restart: unless-stopped
```

---

## Verification

```bash
# 1. Backends alive
curl http://localhost:8082/api/v1/health   # {"status":"healthy"}
curl http://localhost:8085/api/v1/health   # {"status":"healthy"}

# 2. Dev server
cd web-portal && pnpm dev                  # http://localhost:5173

# 3. Proxy works
curl http://localhost:5173/api/v1/stats    # device stats JSON
curl http://localhost:5173/api/v1/telemetry/recent?limit=5  # telemetry JSON

# 4. Docker
docker compose -f deployment/docker-compose/dev.yml up -d web-portal
curl http://localhost:8090                 # SPA HTML
curl http://localhost:8090/api/v1/health   # proxied to device-service

# 5. UI tests
# - Dashboard hiển thị stats + chart + device cards
# - /devices: thêm/xóa/lọc/search thiết bị
# - /devices/:id: xem info + telemetry chart + shadow + firmware
```
