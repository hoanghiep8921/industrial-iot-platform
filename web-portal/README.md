# Web Portal

> Giao diện web giám sát & điều khiển nhà máy công nghiệp

## 🎯 Chức Năng

- **Dashboard**: Real-time overview toàn bộ nhà máy với biểu đồ, widget
- **Device Control**: Giao diện điều khiển máy móc, gửi lệnh, xem trạng thái
- **Alarm Center**: Trung tâm cảnh báo — xem, acknowledge, xử lý alarm
- **Reports**: Báo cáo sản xuất, OEE, tiêu thụ năng lượng, xu hướng
- **Admin**: Quản lý user, phân quyền, cấu hình hệ thống

## 📁 Cấu Trúc

```
web-portal/
├── dashboard/               # Real-time dashboard
│   ├── factory-overview/    #   Tổng quan toàn nhà máy
│   ├── line-detail/         #   Chi tiết từng dây chuyền
│   ├── machine-detail/      #   Chi tiết từng máy
│   └── kpi-widgets/         #   Widgets OEE, throughput, energy
├── device-control/          # Giao diện điều khiển
│   ├── device-tree/         #   Cây thiết bị theo nhà máy/khu vực
│   ├── control-panel/       #   Bảng điều khiển máy
│   ├── command-history/     #   Lịch sử lệnh điều khiển
│   └── batch-control/       #   Điều khiển hàng loạt
├── alarm-center/            # Trung tâm cảnh báo
│   ├── active-alarms/       #   Alarm đang active
│   ├── alarm-history/       #   Lịch sử alarm
│   ├── alarm-config/        #   Cấu hình rule alarm
│   └── escalation-view/     #   Xem trạng thái escalation
├── reports/                 # Báo cáo
│   ├── production-report/   #   Báo cáo sản xuất
│   ├── oee-report/          #   Báo cáo OEE
│   ├── energy-report/       #   Báo cáo năng lượng
│   └── custom-report/       #   Báo cáo tùy chỉnh
└── admin/                   # Quản trị
    ├── user-management/     #   Quản lý người dùng
    ├── device-management/   #   Quản lý thiết bị
    ├── system-config/       #   Cấu hình hệ thống
    └── audit-viewer/        #   Xem audit log
```

## 🔧 Công Nghệ

- **Framework**: React 19 + TypeScript
- **UI Library**: Material UI (MUI) v6
- **State Management**: Zustand / Jotai
- **Charts**: Recharts / ECharts / D3.js
- **Real-time**: WebSocket / Server-Sent Events
- **Build**: Vite
- **Testing**: Vitest + Playwright

## 🚀 Chạy

```bash
cd web-portal
pnpm install
pnpm dev          # http://localhost:5173
pnpm build        # Production build
pnpm test         # Unit tests
```
