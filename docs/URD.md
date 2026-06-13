# User Requirements Document (URD)

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

1. [Giới Thiệu](#1-giới-thiệu)
2. [Phân Tích Người Dùng](#2-phân-tích-người-dùng)
3. [User Personas](#3-user-personas)
4. [User Stories](#4-user-stories)
5. [Use Cases](#5-use-cases)
6. [Luồng Công Việc Người Dùng](#6-luồng-công-việc-người-dùng)
7. [Yêu Cầu Giao Diện Người Dùng](#7-yêu-cầu-giao-diện-người-dùng)
8. [Yêu Cầu Trải Nghiệm Người Dùng](#8-yêu-cầu-trải-nghiệm-người-dùng)
9. [Yêu Cầu Báo Cáo](#9-yêu-cầu-báo-cáo)
10. [Phụ Lục](#10-phụ-lục)

---

## 1. Giới Thiệu

### 1.1 Mục Đích

Tài liệu này mô tả các yêu cầu của người dùng đối với **Industrial IoT Platform**. Nó bao gồm phân tích người dùng (personas), user stories, use cases, và luồng công việc nhằm đảm bảo hệ thống đáp ứng đúng nhu cầu thực tế của từng nhóm người dùng trong nhà máy công nghiệp.

### 1.2 Phạm Vi

Tài liệu bao phủ toàn bộ trải nghiệm người dùng trên cả ba nền tảng:

- **Web Portal**: Giao diện chính cho máy tính để bàn (React + TypeScript)
- **Mobile App**: Ứng dụng di động cho iOS và Android (Flutter)
- **Edge Gateway**: Giao diện cấu hình và giám sát thiết bị edge (dành cho IT)

### 1.3 Đối Tượng Đọc

- Product Owner & Business Analyst
- UI/UX Designer
- Development Team
- QA/Testing Team
- Khách hàng & Người dùng cuối

---

## 2. Phân Tích Người Dùng

### 2.1 Phân Loại Người Dùng

| Nhóm | Vai trò | Số lượng ước tính | Tần suất sử dụng | Nền tảng chính |
|------|--------|-------------------|------------------|----------------|
| **Ban Giám Đốc** | Ra quyết định chiến lược | 3-5 người | 1-2 lần/tuần | Web Portal, Mobile |
| **Giám Đốc Sản Xuất** | Quản lý toàn bộ sản xuất | 1-2 người/nhà máy | Hàng ngày, nhiều lần | Web Portal |
| **Quản Đốc / Ca Trưởng** | Giám sát ca sản xuất | 3-4 người/nhà máy | Liên tục trong ca | Web Portal, Mobile |
| **Kỹ Thuật Viên Vận Hành** | Vận hành máy móc | 10-30 người/nhà máy | Liên tục trong ca | Web Portal, Mobile |
| **Kỹ Sư Bảo Trì** | Bảo trì & sửa chữa | 3-5 người/nhà máy | Khi có sự cố/bảo trì | Web Portal, Mobile |
| **IT/OT Engineer** | Quản trị hạ tầng | 2-3 người/toàn hệ thống | Hàng ngày | Web Portal, CLI |
| **Data Analyst** | Phân tích dữ liệu | 1-2 người/toàn hệ thống | Hàng tuần | Web Portal |
| **System Admin** | Quản trị hệ thống | 1-2 người/toàn hệ thống | Hàng ngày | Web Portal |

### 2.2 Nhu Cầu Theo Vai Trò

| Vai trò | Nhu cầu chính | Nỗi đau (Pain Points) |
|---------|--------------|----------------------|
| Ban Giám Đốc | Báo cáo tổng quan, KPIs | Không có dữ liệu real-time để ra quyết định |
| Giám Đốc SX | OEE, throughput, chất lượng | Phải đợi báo cáo cuối ca mới biết tình hình |
| Quản Đốc | Trạng thái máy real-time, alarm | Phải đi bộ kiểm tra từng máy |
| KT Vận Hành | Điều khiển máy, nhận alarm | Không biết máy hỏng cho đến khi thấy trực tiếp |
| KS Bảo Trì | Lịch sử hỏng hóc, dự đoán | Không biết máy nào sắp hỏng để ưu tiên |
| IT Engineer | Quản lý thiết bị, network | Thiếu công cụ quản lý tập trung |
| Data Analyst | Dữ liệu sạch, API truy vấn | Dữ liệu rải rác, không chuẩn hóa |
| System Admin | Quản lý user, phân quyền | Không có audit trail |

---

## 3. User Personas

### 3.1 Persona 1: Anh Minh — Quản Đốc Ca Sản Xuất

```
┌─────────────────────────────────────────────────────────┐
│  👤 Nguyễn Văn Minh                                      │
│  📋 Quản Đốc Ca Sản Xuất — Nhà máy Hà Nội               │
│  🎂 38 tuổi, 12 năm kinh nghiệm                          │
│  📱 Thành thạo smartphone, cơ bản máy tính               │
│  ⏰ Làm việc theo ca: 6h-14h hoặc 14h-22h               │
├─────────────────────────────────────────────────────────┤
│  Mục tiêu:                                               │
│  • Đảm bảo ca sản xuất đạt KPI sản lượng                 │
│  • Phát hiện và xử lý sự cố nhanh nhất có thể           │
│  • Điều phối công nhân hiệu quả                          │
│  • Báo cáo tình hình sản xuất cho Giám đốc SX           │
├─────────────────────────────────────────────────────────┤
│  Nỗi đau hiện tại:                                       │
│  • Phải đi bộ 5km mỗi ca để kiểm tra từng máy           │
│  • Không biết máy hỏng cho đến khi công nhân báo         │
│  • Ghi chép sản lượng bằng giấy, dễ sai lệch            │
│  • Không có dữ liệu để truy vết nguyên nhân sự cố       │
├─────────────────────────────────────────────────────────┤
│  Kỳ vọng với hệ thống:                                   │
│  • Xem toàn bộ trạng thái máy từ văn phòng              │
│  • Nhận cảnh báo ngay khi máy có dấu hiệu bất thường    │
│  • Dashboard trực quan, dễ hiểu không cần đào tạo nhiều  │
│  • Xem lại dữ liệu ca trước khi bàn giao                │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Persona 2: Chị Hương — Kỹ Thuật Viên Vận Hành

```
┌─────────────────────────────────────────────────────────┐
│  👤 Trần Thị Hương                                       │
│  📋 Kỹ Thuật Viên Vận Hành — Khu vực Lò Nhiệt           │
│  🎂 28 tuổi, 5 năm kinh nghiệm                          │
│  📱 Smartphone thành thạo                                │
│  ⏰ Trực ca 8 tiếng, làm việc gần máy móc               │
├─────────────────────────────────────────────────────────┤
│  Mục tiêu:                                               │
│  • Vận hành máy móc an toàn, đúng quy trình              │
│  • Phát hiện bất thường qua âm thanh, nhiệt độ, rung    │
│  • Ghi nhận thông số vận hành định kỳ                   │
│  • Thực hiện lệnh điều khiển từ Quản đốc                │
├─────────────────────────────────────────────────────────┤
│  Nỗi đau hiện tại:                                       │
│  • Phải đứng cạnh máy đọc đồng hồ, môi trường ồn, nóng │
│  • Không có cảnh báo sớm khi thông số vượt ngưỡng       │
│  • Phải chạy đi tìm Quản đốc khi phát hiện sự cố        │
│  • Không thể điều khiển máy từ xa khi cần dừng khẩn     │
├─────────────────────────────────────────────────────────┤
│  Kỳ vọng với hệ thống:                                   │
│  • Nhận thông số máy real-time trên điện thoại           │
│  • Cảnh báo âm thanh + rung khi có alarm                │
│  • Điều khiển máy từ xa an toàn                         │
│  • Giao diện đơn giản, nút to, dễ thao tác              │
└─────────────────────────────────────────────────────────┘
```

### 3.3 Persona 3: Ông Tuấn — Giám Đốc Sản Xuất

```
┌─────────────────────────────────────────────────────────┐
│  👤 Lê Anh Tuấn                                          │
│  📋 Giám Đốc Sản Xuất — Phụ trách 3 nhà máy             │
│  🎂 45 tuổi, 20 năm kinh nghiệm                         │
│  💻 Thành thạo máy tính, Excel, BI tools                 │
│  ✈️ Di chuyển thường xuyên giữa HN, HCM, ĐN             │
├─────────────────────────────────────────────────────────┤
│  Mục tiêu:                                               │
│  • Tối ưu hiệu suất sản xuất toàn bộ 3 nhà máy          │
│  • So sánh hiệu suất giữa các nhà máy, ca sản xuất      │
│  • Dự báo sản lượng, lập kế hoạch bảo trì               │
│  • Báo cáo KPI cho Ban Giám Đốc                         │
├─────────────────────────────────────────────────────────┤
│  Nỗi đau hiện tại:                                       │
│  • Dữ liệu sản xuất nằm ở Excel, mỗi nhà máy một file   │
│  • Mất 2-3 ngày để tổng hợp báo cáo tháng               │
│  • Không so sánh được real-time giữa các nhà máy        │
│  • Không có công cụ phân tích xu hướng dài hạn          │
├─────────────────────────────────────────────────────────┤
│  Kỳ vọng với hệ thống:                                   │
│  • Dashboard tổng quan tất cả nhà máy trên một màn hình  │
│  • Báo cáo OEE tự động theo ca/ngày/tuần/tháng          │
│  • So sánh hiệu suất giữa các nhà máy, khu vực          │
│  • Export báo cáo PDF/Excel cho Ban Giám Đốc            │
│  • Xem dữ liệu lịch sử đến 2 năm                        │
└─────────────────────────────────────────────────────────┘
```

### 3.4 Persona 4: Anh Dũng — IT/OT Engineer

```
┌─────────────────────────────────────────────────────────┐
│  👤 Phạm Tiến Dũng                                       │
│  📋 IT/OT Engineer — Phụ trách hạ tầng edge & cloud     │
│  🎂 32 tuổi, 8 năm kinh nghiệm                          │
│  💻 Linux, Network, Security thành thạo                  │
│  🔧 Làm việc với cả IT (cloud) và OT (factory floor)     │
├─────────────────────────────────────────────────────────┤
│  Mục tiêu:                                               │
│  • Đảm bảo hạ tầng edge-cloud hoạt động ổn định         │
│  • Quản lý vòng đời thiết bị edge (cài đặt, cập nhật)  │
│  • Bảo mật mạng OT, phát hiện xâm nhập                  │
│  • Troubleshooting khi có sự cố kết nối                 │
├─────────────────────────────────────────────────────────┤
│  Nỗi đau hiện tại:                                       │
│  • Phải đến tận nhà máy để cập nhật firmware edge       │
│  • Không có công cụ giám sát tập trung cho edge gateway │
│  • Khó khăn trong việc quản lý chứng chỉ X.509          │
│  • Thiếu audit trail khi có sự cố bảo mật               │
├─────────────────────────────────────────────────────────┤
│  Kỳ vọng với hệ thống:                                   │
│  • Quản lý tất cả edge gateway từ xa                    │
│  • OTA firmware update cho edge gateway                 │
│  • Giám sát health edge: CPU, RAM, disk, network        │
│  • Cảnh báo khi edge mất kết nối                        │
│  • Quản lý tập trung chứng chỉ X.509                    │
│  • Audit log mọi thao tác hệ thống                      │
└─────────────────────────────────────────────────────────┘
```

### 3.5 Persona 5: Chị Lan — System Administrator

```
┌─────────────────────────────────────────────────────────┐
│  👤 Võ Thị Lan                                           │
│  📋 System Administrator                                 │
│  🎂 35 tuổi, 10 năm kinh nghiệm                         │
│  💻 Thành thạo Linux, K8s, Database, Security           │
├─────────────────────────────────────────────────────────┤
│  Mục tiêu:                                               │
│  • Quản lý user, role, permission toàn hệ thống          │
│  • Cấu hình hệ thống theo yêu cầu nghiệp vụ             │
│  • Đảm bảo an ninh, tuân thủ quy định                   │
│  • Sao lưu & phục hồi dữ liệu                           │
├─────────────────────────────────────────────────────────┤
│  Kỳ vọng với hệ thống:                                   │
│  • Giao diện quản trị đầy đủ, dễ sử dụng                │
│  • Quản lý user, role, permission tập trung             │
│  • Xem audit log mọi hoạt động                          │
│  • Cấu hình alarm rule, notification template           │
│  • Quản lý tenant/nhà máy                               │
│  • Backup/restore database                              │
└─────────────────────────────────────────────────────────┘
```

---

## 4. User Stories

### 4.1 Epic 1: Giám Sát Real-Time (P0 — MVP)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-001 | Là một **Quản đốc**, tôi muốn xem dashboard tổng quan nhà máy để nắm bắt tình hình sản xuất ngay lập tức | Quản Đốc | Dashboard hiển thị: tổng số máy online/offline, trạng thái từng dây chuyền, alarm đang active, OEE hiện tại. Dữ liệu refresh ≤ 2 giây |
| US-002 | Là một **KT Vận hành**, tôi muốn xem thông số real-time của từng máy để phát hiện bất thường | KT Vận Hành | Click vào máy → hiển thị tất cả metrics: nhiệt độ, áp suất, tốc độ, rung... kèm biểu đồ real-time. Dữ liệu refresh ≤ 1 giây |
| US-003 | Là một **Giám đốc SX**, tôi muốn xem dashboard tổng quan tất cả nhà máy để so sánh hiệu suất | GĐ Sản Xuất | Dashboard đa nhà máy: chọn từng nhà máy hoặc xem tất cả. Hiển thị OEE, throughput, alarm count |
| US-004 | Là một **KT Vận hành**, tôi muốn xem biểu đồ xu hướng của một metric trong X giờ qua để phân tích | KT Vận Hành | Chọn metric, chọn khoảng thời gian (1h/6h/24h/7d/30d) → biểu đồ line chart. Hỗ trợ zoom |
| US-005 | Là một **KT Vận hành trên mobile**, tôi muốn xem trạng thái máy trên điện thoại khi đang đi tuần tra | KT Vận Hành (Mobile) | Mobile dashboard: danh sách máy theo khu vực, màu sắc trạng thái (xanh/vàng/đỏ), chạm để xem chi tiết |

### 4.2 Epic 2: Cảnh Báo & Thông Báo (P1 — Core)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-010 | Là một **Quản đốc**, tôi muốn nhận cảnh báo ngay khi máy vượt ngưỡng nguy hiểm | Quản Đốc | Khi metric vượt threshold → hiển thị alarm trên Web Portal + gửi notification. Độ trễ < 5 giây |
| US-011 | Là một **KT Vận hành**, tôi muốn nhận push notification trên điện thoại khi có alarm trong khu vực tôi phụ trách | KT Vận Hành | Mobile push notification kèm âm thanh, rung. Hiển thị: tên máy, metric, giá trị, ngưỡng |
| US-012 | Là một **Quản đốc**, tôi muốn xác nhận (acknowledge) alarm để thông báo tôi đang xử lý | Quản Đốc | Nút "Acknowledge" trên alarm. Ghi lại: ai ack, thời gian. Trạng thái chuyển RAISED → ACKNOWLEDGED |
| US-013 | Là một **Giám đốc SX**, tôi muốn alarm được tự động escalate nếu không được xử lý trong X phút | GĐ Sản Xuất | Cấu hình escalation policy: sau 5 phút → thông báo Quản đốc, sau 15 phút → thông báo GĐ SX, sau 30 phút → gọi điện |
| US-014 | Là một **Quản đốc**, tôi muốn tạm tắt alarm theo lịch bảo trì để không bị spam | Quản Đốc | Cấu hình suppression schedule: chọn máy, chọn khoảng thời gian. Alarm không kích hoạt trong thời gian đó |
| US-015 | Là một **KT Vận hành**, tôi muốn xem lịch sử alarm của máy để hiểu tần suất và nguyên nhân | KT Vận Hành | Alarm history: lọc theo máy, loại alarm, severity, thời gian. Hiển thị timeline |

### 4.3 Epic 3: Điều Khiển Thiết Bị (P1 — Core)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-020 | Là một **KT Vận hành**, tôi muốn gửi lệnh điều khiển đến máy (bật/tắt, đổi tốc độ...) | KT Vận Hành | Bảng điều khiển máy: nút Bật/Tắt, slider tốc độ, input setpoint. Yêu cầu xác nhận trước khi gửi. Hiển thị kết quả |
| US-021 | Là một **Quản đốc**, tôi muốn dừng khẩn cấp tất cả máy trong khu vực khi có sự cố | Quản Đốc | Nút "Emergency Stop" cho khu vực/dây chuyền. Yêu cầu xác nhận 2 bước. Ghi audit log |
| US-022 | Là một **Quản đốc**, tôi muốn xem lịch sử tất cả lệnh điều khiển đã gửi | Quản Đốc | Command history: ai gửi, lệnh gì, máy nào, thời gian, trạng thái (done/failed/timeout) |
| US-023 | Là một **Giám đốc SX**, tôi muốn đảm bảo mọi lệnh điều khiển đều được ghi audit log | GĐ Sản Xuất | Audit log không thể xóa/sửa. Ghi: user, role, command, target device, timestamp, result |
| US-024 | Là một **KT Vận hành (Mobile)**, tôi muốn dừng khẩn cấp máy từ điện thoại khi phát hiện nguy hiểm | KT Vận Hành (Mobile) | Nút Emergency Stop lớn, màu đỏ trên mobile. Vuốt để xác nhận. Gửi lệnh ngay lập tức |

### 4.4 Epic 4: Quản Lý Thiết Bị (P1 — Core)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-030 | Là một **IT Engineer**, tôi muốn đăng ký thiết bị mới vào hệ thống | IT Engineer | Form đăng ký: serial number, tên, model, nhà máy, khu vực, giao thức, capabilities. Tự động generate device ID |
| US-031 | Là một **IT Engineer**, tôi muốn xem danh sách tất cả thiết bị và trạng thái | IT Engineer | Bảng danh sách: lọc theo nhà máy, khu vực, trạng thái. Hiển thị: online/offline/maintenance, last seen |
| US-032 | Là một **KT Vận hành**, tôi muốn xem Digital Twin của máy — trạng thái mong muốn vs thực tế | KT Vận Hành | Device Shadow view: cột "Desired" vs "Reported". Highlight sự khác biệt (delta). Cập nhật real-time |
| US-033 | Là một **IT Engineer**, tôi muốn cập nhật firmware từ xa (OTA) cho thiết bị | IT Engineer | Upload firmware binary, chọn thiết bị đích, lập lịch cập nhật, theo dõi tiến độ. Hỗ trợ rollback |
| US-034 | Là một **Quản đốc**, tôi muốn phân nhóm thiết bị theo khu vực/dây chuyền để quản lý | Quản Đốc | Cây thiết bị: Nhà máy → Khu vực → Dây chuyền → Máy. Kéo thả để di chuyển |

### 4.5 Epic 5: Phân Tích & Báo Cáo (P2 — Advanced)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-040 | Là một **Giám đốc SX**, tôi muốn xem báo cáo OEE theo ca/ngày/tuần/tháng | GĐ Sản Xuất | Báo cáo OEE: Availability, Performance, Quality. Biểu đồ cột theo thời gian. So sánh giữa các ca/nhà máy |
| US-041 | Là một **Data Analyst**, tôi muốn export dữ liệu telemetry ra CSV/Excel để phân tích | Data Analyst | Chọn metrics, thiết bị, khoảng thời gian → Export CSV/Excel. Hỗ trợ lên đến 1M data points |
| US-042 | Là một **Giám đốc SX**, tôi muốn xem báo cáo tiêu thụ năng lượng để tối ưu chi phí | GĐ Sản Xuất | Báo cáo năng lượng: kWh theo máy/khu vực/nhà máy. So sánh với baseline. Biểu đồ trend |
| US-043 | Là một **KS Bảo trì**, tôi muốn hệ thống dự đoán thời điểm cần bảo trì máy | KS Bảo Trì | Predictive maintenance dashboard: danh sách máy có nguy cơ hỏng trong 7/14/30 ngày tới, kèm confidence score |
| US-044 | Là một **Giám đốc SX**, tôi muốn so sánh hiệu suất giữa các nhà máy | GĐ Sản Xuất | Comparative report: side-by-side KPI của các nhà máy. Biểu đồ radar cho OEE thành phần |
| US-045 | Là một **Data Analyst**, tôi muốn tạo báo cáo tùy chỉnh với các metrics tự chọn | Data Analyst | Custom report builder: chọn metrics, group by, filter, time range → preview → export |

### 4.6 Epic 6: Quản Trị Hệ Thống (P1 — Core)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-050 | Là một **System Admin**, tôi muốn quản lý người dùng và phân quyền | System Admin | CRUD user, gán role (Admin/Supervisor/Operator/Viewer), gán scope (global/factory/area) |
| US-051 | Là một **System Admin**, tôi muốn xem audit log mọi hoạt động trong hệ thống | System Admin | Audit log: lọc theo user, action, thời gian, resource. Không thể xóa. Export CSV |
| US-052 | Là một **System Admin**, tôi muốn cấu hình alarm rule mà không cần code | System Admin | Rule builder UI: chọn metric, điều kiện (threshold/trend/state), ngưỡng, severity, group. Test rule với data mẫu |
| US-053 | Là một **IT Engineer**, tôi muốn giám sát health của tất cả edge gateway | IT Engineer | Edge health dashboard: CPU/RAM/Disk/Network của từng edge. Cảnh báo khi vượt ngưỡng |
| US-054 | Là một **System Admin**, tôi muốn cấu hình notification template | System Admin | Template editor: chọn kênh (Email/SMS/Push), soạn nội dung với biến (${device_name}, ${value}...) |

### 4.7 Epic 7: Mobile App (P3 — Scale)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-060 | Là một **KT Vận hành**, tôi muốn xem dashboard real-time trên điện thoại | KT Vận Hành (Mobile) | Mobile dashboard: widget KPI, danh sách máy có trạng thái màu, kéo xuống refresh |
| US-061 | Là một **KT Vận hành**, tôi muốn nhận push notification alarm và thực hiện thao tác nhanh | KT Vận Hành (Mobile) | Push kèm action buttons: "Acknowledge", "Xem chi tiết". Mở app đúng màn hình |
| US-062 | Là một **KT Vận hành**, tôi muốn dừng khẩn cấp máy từ màn hình khóa | KT Vận Hành (Mobile) | Notification có nút "Emergency Stop". Vuốt xác nhận. Không cần mở app |
| US-063 | Là một **KT Vận hành**, tôi muốn quét QR code trên máy để xem thông tin nhanh | KT Vận Hành (Mobile) | Camera scan QR → hiển thị: tên máy, trạng thái, metrics chính, alarm gần đây |
| US-064 | Là một **Quản đốc**, tôi muốn xem dữ liệu đã cache khi không có internet | Quản Đốc (Mobile) | Offline mode: hiển thị dữ liệu cuối cùng được sync. Có indicator "Offline". Tự động sync khi có mạng |

### 4.8 Epic 8: Edge Gateway (P0 — MVP)

| ID | User Story | Actor | Acceptance Criteria |
|----|-----------|-------|---------------------|
| US-070 | Là một **IT Engineer**, tôi muốn cài đặt edge gateway dễ dàng tại nhà máy mới | IT Engineer | Cài đặt < 30 phút: flash image, cấu hình network, đăng ký với cloud. Provisioning tự động |
| US-071 | Là một **IT Engineer**, tôi muốn edge gateway tự động kết nối lại sau khi mất internet | IT Engineer | Auto-reconnect với exponential backoff. Store-and-forward data trong thời gian mất kết nối. Tự động đồng bộ khi có mạng |
| US-072 | Là một **KT Vận hành**, tôi muốn máy vẫn dừng được khi mất kết nối internet | KT Vận Hành | Emergency stop hoạt động local, không phụ thuộc cloud. Rule engine edge xử lý alarm cơ bản offline |

---

## 5. Use Cases

### 5.1 UC-01: Giám Sát Nhà Máy Real-Time

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-01 — Giám Sát Nhà Máy Real-Time                │
├─────────────────────────────────────────────────────────────┤
│ Actor: Quản Đốc, KT Vận Hành, Giám Đốc SX                  │
│ Trigger: Người dùng mở Dashboard                            │
│ Pre-condition: Đã đăng nhập, có quyền xem dashboard         │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. User mở Web Portal / Mobile App                          │
│ 2. Hệ thống hiển thị dashboard tổng quan                    │
│    - Factory Overview: số máy online/offline/maintenance     │
│    - Active Alarms: số alarm đang active                    │
│    - OEE Gauge: OEE hiện tại                                │
│    - Production Throughput: sản lượng ca hiện tại           │
│ 3. User click vào một dây chuyền                            │
│ 4. Hệ thống hiển thị chi tiết dây chuyền:                   │
│    - Danh sách máy với trạng thái (xanh/vàng/đỏ)            │
│    - Metrics chính của từng máy                             │
│ 5. User click vào một máy                                   │
│ 6. Hệ thống hiển thị chi tiết máy:                          │
│    - Tất cả metrics dạng gauge/line chart                   │
│    - Trạng thái hoạt động                                   │
│    - Alarm history gần đây                                  │
│    - Device Shadow (Desired vs Reported)                    │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Không có dữ liệu: Hiển thị "No data" + thời gian last   │
│     seen                                                     │
│ A2. Máy offline: Hiển thị trạng thái offline + thời gian     │
│     mất kết nối                                              │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: User có cái nhìn tổng quan về trạng thái     │
│                 nhà máy và có thể khoanh vùng vấn đề          │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 UC-02: Xử Lý Cảnh Báo (Alarm Handling)

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-02 — Xử Lý Cảnh Báo                             │
├─────────────────────────────────────────────────────────────┤
│ Actor: Quản Đốc, KT Vận Hành                                │
│ Trigger: Một alarm được kích hoạt                            │
│ Pre-condition: Rule engine đánh giá điều kiện alarm = true   │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. Rule Engine phát hiện metric vượt ngưỡng                 │
│ 2. Alarm State Machine tạo alarm với trạng thái RAISED      │
│ 3. Notification Service gửi thông báo:                       │
│    - Web Portal: Alarm xuất hiện trên Alarm Center          │
│    - Mobile: Push notification (nếu user đăng ký)           │
│    - Email: Gửi email (nếu severity = Critical/Emergency)   │
│ 4. Operator nhận thông báo, click vào alarm                 │
│ 5. Hệ thống hiển thị chi tiết alarm:                        │
│    - Máy nào, metric gì, giá trị hiện tại, ngưỡng           │
│    - Timeline: khi nào vượt, bao lâu                        │
│    - Biểu đồ metric trong 30 phút gần nhất                  │
│ 6. Operator nhấn "Acknowledge"                              │
│ 7. Alarm chuyển trạng thái RAISED → ACKNOWLEDGED            │
│ 8. Hệ thống ghi log: ai ack, thời gian                      │
│ 9. Operator xử lý sự cố (có thể gửi lệnh điều khiển)       │
│ 10. Khi metric về bình thường, alarm tự động RESOLVED       │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Không ai ack trong X phút → Escalation:                  │
│     - Sau 5 phút: thông báo Quản đốc                        │
│     - Sau 15 phút: thông báo Giám đốc SX                    │
│     - Sau 30 phút: gọi điện (SMS/voice call)                │
│ A2. Operator nhấn "Snooze": Alarm tạm ẩn trong T phút       │
│ A3. Operator nhấn "Suppress": Chọn khoảng thời gian tắt     │
│     alarm cho máy này                                        │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: Alarm được xử lý, ghi nhận, và resolved     │
└─────────────────────────────────────────────────────────────┘
```

### 5.3 UC-03: Điều Khiển Máy Từ Xa

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-03 — Điều Khiển Máy Từ Xa                       │
├─────────────────────────────────────────────────────────────┤
│ Actor: KT Vận Hành, Quản Đốc                                │
│ Trigger: User muốn gửi lệnh điều khiển đến máy              │
│ Pre-condition: Đã đăng nhập, có quyền điều khiển máy đó      │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. User mở màn hình Device Control                          │
│ 2. User chọn máy từ cây thiết bị                            │
│ 3. Hệ thống hiển thị Control Panel:                         │
│    - Trạng thái hiện tại                                    │
│    - Các nút điều khiển: Start, Stop, Set Speed, ...        │
│    - Input fields: setpoint, parameter                      │
│ 4. User nhập lệnh (vd: Set speed = 1500 RPM)               │
│ 5. Hệ thống hiển thị confirmation dialog:                   │
│    "Bạn có chắc muốn đặt tốc độ máy X thành 1500 RPM?"      │
│ 6. User xác nhận                                            │
│ 7. Command Service tạo command (PENDING)                    │
│ 8. Command Service gửi MQTT message xuống Edge Gateway      │
│ 9. Command chuyển trạng thái: PENDING → SENT               │
│10. Edge Gateway nhận lệnh, gửi xuống PLC                   │
│11. Command chuyển trạng thái: SENT → EXECUTING              │
│12. PLC thực thi lệnh, Edge Gateway gửi phản hồi            │
│13. Command chuyển trạng thái: EXECUTING → DONE (hoặc FAIL)  │
│14. UI cập nhật kết quả và trạng thái mới của máy            │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Emergency Stop:                                           │
│     - User nhấn nút "EMERGENCY STOP" (màu đỏ, nổi bật)      │
│     - Hệ thống yêu cầu xác nhận 2 bước                      │
│     - Gửi lệnh STOP khẩn cấp với priority = HIGHEST         │
│     - Bỏ qua mọi validation, gửi trực tiếp                  │
│     - Ghi audit log đặc biệt                                │
│ A2. Command Timeout:                                         │
│     - Sau X giây không có phản hồi → FAILED                  │
│     - Hiển thị thông báo lỗi cho user                       │
│     - User có thể retry                                     │
│ A3. Không có quyền: Hiển thị "Bạn không có quyền điều khiển  │
│     máy này"                                                 │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: Máy thực hiện lệnh, trạng thái cập nhật,    │
│                 command được ghi audit log                   │
└─────────────────────────────────────────────────────────────┘
```

### 5.4 UC-04: Đăng Ký Thiết Bị Mới

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-04 — Đăng Ký Thiết Bị Mới                       │
├─────────────────────────────────────────────────────────────┤
│ Actor: IT Engineer                                          │
│ Trigger: Có máy mới cần kết nối vào hệ thống                │
│ Pre-condition: Edge Gateway đã được cài đặt và online        │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. IT Engineer mở màn hình Device Management                │
│ 2. Nhấn "Thêm thiết bị"                                     │
│ 3. Điền form đăng ký:                                       │
│    - Serial Number                                          │
│    - Tên thiết bị                                           │
│    - Model/Mã hiệu                                          │
│    - Nhà máy (Hà Nội / HCM / Đà Nẵng)                       │
│    - Khu vực (Area A, B, C...)                              │
│    - Giao thức kết nối (Modbus TCP, OPC-UA, ...)            │
│    - Địa chỉ kết nối (IP:Port, OPC-UA endpoint)             │
│    - Danh sách metrics cần thu thập                         │
│ 4. Nhấn "Lưu"                                               │
│ 5. Hệ thống validate và tạo device record                   │
│ 6. Hệ thống cấp chứng chỉ X.509 cho thiết bị                │
│ 7. Edge Gateway bắt đầu thu thập dữ liệu từ máy             │
│ 8. Dữ liệu bắt đầu xuất hiện trên dashboard                 │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Serial Number đã tồn tại: Hiển thị lỗi, đề xuất kiểm tra│
│ A2. Edge Gateway offline: Thiết bị được tạo với trạng thái  │
│     "pending", sẽ active khi edge online                     │
│ A3. Tự động provisioning: Khi edge gateway phát hiện máy mới │
│     trên mạng Modbus → tự động gửi lên cloud → IT Engineer   │
│     xác nhận                                                 │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: Thiết bị được đăng ký, dữ liệu bắt đầu      │
│                 được thu thập và hiển thị                    │
└─────────────────────────────────────────────────────────────┘
```

### 5.5 UC-05: Cấu Hình Alarm Rule

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-05 — Cấu Hình Alarm Rule                        │
├─────────────────────────────────────────────────────────────┤
│ Actor: System Admin, Quản Đốc                               │
│ Trigger: Cần tạo/chỉnh sửa rule cảnh báo                    │
│ Pre-condition: Có quyền quản lý alarm rule                  │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. User mở màn hình Alarm Configuration                     │
│ 2. Nhấn "Tạo Rule mới"                                      │
│ 3. Chọn loại rule:                                          │
│    - Threshold: giá trị vượt ngưỡng                         │
│    - Rate of Change: giá trị thay đổi quá nhanh             │
│    - State Change: máy thay đổi trạng thái                  │
│    - Anomaly: phát hiện bất thường bằng thuật toán          │
│ 4. Cấu hình rule:                                           │
│    - Tên rule: "Nhiệt độ lò quá cao"                        │
│    - Metric: furnace_01.temperature                         │
│    - Điều kiện: value > 800                                 │
│    - Duration: 30 giây (phải vượt liên tục 30s mới alarm)  │
│    - Severity: Critical                                     │
│    - Group/Scope: Nhà máy Hà Nội, Khu vực A                 │
│    - Notification channels: Push + Email                    │
│    - Escalation policy: sau 5 phút escalate lên Supervisor  │
│ 5. Nhấn "Test Rule" với dữ liệu mẫu                         │
│ 6. Hệ thống mô phỏng kết quả: "Alarm sẽ kích hoạt"          │
│ 7. Nhấn "Lưu & Kích hoạt"                                   │
│ 8. Rule được deploy vào Rule Engine                         │
│ 9. Rule bắt đầu đánh giá real-time                          │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Rule conflict: Hệ thống cảnh báo nếu rule mới xung đột   │
│     với rule hiện có                                         │
│ A2. Test fails: Hiển thị lý do, cho phép điều chỉnh         │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: Rule mới được kích hoạt, sẵn sàng đánh giá   │
└─────────────────────────────────────────────────────────────┘
```

### 5.6 UC-06: Cập Nhật Firmware OTA

```
┌─────────────────────────────────────────────────────────────┐
│ Use Case: UC-06 — Cập Nhật Firmware OTA                      │
├─────────────────────────────────────────────────────────────┤
│ Actor: IT Engineer                                          │
│ Trigger: Có phiên bản firmware mới cho Edge Gateway         │
│ Pre-condition: Edge Gateway đang online                      │
├─────────────────────────────────────────────────────────────┤
│ Main Flow:                                                   │
│ 1. IT Engineer upload firmware binary lên hệ thống          │
│ 2. Hệ thống lưu trữ firmware trên MinIO                     │
│ 3. IT Engineer chọn firmware, chọn thiết bị đích:           │
│    - Một thiết bị cụ thể                                    │
│    - Nhóm thiết bị (vd: tất cả edge Hà Nội)                │
│ 4. IT Engineer chọn lịch cập nhật:                          │
│    - Cập nhật ngay                                          │
│    - Lên lịch (vd: 2:00 AM Chủ Nhật)                        │
│    - Cập nhật tuần tự (rolling update)                      │
│ 5. Hệ thống kiểm tra điều kiện:                             │
│    - Edge online?                                            │
│    - Đủ dung lượng lưu trữ?                                 │
│    - Máy đang không chạy (nếu yêu cầu)?                     │
│ 6. Hệ thống gửi firmware xuống edge                         │
│ 7. Edge tải firmware, verify checksum                       │
│ 8. Edge cài đặt firmware, khởi động lại                     │
│ 9. Edge kết nối lại cloud, báo cáo phiên bản mới           │
│10. Hệ thống cập nhật trạng thái: Success/Fail              │
│11. Nếu fail → tự động rollback firmware cũ                  │
├─────────────────────────────────────────────────────────────┤
│ Alternative Flow:                                            │
│ A1. Checksum mismatch: Hủy cập nhật, thông báo lỗi          │
│ A2. Edge không khởi động lại sau update:                     │
│     - Timeout X phút → đánh dấu FAILED                      │
│     - IT Engineer cần can thiệp thủ công                    │
│ A3. Rolling update: Cập nhật lần lượt từng edge, đảm bảo     │
│     luôn có ít nhất 1 edge hoạt động                        │
├─────────────────────────────────────────────────────────────┤
│ Post-condition: Edge Gateway được cập nhật firmware mới     │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Luồng Công Việc Người Dùng

### 6.1 Ca Sản Xuất Thông Thường

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    WORKFLOW: CA SẢN XUẤT THÔNG THƯỜNG                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ĐẦU CA (6:00 / 14:00 / 22:00)                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ Quản Đốc mở Dashboard → Kiểm tra tất cả máy online               │    │
│  │ Quản Đốc kiểm tra alarm từ ca trước (đã resolved?)               │    │
│  │ Quản Đốc ghi chú đầu ca (nếu cần)                                │    │
│  │ KT Vận Hành mở Mobile App → Xem danh sách máy phụ trách          │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                    ↓                                     │
│  TRONG CA (liên tục)                                                     │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ ⏱ Mỗi 5 phút:                                                     │    │
│  │   • KT Vận Hành liếc Mobile App kiểm tra trạng thái máy           │    │
│  │ ⏱ Mỗi 30 phút:                                                    │    │
│  │   • Quản Đốc kiểm tra Dashboard: OEE, throughput, alarm           │    │
│  │ ⏱ Mỗi giờ:                                                        │    │
│  │   • KT Vận Hành ghi nhận thông số vào checklist (nếu cần)         │    │
│  │ ⏱ Khi có alarm:                                                   │    │
│  │   • Push notification → KT Vận Hành nhận, ack, xử lý              │    │
│  │   • Nếu alarm critical → Quản Đốc cũng nhận được                  │    │
│  │   • Sau xử lý → alarm auto-resolve khi metric về bình thường      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                    ↓                                     │
│  CUỐI CA                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ Quản Đốc xem báo cáo ca: sản lượng, OEE, alarm đã xử lý         │    │
│  │ Quản Đốc in/xuất báo cáo ca (nếu cần)                            │    │
│  │ Ca sau: Quản Đốc mới xem dashboard, nắm tình hình                │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Xử Lý Sự Cố Khẩn Cấp

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    WORKFLOW: XỬ LÝ SỰ CỐ KHẨN CẤP                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PHÁT HIỆN                                                              │
│  ┌─────────────────────────┐                                            │
│  │ Alarm EMERGENCY kích    │ ← Nhiệt độ lò > 1000°C (ngưỡng khẩn)      │
│  │ hoạt                    │ ← hoặc Motor dừng đột ngột                 │
│  └────────────┬────────────┘                                            │
│               ↓                                                         │
│  THÔNG BÁO (đồng thời, < 5 giây)                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ 📱 Push Notification: tất cả operator trong khu vực              │    │
│  │ 🖥 Web Portal: Alarm Center hiển thị alarm đỏ nhấp nháy          │    │
│  │ 🔊 Âm thanh: Còi báo động (nếu tích hợp)                         │    │
│  │ 📧 Email: Gửi cho Giám đốc SX                                   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│               ↓                                                         │
│  PHẢN ỨNG (bởi KT Vận Hành gần nhất)                                   │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ 1. KT Vận Hành nhấn "Acknowledge" trên điện thoại               │    │
│  │ 2. KT Vận Hành chạy đến máy (nếu cần kiểm tra trực tiếp)       │    │
│  │ 3. Nếu nguy hiểm → Nhấn EMERGENCY STOP trên mobile              │    │
│  │    • Vuốt để xác nhận                                            │    │
│  │    • Lệnh STOP gửi ngay xuống máy                                │    │
│  │    • Máy dừng trong < 500ms (kể cả khi mất internet)            │    │
│  │ 4. Gọi điện cho Quản đốc (nếu cần)                              │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│               ↓                                                         │
│  ESCALATION (nếu không ai phản hồi)                                     │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ Sau 1 phút không ack → Thông báo Quản đốc                        │    │
│  │ Sau 5 phút không ack → Thông báo Giám đốc SX                     │    │
│  │ Sau 10 phút không ack → Gọi điện tự động                         │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│               ↓                                                         │
│  SAU SỰ CỐ                                                              │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ 1. Alarm chuyển RESOLVED khi điều kiện về bình thường            │    │
│  │ 2. Quản đốc viết ghi chú nguyên nhân & hành động đã làm         │    │
│  │ 3. Hệ thống lưu alarm history (cho phân tích sau)               │    │
│  │ 4. Tự động tạo báo cáo sự cố (incident report)                  │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Bảo Trì Dự Đoán

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    WORKFLOW: BẢO TRÌ DỰ ĐOÁN                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  HÀNG NGÀY                                                              │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ KS Bảo Trì mở Predictive Maintenance Dashboard                   │    │
│  │ Hệ thống hiển thị:                                               │    │
│  │   🔴 Critical: Máy X — 95% khả năng hỏng trong 7 ngày           │    │
│  │   🟡 Warning:  Máy Y — 70% khả năng hỏng trong 14 ngày          │    │
│  │   🟢 Healthy:  Các máy còn lại                                   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                    ↓                                     │
│  LẬP KẾ HOẠCH                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ KS Bảo Trì click vào Máy X                                       │    │
│  │ Xem chi tiết:                                                    │    │
│  │   • Feature quan trọng nhất dẫn đến dự đoán (vd: rung tăng)     │    │
│  │   • Biểu đồ trend của feature đó trong 30 ngày                  │    │
│  │   • Lịch sử bảo trì trước đây                                   │    │
│  │   • Đề xuất: "Thay vòng bi motor"                               │    │
│  │ KS Bảo Trì lên lịch bảo trì: "Thứ Bảy tuần này, 8:00-10:00"   │    │
│  │ Hệ thống tự động tạo work order                                 │    │
│  │ Hệ thống tự động suppression alarm trong thời gian bảo trì      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                    ↓                                     │
│  SAU BẢO TRÌ                                                            │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ KS Bảo Trì cập nhật work order: "Đã thay vòng bi"               │    │
│  │ Hệ thống ghi nhận vào lịch sử bảo trì                           │    │
│  │ ML model cập nhật với dữ liệu mới                               │    │
│  │ Dự đoán cho Máy X giảm xuống 🟢 Healthy                        │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Yêu Cầu Giao Diện Người Dùng

### 7.1 Web Portal — Các Màn Hình Chính

| Màn hình | Mô tả | Thành phần chính |
|----------|-------|-----------------|
| **Login** | Đăng nhập qua Keycloak SSO | Logo, form login, "Quên mật khẩu", đa ngôn ngữ |
| **Dashboard** | Tổng quan nhà máy real-time | Widget: OEE gauge, active alarms, online/offline count, throughput. Có thể chọn nhà máy |
| **Factory Detail** | Chi tiết một nhà máy | Sơ đồ layout nhà máy, màu trạng thái từng khu vực |
| **Line Detail** | Chi tiết một dây chuyền | Danh sách máy dạng card/grid, trạng thái màu, metrics chính |
| **Machine Detail** | Chi tiết một máy | Gauge/line chart cho từng metric, device shadow, control panel, alarm history |
| **Alarm Center** | Trung tâm cảnh báo | Bảng alarm active (filter, sort, ack), alarm history timeline |
| **Control Panel** | Bảng điều khiển máy | Nút Start/Stop, slider, input setpoint, confirmation dialog |
| **Reports** | Báo cáo & phân tích | OEE report, production report, energy report, custom report builder |
| **Device Management** | Quản lý thiết bị | Device tree, CRUD form, firmware management, device shadow |
| **Admin Panel** | Quản trị hệ thống | User management, role management, audit log, system config |
| **Alarm Config** | Cấu hình alarm rule | Rule builder UI, rule list, test rule |

### 7.2 Mobile App — Các Màn Hình Chính

| Màn hình | Mô tả | Thành phần chính |
|----------|-------|-----------------|
| **Splash/Login** | Màn hình khởi động & đăng nhập | Logo, biometric login (FaceID/vân tay) |
| **Dashboard** | Tổng quan mobile | KPI cards, danh sách máy gần đây, alarm count badge |
| **Machine List** | Danh sách máy theo khu vực | List với status dot (xanh/vàng/đỏ), search, filter |
| **Machine Detail** | Chi tiết máy mobile | Gauge/line chart tối ưu mobile, metrics chính |
| **Alarm Inbox** | Hộp thư alarm | List alarm kèm severity, swipe to ack |
| **Quick Control** | Điều khiển nhanh | Nút lớn: Emergency Stop, Start, Stop |
| **QR Scanner** | Quét QR code máy | Camera view, hiển thị thông tin máy sau khi quét |
| **Settings** | Cài đặt | Notification preferences, chọn nhà máy, ngôn ngữ, đăng xuất |

### 7.3 Nguyên Tắc Thiết Kế UI/UX

1. **Đơn giản & Trực quan**: Giao diện phải dễ hiểu ngay lần đầu sử dụng, không cần đào tạo dài
2. **Màu sắc ngữ nghĩa**: Xanh lá = bình thường, Vàng = cảnh báo, Đỏ = nguy hiểm
3. **Phản hồi tức thì**: Mọi thao tác phải có phản hồi trong < 200ms (loading state)
4. **Mobile-first cho operator**: Tối ưu cho thao tác bằng ngón tay cái, nút bấm lớn
5. **Dark mode**: Hỗ trợ dark mode cho môi trường nhà máy thiếu sáng
6. **Offline resilience**: Hiển thị trạng thái mất kết nối, dữ liệu cache
7. **Accessibility**: Tuân thủ WCAG 2.1 AA, hỗ trợ screen reader
8. **Đa ngôn ngữ**: Tiếng Việt (mặc định) + Tiếng Anh

---

## 8. Yêu Cầu Trải Nghiệm Người Dùng

### 8.1 Yêu Cầu Hiệu Năng UX

| Chỉ số | Mục tiêu | Ghi chú |
|--------|----------|---------|
| Time to First Byte (TTFB) | < 500ms | API Gateway → BFF → Service |
| Dashboard Load Time | < 2 giây | Bao gồm tất cả widgets |
| Real-time Data Refresh | ≤ 2 giây | WebSocket/SSE từ Kafka |
| Chart Rendering | < 500ms | Cho 10,000 data points |
| Mobile Push Notification | < 3 giây | Từ alarm trigger → push received |
| Page Transition | < 200ms | SPA routing |

### 8.2 Yêu Cầu Usability

| Yêu cầu | Mô tả |
|---------|-------|
| **SUS Score** | System Usability Scale ≥ 75 |
| **Error Rate** | < 5% người dùng mắc lỗi khi thực hiện task chính |
| **Learnability** | Người dùng mới có thể hoàn thành task cơ bản trong < 10 phút |
| **Satisfaction** | CSAT ≥ 4.0/5.0 sau 1 tháng sử dụng |

### 8.3 Trạng Thái Hệ Thống

Hệ thống phải luôn hiển thị rõ ràng:

- **Trạng thái kết nối**: Online / Offline / Reconnecting
- **Trạng thái đồng bộ**: Đã đồng bộ / Đang đồng bộ / Lỗi đồng bộ
- **Tuổi dữ liệu**: "Dữ liệu cập nhật cách đây X giây"
- **Trạng thái loading**: Skeleton loader, không phải màn hình trắng
- **Trạng thái lỗi**: Thông báo lỗi rõ ràng, actionable (nút retry, link hỗ trợ)
- **Trạng thái trống**: "Chưa có dữ liệu" thay vì màn hình trống

---

## 9. Yêu Cầu Báo Cáo

### 9.1 Báo Cáo Tiêu Chuẩn

| Báo cáo | Mô tả | Tần suất | Định dạng |
|---------|-------|----------|-----------|
| **Báo cáo ca** | Sản lượng, OEE, alarm trong ca | Mỗi ca | PDF, Web |
| **Báo cáo ngày** | Tổng hợp sản xuất trong ngày | Hàng ngày | PDF, Excel, Web |
| **Báo cáo tuần** | Xu hướng OEE, throughput, downtime | Hàng tuần | PDF, Excel |
| **Báo cáo tháng** | KPIs, so sánh với tháng trước, kế hoạch | Hàng tháng | PDF, Excel |
| **Báo cáo OEE** | OEE breakdown: A × P × Q | Theo yêu cầu | PDF, Web |
| **Báo cáo năng lượng** | Tiêu thụ điện theo máy/khu vực | Hàng tuần | PDF, Web |
| **Báo cáo sự cố** | Incident report cho alarm Critical/Emergency | Tự động khi resolve | PDF |
| **Báo cáo bảo trì** | Lịch sử bảo trì, predictive maintenance status | Hàng tuần | Web |

### 9.2 Báo Cáo Tùy Chỉnh

Người dùng có thể tạo báo cáo tùy chỉnh:

1. Chọn loại báo cáo (bảng, biểu đồ, dashboard)
2. Chọn metrics (nhiệt độ, áp suất, OEE, throughput, ...)
3. Chọn phạm vi (nhà máy, khu vực, dây chuyền, máy)
4. Chọn khoảng thời gian
5. Chọn group by (theo giờ, ca, ngày, tuần, tháng)
6. Chọn filter (chỉ máy có OEE < 80%, ...)
7. Preview → Export (PDF, Excel, CSV)

---

## 10. Phụ Lục

### A. Ma Trận Vai Trò - Quyền Hạn

| Chức năng | Admin | Supervisor | Operator | Viewer |
|-----------|-------|------------|----------|--------|
| Xem Dashboard | ✅ | ✅ | ✅ | ✅ |
| Xem chi tiết máy | ✅ | ✅ | ✅ | ✅ |
| Điều khiển máy cơ bản | ✅ | ✅ | ✅ | ❌ |
| Điều khiển máy nâng cao | ✅ | ✅ | ❌ | ❌ |
| Emergency Stop | ✅ | ✅ | ✅ | ❌ |
| Acknowledge Alarm | ✅ | ✅ | ✅ | ❌ |
| Cấu hình Alarm Rule | ✅ | ✅ | ❌ | ❌ |
| Quản lý thiết bị | ✅ | ❌ | ❌ | ❌ |
| Quản lý người dùng | ✅ | ❌ | ❌ | ❌ |
| Xem báo cáo | ✅ | ✅ | ❌ | ✅ |
| Export báo cáo | ✅ | ✅ | ❌ | ❌ |
| Quản lý firmware | ✅ | ❌ | ❌ | ❌ |
| Xem audit log | ✅ | ❌ | ❌ | ❌ |

### B. Yêu Cầu Notification Theo Vai Trò

| Sự kiện | Operator | Supervisor | Admin | Kênh |
|---------|----------|------------|-------|------|
| Alarm Warning | Push | — | — | Push |
| Alarm Critical | Push + SMS | Push | Email | Push, SMS, Email |
| Alarm Emergency | Push + SMS | Push + SMS | Push + SMS | Push, SMS, Email |
| Alarm Escalated | — | Push + SMS | — | Push, SMS |
| Device Offline | — | Push | Email | Push, Email |
| Firmware Update Failed | — | — | Email | Email |
| System Error | — | — | Push + SMS | Push, SMS |

### C. Tài Liệu Tham Khảo

- [Business Requirements Document (BRD)](BRD.md)
- [Software Requirements Specification (SRS)](SRS.md)
- [Kiến trúc hệ thống (README)](../README.md)
- README của từng service

---

> **Tài liệu này là tài sản trí tuệ của Industrial IoT Platform. Không được sao chép hoặc phân phối khi chưa có sự đồng ý.**
