// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kdsmith18542/pwny/internal/core"
)

type TCPConnect struct {
	*core.BaseModule
}

type scanResult struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	State string `json:"state"`
}

func NewTCPConnect() *TCPConnect {
	m := &TCPConnect{
		BaseModule: core.NewBaseModule(core.TypeAuxiliary, "scanner/tcp_connect"),
	}
	m.SetDescription("TCP Connect port scanner")
	m.SetAuthors([]string{"Pwny Framework"})
	m.RegisterOption("RHOSTS", "Target IP address or hostname", true, "")
	m.RegisterOption("RPORTS", "Ports to scan (e.g. 1-1000 or 22,80,443)", true, "1-1024")
	m.RegisterOption("TIMEOUT", "Connection timeout in seconds", false, 3.0)
	m.RegisterOption("CONCURRENCY", "Number of concurrent connections", false, 50)

	return m
}

func (m *TCPConnect) Run() (interface{}, error) {
	rawHost, _ := m.GetOption("RHOSTS")
	host, ok := rawHost.(string)
	if !ok || host == "" {
		return nil, fmt.Errorf("RHOSTS must be a non-empty string")
	}

	rawPorts, _ := m.GetOption("RPORTS")
	portsStr, ok := rawPorts.(string)
	if !ok || portsStr == "" {
		return nil, fmt.Errorf("RPORTS must be a non-empty string")
	}

	rawTimeout, _ := m.GetOption("TIMEOUT")
	timeoutSec, _ := rawTimeout.(float64)
	if timeoutSec <= 0 {
		timeoutSec = 3.0
	}

	rawConcurrency, _ := m.GetOption("CONCURRENCY")
	concurrency, _ := rawConcurrency.(int)
	if concurrency <= 0 {
		concurrency = 50
	}

	ports, err := parsePorts(portsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RPORTS: %w", err)
	}

	return m.scan(host, ports, time.Duration(timeoutSec)*time.Second, concurrency), nil
}

func (m *TCPConnect) scan(host string, ports []int, timeout time.Duration, concurrency int) []scanResult {
	var mu sync.Mutex
	var results []scanResult
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()

			addr := net.JoinHostPort(host, strconv.Itoa(p))
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				results = append(results, scanResult{Host: host, Port: p, State: "open"})
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return results
}

func parsePorts(input string) ([]int, error) {
	var ports []int
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", rangeParts[1])
			}

			if start > end || start < 1 || end > 65535 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			port, err := strconv.Atoi(part)
			if err != nil || port < 1 || port > 65535 {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

func init() {
	core.Register("auxiliary/scanner/tcp_connect", func() core.Module {
		return NewTCPConnect()
	})
}
