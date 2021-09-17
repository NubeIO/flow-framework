package main

import (
	"github.com/NubeDev/flow-framework/cachestore"
	"github.com/NubeDev/flow-framework/dbhandler"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/patrickmn/go-cache"
)

const name = "lorawan" //must be unique across all plugins
const description = "lorawan api"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "lorawan"
const transportType = "ip" //serial, ip
const ip = "0.0.0.0"
const port = "8080"

// token
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlfa2V5X2lkIjoiNzNmMDMzNjktMWE5Mi00YTBmLTkxYzItOWU1NWFjZjZlYTA4IiwiYXVkIjoiYXMiLCJpc3MiOiJhcyIsIm5iZiI6MTYzMTg0ODE3MSwic3ViIjoiYXBpX2tleSJ9.dqoUsP0tBfEX0rZQEXAWw8sSSlIqo_hfnucujy7PSJ8

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
func NewFlowPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Instance{}
}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
