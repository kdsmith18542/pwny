// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package network

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// HTTPListener manages an HTTP/HTTPS listener for reverse connections
type HTTPListener struct {
	addr    string
	server  *http.Server
	handler http.Handler
}

// NewHTTPListener creates a new HTTP listener
func NewHTTPListener(addr string, handler http.Handler) *HTTPListener {
	return &HTTPListener{
		addr:    addr,
		handler: handler,
	}
}

// Start starts the HTTP listener
func (h *HTTPListener) Start() error {
	h.server = &http.Server{
		Addr:    h.addr,
		Handler: h.handler,
	}

	slog.Info("HTTP listener started", "addr", h.addr)
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP listener failed", "error", err)
		}
	}()

	return nil
}

// Stop stops the HTTP listener
func (h *HTTPListener) Stop() {
	if h.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := h.server.Shutdown(ctx); err != nil {
			slog.Error("HTTP listener shutdown failed", "error", err)
		}
	}
	slog.Info("HTTP listener stopped", "addr", h.addr)
}
