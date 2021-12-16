package main

import (
	"errors"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

//wizard make a network/dev/pnt
func (i *Instance) wizard() (string, error) {
	var net model.Network
	net.Name = "bacnet"
	net.TransportType = model.TransType.IP
	net.PluginPath = "bacnetserver"

	network, err := i.db.CreateNetwork(&net)
	if err != nil {
		return "", err
	}
	if network.UUID == "" {
		return "", errors.New("failed to create a new network")
	}
	var dev model.Device
	dev.NetworkUUID = network.UUID
	dev.Name = "bacnet"

	device, err := i.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if device.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	var pnt model.Point
	pnt.DeviceUUID = device.UUID
	pName := utils.NameIsNil()
	pnt.Name = pName
	pnt.Description = pName
	pnt.AddressID = utils.NewInt(1)
	pnt.ObjectType = "analogValue"
	pnt.COV = utils.NewFloat64(0.5)
	pnt.Fallback = utils.NewFloat64(1)
	pnt.MessageCode = "normal"
	pnt.Unit = "noUnits"
	pnt.Priority = new(model.Priority)
	(*pnt.Priority).P16 = utils.NewFloat64(1)
	point, err := i.db.CreatePoint(&pnt)
	if err != nil {
		return "", err
	}
	if point.UUID == "" {
		return "", errors.New("failed to create a new point")
	}
	return "pass: added network and points", nil
}
