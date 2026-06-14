-- Chèn dữ liệu mẫu cho Alarm Rules
INSERT INTO alarm_rules (id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled)
VALUES
    ('c0a80101-0000-0000-0000-000000000001', 'Nhiệt độ lò quá cao', 'Cảnh báo khi nhiệt độ lò đúc vượt quá 100 độ C', 'threshold', 'temperature', '>', 100.0, 5, 'critical', 'factory-01', true),
    ('c0a80101-0000-0000-0000-000000000002', 'Áp suất nồi hơi nguy hiểm', 'Áp suất hệ thống vượt quá giới hạn an toàn 80 bar', 'threshold', 'pressure', '>', 80.0, 0, 'emergency', 'factory-02', true),
    ('c0a80101-0000-0000-0000-000000000003', 'Tốc độ động cơ thấp', 'Cảnh báo khi vòng quay động cơ giảm dưới 500 RPM', 'threshold', 'rpm', '<', 500.0, 10, 'warning', 'factory-01', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    metric_name = EXCLUDED.metric_name,
    condition_operator = EXCLUDED.condition_operator,
    condition_value = EXCLUDED.condition_value,
    duration_seconds = EXCLUDED.duration_seconds,
    severity = EXCLUDED.severity,
    factory_id = EXCLUDED.factory_id,
    is_enabled = EXCLUDED.is_enabled;

-- Chèn dữ liệu mẫu cho Alarms
INSERT INTO alarms (id, rule_id, machine_id, metric_name, trigger_value, trigger_time, ack_time, ack_by, resolve_time, status, severity)
VALUES
    -- Alarm 1: Đang hoạt động (RAISED)
    ('e0a80101-0000-0000-0000-000000000001', 'c0a80101-0000-0000-0000-000000000001', 'furnace-01', 'temperature', 105.5, NOW() - INTERVAL '15 minutes', NULL, NULL, NULL, 'RAISED', 'critical'),
    
    -- Alarm 2: Đang được xử lý (ACKNOWLEDGED)
    ('e0a80101-0000-0000-0000-000000000002', 'c0a80101-0000-0000-0000-000000000002', 'boiler-02', 'pressure', 85.2, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '20 minutes', 'hoanghiep', NULL, 'ACKNOWLEDGED', 'emergency'),
    
    -- Alarm 3: Đã xử lý xong (RESOLVED) - sẽ không hiển thị ở tab Active
    ('e0a80101-0000-0000-0000-000000000003', 'c0a80101-0000-0000-0000-000000000003', 'motor-01', 'rpm', 420.0, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour 45 minutes', 'operator_01', NOW() - INTERVAL '1 hour', 'RESOLVED', 'warning')
ON CONFLICT (id) DO UPDATE SET
    rule_id = EXCLUDED.rule_id,
    machine_id = EXCLUDED.machine_id,
    metric_name = EXCLUDED.metric_name,
    trigger_value = EXCLUDED.trigger_value,
    trigger_time = EXCLUDED.trigger_time,
    ack_time = EXCLUDED.ack_time,
    ack_by = EXCLUDED.ack_by,
    resolve_time = EXCLUDED.resolve_time,
    status = EXCLUDED.status,
    severity = EXCLUDED.severity;
