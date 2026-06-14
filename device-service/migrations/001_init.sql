-- ============================================================
-- Device Service - Database Schema
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========== Devices ==========
CREATE TABLE IF NOT EXISTS devices (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_number    VARCHAR(100) UNIQUE NOT NULL,
    name             VARCHAR(255) NOT NULL,
    model            VARCHAR(100),
    vendor           VARCHAR(100),
    factory_id       VARCHAR(50) NOT NULL,
    area             VARCHAR(100),
    line             VARCHAR(100),
    protocol         VARCHAR(50) DEFAULT 'mqtt',
    capabilities     JSONB DEFAULT '[]',
    status           VARCHAR(20) DEFAULT 'offline',
    firmware_version VARCHAR(50),
    metadata         JSONB DEFAULT '{}',
    last_seen_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_devices_factory   ON devices(factory_id);
CREATE INDEX IF NOT EXISTS idx_devices_status    ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_protocol  ON devices(protocol);
CREATE INDEX IF NOT EXISTS idx_devices_area      ON devices(factory_id, area);

-- ========== Device Shadows ==========
CREATE TABLE IF NOT EXISTS device_shadows (
    device_id       UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    desired_state   JSONB NOT NULL DEFAULT '{}',
    reported_state  JSONB NOT NULL DEFAULT '{}',
    delta           JSONB NOT NULL DEFAULT '{}',
    version         BIGINT DEFAULT 1,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ========== Device Connections ==========
CREATE TABLE IF NOT EXISTS device_connections (
    device_id       UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    connected       BOOLEAN DEFAULT FALSE,
    connected_at    TIMESTAMPTZ,
    disconnected_at TIMESTAMPTZ,
    ip_address      VARCHAR(45),
    mqtt_client_id  VARCHAR(255),
    last_heartbeat  TIMESTAMPTZ
);

-- ========== Firmwares ==========
CREATE TABLE IF NOT EXISTS firmwares (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version         VARCHAR(50) NOT NULL,
    model           VARCHAR(100),
    description     TEXT,
    file_name       VARCHAR(255),
    file_size       BIGINT DEFAULT 0,
    checksum_sha256 VARCHAR(64),
    minio_path      VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'draft',
    created_by      VARCHAR(100),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    released_at     TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_firmwares_version_model ON firmwares(version, COALESCE(model, ''));
CREATE INDEX IF NOT EXISTS idx_firmwares_status ON firmwares(status);

-- ========== Firmware Updates ==========
CREATE TABLE IF NOT EXISTS firmware_updates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firmware_id     UUID NOT NULL REFERENCES firmwares(id),
    device_id       UUID NOT NULL REFERENCES devices(id),
    campaign_name   VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'pending',
    from_version    VARCHAR(50),
    to_version      VARCHAR(50) NOT NULL,
    priority        INT DEFAULT 0,
    retry_count     INT DEFAULT 0,
    max_retries     INT DEFAULT 3,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    error_message   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fw_updates_device  ON firmware_updates(device_id);
CREATE INDEX IF NOT EXISTS idx_fw_updates_status  ON firmware_updates(status);
CREATE INDEX IF NOT EXISTS idx_fw_updates_campaign ON firmware_updates(campaign_name);
