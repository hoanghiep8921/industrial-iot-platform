package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/industrial-iot/device-service/internal/model"
	"go.uber.org/zap"
)

func getTestRepo(t *testing.T) *PostgresRepo {
	t.Helper()

	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5433")
	db := getEnv("POSTGRES_DATABASE", "industrial_iot")
	user := getEnv("POSTGRES_USER", "iiot")
	pass := getEnv("POSTGRES_PASSWORD", "iiot_dev_2024")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=5", user, pass, host, port, db)

	logger := zap.NewNop()
	repo, err := NewPostgresRepo(connStr, logger)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to PostgreSQL: %v", err)
	}
	return repo
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func TestPostgresRepo_CreateAndGetDevice(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	req := &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408501-001",
		Name:         "Integration Test Device",
		FactoryID:    "HN",
		Area:         "area_a",
		Protocol:     "mqtt",
		Capabilities: []string{"temperature", "speed"},
	}

	device, err := repo.CreateDevice(ctx, req)
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	if device.ID == "" {
		t.Error("device ID should not be empty")
	}
	if device.Status != model.StatusOffline {
		t.Errorf("status = %s, want offline", device.Status)
	}

	// Get device
	got, err := repo.GetDeviceByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("GetDeviceByID failed: %v", err)
	}
	if got == nil {
		t.Fatal("device should exist")
	}
	if got.SerialNumber != "ITEST-T1781408501-001" {
		t.Errorf("serial = %s, want ITEST-T1781408501-001", got.SerialNumber)
	}

	// Get by serial
	gotBySN, err := repo.GetDeviceBySerialNumber(ctx, "ITEST-T1781408501-001")
	if err != nil {
		t.Fatalf("GetDeviceBySerialNumber failed: %v", err)
	}
	if gotBySN.ID != device.ID {
		t.Errorf("IDs don't match: %s vs %s", gotBySN.ID, device.ID)
	}

	// Shadow should exist
	shadow, err := repo.GetShadow(ctx, device.ID)
	if err != nil {
		t.Fatalf("GetShadow failed: %v", err)
	}
	if shadow == nil {
		t.Fatal("shadow should exist after device creation")
	}
	if shadow.Version != 1 {
		t.Errorf("shadow version = %d, want 1", shadow.Version)
	}

	// Cleanup
	repo.DeleteDevice(ctx, device.ID)
}

func TestPostgresRepo_ListDevices(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	// Create test devices
	ids := []string{}
	for i := 1; i <= 3; i++ {
		req := &model.CreateDeviceRequest{
			SerialNumber: fmt.Sprintf("ITEST-T1781408501-LIST-%d", i),
			Name:         fmt.Sprintf("List Device %d", i),
			FactoryID:    "HCM",
			Area:         "area_b",
		}
		d, err := repo.CreateDevice(ctx, req)
		if err != nil {
			t.Fatalf("CreateDevice %d failed: %v", i, err)
		}
		ids = append(ids, d.ID)
	}
	defer func() {
		for _, id := range ids {
			repo.DeleteDevice(ctx, id)
		}
	}()

	// List all
	_, total, err := repo.ListDevices(ctx, model.DeviceFilter{Limit: 10})
	if err != nil {
		t.Fatalf("ListDevices failed: %v", err)
	}
	if total < 3 {
		t.Errorf("total should be >= 3, got %d", total)
	}

	// Filter by factory
	filtered, ftotal, err := repo.ListDevices(ctx, model.DeviceFilter{FactoryID: "HCM", Limit: 10})
	if err != nil {
		t.Fatalf("ListDevices(filtered) failed: %v", err)
	}
	if ftotal < 3 {
		t.Errorf("filtered total should be >= 3, got %d", ftotal)
	}
	if len(filtered) < 3 {
		t.Errorf("filtered results should be >= 3, got %d", len(filtered))
	}

	// Pagination
	paged, ptotal, err := repo.ListDevices(ctx, model.DeviceFilter{FactoryID: "HCM", Page: 1, Limit: 2})
	if err != nil {
		t.Fatalf("ListDevices(paged) failed: %v", err)
	}
	if ptotal < 3 {
		t.Errorf("paged total should be >= 3, got %d", ptotal)
	}
	if len(paged) > 2 {
		t.Errorf("paged results should be <= 2, got %d", len(paged))
	}
}

func TestPostgresRepo_UpdateAndDeleteDevice(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	// Create
	req := &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408501-UPDATE",
		Name:         "Original Name",
		FactoryID:    "DN",
	}
	d, err := repo.CreateDevice(ctx, req)
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	// Update
	newName := "Updated Name"
	updated, err := repo.UpdateDevice(ctx, d.ID, &model.UpdateDeviceRequest{Name: &newName})
	if err != nil {
		t.Fatalf("UpdateDevice failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("name = %s, want Updated Name", updated.Name)
	}

	// Update status
	err = repo.UpdateDeviceStatus(ctx, d.ID, model.StatusOnline)
	if err != nil {
		t.Fatalf("UpdateDeviceStatus failed: %v", err)
	}
	got, _ := repo.GetDeviceByID(ctx, d.ID)
	if got.Status != model.StatusOnline {
		t.Errorf("status = %s, want online", got.Status)
	}
	if got.LastSeenAt == nil {
		t.Error("last_seen_at should be set after status update")
	}

	// Delete
	err = repo.DeleteDevice(ctx, d.ID)
	if err != nil {
		t.Fatalf("DeleteDevice failed: %v", err)
	}
	gone, _ := repo.GetDeviceByID(ctx, d.ID)
	if gone != nil {
		t.Error("device should be gone after delete")
	}
}

func TestPostgresRepo_Shadow(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	req := &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408501-SHADOW",
		Name:         "Shadow Test Device",
		FactoryID:    "HN",
	}
	d, err := repo.CreateDevice(ctx, req)
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}
	defer repo.DeleteDevice(ctx, d.ID)

	// Set desired
	desired := map[string]any{"temp": 200, "speed": 1500}
	s1, err := repo.UpsertDesiredState(ctx, d.ID, desired)
	if err != nil {
		t.Fatalf("UpsertDesiredState failed: %v", err)
	}
	if s1.Version != 2 {
		t.Errorf("version = %d, want 2", s1.Version)
	}
	if len(s1.Delta) != 2 {
		t.Errorf("delta len = %d, want 2", len(s1.Delta))
	}

	// Set reported (partial match)
	reported := map[string]any{"temp": 200, "speed": 1400}
	s2, err := repo.UpsertReportedState(ctx, d.ID, reported)
	if err != nil {
		t.Fatalf("UpsertReportedState failed: %v", err)
	}
	if s2.Version != 3 {
		t.Errorf("version = %d, want 3", s2.Version)
	}
	if len(s2.Delta) != 1 {
		t.Errorf("delta len = %d, want 1 (only speed differs)", len(s2.Delta))
	}

	// Set reported (full match)
	reported2 := map[string]any{"temp": 200, "speed": 1500}
	s3, err := repo.UpsertReportedState(ctx, d.ID, reported2)
	if err != nil {
		t.Fatalf("UpsertReportedState failed: %v", err)
	}
	if len(s3.Delta) != 0 {
		t.Errorf("delta len = %d, want 0 (all match)", len(s3.Delta))
	}
}

func TestPostgresRepo_Stats(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.Total <= 0 {
		t.Log("No devices, stats total might be 0")
	}
	if stats.ByFactory == nil {
		t.Error("ByFactory map should not be nil")
	}
	if stats.ByStatus == nil {
		t.Error("ByStatus map should not be nil")
	}
}

func TestPostgresRepo_DuplicateSerialNumber(t *testing.T) {
	repo := getTestRepo(t)
	ctx := context.Background()

	req := &model.CreateDeviceRequest{
		SerialNumber: "ITEST-T1781408501-DUP",
		Name:         "First",
		FactoryID:    "HN",
	}
	d, err := repo.CreateDevice(ctx, req)
	if err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}
	defer repo.DeleteDevice(ctx, d.ID)

	// Try duplicate
	_, err = repo.CreateDevice(ctx, req)
	if err == nil {
		t.Error("expected error for duplicate serial number")
	}
}
