package main

import (
	"fmt"
	"net/url"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	"github.com/NubeIO/flow-framework/utils"
)

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	out.Add(model.ObjTypeAnalogInput)
	out.Add(model.ObjTypeAnalogOutput)
	out.Add(model.ObjTypeAnalogValue)
	out.Add(model.ObjTypeBinaryInput)
	out.Add(model.ObjTypeBinaryOutput)
	out.Add(model.ObjTypeBinaryValue)
	return out
}

// GetDisplay implements public.Displayer
func (i *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	loc := &url.URL{
		Path: i.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	fmt.Println(loc) //can show the ui the custom endpoints

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
	messageURL := plugin.Response{
		Details: m,
	}
	return messageURL
}
