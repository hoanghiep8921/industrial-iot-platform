package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// TelemetryRecord represents a data point to store
type TelemetryRecord struct {
	Time       time.Time
	DeviceID   string
	FactoryID  string
	Area       string
	MetricName string
	Value      float64
	Unit       string
	Quality    int
	Tags       map[string]string
	RawPayload json.RawMessage
}

// Config holds TimescaleDB connection configuration
type Config struct {
	Host           string
	Port           int
	Database       string
	Username       string
	Password       string
	SSLMode        string
	MaxConnections int
}

// TimescaleDB wraps the database connection pool
type TimescaleDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewTimescaleDB creates a new TimescaleDB connection
func NewTimescaleDB(cfg Config, logger *zap.Logger) (*TimescaleDB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		cfg.Database, cfg.SSLMode, cfg.MaxConnections,
	)

	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxConnections)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to TimescaleDB",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return &TimescaleDB{pool: pool, logger: logger}, nil
}

// Insert stores a telemetry record into TimescaleDB
func (db *TimescaleDB) Insert(ctx context.Context, record *TelemetryRecord) error {
	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		tagsJSON = []byte("{}")
	}

	rawPayloadJSON := record.RawPayload
	if len(rawPayloadJSON) == 0 {
		rawPayloadJSON = json.RawMessage("{}")
	}

	// Insert into hypertable
	_, err = db.pool.Exec(ctx, `
		INSERT INTO telemetry (time, device_id, factory_id, area, metric_name, value, unit, quality, tags, raw_payload)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, record.Time, record.DeviceID, record.FactoryID, record.Area,
		record.MetricName, record.Value, record.Unit, record.Quality,
		tagsJSON, rawPayloadJSON)

	if err != nil {
		return fmt.Errorf("insert telemetry: %w", err)
	}

	// Upsert latest value
	_, err = db.pool.Exec(ctx, `
		INSERT INTO device_latest_value (device_id, metric_name, factory_id, area, value, unit, quality, tags, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (device_id, metric_name) DO UPDATE SET
			factory_id = EXCLUDED.factory_id,
			area = EXCLUDED.area,
			value = EXCLUDED.value,
			unit = EXCLUDED.unit,
			quality = EXCLUDED.quality,
			tags = EXCLUDED.tags,
			updated_at = NOW()
	`, record.DeviceID, record.MetricName, record.FactoryID, record.Area,
		record.Value, record.Unit, record.Quality, tagsJSON)

	if err != nil {
		db.logger.Error("Failed to upsert latest value", zap.Error(err))
		// Don't return error — the main insert already succeeded
	}

	return nil
}

// InsertBatch stores multiple records in a single transaction
func (db *TimescaleDB) InsertBatch(ctx context.Context, records []*TelemetryRecord) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, record := range records {
		tagsJSON, _ := json.Marshal(record.Tags)
		rawPayloadJSON := record.RawPayload
		if len(rawPayloadJSON) == 0 {
			rawPayloadJSON = json.RawMessage("{}")
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO telemetry (time, device_id, factory_id, area, metric_name, value, unit, quality, tags, raw_payload)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, record.Time, record.DeviceID, record.FactoryID, record.Area,
			record.MetricName, record.Value, record.Unit, record.Quality,
			tagsJSON, rawPayloadJSON)

		if err != nil {
			return fmt.Errorf("insert telemetry batch: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// QueryRecent returns the most recent telemetry records
func (db *TimescaleDB) QueryRecent(ctx context.Context, limit int) ([]*TelemetryRecord, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT time, device_id, factory_id, area, metric_name, value, unit, quality, tags
		FROM telemetry
		ORDER BY time DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent: %w", err)
	}
	defer rows.Close()

	var records []*TelemetryRecord
	for rows.Next() {
		var r TelemetryRecord
		var tagsJSON []byte
		if err := rows.Scan(&r.Time, &r.DeviceID, &r.FactoryID, &r.Area,
			&r.MetricName, &r.Value, &r.Unit, &r.Quality, &tagsJSON); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		json.Unmarshal(tagsJSON, &r.Tags)
		records = append(records, &r)
	}

	return records, rows.Err()
}

// QueryByDevice returns telemetry for a specific device in a time range
func (db *TimescaleDB) QueryByDevice(ctx context.Context, deviceID string, from, to time.Time) ([]*TelemetryRecord, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT time, device_id, factory_id, area, metric_name, value, unit, quality, tags
		FROM telemetry
		WHERE device_id = $1 AND time >= $2 AND time <= $3
		ORDER BY time DESC
		LIMIT 1000
	`, deviceID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query by device: %w", err)
	}
	defer rows.Close()

	var records []*TelemetryRecord
	for rows.Next() {
		var r TelemetryRecord
		var tagsJSON []byte
		if err := rows.Scan(&r.Time, &r.DeviceID, &r.FactoryID, &r.Area,
			&r.MetricName, &r.Value, &r.Unit, &r.Quality, &tagsJSON); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		json.Unmarshal(tagsJSON, &r.Tags)
		records = append(records, &r)
	}

	return records, rows.Err()
}

// GetLatestValues returns latest values for all devices
func (db *TimescaleDB) GetLatestValues(ctx context.Context, factoryID string) ([]*TelemetryRecord, error) {
	query := `
		SELECT device_id, metric_name, factory_id, area, value, unit, COALESCE(quality, 0), tags, updated_at
		FROM device_latest_value
	`
	args := []interface{}{}
	if factoryID != "" {
		query += ` WHERE factory_id = $1`
		args = append(args, factoryID)
	}
	query += ` ORDER BY factory_id, area, device_id`

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query latest values: %w", err)
	}
	defer rows.Close()

	var records []*TelemetryRecord
	for rows.Next() {
		var r TelemetryRecord
		var tagsJSON []byte
		if err := rows.Scan(&r.DeviceID, &r.MetricName, &r.FactoryID, &r.Area,
			&r.Value, &r.Unit, &r.Quality, &tagsJSON, &r.Time); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		json.Unmarshal(tagsJSON, &r.Tags)
		records = append(records, &r)
	}

	return records, rows.Err()
}

// GetStorageStats returns approximate hypertable storage size
func (db *TimescaleDB) GetStorageStats(ctx context.Context) (map[string]interface{}, error) {
	var totalSize int64
	err := db.pool.QueryRow(ctx, `
		SELECT COALESCE(
			SUM(COALESCE(total_size, 0))
		, 0)
		FROM timescaledb_information.hypertable_size_stats
	`).Scan(&totalSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_size_bytes": totalSize,
		"total_size_mb":    float64(totalSize) / 1024 / 1024,
		"total_size_gb":    float64(totalSize) / 1024 / 1024 / 1024,
	}, nil
}

// Close closes the database connection pool
func (db *TimescaleDB) Close() {
	db.logger.Info("Closing TimescaleDB connection pool...")
	db.pool.Close()
}
