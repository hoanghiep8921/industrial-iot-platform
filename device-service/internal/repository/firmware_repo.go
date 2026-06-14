package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// FirmwareRepo handles firmware and firmware_update database operations
type FirmwareRepo struct {
	repo   *PostgresRepo
	logger *zap.Logger
}

// NewFirmwareRepo creates a new firmware repository
func NewFirmwareRepo(repo *PostgresRepo) *FirmwareRepo {
	return &FirmwareRepo{repo: repo, logger: repo.logger}
}

// ============================================================
// Firmware CRUD
// ============================================================

// CreateFirmware inserts a new firmware metadata record
func (r *FirmwareRepo) CreateFirmware(ctx context.Context, req *model.CreateFirmwareRequest) (*model.Firmware, error) {
	var fw model.Firmware
	err := r.repo.Pool().QueryRow(ctx, `
		INSERT INTO firmwares (version, model, description, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, model, description, file_name, file_size, checksum_sha256,
		          minio_path, status, created_by, created_at, released_at
	`,
		req.Version, req.Model, req.Description, req.CreatedBy,
	).Scan(
		&fw.ID, &fw.Version, &fw.Model, &fw.Description, &fw.FileName, &fw.FileSize,
		&fw.ChecksumSHA256, &fw.MinioPath, &fw.Status, &fw.CreatedBy, &fw.CreatedAt, &fw.ReleasedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert firmware: %w", err)
	}

	r.logger.Info("Firmware created",
		zap.String("id", fw.ID),
		zap.String("version", fw.Version),
	)

	return &fw, nil
}

// GetFirmwareByID retrieves a firmware by ID
func (r *FirmwareRepo) GetFirmwareByID(ctx context.Context, id string) (*model.Firmware, error) {
	var fw model.Firmware
	err := r.repo.Pool().QueryRow(ctx, `
		SELECT id, version, model, description, file_name, file_size, checksum_sha256,
		       minio_path, status, created_by, created_at, released_at
		FROM firmwares WHERE id = $1
	`, id).Scan(
		&fw.ID, &fw.Version, &fw.Model, &fw.Description, &fw.FileName, &fw.FileSize,
		&fw.ChecksumSHA256, &fw.MinioPath, &fw.Status, &fw.CreatedBy, &fw.CreatedAt, &fw.ReleasedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get firmware: %w", err)
	}
	return &fw, nil
}

// ListFirmwares returns all firmwares, optionally filtered by status or model
func (r *FirmwareRepo) ListFirmwares(ctx context.Context, status, modelFilter string) ([]*model.Firmware, error) {
	query := `SELECT id, version, model, description, file_name, file_size, checksum_sha256,
	                 minio_path, status, created_by, created_at, released_at
	          FROM firmwares WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if modelFilter != "" {
		query += fmt.Sprintf(" AND (model = $%d OR model IS NULL)", argIdx)
		args = append(args, modelFilter)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.repo.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list firmwares: %w", err)
	}
	defer rows.Close()

	var firmwares []*model.Firmware
	for rows.Next() {
		var fw model.Firmware
		if err := rows.Scan(
			&fw.ID, &fw.Version, &fw.Model, &fw.Description, &fw.FileName, &fw.FileSize,
			&fw.ChecksumSHA256, &fw.MinioPath, &fw.Status, &fw.CreatedBy, &fw.CreatedAt, &fw.ReleasedAt,
		); err != nil {
			return nil, fmt.Errorf("scan firmware: %w", err)
		}
		firmwares = append(firmwares, &fw)
	}

	return firmwares, rows.Err()
}

// UpdateFirmwareFile updates firmware file metadata after upload
func (r *FirmwareRepo) UpdateFirmwareFile(ctx context.Context, id, fileName string, fileSize int64, checksum, minioPath string) error {
	tag, err := r.repo.Pool().Exec(ctx, `
		UPDATE firmwares
		SET file_name = $1, file_size = $2, checksum_sha256 = $3, minio_path = $4
		WHERE id = $5
	`, fileName, fileSize, checksum, minioPath, id)
	if err != nil {
		return fmt.Errorf("update firmware file: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("firmware not found")
	}
	return nil
}

// UpdateFirmwareStatus changes the firmware lifecycle status
func (r *FirmwareRepo) UpdateFirmwareStatus(ctx context.Context, id string, status model.FirmwareStatus) error {
	var releasedAt *time.Time
	now := time.Now()
	if status == model.FirmwareReleased {
		releasedAt = &now
	}

	tag, err := r.repo.Pool().Exec(ctx, `
		UPDATE firmwares SET status = $1, released_at = $2 WHERE id = $3
	`, status, releasedAt, id)
	if err != nil {
		return fmt.Errorf("update firmware status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("firmware not found")
	}

	r.logger.Info("Firmware status updated",
		zap.String("id", id),
		zap.String("status", string(status)),
	)

	return nil
}

// DeleteFirmware removes a firmware record (returns minio_path for cleanup)
func (r *FirmwareRepo) DeleteFirmware(ctx context.Context, id string) (minioPath string, err error) {
	// Get minio_path first for cleanup
	var mp *string
	err = r.repo.Pool().QueryRow(ctx, `SELECT minio_path FROM firmwares WHERE id = $1`, id).Scan(&mp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("firmware not found")
		}
		return "", fmt.Errorf("get firmware for delete: %w", err)
	}

	_, err = r.repo.Pool().Exec(ctx, `DELETE FROM firmwares WHERE id = $1`, id)
	if err != nil {
		return "", fmt.Errorf("delete firmware: %w", err)
	}

	if mp != nil {
		minioPath = *mp
	}

	r.logger.Info("Firmware deleted", zap.String("id", id))
	return minioPath, nil
}

// ============================================================
// Firmware Update CRUD
// ============================================================

// CreateUpdate creates one firmware update record
func (r *FirmwareRepo) CreateUpdate(ctx context.Context, firmwareID, deviceID, campaignName, fromVersion, toVersion string, priority, maxRetries int) (*model.FirmwareUpdate, error) {
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var u model.FirmwareUpdate
	err := r.repo.Pool().QueryRow(ctx, `
		INSERT INTO firmware_updates (firmware_id, device_id, campaign_name, from_version, to_version, priority, max_retries)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, firmware_id, device_id, campaign_name, status, from_version, to_version,
		          priority, retry_count, max_retries, started_at, completed_at, error_message, created_at, updated_at
	`,
		firmwareID, deviceID, campaignName, fromVersion, toVersion, priority, maxRetries,
	).Scan(
		&u.ID, &u.FirmwareID, &u.DeviceID, &u.CampaignName, &u.Status, &u.FromVersion, &u.ToVersion,
		&u.Priority, &u.RetryCount, &u.MaxRetries, &u.StartedAt, &u.CompletedAt, &u.ErrorMessage,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert firmware update: %w", err)
	}

	return &u, nil
}

// BatchCreateUpdates creates multiple firmware updates in a single transaction
func (r *FirmwareRepo) BatchCreateUpdates(ctx context.Context, firmwareID, campaignName, fromVersion, toVersion string, deviceIDs []string, priority, maxRetries int) ([]*model.FirmwareUpdate, error) {
	tx, err := r.repo.Pool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if maxRetries <= 0 {
		maxRetries = 3
	}

	var updates []*model.FirmwareUpdate
	for _, deviceID := range deviceIDs {
		var u model.FirmwareUpdate
		err := tx.QueryRow(ctx, `
			INSERT INTO firmware_updates (firmware_id, device_id, campaign_name, from_version, to_version, priority, max_retries)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, firmware_id, device_id, campaign_name, status, from_version, to_version,
			          priority, retry_count, max_retries, started_at, completed_at, error_message, created_at, updated_at
		`,
			firmwareID, deviceID, campaignName, fromVersion, toVersion, priority, maxRetries,
		).Scan(
			&u.ID, &u.FirmwareID, &u.DeviceID, &u.CampaignName, &u.Status, &u.FromVersion, &u.ToVersion,
			&u.Priority, &u.RetryCount, &u.MaxRetries, &u.StartedAt, &u.CompletedAt, &u.ErrorMessage,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("batch insert for device %s: %w", deviceID, err)
		}
		updates = append(updates, &u)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit batch: %w", err)
	}

	r.logger.Info("Firmware update campaign created",
		zap.String("campaign", campaignName),
		zap.Int("devices", len(deviceIDs)),
	)

	return updates, nil
}

// GetUpdateByID retrieves a firmware update by ID
func (r *FirmwareRepo) GetUpdateByID(ctx context.Context, id string) (*model.FirmwareUpdate, error) {
	var u model.FirmwareUpdate
	err := r.repo.Pool().QueryRow(ctx, `
		SELECT id, firmware_id, device_id, campaign_name, status, from_version, to_version,
		       priority, retry_count, max_retries, started_at, completed_at, error_message, created_at, updated_at
		FROM firmware_updates WHERE id = $1
	`, id).Scan(
		&u.ID, &u.FirmwareID, &u.DeviceID, &u.CampaignName, &u.Status, &u.FromVersion, &u.ToVersion,
		&u.Priority, &u.RetryCount, &u.MaxRetries, &u.StartedAt, &u.CompletedAt, &u.ErrorMessage,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get update: %w", err)
	}
	return &u, nil
}

// ListUpdates returns firmware updates, optionally filtered
func (r *FirmwareRepo) ListUpdates(ctx context.Context, deviceID, statusFilter, campaignName interface{}) ([]*model.FirmwareUpdate, error) {
	query := `SELECT id, firmware_id, device_id, campaign_name, status, from_version, to_version,
	                 priority, retry_count, max_retries, started_at, completed_at, error_message, created_at, updated_at
	          FROM firmware_updates WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if deviceID != "" {
		query += fmt.Sprintf(" AND device_id = $%d", argIdx)
		args = append(args, deviceID)
		argIdx++
	}
	if statusFilter != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, statusFilter)
		argIdx++
	}
	if campaignName != "" {
		query += fmt.Sprintf(" AND campaign_name = $%d", argIdx)
		args = append(args, campaignName)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT 500"

	rows, err := r.repo.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list updates: %w", err)
	}
	defer rows.Close()

	var updates []*model.FirmwareUpdate
	for rows.Next() {
		var u model.FirmwareUpdate
		if err := rows.Scan(
			&u.ID, &u.FirmwareID, &u.DeviceID, &u.CampaignName, &u.Status, &u.FromVersion, &u.ToVersion,
			&u.Priority, &u.RetryCount, &u.MaxRetries, &u.StartedAt, &u.CompletedAt, &u.ErrorMessage,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan update: %w", err)
		}
		updates = append(updates, &u)
	}

	return updates, rows.Err()
}

// UpdateProgress updates the status of a firmware update
func (r *FirmwareRepo) UpdateProgress(ctx context.Context, id string, status model.UpdateStatus, errMsg string) error {
	var startedAt, completedAt *time.Time
	now := time.Now()

	if status == model.UpdateDownloading || status == model.UpdateInstalling {
		startedAt = &now
	}
	if status.IsTerminal() {
		completedAt = &now
	}

	tag, err := r.repo.Pool().Exec(ctx, `
		UPDATE firmware_updates
		SET status = $1, error_message = $2, started_at = COALESCE($3, started_at),
		    completed_at = $4, updated_at = NOW()
		WHERE id = $5
	`, status, errMsg, startedAt, completedAt, id)
	if err != nil {
		return fmt.Errorf("update progress: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update not found")
	}

	// If success, update device firmware_version
	if status == model.UpdateSuccess {
		var toVersion, deviceID string
		r.repo.Pool().QueryRow(ctx, `SELECT to_version, device_id FROM firmware_updates WHERE id = $1`, id).Scan(&toVersion, &deviceID)
		r.repo.Pool().Exec(ctx, `UPDATE devices SET firmware_version = $1, updated_at = NOW() WHERE id = $2`, toVersion, deviceID)
	}

	r.logger.Info("Firmware update progress",
		zap.String("id", id),
		zap.String("status", string(status)),
	)

	return nil
}

// RetryUpdate increments retry counter and resets status to pending
func (r *FirmwareRepo) RetryUpdate(ctx context.Context, id string) error {
	var retryCount, maxRetries int
	err := r.repo.Pool().QueryRow(ctx,
		`SELECT retry_count, max_retries FROM firmware_updates WHERE id = $1`, id,
	).Scan(&retryCount, &maxRetries)
	if err != nil {
		return fmt.Errorf("get retry info: %w", err)
	}

	if retryCount >= maxRetries {
		return fmt.Errorf("max retries (%d) exceeded", maxRetries)
	}

	tag, err := r.repo.Pool().Exec(ctx, `
		UPDATE firmware_updates
		SET status = 'pending', retry_count = retry_count + 1, error_message = '', updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("retry update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update not found")
	}

	return nil
}

// RollbackUpdate creates a reverse update back to a previous version
func (r *FirmwareRepo) RollbackUpdate(ctx context.Context, id string) (*model.FirmwareUpdate, error) {
	update, err := r.GetUpdateByID(ctx, id)
	if err != nil || update == nil {
		return nil, fmt.Errorf("update not found")
	}

	if update.Status != model.UpdateSuccess {
		return nil, fmt.Errorf("can only rollback successful updates, current: %s", update.Status)
	}

	// Create reverse update
	rollback, err := r.CreateUpdate(ctx, update.FirmwareID, update.DeviceID, update.CampaignNameStr(),
		update.ToVersion, update.FromVersionStr(), 2, 3) // high priority
	if err != nil {
		return nil, fmt.Errorf("create rollback: %w", err)
	}

	// Mark original as rolled_back
	r.repo.Pool().Exec(ctx, `
		UPDATE firmware_updates SET status = 'rolled_back', updated_at = NOW() WHERE id = $1
	`, id)

	// Revert device firmware_version
	r.repo.Pool().Exec(ctx, `
		UPDATE devices SET firmware_version = $1, updated_at = NOW() WHERE id = $2
	`, update.FromVersionStr(), update.DeviceID)

	r.logger.Info("Firmware update rolled back",
		zap.String("original_id", id),
		zap.String("rollback_id", rollback.ID),
		zap.String("from", update.ToVersion),
		zap.String("to", update.FromVersionStr()),
	)

	return rollback, nil
}
