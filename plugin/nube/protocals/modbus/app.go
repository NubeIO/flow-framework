package main

import (
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//pointUpdate update point present value
func (i *Instance) pointUpdate(uuid string, point *model.Point) (*model.Point, error) {
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	_, _ = i.db.UpdatePointValue(uuid, point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

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
	pnt.ObjectType = string(model.ObjTypeWriteSingleFloat32)

	_, err = i.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus TCP network wizard", err
	}

	return "pass: added network and points", err
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

//listSerialPorts list all serial ports on host
func (i *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
