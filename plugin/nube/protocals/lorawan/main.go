package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwrest"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
)

const path = "lorawan" // must be unique across all plugins
const name = "lorawan" // must be unique across all plugins
const description = "lorawan api"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "lorawan"
const transportType = "ip" // serial, ip
const ip = "0.0.0.0"
const port = "8080"

// Instance is plugin instance
type Instance struct {
	config      *Config
	enabled     bool
	basePath    string
	db          dbhandler.Handler
	store       cachestore.Handler
	bus         eventbus.BusService
	pluginUUID  string
	networkUUID string
	REST        *lwrest.RestClient
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   path,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
		HasNetwork:   true,
		ProtocolType: protocolType,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance() pluginapi.Plugin {
	return &Instance{}
}

// main will not let main run
func main() {
	panic("this should be built as plugin")
}
