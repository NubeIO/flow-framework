package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/labstack/gommon/log"
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
	if err != nil {
		log.Error("error on enable lorawan-plugin", err)
	}

	// TODO: temporary call due to config being broken
	inst.ValidateAndSetConfig(&Config{})

	inst.REST = csrest.CSLogin(inst.config.CSAddress, inst.config.CSPort,
		inst.config.CSUsername, inst.config.CSPassword)
	inst.REST.SetDeviceLimit(inst.config.DeviceLimit)
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
