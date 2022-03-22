package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func (d *GormDatabase) GetDevices(args api.Args) ([]*model.Device, error) {
	var devicesModel []*model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Find(&devicesModel).Error; err != nil {
		return nil, err
	}
	return devicesModel, nil
}

func (d *GormDatabase) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) GetOneDeviceByArgs(args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

//TODO: this function doesn't do what it says it does
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

func (d *GormDatabase) CreateDevicePlugin(body *model.Device) (device *model.Device, err error) {
	fmt.Println(body, body.NetworkUUID)
	network, err := d.GetNetwork(body.NetworkUUID, api.Args{})
	if network == nil {
		errMsg := fmt.Sprintf("model.device failed to find a network with uuid:%s", body.NetworkUUID)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	pluginName := network.PluginPath

	if pluginName == "system" {
		device, err = d.CreateDevice(body)
		if err != nil {
			return nil, err
		}
		return
	}
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	cli := client.NewLocalClient()
	device, err = cli.CreateDevicePlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	existing := d.deviceNameExists(body, body)
	if existing {
		eMsg := fmt.Sprintf("a device with existing name: %s exists", body.Name)
		return nil, errors.New(eMsg)
	}
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Device)
	networkUUID := body.NetworkUUID
	body.Name = nameIsNil(body.Name)
	query := d.DB.Where("uuid = ? ", networkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
	}
	body.ThingClass = model.ThingClass.Device
	body.CommonEnable.Enable = utils.NewTrue()
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, query.Error
	}
	var nModel *model.Network
	query = d.DB.Where("uuid = ?", body.NetworkUUID).First(&nModel)
	if query.Error != nil {
		return nil, query.Error
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, nModel.PluginConfId, body.UUID)
	d.Bus.RegisterTopic(t)
	err := d.Bus.Emit(eventbus.CTX(), t, body)
	if err != nil {
		return nil, errors.New("error on device eventbus")
	}
	return body, query.Error
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device, fromPlugin bool) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Name != "" {
		existing := d.deviceNameExists(deviceModel, body)
		if existing {
			eMsg := fmt.Sprintf("a device with existing name: %s exists", body.Name)
			return nil, errors.New(eMsg)
		}
	}
	if query.Error != nil {
		return nil, query.Error
	}
	if body.CommonEnable.Enable != nil {
		body.CommonEnable.Enable = utils.NewTrue()
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&deviceModel, body.Tags); err != nil {
			return nil, err
		}
	}
	query = d.DB.Model(&deviceModel).Updates(body)

	var nModel *model.Network
	query = d.DB.Where("uuid = ?", deviceModel.NetworkUUID).First(&nModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if !fromPlugin {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, nModel.PluginConfId, uuid)
		d.Bus.RegisterTopic(t)
		err := d.Bus.Emit(eventbus.CTX(), t, deviceModel)
		if err != nil {
			return nil, errors.New("error on device eventbus")
		}
	}

	return deviceModel, nil
}

// DeleteDevice delete a Device.
func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	var deviceModel *model.Device
	network, err := d.GetPluginIDFromDevice(uuid)
	query := d.DB.Where("uuid = ? ", uuid).Delete(&deviceModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsDeleted, network.UUID, uuid) //TODO: second argument in other instances is pluginUUID, but i didn't have access in this case
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, deviceModel)
		if err != nil {
			return false, errors.New("error on device delete eventbus")
		}
		return true, nil
	}
}

// DropDevices delete all devices.
func (d *GormDatabase) DropDevices() (bool, error) {
	var deviceModel *model.Device
	query := d.DB.Where("1 = 1").Delete(&deviceModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsDeleted, "ALL", "ALL") //TODO: second argument in other instances is pluginUUID, but i didn't have access in this case
	d.Bus.RegisterTopic(t)
	err := d.Bus.Emit(eventbus.CTX(), t, deviceModel)
	if err != nil {
		return false, errors.New("error on device eventbus")
	}
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
