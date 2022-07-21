package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/nerrors"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	existing := d.deviceNameExists(body, body)
	if existing {
		eMsg := fmt.Sprintf("a device with existing name: %s exists", body.Name)
		return nil, nerrors.NewErrConflict(eMsg)
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	networkUUID := body.NetworkUUID
	body.Name = nameIsNil(body.Name)
	query := d.DB.Where("uuid = ? ", networkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
	}
	body.ThingClass = model.ThingClass.Device
	body.CommonEnable.Enable = boolean.NewTrue()
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
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

// UpdateDeviceErrors will only update the CommonFault properties of the device, all other properties will not be updated. Does not update `LastOk`.
func (d *GormDatabase) UpdateDeviceErrors(uuid string, body *model.Device) error {
	/* I THINK THE FIRST DB CALL HERE IS NOT REQUIRED
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&deviceModel).Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail").Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}

	*/
	query := d.DB.Model(&body).Where("uuid = ?", uuid).Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail").Updates(&body)
	if query.Error != nil {
		return query.Error
	}
	return nil
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
	if body.CommonEnable.Enable == nil {
		body.CommonEnable.Enable = boolean.NewTrue()
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&deviceModel, body.Tags); err != nil {
			return nil, err
		}
	}
	query = d.DB.Model(&deviceModel).Select("*").Updates(body)

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
	query := d.DB.Where("uuid = ? ", uuid).Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}
