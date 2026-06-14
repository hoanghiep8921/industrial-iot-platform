package model

import "time"

type AlarmStatus string

const (
	StatusRaised       AlarmStatus = "RAISED"
	StatusAcknowledged AlarmStatus = "ACKNOWLEDGED"
	StatusResolved     AlarmStatus = "RESOLVED"
)

type AlarmSeverity string

const (
	SeverityWarning   AlarmSeverity = "warning"
	SeverityCritical  AlarmSeverity = "critical"
	SeverityEmergency AlarmSeverity = "emergency"
)

type Alarm struct {
	ID           string        `json:"id"`
	RuleID       string        `json:"rule_id"`
	RuleName     string        `json:"rule_name,omitempty"` // populated on joins
	MachineID    string        `json:"machine_id"`
	MetricName   string        `json:"metric_name"`
	TriggerValue float64       `json:"trigger_value"`
	TriggerTime  time.Time     `json:"trigger_time"`
	AckTime      *time.Time    `json:"ack_time,omitempty"`
	AckBy        string        `json:"ack_by,omitempty"`
	ResolveTime  *time.Time    `json:"resolve_time,omitempty"`
	Status       AlarmStatus   `json:"status"`
	Severity     AlarmSeverity `json:"severity"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type AlarmEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"` // e.g. "alarm.raised", "alarm.acknowledged", "alarm.resolved"
	Timestamp time.Time `json:"timestamp"`
	Alarm     *Alarm    `json:"alarm"`
}

type AlarmFilter struct {
	FactoryID string      `json:"factory_id"`
	MachineID string      `json:"machine_id"`
	Status    AlarmStatus `json:"status"`
	Severity  string      `json:"severity"`
	Limit     int         `json:"limit"`
	Page      int         `json:"page"`
}
