package ssh

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"pwny/internal/core"
)

type SSHVersionScanner struct {
	*core.BaseModule
}

func init() {
	core.Register("auxiliary/scanner/ssh/ssh_version", func() core.Module {
		m := &SSHVersionScanner{
			BaseModule: core.NewBaseModule(core.TypeAuxiliary, "scanner/ssh/ssh_version"),
		}
		m.Name = "SSH Version Scanner"
		m.Description = "Identify SSH server version and banner."
		m.Author = "Pwny Team"

		m.AddOption("RHOSTS", core.Option{Type: core.TypeString, Required: true, Description: "Target address range or host"})
		m.AddOption("RPORT", core.Option{Type: core.TypeInt, Default: 22, Description: "Target port"})
		m.AddOption("TIMEOUT", core.Option{Type: core.TypeInt, Default: 5, Description: "Timeout in seconds"})

		return m
	})
}

func (m *SSHVersionScanner) Run() (interface{}, error) {
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

	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read banner: %v", err)
	}

	banner = strings.TrimSpace(banner)
	return map[string]interface{}{
		"address": rhosts.(string),
		"port":    rport.(int),
		"banner":  banner,
	}, nil
}
