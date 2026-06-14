package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ============================================================
// Device Shadow Handlers
// ============================================================

func (s *Server) handleGetShadow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	shadow, err := s.shadow.GetShadow(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, shadow)
}

func (s *Server) handleSetDesired(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var state map[string]any
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if len(state) == 0 {
		writeError(w, http.StatusBadRequest, "state cannot be empty")
		return
	}

	shadow, err := s.shadow.SetDesiredState(r.Context(), id, state)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, shadow)
}

func (s *Server) handleSetReported(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var state map[string]any
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if len(state) == 0 {
		writeError(w, http.StatusBadRequest, "state cannot be empty")
		return
	}

	shadow, err := s.shadow.SetReportedState(r.Context(), id, state)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, shadow)
}
