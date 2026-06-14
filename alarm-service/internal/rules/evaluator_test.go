package rules

import (
	"context"
	"testing"
	"time"

	"github.com/industrial-iot/alarm-service/internal/consumer"
	"github.com/industrial-iot/alarm-service/internal/model"
	"go.uber.org/zap"
)

// --- Mocks ---

type MockDB struct {
	rules        []*model.AlarmRule
	createdAlarm *model.Alarm
	resolvedID   string
	resolveTime  time.Time
	err          error
}

func (m *MockDB) GetRulesByMetric(ctx context.Context, metricName string) ([]*model.AlarmRule, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []*model.AlarmRule
	for _, r := range m.rules {
		if r.MetricName == metricName {
			res = append(res, r)
		}
	}
	return res, nil
}

func (m *MockDB) CreateAlarm(ctx context.Context, ruleID string, machineID string, metricName string, triggerVal float64, triggerTime time.Time, severity string) (*model.Alarm, error) {
	alarm := &model.Alarm{
		ID:           "alarm-uuid-123",
		RuleID:       ruleID,
		MachineID:    machineID,
		MetricName:   metricName,
		TriggerValue: triggerVal,
		TriggerTime:  triggerTime,
		Status:       model.StatusRaised,
		Severity:     model.AlarmSeverity(severity),
		CreatedAt:    triggerTime,
		UpdatedAt:    triggerTime,
	}
	m.createdAlarm = alarm
	return alarm, nil
}

func (m *MockDB) ResolveAlarm(ctx context.Context, id string, resolveTime time.Time) (*model.Alarm, error) {
	m.resolvedID = id
	m.resolveTime = resolveTime
	return &model.Alarm{
		ID:          id,
		Status:      model.StatusResolved,
		ResolveTime: &resolveTime,
	}, nil
}

func (m *MockDB) GetActiveAlarm(ctx context.Context, ruleID string, machineID string) (*model.Alarm, error) {
	if m.createdAlarm != nil && m.createdAlarm.RuleID == ruleID && m.createdAlarm.MachineID == machineID && m.createdAlarm.Status != model.StatusResolved {
		return m.createdAlarm, nil
	}
	return nil, nil
}

type MockCache struct {
	failureTimes map[string]time.Time
	activeAlarms map[string]string
}

func NewMockCache() *MockCache {
	return &MockCache{
		failureTimes: make(map[string]time.Time),
		activeAlarms: make(map[string]string),
	}
}

func (m *MockCache) GetFailureStartTime(ctx context.Context, ruleID string, machineID string) (*time.Time, error) {
	key := ruleID + ":" + machineID
	t, exists := m.failureTimes[key]
	if !exists {
		return nil, nil
	}
	return &t, nil
}

func (m *MockCache) SetFailureStartTime(ctx context.Context, ruleID string, machineID string, t time.Time) error {
	key := ruleID + ":" + machineID
	m.failureTimes[key] = t
	return nil
}

func (m *MockCache) ClearFailureStartTime(ctx context.Context, ruleID string, machineID string) error {
	key := ruleID + ":" + machineID
	delete(m.failureTimes, key)
	return nil
}

func (m *MockCache) GetActiveAlarmID(ctx context.Context, ruleID string, machineID string) (string, error) {
	key := ruleID + ":" + machineID
	return m.activeAlarms[key], nil
}

func (m *MockCache) SetActiveAlarmID(ctx context.Context, ruleID string, machineID string, alarmID string) error {
	key := ruleID + ":" + machineID
	m.activeAlarms[key] = alarmID
	return nil
}

func (m *MockCache) ClearActiveAlarmID(ctx context.Context, ruleID string, machineID string) error {
	key := ruleID + ":" + machineID
	delete(m.activeAlarms, key)
	return nil
}

type MockPublisher struct {
	events []model.AlarmEvent
}

func (m *MockPublisher) PublishAlarmEvent(ctx context.Context, eventType string, alarm *model.Alarm) error {
	m.events = append(m.events, model.AlarmEvent{
		EventType: eventType,
		Alarm:     alarm,
	})
	return nil
}

// --- Unit Tests ---

func TestEvaluator_InstantAlarm(t *testing.T) {
	logger := zap.NewNop()
	db := &MockDB{
		rules: []*model.AlarmRule{
			{
				ID:                "rule-1",
				Name:              "High Temperature",
				MetricName:        "temperature",
				ConditionOperator: ">",
				ConditionValue:    100.0,
				DurationSeconds:   0,
				Severity:          "critical",
				IsEnabled:         true,
			},
		},
	}
	cache := NewMockCache()
	pub := &MockPublisher{}

	eval := NewEvaluator(db, cache, pub, logger)

	// Send message that triggers violation
	msg := &consumer.TelemetryMessage{
		FactoryID: "fact-1",
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     120.0,
		Timestamp: time.Now(),
	}

	err := eval.Evaluate(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Alarm was raised in DB
	if db.createdAlarm == nil {
		t.Fatal("expected alarm to be created, but got nil")
	}
	if db.createdAlarm.RuleID != "rule-1" || db.createdAlarm.TriggerValue != 120.0 {
		t.Errorf("unexpected alarm properties: %+v", db.createdAlarm)
	}

	// Verify active alarm is recorded in Cache
	activeID, _ := cache.GetActiveAlarmID(context.Background(), "rule-1", "mach-1")
	if activeID != "alarm-uuid-123" {
		t.Errorf("expected active alarm cache to be alarm-uuid-123, got %s", activeID)
	}

	// Verify Event was published
	if len(pub.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(pub.events))
	}
	if pub.events[0].EventType != "alarm.raised" || pub.events[0].Alarm.ID != "alarm-uuid-123" {
		t.Errorf("unexpected event: %+v", pub.events[0])
	}
}

func TestEvaluator_AutoResolve(t *testing.T) {
	logger := zap.NewNop()
	db := &MockDB{
		rules: []*model.AlarmRule{
			{
				ID:                "rule-1",
				MetricName:        "temperature",
				ConditionOperator: ">",
				ConditionValue:    100.0,
				IsEnabled:         true,
			},
		},
	}
	cache := NewMockCache()
	// Pre-fill cache with active alarm
	cache.SetActiveAlarmID(context.Background(), "rule-1", "mach-1", "alarm-uuid-123")

	pub := &MockPublisher{}

	eval := NewEvaluator(db, cache, pub, logger)

	// Send message that is normal
	msg := &consumer.TelemetryMessage{
		FactoryID: "fact-1",
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     85.0,
		Timestamp: time.Now(),
	}

	err := eval.Evaluate(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify DB ResolveAlarm was called
	if db.resolvedID != "alarm-uuid-123" {
		t.Errorf("expected resolved alarm ID to be alarm-uuid-123, got %s", db.resolvedID)
	}

	// Verify Cache was cleared
	activeID, _ := cache.GetActiveAlarmID(context.Background(), "rule-1", "mach-1")
	if activeID != "" {
		t.Errorf("expected active alarm cache to be cleared, got %s", activeID)
	}

	// Verify Event was published
	if len(pub.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(pub.events))
	}
	if pub.events[0].EventType != "alarm.resolved" {
		t.Errorf("unexpected event: %+v", pub.events[0])
	}
}

func TestEvaluator_DelayedAlarm(t *testing.T) {
	logger := zap.NewNop()
	db := &MockDB{
		rules: []*model.AlarmRule{
			{
				ID:                "rule-1",
				MetricName:        "temperature",
				ConditionOperator: ">",
				ConditionValue:    100.0,
				DurationSeconds:   5,
				Severity:          "critical",
				IsEnabled:         true,
			},
		},
	}
	cache := NewMockCache()
	pub := &MockPublisher{}

	eval := NewEvaluator(db, cache, pub, logger)

	now := time.Now()

	// 1. First violation (t=0)
	msg1 := &consumer.TelemetryMessage{
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     120.0,
		Timestamp: now,
	}
	if err := eval.Evaluate(context.Background(), msg1); err != nil {
		t.Fatal(err)
	}

	// Verify no alarm is raised yet
	if db.createdAlarm != nil {
		t.Fatal("alarm raised prematurely")
	}

	// Verify failure start time is set
	failTime, _ := cache.GetFailureStartTime(context.Background(), "rule-1", "mach-1")
	if failTime == nil || !failTime.Equal(now) {
		t.Fatalf("expected failure start time to be set to %v, got %v", now, failTime)
	}

	// 2. Second violation (t=3s)
	msg2 := &consumer.TelemetryMessage{
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     120.0,
		Timestamp: now.Add(3 * time.Second),
	}
	if err := eval.Evaluate(context.Background(), msg2); err != nil {
		t.Fatal(err)
	}

	// Verify no alarm is raised yet
	if db.createdAlarm != nil {
		t.Fatal("alarm raised prematurely at t=3s")
	}

	// 3. Third violation (t=6s >= 5s duration)
	msg3 := &consumer.TelemetryMessage{
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     120.0,
		Timestamp: now.Add(6 * time.Second),
	}
	if err := eval.Evaluate(context.Background(), msg3); err != nil {
		t.Fatal(err)
	}

	// Verify alarm is now raised!
	if db.createdAlarm == nil {
		t.Fatal("expected alarm to be raised at t=6s, but got nil")
	}

	// Verify failure start time in Cache was cleared
	failTime, _ = cache.GetFailureStartTime(context.Background(), "rule-1", "mach-1")
	if failTime != nil {
		t.Fatalf("expected failure start time cache to be cleared, got %v", failTime)
	}

	// Verify active alarm is in cache
	activeID, _ := cache.GetActiveAlarmID(context.Background(), "rule-1", "mach-1")
	if activeID != "alarm-uuid-123" {
		t.Fatalf("expected active alarm ID in cache, got %s", activeID)
	}
}

func TestEvaluator_NormalRecoveryResetsFailureStartTime(t *testing.T) {
	logger := zap.NewNop()
	db := &MockDB{
		rules: []*model.AlarmRule{
			{
				ID:                "rule-1",
				MetricName:        "temperature",
				ConditionOperator: ">",
				ConditionValue:    100.0,
				DurationSeconds:   5,
				IsEnabled:         true,
			},
		},
	}
	cache := NewMockCache()
	pub := &MockPublisher{}

	eval := NewEvaluator(db, cache, pub, logger)

	now := time.Now()

	// 1. Violation starts (t=0)
	msg1 := &consumer.TelemetryMessage{
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     120.0,
		Timestamp: now,
	}
	eval.Evaluate(context.Background(), msg1)

	failTime, _ := cache.GetFailureStartTime(context.Background(), "rule-1", "mach-1")
	if failTime == nil {
		t.Fatal("expected failure start time to be set")
	}

	// 2. Metric returns to normal (t=2s)
	msg2 := &consumer.TelemetryMessage{
		MachineID: "mach-1",
		Metric:    "temperature",
		Value:     80.0,
		Timestamp: now.Add(2 * time.Second),
	}
	eval.Evaluate(context.Background(), msg2)

	// Verify failure start time was cleared
	failTime, _ = cache.GetFailureStartTime(context.Background(), "rule-1", "mach-1")
	if failTime != nil {
		t.Fatal("expected failure start time to be cleared after recovery")
	}
}

func TestCheckCondition(t *testing.T) {
	tests := []struct {
		val      float64
		op       string
		limit    float64
		expected bool
	}{
		{10.0, ">", 5.0, true},
		{10.0, ">", 15.0, false},
		{10.0, ">=", 10.0, true},
		{10.0, "<", 15.0, true},
		{10.0, "<=", 9.0, false},
		{10.0, "==", 10.0, true},
		{10.0, "!=", 5.0, true},
		{10.0, "invalid", 5.0, false},
	}

	for _, tt := range tests {
		res := checkCondition(tt.val, tt.op, tt.limit)
		if res != tt.expected {
			t.Errorf("checkCondition(%f, %s, %f) = %t, expected %t", tt.val, tt.op, tt.limit, res, tt.expected)
		}
	}
}
