-- ============================================================
-- Telemetry Service - Database Migration
-- ============================================================

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================================
-- Main telemetry hypertable
-- ============================================================
CREATE TABLE IF NOT EXISTS telemetry (
    time        TIMESTAMPTZ NOT NULL,
    device_id   TEXT NOT NULL,          -- machine_id
    factory_id  TEXT NOT NULL,
    area        TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    value       DOUBLE PRECISION,
    unit        TEXT,
    quality     SMALLINT DEFAULT 0,     -- 0=good, 1=uncertain, 2=bad
    tags        JSONB,
    raw_payload JSONB
);

-- Convert to hypertable partitioned by time (7-day chunks)
SELECT create_hypertable('telemetry', 'time',
    chunk_time_interval => INTERVAL '7 days',
    if_not_exists       => TRUE
);

-- ============================================================
-- Indexes for common query patterns
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_telemetry_device_time
    ON telemetry (device_id, time DESC);

CREATE INDEX IF NOT EXISTS idx_telemetry_factory_area
    ON telemetry (factory_id, area, time DESC);

CREATE INDEX IF NOT EXISTS idx_telemetry_metric
    ON telemetry (metric_name, time DESC);

-- ============================================================
-- Latest value table (for real-time dashboard)
-- ============================================================
CREATE TABLE IF NOT EXISTS device_latest_value (
    device_id   TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    factory_id  TEXT NOT NULL,
    area        TEXT NOT NULL,
    value       DOUBLE PRECISION,
    unit        TEXT,
    quality     SMALLINT,
    tags        JSONB,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (device_id, metric_name)
);

-- ============================================================
-- Compression policy (compress data older than 7 days)
-- ============================================================
SELECT add_compression_policy('telemetry', INTERVAL '7 days',
    if_not_exists => TRUE
);

-- ============================================================
-- Retention policy (keep 2 years of data)
-- ============================================================
SELECT add_retention_policy('telemetry', INTERVAL '2 years',
    if_not_exists => TRUE
);
