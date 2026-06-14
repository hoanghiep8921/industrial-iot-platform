-- ============================================================
-- Alarm Service - Database Schema
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========== Alarm Rules ==========
CREATE TABLE IF NOT EXISTS alarm_rules (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name               VARCHAR(255) NOT NULL,
    description        TEXT,
    type               VARCHAR(50) NOT NULL DEFAULT 'threshold',
    metric_name        VARCHAR(100) NOT NULL,
    condition_operator VARCHAR(10) NOT NULL, -- '>', '>=', '<', '<=', '==', '!='
    condition_value    DOUBLE PRECISION NOT NULL,
    duration_seconds   INT NOT NULL DEFAULT 0,
    severity           VARCHAR(20) NOT NULL DEFAULT 'warning', -- 'warning', 'critical', 'emergency'
    factory_id         VARCHAR(50),
    is_enabled         BOOLEAN DEFAULT TRUE,
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alarm_rules_metric ON alarm_rules(metric_name);
CREATE INDEX IF NOT EXISTS idx_alarm_rules_factory ON alarm_rules(factory_id);

-- ========== Alarms ==========
CREATE TABLE IF NOT EXISTS alarms (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id        UUID NOT NULL REFERENCES alarm_rules(id) ON DELETE CASCADE,
    machine_id     VARCHAR(100) NOT NULL,
    metric_name    VARCHAR(100) NOT NULL,
    trigger_value  DOUBLE PRECISION NOT NULL,
    trigger_time   TIMESTAMPTZ NOT NULL,
    ack_time       TIMESTAMPTZ,
    ack_by         VARCHAR(100),
    resolve_time   TIMESTAMPTZ,
    status         VARCHAR(20) NOT NULL DEFAULT 'RAISED', -- 'RAISED', 'ACKNOWLEDGED', 'RESOLVED'
    severity       VARCHAR(20) NOT NULL DEFAULT 'warning',
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alarms_status ON alarms(status);
CREATE INDEX IF NOT EXISTS idx_alarms_machine ON alarms(machine_id);
CREATE INDEX IF NOT EXISTS idx_alarms_rule ON alarms(rule_id);
