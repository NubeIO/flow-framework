package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) BusServ() {
	handlerCreated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					i.OnNetworkCreated(net)
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
					i.OnDeviceCreated(dev)
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
					i.OnPointCreated(pnt)
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
					i.OnNetworkUpdated(net)
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
					i.OnDeviceUpdated(dev)
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
					i.OnPointUpdated(pnt)
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
					netUUID := eventbus.GetTopicPart(e.Topic, 3, "net")
					i.OnNetworkDeleted(netUUID)
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
					if e.Topic == "ALL" {
						log.Error("MODBUS: DropDevices() has been called, check operation of polling.")
						for index, netPollMan := range i.NetworkPollManagers {
							netPollMan.StopPolling()
							//Next remove the NetworkPollManager from the slice in polling instance
							i.NetworkPollManagers[index] = i.NetworkPollManagers[len(i.NetworkPollManagers)-1]
							i.NetworkPollManagers = i.NetworkPollManagers[:len(i.NetworkPollManagers)-1]
						}
					}
					fmt.Println("MODBUS DEVICE DELETED:")
					fmt.Printf("%+v\n", dev)
					fmt.Printf("%+v\n", e)
					netUUID := eventbus.GetTopicPart(e.Topic, 2, "net")
					devUUID := eventbus.GetTopicPart(e.Topic, 3, "dev")
					i.OnDeviceDeleted(devUUID, netUUID)
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
					fmt.Println("MODBUS POINT DELETED:")
					fmt.Printf("%+v\n", pnt)
					fmt.Printf("%+v\n", e)
					netUUID := eventbus.GetTopicPart(e.Topic, 2, "net")
					pntUUID := eventbus.GetTopicPart(e.Topic, 3, "pnt")
					devUUID := eventbus.GetTopicPart(e.Topic, 4, "dev")
					i.OnPointDeleted(pntUUID, devUUID, netUUID)
					//p, err := i.deletePoint(pnt)
					log.Info("MODBUS BUS DELETED IsPoint", " ", pntUUID, "WAS DELETED", " ", "p")
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

func (i *Instance) OnNetworkCreated(net *model.Network) {
	fmt.Println("MODBUS OnNetworkCreated(): ", net.UUID)
	if net == nil {
		log.Error("Modbus: OnNetworkCreated(): cannot find network ", net.UUID)
		return
	}
	if utils.BoolIsNil(net.Enable) == false {
		return
	}

	pollManager := pollqueue.NewPollManager(&i.db, net.UUID, i.pluginUUID)
	pollManager.StartPolling()
	i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
}

func (i *Instance) OnDeviceCreated(dev *model.Device) {
	fmt.Println("MODBUS OnDeviceCreated(): ", dev.UUID)
	if dev == nil {
		log.Error("Modbus: OnDeviceCreated(): cannot find device ", dev.UUID)
		return
	}
	//NOTHING TO DO ON DEVICE CREATED
}

func (i *Instance) OnPointCreated(pnt *model.Point) {
	fmt.Println("MODBUS OnPointCreated(): ", pnt.UUID)
	fmt.Println("MODBUS OnPointCreated(): point")
	fmt.Printf("%+v\n", pnt)
	//NOTHING REQUIRED AS OnPointUpdated() is called whenever a point is created.
	/*
		if pnt == nil {
			log.Error("Modbus: OnPointCreated(): cannot find point ", pnt.UUID)
			return
		}
		if utils.BoolIsNil(pnt.Enable) {
			//DO ADD POINT POLL QUEUE OPERATIONS
			netPollMan, err := i.getNetworkPollManagerByUUID(pnt.NetworkUUID)
			if netPollMan == nil || err != nil {
				log.Error("Modbus: OnPointCreated(): cannot find NetworkPollManager for network: ", pnt.NetworkUUID)
				return
			}
			pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, pnt.NetworkUUID, netPollMan.FFPluginUUID)
			pp.PollPriority = pnt.PollPriority
			netPollMan.PollQueue.AddPollingPoint(pp)
		}

	*/
}

func (i *Instance) OnNetworkUpdated(net *model.Network) {
	fmt.Println("MODBUS OnNetworkUpdated(): ", net.UUID)
	if net == nil {
		log.Error("Modbus: OnNetworkUpdated(): cannot find network ", net.UUID)
		return
	}
	netPollMan, err := i.getNetworkPollManagerByUUID(net.UUID)
	if netPollMan == nil || err != nil {
		log.Error("Modbus: OnNetworkUpdated(): cannot find NetworkPollManager for network: ", net.UUID)
		return
	}

	if utils.BoolIsNil(net.Enable) == false && netPollMan.Enable == true {
		//DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()

	} else if utils.BoolIsNil(net.Enable) == true && netPollMan.Enable == false {
		//DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
	}
}

func (i *Instance) OnDeviceUpdated(dev *model.Device) {
	fmt.Println("MODBUS OnDeviceUpdated(): ", dev.UUID)
	if dev == nil {
		log.Error("Modbus: OnDeviceUpdated(): cannot find device ", dev.UUID)
		return
	}
	netPollMan, err := i.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		log.Error("Modbus: OnDeviceUpdated(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}

	if utils.BoolIsNil(dev.Enable) == false && netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		//DO POLLING DISABLE ACTIONS FOR DEVICE
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)

	} else if utils.BoolIsNil(dev.Enable) == true && !netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(dev.UUID) {
		//DO POLLING ENABLE ACTIONS FOR DEVICE
		for _, pnt := range dev.Points {
			if utils.BoolIsNil(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, pnt.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollQueue.AddPollingPoint(pp)
			}
		}

	} else if utils.BoolIsNil(dev.Enable) == true {
		//TODO: Currently on every device update, all device points are removed, and re-added.
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(dev.UUID)
		for _, pnt := range dev.Points {
			if utils.BoolIsNil(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, pnt.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollQueue.AddPollingPoint(pp)
			}
		}
	}
	// TODO: NEED TO ACCOUNT FOR OTHER CHANGES ON DEVICE.  It would be useful to have a way to know if the device polling rates were changed.
}

func (i *Instance) OnPointUpdated(pnt *model.Point) {
	fmt.Println("MODBUS OnPointUpdated(): ", pnt.UUID)
	fmt.Printf("%+v\n", pnt)
	fmt.Println("MODBUS OnPointUpdated(): PRIORITY")
	fmt.Printf("%+v\n", pnt.Priority)
	if pnt == nil {
		log.Error("Modbus: OnPointUpdated(): cannot find point ", pnt.UUID)
		return
	}
	netPollMan, err := i.getNetworkPollManagerByUUID(pnt.NetworkUUID)
	if netPollMan == nil || err != nil {
		log.Error("Modbus: OnPointUpdated(): cannot find NetworkPollManager for network: ", pnt.NetworkUUID)
		return
	}

	if utils.BoolIsNil(pnt.Enable) == false {
		//DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(pnt.UUID)

	} else if utils.BoolIsNil(pnt.Enable) == true {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(pnt.UUID)
		//DO POLLING ENABLE ACTIONS FOR POINT
		//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, pnt.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = pnt.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode.
		//netPollMan.PollQueue.AddPollingPoint(pp)
		//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	}
}

func (i *Instance) OnNetworkDeleted(netUUID string) {
	fmt.Println("MODBUS OnNetworkDeleted(): ", netUUID)
	found := false
	for index, netPollMan := range i.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == netUUID {
			netPollMan.StopPolling()
			//Next remove the NetworkPollManager from the slice in polling instance
			i.NetworkPollManagers[index] = i.NetworkPollManagers[len(i.NetworkPollManagers)-1]
			i.NetworkPollManagers = i.NetworkPollManagers[:len(i.NetworkPollManagers)-1]
			found = true
		}
	}
	if !found {
		log.Error("Modbus: OnNetworkDeleted(): cannot find NetworkPollManager for network: ", netUUID)
	}
}

func (i *Instance) OnDeviceDeleted(devUUID, netUUID string) {
	fmt.Println("MODBUS OnDeviceDeleted(): devUUID: ", devUUID, "  netUUID: ", netUUID)
	netPollMan, err := i.getNetworkPollManagerByUUID(netUUID)
	if netPollMan == nil || err != nil {
		log.Error("Modbus: OnDeviceDeleted(): cannot find NetworkPollManager for network: ", netUUID)
		return
	}
	netPollMan.PollQueue.RemovePollingPointByDeviceUUID(devUUID)

}

func (i *Instance) OnPointDeleted(pntUUID, devUUID, netUUID string) {
	fmt.Println("MODBUS OnPointDeleted(): ", pntUUID)
	netPollMan, err := i.getNetworkPollManagerByUUID(netUUID)
	if netPollMan == nil || err != nil {
		log.Error("Modbus: OnPointDeleted(): cannot find NetworkPollManager for network: ", netUUID)
		return
	}
	netPollMan.PollQueue.RemovePollingPointByPointUUID(pntUUID)
	otherPointsOnSameDeviceExist := netPollMan.PollQueue.CheckPollingQueueForDevUUID(devUUID)
	if !otherPointsOnSameDeviceExist {
		netPollMan.PollQueue.RemoveDeviceFromActiveDevicesList(devUUID)
	}
}
