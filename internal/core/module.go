package core

import (
	"encoding/json"
	"fmt"
)

// ModuleType represents the type of module
type ModuleType string

const (
	TypeExploit  ModuleType = "exploit"
	TypePayload  ModuleType = "payload"
	TypeAuxiliary ModuleType = "auxiliary"
	TypePost     ModuleType = "post"
	TypeEncoder  ModuleType = "encoder"
	TypeNOP      ModuleType = "nop"
)

// ModuleInfo contains metadata about a module
type ModuleInfo struct {
	Name        string     `json:"name"`
	Type        ModuleType `json:"type"`
	Description string     `json:"description"`
	Authors     []string   `json:"authors"`
	References  []string   `json:"references"`
	Platforms   []string   `json:"platforms"`
	Arch        []string   `json:"arch"`
}

// ModuleOption defines a configurable option for a module
type ModuleOption struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
	Value       interface{} `json:"value"`
}

// Module is the base interface that all modules must implement
type Module interface {
	// Info returns the module's metadata
	Info() ModuleInfo

	// Options returns the module's configurable options
	Options() map[string]ModuleOption

	// SetOption configures a module option
	SetOption(name string, value interface{}) error

	// Validate checks if the module is properly configured
	Validate() error

	// Run executes the module
	Run() (interface{}, error)
}

// BaseModule provides a default implementation of common Module methods
type BaseModule struct {
	info    ModuleInfo
	options map[string]ModuleOption
}

// NewBaseModule creates a new BaseModule with default values
func NewBaseModule(moduleType ModuleType, name string) *BaseModule {
	return &BaseModule{
		info: ModuleInfo{
			Name:        name,
			Type:        moduleType,
			Description: "No description available",
			Authors:     []string{"Unknown"},
			Platforms:   []string{"*"},
			Arch:        []string{"*"},
		},
		options: make(map[string]ModuleOption),
	}
}

// Info returns the module's metadata
func (m *BaseModule) Info() ModuleInfo {
	return m.info
}

// Options returns the module's configurable options
func (m *BaseModule) Options() map[string]ModuleOption {
	return m.options
}

// SetOption configures a module option
func (m *BaseModule) SetOption(name string, value interface{}) error {
	if opt, exists := m.options[name]; exists {
		opt.Value = value
		m.options[name] = opt
		return nil
	}
	return fmt.Errorf("unknown option: %s", name)
}

// Validate checks if the module is properly configured
func (m *BaseModule) Validate() error {
	// Check required options
	for name, opt := range m.options {
		if opt.Required && opt.Value == nil {
			return fmt.Errorf("missing required option: %s", name)
		}
	}
	return nil
}

// Run is a placeholder that should be overridden by actual implementations
func (m *BaseModule) Run() (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// RegisterOption adds a new configurable option to the module
func (m *BaseModule) RegisterOption(name, description string, required bool, defaultValue interface{}) {
	m.options[name] = ModuleOption{
		Name:        name,
		Description: description,
		Required:    required,
		Default:     defaultValue,
		Value:       defaultValue,
	}
}

// GetOption returns the current value of an option
func (m *BaseModule) GetOption(name string) (interface{}, error) {
	if opt, exists := m.options[name]; exists {
		return opt.Value, nil
	}
	return nil, fmt.Errorf("unknown option: %s", name)
}

// ToJSON returns the module's configuration as JSON
func (m *BaseModule) ToJSON() ([]byte, error) {
	return json.MarshalIndent(struct {
		Info    ModuleInfo
		Options map[string]ModuleOption
	}{
		Info:    m.info,
		Options: m.options,
	}, "", "  ")
}
