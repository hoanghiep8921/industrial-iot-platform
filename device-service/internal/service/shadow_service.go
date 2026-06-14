package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/industrial-iot/device-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ShadowService handles device shadow business logic
type ShadowService struct {
	repo   *repository.PostgresRepo
	redis  *redis.Client
	logger *zap.Logger
}

// NewShadowService creates a new shadow service
func NewShadowService(repo *repository.PostgresRepo, redisClient *redis.Client, logger *zap.Logger) *ShadowService {
	return &ShadowService{repo: repo, redis: redisClient, logger: logger}
}

// GetShadow returns the current device shadow
func (s *ShadowService) GetShadow(ctx context.Context, deviceID string) (*model.DeviceShadow, error) {
	shadow, err := s.repo.GetShadow(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if shadow == nil {
		return nil, fmt.Errorf("shadow not found for device %s", deviceID)
	}
	return shadow, nil
}

// SetDesiredState updates the desired cloud state
func (s *ShadowService) SetDesiredState(ctx context.Context, deviceID string, state map[string]any) (*model.DeviceShadow, error) {
	shadow, err := s.repo.UpsertDesiredState(ctx, deviceID, state)
	if err != nil {
		return nil, fmt.Errorf("set desired state: %w", err)
	}

	// Cache in Redis with 5-minute TTL
	key := fmt.Sprintf("device:%s:desired", deviceID)
	data, _ := json.Marshal(state)
	s.redis.Set(ctx, key, data, 5*time.Minute)

	s.logger.Info("Desired state updated",
		zap.String("device_id", deviceID),
		zap.Int64("version", shadow.Version),
		zap.Any("delta", shadow.Delta),
	)

	return shadow, nil
}

// SetReportedState processes state reported by the device
func (s *ShadowService) SetReportedState(ctx context.Context, deviceID string, state map[string]any) (*model.DeviceShadow, error) {
	shadow, err := s.repo.UpsertReportedState(ctx, deviceID, state)
	if err != nil {
		return nil, fmt.Errorf("set reported state: %w", err)
	}

	// Cache in Redis
	key := fmt.Sprintf("device:%s:reported", deviceID)
	data, _ := json.Marshal(state)
	s.redis.Set(ctx, key, data, 5*time.Minute)

	s.logger.Info("Reported state updated",
		zap.String("device_id", deviceID),
		zap.Int64("version", shadow.Version),
		zap.Any("delta", shadow.Delta),
	)

	return shadow, nil
}
