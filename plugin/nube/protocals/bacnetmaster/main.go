package main

import (
	"container/heap"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

const path = "bacnetmaster" // must be unique across all plugins
const name = "bacnetmaster" // must be unique across all plugins
const description = "bacnet bserver api to nube bacnet stack"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "bacnet"
const transportType = "ip" // serial, ip

// Instance is plugin instance
type Instance struct {
	config              *Config
	enabled             bool
	basePath            string
	db                  dbhandler.Handler
	store               cachestore.Handler
	bus                 eventbus.BusService
	pluginUUID          string
	networks            []*model.Network
	pollingEnabled      bool
	BacStore            *network.Store
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

func NewPollManager(conf *pollqueue.Config, dbHandler *dbhandler.Handler, ffNetworkUUID, ffPluginUUID string, maxPollRate float64) *pollqueue.NetworkPollManager {
	// Make the main priority polling queue
	queue := make([]*pollqueue.PollingPoint, 0)
	pq := &pollqueue.PriorityPollQueue{PriorityQueue: queue}
	heap.Init(pq)                                  // Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	refQueue := make([]*pollqueue.PollingPoint, 0) // Make the reference slice that contains points that are not in the current polling queue.
	rq := &pollqueue.PriorityPollQueue{PriorityQueue: refQueue}
	heap.Init(rq)                                          // Init needs to be called on the main PriorityQueue so that it is maintained by PollingPriority.
	outstandingQueue := make([]*pollqueue.PollingPoint, 0) // Make the reference slice that contains points that are not in the current polling queue.
	opq := &pollqueue.PriorityPollQueue{PriorityQueue: outstandingQueue}
	heap.Init(opq)
	adl := make([]string, 0)
	pqu := &pollqueue.QueueUnloader{NextPollPoint: nil, NextUnloadTimer: nil, CancelChannel: nil}
	puwp := make(map[string]bool)
	npq := &pollqueue.NetworkPriorityPollQueue{Config: conf, PriorityQueue: pq, StandbyPollingPoints: rq, OutstandingPollingPoints: opq, PointsUpdatedWhilePolling: puwp, QueueUnloader: pqu, FFPluginUUID: ffPluginUUID, FFNetworkUUID: ffNetworkUUID, ActiveDevicesList: adl}
	pm := new(pollqueue.NetworkPollManager)
	pm.Config = conf
	pm.PollQueue = npq
	pm.PluginQueueUnloader = pqu
	pm.DBHandlerRef = dbHandler
	pm.MaxPollRate = time.Duration(maxPollRate) * time.Second
	pm.FFNetworkUUID = ffNetworkUUID
	pm.FFPluginUUID = ffPluginUUID
	pm.ASAPPriorityMaxCycleTime, _ = time.ParseDuration("2m")
	pm.HighPriorityMaxCycleTime, _ = time.ParseDuration("5m")
	pm.NormalPriorityMaxCycleTime, _ = time.ParseDuration("15m")
	pm.LowPriorityMaxCycleTime, _ = time.ParseDuration("60m")
	return pm
}

// main will not let main run
func main() {
	panic("this should be built as plugin")
}
