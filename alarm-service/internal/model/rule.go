package model

import (
	"errors"
	"time"
)

type AlarmRule struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Type              string    `json:"type"` // e.g. "threshold"
	MetricName        string    `json:"metric_name"`
	ConditionOperator string    `json:"condition_operator"` // '>', '>=', '<', '<=', '==', '!='
	ConditionValue    float64   `json:"condition_value"`
	DurationSeconds   int       `json:"duration_seconds"`
	Severity          string    `json:"severity"` // "warning", "critical", "emergency"
	FactoryID         string    `json:"factory_id"`
	IsEnabled         bool      `json:"is_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateRuleRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Type              string  `json:"type"`
	MetricName        string  `json:"metric_name"`
	ConditionOperator string  `json:"condition_operator"`
	ConditionValue    float64 `json:"condition_value"`
	DurationSeconds   int     `json:"duration_seconds"`
	Severity          string  `json:"severity"`
	FactoryID         string  `json:"factory_id"`
}

func (r *CreateRuleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.MetricName == "" {
		return errors.New("metric_name is required")
	}
	if r.ConditionOperator == "" {
		return errors.New("condition_operator is required")
	}
	switch r.ConditionOperator {
	case ">", ">=", "<", "<=", "==", "!=":
		// Valid
	default:
		return errors.New("invalid condition_operator: must be one of >, >=, <, <=, ==, !=")
	}
	if r.Severity == "" {
		r.Severity = "warning"
	}
	switch r.Severity {
	case "warning", "critical", "emergency":
		// Valid
	default:
		return errors.New("invalid severity: must be one of warning, critical, emergency")
	}
	if r.Type == "" {
		r.Type = "threshold"
	}
	return nil
}

type UpdateRuleRequest struct {
	Name              *string  `json:"name"`
	Description       *string  `json:"description"`
	Type              *string  `json:"type"`
	MetricName        *string  `json:"metric_name"`
	ConditionOperator *string  `json:"condition_operator"`
	ConditionValue    *float64 `json:"condition_value"`
	DurationSeconds   *int     `json:"duration_seconds"`
	Severity          *string  `json:"severity"`
	FactoryID         *string  `json:"factory_id"`
	IsEnabled         *bool    `json:"is_enabled"`
}

func (r *UpdateRuleRequest) Validate() error {
	if r.ConditionOperator != nil {
		switch *r.ConditionOperator {
		case ">", ">=", "<", "<=", "==", "!=":
			// Valid
		default:
			return errors.New("invalid condition_operator: must be one of >, >=, <, <=, ==, !=")
		}
	}
	if r.Severity != nil {
		switch *r.Severity {
		case "warning", "critical", "emergency":
			// Valid
		default:
			return errors.New("invalid severity: must be one of warning, critical, emergency")
		}
	}
	return nil
}
