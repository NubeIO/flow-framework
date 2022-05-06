package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) BusServ() {
	handlerCreated := bus.Handler{ //CREATED
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil { //_, err = inst.addNetwork(net)
					log.Info("BACNET-MASTER BUS PluginsCreated isNetwork", " ", net.UUID)
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
					log.Info("BACNET-MASTER BUS PluginsCreated IsDevice", " ", dev.UUID)
					//_, err = inst.addDevice(dev)
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
					log.Info("BACNET-MASTER BUS PluginsCreated IsPoint", " ", pnt.UUID)
					//_, err = inst.addPoint(pnt)
					if err != nil {
						return
					}
					if pnt != nil {
						log.Info("BACNET-MASTER BUS PluginsCreated IsPoint", " ", pnt.UUID)
						if err != nil {
							return
						}
						return
					}
				}

			}()
		},
		Matcher: eventbus.PluginsCreated,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerCreated)
	handlerUpdated := bus.Handler{ //UPDATED
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("BACNET-MASTER BUS PluginsUpdated isNetwork", " ", net.UUID)
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
					//_, err = inst.addPoints(dev)
					log.Info("BACNET-MASTER BUS PluginsUpdated IsDevice", " ", dev.UUID)
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
					//_, err = inst.pointPatch(pnt)
					log.Info("BACNET-MASTER BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsUpdated,
	}
	u, _ = nuuid.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerUpdated)
	handlerDeleted := bus.Handler{ //DELETED
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				log.Info("BACNET-MASTER BUS DELETED NEW MSG", " ", e.Topic)
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("BACNET-MASTER BUS DELETED isNetwork", " ", net.UUID)
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
					//_, err = inst.addPoints(dev)
					log.Info("BACNET-MASTER BUS DELETED IsDevice", " ", dev.UUID)
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
				log.Info("BACNET-MASTER BUS DELETED IsPoint", " ")
				if pnt != nil {
					//p, err := inst.deletePoint(pnt)
					//log.Info("BACNET-MASTER BUS DELETED IsPoint", " ", pnt.UUID, "WAS DELETED", " ", p)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsDeleted,
	}
	u, _ = nuuid.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerDeleted)

	handlerMQTT := bus.Handler{ //MQTT UPDATE (got as msg over from bacnet stack)
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				p, _ := e.Data.(mqtt.Message)
				inst.bacnetUpdate(p)

			}()
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ = nuuid.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)

}
