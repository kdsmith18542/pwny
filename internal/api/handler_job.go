// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kdsmith18542/pwny/internal/core"
)

type jobResponse struct {
	ID          string                 `json:"id"`
	ModuleName  string                 `json:"module_name"`
	Status      core.JobStatus         `json:"status"`
	Options     map[string]interface{} `json:"options"`
	Result      interface{}            `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

func jobToResponse(j *core.Job) jobResponse {
	var errStr string
	if j.Error != "" {
		errStr = j.Error
	}
	return jobResponse{
		ID:          j.ID,
		ModuleName:  j.ModuleName,
		Status:      j.Status,
		Options:     j.Options,
		Result:      j.Result,
		Error:       errStr,
		CreatedAt:   j.CreatedAt,
		StartedAt:   j.StartedAt,
		CompletedAt: j.CompletedAt,
	}
}

func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	jobs := s.jobs.List()
	resp := make([]jobResponse, 0, len(jobs))
	for _, j := range jobs {
		resp = append(resp, jobToResponse(j))
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: resp})
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	j, err := s.jobs.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: jobToResponse(j)})
}

func (s *Server) handleCancelJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.jobs.Cancel(id); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"cancelled": id}})
}
