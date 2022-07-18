package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

//TODO: add function to set/clear an error on all points in a device
