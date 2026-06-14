package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/alarm-service/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	db     *storage.PostgresRepo
	logger *zap.Logger
	router *mux.Router
	port   int
}

func NewServer(db *storage.PostgresRepo, port int, logger *zap.Logger) *Server {
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
	s.router.Use(recoveryMiddleware(s.logger))
	s.router.Use(corsMiddleware)
	s.router.Use(loggingMiddleware(s.logger))

	// Catch-all OPTIONS for CORS preflight
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Alarm Rules CRUD
	api.HandleFunc("/alarms/rules", s.handleCreateRule).Methods("POST")
	api.HandleFunc("/alarms/rules", s.handleListRules).Methods("GET")
	api.HandleFunc("/alarms/rules/{id}", s.handleGetRule).Methods("GET")
	api.HandleFunc("/alarms/rules/{id}", s.handleUpdateRule).Methods("PUT")
	api.HandleFunc("/alarms/rules/{id}", s.handleDeleteRule).Methods("DELETE")

	// Alarms Management
	api.HandleFunc("/alarms/active", s.handleListActiveAlarms).Methods("GET")
	api.HandleFunc("/alarms/{id}/acknowledge", s.handleAcknowledgeAlarm).Methods("POST")
	api.HandleFunc("/alarms/{id}/resolve", s.handleResolveAlarm).Methods("POST")
	api.HandleFunc("/alarms/history", s.handleAlarmHistory).Methods("GET")
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("Starting Alarm Service API server", zap.String("addr", addr))
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
