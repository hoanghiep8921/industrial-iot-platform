package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/device-service/internal/model"
)

// ============================================================
// Device Registry Handlers
// ============================================================

func (s *Server) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	var req model.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	device, err := s.registry.RegisterDevice(r.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, device)
}

func (s *Server) handleGetDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	device, err := s.registry.GetDevice(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, device)
}

func (s *Server) handleListDevices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := model.DeviceFilter{
		FactoryID: q.Get("factory_id"),
		Area:      q.Get("area"),
		Status:    model.DeviceStatus(q.Get("status")),
		Protocol:  q.Get("protocol"),
		Search:    q.Get("search"),
		Page:      page,
		Limit:     limit,
	}

	devices, total, err := s.registry.ListDevices(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if devices == nil {
		devices = []*model.Device{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   total,
		"page":    page,
		"limit":   limit,
		"devices": devices,
	})
}

func (s *Server) handleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req model.UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	device, err := s.registry.UpdateDevice(r.Context(), id, &req)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, device)
}

func (s *Server) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.registry.DeleteDevice(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req model.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := s.registry.UpdateDeviceStatus(r.Context(), id, req.Status); err != nil {
		if err.Error()[:8] == "invalid " {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	// Fetch updated device
	device, _ := s.registry.GetDevice(r.Context(), id)
	writeJSON(w, http.StatusOK, device)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.registry.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "device-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
