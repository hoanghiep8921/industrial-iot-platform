# Device Service — Goals & Acceptance Criteria

> Dùng để đánh giá mức độ hoàn thành. Mỗi goal có test case cụ thể kèm lệnh `curl` hoặc `docker logs` để xác minh.

---

## G1: Device Registry hoạt động đầy đủ (12 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G1.1 | `curl -X POST localhost:8083/api/v1/devices -H 'Content-Type: application/json' -d '{"serial_number":"SN-TEST-001","name":"Test Device","factory_id":"HN"}'` | `201 Created` + JSON có `id`, `status="offline"` |
| G1.2 | Chạy lại G1.1 với cùng serial_number | `409 Conflict` + `{"error": "serial_number already exists"}` |
| G1.3 | `curl -X POST localhost:8083/api/v1/devices -H 'Content-Type: application/json' -d '{"name":"Missing SN"}'` | `400 Bad Request` + message lỗi chi tiết |
| G1.4 | `curl localhost:8083/api/v1/devices` | `200 OK` + `{"count": N, "devices": [...]}` |
| G1.5 | `curl 'localhost:8083/api/v1/devices?factory_id=HN&status=online'` | `200 OK` + chỉ trả về device khớp filter |
| G1.6 | `curl 'localhost:8083/api/v1/devices?page=1&limit=10'` | `200 OK` + phân trang đúng, `count` là tổng số |
| G1.7 | `curl localhost:8083/api/v1/devices/{id}` | `200 OK` + chi tiết thiết bị |
| G1.8 | `curl localhost:8083/api/v1/devices/00000000-0000-0000-0000-000000000000` | `404 Not Found` |
| G1.9 | `curl -X PUT localhost:8083/api/v1/devices/{id} -H 'Content-Type: application/json' -d '{"name":"Updated Name"}'` | `200 OK` + `name="Updated Name"` |
| G1.10 | `curl -X DELETE localhost:8083/api/v1/devices/{id}` | `204 No Content`, shadow + connections bị cascade |
| G1.11 | `curl -X PATCH localhost:8083/api/v1/devices/{id}/status -H 'Content-Type: application/json' -d '{"status":"online"}'` | `200 OK` + `status="online"`, `last_seen_at` được cập nhật |
| G1.12 | `curl localhost:8083/api/v1/stats` | `200 OK` + `{"total": N, "by_factory": {...}, "by_status": {...}}` |

---

## G2: Device Shadow hoạt động đúng (5 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G2.1 | `curl localhost:8083/api/v1/devices/{id}/shadow` (thiết bị mới tạo) | `200 OK` + `desired={}`, `reported={}`, `delta={}`, `version=1` |
| G2.2 | `curl -X PUT localhost:8083/api/v1/devices/{id}/shadow/desired -H 'Content-Type: application/json' -d '{"temp":200}'` | `200 OK` + `desired={"temp":200}`, `delta={"temp":200}`, version tăng lên 2 |
| G2.3 | `curl -X PUT localhost:8083/api/v1/devices/{id}/shadow/reported -H 'Content-Type: application/json' -d '{"temp":195}'` | `200 OK` + `reported={"temp":195}`, `delta={"temp":200}` (còn khác biệt) |
| G2.4 | `curl -X PUT localhost:8083/api/v1/devices/{id}/shadow/reported -H 'Content-Type: application/json' -d '{"temp":200}'` | `200 OK` + `delta={}` (không còn khác biệt) |
| G2.5 | `docker exec iiot-redis redis-cli -a iiot_dev_2024 GET device:{id}:desired` | Redis có cache giá trị đúng |

---

## G3: Firmware Management hoạt động (7 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G3.1 | `curl -X POST localhost:8083/api/v1/firmwares -H 'Content-Type: application/json' -d '{"version":"v1.0.0","model":"CNC-L500","description":"Initial release"}'` | `201 Created` + `status="draft"` |
| G3.2 | `curl -X POST localhost:8083/api/v1/firmwares/{id}/upload -F 'file=@test.bin'` | `200 OK` + `file_size` và `checksum_sha256` được điền, file có trong MinIO |
| G3.3 | `curl -O localhost:8083/api/v1/firmwares/{id}/download` | `200 OK` + file tải về đúng checksum |
| G3.4 | `curl -X PUT localhost:8083/api/v1/firmwares/{id}/release` | `200 OK` + `status="released"`, `released_at` khác null |
| G3.5 | `curl -X PUT localhost:8083/api/v1/firmwares/{id}/deprecate` | `200 OK` + `status="deprecated"` |
| G3.6 | `curl -X DELETE localhost:8083/api/v1/firmwares/{id}` | `204 No Content`, file trong MinIO cũng bị xóa |
| G3.7 | Upload file > 100MB lên G3.2 | `413 Request Entity Too Large` |

---

## G4: OTA Firmware Update hoạt động (7 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G4.1 | `curl -X POST localhost:8083/api/v1/firmware-updates -H 'Content-Type: application/json' -d '{"firmware_id":"...","device_id":"..."}'` | `201 Created` + `status="pending"` |
| G4.2 | `curl -X POST localhost:8083/api/v1/firmware-updates -H 'Content-Type: application/json' -d '{"firmware_id":"...","device_ids":["id1","id2","id3","id4","id5"],"campaign_name":"test-rollout"}'` | `201 Created` + `campaign_name="test-rollout"`, 5 bản ghi được tạo |
| G4.3 | Lần lượt PATCH status: `"downloading"` → `"installing"` → `"success"` | Mỗi lần PUT → response đúng status mới |
| G4.4 | `curl localhost:8083/api/v1/devices/{id}` sau update success | `firmware_version` được cập nhật thành version mới |
| G4.5 | `curl -X POST localhost:8083/api/v1/firmware-updates/{id}/retry` với update đã failed | `200 OK` + `retry_count` tăng 1, `status="pending"` |
| G4.6 | `curl -X POST localhost:8083/api/v1/firmware-updates/{id}/rollback` với update đã success | `200 OK` + `status="rolled_back"`, firmware_version về version cũ |
| G4.7 | Retry đến khi vượt `max_retries` | `status="failed"`, không retry thêm được nữa |

---

## G5: Tích hợp hệ thống (6 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G5.1 | `curl localhost:8083/api/v1/health` | `200 OK` + `{"status":"healthy","service":"device-service"}` |
| G5.2 | `docker compose -f deployment/docker-compose/dev.yml up -d` | Container `iiot-device` và `iiot-minio` đều healthy |
| G5.3 | `docker logs iiot-device \| grep -i migration` | Log: `"Migration 001_init.sql applied successfully"` |
| G5.4 | `curl localhost:8083/api/v1/health` từ host machine | Kết nối thành công qua port mapping |
| G5.5 | `docker compose restart device-service` rồi query danh sách devices | Data không mất, device đã tạo vẫn query được |
| G5.6 | `docker logs iiot-device --tail 5` | Log đúng zap format: timestamp + level + message + fields |

---

## G6: Chất lượng code (4 tests)

| ID | Lệnh kiểm tra | Expected Result |
|----|---------------|-----------------|
| G6.1 | `cd device-service && go test ./... -cover` | Tất cả test pass, coverage ≥ 70% |
| G6.2 | `cd device-service && go vet ./...` | Không có warning nào |
| G6.3 | Review thủ công: so sánh với gateway-service và telemetry-service | Cùng pattern: `getEnv()`, `zap.Logger`, `pgxpool`, graceful shutdown |
| G6.4 | `grep -r 'localhost\|9092\|5432\|6379' device-service/cmd/` (tìm hardcode) | Không có hardcode — tất cả config qua env vars kèm default |

---

## Tổng kết

| Goal | Mô tả | Tests |
|------|-------|:----:|
| **G1** | Device Registry | 12 |
| **G2** | Device Shadow | 5 |
| **G3** | Firmware Management | 7 |
| **G4** | OTA Firmware Update | 7 |
| **G5** | Tích hợp hệ thống | 6 |
| **G6** | Chất lượng code | 4 |
| | **Tổng cộng** | **41** |

### Cách đánh giá

```
Hoàn thành = 41/41 test pass
Đạt yêu cầu tối thiểu = G1 + G2 + G5 + G6 (27/27)
Có thể demo = G1 + G2 + G5 (23/23)
```
