package main

import (
	"encoding/json"
	model "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	bm "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// mqttUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) handleMQTT(body mqtt.Message, devEUI string) (*bm.Point, error) {

	log.Infof("lorawan: got new mqtt uplink msg from dev-id:%s", devEUI)

	payload := new(model.BasePayload)
	err := json.Unmarshal(body.Payload(), &payload)

	dev, err := inst.REST.GetDevice(payload.DevEUI)

	if err != nil {
		log.Errorf("lorawan: check device on chirpstack exists dev-id:%s  %v", devEUI, err)
		return nil, err
	}
	// check the payload for how to decode from
	if dev.Device.DeviceProfileID == elsysAPB {
		decoded := new(model.ElsysAPB)
		err = json.Unmarshal(body.Payload(), &decoded)
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}
