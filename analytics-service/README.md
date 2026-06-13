# Analytics Service

> Phân tích dữ liệu — Real-time Analytics, Batch Processing & Machine Learning

## 🎯 Chức Năng

- **Real-time Analytics**: Stream processing với Kafka Streams — OEE, throughput, yield
- **Batch Analytics**: Spark/Presto jobs cho báo cáo cuối ngày/tuần/tháng
- **ML Pipeline**: Predictive maintenance, anomaly detection, tối ưu hóa sản xuất
- **Data Export**: Export báo cáo ra CSV, Excel, PDF
- **Dashboard Metrics**: Pre-aggregated metrics cho dashboard real-time

## 📁 Cấu Trúc

```
analytics-service/
├── real-time-analytics/     # Stream processing
│   ├── oee-calculator/      #   Overall Equipment Effectiveness
│   ├── throughput-monitor/  #   Production throughput
│   ├── energy-monitor/      #   Tiêu thụ năng lượng
│   └── kpi-aggregator/      #   KPI real-time aggregation
├── batch-analytics/         # Batch processing
│   ├── daily-report/        #   Báo cáo sản xuất hàng ngày
│   ├── trend-analysis/      #   Phân tích xu hướng dài hạn
│   ├── comparative-report/  #   So sánh giữa các nhà máy/ca
│   └── export-engine/       #   Export CSV, Excel, PDF
└── ml-pipeline/             # Machine Learning
    ├── model-training/      #   Training pipeline
    ├── model-serving/       #   Model inference API
    ├── anomaly-detection/   #   Phát hiện bất thường
    └── predictive-maintenance/ # Dự đoán thời điểm bảo trì
```

## 🔧 Công Nghệ

- **Stream Processing**: Kafka Streams / Apache Flink
- **Batch**: Apache Spark / Trino (Presto)
- **ML**: Python (scikit-learn, TensorFlow) + MLflow
- **Database**: TimescaleDB (đọc telemetry), PostgreSQL (metadata)
- **Object Storage**: MinIO / S3 (model artifacts, reports)

## 📊 OEE Calculation

```
OEE = Availability × Performance × Quality

Availability = Run Time / Planned Production Time
Performance  = (Ideal Cycle Time × Total Count) / Run Time
Quality      = Good Count / Total Count
```
