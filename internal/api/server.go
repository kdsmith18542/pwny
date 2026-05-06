// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kdsmith18542/pwny/internal/core"
	"github.com/kdsmith18542/pwny/internal/db"
)

type Server struct {
	router   chi.Router
	httpSrv  *http.Server
	cfg      Config
	started  time.Time
	sessions *core.SessionManager
	jobs     *core.JobManager
	events   *core.EventBus
	database *db.Database
}

type Config struct {
	Host    string
	Port    int
	Allowed []string
}

func New(cfg Config, sm *core.SessionManager, jm *core.JobManager, bus *core.EventBus, database *db.Database) *Server {
	r := chi.NewRouter()
	s := &Server{
		router:   r,
		cfg:      cfg,
		started:  time.Now(),
		sessions: sm,
		jobs:     jm,
		events:   bus,
		database: database,
	}
	s.registerRoutes()
	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("api server starting", "addr", addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("api server error: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("api server shutting down")

	if s.database != nil {
		if err := s.database.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}

	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) Router() chi.Router {
	return s.router
}
