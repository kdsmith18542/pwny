// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kdsmith18542/pwny/internal/db"
)

type createServiceRequest struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"`
}

type updateServiceRequest struct {
	Name  *string `json:"name,omitempty"`
	Info  *string `json:"info,omitempty"`
	State *string `json:"state,omitempty"`
}

func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	hostID := chi.URLParam(r, "id")

	services, err := s.database.ListServices(hostID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if services == nil {
		services = []db.Service{}
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: services})
}

func (s *Server) handleCreateService(w http.ResponseWriter, r *http.Request) {
	hostID := chi.URLParam(r, "id")

	var req createServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Port < 1 || req.Port > 65535 {
		writeError(w, http.StatusBadRequest, strconv.ErrRange)
		return
	}
	if req.Proto == "" {
		req.Proto = "tcp"
	}

	service, err := s.database.CreateService(hostID, req.Port, req.Proto)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Status: "ok", Data: service})
}

func (s *Server) handleGetService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	service, err := s.database.GetService(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: service})
}

func (s *Server) handleDeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.database.DeleteService(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"deleted": id}})
}

func (s *Server) handleUpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	service, err := s.database.UpdateService(id, func(svc *db.Service) {
		if req.Name != nil {
			svc.Name = *req.Name
		}
		if req.Info != nil {
			svc.Info = *req.Info
		}
		if req.State != nil {
			svc.State = *req.State
		}
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: service})
}
