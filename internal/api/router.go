package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIError struct {
	Error  string `json:"error,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *APIError   `json:"error,omitempty"`
}

func (s *Server) registerRoutes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(loggingMiddleware)
	s.router.Use(recoveryMiddleware)
	s.router.Use(corsMiddleware(s.cfg.Allowed))

	s.router.Get("/api/v1/status", s.handleStatus)

	s.router.Route("/api/v1/modules", func(r chi.Router) {
		r.Get("/", s.handleListModules)
		r.Get("/{name}", s.handleGetModule)
		r.Post("/{name}/validate", s.handleValidateModule)
		r.Post("/{name}/run", s.handleRunModule)
	})

	s.router.Route("/api/v1/sessions", func(r chi.Router) {
		r.Get("/", s.handleListSessions)
		r.Get("/{id}", s.handleGetSession)
		r.Delete("/{id}", s.handleCloseSession)
		r.Get("/{id}/ws", s.handleSessionWS)
	})

	s.router.Get("/api/v1/events", s.handleEventStream)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		Status: "ok",
		Data: map[string]interface{}{
			"version":   "0.1.0",
			"uptime":    time.Since(s.started).String(),
			"started":   s.started.Format(time.RFC3339),
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, APIResponse{
		Status: "error",
		Error:  &APIError{Error: err.Error()},
	})
}

func corsMiddleware(origins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range origins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func moduleKey(mType, mName string) string {
	return fmt.Sprintf("%s/%s", mType, mName)
}
