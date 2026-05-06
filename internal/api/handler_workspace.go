// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type createWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (s *Server) handleListWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := s.database.ListWorkspaces()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: workspaces})
}

func (s *Server) handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, errNameRequired())
		return
	}

	wksp, err := s.database.CreateWorkspace(req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Status: "ok", Data: wksp})
}

func (s *Server) handleGetWorkspace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	wksp, err := s.database.GetWorkspace(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: wksp})
}

func (s *Server) handleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.database.DeleteWorkspace(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"deleted": id}})
}

func errNameRequired() error {
	return &apiArgError{msg: "name is required"}
}

type apiArgError struct {
	msg string
}

func (e *apiArgError) Error() string { return e.msg }
