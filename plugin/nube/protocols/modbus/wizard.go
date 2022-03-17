package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

//wizard make a network/dev/pnt
func (i *Instance) wizardTCP(body wizard) (string, error) {
	ip := "192.168.15.202"
	if body.IP != "" {
		ip = body.IP
	}
	p := 502
	if body.Port != 0 {
		p = body.Port
	}
	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}
	wizVersion := 0
	if body.WizardVersion != 0 {
		wizVersion = int(body.WizardVersion)
	}

	switch wizVersion {
	case 0:
		var net model.Network
		net.Name = "modbus"
		net.TransportType = model.TransType.IP
		net.PluginPath = "modbus"

		var dev model.Device
		dev.Name = "modbus"
		dev.CommonIP.Host = ip
		dev.CommonIP.Port = p
		dev.AddressId = da
		dev.ZeroMode = utils.NewTrue()
		dev.PollDelayPointsMS = 5000

		var pnt model.Point
		pnt.Name = "modbus"
		pnt.Description = "modbus"
		pnt.AddressID = utils.NewInt(1) //TODO check conversion
		pnt.ObjectType = string(model.ObjTypeWriteFloat32)

		_, err = i.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
		if err != nil {
			return "modbus wizard 0 error: on flow-framework add modbus TCP network wizard", err
		}

		return "modbus wizard 0: added network, device, and point", err

	case 1:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.IP
		net.PluginPath = "modbus"

		net.PluginConfId = i.pluginUUID
		_, err := i.db.CreateNetwork(&net, false)
		if err != nil {
			fmt.Errorf("network creation failure: %s", err)
		}
		log.Info("Created a Network")

		for j := 1; j < 4; j++ {
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = utils.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			_, err := i.db.CreateDevice(&dev)
			if err != nil {
				fmt.Errorf("device creation failure: %s", err)
			}
			log.Info("Created a Device: ", dev)

			var pnt model.Point
			pnt.Name = "Modbus Pnt " + strconv.Itoa(j)
			pnt.Description = "modbus"
			pnt.AddressID = utils.NewInt(j) //TODO check conversion
			pnt.ObjectType = string(model.ObjTypeWriteFloat32)
			pnt.DeviceUUID = dev.UUID
			pnt.NetworkUUID = dev.NetworkUUID
			pnt.PollPriority = poller.PRIORITY_NORMAL
			pnt.PollRate = poller.RATE_NORMAL
			pnt.WriteMode = poller.ReadOnly
			_, err = i.db.CreatePoint(&pnt, false, true)
			if err != nil {
				fmt.Errorf("consumer point creation failure: %s", err)
			}
			log.Info("Created a Point for Consumer", pnt)

		}
		return "modbus wizard 1: added networks, devices, and points", err

	case 2:
		for j := 1; j < 4; j++ {
			var net model.Network
			net.Name = "Modbus Net " + strconv.Itoa(j)
			net.TransportType = model.TransType.IP
			net.PluginPath = "modbus"

			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = utils.NewTrue()
			dev.PollDelayPointsMS = 5000

			var pnt model.Point
			pnt.Name = "Modbus Pnt " + strconv.Itoa(j)
			pnt.Description = "modbus"
			pnt.AddressID = utils.NewInt(j) //TODO check conversion
			pnt.ObjectType = string(model.ObjTypeWriteFloat32)

			_, err = i.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
			if err != nil {
				return "modbus wizard 1: on flow-framework add modbus TCP network wizard", err
			}
		}
		return "modbus wizard 1: added networks, devices, and points", err

	case 3:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.Serial
		net.PluginPath = "modbus"

		net.PluginConfId = i.pluginUUID
		_, err := i.db.CreateNetwork(&net, false)
		if err != nil {
			fmt.Errorf("network creation failure: %s", err)
		}
		log.Info("Created a Network")

		for j := 1; j < 4; j++ {
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = utils.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			fastDuration, err := time.ParseDuration("5s")
			dev.FastPollRate = &fastDuration
			normalDuration, err := time.ParseDuration("30s")
			dev.NormalPollRate = &normalDuration
			slowDuration, err := time.ParseDuration("120s")
			dev.SlowPollRate = &slowDuration
			_, err = i.db.CreateDevice(&dev)
			if err != nil {
				fmt.Errorf("device creation failure: %s", err)
			}
			log.Info("Created a Device: ", dev)
			for l := 1; l < 6; l++ {
				var pnt model.Point
				pnt.Name = "Modbus Pnt " + strconv.Itoa(l)
				pnt.Description = "modbus"
				pnt.AddressID = utils.NewInt(l) //TODO check conversion
				pnt.ObjectType = string(model.ObjTypeWriteFloat32)
				pnt.DeviceUUID = dev.UUID
				pnt.NetworkUUID = dev.NetworkUUID
				pnt.PollPriority = poller.PRIORITY_NORMAL
				pnt.PollRate = poller.RATE_NORMAL
				if l == 1 {
					pnt.PollPriority = poller.PRIORITY_LOW
				} else if l == 3 {
					pnt.PollPriority = poller.PRIORITY_HIGH
				}
				pnt.WriteMode = poller.ReadOnly
				_, err = i.db.CreatePoint(&pnt, false, true)
				if err != nil {
					fmt.Errorf("consumer point creation failure: %s", err)
				}
				log.Info("Created a Point for Consumer", pnt)
			}
		}
		return "modbus wizard 3: added a network, 3 devices, and 3 points per device", err

	case 4:
		var net model.Network
		net.Name = "Modbus Net"
		net.TransportType = model.TransType.Serial
		net.PluginPath = "modbus"

		net.PluginConfId = i.pluginUUID
		_, err := i.db.CreateNetwork(&net, false)
		if err != nil {
			fmt.Errorf("network creation failure: %s", err)
		}
		log.Info("Created a Network")

		for j := 1; j < 2; j++ {
			var dev model.Device
			dev.Name = "Modbus Dev " + strconv.Itoa(j)
			dev.CommonIP.Host = "0.0.0.0"
			dev.CommonIP.Port = p
			dev.AddressId = j
			dev.ZeroMode = utils.NewTrue()
			dev.PollDelayPointsMS = 5000
			dev.NetworkUUID = net.UUID
			fastDuration, err := time.ParseDuration("5s")
			dev.FastPollRate = &fastDuration
			normalDuration, err := time.ParseDuration("30s")
			dev.NormalPollRate = &normalDuration
			slowDuration, err := time.ParseDuration("120s")
			dev.SlowPollRate = &slowDuration
			_, err = i.db.CreateDevice(&dev)
			if err != nil {
				fmt.Errorf("device creation failure: %s", err)
			}
			log.Info("Created a Device: ", dev)
			pointsArray := [4]int{401, 403, 405, 407}
			for _, l := range pointsArray {
				var pnt model.Point
				pnt.Name = "Modbus Pnt " + strconv.Itoa(l)
				pnt.Description = "modbus"
				pnt.AddressID = utils.NewInt(l) //TODO check conversion
				pnt.ObjectType = string(model.ObjTypeWriteHolding)
				pnt.DataType = string(model.TypeFloat32)
				pnt.DeviceUUID = dev.UUID
				pnt.NetworkUUID = dev.NetworkUUID
				pnt.PollPriority = poller.PRIORITY_NORMAL
				pnt.PollRate = poller.RATE_NORMAL
				pnt.WriteMode = poller.ReadOnly
				pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
				_, err = i.db.CreatePoint(&pnt, false, true)
				if err != nil {
					fmt.Errorf("consumer point creation failure: %s", err)
				}
				log.Info("Created a Point for Consumer", pnt)
			}
		}
		return "modbus wizard 4: added a network, 1 device, and 4 points", err

	}
	return "modbus wizard error: unknown wizard version", err
}

//wizard make a network/dev/pnt
func (i *Instance) wizardSerial(body wizard) (string, error) {

	sp := "/dev/ttyUSB0"
	if body.SerialPort != "" {
		sp = body.SerialPort
	}
	br := 9600
	if body.BaudRate != 0 {
		br = int(body.BaudRate)
	}

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.Serial
	net.PluginPath = "modbus"
	net.SerialPort = &sp
	net.SerialBaudRate = utils.NewUint(uint(br))
	net.SerialParity = utils.NewStr("none")
	net.SerialDataBits = utils.NewUint(1)
	net.SerialStopBits = utils.NewUint(1)

	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}

	var dev model.Device
	dev.Name = "modbus"
	dev.AddressId = da
	dev.ZeroMode = utils.NewTrue()
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressID = utils.NewInt(1) //TODO check conversion
	pnt.ObjectType = string(model.ObjTypeWriteCoil)

	pntRet, err := i.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus serial network wizard", err
	}

	log.Println(pntRet, err)
	return "pass: added network and points", err
}
