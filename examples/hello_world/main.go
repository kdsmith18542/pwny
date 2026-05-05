// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package main

import (
	"fmt"

	"github.com/msfgo/msfgo/internal/core"
)

type HelloWorldModule struct {
	*core.BaseModule
}

func NewHelloWorldModule() *HelloWorldModule {
	module := &HelloWorldModule{
		BaseModule: core.NewBaseModule(core.TypeAuxiliary, "hello_world"),
	}

	module.SetDescription("A simple hello world module")
	module.SetAuthors([]string{"MSF-Go Team"})

	module.RegisterOption("NAME", "Name to greet", false, "World")
	module.RegisterOption("EXCLAMATION", "Add an exclamation mark", false, true)

	return module
}

func (m *HelloWorldModule) Run() (interface{}, error) {
	name, _ := m.GetOption("NAME")
	exclamation, _ := m.GetOption("EXCLAMATION")

	greeting := fmt.Sprintf("Hello, %s", name)
	if exclamation.(bool) {
		greeting += "!"
	}

	return greeting, nil
}

func init() {
	core.Register("auxiliary/hello_world", func() core.Module {
		return NewHelloWorldModule()
	})
}

func main() {
	m, _ := core.GetModule("auxiliary/hello_world")
	result, err := m.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}
