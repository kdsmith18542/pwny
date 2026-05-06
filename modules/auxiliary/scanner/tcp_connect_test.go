// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package scanner

import (
	"testing"

	"github.com/kdsmith18542/pwny/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTCPConnect(t *testing.T) {
	m := NewTCPConnect()
	assert.NotNil(t, m)

	info := m.Info()
	assert.Equal(t, "scanner/tcp_connect", info.Name)
	assert.Equal(t, core.TypeAuxiliary, info.Type)
	assert.NotEmpty(t, info.Description)
	assert.NotEmpty(t, info.Authors)
}

func TestTCPConnectOptions(t *testing.T) {
	m := NewTCPConnect()

	opts := m.Options()
	assert.Contains(t, opts, "RHOSTS")
	assert.Contains(t, opts, "RPORTS")
	assert.Contains(t, opts, "TIMEOUT")
	assert.Contains(t, opts, "CONCURRENCY")

	assert.True(t, opts["RHOSTS"].Required)
	assert.True(t, opts["RPORTS"].Required)
	assert.False(t, opts["TIMEOUT"].Required)
	assert.False(t, opts["CONCURRENCY"].Required)
}

func TestTCPConnectDefaults(t *testing.T) {
	m := NewTCPConnect()

	rhosts, err := m.GetOption("RHOSTS")
	require.NoError(t, err)
	assert.Equal(t, "", rhosts)

	rports, err := m.GetOption("RPORTS")
	require.NoError(t, err)
	assert.Equal(t, "1-1024", rports)

	timeout, err := m.GetOption("TIMEOUT")
	require.NoError(t, err)
	assert.Equal(t, 3.0, timeout)

	concurrency, err := m.GetOption("CONCURRENCY")
	require.NoError(t, err)
	assert.Equal(t, 50, concurrency)
}

func TestTCPConnectSetOptions(t *testing.T) {
	m := NewTCPConnect()

	err := m.SetOption("RHOSTS", "192.168.1.1")
	assert.NoError(t, err)

	err = m.SetOption("RPORTS", "80,443")
	assert.NoError(t, err)

	v, _ := m.GetOption("RHOSTS")
	assert.Equal(t, "192.168.1.1", v)
}

func TestTCPConnectValidateSuccess(t *testing.T) {
	m := NewTCPConnect()
	m.SetOption("RHOSTS", "10.0.0.1")
	m.SetOption("RPORTS", "80")
	err := m.Validate()
	assert.NoError(t, err)
}

func TestTCPConnectValidateMissingRequired(t *testing.T) {
	m := NewTCPConnect()
	_, err := m.GetOption("RHOSTS")
	require.NoError(t, err)

	opts := m.Options()
	opt := opts["RHOSTS"]
	assert.True(t, opt.Required)

	m.SetOption("RHOSTS", nil)
	err = m.Validate()
	assert.Error(t, err)
}

func TestTCPConnectRunNoTarget(t *testing.T) {
	m := NewTCPConnect()
	_, err := m.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RHOSTS")
}

func TestParsePortsSingle(t *testing.T) {
	ports, err := parsePorts("80")
	assert.NoError(t, err)
	assert.Equal(t, []int{80}, ports)
}

func TestParsePortsRange(t *testing.T) {
	ports, err := parsePorts("1-5")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, ports)
}

func TestParsePortsMixed(t *testing.T) {
	ports, err := parsePorts("22,80-82,443")
	assert.NoError(t, err)
	assert.Equal(t, []int{22, 80, 81, 82, 443}, ports)
}

func TestParsePortsInvalid(t *testing.T) {
	_, err := parsePorts("invalid")
	assert.Error(t, err)
}

func TestParsePortsInvalidRange(t *testing.T) {
	_, err := parsePorts("100-50")
	assert.Error(t, err)
}

func TestParsePortsOutOfRange(t *testing.T) {
	_, err := parsePorts("0")
	assert.Error(t, err)

	_, err = parsePorts("70000")
	assert.Error(t, err)
}
