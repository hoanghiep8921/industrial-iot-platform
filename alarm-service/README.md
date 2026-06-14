# Alarm Service

> Cảnh báo thông minh — Rule Engine + Alarm Lifecycle Management

## 🎯 Chức Năng

- **Rule Engine**: Đánh giá telemetry data theo các rule định nghĩa trước
- **Alarm State Machine**: Quản lý vòng đời alarm (raised → acknowledged → resolved)
- **Escalation Policy**: Tự động escalate alarm chưa được xử lý
- **Alarm Correlation**: Gộp các alarm liên quan để tránh alarm storm
- **Scheduled Suppression**: Tạm tắt alarm theo lịch bảo trì

## 📁 Cấu Trúc

```
alarm-service/
├── rule-engine/             # Đánh giá điều kiện alarm
│   ├── rule-definition/     #   Định nghĩa rule (threshold, trend, state)
│   ├── rule-evaluator/      #   Engine đánh giá rule real-time
│   ├── window-aggregator/   #   Aggregation theo time window
│   └── expression-parser/   #   Parse biểu thức điều kiện
├── alarm-state-machine/     # Vòng đời alarm
│   ├── state-handler/       #   Xử lý chuyển trạng thái
│   ├── ack-manager/         #   Acknowledge alarm
│   └── auto-resolver/       #   Tự động resolve khi điều kiện hết
└── escalation-policy/       # Chính sách escalation
    ├── level-definition/    #   Các mức alarm (Warning/Critical/Emergency)
    ├── escalation-chain/    #   Chuỗi escalation theo thời gian
    └── schedule-manager/    #   Lịch on-call & suppression
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go hoặc Java (Spring Boot)
- **Database**: PostgreSQL
- **Message Queue**: Kafka (nhận telemetry stream, publish alarm events)
- **Cache**: Redis (real-time alarm state)

## 📊 Alarm Rule Types

```yaml
rules:
  - name: "Nhiệt độ lò quá cao"
    type: threshold
    metric: furnace.temperature
    condition: "> 800"
    duration: 30s         # Phải vượt ngưỡng liên tục 30s
    severity: critical
    group: furnace_hn_01

  - name: "Áp suất tăng đột ngột"
    type: rate_of_change
    metric: boiler.pressure
    condition: "delta > 50 over 10s"
    severity: warning

  - name: "Motor rung bất thường"
    type: anomaly
    metric: motor.vibration
    algorithm: "3-sigma"
    training_window: 7d
    severity: warning

  - name: "Máy ngừng hoạt động"
    type: state_change
    metric: machine.status
    condition: "RUNNING → STOPPED"
    severity: emergency
```

## 🔄 Alarm Lifecycle

```
    [Rule Engine phát hiện điều kiện]
                │
         ┌──────▼──────┐
         │   RAISED    │  ← Alarm được tạo, gửi notification
         └──────┬──────┘
                │ Operator acknowledges
         ┌──────▼──────┐
         │ ACKNOWLEDGED│  ← Đã có người nhận xử lý
         └──────┬──────┘
                │ Condition clears
         ┌──────▼──────┐
         │  RESOLVED   │  ← Alarm kết thúc
         └─────────────┘

    RAISED ──(timeout)──► ESCALATED ──► Gửi cho supervisor
```

---

## 🧪 Hướng Dẫn Kiểm Thử Tích Hợp (Integration Testing Guide)

Để kiểm thử hoạt động liên thông của Alarm Service (HTTP REST API, Kafka Telemetry Ingestion, Redis active state cache, Postgres database, và Kafka Alarm Event Publisher), thực hiện các bước sau:

### Bước 1: Tạo Rule Cảnh Báo
Gửi yêu cầu POST tạo rule nhiệt độ quá cao (> 100°C) với thời gian trì hoãn 5 giây (`duration_seconds`):
```bash
curl.exe -X POST -H "Content-Type: application/json" -d "{\"name\":\"Nhiet do qua cao Test\",\"description\":\"Canh bao khi nhiet do lo vuot 100 do C\",\"type\":\"threshold\",\"metric_name\":\"temperature\",\"condition_operator\":\"\u003e\",\"condition_value\":100,\"duration_seconds\":5,\"severity\":\"critical\",\"factory_id\":\"factory-01\",\"is_enabled\":true}" http://localhost:8086/api/v1/alarms/rules
```
*Ghi lại ID của rule mới được trả về trong JSON phản hồi (ví dụ: `16ddb96a-563f-4cdb-8968-980590ce927e`).*

### Bước 2: Giả lập Telemetry vi phạm ngưỡng (Bắt đầu vi phạm)
Gửi bản tin đo lường nhiệt độ 105.5°C (> 100) qua Kafka topic `telemetry.raw`:
```bash
echo '{"factory_id": "factory-01", "area": "machining", "machine_id": "machine-01", "metric": "temperature", "topic": "devices/machine-01/telemetry", "timestamp": "2026-06-14T18:05:00Z", "value": 105.5}' | docker exec -i iiot-kafka kafka-console-producer --bootstrap-server localhost:9092 --topic telemetry.raw
```
*(Trạng thái vi phạm bắt đầu được lưu vết vào Redis với thời gian bắt đầu lỗi là `18:05:00Z`, nhưng chưa tạo alarm ngay vì chưa vượt quá 5s trì hoãn).*

### Bước 3: Giả lập Telemetry vi phạm ngưỡng lần 2 (Đủ điều kiện tạo Alarm)
Gửi tiếp bản tin đo lường tại thời điểm sau đó 6 giây (vượt qua trì hoãn 5s):
```bash
echo '{"factory_id": "factory-01", "area": "machining", "machine_id": "machine-01", "metric": "temperature", "topic": "devices/machine-01/telemetry", "timestamp": "2026-06-14T18:05:06Z", "value": 105.5}' | docker exec -i iiot-kafka kafka-console-producer --bootstrap-server localhost:9092 --topic telemetry.raw
```
*(Hệ thống sẽ tự động tạo cảnh báo `alarm.raised` và đẩy lên Kafka topic `alarm.events` đồng thời lưu vào Postgres).*

### Bước 4: Kiểm tra Cảnh báo đang hoạt động
Kiểm tra danh sách Active Alarm qua API:
```bash
curl.exe -s http://localhost:8086/api/v1/alarms/active
```
*Bạn sẽ nhận về thông tin cảnh báo ở trạng thái `RAISED`. Hãy lấy giá trị `id` của alarm này.*

### Bước 5: Xác nhận (Acknowledge) Cảnh báo
Gửi yêu cầu xác nhận cảnh báo:
```bash
curl.exe -X POST -H "Content-Type: application/json" -d "{\"ack_by\":\"admin\"}" http://localhost:8086/api/v1/alarms/<ALARM_ID>/acknowledge
```

### Bước 6: Giả lập Telemetry bình thường (Tự động phục hồi)
Gửi dữ liệu đo lường nhiệt độ bình thường trở lại (ví dụ 80.0°C <= 100):
```bash
echo '{"factory_id": "factory-01", "area": "machining", "machine_id": "machine-01", "metric": "temperature", "topic": "devices/machine-01/telemetry", "timestamp": "2026-06-14T18:05:10Z", "value": 80.0}' | docker exec -i iiot-kafka kafka-console-producer --bootstrap-server localhost:9092 --topic telemetry.raw
```
*(Hệ thống sẽ tự phục hồi trạng thái cảnh báo sang `RESOLVED` trong DB).*

### Bước 7: Kiểm tra Sự kiện bắn ra trên Kafka
Đọc các sự kiện thay đổi trạng thái (`alarm.raised`, `alarm.resolved`) đã được đẩy lên Kafka:
```bash
docker exec iiot-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic alarm.events --from-beginning --max-messages 5
```

