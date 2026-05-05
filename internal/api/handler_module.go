// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/msfgo/msfgo/internal/core"
)

type moduleResponse struct {
	Name        string               `json:"name"`
	Type        core.ModuleType      `json:"type"`
	Description string               `json:"description"`
	Authors     []string             `json:"authors"`
	Platforms   []string             `json:"platforms"`
	Arch        []string             `json:"arch"`
	Options     map[string]optionDef `json:"options,omitempty"`
}

type optionDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
}

type runRequest struct {
	Options map[string]interface{} `json:"options"`
}

func (s *Server) handleListModules(w http.ResponseWriter, r *http.Request) {
	filterType := r.URL.Query().Get("type")
	infos := core.ListModules(core.ModuleType(filterType))
	modules := make([]moduleResponse, 0, len(infos))
	for _, info := range infos {
		modules = append(modules, moduleResponse{
			Name:        info.Name,
			Type:        info.Type,
			Description: info.Description,
			Authors:     info.Authors,
			Platforms:   info.Platforms,
			Arch:        info.Arch,
		})
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: modules})
}

func (s *Server) handleGetModule(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	m, err := core.GetModule(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	opts := make(map[string]optionDef)
	for k, v := range m.Options() {
		opts[k] = optionDef{
			Name:        v.Name,
			Description: v.Description,
			Required:    v.Required,
			Default:     v.Default,
		}
	}

	info := m.Info()
	resp := moduleResponse{
		Name:        info.Name,
		Type:        info.Type,
		Description: info.Description,
		Authors:     info.Authors,
		Platforms:   info.Platforms,
		Arch:        info.Arch,
		Options:     opts,
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: resp})
}

func (s *Server) handleValidateModule(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	m, err := core.GetModule(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	for k, v := range req.Options {
		if err := m.SetOption(k, v); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}

	if err := m.Validate(); err != nil {
		writeJSON(w, http.StatusOK, APIResponse{
			Status: "ok",
			Data:   map[string]interface{}{"valid": false, "error": err.Error()},
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "ok",
		Data:   map[string]interface{}{"valid": true},
	})
}

func (s *Server) handleRunModule(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	m, err := core.GetModule(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	for k, v := range req.Options {
		if err := m.SetOption(k, v); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}

	if err := m.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := m.Run()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "ok",
		Data:   map[string]interface{}{"result": result},
	})
}
