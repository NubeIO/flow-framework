package main

import (
	"github.com/NubeIO/flow-framework/api"
)

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if err == nil {
		inst.networkUUID = q.UUID
		inst.interruptChan = make(chan struct{}, 1)
		go inst.run()
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	select {
	case inst.interruptChan <- struct{}{}:
	default:
	}
	return nil
}
