-- 1. Chèn thiết bị Máy cắt Laser 3015-3000W vào Device Registry
INSERT INTO devices (id, serial_number, name, model, vendor, factory_id, area, line, protocol, capabilities, status, firmware_version, metadata)
VALUES (
    'd0a80101-0000-0000-0000-000000000001',
    'LASER-3015-3000W-001',
    'Máy Cắt Laser 3015-3000W (DT-3015)',
    'DT-3015',
    'Datyso',
    'factory-01',
    'cutting',
    'Line 1',
    'mqtt',
    '[
        {"name": "laser_power", "unit": "W", "type": "float"},
        {"name": "laser_temp", "unit": "°C", "type": "float"},
        {"name": "gas_pressure", "unit": "bar", "type": "float"},
        {"name": "cutting_speed", "unit": "m/min", "type": "float"},
        {"name": "nozzle_height", "unit": "mm", "type": "float"}
    ]'::jsonb,
    'online',
    'v1.0.0',
    '{
        "laser_type": "Fiber",
        "power_rating": "3000W",
        "dimensions": "3000mm x 1500mm",
        "positioning_accuracy": "±0.03mm",
        "max_speed": "120m/min",
        "max_steel_thickness_mm": 20,
        "max_inox_thickness_mm": 12
    }'::jsonb
)
ON CONFLICT (serial_number) DO UPDATE SET
    name = EXCLUDED.name,
    model = EXCLUDED.model,
    vendor = EXCLUDED.vendor,
    capabilities = EXCLUDED.capabilities,
    status = EXCLUDED.status,
    metadata = EXCLUDED.metadata;

-- 2. Khởi tạo Device Shadow (Digital Twin) cho thiết bị
INSERT INTO device_shadows (device_id, desired_state, reported_state, delta, version)
VALUES (
    'd0a80101-0000-0000-0000-000000000001',
    '{"assist_gas": "oxygen", "laser_power_limit": 3000}'::jsonb,
    '{"assist_gas": "oxygen", "laser_power_limit": 3000, "laser_power": 2850, "laser_temp": 45.2, "gas_pressure": 12.5, "cutting_speed": 15.0}'::jsonb,
    '{}'::jsonb,
    1
)
ON CONFLICT (device_id) DO UPDATE SET
    reported_state = EXCLUDED.reported_state;

-- 3. Chèn Alarm Rules mẫu cho máy cắt laser này
INSERT INTO alarm_rules (id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled)
VALUES
    -- Cảnh báo quá nhiệt nguồn Laser (> 55°C trong 5 giây)
    ('c0a80101-0000-0000-0000-000000001001', 'Nhiệt độ nguồn Laser quá cao', 'Cảnh báo khi nhiệt độ buồng laser vượt quá 55 độ C liên tục 5s', 'threshold', 'laser_temp', '>', 55.0, 5, 'critical', 'factory-01', true),
    
    -- Cảnh báo áp suất khí cắt quá thấp (< 4.0 bar)
    ('c0a80101-0000-0000-0000-000000001002', 'Áp suất khí hỗ trợ thấp', 'Áp suất khí cắt (N2/O2) quá thấp không đảm bảo mạch cắt sạch', 'threshold', 'gas_pressure', '<', 4.0, 0, 'warning', 'factory-01', true),
    
    -- Cảnh báo công suất laser vượt ngưỡng cho phép (> 3000W)
    ('c0a80101-0000-0000-0000-000000001003', 'Quá tải công suất Laser', 'Công suất laser vượt quá định mức 3000W', 'threshold', 'laser_power', '>', 3000.0, 0, 'emergency', 'factory-01', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    condition_value = EXCLUDED.condition_value,
    duration_seconds = EXCLUDED.duration_seconds,
    severity = EXCLUDED.severity,
    is_enabled = EXCLUDED.is_enabled;

-- 4. Chèn 1 Cảnh Báo Đang Hoạt Động (RAISED) mẫu: Nhiệt độ laser đang ở mức 58.3°C (> 55°C)
INSERT INTO alarms (id, rule_id, machine_id, metric_name, trigger_value, trigger_time, status, severity)
VALUES (
    'e0a80101-0000-0000-0000-000000001001',
    'c0a80101-0000-0000-0000-000000001001',
    'LASER-3015-3000W-001',
    'laser_temp',
    58.3,
    NOW() - INTERVAL '3 minutes',
    'RAISED',
    'critical'
)
ON CONFLICT (id) DO UPDATE SET
    trigger_value = EXCLUDED.trigger_value,
    trigger_time = EXCLUDED.trigger_time,
    status = EXCLUDED.status;
