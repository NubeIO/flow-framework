package main

import (
	"errors"
	"math/rand"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

//wizard make a network/dev/pnt
func (inst *Instance) wizard() (string, error) {
	var net model.Network
	net.Name = "bacnet"
	net.TransportType = model.TransType.IP
	net.PluginPath = "bacnetserver"

	network, err := inst.db.CreateNetwork(&net, false)
	if err != nil {
		return "", err
	}
	if network.UUID == "" {
		return "", errors.New("failed to create a new network")
	}
	var dev model.Device
	dev.NetworkUUID = network.UUID
	dev.Name = "bacnet"

	device, err := inst.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if device.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	min := 1
	max := 1000
	a := rand.Intn(max-min) + min

	var pnt model.Point
	pnt.DeviceUUID = device.UUID
	pnt.NetworkUUID = device.NetworkUUID
	pName := utils.NameIsNil()
	pnt.Name = pName
	pnt.Description = pName
	pnt.AddressID = utils.NewInt(a)
	pnt.ObjectType = "analogValue"
	pnt.COV = utils.NewFloat64(0.5)
	pnt.Fallback = utils.NewFloat64(1)
	pnt.MessageCode = "normal"
	pnt.Unit = utils.NewStringAddress("noUnits")
	pnt.Priority = new(model.Priority)
	(*pnt.Priority).P16 = utils.NewFloat64(1)
	point, err := inst.db.CreatePoint(&pnt, false, false)
	if err != nil {
		return "", err
	}
	if point.UUID == "" {
		return "", errors.New("failed to create a new point")
	}
	return "pass: added network and points", nil
}
