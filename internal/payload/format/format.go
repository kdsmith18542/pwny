// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package format

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type PythonFormatter struct{}
func (f *PythonFormatter) Format(data []byte) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("buf = b\"")
	for i, b := range data {
		if i > 0 && i%16 == 0 {
			sb.WriteString("\"\nbuf += b\"")
		}
		sb.WriteString(fmt.Sprintf("\\x%02x", b))
	}
	sb.WriteString("\"")
	return []byte(sb.String()), nil
}

type PowerShellFormatter struct{}
func (f *PowerShellFormatter) Format(data []byte) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("[Byte[]] $buf = ")
	for i, b := range data {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("0x%02x", b))
	}
	return []byte(sb.String()), nil
}

type CFormatter struct{}
func (f *CFormatter) Format(data []byte) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("unsigned char buf[] = { ")
	for i, b := range data {
		if i > 0 {
			sb.WriteString(", ")
		}
		if i > 0 && i%12 == 0 {
			sb.WriteString("\n  ")
		}
		sb.WriteString(fmt.Sprintf("0x%02x", b))
	}
	sb.WriteString(" };")
	return []byte(sb.String()), nil
}

type HexFormatter struct{}
func (f *HexFormatter) Format(data []byte) ([]byte, error) {
	return []byte(hex.EncodeToString(data)), nil
}
