package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// PostgresRepo handles all PostgreSQL operations for the device service
type PostgresRepo struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresRepo creates a new PostgreSQL repository
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

	logger.Info("Connected to PostgreSQL",
		zap.Int32("max_connections", poolCfg.MaxConns),
	)

	return &PostgresRepo{pool: pool, logger: logger}, nil
}

// Close closes the database connection pool
func (r *PostgresRepo) Close() {
	r.logger.Info("Closing PostgreSQL connection pool...")
	r.pool.Close()
}

// Pool returns the underlying connection pool (for migrations)
func (r *PostgresRepo) Pool() *pgxpool.Pool {
	return r.pool
}

// ============================================================
// Device CRUD
// ============================================================

// CreateDevice inserts a new device and its shadow + connection records
func (r *PostgresRepo) CreateDevice(ctx context.Context, req *model.CreateDeviceRequest) (*model.Device, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Default values
	protocol := req.Protocol
	if protocol == "" {
		protocol = "mqtt"
	}
	caps := req.Capabilities
	if caps == nil {
		caps = []string{}
	}
	meta := req.Metadata
	if meta == nil {
		meta = map[string]any{}
	}

	capsJSON, _ := json.Marshal(caps)
	metaJSON, _ := json.Marshal(meta)

	var device model.Device
	err = tx.QueryRow(ctx, `
		INSERT INTO devices (serial_number, name, model, vendor, factory_id, area, line, protocol, capabilities, firmware_version, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, serial_number, name, model, vendor, factory_id, area, line, protocol, capabilities,
		          status, firmware_version, metadata, last_seen_at, created_at, updated_at
	`,
		req.SerialNumber, req.Name, req.Model, req.Vendor, req.FactoryID, req.Area, req.Line,
		protocol, capsJSON, req.FirmwareVersion, metaJSON,
	).Scan(
		&device.ID, &device.SerialNumber, &device.Name, &device.Model, &device.Vendor,
		&device.FactoryID, &device.Area, &device.Line, &device.Protocol, &capsJSON,
		&device.Status, &device.FirmwareVersion, &metaJSON, &device.LastSeenAt,
		&device.CreatedAt, &device.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}

	json.Unmarshal(capsJSON, &device.Capabilities)
	json.Unmarshal(metaJSON, &device.Metadata)

	// Create shadow record
	_, err = tx.Exec(ctx, `INSERT INTO device_shadows (device_id) VALUES ($1)`, device.ID)
	if err != nil {
		return nil, fmt.Errorf("insert shadow: %w", err)
	}

	// Create connection record
	_, err = tx.Exec(ctx, `INSERT INTO device_connections (device_id) VALUES ($1)`, device.ID)
	if err != nil {
		return nil, fmt.Errorf("insert connection: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	r.logger.Info("Device created",
		zap.String("id", device.ID),
		zap.String("serial_number", device.SerialNumber),
		zap.String("factory_id", device.FactoryID),
	)

	return &device, nil
}

// GetDeviceByID retrieves a device by its UUID
func (r *PostgresRepo) GetDeviceByID(ctx context.Context, id string) (*model.Device, error) {
	var device model.Device
	var capsJSON, metaJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, serial_number, name, model, vendor, factory_id, area, line, protocol,
		       capabilities, status, firmware_version, metadata, last_seen_at, created_at, updated_at
		FROM devices WHERE id = $1
	`, id).Scan(
		&device.ID, &device.SerialNumber, &device.Name, &device.Model, &device.Vendor,
		&device.FactoryID, &device.Area, &device.Line, &device.Protocol,
		&capsJSON, &device.Status, &device.FirmwareVersion, &metaJSON, &device.LastSeenAt,
		&device.CreatedAt, &device.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get device: %w", err)
	}

	json.Unmarshal(capsJSON, &device.Capabilities)
	json.Unmarshal(metaJSON, &device.Metadata)

	return &device, nil
}

// GetDeviceBySerialNumber retrieves a device by serial number
func (r *PostgresRepo) GetDeviceBySerialNumber(ctx context.Context, sn string) (*model.Device, error) {
	var device model.Device
	var capsJSON, metaJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, serial_number, name, model, vendor, factory_id, area, line, protocol,
		       capabilities, status, firmware_version, metadata, last_seen_at, created_at, updated_at
		FROM devices WHERE serial_number = $1
	`, sn).Scan(
		&device.ID, &device.SerialNumber, &device.Name, &device.Model, &device.Vendor,
		&device.FactoryID, &device.Area, &device.Line, &device.Protocol,
		&capsJSON, &device.Status, &device.FirmwareVersion, &metaJSON, &device.LastSeenAt,
		&device.CreatedAt, &device.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get device by SN: %w", err)
	}

	json.Unmarshal(capsJSON, &device.Capabilities)
	json.Unmarshal(metaJSON, &device.Metadata)

	return &device, nil
}

// ListDevices returns devices matching the filter with pagination
func (r *PostgresRepo) ListDevices(ctx context.Context, filter model.DeviceFilter) ([]*model.Device, int, error) {
	// Build dynamic WHERE clause
	where := []string{}
	args := []interface{}{}
	argIdx := 1

	if filter.FactoryID != "" {
		where = append(where, fmt.Sprintf("factory_id = $%d", argIdx))
		args = append(args, filter.FactoryID)
		argIdx++
	}
	if filter.Area != "" {
		where = append(where, fmt.Sprintf("area = $%d", argIdx))
		args = append(args, filter.Area)
		argIdx++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Protocol != "" {
		where = append(where, fmt.Sprintf("protocol = $%d", argIdx))
		args = append(args, filter.Protocol)
		argIdx++
	}
	if filter.Search != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR serial_number ILIKE $%d)", argIdx, argIdx+1))
		like := "%" + filter.Search + "%"
		args = append(args, like, like)
		argIdx += 2
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM devices %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count devices: %w", err)
	}

	// Pagination defaults
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Query
	query := fmt.Sprintf(`
		SELECT id, serial_number, name, model, vendor, factory_id, area, line, protocol,
		       capabilities, status, firmware_version, metadata, last_seen_at, created_at, updated_at
		FROM devices %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()

	var devices []*model.Device
	for rows.Next() {
		var d model.Device
		var capsJSON, metaJSON []byte
		if err := rows.Scan(
			&d.ID, &d.SerialNumber, &d.Name, &d.Model, &d.Vendor,
			&d.FactoryID, &d.Area, &d.Line, &d.Protocol,
			&capsJSON, &d.Status, &d.FirmwareVersion, &metaJSON, &d.LastSeenAt,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan device: %w", err)
		}
		json.Unmarshal(capsJSON, &d.Capabilities)
		json.Unmarshal(metaJSON, &d.Metadata)
		devices = append(devices, &d)
	}

	return devices, total, rows.Err()
}

// UpdateDevice updates a device's metadata
func (r *PostgresRepo) UpdateDevice(ctx context.Context, id string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	// Build dynamic SET clause
	sets := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Model != nil {
		sets = append(sets, fmt.Sprintf("model = $%d", argIdx))
		args = append(args, *req.Model)
		argIdx++
	}
	if req.Vendor != nil {
		sets = append(sets, fmt.Sprintf("vendor = $%d", argIdx))
		args = append(args, *req.Vendor)
		argIdx++
	}
	if req.FactoryID != nil {
		sets = append(sets, fmt.Sprintf("factory_id = $%d", argIdx))
		args = append(args, *req.FactoryID)
		argIdx++
	}
	if req.Area != nil {
		sets = append(sets, fmt.Sprintf("area = $%d", argIdx))
		args = append(args, *req.Area)
		argIdx++
	}
	if req.Line != nil {
		sets = append(sets, fmt.Sprintf("line = $%d", argIdx))
		args = append(args, *req.Line)
		argIdx++
	}
	if req.Protocol != nil {
		sets = append(sets, fmt.Sprintf("protocol = $%d", argIdx))
		args = append(args, *req.Protocol)
		argIdx++
	}
	if req.Capabilities != nil {
		sets = append(sets, fmt.Sprintf("capabilities = $%d", argIdx))
		capsJSON, _ := json.Marshal(*req.Capabilities)
		args = append(args, capsJSON)
		argIdx++
	}
	if req.FirmwareVersion != nil {
		sets = append(sets, fmt.Sprintf("firmware_version = $%d", argIdx))
		args = append(args, *req.FirmwareVersion)
		argIdx++
	}
	if req.Metadata != nil {
		sets = append(sets, fmt.Sprintf("metadata = $%d", argIdx))
		metaJSON, _ := json.Marshal(*req.Metadata)
		args = append(args, metaJSON)
		argIdx++
	}

	if len(sets) == 0 {
		return r.GetDeviceByID(ctx, id)
	}

	sets = append(sets, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE devices SET %s WHERE id = $%d", strings.Join(sets, ", "), argIdx)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update device: %w", err)
	}

	return r.GetDeviceByID(ctx, id)
}

// DeleteDevice removes a device (shadows + connections cascade)
func (r *PostgresRepo) DeleteDevice(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM devices WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete device: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("device not found")
	}
	r.logger.Info("Device deleted", zap.String("id", id))
	return nil
}

// UpdateDeviceStatus changes the device status and updates last_seen_at
func (r *PostgresRepo) UpdateDeviceStatus(ctx context.Context, id string, status model.DeviceStatus) error {
	now := time.Now()
	tag, err := r.pool.Exec(ctx, `
		UPDATE devices SET status = $1, last_seen_at = $2, updated_at = NOW()
		WHERE id = $3
	`, status, now, id)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("device not found")
	}

	// Also update connection tracking
	connected := status == model.StatusOnline
	if connected {
		r.pool.Exec(ctx, `
			INSERT INTO device_connections (device_id, connected, connected_at, last_heartbeat)
			VALUES ($1, TRUE, $2, $2)
			ON CONFLICT (device_id) DO UPDATE SET
				connected = TRUE, connected_at = $2, last_heartbeat = $2
		`, id, now)
	} else {
		r.pool.Exec(ctx, `
			UPDATE device_connections SET
				connected = FALSE, disconnected_at = $2 WHERE device_id = $1
		`, id, now)
	}

	return nil
}

// GetStats returns device statistics
func (r *PostgresRepo) GetStats(ctx context.Context) (*model.DeviceStats, error) {
	stats := &model.DeviceStats{
		ByFactory: make(map[string]int),
		ByStatus:  make(map[model.DeviceStatus]int),
		ByModel:   make(map[string]int),
	}

	// Total
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM devices`).Scan(&stats.Total)

	// By factory
	rows, _ := r.pool.Query(ctx, `SELECT factory_id, COUNT(*) FROM devices GROUP BY factory_id`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var k string
			var v int
			rows.Scan(&k, &v)
			stats.ByFactory[k] = v
		}
	}

	// By status
	rows2, _ := r.pool.Query(ctx, `SELECT status, COUNT(*) FROM devices GROUP BY status`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var k string
			var v int
			rows2.Scan(&k, &v)
			stats.ByStatus[model.DeviceStatus(k)] = v
		}
	}

	// By model
	rows3, _ := r.pool.Query(ctx, `SELECT COALESCE(model, 'unknown'), COUNT(*) FROM devices GROUP BY model`)
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var k string
			var v int
			rows3.Scan(&k, &v)
			stats.ByModel[k] = v
		}
	}

	return stats, nil
}

// ============================================================
// Device Shadow
// ============================================================

// GetShadow retrieves the device shadow
func (r *PostgresRepo) GetShadow(ctx context.Context, deviceID string) (*model.DeviceShadow, error) {
	var s model.DeviceShadow
	var desiredJSON, reportedJSON, deltaJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT device_id, desired_state, reported_state, delta, version, updated_at
		FROM device_shadows WHERE device_id = $1
	`, deviceID).Scan(&s.DeviceID, &desiredJSON, &reportedJSON, &deltaJSON, &s.Version, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get shadow: %w", err)
	}

	json.Unmarshal(desiredJSON, &s.DesiredState)
	json.Unmarshal(reportedJSON, &s.ReportedState)
	json.Unmarshal(deltaJSON, &s.Delta)

	return &s, nil
}

// UpsertDesiredState updates the desired state and recalculates delta
func (r *PostgresRepo) UpsertDesiredState(ctx context.Context, deviceID string, desired map[string]any) (*model.DeviceShadow, error) {
	// Fetch current reported state
	shadow, err := r.GetShadow(ctx, deviceID)
	if err != nil || shadow == nil {
		return nil, fmt.Errorf("shadow not found for device %s", deviceID)
	}

	desiredJSON, _ := json.Marshal(desired)
	reportedJSON, _ := json.Marshal(shadow.ReportedState)
	delta := computeDelta(desired, shadow.ReportedState)
	deltaJSON, _ := json.Marshal(delta)

	var newVersion int64
	err = r.pool.QueryRow(ctx, `
		UPDATE device_shadows
		SET desired_state = $1, reported_state = $2, delta = $3, version = version + 1, updated_at = NOW()
		WHERE device_id = $4
		RETURNING version
	`, desiredJSON, reportedJSON, deltaJSON, deviceID).Scan(&newVersion)
	if err != nil {
		return nil, fmt.Errorf("upsert desired: %w", err)
	}

	return &model.DeviceShadow{
		DeviceID:      deviceID,
		DesiredState:  desired,
		ReportedState: shadow.ReportedState,
		Delta:         delta,
		Version:       newVersion,
		UpdatedAt:     time.Now(),
	}, nil
}

// UpsertReportedState updates the reported state and recalculates delta
func (r *PostgresRepo) UpsertReportedState(ctx context.Context, deviceID string, reported map[string]any) (*model.DeviceShadow, error) {
	shadow, err := r.GetShadow(ctx, deviceID)
	if err != nil || shadow == nil {
		return nil, fmt.Errorf("shadow not found for device %s", deviceID)
	}

	reportedJSON, _ := json.Marshal(reported)
	desiredJSON, _ := json.Marshal(shadow.DesiredState)
	delta := computeDelta(shadow.DesiredState, reported)
	deltaJSON, _ := json.Marshal(delta)

	var newVersion int64
	err = r.pool.QueryRow(ctx, `
		UPDATE device_shadows
		SET reported_state = $1, desired_state = $2, delta = $3, version = version + 1, updated_at = NOW()
		WHERE device_id = $4
		RETURNING version
	`, reportedJSON, desiredJSON, deltaJSON, deviceID).Scan(&newVersion)
	if err != nil {
		return nil, fmt.Errorf("upsert reported: %w", err)
	}

	// Update device heartbeat
	r.pool.Exec(ctx, `
		UPDATE devices SET last_seen_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, deviceID)

	return &model.DeviceShadow{
		DeviceID:      deviceID,
		DesiredState:  shadow.DesiredState,
		ReportedState: reported,
		Delta:         delta,
		Version:       newVersion,
		UpdatedAt:     time.Now(),
	}, nil
}

// computeDelta finds keys in desired that differ from reported
func computeDelta(desired, reported map[string]any) map[string]any {
	delta := map[string]any{}
	for key, desiredVal := range desired {
		reportedVal, exists := reported[key]
		if !exists || !valueEqual(desiredVal, reportedVal) {
			delta[key] = desiredVal
		}
	}
	return delta
}

// valueEqual compares two values, normalizing numeric types (int vs float64 from JSON)
func valueEqual(a, b any) bool {
	// Normalize numeric types
	aNum, aOk := toFloat64(a)
	bNum, bOk := toFloat64(b)
	if aOk && bOk {
		return aNum == bNum
	}
	// Fall back to string comparison for non-numeric types
	return fmt.Sprint(a) == fmt.Sprint(b)
}

// toFloat64 attempts to convert a value to float64
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}
