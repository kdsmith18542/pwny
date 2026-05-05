package core

import (
	"fmt"
	"log/slog"
	"sync"
)

type ModuleFactory func() Module

var (
	registry   = make(map[string]ModuleFactory)
	registryMu sync.RWMutex
)

func Register(name string, factory ModuleFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[name]; exists {
		slog.Warn("module already registered, overwriting", "module", name)
	}
	registry[name] = factory
	slog.Debug("module registered", "module", name)
}

func GetModule(name string) (Module, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("module not found: %s", name)
	}
	return factory(), nil
}

func ListModules(filterType ModuleType) []ModuleInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()
	var infos []ModuleInfo
	seen := make(map[string]bool)
	for _, factory := range registry {
		m := factory()
		info := m.Info()
		key := fmt.Sprintf("%s/%s", info.Type, info.Name)
		if seen[key] {
			continue
		}
		seen[key] = true
		if filterType == "" || info.Type == filterType {
			infos = append(infos, info)
		}
	}
	return infos
}

func Unregister(name string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	delete(registry, name)
	slog.Debug("module unregistered", "module", name)
}

func ModuleCount() int {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return len(registry)
}
