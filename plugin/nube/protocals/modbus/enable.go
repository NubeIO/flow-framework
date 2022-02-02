package main

import (
	"github.com/NubeIO/flow-framework/api"
	log "github.com/sirupsen/logrus"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, api.Args{})
	if nets != nil {
		i.networks = nets
	} else {
		i.networks = nil
	}
	if i.config.EnablePolling {
		if !i.pollingEnabled {
			i.pollingEnabled = true
			for _, net := range nets {
				net.PollManager.StartPolling()
			}
			func() error {
				err := i.ModbusPolling()
				if err != nil {
					log.Errorf("modbus: POLLING ERROR on routine: %v\n", err)
				}
				return nil
			}()
			if err != nil {
				log.Errorf("modbus: POLLING ERROR: %v\n", err)
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	if i.pollingEnabled {
		i.pollingEnabled = false
	}
	if i.pollingCancel != nil {
		i.pollingCancel()
	}
	return nil
}

/*
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
	if i.config.EnablePolling {
		if !i.pollingEnabled {
			var arg polling
			i.pollingEnabled = true
			arg.enable = true
			go func() error {
				err := i.PollingTCP(arg)
				if err != nil {
					log.Errorf("modbus: POLLING ERROR on routine: %v\n", err)
				}
				return nil
			}()
			if err != nil {
				log.Errorf("modbus: POLLING ERROR: %v\n", err)
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	if i.pollingEnabled {
		var arg polling
		i.pollingEnabled = false
		arg.enable = false
		go func() {
			err := i.PollingTCP(arg)
			if err != nil {

			}
		}()
		if err != nil {
			return errors.New("error on starting polling")
		}
	}
	return nil
}
*/
