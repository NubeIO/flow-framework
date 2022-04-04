package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

//GetDeviceByPointUUID get a device by its pointUUID
func (d *GormDatabase) GetDeviceByPointUUID(point *model.Point) (*model.Device, error) {
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

// GetPluginIDFromDevice returns the pluginUUID by using the deviceUUID to query the network.
func (d *GormDatabase) GetPluginIDFromDevice(uuid string) (*model.Network, error) {
	device, err := d.GetDevice(uuid, api.Args{})
	if err != nil {
		return nil, err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	return network, err
}

func (d *GormDatabase) deviceNameExists(dev *model.Device, body *model.Device) bool {
	var arg api.Args
	arg.WithDevices = true
	device, err := d.GetNetwork(dev.NetworkUUID, arg)
	if err != nil {
		return false
	}
	for _, p := range device.Devices {
		if p.Name == body.Name {
			if p.UUID == dev.UUID {
				return false
			} else {
				return true
			}
		}
	}
	return false
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
