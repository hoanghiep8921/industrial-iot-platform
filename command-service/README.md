# Command Service

> Điều khiển máy móc từ cloud — Nhận lệnh từ UI, gửi xuống Edge Gateway

## 🎯 Chức Năng

- **Command Dispatch**: Nhận lệnh điều khiển từ Web/Mobile UI, gửi qua MQTT xuống edge
- **Command Tracking**: Theo dõi trạng thái từng lệnh (pending → sent → executing → done/failed)
- **Command History**: Audit log toàn bộ lệnh — ai gửi, khi nào, kết quả ra sao
- **Command Validation**: Kiểm tra quyền và tính hợp lệ của lệnh trước khi gửi
- **Timeout & Retry**: Timeout handling và automatic retry khi lệnh thất bại
- **Batch Commands**: Hỗ trợ gửi lệnh hàng loạt (vd: dừng tất cả máy trong khu vực)

## 📁 Cấu Trúc

```
command-service/
├── command-dispatcher/      # Điều phối lệnh
│   ├── api/                 #   REST/gRPC API nhận lệnh từ UI
│   ├── validator/           #   Validate lệnh (quyền, tham số, trạng thái thiết bị)
│   ├── mqtt-publisher/      #   Gửi lệnh qua MQTT
│   └── batch-handler/       #   Xử lý batch commands
├── command-tracker/         # Theo dõi trạng thái lệnh
│   ├── state-machine/       #   Trạng thái: pending→sent→executing→done/failed
│   ├── timeout-watchdog/    #   Phát hiện lệnh timeout
│   └── retry-manager/       #   Retry với exponential backoff
└── command-history/         # Lịch sử & Audit
    ├── repository/          #   Lưu trữ lịch sử lệnh
    ├── audit-log/           #   Audit trail cho compliance
    └── report-generator/    #   Báo cáo thống kê lệnh
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go
- **Database**: PostgreSQL (command history)
- **Message Queue**: Kafka (command events) + MQTT (device communication)
- **Cache**: Redis (real-time command state)

## 📊 Command State Machine

```
                    ┌──────────┐
                    │  PENDING │  ← Lệnh được tạo từ UI
                    └────┬─────┘
                         │ dispatch
                    ┌────▼─────┐
                    │   SENT   │  ← Đã gửi MQTT xuống edge
                    └────┬─────┘
                         │ edge nhận
                    ┌────▼─────┐
                    │EXECUTING │  ← Edge đang thực thi
                    └──┬───┬───┘
                       │   │
                  ┌────▼┐ ┌▼────┐
                  │DONE │ │FAIL │  ← Kết quả cuối
                  └─────┘ └─────┘

         Timeout: SENT/EXECUTING → FAILED (sau N giây)
```

## 🔒 Security

- Mọi lệnh điều khiển đều được ghi audit log
- Xác thực user + RBAC trước khi dispatch
- Command confirmation dialog trên UI cho lệnh nguy hiểm
- Rate limiting: tối đa N lệnh/giây/user
