package main

import (
	"errors"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nums"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/system/networking"
	"math/rand"
)

//wizard make a network/dev/pnt
func (i *Instance) wizard(network *Network) (string, error) {
	var net model.Network
	net.Name = network.NetworkName
	net.TransportType = model.TransType.IP
	net.PluginPath = "bacnetmaster"

	if network.InterfaceName != "" {
		_net, _ := networking.GetInterfaceByName(network.InterfaceName)
		if _net == nil {
			return "", errors.New("failed to find a valid network interface")
		}
		network.NetworkIp = _net.IP
		network.NetworkMask = _net.NetMaskLength
	} else {
		net.IP = network.NetworkIp
		net.Port = nums.NewInt(network.NetworkPort)
		net.NetworkMask = nums.NewInt(network.NetworkMask)
	}
	_network, err := i.db.CreateNetwork(&net, false)
	if err != nil {
		return "", err
	}
	if _network.UUID == "" {
		return "", errors.New("failed to create a new network")
	}
	var dev model.Device
	dev.NetworkUUID = _network.UUID
	dev.Name = "bacnet"

	d, err := i.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if d.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	min := 1
	max := 1000
	a := rand.Intn(max-min) + min

	var pnt model.Point
	pnt.DeviceUUID = d.UUID
	pnt.NetworkUUID = d.NetworkUUID
	pName := utils.NameIsNil()
	pnt.Name = pName
	pnt.Description = pName
	pnt.AddressID = nums.NewInt(a)
	pnt.ObjectType = "analogValue"
	pnt.COV = nums.NewFloat64(0.5)
	pnt.Fallback = nums.NewFloat64(1)
	pnt.MessageCode = "normal"
	pnt.Unit = utils.NewStringAddress("noUnits")
	pnt.Priority = new(model.Priority)
	(*pnt.Priority).P16 = nums.NewFloat64(1)
	point, err := i.db.CreatePoint(&pnt, false, false)
	if err != nil {
		return "", err
	}
	if point.UUID == "" {
		return "", errors.New("failed to create a new point")
	}
	return "pass: added network and points", nil
}
