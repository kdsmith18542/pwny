// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestRegistry(t *testing.T) func() {
	t.Helper()

	Register("test/existing", func() Module {
		m := NewBaseModule(TypeAuxiliary, "existing")
		m.SetDescription("Existing test module")
		return m
	})

	return func() {
		Unregister("test/existing")
	}
}

func TestRegisterAndGetModule(t *testing.T) {
	cleanup := setupTestRegistry(t)
	defer cleanup()

	m, err := GetModule("test/existing")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, "existing", m.Info().Name)
}

func TestGetModuleNotFound(t *testing.T) {
	_, err := GetModule("nonexistent/module")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetModuleReturnsNewInstance(t *testing.T) {
	cleanup := setupTestRegistry(t)
	defer cleanup()

	m1, _ := GetModule("test/existing")
	m2, _ := GetModule("test/existing")

	assert.NotSame(t, m1, m2, "each GetModule call should return a new instance")
}

func TestRegisterOverwritesExisting(t *testing.T) {
	Register("test/overwrite", func() Module {
		return NewBaseModule(TypeAuxiliary, "original")
	})

	Register("test/overwrite", func() Module {
		return NewBaseModule(TypeExploit, "overwritten")
	})
	defer Unregister("test/overwrite")

	m, _ := GetModule("test/overwrite")
	assert.Equal(t, TypeExploit, m.Info().Type)
	assert.Equal(t, "overwritten", m.Info().Name)
}

func TestListModules(t *testing.T) {
	Register("test/alpha", func() Module { return NewBaseModule(TypeExploit, "alpha") })
	Register("test/beta", func() Module { return NewBaseModule(TypeAuxiliary, "beta") })
	defer func() {
		Unregister("test/alpha")
		Unregister("test/beta")
	}()

	all := ListModules("")
	assert.GreaterOrEqual(t, len(all), 2)

	exploits := ListModules(TypeExploit)
	for _, info := range exploits {
		assert.Equal(t, TypeExploit, info.Type)
	}
}

func TestListModulesFilterType(t *testing.T) {
	Register("test/filter_exploit", func() Module { return NewBaseModule(TypeExploit, "filter_exploit") })
	Register("test/filter_post", func() Module { return NewBaseModule(TypePost, "filter_post") })
	defer func() {
		Unregister("test/filter_exploit")
		Unregister("test/filter_post")
	}()

	postModules := ListModules(TypePost)
	assert.GreaterOrEqual(t, len(postModules), 1)
	for _, info := range postModules {
		assert.Equal(t, TypePost, info.Type)
	}

	payloadModules := ListModules(TypePayload)
	for _, info := range payloadModules {
		assert.Equal(t, TypePayload, info.Type)
	}
}

func TestModuleCount(t *testing.T) {
	before := ModuleCount()

	Register("test/count1", func() Module { return NewBaseModule(TypeAuxiliary, "count1") })
	Register("test/count2", func() Module { return NewBaseModule(TypeAuxiliary, "count2") })

	assert.Equal(t, before+2, ModuleCount())

	Unregister("test/count1")
	assert.Equal(t, before+1, ModuleCount())

	Unregister("test/count2")
	assert.Equal(t, before, ModuleCount())
}

func TestUnregister(t *testing.T) {
	Register("test/unregister", func() Module { return NewBaseModule(TypeAuxiliary, "unregister") })
	_, err := GetModule("test/unregister")
	assert.NoError(t, err)

	Unregister("test/unregister")
	_, err = GetModule("test/unregister")
	assert.Error(t, err)
}

func TestFactoryCreatesConfiguredModule(t *testing.T) {
	Register("test/configured", func() Module {
		m := NewBaseModule(TypeExploit, "configured")
		m.SetDescription("Fully configured")
		m.SetAuthors([]string{"Tester"})
		m.RegisterOption("RHOST", "Remote host", true, nil)
		return m
	})
	defer Unregister("test/configured")

	m, _ := GetModule("test/configured")
	assert.Equal(t, "Fully configured", m.Info().Description)
	assert.Contains(t, m.Options(), "RHOST")
}

func TestListModulesDeduplicatesByKey(t *testing.T) {
	Register("test/dup", func() Module { return NewBaseModule(TypeAuxiliary, "dup") })
	Register("test/dup2", func() Module { return NewBaseModule(TypeAuxiliary, "dup") })
	defer func() {
		Unregister("test/dup")
		Unregister("test/dup2")
	}()

	all := ListModules(TypeAuxiliary)
	count := 0
	for _, info := range all {
		if info.Name == "dup" {
			count++
		}
	}
	assert.Equal(t, 1, count, "duplicate module names should be filtered out")
}

func TestModuleFactoryImmutability(t *testing.T) {
	Register("test/factory_immut", func() Module {
		m := NewBaseModule(TypeAuxiliary, "factory_immut")
		m.RegisterOption("OPT", "An option", false, "initial")
		return m
	})
	defer Unregister("test/factory_immut")

	m1, _ := GetModule("test/factory_immut")
	m2, _ := GetModule("test/factory_immut")

	bm1 := m1.(*BaseModule)
	bm2 := m2.(*BaseModule)
	bm1.SetOption("OPT", "modified")

	v1, _ := bm1.GetOption("OPT")
	v2, _ := bm2.GetOption("OPT")

	assert.Equal(t, "modified", v1)
	assert.Equal(t, "initial", v2, "module instances should be independent")
}
