package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/utils/float"
)

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.modbusDebugMsg("MODBUS Enable()")
	inst.enabled = true
	inst.setUUID()

	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.networks = nil
	}
	if inst.config.EnablePolling {
		if !inst.pollingEnabled {
			var arg polling
			inst.pollingEnabled = true
			arg.enable = true
			inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0) // This will delete any existing NetworkPollManagers (if enable is called multiple times, it will rebuild the queues).
			for _, net := range nets {                                          // Create a new Poll Manager for each network in the plugin.
				conf := inst.GetConfig().(*Config)
				pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.LogLevel}
				pollManager := NewPollManager(&pollQueueConfig, &inst.db, net.UUID, inst.pluginUUID, float.NonNil(net.MaxPollRate))
				// inst.modbusDebugMsg("net")
				// inst.modbusDebugMsg("%+v\n", net)
				// inst.modbusDebugMsg("pollManager")
				// inst.modbusDebugMsg("%+v\n", pollManager)
				pollManager.StartPolling()
				inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
			}

			// TODO: VERIFY POLLING WITHOUT GO ROUTINE WRAPPER
			err := inst.ModbusPolling()
			if err != nil {
				inst.modbusErrorMsg("POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.modbusDebugMsg("MODBUS Disable()")
	inst.enabled = false
	if inst.pollingEnabled {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
		for _, pollMan := range inst.NetworkPollManagers {
			pollMan.StopPolling()
		}
		inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0)
	}
	return nil
}
