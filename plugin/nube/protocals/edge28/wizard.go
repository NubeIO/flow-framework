package main

import (
	"errors"
	"fmt"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

// wizard
type wizard struct {
	IP    string  `json:"ip"`
	Port  int     `json:"port"`
	IONum string  `json:"io_num"`
	Value float64 `json:"value"`
}

//wizard make a network/dev/pnt
func (inst *Instance) wizard(body wizard) (string, error) {
	ip := "192.168.15.63"
	if body.IP != "" {
		ip = body.IP
	}
	p := 5000
	if body.Port != 0 {
		p = body.Port
	}

	var net model.Network
	net.Name = "edge28"
	net.TransportType = model.TransType.IP
	net.PluginPath = "edge28"

	network, err := inst.db.CreateNetwork(&net, false)
	if err != nil {
		return "", err
	}
	if network.UUID == "" {
		return "", errors.New("failed to create a new network")
	}
	var dev model.Device
	dev.NetworkUUID = network.UUID
	dev.Name = "edge28"
	dev.CommonIP.Host = ip
	dev.CommonIP.Port = p
	dev.PollDelayPointsMS = 2500

	device, err := inst.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if device.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	var pnt model.Point
	pnt.DeviceUUID = device.UUID
	pnt.NetworkUUID = device.NetworkUUID

	for _, e := range pointsAll() {
		pName := fmt.Sprintf("edge28_%s", e)
		pnt.Name = pName
		pnt.Description = pName
		pnt.IoNumber = e

		pnt.Fallback = utils.NewFloat64(0)
		if pnt.IoNumber == "UO1" || pnt.IoNumber == "UO2" {
			pnt.IoType = UOTypes.DIGITAL
		} else if pnt.IoNumber == "UO3" || pnt.IoNumber == "UO4" {
			pnt.IoType = UOTypes.VOLTSDC
		} else if pnt.IoNumber == "UO5" || pnt.IoNumber == "UO6" || pnt.IoNumber == "UO7" {
			pnt.IoType = UOTypes.PERCENT
		} else if pnt.IoNumber == "UI1" || pnt.IoNumber == "UI2" {
			pnt.IoType = UITypes.DIGITAL
		} else if pnt.IoNumber == "UI3" || pnt.IoNumber == "UI4" {
			pnt.IoType = UITypes.PERCENT
		} else if pnt.IoNumber == "UI5" {
			pnt.IoType = UITypes.VOLTSDC
		} else if pnt.IoNumber == "UI6" {
			pnt.IoType = UITypes.RESISTANCE
		} else if pnt.IoNumber == "UI7" {
			pnt.IoType = UITypes.THERMISTOR10KT2
		} else {
			pnt.IoType = UITypes.DIGITAL
		}
		pnt.IoType = string(model.IOTypeRAW)
		pnt.COV = utils.NewFloat64(0.5)
		point, err := inst.db.CreatePoint(&pnt, false, false)
		if err != nil {
			return "", err
		}
		if point.UUID == "" {
			return "", errors.New("failed to create a new point")
		}

	}
	return "pass: added network and points", nil
}
