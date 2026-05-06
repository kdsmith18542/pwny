package smb

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"pwny/internal/core"
)

type SMBVersionScanner struct {
	*core.BaseModule
}

func init() {
	core.Register("auxiliary/scanner/smb/smb_version", func() core.Module {
		m := &SMBVersionScanner{
			BaseModule: core.NewBaseModule(core.TypeAuxiliary, "scanner/smb/smb_version"),
		}
		m.Name = "SMB Version Scanner"
		m.Description = "Identify SMB version and target OS information."
		m.Author = "Pwny Team"

		m.AddOption("RHOSTS", core.Option{Type: core.TypeString, Required: true, Description: "Target address range or host"})
		m.AddOption("RPORT", core.Option{Type: core.TypeInt, Default: 445, Description: "Target port"})
		m.AddOption("TIMEOUT", core.Option{Type: core.TypeInt, Default: 5, Description: "Timeout in seconds"})

		return m
	})
}

func (m *SMBVersionScanner) Run() (interface{}, error) {
	rhosts, _ := m.GetOption("RHOSTS")
	rport, _ := m.GetOption("RPORT")
	timeoutVal, _ := m.GetOption("TIMEOUT")
	timeout := time.Duration(timeoutVal.(int)) * time.Second

	addr := fmt.Sprintf("%s:%d", rhosts.(string), rport.(int))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	// 1. Negotiate Protocol
	negotiate := []byte{
		0x00, 0x00, 0x00, 0x2f,
		0xff, 0x53, 0x4d, 0x42, 0x72, 0x00, 0x00, 0x00, 0x00, 0x18, 0x53, 0xc8,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x02,
		0x4e, 0x54, 0x20, 0x4c, 0x4d, 0x20, 0x30, 0x2e, 0x31, 0x32, 0x00,
	}
	conn.Write(negotiate)

	reply := make([]byte, 1024)
	n, _ := conn.Read(reply)
	if n < 36 || !bytes.Equal(reply[4:8], []byte("\xffSMB")) {
		return nil, fmt.Errorf("invalid SMB response")
	}

	// 2. Session Setup AndX (Anonymous)
	// This often triggers the server to send its OS and LANMAN strings in the response data
	sessionSetup := []byte{
		0x00, 0x00, 0x00, 0x48,
		0xff, 0x53, 0x4d, 0x42, 0x73, 0x00, 0x00, 0x00, 0x00, 0x18, 0x07, 0xc8,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x40, 0x00, 0x0d, 0xff, 0x00, 0x00,
		0x00, 0xff, 0xff, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00,
		0x00, 0x00, 0x26, 0x00, 0x00, 0x2e, 0x00, 0x57, 0x69, 0x6e, 0x64, 0x6f,
		0x77, 0x73, 0x20, 0x32, 0x30, 0x30, 0x30, 0x00, 0x57, 0x69, 0x6e, 0x64,
		0x6f, 0x77, 0x73, 0x20, 0x32, 0x30, 0x30, 0x30, 0x20, 0x4c, 0x41, 0x4e,
		0x4d, 0x41, 0x4e, 0x00,
	}
	conn.Write(sessionSetup)

	n, _ = conn.Read(reply)
	osInfo := "SMBv1 (Unknown OS)"
	if n > 50 {
		// NativeOS and NativeLanMan are null-terminated strings at the end of the response
		// Let's try to extract them
		data := reply[36:]
		// Skip WordCount and parameters
		if len(data) > 30 {
			parts := bytes.Split(data[30:], []byte{0x00})
			if len(parts) > 0 && len(parts[0]) > 0 {
				osInfo = string(parts[0])
			}
		}
	}

	return map[string]interface{}{
		"address": rhosts.(string),
		"port":    rport.(int),
		"os":      strings.TrimSpace(osInfo),
		"version": "1.0",
	}, nil
}
