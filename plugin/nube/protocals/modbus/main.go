package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/config"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/pollqueue"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/patrickmn/go-cache"
)

const path = "modbus" //must be unique across all plugins
const name = "modbus" //must be unique across all plugins
const description = "modbus api"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "modbus"
const transportType = "ip" //serial, ip

// Instance is plugin instance
type Instance struct {
	config              *config.Config
	enabled             bool
	basePath            string
	db                  dbhandler.Handler
	store               cachestore.Handler
	bus                 eventbus.BusService
	pluginUUID          string
	networks            []*model.Network
	pollingEnabled      bool
	pollingCancel       func()
	NetworkPollManagers []*pollqueue.NetworkPollManager
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

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
