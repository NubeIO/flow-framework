package plugin

// Plugin is the interface every plugin need to implement
type Plugin interface {
	// Enable is called every time a plugin is started. Spawn custom goroutines here for polling, etc.
	// It is always called after ^Set.*Handler$
	Enable() error
	// Disable is called every time a plugin is disabled. Plugins should stop all custom goroutines here.
	Disable() error
}

// Info is returned by the exported plugin function GetPluginInfo() for identification
// plugins are identified by their ModulePath, flow will refuse to load plugins with empty ModulePath
type Info struct {
	Version      string
	Author       string
	Name         string
	Website      string
	Description  string
	License      string
	ModulePath   string
	HasNetwork   bool
	ProtocolType string
}
