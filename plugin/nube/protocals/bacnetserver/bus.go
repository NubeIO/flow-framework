package main

import (
	"context"
	"fmt"
	eventbus2 "github.com/NubeDev/flow-framework/src/eventbus"
	"github.com/NubeDev/flow-framework/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) BusServ() {
	handlerCreated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus2.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("BACNET BUS PluginsCreated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus2.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("BACNET BUS PluginsCreated IsDevice", " ", dev.UUID)
					//_, err = i.addPoints(dev)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus2.IsPoint(e.Topic, e)
				fmt.Println("BACNET ADD POINT ON BUS")
				if err != nil {
					return
				}
				_, err = i.addPoint(pnt)
				if err != nil {
					log.Error("BACNET ADD POINT BUS issue on add rest")
					return
				}
				if pnt != nil {
					log.Info("BACNET BUS PluginsCreated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus2.PluginsCreated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerCreated)
	handlerUpdated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus2.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("BACNET BUS PluginsUpdated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus2.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("BACNET BUS PluginsUpdated IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus2.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					_, err = i.pointPatch(pnt)
					log.Info("BACNET BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus2.PluginsUpdated,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerUpdated)
	handlerDeleted := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				log.Info("BACNET BUS DELETED NEW MSG", " ", e.Topic)
				//try and match is network
				net, err := eventbus2.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("BACNET BUS DELETED isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus2.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("BACNET BUS DELETED IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus2.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				log.Info("BACNET BUS DELETED IsPoint", " ")
				if pnt != nil {
					p, err := i.deletePoint(pnt)
					log.Info("BACNET BUS DELETED IsPoint", " ", pnt.UUID, "WAS DELETED", " ", p)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus2.PluginsDeleted,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerDeleted)
	handlerMQTT := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				p, _ := e.Data.(mqtt.Message)
				_, err := i.bacnetUpdate(p)
				if err != nil {
					return
				}
			}()
		},
		Matcher: eventbus2.MQTTUpdated,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerMQTT)

}
