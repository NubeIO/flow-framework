package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

func (d *GormDatabase) CreatePointPlugin(body *model.Point) (point *model.Point, err error) {
	network, err := d.GetNetworkByPoint(body, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		point, err = d.CreatePoint(body, false)
		if err != nil {
			return nil, err
		}
		point, err = d.UpdatePoint(point.UUID, point, false)
		return
	}
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	body.CommonFault.InFault = true
	// if plugin like bacnet then call the api direct on the plugin as the plugin knows best how to add a point to keep things in sync
	cli := client.NewLocalClient()
	point, err = cli.CreatePointPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) UpdatePointPlugin(uuid string, body *model.Point) (point *model.Point, err error) {
	network, err := d.GetNetworkByPoint(body, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		point, err = d.UpdatePoint(uuid, body, false)
		if err != nil {
			return nil, err
		}
		return
	}
	cli := client.NewLocalClient()
	point, err = cli.UpdatePointPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) WritePointPlugin(uuid string, body *model.PointWriter) (point *model.Point, err error) {
	network, err := d.GetNetworkByPointUUID(uuid, api.Args{})
	if err != nil || network == nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		point, _, _, _, err = d.PointWrite(uuid, body, false)
		if err != nil {
			return nil, err
		}
		return
	}

	cli := client.NewLocalClient()
	point, err = cli.WritePointPlugin(uuid, body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) DeletePointPlugin(uuid string) (ok bool, err error) {
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return ok, err
	}
	network, err := d.GetNetworkByPoint(point, api.Args{})
	if err != nil {
		return ok, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		ok, err = d.DeletePoint(uuid)
		if err != nil {
			return ok, err
		}
		return
	}
	cli := client.NewLocalClient()
	ok, err = cli.DeletePointPlugin(point, pluginName)
	if err != nil {
		return ok, err
	}
	return
}
