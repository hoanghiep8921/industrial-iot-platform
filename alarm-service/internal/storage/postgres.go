package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/industrial-iot/alarm-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepo struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresRepo(connStr string, logger *zap.Logger) (*PostgresRepo, error) {
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	poolCfg.MaxConns = 25

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL", zap.Int32("max_connections", poolCfg.MaxConns))

	return &PostgresRepo{pool: pool, logger: logger}, nil
}

func (r *PostgresRepo) Close() {
	r.logger.Info("Closing PostgreSQL connection pool...")
	r.pool.Close()
}

func (r *PostgresRepo) Pool() *pgxpool.Pool {
	return r.pool
}

// ============================================================
// Alarm Rules CRUD
// ============================================================

func (r *PostgresRepo) CreateRule(ctx context.Context, req *model.CreateRuleRequest) (*model.AlarmRule, error) {
	ruleType := req.Type
	if ruleType == "" {
		ruleType = "threshold"
	}
	severity := req.Severity
	if severity == "" {
		severity = "warning"
	}

	var rule model.AlarmRule
	err := r.pool.QueryRow(ctx, `
		INSERT INTO alarm_rules (name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE)
		RETURNING id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
	`,
		req.Name, req.Description, ruleType, req.MetricName, req.ConditionOperator, req.ConditionValue, req.DurationSeconds, severity, req.FactoryID,
	).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.MetricName, &rule.ConditionOperator, &rule.ConditionValue, &rule.DurationSeconds, &rule.Severity, &rule.FactoryID, &rule.IsEnabled, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create rule: %w", err)
	}

	r.logger.Info("Alarm rule created", zap.String("id", rule.ID), zap.String("name", rule.Name))
	return &rule, nil
}

func (r *PostgresRepo) GetRuleByID(ctx context.Context, id string) (*model.AlarmRule, error) {
	var rule model.AlarmRule
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
		FROM alarm_rules WHERE id = $1
	`, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.MetricName, &rule.ConditionOperator, &rule.ConditionValue, &rule.DurationSeconds, &rule.Severity, &rule.FactoryID, &rule.IsEnabled, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get rule by id: %w", err)
	}
	return &rule, nil
}

func (r *PostgresRepo) ListRules(ctx context.Context, factoryID string) ([]*model.AlarmRule, error) {
	var rows pgx.Rows
	var err error

	if factoryID != "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
			FROM alarm_rules WHERE factory_id = $1
			ORDER BY created_at DESC
		`, factoryID)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
			FROM alarm_rules
			ORDER BY created_at DESC
		`)
	}
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	defer rows.Close()

	var rules []*model.AlarmRule
	for rows.Next() {
		var rule model.AlarmRule
		err = rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.MetricName, &rule.ConditionOperator, &rule.ConditionValue, &rule.DurationSeconds, &rule.Severity, &rule.FactoryID, &rule.IsEnabled, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rule: %w", err)
		}
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

func (r *PostgresRepo) GetRulesByMetric(ctx context.Context, metricName string) ([]*model.AlarmRule, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
		FROM alarm_rules WHERE metric_name = $1 AND is_enabled = TRUE
	`, metricName)
	if err != nil {
		return nil, fmt.Errorf("get rules by metric: %w", err)
	}
	defer rows.Close()

	var rules []*model.AlarmRule
	for rows.Next() {
		var rule model.AlarmRule
		err = rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.MetricName, &rule.ConditionOperator, &rule.ConditionValue, &rule.DurationSeconds, &rule.Severity, &rule.FactoryID, &rule.IsEnabled, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rule: %w", err)
		}
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

func (r *PostgresRepo) UpdateRule(ctx context.Context, id string, req *model.UpdateRuleRequest) (*model.AlarmRule, error) {
	current, err := r.GetRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("rule not found")
	}

	name := current.Name
	if req.Name != nil {
		name = *req.Name
	}
	desc := current.Description
	if req.Description != nil {
		desc = *req.Description
	}
	ruleType := current.Type
	if req.Type != nil {
		ruleType = *req.Type
	}
	metricName := current.MetricName
	if req.MetricName != nil {
		metricName = *req.MetricName
	}
	op := current.ConditionOperator
	if req.ConditionOperator != nil {
		op = *req.ConditionOperator
	}
	val := current.ConditionValue
	if req.ConditionValue != nil {
		val = *req.ConditionValue
	}
	dur := current.DurationSeconds
	if req.DurationSeconds != nil {
		dur = *req.DurationSeconds
	}
	sev := current.Severity
	if req.Severity != nil {
		sev = *req.Severity
	}
	fac := current.FactoryID
	if req.FactoryID != nil {
		fac = *req.FactoryID
	}
	enabled := current.IsEnabled
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}

	var rule model.AlarmRule
	err = r.pool.QueryRow(ctx, `
		UPDATE alarm_rules
		SET name = $1, description = $2, type = $3, metric_name = $4, condition_operator = $5, condition_value = $6, duration_seconds = $7, severity = $8, factory_id = $9, is_enabled = $10, updated_at = NOW()
		WHERE id = $11
		RETURNING id, name, description, type, metric_name, condition_operator, condition_value, duration_seconds, severity, factory_id, is_enabled, created_at, updated_at
	`, name, desc, ruleType, metricName, op, val, dur, sev, fac, enabled, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.MetricName, &rule.ConditionOperator, &rule.ConditionValue, &rule.DurationSeconds, &rule.Severity, &rule.FactoryID, &rule.IsEnabled, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update rule: %w", err)
	}

	r.logger.Info("Alarm rule updated", zap.String("id", id))
	return &rule, nil
}

func (r *PostgresRepo) DeleteRule(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM alarm_rules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("rule not found")
	}
	r.logger.Info("Alarm rule deleted", zap.String("id", id))
	return nil
}

// ============================================================
// Alarms Management
// ============================================================

func (r *PostgresRepo) GetAlarmByID(ctx context.Context, id string) (*model.Alarm, error) {
	var a model.Alarm
	err := r.pool.QueryRow(ctx, `
		SELECT a.id, a.rule_id, r.name as rule_name, a.machine_id, a.metric_name, a.trigger_value, a.trigger_time, a.ack_time, COALESCE(a.ack_by, '') as ack_by, a.resolve_time, a.status, a.severity, a.created_at, a.updated_at
		FROM alarms a
		JOIN alarm_rules r ON a.rule_id = r.id
		WHERE a.id = $1
	`, id).Scan(
		&a.ID, &a.RuleID, &a.RuleName, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get alarm by id: %w", err)
	}
	return &a, nil
}

func (r *PostgresRepo) CreateAlarm(ctx context.Context, ruleID string, machineID string, metricName string, triggerVal float64, triggerTime time.Time, severity string) (*model.Alarm, error) {
	var a model.Alarm
	err := r.pool.QueryRow(ctx, `
		INSERT INTO alarms (rule_id, machine_id, metric_name, trigger_value, trigger_time, status, severity)
		VALUES ($1, $2, $3, $4, $5, 'RAISED', $6)
		RETURNING id, rule_id, machine_id, metric_name, trigger_value, trigger_time, ack_time, COALESCE(ack_by, '') as ack_by, resolve_time, status, severity, created_at, updated_at
	`, ruleID, machineID, metricName, triggerVal, triggerTime, severity).Scan(
		&a.ID, &a.RuleID, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create alarm: %w", err)
	}

	// Fetch rule name for populated model
	rule, _ := r.GetRuleByID(ctx, ruleID)
	if rule != nil {
		a.RuleName = rule.Name
	}

	r.logger.Info("Alarm created", zap.String("id", a.ID), zap.String("rule_id", ruleID), zap.String("machine_id", machineID), zap.Float64("val", triggerVal))
	return &a, nil
}

func (r *PostgresRepo) AcknowledgeAlarm(ctx context.Context, id string, ackBy string) (*model.Alarm, error) {
	now := time.Now()
	var a model.Alarm
	err := r.pool.QueryRow(ctx, `
		UPDATE alarms
		SET status = 'ACKNOWLEDGED', ack_time = $1, ack_by = $2, updated_at = NOW()
		WHERE id = $3 AND status = 'RAISED'
		RETURNING id, rule_id, machine_id, metric_name, trigger_value, trigger_time, ack_time, COALESCE(ack_by, '') as ack_by, resolve_time, status, severity, created_at, updated_at
	`, now, ackBy, id).Scan(
		&a.ID, &a.RuleID, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("alarm not found or already acknowledged/resolved")
		}
		return nil, fmt.Errorf("acknowledge alarm: %w", err)
	}

	rule, _ := r.GetRuleByID(ctx, a.RuleID)
	if rule != nil {
		a.RuleName = rule.Name
	}

	r.logger.Info("Alarm acknowledged", zap.String("id", id), zap.String("by", ackBy))
	return &a, nil
}

func (r *PostgresRepo) ResolveAlarm(ctx context.Context, id string, resolveTime time.Time) (*model.Alarm, error) {
	var a model.Alarm
	err := r.pool.QueryRow(ctx, `
		UPDATE alarms
		SET status = 'RESOLVED', resolve_time = $1, updated_at = NOW()
		WHERE id = $2 AND status IN ('RAISED', 'ACKNOWLEDGED')
		RETURNING id, rule_id, machine_id, metric_name, trigger_value, trigger_time, ack_time, COALESCE(ack_by, '') as ack_by, resolve_time, status, severity, created_at, updated_at
	`, resolveTime, id).Scan(
		&a.ID, &a.RuleID, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("alarm not found or already resolved")
		}
		return nil, fmt.Errorf("resolve alarm: %w", err)
	}

	rule, _ := r.GetRuleByID(ctx, a.RuleID)
	if rule != nil {
		a.RuleName = rule.Name
	}

	r.logger.Info("Alarm resolved", zap.String("id", id))
	return &a, nil
}

func (r *PostgresRepo) ListActiveAlarms(ctx context.Context) ([]*model.Alarm, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.id, a.rule_id, r.name as rule_name, a.machine_id, a.metric_name, a.trigger_value, a.trigger_time, a.ack_time, COALESCE(a.ack_by, '') as ack_by, a.resolve_time, a.status, a.severity, a.created_at, a.updated_at
		FROM alarms a
		JOIN alarm_rules r ON a.rule_id = r.id
		WHERE a.status IN ('RAISED', 'ACKNOWLEDGED')
		ORDER BY a.trigger_time DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list active alarms: %w", err)
	}
	defer rows.Close()

	var alarms []*model.Alarm
	for rows.Next() {
		var a model.Alarm
		err = rows.Scan(
			&a.ID, &a.RuleID, &a.RuleName, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan alarm: %w", err)
		}
		alarms = append(alarms, &a)
	}
	return alarms, rows.Err()
}

func (r *PostgresRepo) ListAlarmHistory(ctx context.Context, filter model.AlarmFilter) ([]*model.Alarm, int, error) {
	var conditions []string
	var args []interface{}
	idx := 1

	if filter.FactoryID != "" {
		conditions = append(conditions, fmt.Sprintf("r.factory_id = $%d", idx))
		args = append(args, filter.FactoryID)
		idx++
	}
	if filter.MachineID != "" {
		conditions = append(conditions, fmt.Sprintf("a.machine_id = $%d", idx))
		args = append(args, filter.MachineID)
		idx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("a.status = $%d", idx))
		args = append(args, filter.Status)
		idx++
	}
	if filter.Severity != "" {
		conditions = append(conditions, fmt.Sprintf("a.severity = $%d", idx))
		args = append(args, filter.Severity)
		idx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM alarms a
		JOIN alarm_rules r ON a.rule_id = r.id
		%s
	`, whereClause)

	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count alarm history: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT a.id, a.rule_id, r.name as rule_name, a.machine_id, a.metric_name, a.trigger_value, a.trigger_time, a.ack_time, COALESCE(a.ack_by, '') as ack_by, a.resolve_time, a.status, a.severity, a.created_at, a.updated_at
		FROM alarms a
		JOIN alarm_rules r ON a.rule_id = r.id
		%s
		ORDER BY a.trigger_time DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query alarm history: %w", err)
	}
	defer rows.Close()

	var alarms []*model.Alarm
	for rows.Next() {
		var a model.Alarm
		err = rows.Scan(
			&a.ID, &a.RuleID, &a.RuleName, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan alarm history row: %w", err)
		}
		alarms = append(alarms, &a)
	}

	return alarms, total, rows.Err()
}

func (r *PostgresRepo) GetActiveAlarm(ctx context.Context, ruleID string, machineID string) (*model.Alarm, error) {
	var a model.Alarm
	err := r.pool.QueryRow(ctx, `
		SELECT a.id, a.rule_id, r.name as rule_name, a.machine_id, a.metric_name, a.trigger_value, a.trigger_time, a.ack_time, COALESCE(a.ack_by, '') as ack_by, a.resolve_time, a.status, a.severity, a.created_at, a.updated_at
		FROM alarms a
		JOIN alarm_rules r ON a.rule_id = r.id
		WHERE a.rule_id = $1 AND a.machine_id = $2 AND a.status IN ('RAISED', 'ACKNOWLEDGED')
		LIMIT 1
	`, ruleID, machineID).Scan(
		&a.ID, &a.RuleID, &a.RuleName, &a.MachineID, &a.MetricName, &a.TriggerValue, &a.TriggerTime, &a.AckTime, &a.AckBy, &a.ResolveTime, &a.Status, &a.Severity, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get active alarm: %w", err)
	}
	return &a, nil
}

