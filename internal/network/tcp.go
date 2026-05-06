// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package network

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
)

// Handler handles incoming connections
type Handler func(net.Conn)

// TCPListener manages a TCP listener for reverse connections
type TCPListener struct {
	addr     string
	listener net.Listener
	handler  Handler
	quit     chan struct{}
	wg       sync.WaitGroup
}

// NewTCPListener creates a new TCP listener
func NewTCPListener(addr string, handler Handler) *TCPListener {
	return &TCPListener{
		addr:    addr,
		handler: handler,
		quit:    make(chan struct{}),
	}
}

// Start starts the TCP listener
func (t *TCPListener) Start() error {
	l, err := net.Listen("tcp", t.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", t.addr, err)
	}
	t.listener = l
	slog.Info("TCP listener started", "addr", t.addr)

	t.wg.Add(1)
	go t.acceptLoop()

	return nil
}

// Stop stops the TCP listener
func (t *TCPListener) Stop() {
	close(t.quit)
	if t.listener != nil {
		t.listener.Close()
	}
	t.wg.Wait()
	slog.Info("TCP listener stopped", "addr", t.addr)
}

func (t *TCPListener) acceptLoop() {
	defer t.wg.Done()
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.quit:
				return
			default:
				slog.Error("failed to accept connection", "error", err)
				continue
			}
		}

		slog.Info("new connection accepted", "remote", conn.RemoteAddr())
		go t.handler(conn)
	}
}
