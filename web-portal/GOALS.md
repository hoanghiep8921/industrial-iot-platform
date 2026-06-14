# Web Portal — Goals & Acceptance Criteria

> Dùng để đánh giá mức độ hoàn thành. Mỗi goal có test case cụ thể để xác minh.

---

## G1: Project hoạt động ở dev mode (5 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G1.1 | `cd web-portal && pnpm install && pnpm dev` | Dev server chạy tại http://localhost:5173 |
| G1.2 | Mở browser `http://localhost:5173/` | Hiển thị Dashboard page (không crash, không console error) |
| G1.3 | Mở `http://localhost:5173/devices` | Hiển thị Device List page |
| G1.4 | `curl http://localhost:5173/api/v1/stats` | Trả về device stats JSON (proxy hoạt động) |
| G1.5 | `curl http://localhost:5173/api/v1/telemetry/recent?limit=5` | Trả về telemetry JSON (proxy hoạt động) |

---

## G2: Dashboard page hoạt động (5 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G2.1 | Vào `/` khi backend đang chạy | Hiển thị 4 stats cards: Total Devices, Online, Offline, Errors |
| G2.2 | Stats cards hiển thị đúng số liệu | Số khớp với `curl http://localhost:8085/api/v1/stats` |
| G2.3 | TelemetryChart có dữ liệu | Line chart hiển thị các điểm dữ liệu từ recent telemetry |
| G2.4 | DeviceStatusCards hiển thị | Mỗi card: device name + latest metric value + status color |
| G2.5 | Auto-refresh hoạt động | Chờ 15s → chart + cards tự cập nhật dữ liệu mới |

---

## G3: Device List page hoạt động (6 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G3.1 | Vào `/devices` | Hiển thị bảng danh sách devices (có thể rỗng) |
| G3.2 | Click "Add Device" → điền form → Submit | Device mới xuất hiện trong bảng, status=offline |
| G3.3 | Chọn filter factory | Bảng chỉ hiện devices của factory đó |
| G3.4 | Chọn filter status chips (online/offline) | Bảng filter đúng |
| G3.5 | Gõ vào search box (tên hoặc serial) | Bảng lọc theo text |
| G3.6 | Click vào 1 row | Navigate đến `/devices/:id` |

---

## G4: Device Detail page hoạt động (4 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G4.1 | Vào `/devices/:id` | DeviceInfoCard hiển thị name, model, serial, factory, status |
| G4.2 | TelemetryHistoryChart | Line chart hiển thị data 24h của device (hoặc empty state nếu không có) |
| G4.3 | ShadowStateViewer | Hiển thị desired_state vs reported_state dạng JSON |
| G4.4 | FirmwareInfoCard | Hiển thị firmware version hiện tại |

---

## G5: Error & Edge Cases (4 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G5.1 | Tắt backend → reload trang | Hiển thị ErrorAlert với message lỗi + nút Retry |
| G5.2 | Vào `/devices/00000000-0000-0000-0000-000000000000` | Hiển thị error state (device not found) |
| G5.3 | Bấm Retry sau khi backend bật lại | Data load lại thành công |
| G5.4 | Dashboard khi chưa có device nào | Hiển thị EmptyState "No devices registered yet" |

---

## G6: Docker & Production (4 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G6.1 | `cd web-portal && docker build -t iiot-web-portal .` | Build thành công |
| G6.2 | `docker run -d -p 8090:80 iiot-web-portal && curl http://localhost:8090/` | Trả về HTML SPA |
| G6.3 | `curl http://localhost:8090/api/v1/health` qua container | Proxy hoạt động (qua nginx → device-service) |
| G6.4 | `docker compose -f deployment/docker-compose/dev.yml up -d web-portal` | Container `iiot-web-portal` chạy healthy |

---

## G7: Chất lượng code (3 tests)

| ID | Test | Expected Result |
|----|------|-----------------|
| G7.1 | `pnpm typecheck` (tsc --noEmit) | Không có lỗi TypeScript |
| G7.2 | `pnpm lint` | Không có warning |
| G7.3 | `pnpm build` | Build thành công, tạo thư mục `dist/` |

---

## Tổng kết

| Goal | Mô tả | Tests |
|------|-------|:----:|
| **G1** | Dev mode hoạt động | 5 |
| **G2** | Dashboard page | 5 |
| **G3** | Device List page | 6 |
| **G4** | Device Detail page | 4 |
| **G5** | Error & Edge Cases | 4 |
| **G6** | Docker & Production | 4 |
| **G7** | Chất lượng code | 3 |
| | **Tổng** | **31** |

### Cách đánh giá

```
Hoàn thành tuyệt đối  = 31/31 test pass
Đạt yêu cầu tối thiểu = G1 + G2 + G3 + G4 + G6 + G7 (27/27)
Có thể demo           = G1 + G2 + G3 + G4 (20/20)
```
