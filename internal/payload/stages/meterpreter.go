package stages

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kdsmith18542/pwny/internal/payload"
)

type MeterpreterStage struct{}

func (g *MeterpreterStage) Generate(opts payload.PayloadOptions) ([]byte, error) {
	switch opts.Platform {
	case payload.PlatformWindows:
		return g.generateWindows(opts)
	case payload.PlatformLinux:
		return g.generateLinux(opts)
	default:
		return nil, fmt.Errorf("unsupported platform for meterpreter stage: %s", opts.Platform)
	}
}

func (g *MeterpreterStage) generateWindows(opts payload.PayloadOptions) ([]byte, error) {
	arch := "x86"
	if opts.Arch == payload.ArchX64 {
		arch = "x64"
	}

	filename := fmt.Sprintf("metsrv.%s.dll", arch)
	path := filepath.Join("data", "meterpreter", filename)

	data, err := os.ReadFile(path)
	if err != nil {
		// Fallback to a small placeholder if file doesn't exist
		return []byte(fmt.Sprintf("REPLACE_WITH_ACTUAL_METSRV_DLL_%s", arch)), nil
	}

	return data, nil
}

func (g *MeterpreterStage) generateLinux(opts payload.PayloadOptions) ([]byte, error) {
	arch := "x86"
	if opts.Arch == payload.ArchX64 {
		arch = "x64"
	}

	filename := fmt.Sprintf("metsrv.%s.so", arch)
	path := filepath.Join("data", "meterpreter", filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return []byte(fmt.Sprintf("REPLACE_WITH_ACTUAL_METSRV_SO_%s", arch)), nil
	}

	return data, nil
}
