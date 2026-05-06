// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package portscan

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/kdsmith18542/pwny/internal/core"
)

type TCPScanner struct {
	*core.BaseModule
}

func NewTCPScanner() *TCPScanner {
	m := &TCPScanner{BaseModule: core.NewBaseModule(core.TypeAuxiliary, "scanner/portscan/tcp")}
	m.SetDescription("TCP Connect Port Scanner")
	m.SetAuthors([]string{"kdsmith18542"})

	m.RegisterOption("RHOSTS", "Target address range", true, nil)
	m.RegisterOption("PORTS", "Target ports (e.g. 21-25,80,443)", true, "21,22,23,25,53,80,110,135,139,443,445,3306,3389,8080")
	m.RegisterOption("THREADS", "Number of concurrent threads", false, 10)
	m.RegisterOption("TIMEOUT", "Connection timeout in milliseconds", false, 1000)

	return m
}

func (m *TCPScanner) Run() (interface{}, error) {
	rhosts, _ := m.GetOption("RHOSTS")
	portsStr, _ := m.GetOption("PORTS")
	threads, _ := m.GetOption("THREADS")
	timeoutMs, _ := m.GetOption("TIMEOUT")

	// Parse ports
	ports := m.parsePorts(portsStr.(string))
	if len(ports) == 0 {
		return nil, fmt.Errorf("no valid ports specified")
	}

	// For simplicity, we assume RHOSTS is a single IP or a comma-separated list
	hosts := strings.Split(rhosts.(string), ",")
	timeout := time.Duration(timeoutMs.(int)) * time.Millisecond

	var wg sync.WaitGroup
	sem := make(chan struct{}, threads.(int))
	results := make(map[string][]int)
	var mu sync.Mutex

	for _, host := range hosts {
		host = strings.TrimSpace(host)
		for _, port := range ports {
			wg.Add(1)
			sem <- struct{}{}
			go func(h string, p int) {
				defer wg.Done()
				defer func() { <-sem }()

				target := fmt.Sprintf("%s:%d", h, p)
				conn, err := net.DialTimeout("tcp", target, timeout)
				if err == nil {
					conn.Close()
					mu.Lock()
					results[h] = append(results[h], p)
					mu.Unlock()
					fmt.Printf("[+] %s:%d - OPEN\n", h, p)
				}
			}(host, port)
		}
	}

	wg.Wait()
	return results, nil
}

func (m *TCPScanner) parsePorts(s string) []int {
	var ports []int
	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			var start, end int
			fmt.Sscanf(part, "%d-%d", &start, &end)
			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		} else {
			var port int
			fmt.Sscanf(part, "%d", &port)
			if port > 0 {
				ports = append(ports, port)
			}
		}
	}
	return ports
}

func init() {
	core.Register("auxiliary/scanner/portscan/tcp", func() core.Module {
		return NewTCPScanner()
	})
}
