package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"time"

	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//addDevice add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.CreatePoint(body, true, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//updateNetwork update network
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//updateDevice update device
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//updatePoint update point
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//deleteNetwork delete network
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deleteNetwork delete device
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deletePoint delete point
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdate(point *model.Point, value float64) (*model.Point, error) {
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.OriginalValue = utils.NewFloat64(value)
	point.InSync = utils.NewTrue()
	_, err = inst.db.UpdatePointValue(point.UUID, point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT UpdatePointValue()", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

//listSerialPorts list all serial ports on host
func (inst *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
