package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) GetNetworks(args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

func (d *GormDatabase) GetNetwork(uuid string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

// GetNetworkByField returns the network for the given field ie name or nil.
func (d *GormDatabase) GetNetworkByField(field string, value string, withDevices bool) (*model.Network, error) {
	var networkModel *model.Network
	f := fmt.Sprintf("%s = ? ", field)
	if withDevices { // drop child to reduce json size
		query := d.DB.Where(f, value).Preload("Devices").First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	} else {
		query := d.DB.Where(f, value).First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	}
}

// CreateNetwork creates a device.
func (d *GormDatabase) CreateNetwork(body *model.Network, fromPlugin bool) (*model.Network, error) {
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Network)
	body.Name = nameIsNil(body.Name)
	body.ThingClass = model.ThingClass.Network
	body.CommonEnable.Enable = utils.NewTrue()
	transport, err := checkTransport(body.TransportType) //set to ip by default
	if err != nil {
		return nil, err
	}
	body.TransportType = transport
	if body.PluginPath != "" || body.PluginConfId != "" {
		if body.PluginConfId == "" {
			plugin, err := d.GetPluginByPath(body.PluginPath)
			if err != nil {
				return nil, errors.New("failed to find a valid plugin")
			}
			if plugin.UUID == "" && body.PluginConfId != "" {
				return nil, errors.New("failed to find a valid plugin uuid")
			}
			body.PluginConfId = plugin.UUID
		}
	} else {
		return nil, errors.New("provide a plugin name ie: system, lora, modbus, lorawan, bacnet")
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	if !fromPlugin {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, body.PluginConfId, body.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, body)
	}
	if err != nil {
		return nil, errors.New("error on device eventbus")
	}
	return body, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network, fromPlugin bool) (*model.Network, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ?", uuid).First(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&networkModel, body.Tags); err != nil {
			return nil, err
		}
	}
	query = d.DB.Model(&networkModel).Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	if !fromPlugin {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, networkModel.PluginConfId, networkModel.UUID)
		d.Bus.RegisterTopic(t)
		err := d.Bus.Emit(eventbus.CTX(), t, networkModel)
		if err != nil {
			return nil, errors.New("error on network eventbus")
		}
	}
	return networkModel, nil

}

func (d *GormDatabase) DeleteNetwork(uuid string) (ok bool, err error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ? ", uuid).Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) DropNetworks() (bool, error) {
	var networkModel *model.Network
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) getPluginConf(body *model.Network) compat.Info {
	var pluginConf *model.PluginConf
	query := d.DB.Where("uuid = ?", body.PluginConfId).First(&pluginConf)
	if query.Error != nil {
		return compat.Info{}
	}
	info := d.PluginManager.PluginInfo(pluginConf.ModulePath)
	return info
}
