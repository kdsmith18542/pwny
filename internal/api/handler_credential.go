// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kdsmith18542/pwny/internal/db"
)

type createCredentialRequest struct {
	HostID      string `json:"host_id,omitempty"`
	Username    string `json:"username"`
	PasswordEnc string `json:"password_enc,omitempty"`
	HashEnc     string `json:"hash_enc,omitempty"`
	Type        string `json:"type"`
	Module      string `json:"module,omitempty"`
}

func (s *Server) handleListCredentials(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "id")

	if _, err := s.database.GetWorkspace(workspaceID); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	creds, err := s.database.ListCredentials(workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if creds == nil {
		creds = []db.Credential{}
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: creds})
}

func (s *Server) handleCreateCredential(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "id")

	if _, err := s.database.GetWorkspace(workspaceID); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	var req createCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Username == "" {
		writeError(w, http.StatusBadRequest, errNameRequired())
		return
	}

	cred, err := s.database.CreateCredential(workspaceID, req.HostID, req.Username, req.PasswordEnc, req.HashEnc, req.Type, req.Module)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Status: "ok", Data: cred})
}

func (s *Server) handleGetCredential(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cred, err := s.database.GetCredential(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: cred})
}

func (s *Server) handleDeleteCredential(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := s.database.DeleteCredential(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"deleted": id}})
}
