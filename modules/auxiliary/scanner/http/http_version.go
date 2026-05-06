package http

import (
	"fmt"
	"net/http"
	"time"

	"pwny/internal/core"
)

type HTTPVersionScanner struct {
	*core.BaseModule
}

func init() {
	core.Register("auxiliary/scanner/http/http_version", func() core.Module {
		m := &HTTPVersionScanner{
			BaseModule: core.NewBaseModule(core.TypeAuxiliary, "scanner/http/http_version"),
		}
		m.Name = "HTTP Version Scanner"
		m.Description = "Identify HTTP server version and technologies."
		m.Author = "Pwny Team"

		m.AddOption("RHOSTS", core.Option{Type: core.TypeString, Required: true, Description: "Target address range or host"})
		m.AddOption("RPORT", core.Option{Type: core.TypeInt, Default: 80, Description: "Target port"})
		m.AddOption("SSL", core.Option{Type: core.TypeBool, Default: false, Description: "Use HTTPS"})
		m.AddOption("TIMEOUT", core.Option{Type: core.TypeInt, Default: 5, Description: "Timeout in seconds"})

		return m
	})
}

func (m *HTTPVersionScanner) Run() (interface{}, error) {
	rhosts, _ := m.GetOption("RHOSTS")
	rport, _ := m.GetOption("RPORT")
	ssl, _ := m.GetOption("SSL")
	timeoutVal, _ := m.GetOption("TIMEOUT")
	timeout := time.Duration(timeoutVal.(int)) * time.Second

	protocol := "http"
	if ssl.(bool) {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/", protocol, rhosts.(string), rport.(int))
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Head(url)
	if err != nil {
		// Try GET if HEAD fails
		resp, err = client.Get(url)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	server := resp.Header.Get("Server")
	if server == "" {
		server = "Unknown"
	}

	poweredBy := resp.Header.Get("X-Powered-By")

	return map[string]interface{}{
		"address":    rhosts.(string),
		"port":       rport.(int),
		"server":     server,
		"powered_by": poweredBy,
	}, nil
}
