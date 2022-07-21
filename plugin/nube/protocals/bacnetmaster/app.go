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

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body.NetworkInterface == "" {
		return nil, errors.New("network interface can not be empty try, eth0")
	}
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetStoreNetwork(network)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("issue on add bacnet-device to store err:%s", err.Error()))
	}
	return body, nil
}

func (inst *Instance) getNetworks() ([]*model.Network, error) {
	return inst.db.GetNetworks(api.Args{})
}

// addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetStoreDevice(device)
	if err != nil {
		return nil, errors.New("issue on add bacnet-device to store")
	}
	return device, nil
}

// addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.ObjectType == "" {
		errMsg := fmt.Sprintf("bacnet-master: point object type can not be empty")
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
	err = inst.bacnetStoreDevice(device)
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
	ok, err = inst.closeBacnetStoreNetwork(body.UUID)
	return ok, err
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, _, _, _, err = inst.db.WritePoint(pntUUID, body, true)
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
		log.Error("bacnet-master: UpdatePoint()", err)
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
	_, _, _, _, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("bacnet-master: pointUpdateValue()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) error {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteError
	point.CommonFault.Message = fmt.Sprintf("error-time: %s msg:%s", utilstime.TimeStamp(), err.Error())
	point.CommonFault.LastFail = time.Now().UTC()
	point.InSync = boolean.NewFalse()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("bacnet-master: pointUpdateErr()", err)
		return err
	}
	return nil
}
