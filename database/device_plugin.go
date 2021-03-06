package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func (d *GormDatabase) CreateDevicePlugin(body *model.Device) (device *model.Device, err error) {
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

func (d *GormDatabase) UpdateDevicePlugin(uuid string, body *model.Device) (device *model.Device, err error) {
	network, err := d.GetNetwork(body.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		device, err = d.UpdateDevice(body.UUID, body, false)
		if err != nil {
			return nil, err
		}
		return
	}
	cli := client.NewLocalClient()
	device, err = cli.UpdateDevicePlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) DeleteDevicePlugin(uuid string) (ok bool, err error) {
	device, err := d.GetDevice(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	getNetwork, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return false, err
	}
	pluginName := getNetwork.PluginPath
	if pluginName == "system" {
		ok, err = d.DeleteDevice(uuid)
		if err != nil {
			return ok, err
		}
		return
	}
	cli := client.NewLocalClient()
	ok, err = cli.DeleteDevicePlugin(device, pluginName)
	if err != nil {
		return ok, err
	}
	return
}
