package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/alarm-service/internal/model"
)

type AcknowledgeRequest struct {
	AckBy string `json:"ack_by"`
}

func (s *Server) handleListActiveAlarms(w http.ResponseWriter, r *http.Request) {
	alarms, err := s.db.ListActiveAlarms(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, alarms)
}

func (s *Server) handleAcknowledgeAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ackBy := "operator"
	if r.Body != nil && r.ContentLength > 0 {
		var req AcknowledgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.AckBy != "" {
			ackBy = req.AckBy
		}
	}

	alarm, err := s.db.AcknowledgeAlarm(r.Context(), id, ackBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, alarm)
}

func (s *Server) handleResolveAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	alarm, err := s.db.ResolveAlarm(r.Context(), id, time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, alarm)
}

func (s *Server) handleAlarmHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	page, _ := strconv.Atoi(query.Get("page"))

	filter := model.AlarmFilter{
		FactoryID: query.Get("factory_id"),
		MachineID: query.Get("machine_id"),
		Status:    model.AlarmStatus(query.Get("status")),
		Severity:  query.Get("severity"),
		Limit:     limit,
		Page:      page,
	}

	alarms, total, err := s.db.ListAlarmHistory(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  alarms,
		"total": total,
		"limit": filter.Limit,
		"page":  filter.Page,
	})
}
