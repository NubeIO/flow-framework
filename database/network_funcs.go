package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetNetworkByPluginName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByPluginName(name string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_path = ? ", name).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

// GetNetworksByPluginName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByPluginName(name string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_path = ? ", name).Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworkByPlugin returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByPlugin(pluginUUID string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_conf_id = ? ", pluginUUID).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

// GetNetworksByPlugin returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByPlugin(pluginUUID string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_conf_id = ? ", pluginUUID).Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworksByName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByName(name string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Find(&networksModel).Where("name = ? ", name).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworkByName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByName(name string, args api.Args) (*model.Network, error) {
	var networksModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("name = ? ", name).First(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworkByPoint returns a network by passing in the pointUUID.
func (d *GormDatabase) GetNetworkByPoint(point *model.Point, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDeviceByPoint(point)
	if err != nil {
		return nil, err
	}
	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}

// GetNetworkByPointUUID returns a network by passing in the pointUUID.
func (d *GormDatabase) GetNetworkByPointUUID(pntUUID string, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDeviceByPointUUID(pntUUID)
	if err != nil {
		return nil, err
	}
	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}

// GetNetworkByDeviceUUID returns a network by passing in the device UUID.
func (d *GormDatabase) GetNetworkByDeviceUUID(devUUID string, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDevice(devUUID, args)
	if err != nil && device == nil {
		return nil, err
	}

	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}

// SetErrorsForAllDevicesOnNetwork sets the fault/error properties of all devices for a specific network. Optional to set the points of each device also.
// messageLevel = model.MessageLevel
// messageCode = model.CommonFaultCode
func (d *GormDatabase) SetErrorsForAllDevicesOnNetwork(networkUUID string, message string, messageLevel string, messageCode string, doPoints bool) error {
	network, err := d.GetNetwork(networkUUID, api.Args{WithDevices: true, WithPoints: doPoints})
	if network != nil && err != nil {
		return err
	}
	for _, device := range network.Devices {
		device.CommonFault.InFault = true
		device.CommonFault.MessageLevel = messageLevel
		device.CommonFault.MessageCode = messageCode
		device.CommonFault.Message = message
		device.CommonFault.LastFail = time.Now().UTC()
		err = d.UpdateDeviceErrors(device.UUID, device)
		if err != nil {
			log.Infof("setErrorsForAllDevicesOnNetwork() Error: %s\n", err.Error())
		}
		if doPoints {
			err = d.SetErrorsForAllPointsOnDevice(device.UUID, message, messageLevel, messageCode)
		}
	}
	return nil
}

// ClearErrorsForAllDevicesOnNetwork clears the fault/error properties of all devices for a specific network. Optional to clear the points of each device also.
func (d *GormDatabase) ClearErrorsForAllDevicesOnNetwork(networkUUID string, doPoints bool) error {
	network, err := d.GetNetwork(networkUUID, api.Args{WithDevices: true, WithPoints: doPoints})
	if network != nil && err != nil {
		return err
	}
	for _, device := range network.Devices {
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Normal
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
		err = d.UpdateDeviceErrors(device.UUID, device)
		if err != nil {
			log.Infof("clearErrorsForAllDevicesOnNetwork() Error: %s\n", err.Error())
		}
		if doPoints {
			err = d.ClearErrorsForAllPointsOnDevice(device.UUID)
		}
	}
	return nil
}
