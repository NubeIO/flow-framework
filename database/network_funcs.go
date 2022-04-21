package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

// GetNetworkByPointUUID returns a network by passing in the point UUID.
func (d *GormDatabase) GetNetworkByPointUUID(point *model.Point, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDeviceByPointUUID(point)
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
func (d *GormDatabase) GetNetworkByDeviceUUID(device *model.Device, args api.Args) (network *model.Network, err error) {
	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}
