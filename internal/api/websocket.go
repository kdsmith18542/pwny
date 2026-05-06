// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/kdsmith18542/pwny/internal/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`
}

type wsSession struct {
	conn    *websocket.Conn
	mu      sync.Mutex
	session core.Session
}

func (s *Server) registerWSSessionRoutes(r chi.Router) {
	r.Get("/session/{id}/ws", s.handleSessionWS)
	r.Get("/events", s.handleEventStream)
}

func (s *Server) handleSessionWS(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sess, err := s.sessions.GetSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	ws := &wsSession{conn: conn, session: sess}
	slog.Info("websocket session opened", "session_id", id)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ws.readLoop()
	}()

	go func() {
		defer wg.Done()
		ws.writeLoop()
	}()

	wg.Wait()
	slog.Info("websocket session closed", "session_id", id)
}

func (ws *wsSession) readLoop() {
	defer ws.conn.Close()
	for {
		_, msgBytes, err := ws.conn.ReadMessage()
		if err != nil {
			return
		}

		var msg wsMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "session:input":
			ws.session.Write([]byte(msg.Data))
		case "session:execute":
			output, err := ws.session.Execute(msg.Data)
			if err != nil {
				ws.sendJSON(wsMessage{Type: "session:error", Data: err.Error()})
			} else {
				ws.sendJSON(wsMessage{Type: "session:output", Data: output})
			}
		}
	}
}

func (ws *wsSession) writeLoop() {
	defer ws.conn.Close()
	for {
		data, err := ws.session.Read(4096)
		if err != nil {
			return
		}
		ws.sendJSON(wsMessage{
			Type: "session:output",
			Data: string(data),
		})
	}
}

func (ws *wsSession) sendJSON(v interface{}) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.conn.WriteJSON(v)
}

func (s *Server) handleEventStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("event stream upgrade failed", "error", err)
		return
	}

	slog.Info("event stream client connected", "remote", r.RemoteAddr)

	ch := s.events.Subscribe(64)
	defer s.events.Unsubscribe(ch)

	for evt := range ch {
		data, err := json.Marshal(evt)
		if err != nil {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}
