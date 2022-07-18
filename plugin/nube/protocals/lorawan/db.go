package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) csConvertDevice(dev *model.Device, csDev *csmodel.Device) {
	dev.NetworkUUID = inst.networkUUID
	dev.CommonName.Name = csDev.Name
	dev.CommonDescription.Description = csDev.Description
	dev.CommonDevice.AddressUUID = &csDev.DevEUI
}

func (inst *Instance) getNetwork() (network *model.Network, err error) {
	net, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if len(net) == 0 {
		return nil, err
	}
	return net[0], err
}

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	if len(nets) > 0 {
		errMsg := "lorawan: only max one network is allowed with lora"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	body, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// createDeviceFromCSDevice
func (inst *Instance) createDeviceFromCSDevice(csDev *csmodel.Device) (device *model.Device, err error) {
	newDev := new(model.Device)
	inst.csConvertDevice(newDev, csDev)
	_, err = inst.db.CreateDevice(newDev)
	if err != nil {
		log.Error("lorawan: Error adding new device: ", err)
		log.Warn("lorawan: Trying with new name: ", csDev.DevEUI)
		newDev.Name = fmt.Sprintf("%s_%s", csDev.Name, csDev.DevEUI)
		_, err = inst.db.CreateDevice(newDev)
		if err != nil {
			log.Error("lorawan: Error adding new device: ", err)
			return nil, err
		}
	}
	log.Info("lorawan: Added device ", csDev.DevEUI)
	return newDev, nil
}

// createPointAddressUUID combines name and deviceEUI to form unique string to search
func createPointAddressUUID(name string, deviceEUI string) string {
	addressUUID := fmt.Sprintf("%s_%s", deviceEUI, name)
	return addressUUID
}

// getPointByAddressUUID find point in local Point slice
func (inst *Instance) getPointByAddressUUID(name string, deviceEUI string, points []*model.Point) *model.Point {
	addressUUID := createPointAddressUUID(name, deviceEUI)
	for _, point := range points {
		if *point.AddressUUID == addressUUID {
			return point
		}
	}
	return nil
}

// createNewPoint create default lorawan point
func (inst *Instance) createNewPoint(name string, deviceEUI string, deviceUUID string) (point *model.Point, err error) {
	b := false
	addressUUID := createPointAddressUUID(name, deviceEUI)
	point = &model.Point{
		CommonName:      model.CommonName{Name: name},
		AddressUUID:     &addressUUID,
		DeviceUUID:      deviceUUID,
		EnableWriteable: &b,
		ObjectType:      string(model.ObjTypeAnalogValue),
		// PointPriorityArrayMode: model.ReadOnlyNoPriorityArrayRequired,
	}
	point, err = inst.db.CreatePoint(point, true, true)
	if err != nil {
		log.Errorf("lorawan: Error creating point %s. Error: %s", addressUUID, err)
	} else {
		log.Debug("lorawan: created point ", addressUUID)
	}
	return point, err
}

// pointUpdateValue update point present value
func (inst *Instance) pointUpdateValue(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.InSync = boolean.NewTrue()
	priority := map[string]*float64{"_16": &value}
	_, _, _, _, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("lorawan: pointUpdateValue ", err)
		return nil, err
	}
	return nil, nil
}
