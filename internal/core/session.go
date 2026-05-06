// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// SessionType represents the type of session
type SessionType string

const (
	SessionTypeMeterpreter    SessionType = "meterpreter"
	SessionTypeShell          SessionType = "shell"
	SessionTypeMeterpreterSSL SessionType = "meterpreter_ssl"
	SessionTypeHTTP           SessionType = "http"
	SessionTypeHTTPS          SessionType = "https"
)

// SessionStatus represents the status of a session
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusClosed    SessionStatus = "closed"
	SessionStatusError     SessionStatus = "error"
	SessionStatusUpgrading SessionStatus = "upgrading"
)

// SessionInfo contains metadata about a session
type SessionInfo struct {
	ID           string        `json:"id"`
	Type         SessionType   `json:"type"`
	Status       SessionStatus `json:"status"`
	ViaExploit   string        `json:"via_exploit"`
	ViaPayload   string        `json:"via_payload"`
	Target       string        `json:"target"`
	Username     string        `json:"username"`
	UUID         string        `json:"uuid"`
	Platform     string        `json:"platform"`
	Arch         string        `json:"arch"`
	Workspace    string        `json:"workspace"`
	OpenedAt     time.Time     `json:"opened_at"`
	LastActivity time.Time     `json:"last_activity"`
}

// Session represents an active session with a target
type Session interface {
	// Info returns session metadata
	Info() SessionInfo

	// Write sends data to the session
	Write(data []byte) (int, error)

	// Read receives data from the session
	Read(length int) ([]byte, error)

	// Close terminates the session
	Close() error

	// Interact starts an interactive session
	Interact() error

	// Execute runs a command and returns the output
	Execute(cmd string) (string, error)

	// Upload uploads a file to the target
	Upload(src, dst string) error

	// Download downloads a file from the target
	Download(src, dst string) error

	// GetPID gets the process ID of the session
	GetPID() (int, error)

	// GetUID gets the user ID of the session
	GetUID() (string, error)

	// GetSid gets the security identifier (Windows only)
	GetSid() (string, error)

	// IsAdmin checks if the session has administrative privileges
	IsAdmin() (bool, error)

	// GetProcesses returns a list of running processes
	GetProcesses() ([]map[string]interface{}, error)

	// GetInterfaces returns a list of network interfaces
	GetInterfaces() ([]map[string]interface{}, error)
}

// SessionManager handles session creation and management
type SessionManager struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
	}
}

// NewSession creates a new session with the given type and connection
func (sm *SessionManager) NewSession(sessionType SessionType, conn interface{}) (Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := generateSessionID()

	var session Session
	var err error

	switch sessionType {
	case SessionTypeShell:
		session, err = newShellSession(sessionID, conn)
	case SessionTypeMeterpreter, SessionTypeMeterpreterSSL:
		session, err = newMeterpreterSession(sessionID, conn, sessionType == SessionTypeMeterpreterSSL)
	default:
		return nil, fmt.Errorf("unsupported session type: %s", sessionType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	sm.sessions[sessionID] = session
	slog.Info("session created", "session_id", sessionID, "type", sessionType)
	return session, nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session, nil
}

// ListSessions returns all active sessions
func (sm *SessionManager) ListSessions() []SessionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var sessions []SessionInfo
	for _, s := range sm.sessions {
		sessions = append(sessions, s.Info())
	}

	return sessions
}

// CloseSession terminates a session
func (sm *SessionManager) CloseSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	if err := session.Close(); err != nil {
		slog.Error("error closing session", "session_id", sessionID, "error", err)
		return fmt.Errorf("error closing session: %v", err)
	}

	delete(sm.sessions, sessionID)
	slog.Info("session closed", "session_id", sessionID)
	return nil
}

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		slog.Error("failed to generate session ID", "error", err)
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(buf)
}

// BaseSession provides common session functionality
type BaseSession struct {
	info SessionInfo
	mu   sync.RWMutex
}

// Info returns session metadata
func (s *BaseSession) Info() SessionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.info
}

// UpdateInfo updates session metadata
func (s *BaseSession) UpdateInfo(updater func(*SessionInfo)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	updater(&s.info)
	s.info.LastActivity = time.Now()
}

// NewBaseSession creates a new base session
func NewBaseSession(sessionType SessionType) *BaseSession {
	now := time.Now()
	return &BaseSession{
		info: SessionInfo{
			ID:           generateSessionID(),
			Type:         sessionType,
			Status:       SessionStatusActive,
			OpenedAt:     now,
			LastActivity: now,
		},
	}
}
