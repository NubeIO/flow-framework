package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("rubix-io: only max one network is allowed")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	body.NumberOfNetworksPermitted = nils.NewInt(1)
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	network, err := inst.db.GetNetwork(body.NetworkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, err
	}
	if len(network.Devices) >= 1 {
		errMsg := fmt.Sprintf("rubix-ior: only max one device is allowed")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	body.NumberOfDevicesPermitted = nils.NewInt(1)
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.IoNumber == "" {
		body.IoNumber = "UI1"
	}
	if body.IoType == "" {
		body.IoType = string(model.IOTypeDigital)
	}
	objectType, isOutput, isTypeBool := selectObjectType(body.IoNumber)
	body.ObjectType = objectType
	if objectType == "" {
		errMsg := fmt.Sprintf("rubix-io: point object type can not be empty")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	body.IsOutput = nils.NewBool(isOutput)
	body.IsTypeBool = nils.NewBool(isTypeBool)
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
	return ok, nil
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
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	_, err := inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("edge28-app: UpdatePoint()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateValue(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	priority := map[string]*float64{"_16": &value}
	point.InSync = boolean.NewTrue()
	_, _, _, _, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("edge28-app: pointUpdateValue()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) error {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	err = inst.db.UpdatePointErrors(uuid, &point)
	if err != nil {
		log.Error("edge28-app: pointUpdateErr()", err)
		return err
	}
	return nil
}

func selectObjectType(selectedPlugin string) (objectType string, isOutput, isTypeBool bool) {
	isOutput = false
	isOutput = false
	switch selectedPlugin {
	case PointsList.DO1.IoNumber, PointsList.DO2.IoNumber:
		objectType = PointsList.DO1.ObjectType
		isOutput = true
		isTypeBool = true
	case PointsList.UO1.IoNumber, PointsList.UO2.IoNumber, PointsList.UO3.IoNumber, PointsList.UO4.IoNumber, PointsList.UO5.IoNumber, PointsList.UO6.IoNumber:
		objectType = PointsList.UO1.ObjectType
		isOutput = true
	case PointsList.UI1.IoNumber, PointsList.UI2.IoNumber, PointsList.UI3.IoNumber, PointsList.UI4.IoNumber, PointsList.UI5.IoNumber, PointsList.UI6.IoNumber, PointsList.UI7.IoNumber, PointsList.UI8.IoNumber:
		objectType = PointsList.UI1.ObjectType
	}
	return

}

type Point struct {
	IoNumber   string // R1
	ObjectType string // binary_output
	IsOutput   *bool
	IsTypeBool *bool
}

var PointsList = struct {
	UO1 Point `json:"UO1"`
	UO2 Point `json:"UO2"`
	UO3 Point `json:"UO3"`
	UO4 Point `json:"UO4"`
	UO5 Point `json:"UO5"`
	UO6 Point `json:"UO6"`
	DO1 Point `json:"DO1"`
	DO2 Point `json:"DO2"`
	UI1 Point `json:"UI1"`
	UI2 Point `json:"UI2"`
	UI3 Point `json:"UI3"`
	UI4 Point `json:"UI4"`
	UI5 Point `json:"UI5"`
	UI6 Point `json:"UI6"`
	UI7 Point `json:"UI7"`
	UI8 Point `json:"UI8"`
}{

	UO1: Point{IoNumber: "UO1", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO2: Point{IoNumber: "UO2", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO3: Point{IoNumber: "UO3", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO4: Point{IoNumber: "UO4", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO5: Point{IoNumber: "UO5", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO6: Point{IoNumber: "UO6", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	DO1: Point{IoNumber: "DO1", ObjectType: "binary_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	DO2: Point{IoNumber: "DO2", ObjectType: "binary_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},

	UI1: Point{IoNumber: "UI1", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI2: Point{IoNumber: "UI2", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI3: Point{IoNumber: "UI3", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI4: Point{IoNumber: "UI4", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI5: Point{IoNumber: "UI5", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI6: Point{IoNumber: "UI6", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI7: Point{IoNumber: "UI7", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI8: Point{IoNumber: "UI8", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
}
