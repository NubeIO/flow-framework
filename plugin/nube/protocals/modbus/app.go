package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/utilstime"
	"go.bug.st/serial"
	"time"
)

//addDevice add network
func (i *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = i.db.CreateNetwork(body, false)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//addDevice add device
func (i *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, err = i.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//addPoint add point
func (i *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	point, err = i.db.CreatePoint(body, false, false)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//pointUpdate update point present value
func (i *Instance) pointUpdate(point *model.Point, value float64, writeSuccess, readSuccess bool) (*model.Point, error) {
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()

	if readSuccess {
		if value != utils.Float64IsNil(point.PresentValue) {
			point.ValueUpdatedFlag = utils.NewTrue() //Flag so that UpdatePointValue() will broadcast new value to producers.
		}
		modbusDebugMsg("pointUpdate() value: ", value)
		point.PresentValue = utils.NewFloat64(value)
	}
	point.InSync = utils.NewTrue()

	modbusDebugMsg("pointUpdate(): AFTER READ AND BEFORE DB UPDATE")
	point.PrintPointValues()

	//_, err = i.db.UpdatePointPresentValue(&point, true)
	_, err = i.db.UpdatePoint(point.UUID, point, true) //Changed so that Faults will update too
	if err != nil {
		modbusErrorMsg("MODBUS UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (i *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = i.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		modbusErrorMsg("MODBUS UPDATE POINT pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
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
