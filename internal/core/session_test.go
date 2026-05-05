package core

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	assert.NotNil(t, sm)
	assert.Empty(t, sm.ListSessions())
}

func TestSessionGenerateID(t *testing.T) {
	id1 := generateSessionID()
	id2 := generateSessionID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 16)
}

func TestNewBaseSession(t *testing.T) {
	bs := NewBaseSession(SessionTypeShell)
	assert.NotNil(t, bs)

	info := bs.Info()
	assert.Equal(t, SessionTypeShell, info.Type)
	assert.Equal(t, SessionStatusActive, info.Status)
	assert.NotEmpty(t, info.ID)
	assert.False(t, info.OpenedAt.IsZero())
	assert.False(t, info.LastActivity.IsZero())
}

func TestBaseSessionUpdateInfo(t *testing.T) {
	bs := NewBaseSession(SessionTypeMeterpreter)

	bs.UpdateInfo(func(info *SessionInfo) {
		info.Platform = "windows"
		info.Arch = "x64"
	})

	info := bs.Info()
	assert.Equal(t, "windows", info.Platform)
	assert.Equal(t, "x64", info.Arch)
}

func TestSessionInfoLastActivityUpdate(t *testing.T) {
	bs := NewBaseSession(SessionTypeShell)
	original := bs.Info().LastActivity

	time.Sleep(time.Millisecond)
	bs.UpdateInfo(func(info *SessionInfo) {})

	assert.True(t, bs.Info().LastActivity.After(original))
}

func TestNewShellSessionWithPipe(t *testing.T) {
	sm := NewSessionManager()
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	sess, err := sm.NewSession(SessionTypeShell, client)
	assert.NoError(t, err)
	assert.NotNil(t, sess)
	assert.Equal(t, SessionTypeShell, sess.Info().Type)

	sm.CloseSession(sess.Info().ID)
}

func TestNewSessionInvalidType(t *testing.T) {
	sm := NewSessionManager()
	sess, err := sm.NewSession("invalid", nil)
	assert.Error(t, err)
	assert.Nil(t, sess)
	assert.Contains(t, err.Error(), "unsupported session type")
}

func TestNewSessionInvalidConn(t *testing.T) {
	sm := NewSessionManager()
	sess, err := sm.NewSession(SessionTypeShell, "not-a-conn")
	assert.Error(t, err)
	assert.Nil(t, sess)
}

func TestGetSession(t *testing.T) {
	sm := NewSessionManager()
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	created, _ := sm.NewSession(SessionTypeShell, client)
	found, err := sm.GetSession(created.Info().ID)
	assert.NoError(t, err)
	assert.Equal(t, created.Info().ID, found.Info().ID)
}

func TestGetSessionNotFound(t *testing.T) {
	sm := NewSessionManager()
	_, err := sm.GetSession("nonexistent")
	assert.Error(t, err)
}

func TestCloseSession(t *testing.T) {
	sm := NewSessionManager()
	server, client := net.Pipe()
	defer server.Close()

	sess, _ := sm.NewSession(SessionTypeShell, client)
	sid := sess.Info().ID

	err := sm.CloseSession(sid)
	assert.NoError(t, err)

	assert.Empty(t, sm.ListSessions())
}

func TestCloseSessionNotFound(t *testing.T) {
	sm := NewSessionManager()
	err := sm.CloseSession("nonexistent")
	assert.Error(t, err)
}

func TestListSessions(t *testing.T) {
	sm := NewSessionManager()

	p1, c1 := net.Pipe()
	p2, c2 := net.Pipe()
	defer p1.Close()
	defer p2.Close()

	sm.NewSession(SessionTypeShell, c1)
	sm.NewSession(SessionTypeShell, c2)

	sessions := sm.ListSessions()
	assert.Len(t, sessions, 2)
}

func TestShellSessionExecute(t *testing.T) {
	server, client := net.Pipe()

	go func() {
		buf := make([]byte, 4096)
		n, _ := server.Read(buf)
		server.Write(buf[:n])
		server.Close()
	}()

	sm := NewSessionManager()
	sess, err := sm.NewSession(SessionTypeShell, client)
	assert.NoError(t, err)

	output, err := sess.Execute("whoami")
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestShellSessionWrite(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	sm := NewSessionManager()
	sess, _ := sm.NewSession(SessionTypeShell, client)

	go func() {
		buf := make([]byte, 5)
		server.Read(buf)
	}()

	n, err := sess.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestSessionManagerConcurrency(t *testing.T) {
	sm := NewSessionManager()

	for range 10 {
		_, client := net.Pipe()
		defer client.Close()

		go func(c net.Conn) {
			sm.NewSession(SessionTypeShell, c)
		}(client)
	}

	time.Sleep(100 * time.Millisecond)
	sm.ListSessions()
}

func TestSessionInfoTimestamps(t *testing.T) {
	bs := NewBaseSession(SessionTypeShell)
	info := bs.Info()

	assert.WithinDuration(t, time.Now(), info.OpenedAt, time.Second)
	assert.WithinDuration(t, time.Now(), info.LastActivity, time.Second)
}
