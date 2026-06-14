package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/industrial-iot/device-service/internal/model"
)

// ============================================================
// Firmware Handlers
// ============================================================

func (s *Server) handleCreateFirmware(w http.ResponseWriter, r *http.Request) {
	var req model.CreateFirmwareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	fw, err := s.firmware.CreateFirmware(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, fw)
}

func (s *Server) handleGetFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	fw, err := s.firmware.GetFirmware(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, fw)
}

func (s *Server) handleListFirmwares(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	firmwares, err := s.firmware.ListFirmwares(r.Context(), q.Get("status"), q.Get("model"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if firmwares == nil {
		firmwares = []*model.Firmware{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":     len(firmwares),
		"firmwares": firmwares,
	})
}

func (s *Server) handleUploadFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Limit upload to 100MB
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "file too large (max 100MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file: use field name 'file'")
		return
	}
	defer file.Close()

	fw, err := s.firmware.UploadFirmware(r.Context(), id, file, header.Filename, header.Size, 100)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, fw)
}

func (s *Server) handleDownloadFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	reader, fw, err := s.firmware.DownloadFirmware(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fw.FileNameStr()+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fw.FileSizeVal(), 10))

	io.Copy(w, reader)
}

func (s *Server) handleReleaseFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.firmware.ReleaseFirmware(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	fw, _ := s.firmware.GetFirmware(r.Context(), id)
	writeJSON(w, http.StatusOK, fw)
}

func (s *Server) handleDeprecateFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.firmware.DeprecateFirmware(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	fw, _ := s.firmware.GetFirmware(r.Context(), id)
	writeJSON(w, http.StatusOK, fw)
}

func (s *Server) handleDeleteFirmware(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.firmware.DeleteFirmware(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// Firmware Update Handlers
// ============================================================

func (s *Server) handleCreateUpdate(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	updates, err := s.firmware.CreateUpdateCampaign(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updates) == 1 {
		writeJSON(w, http.StatusCreated, updates[0])
	} else {
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"campaign_name":  req.CampaignName,
			"total_devices":  len(updates),
			"updates":        updates,
		})
	}
}

func (s *Server) handleGetUpdate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	update, err := s.firmware.GetUpdate(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, update)
}

func (s *Server) handleListUpdates(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	updates, err := s.firmware.ListUpdates(r.Context(),
		q.Get("device_id"),
		q.Get("status"),
		q.Get("campaign_name"),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if updates == nil {
		updates = []*model.FirmwareUpdate{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(updates),
		"updates": updates,
	})
}

func (s *Server) handleUpdateProgress(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req model.UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := s.firmware.UpdateProgress(r.Context(), id, req.Status, req.ErrorMessage); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	update, _ := s.firmware.GetUpdate(r.Context(), id)
	writeJSON(w, http.StatusOK, update)
}

func (s *Server) handleRetryUpdate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := s.firmware.RetryUpdate(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	update, _ := s.firmware.GetUpdate(r.Context(), id)
	writeJSON(w, http.StatusOK, update)
}

func (s *Server) handleRollbackUpdate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	rollback, err := s.firmware.RollbackUpdate(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, rollback)
}
