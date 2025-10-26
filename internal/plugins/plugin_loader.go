package plugins

import (
	"plugin"
	"path/filepath"
)

// PluginLoader handles dynamic loading of plugins
type PluginLoader struct {
	pluginDir string
}

// NewPluginLoader creates a new plugin loader instance
func NewPluginLoader(dir string) *PluginLoader {
	return &PluginLoader{
		pluginDir: dir,
	}
}

// LoadPlugin loads a plugin from the specified path
func (p *PluginLoader) LoadPlugin(name string) (plugin.Symbol, error) {
	path := filepath.Join(p.pluginDir, name+".so")
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	
	symbol, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, err
	}
	
	return symbol, nil
}
