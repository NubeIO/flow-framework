package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	elsysAPB = "ELSYS-ABP"
)

//mqttUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) mqttUpdate(body mqtt.Message, devEUI, appID string) (*model.Point, error) {
	//do an api call to chirpstack to get the device profile
	//decode the mqtt payload based of the device profile
	//if deviceProfileName

	//payload := new(lwmodel.BasePayload)
	//err := json.Unmarshal(body.Payload(), &payload)


	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {
	return nil, nil

}

//pointPatch from rest
func (i *Instance) pointPatch(body *model.Point) (*model.Point, error) {
	return nil, nil

}

//delete point make sure
func (i *Instance) deletePoint(body *model.Point) (bool, error) {
	return true, nil
}

//DropDevices drop all devices
func (i *Instance) DropDevices() (bool, error) {

	return true, nil
}

//delete point make sure
func (i *Instance) serverDeletePoint(body *pkgmodel.BacnetPoint) (bool, error) {
	return true, nil
}
