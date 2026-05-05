// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/msfgo/msfgo/internal/core"
)

type sessionResponse struct {
	ID       string            `json:"id"`
	Type     core.SessionType  `json:"type"`
	Status   core.SessionStatus `json:"status"`
	Target   string            `json:"target"`
	Platform string            `json:"platform"`
	OpenedAt string            `json:"opened_at"`
}

func sessionToResponse(info core.SessionInfo) sessionResponse {
	return sessionResponse{
		ID:       info.ID,
		Type:     info.Type,
		Status:   info.Status,
		Target:   info.Target,
		Platform: info.Platform,
		OpenedAt: info.OpenedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *Server) handleListSessions(w http.ResponseWriter, r *http.Request) {
	infos := sessionManager.ListSessions()
	sessions := make([]sessionResponse, 0, len(infos))
	for _, info := range infos {
		sessions = append(sessions, sessionToResponse(info))
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: sessions})
}

func (s *Server) handleGetSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sess, err := sessionManager.GetSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: sessionToResponse(sess.Info())})
}

func (s *Server) handleCloseSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := sessionManager.CloseSession(id); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: map[string]string{"closed": id}})
}
