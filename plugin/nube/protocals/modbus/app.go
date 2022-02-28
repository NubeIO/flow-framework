package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/utilstime"
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//pointUpdate update point present value
func (i *Instance) pointUpdate(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	var pri model.Priority
	pri.P16 = &value
	point.Priority = &pri
	point.InSync = utils.NewTrue()
	_, err = i.db.UpdatePointValue(uuid, &point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT UpdatePointValue()", err)
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
		log.Error("MODBUS UPDATE POINT pointUpdateErr()", err)
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
