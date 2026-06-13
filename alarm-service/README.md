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
