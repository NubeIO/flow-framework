package main

import (
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"net/url"
)

// GetDisplay implements public.Displayer
func (inst *Instance) GetDisplay(baseURL *url.URL) pluginapi.Response {
	loc := &url.URL{
		Path: inst.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	baseURL.Path = inst.basePath
	m := pluginapi.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	return pluginapi.Response{
		Details: m,
	}
}
