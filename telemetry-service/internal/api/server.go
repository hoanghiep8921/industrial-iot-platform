package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/telemetry-service/internal/storage"
	"go.uber.org/zap"
)

// Server provides the REST API for telemetry data
type Server struct {
	db     *storage.TimescaleDB
	logger *zap.Logger
	router *mux.Router
	port   int
}

// NewServer creates a new API server
func NewServer(db *storage.TimescaleDB, port int, logger *zap.Logger) *Server {
	s := &Server{
		db:     db,
		logger: logger,
		port:   port,
		router: mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Query endpoints
	api.HandleFunc("/telemetry/recent", s.handleGetRecent).Methods("GET")
	api.HandleFunc("/telemetry/device/{deviceId}", s.handleGetByDevice).Methods("GET")
	api.HandleFunc("/telemetry/latest", s.handleGetLatest).Methods("GET")

	// Health
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Stats
	api.HandleFunc("/stats", s.handleStats).Methods("GET")

	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("Starting API server", zap.String("addr", addr))
	return http.ListenAndServe(addr, s.router)
}

// ============================================================
// Handlers
// ============================================================

func (s *Server) handleGetRecent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
		if limit > 1000 {
			limit = 1000
		}
	}

	records, err := s.db.QueryRecent(ctx, limit)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(records),
		"records": records,
	})
}

func (s *Server) handleGetByDevice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	deviceID := vars["deviceId"]

	// Parse time range
	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	if f := r.URL.Query().Get("from"); f != "" {
		if parsed, err := time.Parse(time.RFC3339, f); err == nil {
			from = parsed
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			to = parsed
		}
	}

	records, err := s.db.QueryByDevice(ctx, deviceID, from, to)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"device_id": deviceID,
		"from":      from.Format(time.RFC3339),
		"to":        to.Format(time.RFC3339),
		"count":     len(records),
		"records":   records,
	})
}

func (s *Server) handleGetLatest(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	factoryID := r.URL.Query().Get("factory_id")

	records, err := s.db.GetLatestValues(ctx, factoryID)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"factory_id": factoryID,
		"count":      len(records),
		"records":    records,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "telemetry-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	stats, err := s.db.GetStorageStats(ctx)
	if err != nil {
		stats = map[string]interface{}{
			"service": "telemetry-service",
			"error":   err.Error(),
		}
	} else {
		stats["service"] = "telemetry-service"
	}

	s.writeJSON(w, http.StatusOK, stats)
}

// ============================================================
// Helpers
// ============================================================

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
	s.writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Debug("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
