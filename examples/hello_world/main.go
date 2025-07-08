package main

import (
	"fmt"

	"github.com/msfgo/msfgo/internal/core"
)

// HelloWorldModule is a simple example module
type HelloWorldModule struct {
	*core.BaseModule
}

// NewHelloWorldModule creates a new HelloWorldModule instance
func NewHelloWorldModule() *HelloWorldModule {
	module := &HelloWorldModule{
		BaseModule: core.NewBaseModule(core.TypeAuxiliary, "hello_world"),
	}

	// Set module info
	module.Info().Description = "A simple hello world module"
	module.Info().Authors = []string{"MSF-Go Team"}

	// Register options
	module.RegisterOption("NAME", "Name to greet", false, "World")
	module.RegisterOption("EXCLAMATION", "Add an exclamation mark", false, true)

	return module
}

// Run implements the Module interface
func (m *HelloWorldModule) Run() (interface{}, error) {
	// Get options
	name, _ := m.GetOption("NAME")
	exclamation, _ := m.GetOption("EXCLAMATION")

	// Generate greeting
	greeting := fmt.Sprintf("Hello, %s", name)
	if exclamation.(bool) {
		greeting += "!"
	}

	return greeting, nil
}

// Module is the exported symbol that the loader looks for
var Module core.Module = NewHelloWorldModule()
