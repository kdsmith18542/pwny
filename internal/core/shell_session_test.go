// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellSession_ReadWrite(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	session, err := newShellSession("test-shell", client)
	require.NoError(t, err)
	defer session.Close()

	// Test Write
	go func() {
		buf := make([]byte, 1024)
		n, _ := server.Read(buf)
		server.Write(buf[:n])
	}()

	msg := []byte("hello shell")
	n, err := session.Write(msg)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)

	// Test Read
	resp, err := session.Read(len(msg))
	assert.NoError(t, err)
	assert.Equal(t, msg, resp)
}

func TestShellSession_Execute(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	session, err := newShellSession("test-shell", client)
	require.NoError(t, err)

	go func() {
		buf := make([]byte, 1024)
		n, _ := server.Read(buf)
		if string(buf[:n]) == "whoami\n" {
			server.Write([]byte("root"))
		}
	}()

	output, err := session.Execute("whoami")
	assert.NoError(t, err)
	assert.Equal(t, "root", output)
}

func TestShellSession_Close(t *testing.T) {
	client, server := net.Pipe()
	session, err := newShellSession("test-shell", client)
	require.NoError(t, err)

	err = session.Close()
	assert.NoError(t, err)
	assert.Equal(t, SessionStatusClosed, session.Info().Status)

	// Verify connection is closed
	_, err = server.Read(make([]byte, 1))
	assert.Error(t, err)
}

func TestShellSession_Unimplemented(t *testing.T) {
	client, _ := net.Pipe()
	session, _ := newShellSession("test-shell", client)

	assert.Error(t, session.Interact())
	assert.Error(t, session.Upload("", ""))
	assert.Error(t, session.Download("", ""))
	
	pid, err := session.GetPID()
	assert.Error(t, err)
	assert.Equal(t, 0, pid)

	uid, err := session.GetUID()
	assert.Error(t, err)
	assert.Equal(t, "", uid)

	sid, err := session.GetSid()
	assert.Error(t, err)
	assert.Equal(t, "", sid)

	admin, err := session.IsAdmin()
	assert.Error(t, err)
	assert.False(t, admin)
}

func TestShellSession_InvalidConn(t *testing.T) {
	_, err := newShellSession("test", "not-a-conn")
	assert.Error(t, err)
}

func TestShellSession_Keepalive(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	// Use a very short interval for testing
	session, _ := newShellSession("test-keepalive", client)
	s := session.(*shellSession)
	s.interval = 10 * time.Millisecond

	// Start reading on server side
	done := make(chan bool)
	go func() {
		buf := make([]byte, 10)
		n, _ := server.Read(buf)
		if n > 0 && buf[0] == '\n' {
			done <- true
		}
	}()

	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Keepalive newline not received")
	}
	
	session.Close()
}
