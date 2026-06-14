package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/device-service/internal/service"
	"go.uber.org/zap"
)

// Server provides the REST API for the device service
type Server struct {
	registry *service.RegistryService
	shadow   *service.ShadowService
	firmware *service.FirmwareService
	logger   *zap.Logger
	router   *mux.Router
	port     int
}

// NewServer creates a new API server
func NewServer(registry *service.RegistryService, shadow *service.ShadowService, firmware *service.FirmwareService, port int, logger *zap.Logger) *Server {
	s := &Server{
		registry: registry,
		shadow:   shadow,
		firmware: firmware,
		logger:   logger,
		port:     port,
		router:   mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(recoveryMiddleware(s.logger))
	s.router.Use(corsMiddleware)
	s.router.Use(loggingMiddleware(s.logger))

	// Catch-all OPTIONS for CORS preflight (must be before api subrouter)
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Health
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Device Registry
	api.HandleFunc("/devices", s.handleCreateDevice).Methods("POST")
	api.HandleFunc("/devices", s.handleListDevices).Methods("GET")
	api.HandleFunc("/devices/{id}", s.handleGetDevice).Methods("GET")
	api.HandleFunc("/devices/{id}", s.handleUpdateDevice).Methods("PUT")
	api.HandleFunc("/devices/{id}", s.handleDeleteDevice).Methods("DELETE")
	api.HandleFunc("/devices/{id}/status", s.handleUpdateStatus).Methods("PATCH")
	api.HandleFunc("/stats", s.handleStats).Methods("GET")

	// Device Shadow
	api.HandleFunc("/devices/{id}/shadow", s.handleGetShadow).Methods("GET")
	api.HandleFunc("/devices/{id}/shadow/desired", s.handleSetDesired).Methods("PUT")
	api.HandleFunc("/devices/{id}/shadow/reported", s.handleSetReported).Methods("PUT")

	// Firmware Management
	api.HandleFunc("/firmwares", s.handleCreateFirmware).Methods("POST")
	api.HandleFunc("/firmwares", s.handleListFirmwares).Methods("GET")
	api.HandleFunc("/firmwares/{id}", s.handleGetFirmware).Methods("GET")
	api.HandleFunc("/firmwares/{id}/upload", s.handleUploadFirmware).Methods("POST")
	api.HandleFunc("/firmwares/{id}/download", s.handleDownloadFirmware).Methods("GET")
	api.HandleFunc("/firmwares/{id}/release", s.handleReleaseFirmware).Methods("PUT")
	api.HandleFunc("/firmwares/{id}/deprecate", s.handleDeprecateFirmware).Methods("PUT")
	api.HandleFunc("/firmwares/{id}", s.handleDeleteFirmware).Methods("DELETE")

	// Firmware Updates
	api.HandleFunc("/firmware-updates", s.handleCreateUpdate).Methods("POST")
	api.HandleFunc("/firmware-updates", s.handleListUpdates).Methods("GET")
	api.HandleFunc("/firmware-updates/{id}", s.handleGetUpdate).Methods("GET")
	api.HandleFunc("/firmware-updates/{id}/progress", s.handleUpdateProgress).Methods("PATCH")
	api.HandleFunc("/firmware-updates/{id}/retry", s.handleRetryUpdate).Methods("POST")
	api.HandleFunc("/firmware-updates/{id}/rollback", s.handleRollbackUpdate).Methods("POST")
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("Starting API server", zap.String("addr", addr))
	return http.ListenAndServe(addr, s.router)
}

// ============================================================
// Helpers
// ============================================================

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
