package core

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

// shellSession implements a basic shell session
type shellSession struct {
	*BaseSession
	conn   net.Conn
	reader io.Reader
	writer io.Writer
	mu     sync.Mutex
}

// newShellSession creates a new shell session
func newShellSession(sessionID string, conn interface{}) (Session, error) {
	socket, ok := conn.(net.Conn)
	if !ok {
		return nil, fmt.Errorf("invalid connection type: expected net.Conn")
	}

	session := &shellSession{
		BaseSession: NewBaseSession(SessionTypeShell),
		conn:        socket,
		reader:      socket,
		writer:      socket,
	}

	// Update session info
	session.UpdateInfo(func(info *SessionInfo) {
		info.ID = sessionID
		info.Platform = "unknown"
		info.Arch = "unknown"
	})

	slog.Info("shell session created", "session_id", sessionID)
	go session.keepalive()

	return session, nil
}

// Write sends data to the shell session
func (s *shellSession) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, err := s.writer.Write(data)
	if err != nil {
		slog.Warn("shell write error", "session_id", s.info.ID, "error", err)
		return n, err
	}

	s.UpdateInfo(func(info *SessionInfo) {
		info.LastActivity = time.Now()
	})
	return n, nil
}

// Read receives data from the shell session
func (s *shellSession) Read(length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := s.reader.Read(buf)
	if err != nil {
		return nil, err
	}

	s.UpdateInfo(func(info *SessionInfo) {
		info.LastActivity = time.Now()
	})

	return buf[:n], nil
}

// Close terminates the shell session
func (s *shellSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.UpdateInfo(func(info *SessionInfo) {
		info.Status = SessionStatusClosed
	})

	return s.conn.Close()
}

// Interact starts an interactive shell session
func (s *shellSession) Interact() error {
	// TODO: Implement terminal handling
	return fmt.Errorf("interactive mode not yet implemented")
}

// Execute runs a command and returns the output
func (s *shellSession) Execute(cmd string) (string, error) {
	// Send command with newline
	if _, err := s.Write([]byte(cmd + "\n")); err != nil {
		return "", fmt.Errorf("failed to send command: %v", err)
	}

	// Read response
	buf := make([]byte, 4096)
	n, err := s.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(buf[:n]), nil
}

// Upload uploads a file to the target
func (s *shellSession) Upload(src, dst string) error {
	// TODO: Implement file upload
	return fmt.Errorf("file upload not yet implemented")
}

// Download downloads a file from the target
func (s *shellSession) Download(src, dst string) error {
	// TODO: Implement file download
	return fmt.Errorf("file download not yet implemented")
}

// GetPID gets the process ID of the session
func (s *shellSession) GetPID() (int, error) {
	// TODO: Implement PID detection
	return 0, fmt.Errorf("PID detection not implemented")
}

// GetUID gets the user ID of the session
func (s *shellSession) GetUID() (string, error) {
	// TODO: Implement UID detection
	return "", fmt.Errorf("UID detection not implemented")
}

// GetSid gets the security identifier (Windows only)
func (s *shellSession) GetSid() (string, error) {
	// TODO: Implement SID detection
	return "", fmt.Errorf("SID detection not implemented")
}

// IsAdmin checks if the session has administrative privileges
func (s *shellSession) IsAdmin() (bool, error) {
	// TODO: Implement admin check
	return false, fmt.Errorf("admin check not implemented")
}

// keepalive sends periodic keepalive messages to prevent timeouts
func (s *shellSession) keepalive() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send a newline as keepalive
			s.Write([]byte("\n"))
		}
	}
}
