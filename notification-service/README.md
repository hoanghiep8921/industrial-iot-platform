# Notification Service

> Gửi thông báo đa kênh — Email, SMS, Push Notification, Webhook

## 🎯 Chức Năng

- **Multi-channel**: Email, SMS, Push Notification (FCM/APNs), Webhook
- **Template Engine**: Tùy chỉnh nội dung thông báo theo template
- **User Preferences**: Người dùng chọn kênh nhận thông báo, giờ yên lặng
- **Rate Limiting**: Tránh spam notification
- **Delivery Tracking**: Theo dõi trạng thái gửi (sent/delivered/failed)

## 📁 Cấu Trúc

```
notification-service/
├── email-notifier/          # Email notification
│   ├── smtp-client/         #   SMTP client
│   ├── template-engine/     #   HTML email templates
│   └── attachment-handler/  #   File đính kèm (báo cáo PDF)
├── sms-notifier/            # SMS notification
│   ├── twilio-adapter/      #   Twilio integration
│   └── vnpt-adapter/        #   VNPT/Viettel SMS gateway
├── push-notifier/           # Mobile push notification
│   ├── fcm-client/          #   Firebase Cloud Messaging
│   ├── apns-client/         #   Apple Push Notification Service
│   └── badge-manager/       #   App icon badge count
└── webhook-notifier/        # Webhook tích hợp
    ├── http-dispatcher/      #   Gửi HTTP POST đến URL đích
    └── retry-handler/        #   Retry với exponential backoff
```

## 🔧 Công Nghệ

- **Ngôn ngữ**: Go
- **Email**: SMTP / AWS SES / SendGrid
- **SMS**: Twilio / AWS SNS
- **Push**: Firebase Cloud Messaging + Apple APNs
- **Database**: PostgreSQL (notification history)
- **Message Queue**: Kafka consumer (nhận notification request)

## 📊 Notification Template

```json
{
  "alarm_raised": {
    "channels": ["push", "email"],
    "email": {
      "subject": "🚨 [{severity}] {alarm_name} - {device_name}",
      "template": "alarm_raised.html"
    },
    "push": {
      "title": "🚨 {alarm_name}",
      "body": "{metric_name} = {value} (threshold: {threshold})"
    },
    "sms": {
      "template": "ALARM: {alarm_name} at {device_name}. {metric}={value}"
    }
  }
}
```
