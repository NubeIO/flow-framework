package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

// decodeMQTT get devID from topic "application/1/device/a81758fffe05932d/event/up"
func decodeMQTT(topic string) string {
	s := strings.SplitAfter(topic, "/")
	var matchDevice bool
	var matchEvent bool
	var matchUp bool
	var deviceID string
	for i, t := range s {
		t = strings.ReplaceAll(t, "/", "")
		if t == "device" {
			matchDevice = true
		}
		if t == "event" {
			matchEvent = true
			deviceID = s[i-1]
			deviceID = strings.ReplaceAll(deviceID, "/", "")
		}
		if t == "up" {
			matchUp = true
		}
	}
	if matchDevice && matchEvent && matchUp {
		return deviceID

	}
	return ""
}

func (inst *Instance) BusServ() {
	handlerMQTT := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				p, _ := e.Data.(mqtt.Message)
				devEUI := decodeMQTT(p.Topic())
				if devEUI != "" {
					_, err := inst.handleMQTT(p, devEUI)
					if err != nil {
						return
					}
				}

			}()
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)
}
