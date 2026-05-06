// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kdsmith18542/pwny/internal/db"
)

type createHostRequest struct {
	Address string `json:"address"`
}

type updateHostRequest struct {
	MAC      *string `json:"mac,omitempty"`
	OSName   *string `json:"os_name,omitempty"`
	OSFlavor *string `json:"os_flavor,omitempty"`
	OSSP     *string `json:"os_sp,omitempty"`
	Arch     *string `json:"arch,omitempty"`
	Purpose  *string `json:"purpose,omitempty"`
	Info     *string `json:"info,omitempty"`
	State    *string `json:"state,omitempty"`
}

func (s *Server) handleListHosts(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "id")

	if _, err := s.database.GetWorkspace(workspaceID); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	hosts, err := s.database.ListHosts(workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if hosts == nil {
		hosts = []db.Host{}
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: hosts})
}

func (s *Server) handleCreateHost(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "id")

	if _, err := s.database.GetWorkspace(workspaceID); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	var req createHostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Address == "" {
		writeError(w, http.StatusBadRequest, errNameRequired())
		return
	}

	host, err := s.database.CreateHost(workspaceID, req.Address)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Status: "ok", Data: host})
}

func (s *Server) handleGetHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	host, err := s.database.GetHost(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: host})
}

func (s *Server) handleDeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.database.DeleteHost(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"deleted": id}})
}

func (s *Server) handleUpdateHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateHostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	host, err := s.database.UpdateHost(id, func(h *db.Host) {
		if req.MAC != nil {
			h.MAC = *req.MAC
		}
		if req.OSName != nil {
			h.OSName = *req.OSName
		}
		if req.OSFlavor != nil {
			h.OSFlavor = *req.OSFlavor
		}
		if req.OSSP != nil {
			h.OSSP = *req.OSSP
		}
		if req.Arch != nil {
			h.Arch = *req.Arch
		}
		if req.Purpose != nil {
			h.Purpose = *req.Purpose
		}
		if req.Info != nil {
			h.Info = *req.Info
		}
		if req.State != nil {
			h.State = *req.State
		}
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: host})
}
