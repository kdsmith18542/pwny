package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/msfgo/msfgo/internal/core"
)

func main() {
	// Create a session manager
	sessionMgr := core.NewSessionManager()

	// Example 1: Create a shell session
	shellConn, _ := net.Dial("tcp", "example.com:4444") // Replace with actual target
	shellSession, err := sessionMgr.NewSession(core.SessionTypeShell, shellConn)
	if err != nil {
		log.Fatalf("Failed to create shell session: %v", err)
	}

	// Execute a command
	output, err := shellSession.Execute("whoami")
	if err != nil {
		log.Printf("Command failed: %v", err)
	} else {
		fmt.Printf("Command output: %s\n", output)
	}

	// Example 2: Create a meterpreter session
	meterpreterConn, _ := net.Dial("tcp", "example.com:5555") // Replace with actual target

	meterpreterSession, err := sessionMgr.NewSession(core.SessionTypeMeterpreter, meterpreterConn)
	if err != nil {
		log.Fatalf("Failed to create meterpreter session: %v", err)
	}

	// List all active sessions
	fmt.Println("\nActive sessions:")
	for _, info := range sessionMgr.ListSessions() {
		fmt.Printf("- %s (%s) - %s\n", info.ID, info.Type, info.Status)
	}

	// Keep the program running
	select {}
}

// Example server for testing
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

// Example shell handler for testing
func exampleShellHandler(conn net.Conn) {
	defer conn.Close()
	
	// Simple echo server for testing
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		
		// Echo back the command
		conn.Write(buf[:n])
	}
}
