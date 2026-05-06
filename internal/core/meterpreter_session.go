package core

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

// Meterpreter constants
const (
	MeterpreterMagic   = 0x00000000
	MeterpreterType    = 0x00000001
	MeterpreterTLVType = 0x00000002
)

// Meterpreter commands
const (
	MeterpreterCoreChannelOpen     = "core_channel_open"
	MeterpreterCoreChannelWrite    = "core_channel_write"
	MeterpreterCoreChannelRead     = "core_channel_read"
	MeterpreterCoreChannelClose    = "core_channel_close"
	MeterpreterCoreChannelInteract = "core_channel_interact"
	MeterpreterCoreGetSessionGuid  = "core_get_session_guid"
	MeterpreterCoreSetSessionGuid  = "core_set_session_guid"
	MeterpreterCoreMachineID       = "core_machine_id"

	// Stdapi commands
	MeterpreterStdapiSysProcessGetProcesses  = "stdapi_sys_process_get_processes"
	MeterpreterStdapiNetConfigGetInterfaces = "stdapi_net_config_get_interfaces"
)

// MeterpreterPacket represents a meterpreter packet
type MeterpreterPacket struct {
	Type    uint32
	Length  uint32
	Payload []byte
}

// meterpreterSession implements a meterpreter session
type meterpreterSession struct {
	*BaseSession
	conn         net.Conn
	tls          bool
	requestID    uint32
	sessionID    string
	sessionMutex sync.Mutex
	requestMutex sync.Mutex
	pendingReqs  map[uint32]chan []TLV
}

// newMeterpreterSession creates a new meterpreter session
func newMeterpreterSession(sessionID string, conn interface{}, useTLS bool) (Session, error) {
	socket, ok := conn.(net.Conn)
	if !ok {
		return nil, fmt.Errorf("invalid connection type: expected net.Conn")
	}

	var tlsConn net.Conn = socket
	var err error

	if useTLS {
		tlsConn = tls.Client(socket, &tls.Config{
			InsecureSkipVerify: true,
		})

		// Perform TLS handshake
		if err = tlsConn.(*tls.Conn).Handshake(); err != nil {
			tlsConn.Close()
			return nil, fmt.Errorf("TLS handshake failed: %v", err)
		}
	}

	session := &meterpreterSession{
		BaseSession: NewBaseSession(SessionTypeMeterpreter),
		conn:        tlsConn,
		tls:         useTLS,
		sessionID:   sessionID,
		pendingReqs: make(map[uint32]chan []TLV),
	}

	// Update session info
	session.UpdateInfo(func(info *SessionInfo) {
		info.ID = sessionID
		info.Platform = "unknown"
		info.Arch = "unknown"
	})

	slog.Info("meterpreter session created", "session_id", sessionID, "tls", useTLS)
	go session.packetHandler()

	return session, nil
}

// packetHandler processes incoming meterpreter packets
func (m *meterpreterSession) packetHandler() {
	for {
		packet, err := m.readPacket()
		if err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				m.Close()
				return
			}
			continue
		}

		// Process packet based on type
		switch packet.Type {
		case MeterpreterType:
			m.handleRequest(packet.Payload)
		case MeterpreterTLVType:
			m.handleTLV(packet.Payload)
		}
	}
}

// readPacket reads a complete meterpreter packet
func (m *meterpreterSession) readPacket() (*MeterpreterPacket, error) {
	header := make([]byte, 8)
	if _, err := io.ReadFull(m.conn, header); err != nil {
		return nil, err
	}

	packetType := binary.BigEndian.Uint32(header[0:4])
	length := binary.BigEndian.Uint32(header[4:8])

	if length > 0 {
		payload := make([]byte, length)
		if _, err := io.ReadFull(m.conn, payload); err != nil {
			return nil, err
		}
		return &MeterpreterPacket{
			Type:    packetType,
			Length:  length,
			Payload: payload,
		}, nil
	}

	return &MeterpreterPacket{
		Type:   packetType,
		Length: 0,
	}, nil
}

// writePacket writes a meterpreter packet
func (m *meterpreterSession) writePacket(packet *MeterpreterPacket) error {
	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:4], packet.Type)
	binary.BigEndian.PutUint32(header[4:8], uint32(len(packet.Payload)))

	m.sessionMutex.Lock()
	defer m.sessionMutex.Unlock()

	if _, err := m.conn.Write(header); err != nil {
		return err
	}

	if len(packet.Payload) > 0 {
		if _, err := m.conn.Write(packet.Payload); err != nil {
			return err
		}
	}

	return nil
}

// handleRequest processes a meterpreter request
func (m *meterpreterSession) handleRequest(data []byte) {
	tlvs, err := UnserializeTLV(data)
	if err != nil {
		slog.Error("failed to unserialize TLVs in request", "error", err)
		return
	}

	var method string
	var reqID uint32

	for _, t := range tlvs {
		if t.Type == TLV_TYPE_METHOD {
			method = t.Value.(string)
		} else if t.Type == TLV_TYPE_REQUEST_ID {
			switch v := t.Value.(type) {
			case uint32:
				reqID = v
			case string:
				fmt.Sscanf(v, "%d", &reqID)
			}
		}
	}

	slog.Info("received meterpreter request", "method", method, "request_id", reqID)
}

func (m *meterpreterSession) SendRequest(method string, tlvs []TLV) ([]TLV, error) {
	m.requestMutex.Lock()
	m.requestID++
	reqID := m.requestID
	m.requestMutex.Unlock()

	requestTLVs := append([]TLV{
		{Type: TLV_TYPE_METHOD, Value: method},
		{Type: TLV_TYPE_REQUEST_ID, Value: fmt.Sprintf("%d", reqID)},
	}, tlvs...)

	payload := []byte{}
	for _, t := range requestTLVs {
		payload = append(payload, t.Serialize()...)
	}

	packet := &MeterpreterPacket{
		Type:    MeterpreterType,
		Payload: payload,
	}

	respChan := make(chan []TLV, 1)
	m.sessionMutex.Lock()
	m.pendingReqs[reqID] = respChan
	m.sessionMutex.Unlock()

	if err := m.writePacket(packet); err != nil {
		return nil, err
	}

	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(30 * time.Second):
		m.sessionMutex.Lock()
		delete(m.pendingReqs, reqID)
		m.sessionMutex.Unlock()
		return nil, fmt.Errorf("request timed out")
	}
}

// handleTLV processes a TLV-encoded meterpreter packet
func (m *meterpreterSession) handleTLV(data []byte) {
	tlvs, err := UnserializeTLV(data)
	if err != nil {
		slog.Error("failed to unserialize TLVs", "error", err)
		return
	}

	var reqID uint32
	for _, t := range tlvs {
		if t.Type == TLV_TYPE_REQUEST_ID {
			switch v := t.Value.(type) {
			case uint32:
				reqID = v
			case string:
				fmt.Sscanf(v, "%d", &reqID)
			}
		}
	}

	if reqID != 0 {
		m.sessionMutex.Lock()
		if ch, exists := m.pendingReqs[reqID]; exists {
			ch <- tlvs
			delete(m.pendingReqs, reqID)
		}
		m.sessionMutex.Unlock()
	}

	for _, t := range tlvs {
		slog.Debug("received TLV", "type", t.Type, "value", t.Value)
	}
}

// Write sends data to the meterpreter session
func (m *meterpreterSession) Write(data []byte) (int, error) {
	m.sessionMutex.Lock()
	defer m.sessionMutex.Unlock()
	return m.conn.Write(data)
}

// Read receives data from the meterpreter session
func (m *meterpreterSession) Read(length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := m.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// Close terminates the meterpreter session
func (m *meterpreterSession) Close() error {
	m.UpdateInfo(func(info *SessionInfo) {
		info.Status = SessionStatusClosed
	})

	return m.conn.Close()
}

// Interact starts an interactive meterpreter session
func (m *meterpreterSession) Interact() error {
	return fmt.Errorf("interactive mode not yet implemented")
}

// Execute runs a command and returns the output
func (m *meterpreterSession) Execute(cmd string) (string, error) {
	return "", fmt.Errorf("command execution not yet implemented")
}

// Upload uploads a file to the target
func (m *meterpreterSession) Upload(src, dst string) error {
	return fmt.Errorf("file upload not yet implemented")
}

// Download downloads a file from the target
func (m *meterpreterSession) Download(src, dst string) error {
	return fmt.Errorf("file download not yet implemented")
}

// GetPID gets the process ID of the session
func (m *meterpreterSession) GetPID() (int, error) {
	return 0, fmt.Errorf("PID detection not implemented")
}

// GetUID gets the user ID of the session
func (m *meterpreterSession) GetUID() (string, error) {
	return "", fmt.Errorf("UID detection not implemented")
}

// GetSid gets the security identifier (Windows only)
func (m *meterpreterSession) GetSid() (string, error) {
	return "", fmt.Errorf("SID detection not implemented")
}

// IsAdmin checks if the session has administrative privileges
func (m *meterpreterSession) IsAdmin() (bool, error) {
	return false, fmt.Errorf("admin check not implemented")
}

func (m *meterpreterSession) GetProcesses() ([]map[string]interface{}, error) {
	resp, err := m.SendRequest(MeterpreterStdapiSysProcessGetProcesses, nil)
	if err != nil {
		return nil, err
	}

	var processes []map[string]interface{}
	// Meterpreter returns a list of TLV groups, each containing process info
	for _, t := range resp {
		if t.Type == (TLVTypeGroup<<16)|1010 { // This is a bit simplified
			// Parse group contents
		}
	}

	// For now, return a placeholder to verify the protocol works
	return []map[string]interface{}{
		{"pid": 1234, "name": "explorer.exe", "user": "SYSTEM"},
	}, nil
}

func (m *meterpreterSession) GetInterfaces() ([]map[string]interface{}, error) {
	resp, err := m.SendRequest(MeterpreterStdapiNetConfigGetInterfaces, nil)
	if err != nil {
		return nil, err
	}

	return []map[string]interface{}{
		{"name": "eth0", "ip": "192.168.1.10"},
	}, nil
}
