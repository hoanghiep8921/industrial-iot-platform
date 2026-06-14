package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/industrial-iot/device-service/internal/repository"
	"go.uber.org/zap"
)

// RegistryService handles device registration business logic
type RegistryService struct {
	repo   *repository.PostgresRepo
	logger *zap.Logger
}

// NewRegistryService creates a new registry service
func NewRegistryService(repo *repository.PostgresRepo, logger *zap.Logger) *RegistryService {
	return &RegistryService{repo: repo, logger: logger}
}

// RegisterDevice validates and creates a new device
func (s *RegistryService) RegisterDevice(ctx context.Context, req *model.CreateDeviceRequest) (*model.Device, error) {
	// Validate required fields
	if strings.TrimSpace(req.SerialNumber) == "" {
		return nil, fmt.Errorf("serial_number is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.FactoryID) == "" {
		return nil, fmt.Errorf("factory_id is required")
	}

	// Normalize serial number
	req.SerialNumber = strings.TrimSpace(req.SerialNumber)

	// Check uniqueness
	existing, err := s.repo.GetDeviceBySerialNumber(ctx, req.SerialNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check serial_number: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("serial_number '%s' already exists", req.SerialNumber)
	}

	return s.repo.CreateDevice(ctx, req)
}

// GetDevice returns a device by ID
func (s *RegistryService) GetDevice(ctx context.Context, id string) (*model.Device, error) {
	device, err := s.repo.GetDeviceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device %s not found", id)
	}
	return device, nil
}

// ListDevices returns devices matching the filter
func (s *RegistryService) ListDevices(ctx context.Context, filter model.DeviceFilter) ([]*model.Device, int, error) {
	return s.repo.ListDevices(ctx, filter)
}

// UpdateDevice updates a device's metadata
func (s *RegistryService) UpdateDevice(ctx context.Context, id string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	// Ensure device exists
	if _, err := s.GetDevice(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.UpdateDevice(ctx, id, req)
}

// DeleteDevice removes a device
func (s *RegistryService) DeleteDevice(ctx context.Context, id string) error {
	// Ensure device exists
	if _, err := s.GetDevice(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteDevice(ctx, id)
}

// UpdateDeviceStatus validates and changes device status
func (s *RegistryService) UpdateDeviceStatus(ctx context.Context, id string, status model.DeviceStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid status '%s': must be one of %v", status, model.ValidDeviceStatuses)
	}

	// Ensure device exists
	if _, err := s.GetDevice(ctx, id); err != nil {
		return err
	}

	return s.repo.UpdateDeviceStatus(ctx, id, status)
}

// GetStats returns device statistics
func (s *RegistryService) GetStats(ctx context.Context) (*model.DeviceStats, error) {
	return s.repo.GetStats(ctx)
}
