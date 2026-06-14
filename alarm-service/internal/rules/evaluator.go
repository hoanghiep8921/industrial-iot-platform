package rules

import (
	"context"
	"time"

	"github.com/industrial-iot/alarm-service/internal/consumer"
	"github.com/industrial-iot/alarm-service/internal/model"
	"go.uber.org/zap"
)

type DBRepository interface {
	GetRulesByMetric(ctx context.Context, metricName string) ([]*model.AlarmRule, error)
	CreateAlarm(ctx context.Context, ruleID string, machineID string, metricName string, triggerVal float64, triggerTime time.Time, severity string) (*model.Alarm, error)
	ResolveAlarm(ctx context.Context, id string, resolveTime time.Time) (*model.Alarm, error)
	GetActiveAlarm(ctx context.Context, ruleID string, machineID string) (*model.Alarm, error)
}

type CacheRepository interface {
	GetFailureStartTime(ctx context.Context, ruleID string, machineID string) (*time.Time, error)
	SetFailureStartTime(ctx context.Context, ruleID string, machineID string, t time.Time) error
	ClearFailureStartTime(ctx context.Context, ruleID string, machineID string) error
	GetActiveAlarmID(ctx context.Context, ruleID string, machineID string) (string, error)
	SetActiveAlarmID(ctx context.Context, ruleID string, machineID string, alarmID string) error
	ClearActiveAlarmID(ctx context.Context, ruleID string, machineID string) error
}

type EventPublisher interface {
	PublishAlarmEvent(ctx context.Context, eventType string, alarm *model.Alarm) error
}

type Evaluator struct {
	db        DBRepository
	cache     CacheRepository
	publisher EventPublisher
	logger    *zap.Logger
}

func NewEvaluator(db DBRepository, cache CacheRepository, pub EventPublisher, logger *zap.Logger) *Evaluator {
	return &Evaluator{
		db:        db,
		cache:     cache,
		publisher: pub,
		logger:    logger,
	}
}

// Evaluate checks an incoming telemetry message against active alarm rules
func (e *Evaluator) Evaluate(ctx context.Context, msg *consumer.TelemetryMessage) error {
	// Find rules that match this metric
	rules, err := e.db.GetRulesByMetric(ctx, msg.Metric)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if !rule.IsEnabled {
			continue
		}

		// Filter rule by FactoryID if factory is specified in rule
		if rule.FactoryID != "" && rule.FactoryID != msg.FactoryID {
			continue
		}

		isViolated := checkCondition(msg.Value, rule.ConditionOperator, rule.ConditionValue)

		if isViolated {
			if err := e.handleViolation(ctx, rule, msg); err != nil {
				e.logger.Error("Failed to handle violation", zap.Error(err), zap.String("rule_id", rule.ID), zap.String("machine_id", msg.MachineID))
			}
		} else {
			if err := e.handleNormal(ctx, rule, msg); err != nil {
				e.logger.Error("Failed to handle normal recovery", zap.Error(err), zap.String("rule_id", rule.ID), zap.String("machine_id", msg.MachineID))
			}
		}
	}

	return nil
}

func (e *Evaluator) handleViolation(ctx context.Context, rule *model.AlarmRule, msg *consumer.TelemetryMessage) error {
	// Check if alarm is already active (RAISED or ACKNOWLEDGED)
	activeAlarmID, err := e.cache.GetActiveAlarmID(ctx, rule.ID, msg.MachineID)
	if err != nil {
		return err
	}

	if activeAlarmID != "" {
		// Alarm is already active, no need to trigger again
		return nil
	}

	// If duration is 0, raise immediately
	if rule.DurationSeconds == 0 {
		return e.raiseAlarm(ctx, rule, msg)
	}

	// For duration > 0, check failure start time
	failStart, err := e.cache.GetFailureStartTime(ctx, rule.ID, msg.MachineID)
	if err != nil {
		return err
	}

	if failStart == nil {
		// Record the first failure timestamp
		err = e.cache.SetFailureStartTime(ctx, rule.ID, msg.MachineID, msg.Timestamp)
		if err != nil {
			return err
		}
		e.logger.Debug("Metric violation started",
			zap.String("rule_id", rule.ID),
			zap.String("machine_id", msg.MachineID),
			zap.Float64("value", msg.Value),
			zap.Time("start_time", msg.Timestamp),
		)
		return nil
	}

	// Compare time elapsed since first failure
	elapsed := msg.Timestamp.Sub(*failStart)
	threshold := time.Duration(rule.DurationSeconds) * time.Second

	if elapsed >= threshold {
		// Condition met, raise the alarm!
		if err := e.raiseAlarm(ctx, rule, msg); err != nil {
			return err
		}
		// Clear failure start time once alarm is active
		return e.cache.ClearFailureStartTime(ctx, rule.ID, msg.MachineID)
	}

	return nil
}

func (e *Evaluator) handleNormal(ctx context.Context, rule *model.AlarmRule, msg *consumer.TelemetryMessage) error {
	// 1. Clear any transient failure state tracker in Redis
	err := e.cache.ClearFailureStartTime(ctx, rule.ID, msg.MachineID)
	if err != nil {
		e.logger.Warn("Failed to clear failure start time", zap.Error(err))
	}

	// 2. Check if we have an active alarm that needs to be resolved
	activeAlarmID, err := e.cache.GetActiveAlarmID(ctx, rule.ID, msg.MachineID)
	if err != nil {
		return err
	}

	if activeAlarmID == "" {
		// Fallback: Check PostgreSQL database directly for active alarm
		activeAlarm, err := e.db.GetActiveAlarm(ctx, rule.ID, msg.MachineID)
		if err != nil {
			return err
		}
		if activeAlarm == nil {
			return nil
		}
		activeAlarmID = activeAlarm.ID
		// Self-healing: restore the active alarm ID to Redis cache
		e.cache.SetActiveAlarmID(ctx, rule.ID, msg.MachineID, activeAlarmID)
	}

	// Resolve the alarm in DB
	resolvedAlarm, err := e.db.ResolveAlarm(ctx, activeAlarmID, msg.Timestamp)
	if err != nil {
		return err
	}

	// Clear active alarm in Redis cache
	if err := e.cache.ClearActiveAlarmID(ctx, rule.ID, msg.MachineID); err != nil {
		e.logger.Error("Failed to clear active alarm cache", zap.Error(err))
	}

	// Publish resolved event to Kafka
	return e.publisher.PublishAlarmEvent(ctx, "alarm.resolved", resolvedAlarm)
}

func (e *Evaluator) raiseAlarm(ctx context.Context, rule *model.AlarmRule, msg *consumer.TelemetryMessage) error {
	alarm, err := e.db.CreateAlarm(ctx, rule.ID, msg.MachineID, msg.Metric, msg.Value, msg.Timestamp, rule.Severity)
	if err != nil {
		return err
	}

	// Store active alarm ID in Redis
	if err := e.cache.SetActiveAlarmID(ctx, rule.ID, msg.MachineID, alarm.ID); err != nil {
		e.logger.Error("Failed to cache active alarm ID", zap.Error(err))
	}

	// Publish raised event to Kafka
	return e.publisher.PublishAlarmEvent(ctx, "alarm.raised", alarm)
}

func checkCondition(val float64, op string, limit float64) bool {
	switch op {
	case ">":
		return val > limit
	case ">=":
		return val >= limit
	case "<":
		return val < limit
	case "<=":
		return val <= limit
	case "==":
		return val == limit
	case "!=":
		return val != limit
	default:
		return false
	}
}
