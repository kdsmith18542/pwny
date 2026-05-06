// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"net/http"

	"github.com/kdsmith18542/pwny/internal/payload"
)

type GeneratePayloadRequest struct {
	Name       string                 `json:"name"`
	Platform   payload.Platform       `json:"platform"`
	Arch       payload.Arch           `json:"arch"`
	LHOST      string                 `json:"lhost"`
	LPORT      int                    `json:"lport"`
	Format     string                 `json:"format"`
	Encoder    string                 `json:"encoder"`
	Iterations int                    `json:"iterations"`
	Options    map[string]interface{} `json:"options"`
}

type GeneratePayloadResponse struct {
	Payload []byte `json:"payload"`
	Size    int    `json:"size"`
}

func (s *Server) handleGeneratePayload(w http.ResponseWriter, r *http.Request) {
	var req GeneratePayloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	opts := payload.PayloadOptions{
		Platform:   req.Platform,
		Arch:       req.Arch,
		LHOST:      req.LHOST,
		LPORT:      req.LPORT,
		Format:     req.Format,
		Encoder:    req.Encoder,
		Iterations: req.Iterations,
		Options:    req.Options,
	}

	data, err := s.payloads.Generate(req.Name, opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "ok",
		Data: GeneratePayloadResponse{
			Payload: data,
			Size:    len(data),
		},
	})
}
