package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisCache(host string, port int, password string, logger *zap.Logger) (*RedisCache, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis at %s: %w", addr, err)
	}

	logger.Info("Connected to Redis", zap.String("addr", addr))
	return &RedisCache{client: client, logger: logger}, nil
}

func (c *RedisCache) Close() error {
	c.logger.Info("Closing Redis client...")
	return c.client.Close()
}

func (c *RedisCache) Client() *redis.Client {
	return c.client
}

// Active failure tracking for delay (duration_seconds) logic
func (c *RedisCache) GetFailureStartTime(ctx context.Context, ruleID string, machineID string) (*time.Time, error) {
	key := fmt.Sprintf("alarm:failure:%s:%s", ruleID, machineID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *RedisCache) SetFailureStartTime(ctx context.Context, ruleID string, machineID string, t time.Time) error {
	key := fmt.Sprintf("alarm:failure:%s:%s", ruleID, machineID)
	return c.client.Set(ctx, key, t.Format(time.RFC3339), 24*time.Hour).Err()
}

func (c *RedisCache) ClearFailureStartTime(ctx context.Context, ruleID string, machineID string) error {
	key := fmt.Sprintf("alarm:failure:%s:%s", ruleID, machineID)
	return c.client.Del(ctx, key).Err()
}

// Active alarm ID mapping (stores UUID of raised alarm)
func (c *RedisCache) GetActiveAlarmID(ctx context.Context, ruleID string, machineID string) (string, error) {
	key := fmt.Sprintf("alarm:active:%s:%s", ruleID, machineID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *RedisCache) SetActiveAlarmID(ctx context.Context, ruleID string, machineID string, alarmID string) error {
	key := fmt.Sprintf("alarm:active:%s:%s", ruleID, machineID)
	return c.client.Set(ctx, key, alarmID, 0).Err() // No TTL, deleted on resolve
}

func (c *RedisCache) ClearActiveAlarmID(ctx context.Context, ruleID string, machineID string) error {
	key := fmt.Sprintf("alarm:active:%s:%s", ruleID, machineID)
	return c.client.Del(ctx, key).Err()
}
