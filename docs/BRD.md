# Business Requirements Document (BRD)

## Industrial IoT Platform — Hệ Thống Giám Sát & Điều Khiển Nhà Máy Công Nghiệp Đa Vùng

| Thuộc tính | Giá trị |
|---|---|
| **Phiên bản** | 1.0 |
| **Ngày phát hành** | 13/06/2026 |
| **Tác giả** | Đội ngũ Industrial IoT Platform |
| **Trạng thái** | Approved |
| **Phân loại** | Internal — Restricted |

---

## Mục Lục

1. [Tổng Quan Dự Án](#1-tổng-quan-dự-án)
2. [Mục Tiêu Kinh Doanh](#2-mục-tiêu-kinh-doanh)
3. [Phạm Vi Dự Án](#3-phạm-vi-dự-án)
4. [Phân Tích Các Bên Liên Quan](#4-phân-tích-các-bên-liên-quan)
5. [Phân Tích Thị Trường & Bài Toán](#5-phân-tích-thị-trường--bài-toán)
6. [Yêu Cầu Nghiệp Vụ Cấp Cao](#6-yêu-cầu-nghiệp-vụ-cấp-cao)
7. [Lộ Trình Triển Khai](#7-lộ-trình-triển-khai)
8. [Phân Tích Rủi Ro](#8-phân-tích-rủi-ro)
9. [Phân Tích Chi Phí — Lợi Ích](#9-phân-tích-chi-phí--lợi-ích)
10. [Ràng Buộc & Giả Định](#10-ràng-buộc--giả-định)
11. [Tiêu Chí Thành Công](#11-tiêu-chí-thành-công)
12. [Phê Duyệt](#12-phê-duyệt)

---

## 1. Tổng Quan Dự Án

### 1.1 Bối Cảnh

Ngành sản xuất công nghiệp Việt Nam đang trong quá trình chuyển đổi số mạnh mẽ, với nhu cầu cấp thiết về giám sát và điều khiển nhà máy từ xa theo thời gian thực. Các doanh nghiệp sản xuất vận hành nhiều nhà máy tại các khu vực địa lý khác nhau (Hà Nội, Hồ Chí Minh, Đà Nẵng) đang gặp các thách thức:

- Không có cái nhìn tổng thể về hoạt động sản xuất trên toàn bộ nhà máy
- Dữ liệu máy móc nằm rải rác, không được thu thập và lưu trữ tập trung
- Phát hiện sự cố chậm, dẫn đến thời gian dừng máy kéo dài
- Không có khả năng điều khiển máy móc từ xa
- Thiếu dữ liệu lịch sử để phân tích xu hướng và dự đoán bảo trì

### 1.2 Mô Tả Dự Án

**Industrial IoT Platform** là nền tảng giám sát và điều khiển công nghiệp toàn diện, cho phép:

- **Thu thập dữ liệu** thời gian thực từ máy móc/PLC qua các giao thức công nghiệp (Modbus RTU/TCP, OPC-UA, MQTT, Siemens S7)
- **Giám sát real-time** trạng thái máy móc, dây chuyền sản xuất qua Web Portal và Mobile App
- **Điều khiển từ xa** máy móc với cơ chế xác thực và audit log đầy đủ
- **Cảnh báo thông minh** đa kênh (Email, SMS, Push Notification) khi có sự cố hoặc vượt ngưỡng
- **Phân tích dữ liệu** để dự đoán bảo trì (predictive maintenance) và tối ưu hiệu suất sản xuất (OEE)
- **Triển khai đa vùng (Multi-Region)** hỗ trợ nhà máy tại Hà Nội, Hồ Chí Minh, Đà Nẵng và mở rộng

### 1.3 Tầm Nhìn

Trở thành nền tảng IoT công nghiệp hàng đầu tại Việt Nam, giúp các nhà máy sản xuất đạt được:

- **99.9% uptime** hệ thống giám sát
- **Giảm 30% thời gian dừng máy** không kế hoạch
- **Tăng 15% hiệu suất tổng thể thiết bị (OEE)**
- **Phản hồi sự cố trong vòng 30 giây**

---

## 2. Mục Tiêu Kinh Doanh

### 2.1 Mục Tiêu Chiến Lược

| ID | Mục Tiêu | KPI | Thời hạn |
|----|----------|-----|----------|
| BG-01 | Nâng cao hiệu suất sản xuất | OEE tăng ≥ 15% | 12 tháng sau go-live |
| BG-02 | Giảm thời gian dừng máy ngoài kế hoạch | Downtime giảm ≥ 30% | 6 tháng sau go-live |
| BG-03 | Phản ứng nhanh với sự cố | MTTR (Mean Time to Respond) < 30 giây | Ngay khi go-live |
| BG-04 | Ra quyết định dựa trên dữ liệu | 100% máy móc được giám sát | Giai đoạn 2 |
| BG-05 | Tối ưu chi phí bảo trì | Chi phí bảo trì giảm ≥ 20% | 18 tháng sau go-live |
| BG-06 | Mở rộng quy mô linh hoạt | Hỗ trợ ≥ 10 nhà máy, ≥ 10,000 thiết bị | Giai đoạn 4 |

### 2.2 Đề Xuất Giá Trị

- **Giám sát tập trung**: Một nền tảng duy nhất cho tất cả nhà máy, mọi lúc, mọi nơi
- **Phản ứng tức thì**: Cảnh báo real-time giúp đội ngũ vận hành phản ứng trong tích tắc
- **Dự đoán thông minh**: AI/ML dự đoán hỏng hóc trước khi xảy ra, chuyển từ bảo trì phản ứng sang bảo trì chủ động
- **Điều khiển an toàn**: Điều khiển máy móc từ xa với RBAC, audit log và cơ chế xác nhận an toàn
- **Tiết kiệm chi phí**: Giảm downtime, tối ưu bảo trì, tiết kiệm năng lượng

---

## 3. Phạm Vi Dự Án

### 3.1 Phạm Vi Bao Gồm (In-Scope)

| Hạng mục | Mô tả |
|----------|-------|
| **Edge Gateway** | Thiết bị/cổng kết nối tại mỗi nhà máy, giao tiếp với PLC/máy móc qua Modbus, OPC-UA, S7 |
| **Gateway Service** | IoT Gateway cloud — xác thực thiết bị, định tuyến MQTT, cầu nối MQTT→Kafka |
| **Device Service** | Quản lý vòng đời thiết bị, Digital Twin, OTA firmware update |
| **Telemetry Service** | Thu thập, lưu trữ (TimescaleDB), truy vấn dữ liệu đo lường |
| **Command Service** | Điều khiển máy móc từ cloud → edge, theo dõi trạng thái & audit log |
| **Alarm Service** | Rule engine cảnh báo, quản lý vòng đời alarm, chính sách escalation |
| **Notification Service** | Gửi thông báo đa kênh: Email, SMS, Push Notification, Webhook |
| **Analytics Service** | Phân tích real-time, batch, ML pipeline cho predictive maintenance |
| **Auth Service** | Keycloak IdP, RBAC, audit log người dùng |
| **API Gateway** | Kong API Gateway + Web BFF + Mobile BFF |
| **Web Portal** | Giao diện web React — dashboard, điều khiển, alarm center, báo cáo, quản trị |
| **Mobile App** | Ứng dụng Flutter — giám sát, push alarm, điều khiển nhanh |
| **Deployment** | Kubernetes, Terraform, Ansible, Docker Compose |
| **Observability** | Prometheus + Grafana + Loki stack |

### 3.2 Phạm Vi Không Bao Gồm (Out-of-Scope)

- Tích hợp với hệ thống ERP/MES của khách hàng (sẽ làm ở giai đoạn sau)
- Streaming video từ camera giám sát
- Digital Twin 3D rendering
- Blockchain cho supply chain traceability
- Hỗ trợ giao thức công nghiệp đặc thù khác ngoài Modbus, OPC-UA, S7, MQTT
- On-premise cloud deployment (chỉ hỗ trợ public cloud + edge)

---

## 4. Phân Tích Các Bên Liên Quan

### 4.1 Danh Sách Stakeholder

| Stakeholder | Vai trò | Mối quan tâm chính | Mức độ ảnh hưởng |
|-------------|--------|-------------------|-------------------|
| **Ban Giám Đốc** | Sponsor, ra quyết định chiến lược | ROI, hiệu suất tổng thể, báo cáo điều hành | Rất cao |
| **Giám Đốc Sản Xuất** | Quản lý vận hành nhà máy | OEE, throughput, chất lượng sản phẩm | Cao |
| **Quản Đốc/Ca Trưởng** | Giám sát ca sản xuất | Trạng thái máy, alarm, điều động nhân sự | Cao |
| **Kỹ Thuật Viên Vận Hành** | Vận hành máy móc hàng ngày | Điều khiển máy, nhận alarm, xử lý sự cố | Trung bình |
| **Kỹ Sư Bảo Trì** | Bảo trì, sửa chữa máy móc | Predictive maintenance, lịch sử hỏng hóc, phụ tùng | Trung bình |
| **IT/OT Engineer** | Vận hành hạ tầng kỹ thuật | Hạ tầng edge/cloud, network, security | Cao |
| **Security Officer** | An ninh thông tin | RBAC, audit, X.509, network security | Cao |
| **Data Analyst** | Phân tích dữ liệu sản xuất | Data quality, API truy vấn, xuất báo cáo | Thấp |
| **Đối tác tích hợp** | Tích hợp với hệ thống khác | API documentation, webhook, SDK | Thấp |

### 4.2 Ma Trận RACI (Tóm tắt)

| Hoạt động | Ban GĐ | GĐ SX | Quản Đốc | KT Vận Hành | IT Engineer |
|-----------|--------|-------|----------|-------------|-------------|
| Phê duyệt ngân sách | A,R | C | I | I | I |
| Xác định yêu cầu nghiệp vụ | I | A,R | C | C | C |
| Định nghĩa kiến trúc kỹ thuật | I | I | I | I | A,R |
| Triển khai edge gateway | I | I | I | C | A,R |
| Đào tạo người dùng cuối | I | A | R | C | C |
| Vận hành hệ thống | I | A | R | C | R |

---

## 5. Phân Tích Thị Trường & Bài Toán

### 5.1 Xu Hướng Thị Trường

- **Industry 4.0**: Chính phủ Việt Nam đẩy mạnh chuyển đổi số trong sản xuất (Quyết định 749/QĐ-TTg)
- **IIoT Market**: Thị trường IIoT toàn cầu dự kiến đạt $1.1 trillion vào 2028 (CAGR 23%)
- **Smart Manufacturing**: 70% doanh nghiệp sản xuất lớn tại Việt Nam đang tìm kiếm giải pháp IIoT
- **Predictive Maintenance**: Giảm 25-30% chi phí bảo trì, giảm 70-75% sự cố đột xuất (McKinsey)

### 5.2 Bài Toán Kinh Doanh

| Vấn đề hiện tại | Tác động | Cách giải quyết của nền tảng |
|-----------------|----------|---------------------------|
| Giám sát thủ công, ghi chép giấy | Dữ liệu sai lệch, chậm trễ | Thu thập tự động, real-time từ PLC |
| Phát hiện sự cố chậm (>15 phút) | Downtime kéo dài, tổn thất sản lượng | Cảnh báo tức thì < 30 giây |
| Không có dữ liệu lịch sử | Không phân tích được xu hướng | TimescaleDB lưu trữ dài hạn, data retention policy |
| Bảo trì theo lịch cố định | Lãng phí (bảo trì thừa) hoặc hỏng hóc (bảo trì thiếu) | Predictive maintenance dựa trên ML |
| Nhà máy ở xa không giám sát được | Quản lý phân tán, không đồng bộ | Multi-region deployment, quản lý tập trung |
| Điều khiển thủ công tại chỗ | Chậm, nguy hiểm trong tình huống khẩn cấp | Điều khiển từ xa an toàn, emergency stop |

### 5.3 Phân Tích Cạnh Tranh

| Giải pháp | Điểm mạnh | Điểm yếu | Khác biệt của chúng tôi |
|-----------|-----------|----------|------------------------|
| Siemens MindSphere | Thương hiệu mạnh, tích hợp Siemens | Đắt, lock-in Siemens | Chi phí hợp lý, multi-protocol, linh hoạt |
| PTC ThingWorx | Nhiều connector, low-code | Phức tạp, chi phí cao | Tập trung vào thị trường Việt Nam, hỗ trợ nội địa |
| Open-source (ThingsBoard) | Miễn phí, cộng đồng | Thiếu enterprise features | Đầy đủ tính năng enterprise: RBAC, audit, multi-region |
| Giải pháp tự phát triển | Tùy chỉnh cao | Tốn thời gian, thiếu chuẩn hóa | Nền tảng có sẵn, triển khai nhanh |

---

## 6. Yêu Cầu Nghiệp Vụ Cấp Cao

### 6.1 Năng Lực Cốt Lõi

| ID | Năng lực | Mô tả | Ưu tiên |
|----|----------|-------|---------|
| BC-01 | Kết nối đa giao thức | Kết nối PLC/máy móc qua Modbus, OPC-UA, S7, MQTT | P0 — MVP |
| BC-02 | Thu thập dữ liệu real-time | Thu thập telemetry data từ hàng nghìn thiết bị với tần suất đến 1 giây | P0 — MVP |
| BC-03 | Dashboard giám sát | Hiển thị trạng thái nhà máy, dây chuyền, máy móc real-time | P0 — MVP |
| BC-04 | Cảnh báo thông minh | Rule engine: threshold, rate-of-change, anomaly, state-change | P1 — Core |
| BC-05 | Điều khiển từ xa | Gửi lệnh điều khiển từ Web/Mobile xuống máy móc | P1 — Core |
| BC-06 | Digital Twin | Device shadow: desired/reported state, delta detection | P1 — Core |
| BC-07 | Phân tích OEE | Tính toán OEE real-time và báo cáo | P2 — Advanced |
| BC-08 | Predictive Maintenance | ML models dự đoán thời điểm hỏng hóc | P2 — Advanced |
| BC-09 | Multi-Region | Triển khai edge gateway tại nhiều nhà máy, quản lý tập trung | P3 — Scale |
| BC-10 | Mobile App | Giám sát + push alarm + điều khiển nhanh trên mobile | P3 — Scale |

### 6.2 Yêu Cầu Phi Chức Năng Cấp Cao

| ID | Yêu cầu | Tiêu chí |
|----|---------|----------|
| NFR-01 | Hiệu năng | Hỗ trợ 10,000+ thiết bị kết nối đồng thời, 100K msg/s |
| NFR-02 | Độ trễ | Telemetry latency < 2 giây, command latency < 500ms |
| NFR-03 | Sẵn sàng | Hệ thống cloud 99.9% uptime (SLAs) |
| NFR-04 | Bảo mật | MQTT TLS, X.509 device auth, OAuth2/OIDC, RBAC |
| NFR-05 | Khả năng mở rộng | Horizontal scaling cho tất cả services |
| NFR-06 | Khả năng phục hồi | Store-and-forward tại edge, auto-reconnect |
| NFR-07 | Dữ liệu | Data retention 2 năm, compression sau 7 ngày |

---

## 7. Lộ Trình Triển Khai

### 7.1 Phân Pha

```
Phase 1: MVP (3 tháng)
├── Kiến trúc nền tảng core
├── Edge Gateway: Modbus + MQTT
├── Gateway Service: EMQX → Kafka bridge
├── Telemetry Service: TimescaleDB storage + Query API
├── Web Portal: Dashboard real-time cơ bản
└── Infrastructure: Docker Compose dev, K8s manifests

Phase 2: Core (3 tháng)
├── Device Service: Registry + Digital Twin
├── Command Service: Điều khiển + Audit log
├── Alarm Service: Rule engine + Notification
├── Auth Service: Keycloak + RBAC
├── Web Portal: Device control + Alarm center + Admin
└── API Gateway: Kong + Web BFF

Phase 3: Advanced (3 tháng)
├── Analytics Service: OEE + Reports
├── ML Pipeline: Predictive maintenance
├── Notification Service: Push + Webhook
├── Web Portal: Reports + Advanced dashboard
└── Observability: Prometheus + Grafana + Loki

Phase 4: Multi-Region (3 tháng)
├── Multi-region K8s deployment
├── Edge HA (Active/Standby)
├── Database HA & DR
├── Mobile App: Flutter iOS + Android
└── Performance testing & optimization

Phase 5: Launch & Scale (3 tháng)
├── Production go-live từng nhà máy
├── Load testing 10,000+ thiết bị
├── Tài liệu đào tạo & vận hành
├── CI/CD pipeline hoàn chỉnh
└── SLA monitoring & alerting
```

### 7.2 Các Mốc Quan Trọng

| Mốc | Thời gian dự kiến | Deliverables |
|-----|------------------|--------------|
| M1 — MVP Ready | Tháng 3 | Edge GW + Telemetry + Dashboard cơ bản |
| M2 — Core Ready | Tháng 6 | Device + Command + Alarm + Auth + Web Portal |
| M3 — Advanced Ready | Tháng 9 | Analytics + ML + Reports |
| M4 — Multi-Region | Tháng 12 | HA deployment + Mobile App |
| M5 — Production Go-Live | Tháng 15 | Triển khai thực tế, đào tạo, bàn giao |

---

## 8. Phân Tích Rủi Ro

### 8.1 Ma Trận Rủi Ro

| ID | Rủi ro | Khả năng | Tác động | Mức độ | Biện pháp giảm thiểu |
|----|--------|----------|----------|--------|---------------------|
| R-01 | Độ trễ mạng giữa edge và cloud | Trung bình | Cao | **Cao** | Store-and-forward, local cache, edge rules engine |
| R-02 | Thiết bị PLC đời cũ không hỗ trợ protocol hiện đại | Cao | Trung bình | **Cao** | Hỗ trợ Modbus RTU (serial), phát triển adapter tùy chỉnh |
| R-03 | Bảo mật — tấn công vào edge gateway | Trung bình | Rất cao | **Rất cao** | X.509 certs, TLS, VPN, firewall, regular security audit |
| R-04 | Không đủ nhân sự kỹ thuật triển khai | Trung bình | Cao | **Cao** | Tài liệu chi tiết, đào tạo, hỗ trợ triển khai từ xa |
| R-05 | Dữ liệu telemetry quá lớn gây quá tải | Trung bình | Trung bình | **Trung bình** | Compression policy, down-sampling, retention policy |
| R-06 | Người dùng không chấp nhận hệ thống mới | Trung bình | Cao | **Cao** | UX đơn giản, đào tạo, triển khai từng bước, feedback loop |
| R-07 | Vendor lock-in với cloud provider | Thấp | Trung bình | **Thấp** | Kiến trúc cloud-agnostic, Kubernetes, Terraform |
| R-08 | Sự cố mất kết nối internet tại nhà máy | Cao | Trung bình | **Cao** | Local buffer (SQLite), store-and-forward, edge rules engine |

---

## 9. Phân Tích Chi Phí — Lợi Ích

### 9.1 Chi Phí Dự Kiến

| Hạng mục | Chi phí ước tính (USD) | Ghi chú |
|----------|------------------------|---------|
| Phát triển phần mềm | $XXX,XXX | Đội ngũ 5-8 kỹ sư, 12-15 tháng |
| Hạ tầng cloud (năm đầu) | $XX,XXX/tháng | K8s cluster, DB, network, storage |
| Edge gateway hardware | $X,XXX/nhà máy | Raspberry Pi 4 / Mini PC |
| Bản quyền phần mềm | $XX,XXX | Keycloak (free), EMQX (community), Kong (community) |
| Đào tạo & triển khai | $XX,XXX | 2 tuần/nhà máy |
| Bảo trì & hỗ trợ (hàng năm) | $XX,XXX | 20% chi phí phát triển |

### 9.2 Lợi Ích Dự Kiến

| Lợi ích | Ước tính/năm | Cách tính |
|----------|-------------|-----------|
| Giảm downtime | $XXX,XXX | 30% × chi phí downtime hiện tại |
| Tối ưu bảo trì | $XXX,XXX | 20% × ngân sách bảo trì hiện tại |
| Tăng OEE | $XXX,XXX | 15% tăng sản lượng |
| Tiết kiệm nhân công giám sát | $XX,XXX | Giảm 2-3 người/ca/nhà máy |
| Tiết kiệm năng lượng | $XX,XXX | 10% nhờ tối ưu vận hành |

### 9.3 ROI Dự Kiến

- **Payback period**: 18-24 tháng
- **3-Year ROI**: 250-350%
- **NPV** (10% discount rate): Dương

> *Số liệu chi tiết sẽ được tính toán dựa trên quy mô cụ thể của từng nhà máy.*

---

## 10. Ràng Buộc & Giả Định

### 10.1 Ràng Buộc Kỹ Thuật

- Edge gateway phải chạy trên Linux (Raspberry Pi 4 hoặc Mini PC x86)
- MQTT Broker phải hỗ trợ ít nhất 1M concurrent connections
- Cloud deployment trên Kubernetes (EKS/GKE/AKS)
- Hỗ trợ TLS 1.2+ cho mọi kết nối edge-cloud
- Edge → Cloud latency trung bình < 50ms (cùng quốc gia)

### 10.2 Ràng Buộc Nghiệp Vụ

- Tuân thủ quy định an toàn lao động Việt Nam
- Tuân thủ quy định bảo mật dữ liệu (Nghị định 53/2022/NĐ-CP)
- Hệ thống phải hoạt động 24/7, kể cả ngày lễ
- Ngôn ngữ giao diện: Tiếng Việt (ưu tiên) + Tiếng Anh

### 10.3 Giả Định

- Nhà máy có kết nối internet ổn định (≥ 10 Mbps uplink)
- PLC/máy móc hỗ trợ ít nhất một trong các giao thức: Modbus RTU/TCP, OPC-UA
- Có sẵn nhân sự IT tại mỗi nhà máy để cài đặt edge gateway
- Người dùng cuối có smartphone (cho Mobile App)
- Ban lãnh đạo cam kết thay đổi quy trình vận hành

---

## 11. Tiêu Chí Thành Công

### 11.1 Tiêu Chí Định Lượng

| ID | Tiêu chí | Mục tiêu | Cách đo |
|----|----------|----------|---------|
| SC-01 | Tỷ lệ thiết bị được kết nối | 100% thiết bị target | Device Registry vs thực tế |
| SC-02 | Data ingestion rate | 100K messages/giây | Kafka consumer lag |
| SC-03 | Telemetry latency (end-to-end) | < 2 giây | MQTT publish → UI display |
| SC-04 | Command latency | < 500ms | UI click → MQTT publish |
| SC-05 | Alarm response time (MTTR) | < 30 giây | Alarm raised → User notified |
| SC-06 | System uptime (cloud) | 99.9% | Prometheus uptime metrics |
| SC-07 | User adoption rate | 90% sau 3 tháng | Active users / total users |
| SC-08 | Data loss rate | < 0.01% | MQTT published vs Kafka consumed vs DB stored |

### 11.2 Tiêu Chí Định Tính

- Người dùng đánh giá giao diện dễ sử dụng (SUS score ≥ 75)
- Quản đốc có thể ra quyết định dựa trên dashboard mà không cần báo cáo thủ công
- Kỹ thuật viên phản ứng với alarm trong vòng 30 giây
- IT team có thể tự triển khai edge gateway cho nhà máy mới trong < 1 ngày

---

## 12. Phê Duyệt

| Vai trò | Họ tên | Chữ ký | Ngày |
|---------|--------|--------|------|
| **Project Sponsor** | _________________ | _______________ | __/__/____ |
| **Business Owner** | _________________ | _______________ | __/__/____ |
| **Technical Lead** | _________________ | _______________ | __/__/____ |
| **Security Officer** | _________________ | _______________ | __/__/____ |

---

## Phụ Lục

### A. Thuật Ngữ & Viết Tắt

| Thuật ngữ | Định nghĩa |
|-----------|-----------|
| **IIoT** | Industrial Internet of Things — IoT trong công nghiệp |
| **PLC** | Programmable Logic Controller — Bộ điều khiển logic khả trình |
| **OEE** | Overall Equipment Effectiveness — Hiệu suất tổng thể thiết bị |
| **MTTR** | Mean Time to Respond — Thời gian phản hồi trung bình |
| **MQTT** | Message Queuing Telemetry Transport — Giao thức truyền tin nhẹ cho IoT |
| **OPC-UA** | Open Platform Communications Unified Architecture — Giao thức công nghiệp |
| **Digital Twin** | Bản sao số của thiết bị vật lý |
| **BFF** | Backend-for-Frontend — Lớp API chuyên biệt cho từng client |
| **RBAC** | Role-Based Access Control — Phân quyền theo vai trò |
| **OTA** | Over-the-Air — Cập nhật firmware từ xa |

### B. Tài Liệu Tham Khảo

- [User Requirements Document (URD)](URD.md)
- [Software Requirements Specification (SRS)](SRS.md)
- [Kiến trúc hệ thống (README)](../README.md)
- README của từng service trong repository

---

> **Tài liệu này là tài sản trí tuệ của Industrial IoT Platform. Không được sao chép hoặc phân phối khi chưa có sự đồng ý.**
