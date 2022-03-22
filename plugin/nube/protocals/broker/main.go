package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
)

const name = "broker" //must be unique across all plugins
const description = "broker"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "broker"

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = false
const maxAllowedNetworks = 1
const networkType = "mqtt"
const transportType = "ip" //serial, ip

// Instance is plugin instance
type Instance struct {
	config        *Config
	enabled       bool
	basePath      string
	db            dbhandler.Handler
	store         cachestore.Handler
	bus           eventbus.BusService
	pluginUUID    string
	networkUUID   string
	brokerEnabled bool
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:   name,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
		ProtocolType: protocolType,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance() plugin.Plugin {
	return &Instance{}

}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
