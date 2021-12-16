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
func (i *Instance) wizard(body wizard) (string, error) {
	ip := "192.168.15.101"
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

	network, err := i.db.CreateNetwork(&net)
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

	device, err := i.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if device.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	var pnt model.Point
	pnt.DeviceUUID = device.UUID

	for _, e := range pointsAll() {
		pName := fmt.Sprintf("edge28_%s", e)
		pnt.Name = pName
		pnt.Description = pName
		pnt.IoID = e

		pnt.Fallback = utils.NewFloat64(0)
		if pnt.IoID == "UO1" || pnt.IoID == "UO2" {
			pnt.IoType = UOTypes.DIGITAL
		} else if pnt.IoID == "UO3" || pnt.IoID == "UO4" {
			pnt.IoType = UOTypes.VOLTSDC
		} else if pnt.IoID == "UO5" || pnt.IoID == "UO6" || pnt.IoID == "UO7" {
			pnt.IoType = UOTypes.PERCENT
		} else if pnt.IoID == "UI1" || pnt.IoID == "UI2" {
			pnt.IoType = UITypes.DIGITAL
		} else if pnt.IoID == "UI3" || pnt.IoID == "UI4" {
			pnt.IoType = UITypes.PERCENT
		} else if pnt.IoID == "UI5" {
			pnt.IoType = UITypes.VOLTSDC
		} else if pnt.IoID == "UI6" {
			pnt.IoType = UITypes.RESISTANCE
		} else if pnt.IoID == "UI7" {
			pnt.IoType = UITypes.THERMISTOR10KT2
		} else {
			pnt.IoType = UITypes.DIGITAL
		}
		pnt.IoType = string(model.IOTypeRAW)
		pnt.COV = utils.NewFloat64(0.5)
		point, err := i.db.CreatePoint(&pnt)
		if err != nil {
			return "", err
		}
		if point.UUID == "" {
			return "", errors.New("failed to create a new point")
		}

	}
	return "pass: added network and points", nil
}
