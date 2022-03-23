package main

import (
	"github.com/NubeIO/flow-framework/api"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/poll-queue"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	modbusDebugMsg("MODBUS Enable()")
	i.enabled = true
	i.setUUID()
	i.BusServ()
	nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, api.Args{})
	if nets != nil {
		i.networks = nets
	} else if nets == nil || err != nil {
		i.networks = nil
	}

	if i.config.EnablePolling {
		if !i.pollingEnabled {
			var arg polling
			i.pollingEnabled = true
			arg.enable = true
			i.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0) //This will delete any existing NetworkPollManagers (if enable is called multiple times, it will rebuild the queues).
			for _, net := range nets {                                       //Create a new Poll Manager for each network in the plugin.
				pollManager := pollqueue.NewPollManager(&i.db, net.UUID, i.pluginUUID)
				//modbusDebugMsg("net")
				//modbusDebugMsg("%+v\n", net)
				//modbusDebugMsg("pollManager")
				//modbusDebugMsg("%+v\n", pollManager)
				pollManager.StartPolling()
				i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
			}

			/*
				//TODO: CHECK IMPLEMENTATION OF POLLING ROUTINES
				go func() error {
					//err := i.PollingTCP(arg)
					err := i.ModbusPolling()
					if err != nil {
						log.Errorf("modbus: POLLING ERROR on routine: %v\n", err)
					}
					return nil
				}()
				if err != nil {
					log.Errorf("modbus: POLLING ERROR: %v\n", err)
				}
			*/

			//MARC TEST POLLING WITHOUT GO ROUTINE WRAPPER
			err := i.ModbusPolling()
			if err != nil {
				modbusErrorMsg("modbus: POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	modbusDebugMsg("MODBUS Disable()")
	i.enabled = false
	if i.pollingEnabled {
		var arg polling
		i.pollingEnabled = false
		arg.enable = false
		i.pollingCancel()
		i.pollingCancel = nil
		for _, pollMan := range i.NetworkPollManagers {
			pollMan.StopPolling()
		}
		i.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0)

		/*
			go func() {
				err := i.PollingTCP(arg)
				if err != nil {

				}
			}()
			if err != nil {
				return errors.New("error on starting polling")
			}
		*/
	}
	return nil
}
