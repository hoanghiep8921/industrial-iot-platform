package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/industrial-iot/device-service/internal/repository"
	"github.com/industrial-iot/device-service/internal/storage"
	"go.uber.org/zap"
)

// FirmwareService handles firmware management business logic
type FirmwareService struct {
	fwRepo *repository.FirmwareRepo
	store  *storage.MinioClient
	logger *zap.Logger
}

// NewFirmwareService creates a new firmware service
func NewFirmwareService(fwRepo *repository.FirmwareRepo, store *storage.MinioClient, logger *zap.Logger) *FirmwareService {
	return &FirmwareService{fwRepo: fwRepo, store: store, logger: logger}
}

// CreateFirmware registers firmware metadata
func (s *FirmwareService) CreateFirmware(ctx context.Context, req *model.CreateFirmwareRequest) (*model.Firmware, error) {
	if strings.TrimSpace(req.Version) == "" {
		return nil, fmt.Errorf("version is required")
	}
	return s.fwRepo.CreateFirmware(ctx, req)
}

// GetFirmware returns a firmware by ID
func (s *FirmwareService) GetFirmware(ctx context.Context, id string) (*model.Firmware, error) {
	fw, err := s.fwRepo.GetFirmwareByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if fw == nil {
		return nil, fmt.Errorf("firmware %s not found", id)
	}
	return fw, nil
}

// ListFirmwares returns all firmwares with optional filtering
func (s *FirmwareService) ListFirmwares(ctx context.Context, status, modelFilter string) ([]*model.Firmware, error) {
	return s.fwRepo.ListFirmwares(ctx, status, modelFilter)
}

// UploadFirmware validates and stores a firmware binary
func (s *FirmwareService) UploadFirmware(ctx context.Context, id string, reader io.Reader, fileName string, fileSize int64, maxSizeMB int64) (*model.Firmware, error) {
	// Validate file size
	if maxSizeMB <= 0 {
		maxSizeMB = 100
	}
	maxBytes := maxSizeMB * 1024 * 1024
	if fileSize > maxBytes {
		return nil, fmt.Errorf("file too large: %d bytes (max %d MB)", fileSize, maxSizeMB)
	}
	if fileSize == 0 {
		return nil, fmt.Errorf("empty file not allowed")
	}

	// Get firmware metadata
	fw, err := s.GetFirmware(ctx, id)
	if err != nil {
		return nil, err
	}

	// Read into memory to compute SHA256 (also needed for upload)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read firmware file: %w", err)
	}

	// Compute SHA256 checksum
	hash := sha256.Sum256(data)
	checksum := hex.EncodeToString(hash[:])

	// Determine MinIO object path
	ext := filepath.Ext(fileName)
	objectName := fmt.Sprintf("%s/%s%s", fw.Version, id, ext)

	// Upload to MinIO
	uploadReader := strings.NewReader(string(data)) // re-wrap for upload
	_, err = s.store.UploadFirmware(ctx, uploadReader, objectName, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to upload firmware: %w", err)
	}

	// Update metadata in DB
	err = s.fwRepo.UpdateFirmwareFile(ctx, id, fileName, int64(len(data)), checksum, objectName)
	if err != nil {
		// Rollback: delete from MinIO
		s.store.DeleteFirmware(ctx, objectName)
		return nil, fmt.Errorf("failed to update firmware metadata: %w", err)
	}

	s.logger.Info("Firmware binary uploaded",
		zap.String("id", id),
		zap.String("version", fw.Version),
		zap.String("checksum", checksum),
		zap.Int64("size", int64(len(data))),
	)

	return s.fwRepo.GetFirmwareByID(ctx, id)
}

// DownloadFirmware returns a reader for the firmware binary
func (s *FirmwareService) DownloadFirmware(ctx context.Context, id string) (io.ReadCloser, *model.Firmware, error) {
	fw, err := s.GetFirmware(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if fw.MinioPathStr() == "" {
		return nil, nil, fmt.Errorf("firmware binary not uploaded yet")
	}

	reader, err := s.store.DownloadFirmware(ctx, fw.MinioPathStr())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download firmware: %w", err)
	}

	return reader, fw, nil
}

// ReleaseFirmware publishes a firmware for OTA updates
func (s *FirmwareService) ReleaseFirmware(ctx context.Context, id string) error {
	fw, err := s.GetFirmware(ctx, id)
	if err != nil {
		return err
	}
	// Can only release if binary is uploaded
	if fw.FileNameStr() == "" {
		return fmt.Errorf("cannot release firmware without uploaded binary")
	}
	if fw.Status == model.FirmwareReleased {
		return fmt.Errorf("firmware already released")
	}
	if fw.Status == model.FirmwareDeprecated {
		return fmt.Errorf("cannot re-release deprecated firmware")
	}
	return s.fwRepo.UpdateFirmwareStatus(ctx, id, model.FirmwareReleased)
}

// DeprecateFirmware marks a firmware as deprecated
func (s *FirmwareService) DeprecateFirmware(ctx context.Context, id string) error {
	fw, err := s.GetFirmware(ctx, id)
	if err != nil {
		return err
	}
	if fw.Status != model.FirmwareReleased {
		return fmt.Errorf("can only deprecate released firmware, current: %s", fw.Status)
	}
	return s.fwRepo.UpdateFirmwareStatus(ctx, id, model.FirmwareDeprecated)
}

// DeleteFirmware removes a firmware and its binary
func (s *FirmwareService) DeleteFirmware(ctx context.Context, id string) error {
	if _, err := s.GetFirmware(ctx, id); err != nil {
		return err
	}

	minioPath, err := s.fwRepo.DeleteFirmware(ctx, id)
	if err != nil {
		return err
	}

	// Clean up MinIO
	if minioPath != "" {
		s.store.DeleteFirmware(ctx, minioPath)
	}

	return nil
}

// CreateUpdateCampaign creates firmware updates for one or more devices
func (s *FirmwareService) CreateUpdateCampaign(ctx context.Context, req *model.CreateUpdateRequest) ([]*model.FirmwareUpdate, error) {
	// Get firmware
	fw, err := s.GetFirmware(ctx, req.FirmwareID)
	if err != nil {
		return nil, err
	}
	if fw.Status != model.FirmwareReleased {
		return nil, fmt.Errorf("firmware must be released before creating updates (current: %s)", fw.Status)
	}

	if req.MaxRetries <= 0 {
		req.MaxRetries = 3
	}

	// Single device update
	if req.DeviceID != "" {
		update, err := s.fwRepo.CreateUpdate(ctx, fw.ID, req.DeviceID, req.CampaignName, "", fw.Version, req.Priority, req.MaxRetries)
		if err != nil {
			return nil, err
		}
		return []*model.FirmwareUpdate{update}, nil
	}

	// Campaign (batch) update
	if len(req.DeviceIDs) == 0 {
		return nil, fmt.Errorf("either device_id or device_ids is required")
	}

	campaignName := req.CampaignName
	if campaignName == "" {
		modelName := ""
		if fw.Model != nil { modelName = *fw.Model }
		campaignName = fmt.Sprintf("%s-%s-rollout", modelName, fw.Version)
	}

	return s.fwRepo.BatchCreateUpdates(ctx, fw.ID, campaignName, "", fw.Version, req.DeviceIDs, req.Priority, req.MaxRetries)
}

// UpdateProgress reports firmware update status from edge
func (s *FirmwareService) UpdateProgress(ctx context.Context, id string, status model.UpdateStatus, errMsg string) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid update status: %s", status)
	}
	return s.fwRepo.UpdateProgress(ctx, id, status, errMsg)
}

// RetryUpdate resets a failed update to pending
func (s *FirmwareService) RetryUpdate(ctx context.Context, id string) error {
	update, err := s.fwRepo.GetUpdateByID(ctx, id)
	if err != nil || update == nil {
		return fmt.Errorf("update not found")
	}
	if update.Status != model.UpdateFailed {
		return fmt.Errorf("can only retry failed updates, current: %s", update.Status)
	}
	return s.fwRepo.RetryUpdate(ctx, id)
}

// RollbackUpdate rolls back a successful update
func (s *FirmwareService) RollbackUpdate(ctx context.Context, id string) (*model.FirmwareUpdate, error) {
	return s.fwRepo.RollbackUpdate(ctx, id)
}

// GetUpdate returns a single firmware update by ID
func (s *FirmwareService) GetUpdate(ctx context.Context, id string) (*model.FirmwareUpdate, error) {
	update, err := s.fwRepo.GetUpdateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if update == nil {
		return nil, fmt.Errorf("firmware update %s not found", id)
	}
	return update, nil
}

// ListUpdates returns firmware updates with optional filters
func (s *FirmwareService) ListUpdates(ctx context.Context, deviceID, statusFilter, campaignName string) ([]*model.FirmwareUpdate, error) {
	return s.fwRepo.ListUpdates(ctx, deviceID, statusFilter, campaignName)
}
