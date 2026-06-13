# Mobile App

> Ứng dụng di động giám sát & điều khiển nhanh cho Industrial IoT Platform

## 🎯 Chức Năng

- **Monitoring**: Xem dashboard real-time, trạng thái máy móc trên mobile
- **Alarm Push**: Nhận push notification khi có alarm
- **Quick Control**: Điều khiển nhanh (bật/tắt khẩn cấp) từ mobile
- **Offline Mode**: Xem dữ liệu đã cache khi không có kết nối
- **QR Scanner**: Quét QR code trên máy để xem thông tin nhanh

## 📁 Cấu Trúc

```
mobile-app/
├── monitoring/              # Giám sát trên mobile
│   ├── dashboard/           #   Mobile-friendly dashboard
│   ├── machine-view/        #   Xem chi tiết máy
│   ├── realtime-charts/     #   Biểu đồ real-time (tối ưu mobile)
│   └── multi-factory/       #   Chuyển đổi giữa các nhà máy
├── alarm-push/              # Push notification alarm
│   ├── notification-handler/#   Xử lý FCM/APNs push
│   ├── alarm-inbox/         #   Hộp thư alarm
│   └── quick-actions/       #   Thao tác nhanh từ notification
└── quick-control/           # Điều khiển nhanh
    ├── emergency-stop/      #   Dừng khẩn cấp
    ├── quick-commands/      #   Lệnh điều khiển nhanh
    └── barcode-scanner/     #   Quét QR/Barcode máy
```

## 🔧 Công Nghệ

- **Framework**: Flutter 3.x (Dart)
- **State Management**: Riverpod / Bloc
- **Push Notification**: Firebase Cloud Messaging + APNs
- **Local Cache**: SQLite / Hive
- **Charts**: fl_chart / syncfusion_flutter_charts
- **Platform**: iOS 15+ & Android 8+

## 🚀 Chạy

```bash
cd mobile-app
flutter pub get
flutter run          # Development
flutter build apk    # Android build
flutter build ios    # iOS build
```
