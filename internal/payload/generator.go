// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package payload

import (
	"fmt"
)

// Platform represents a target operating system
type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
	PlatformMacOS   Platform = "macos"
)

// Arch represents a target architecture
type Arch string

const (
	ArchX64 Arch = "x64"
	ArchX86 Arch = "x86"
	ArchARM Arch = "arm"
)

// PayloadOptions contains configuration for payload generation
type PayloadOptions struct {
	Platform Platform
	Arch     Arch
	LHOST    string
	LPORT    int
	Format   string
	Encoder  string
	Iterations int
	Options    map[string]interface{}
}

// Generator is the interface for generating various payloads
type Generator interface {
	Generate(opts PayloadOptions) ([]byte, error)
}

// Encoder is the interface for encoding shellcode
type Encoder interface {
	Encode(data []byte, opts PayloadOptions) ([]byte, error)
}

// Formatter is the interface for formatting shellcode output
type Formatter interface {
	Format(data []byte) ([]byte, error)
}

// Registry stores payload generators, encoders, and formatters
type Registry struct {
	stagers    map[string]Generator
	stages     map[string]Generator
	encoders   map[string]Encoder
	formatters map[string]Formatter
}

var (
	GlobalRegistry = NewRegistry()
)

func NewRegistry() *Registry {
	return &Registry{
		stagers:    make(map[string]Generator),
		stages:     make(map[string]Generator),
		encoders:   make(map[string]Encoder),
		formatters: make(map[string]Formatter),
	}
}

func (r *Registry) RegisterStager(name string, g Generator) {
	r.stagers[name] = g
}

func (r *Registry) RegisterStage(name string, g Generator) {
	r.stages[name] = g
}

func (r *Registry) RegisterEncoder(name string, e Encoder) {
	r.encoders[name] = e
}

func (r *Registry) RegisterFormatter(name string, f Formatter) {
	r.formatters[name] = f
}

func (r *Registry) GetStager(name string) (Generator, error) {
	g, ok := r.stagers[name]
	if !ok {
		return nil, fmt.Errorf("stager not found: %s", name)
	}
	return g, nil
}

func (r *Registry) GetEncoder(name string) (Encoder, error) {
	e, ok := r.encoders[name]
	if !ok {
		return nil, fmt.Errorf("encoder not found: %s", name)
	}
	return e, nil
}

func (r *Registry) GetFormatter(name string) (Formatter, error) {
	f, ok := r.formatters[name]
	if !ok {
		return nil, fmt.Errorf("formatter not found: %s", name)
	}
	return f, nil
}

// DefaultGenerator coordinates stager/stage generation, encoding and formatting
type DefaultGenerator struct {
	registry *Registry
}

func NewDefaultGenerator(r *Registry) *DefaultGenerator {
	return &DefaultGenerator{registry: r}
}

func (d *DefaultGenerator) Generate(name string, opts PayloadOptions) ([]byte, error) {
	g, err := d.registry.GetStager(name)
	if err != nil {
		return nil, err
	}

	data, err := g.Generate(opts)
	if err != nil {
		return nil, err
	}

	// Apply encoder if requested
	if opts.Encoder != "" {
		encoder, err := d.registry.GetEncoder(opts.Encoder)
		if err != nil {
			return nil, err
		}

		iterations := opts.Iterations
		if iterations <= 0 {
			iterations = 1
		}

		for i := 0; i < iterations; i++ {
			data, err = encoder.Encode(data, opts)
			if err != nil {
				return nil, fmt.Errorf("encoding failed at iteration %d: %w", i+1, err)
			}
		}
	}

	// Apply output formatter if requested
	if opts.Format != "" && opts.Format != "raw" {
		formatter, err := d.registry.GetFormatter(opts.Format)
		if err != nil {
			return nil, err
		}

		data, err = formatter.Format(data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
