package main

import (
	"encoding/json"
	baseModel "github.com/NubeIO/flow-framework/model"
	bacnetServerModel "github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	model "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lorawan_model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	elsysAPB = "ELSYS-ABP"
)

//mqttUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) mqttUpdate(body mqtt.Message, devEUI, appID string) (*baseModel.Point, error) {
	//do an api call to chirpstack to get the device profile
	//decode the mqtt payload based of the device profile
	//if deviceProfileName

	payload := new(model.BasePayload)
	err := json.Unmarshal(body.Payload(), &payload)

	dev, err := i.REST.GetDevice(payload.DevEUI)
	if err != nil {
		return nil, err
	}
	//check the payload for how to decode from
	if dev.Device.DeviceProfileID == elsysAPB {
		decoded := new(model.ElsysAPB)
		err = json.Unmarshal(body.Payload(), &decoded)
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *baseModel.Point) (*baseModel.Point, error) {
	return nil, nil
}

//pointPatch from rest
func (i *Instance) pointPatch(body *baseModel.Point) (*baseModel.Point, error) {
	return nil, nil
}

//delete point make sure
func (i *Instance) deletePoint(body *baseModel.Point) (bool, error) {
	return true, nil
}

//DropDevices drop all devices
func (i *Instance) DropDevices() (bool, error) {
	cli := i.REST
	devices, err := cli.GetDevices()
	if err != nil {
		return false, err
	}
	for _, dev := range devices.Result {
		_, err := cli.DeleteDevice(dev.DevEUI)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

//delete point make sure
func (i *Instance) serverDeletePoint(body *bacnetServerModel.BacnetPoint) (bool, error) {
	return true, nil
}
