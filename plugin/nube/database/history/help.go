package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"net/url"
)

// GetDisplay implements public.Displayer
func (i *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	loc := &url.URL{
		Path: i.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	baseURL.Path = i.basePath
	m := plugin.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	return plugin.Response{
		Details: m,
	}
}
