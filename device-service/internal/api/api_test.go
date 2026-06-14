package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/industrial-iot/device-service/internal/model"
	"github.com/industrial-iot/device-service/internal/repository"
	"github.com/industrial-iot/device-service/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupTestServer(t *testing.T) *Server {
	t.Helper()

	// Connect to PostgreSQL
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5433")
	db := getEnv("POSTGRES_DATABASE", "industrial_iot")
	user := getEnv("POSTGRES_USER", "iiot")
	pass := getEnv("POSTGRES_PASSWORD", "iiot_dev_2024")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=5", user, pass, host, port, db)

	logger := zap.NewNop()
	repo, err := repository.NewPostgresRepo(connStr, logger)
	if err != nil {
		t.Skipf("Skipping API integration test: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: getEnv("REDIS_PASSWORD", "iiot_dev_2024"),
	})

	registryService := service.NewRegistryService(repo, logger)
	shadowService := service.NewShadowService(repo, redisClient, logger)
	// MinIO is nil for tests - firmware upload/download won't work
	var fwService *service.FirmwareService
	_ = repository.NewFirmwareRepo(repo) // available for future firmware tests

	return NewServer(registryService, shadowService, fwService, 0, logger)
}

func TestAPI_Health(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "healthy" {
		t.Errorf("status = %s, want healthy", body["status"])
	}
	if body["service"] != "device-service" {
		t.Errorf("service = %s, want device-service", body["service"])
	}
}

func TestAPI_CreateDevice(t *testing.T) {
	server := setupTestServer(t)

	body := map[string]interface{}{
		"serial_number": "API-TEST-001",
		"name":          "API Test Device",
		"factory_id":    "HN",
		"area":          "area_a",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}

	var device model.Device
	json.NewDecoder(w.Body).Decode(&device)
	if device.SerialNumber != "API-TEST-001" {
		t.Errorf("serial = %s, want API-TEST-001", device.SerialNumber)
	}
	if device.Status != model.StatusOffline {
		t.Errorf("status = %s, want offline", device.Status)
	}
	if device.ID == "" {
		t.Error("ID should not be empty")
	}

	// Cleanup
	server.registry.DeleteDevice(req.Context(), device.ID)
}

func TestAPI_DuplicateDevice(t *testing.T) {
	server := setupTestServer(t)

	body := map[string]interface{}{
		"serial_number": "API-TEST-DUP",
		"name":          "Duplicate Test",
		"factory_id":    "HN",
	}
	bodyBytes, _ := json.Marshal(body)

	// First request
	req1 := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(bodyBytes))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	server.router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Skipf("First create failed, skipping duplicate test")
		return
	}
	var device model.Device
	json.NewDecoder(w1.Body).Decode(&device)
	defer server.registry.DeleteDevice(req1.Context(), device.ID)

	// Duplicate
	req2 := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	server.router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("status = %d, want 409", w2.Code)
	}

	var errResp map[string]string
	json.NewDecoder(w2.Body).Decode(&errResp)
	if errResp["error"] == "" {
		t.Error("expected error message in response")
	}
}

func TestAPI_GetAndListDevices(t *testing.T) {
	server := setupTestServer(t)

	// Create device first
	device, _ := server.registry.RegisterDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		&model.CreateDeviceRequest{SerialNumber: "API-TEST-LIST", Name: "List Test", FactoryID: "HCM"},
	)
	if device == nil {
		t.Skip("could not create test device")
	}
	defer server.registry.DeleteDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		device.ID,
	)

	// GET by ID
	req := httptest.NewRequest("GET", "/api/v1/devices/"+device.ID, nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /devices/{id} status = %d, want 200", w.Code)
	}

	// LIST
	req2 := httptest.NewRequest("GET", "/api/v1/devices?factory_id=HCM", nil)
	w2 := httptest.NewRecorder()
	server.router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("GET /devices status = %d, want 200", w2.Code)
	}
}

func TestAPI_UpdateDeviceStatus(t *testing.T) {
	server := setupTestServer(t)

	device, _ := server.registry.RegisterDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		&model.CreateDeviceRequest{SerialNumber: "API-TEST-STATUS", Name: "Status Test", FactoryID: "DN"},
	)
	if device == nil {
		t.Skip("could not create test device")
	}
	defer server.registry.DeleteDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		device.ID,
	)

	statusBody := map[string]interface{}{"status": "online"}
	bodyBytes, _ := json.Marshal(statusBody)
	req := httptest.NewRequest("PATCH", "/api/v1/devices/"+device.ID+"/status", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var updated model.Device
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Status != model.StatusOnline {
		t.Errorf("status = %s, want online", updated.Status)
	}
	if updated.LastSeenAt == nil {
		t.Error("last_seen_at should be set")
	}
}

func TestAPI_InvalidStatus(t *testing.T) {
	server := setupTestServer(t)

	device, _ := server.registry.RegisterDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		&model.CreateDeviceRequest{SerialNumber: "API-TEST-INVSTAT", Name: "Invalid Status Test", FactoryID: "HN"},
	)
	if device == nil {
		t.Skip("could not create test device")
	}
	defer server.registry.DeleteDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		device.ID,
	)

	statusBody := map[string]interface{}{"status": "invalid_status"}
	bodyBytes, _ := json.Marshal(statusBody)
	req := httptest.NewRequest("PATCH", "/api/v1/devices/"+device.ID+"/status", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestAPI_GetShadow(t *testing.T) {
	server := setupTestServer(t)

	device, _ := server.registry.RegisterDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		&model.CreateDeviceRequest{SerialNumber: "API-TEST-SHADOW", Name: "Shadow Test", FactoryID: "HN"},
	)
	if device == nil {
		t.Skip("could not create test device")
	}
	defer server.registry.DeleteDevice(
		httptest.NewRequest("GET", "/", nil).Context(),
		device.ID,
	)

	// Get shadow
	req := httptest.NewRequest("GET", "/api/v1/devices/"+device.ID+"/shadow", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}

	var shadow model.DeviceShadow
	json.NewDecoder(w.Body).Decode(&shadow)
	if shadow.Version != 1 {
		t.Errorf("version = %d, want 1", shadow.Version)
	}
}

func TestAPI_CORS(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest("OPTIONS", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("OPTIONS status = %d, want 200", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header should be set")
	}
}
