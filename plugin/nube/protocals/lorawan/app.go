package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	elsysAPB = "ELSYS-ABP"
)

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("lorawan: only max one network is allowed with lora")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err:%s", err.Error()))
	}
	return body, nil
}

// addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("issue on add lorawan-device to store")
	}
	return device, nil
}

// addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.ObjectType == "" {
		errMsg := fmt.Sprintf("lorawan: point object type can not be empty")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	point, err = inst.db.CreatePoint(body, true, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

// updateNetwork update network
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

// updateDevice update device
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return device, nil
}

// updatePoint update point
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

// deleteNetwork delete network
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, err
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		return nil, err
	}
	return point, nil
}

// deleteNetwork delete device
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// deletePoint delete point
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdate(uuid string) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.InSync = boolean.NewTrue()
	_, err := inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("lorawan: UpdatePoint()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateValue(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	priority := map[string]*float64{"_16": &value}
	point.InSync = boolean.NewTrue()
	_, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("lorawan: pointUpdateValue()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteError
	point.CommonFault.Message = fmt.Sprintf("error-time: %s msg:%s", utilstime.TimeStamp(), err.Error())
	point.CommonFault.LastFail = time.Now().UTC()
	point.InSync = boolean.NewFalse()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("lorawan: pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}
