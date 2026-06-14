package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/industrial-iot/device-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupTestServices(t *testing.T) (*RegistryService, *ShadowService) {
	t.Helper()

	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5433")
	db := getEnv("POSTGRES_DATABASE", "industrial_iot")
	user := getEnv("POSTGRES_USER", "iiot")
	pass := getEnv("POSTGRES_PASSWORD", "iiot_dev_2024")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=5", user, pass, host, port, db)

	logger := zap.NewNop()
	repo, err := repository.NewPostgresRepo(connStr, logger)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: getEnv("REDIS_PASSWORD", "iiot_dev_2024"),
	})

	return NewRegistryService(repo, logger), NewShadowService(repo, redisClient, logger)
}

func TestRegistryService_RegisterAndGet(t *testing.T) {
	reg, _ := setupTestServices(t)
	ctx := context.Background()

	req := &model.CreateDeviceRequest{
		SerialNumber: "SVC-TEST-REGISTER",
		Name:         "Service Test Device",
		FactoryID:    "HN",
	}

	d, err := reg.RegisterDevice(ctx, req)
	if err != nil {
		t.Fatalf("RegisterDevice failed: %v", err)
	}
	if d.ID == "" {
		t.Error("ID should not be empty")
	}
	defer reg.DeleteDevice(ctx, d.ID)

	// Get
	got, err := reg.GetDevice(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDevice failed: %v", err)
	}
	if got.SerialNumber != "SVC-TEST-REGISTER" {
		t.Errorf("serial = %s, want SVC-TEST-REGISTER", got.SerialNumber)
	}
}

func TestRegistryService_ListAndFilter(t *testing.T) {
	reg, _ := setupTestServices(t)
	ctx := context.Background()

	ids := []string{}
	for i := 1; i <= 2; i++ {
		d, err := reg.RegisterDevice(ctx, &model.CreateDeviceRequest{
			SerialNumber: fmt.Sprintf("SVC-TEST-LIST-%d", i),
			Name:         fmt.Sprintf("List %d", i),
			FactoryID:    "HCM",
		})
		if err != nil {
			t.Fatalf("RegisterDevice %d failed: %v", i, err)
		}
		ids = append(ids, d.ID)
	}
	defer func() {
		for _, id := range ids {
			reg.DeleteDevice(ctx, id)
		}
	}()

	devices, total, err := reg.ListDevices(ctx, model.DeviceFilter{FactoryID: "HCM", Limit: 10})
	if err != nil {
		t.Fatalf("ListDevices failed: %v", err)
	}
	if total < 2 {
		t.Errorf("total should be >= 2, got %d", total)
	}
	_ = devices
}

func TestRegistryService_UpdateAndStatus(t *testing.T) {
	reg, _ := setupTestServices(t)
	ctx := context.Background()

	d, err := reg.RegisterDevice(ctx, &model.CreateDeviceRequest{
		SerialNumber: "SVC-TEST-UPDATE",
		Name:         "Original",
		FactoryID:    "DN",
	})
	if err != nil {
		t.Fatalf("RegisterDevice failed: %v", err)
	}
	defer reg.DeleteDevice(ctx, d.ID)

	// Update
	newName := "Updated"
	updated, err := reg.UpdateDevice(ctx, d.ID, &model.UpdateDeviceRequest{Name: &newName})
	if err != nil {
		t.Fatalf("UpdateDevice failed: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("name = %s, want Updated", updated.Name)
	}

	// Status
	err = reg.UpdateDeviceStatus(ctx, d.ID, model.StatusMaintenance)
	if err != nil {
		t.Fatalf("UpdateDeviceStatus failed: %v", err)
	}

	got, _ := reg.GetDevice(ctx, d.ID)
	if got.Status != model.StatusMaintenance {
		t.Errorf("status = %s, want maintenance", got.Status)
	}

	// Stats
	stats, err := reg.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.Total <= 0 {
		t.Error("stats should have at least 1 device")
	}
}

func TestRegistryService_Delete(t *testing.T) {
	reg, _ := setupTestServices(t)
	ctx := context.Background()

	d, err := reg.RegisterDevice(ctx, &model.CreateDeviceRequest{
		SerialNumber: "SVC-TEST-DELETE",
		Name:         "To Delete",
		FactoryID:    "HN",
	})
	if err != nil {
		t.Fatalf("RegisterDevice failed: %v", err)
	}

	err = reg.DeleteDevice(ctx, d.ID)
	if err != nil {
		t.Fatalf("DeleteDevice failed: %v", err)
	}

	_, err = reg.GetDevice(ctx, d.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestShadowService_DesiredAndReported(t *testing.T) {
	reg, shadow := setupTestServices(t)
	ctx := context.Background()

	d, err := reg.RegisterDevice(ctx, &model.CreateDeviceRequest{
		SerialNumber: "SVC-TEST-SHADOW",
		Name:         "Shadow Test",
		FactoryID:    "HN",
	})
	if err != nil {
		t.Fatalf("RegisterDevice failed: %v", err)
	}
	defer reg.DeleteDevice(ctx, d.ID)

	// Get initial shadow
	s, err := shadow.GetShadow(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetShadow failed: %v", err)
	}
	if s.Version != 1 {
		t.Errorf("version = %d, want 1", s.Version)
	}

	// Set desired
	desired := map[string]any{"target": 100.0}
	s2, err := shadow.SetDesiredState(ctx, d.ID, desired)
	if err != nil {
		t.Fatalf("SetDesiredState failed: %v", err)
	}
	if s2.Version != 2 {
		t.Errorf("version = %d, want 2", s2.Version)
	}
	if len(s2.Delta) != 1 {
		t.Errorf("delta len = %d, want 1", len(s2.Delta))
	}

	// Set reported matching
	reported := map[string]any{"target": 100.0}
	s3, err := shadow.SetReportedState(ctx, d.ID, reported)
	if err != nil {
		t.Fatalf("SetReportedState failed: %v", err)
	}
	if len(s3.Delta) != 0 {
		t.Errorf("delta len = %d, want 0 (all match)", len(s3.Delta))
	}
}
