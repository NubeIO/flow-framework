package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/labstack/gommon/log"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, api.Args{})
	if q != nil {
		i.networkUUID = q.UUID
	} else {
		i.networkUUID = "NA"
	}
	if err == nil {
		i.networkUUID = q.UUID
		err = i.SerialOpen()
		if err != nil {
			log.Error("error on enable lora-plugin")
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	err := i.SerialClose()
	if err != nil {
		return err
	}
	return nil
}
