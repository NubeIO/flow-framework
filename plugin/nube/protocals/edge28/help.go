package main

import (
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	"net/url"
)

// GetDisplay implements public.Displayer
func (inst *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	loc := &url.URL{
		Path: inst.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})

	baseURL.Path = inst.basePath
	m := plugin.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	messageURL := plugin.Response{
		Details: m,
	}
	return messageURL
}
