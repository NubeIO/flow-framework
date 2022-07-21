package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetDeviceByPoint get a device by point object
func (d *GormDatabase) GetDeviceByPoint(point *model.Point) (*model.Device, error) {
	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	return device, nil
}

// GetDeviceByPointUUID get a device by its pointUUID
func (d *GormDatabase) GetDeviceByPointUUID(pntUUID string) (*model.Device, error) {
	point, err := d.GetPoint(pntUUID, api.Args{})
	if err != nil || point == nil {
		return nil, err
	}

	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (d *GormDatabase) GetOneDeviceByArgs(args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) deviceNameExistsInNetwork(deviceName, networkUUID string) (device *model.Device, existing bool) {
	network, err := d.GetNetwork(networkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, false
	}
	for _, dev := range network.Devices {
		if dev.Name == deviceName {
			return dev, true
		}
	}

	return nil, false
}

// SetErrorsForAllPointsOnDevice sets the fault/error properties of all points for a specific device
// messageLevel = model.MessageLevel
// messageCode = model.CommonFaultCode
func (d *GormDatabase) SetErrorsForAllPointsOnDevice(deviceUUID string, message string, messageLevel string, messageCode string) error {
	device, err := d.GetDevice(deviceUUID, api.Args{WithPoints: true})
	if device != nil && err != nil {
		return err
	}
	for _, point := range device.Points {
		point.CommonFault.InFault = true
		point.CommonFault.MessageLevel = messageLevel
		point.CommonFault.MessageCode = messageCode
		point.CommonFault.Message = message
		point.CommonFault.LastFail = time.Now().UTC()
		err = d.UpdatePointErrors(point.UUID, point)
		if err != nil {
			log.Infof("setErrorsForAllPointsOnDevice() Error: %s\n", err.Error())
		}
	}
	return nil
}

// ClearErrorsForAllPointsOnDevice clears the fault/error properties of all points for a specific device
func (d *GormDatabase) ClearErrorsForAllPointsOnDevice(deviceUUID string) error {
	device, err := d.GetDevice(deviceUUID, api.Args{WithPoints: true})
	if device != nil && err != nil {
		return err
	}
	for _, point := range device.Points {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Normal
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = ""
		point.CommonFault.LastOk = time.Now().UTC()
		err = d.UpdatePointErrors(point.UUID, point)
		if err != nil {
			log.Infof("clearErrorsForAllPointsOnDevice() Error: %s\n", err.Error())
		}
	}
	return nil
}
