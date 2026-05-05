// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testModule struct {
	*BaseModule
}

func newTestModule() *testModule {
	m := &testModule{BaseModule: NewBaseModule(TypeAuxiliary, "test/module")}
	m.SetDescription("A test module")
	m.SetAuthors([]string{"Test Author"})
	return m
}

type helloModule struct {
	*BaseModule
}

func newHelloModule() *helloModule {
	m := &helloModule{BaseModule: NewBaseModule(TypeAuxiliary, "custom/hello")}
	m.RegisterOption("NAME", "Name to greet", false, "World")
	return m
}

func (m *helloModule) Run() (interface{}, error) {
	name, _ := m.GetOption("NAME")
	return "Hello, " + name.(string) + "!", nil
}

func TestNewBaseModule(t *testing.T) {
	m := NewBaseModule(TypeExploit, "test/exploit")
	assert.Equal(t, TypeExploit, m.info.Type)
	assert.Equal(t, "test/exploit", m.info.Name)
	assert.Equal(t, "No description available", m.info.Description)
	assert.Empty(t, m.options)
}

func TestModuleInfo(t *testing.T) {
	m := newTestModule()
	info := m.Info()
	assert.Equal(t, "test/module", info.Name)
	assert.Equal(t, TypeAuxiliary, info.Type)
	assert.Equal(t, "A test module", info.Description)
	assert.Equal(t, []string{"Test Author"}, info.Authors)
}

func TestRegisterOption(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("RHOST", "Remote host", true, nil)
	m.RegisterOption("RPORT", "Remote port", true, 443)
	m.RegisterOption("VERBOSE", "Verbose output", false, false)

	opts := m.Options()
	assert.Len(t, opts, 3)

	rhost, exists := opts["RHOST"]
	assert.True(t, exists)
	assert.True(t, rhost.Required)
	assert.Nil(t, rhost.Default)
	assert.Nil(t, rhost.Value)
}

func TestSetOption(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("RHOST", "Remote host", true, nil)

	err := m.SetOption("RHOST", "192.168.1.1")
	assert.NoError(t, err)

	val, err := m.GetOption("RHOST")
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", val)
}

func TestSetOptionUnknown(t *testing.T) {
	m := newTestModule()
	err := m.SetOption("NONEXISTENT", "value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown option")
}

func TestGetOptionUnknown(t *testing.T) {
	m := newTestModule()
	_, err := m.GetOption("NONEXISTENT")
	assert.Error(t, err)
}

func TestValidateSuccess(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("RHOST", "Remote host", true, nil)
	m.SetOption("RHOST", "192.168.1.1")
	assert.NoError(t, m.Validate())
}

func TestValidateMissingRequired(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("RHOST", "Remote host", true, nil)
	err := m.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required option")
}

func TestValidateOptionalNotSet(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("VERBOSE", "Verbose output", false, false)
	assert.NoError(t, m.Validate())
}

func TestBaseModuleRun(t *testing.T) {
	m := newTestModule()
	_, err := m.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestCustomModuleRun(t *testing.T) {
	m := newHelloModule()
	result, err := m.Run()
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", result)
}

func TestToJSON(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("RHOST", "Remote host", true, nil)
	data, err := m.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(data), "test/module")
	assert.Contains(t, string(data), "RHOST")
}

func TestSetDescription(t *testing.T) {
	m := newTestModule()
	m.SetDescription("Updated description")
	assert.Equal(t, "Updated description", m.info.Description)
}

func TestSetAuthors(t *testing.T) {
	m := newTestModule()
	m.SetAuthors([]string{"Author1", "Author2"})
	assert.Equal(t, []string{"Author1", "Author2"}, m.info.Authors)
}

func TestOptionDefaultValuePreserved(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("PORT", "Port", false, 8080)
	opt := m.Options()["PORT"]
	assert.Equal(t, 8080, opt.Default)
	assert.Equal(t, 8080, opt.Value)
}

func TestValidateEmptyOptions(t *testing.T) {
	m := newTestModule()
	assert.NoError(t, m.Validate())
}

func TestModuleTypeConstants(t *testing.T) {
	assert.Equal(t, ModuleType("exploit"), TypeExploit)
	assert.Equal(t, ModuleType("payload"), TypePayload)
	assert.Equal(t, ModuleType("auxiliary"), TypeAuxiliary)
	assert.Equal(t, ModuleType("post"), TypePost)
	assert.Equal(t, ModuleType("encoder"), TypeEncoder)
	assert.Equal(t, ModuleType("nop"), TypeNOP)
}

func TestOptionValueChange(t *testing.T) {
	m := newTestModule()
	m.RegisterOption("COUNT", "Count", false, 1)
	m.SetOption("COUNT", 100)

	val, _ := m.GetOption("COUNT")
	assert.Equal(t, 100, val)

	opt := m.Options()["COUNT"]
	assert.Equal(t, 1, opt.Default, "default should be unchanged")
	assert.Equal(t, 100, opt.Value, "value should be updated")
}

func TestMultipleOptionsRegistration(t *testing.T) {
	m := newTestModule()
	names := []string{"OPT1", "OPT2", "OPT3", "OPT4", "OPT5"}
	for i, name := range names {
		m.RegisterOption(name, "Option "+name, i%2 == 0, i)
	}

	opts := m.Options()
	assert.Len(t, opts, 5)

	for i, name := range names {
		opt, exists := opts[name]
		require.True(t, exists, "option %s should exist", name)
		assert.Equal(t, i%2 == 0, opt.Required)
		assert.Equal(t, i, opt.Default)
		assert.Equal(t, i, opt.Value)
	}
}
