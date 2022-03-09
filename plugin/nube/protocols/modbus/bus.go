package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocols/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) BusServ() {
	handlerCreated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("MODBUS BUS PluginsCreated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("MODBUS BUS PluginsCreated IsDevice", " ", dev.UUID)
					//_, err = i.addPoints(dev)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				fmt.Println("ADD POINT ON BUS")
				if err != nil {
					return
				}
				//_, err = i.addPoint(pnt)
				if err != nil {
					return
				}
				if pnt != nil {
					log.Info("MODBUS BUS PluginsCreated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsCreated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerCreated)
	handlerUpdated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("MODBUS BUS PluginsUpdated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("MODBUS BUS PluginsUpdated IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					//_, err = i.pointPatch(pnt)
					log.Info("MODBUS BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsUpdated,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerUpdated)
	handlerDeleted := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				log.Info("MODBUS BUS DELETED NEW MSG", " ", e.Topic)
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("MODBUS BUS DELETED isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("MODBUS BUS DELETED IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				log.Info("MODBUS BUS DELETED IsPoint", " ")
				if pnt != nil {
					//p, err := i.deletePoint(pnt)
					log.Info("MODBUS BUS DELETED IsPoint", " ", pnt.UUID, "WAS DELETED", " ", "p")
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsDeleted,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerDeleted)

}

func (i *Instance) OnNetworkAdded(networkUUID *string) {
	net, err := i.db.GetNetwork(*networkUUID, api.Args{})
	if net == nil || err != nil {
		log.Error("Modbus: OnNetworkAdded(): cannot find network ", *networkUUID)
		return
	}
	if *net.Enable == false {
		return
	}

	pollManager := pollqueue.NewPollManager(&i.db, *networkUUID, i.pluginUUID)
	pollManager.StartPolling()
	i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
}

func (i *Instance) OnDeviceAdded(deviceUUID *string) {
	dev, err := i.db.GetDevice(*deviceUUID, api.Args{})
	if dev == nil || err != nil {
		log.Error("Modbus: OnDeviceAdded(): cannot find network ", *deviceUUID)
		return
	}
	if *dev.Enable == false {
		return
	}

	//DO ADD DEVICE POLL QUEUE OPERATIONS

}

func (i *Instance) OnPointAdded(pointUUID *string) {
	pnt, err := i.db.GetPoint(*pointUUID)
	if pnt == nil || err != nil {
		log.Error("Modbus: OnPointAdded(): cannot find network ", *pointUUID)
		return
	}
	if *pnt.Enable == false {
		return
	}

	//DO ADD POINT POLL QUEUE OPERATIONS

}

func (i *Instance) OnNetworkUpdated(networkUUID *string) {
	net, err := i.db.GetNetwork(*networkUUID, api.Args{})
	if net == nil || err != nil {
		log.Error("Modbus: OnNetworkUpdated(): cannot find network ", *networkUUID)
		return
	}
	if *net.Enable == false {
		return
	}

	pollManager := pollqueue.NewPollManager(&i.db, *networkUUID, i.pluginUUID)
	pollManager.StartPolling()
	i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
}

func (i *Instance) OnDeviceUpdated(deviceUUID *string) {
	dev, err := i.db.GetDevice(*deviceUUID, api.Args{})
	if dev == nil || err != nil {
		log.Error("Modbus: OnDeviceUpdated(): cannot find network ", *deviceUUID)
		return
	}
	if *dev.Enable == false {
		return
	}

	//DO ADD DEVICE POLL QUEUE OPERATIONS

}

func (i *Instance) OnPointUpdated(pointUUID *string) {
	pnt, err := i.db.GetDevice(*pointUUID, api.Args{})
	if pnt == nil || err != nil {
		log.Error("Modbus: OnPointUpdated(): cannot find network ", *pointUUID)
		return
	}
	if *pnt.Enable == false {
		return
	}

	//DO ADD POINT POLL QUEUE OPERATIONS

}

func (i *Instance) OnNetworkDeleted(networkUUID *string) {
	net, err := i.db.GetNetwork(*networkUUID, api.Args{})
	if net == nil || err != nil {
		log.Error("Modbus: OnNetworkDeleted(): cannot find network ", *networkUUID)
		return
	}
	if *net.Enable == false {
		return
	}

	pollManager := pollqueue.NewPollManager(&i.db, *networkUUID, i.pluginUUID)
	pollManager.StartPolling()
	i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
}

func (i *Instance) OnDeviceDeleted(deviceUUID *string) {
	dev, err := i.db.GetDevice(*deviceUUID, api.Args{})
	if dev == nil || err != nil {
		log.Error("Modbus: OnDeviceDeleted(): cannot find network ", *deviceUUID)
		return
	}
	if *dev.Enable == false {
		return
	}

	//DO ADD DEVICE POLL QUEUE OPERATIONS

}

func (i *Instance) OnPointDeleted(pointUUID *string) {
	pnt, err := i.db.GetPoint(*pointUUID)
	if pnt == nil || err != nil {
		log.Error("Modbus: OnPointDeleted(): cannot find network ", *pointUUID)
		return
	}
	if *pnt.Enable == false {
		return
	}
	/*
		RemovePollingPointByPointUUID
		pnt, err := i.db.GetDevice(*pointUUID, api.Args{})
		if pnt == nil || err != nil {
			log.Error("Modbus: OnPointDeleted(): cannot find network ", *pointUUID)
			return
		}
		if *pnt.Enable == false {
			return
		}

		//DO ADD POINT POLL QUEUE OPERATIONS
	*/
}
