package core

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"sync"
)

// ModuleLoader handles loading and managing modules
type ModuleLoader struct {
	modules map[string]Module
	mu      sync.RWMutex
}

// NewModuleLoader creates a new ModuleLoader instance
func NewModuleLoader() *ModuleLoader {
	return &ModuleLoader{
		modules: make(map[string]Module),
	}
}

// LoadModule loads a module from a Go plugin (.so file)
func (l *ModuleLoader) LoadModule(path string) (Module, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("module file does not exist: %s", absPath)
	}

	// Try to load the plugin
	plug, err := plugin.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %v", err)
	}

	// Look up the Module symbol
	symModule, err := plug.Lookup("Module")
	if err != nil {
		return nil, fmt.Errorf("module does not export 'Module' symbol: %v", err)
	}

	// Verify the symbol is a Module
	module, ok := symModule.(Module)
	if !ok {
		return nil, fmt.Errorf("exported symbol 'Module' is not of type Module")
	}

	// Register the module
	l.mu.Lock()
	defer l.mu.Unlock()

	info := module.Info()
	moduleID := fmt.Sprintf("%s/%s", info.Type, info.Name)
	
	if _, exists := l.modules[moduleID]; exists {
		return nil, fmt.Errorf("module already loaded: %s", moduleID)
	}

	l.modules[moduleID] = module
	return module, nil
}

// GetModule retrieves a loaded module by type and name
func (l *ModuleLoader) GetModule(moduleType, name string) (Module, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	moduleID := fmt.Sprintf("%s/%s", moduleType, name)
	module, exists := l.modules[moduleID]
	if !exists {
		return nil, fmt.Errorf("module not found: %s", moduleID)
	}

	// Return a new instance of the module to avoid state sharing
	return l.cloneModule(module)
}

// ListModules returns information about all loaded modules
func (l *ModuleLoader) ListModules(filterType string) []ModuleInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var modules []ModuleInfo
	for _, module := range l.modules {
		info := module.Info()
		if filterType == "" || info.Type == ModuleType(filterType) {
			modules = append(modules, info)
		}
	}
	return modules
}

// cloneModule creates a deep copy of a module
func (l *ModuleLoader) cloneModule(original Module) (Module, error) {
	// Create a new instance of the same type
	originalType := reflect.TypeOf(original)
	if originalType.Kind() == reflect.Ptr {
		originalType = originalType.Elem()
	}

	newModule := reflect.New(originalType).Interface().(Module)

	// Copy options
	for name, opt := range original.Options() {
		if err := newModule.SetOption(name, opt.Value); err != nil {
			return nil, fmt.Errorf("failed to copy option %s: %v", name, err)
		}
	}

	return newModule, nil
}

// UnloadModule removes a module from the loader
func (l *ModuleLoader) UnloadModule(moduleType, name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	moduleID := fmt.Sprintf("%s/%s", moduleType, name)
	if _, exists := l.modules[moduleID]; !exists {
		return fmt.Errorf("module not found: %s", moduleID)
	}

	delete(l.modules, moduleID)
	return nil
}
