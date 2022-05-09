package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetFlowNetworks(args api.Args) ([]*model.FlowNetwork, error) {
	var flowNetworksModel []*model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Find(&flowNetworksModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworksModel, nil
}

func (d *GormDatabase) GetFlowNetwork(uuid string, args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkByArgs(args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkModel, nil
}

/*
CreateFlowNetwork
- Create UUID
- Create Name if doesn't exist
- Create SyncUUID
- If it's pointing local device make that is_remote=false forcefully (straightforward: 0.0.0.0 can be remote)
- If local then don't apply rollback feature, it leads deadlock
- Create flow_network
- Edit body with localstorage FlowNetwork details
- Edit body with device_info
- Create FlowNetworkClone
- Update FlowNetwork with FlowNetworkClone details
- Update sync_uuid with FlowNetworkClone's sync_uuid
*/
func (d *GormDatabase) CreateFlowNetwork(body *model.FlowNetwork) (*model.FlowNetwork, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.FlowNetwork)
	isMasterSlave, cli, isRemote, tx, err := d.editFlowNetworkBody(body)
	if err != nil {
		return nil, err
	}
	// restrict to create multiple flow-networks, so it doesn't break existing local flows
	// this is needed because local flows don't have rollback
	if boolean.IsFalse(body.IsRemote) {
		fns, _ := d.GetFlowNetworks(api.Args{IsRemote: boolean.NewFalse()})
		if fns != nil && len(fns) != 0 {
			return nil, errors.New("you already have a local flow-network")
		}
	}
	if err := tx.Create(&body).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	return d.afterCreateUpdateFlowNetwork(body, isMasterSlave, cli, isRemote, tx)
}

func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var fn *model.FlowNetwork
	if err := d.DB.Where("uuid = ?", uuid).First(&fn).Error; err != nil {
		return nil, err
	}
	if len(body.Streams) > 0 {
		// we just either edit flow_network or assign stream
		return d.updateStreamsOnFlowNetwork(fn, body.Streams)
	}
	// restrict to create multiple flow-networks, so it doesn't break existing local flows
	// this is needed because local flows don't have rollback
	if boolean.IsFalse(body.IsRemote) {
		fns, _ := d.GetFlowNetworks(api.Args{IsRemote: boolean.NewFalse()})
		if fns != nil && len(fns) != 0 {
			if fns[0].UUID != uuid {
				return nil, errors.New("you already have a local flow-network")
			}
		}
	}
	isMasterSlave, cli, isRemote, tx, err := d.editFlowNetworkBody(body)
	if err != nil {
		return nil, err
	}
	if err := tx.Model(&fn).Updates(&body).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	return d.afterCreateUpdateFlowNetwork(body, isMasterSlave, cli, isRemote, tx)
}

func (d *GormDatabase) DeleteFlowNetwork(uuid string) (bool, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.DB.Where("uuid = ? ", uuid).Delete(&flowNetworkModel)
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

func (d *GormDatabase) DropFlowNetworks() (bool, error) {
	var networkModel *model.FlowNetwork
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

func (d *GormDatabase) RefreshFlowNetworksConnections() (*bool, error) {
	var flowNetworks []*model.FlowNetwork
	d.DB.Where("is_master_slave IS NOT TRUE AND is_remote IS TRUE").Find(&flowNetworks)
	for _, fn := range flowNetworks {
		accessToken, err := client.GetFlowToken(*fn.FlowIP, *fn.FlowPort, *fn.FlowUsername, *fn.FlowPassword)
		fnModel := model.FlowNetworkClone{}
		if err != nil {
			fnModel.IsError = boolean.NewTrue()
			fnModel.ErrorMsg = nstring.NewStringAddress(err.Error())
			fnModel.FlowToken = fn.FlowToken
		} else {
			fnModel.IsError = boolean.NewFalse()
			fnModel.ErrorMsg = nil
			fnModel.FlowToken = accessToken
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fn).Select("IsError", "ErrorMsg", "FlowToken").Updates(&fnModel).Error; err != nil {
			log.Error(err)
		}
	}
	return boolean.NewTrue(), nil
}

func (d *GormDatabase) editFlowNetworkBody(body *model.FlowNetwork) (bool, *client.FlowClient, bool, *gorm.DB, error) {
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = nuuid.MakeUUID()
	isMasterSlave := boolean.IsTrue(body.IsMasterSlave)
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return false, nil, false, nil, err
	}
	if isMasterSlave {
		d.resetHostAndCredential(body)
	} else if boolean.IsFalse(body.IsRemote) {
		d.resetHostAndCredential(body)
	} else {
		conf := config.Get()
		if body.FlowIP == nil {
			body.FlowIP = &conf.Server.ListenAddr
		}
		if body.FlowPort == nil {
			body.FlowPort = &conf.Server.Port
		}
		if body.FlowUsername == nil {
			return false, nil, false, nil, errors.New("FlowUsername can't be null when we it's not master/slave flow network")
		}
		if body.FlowPassword == nil {
			return false, nil, false, nil, errors.New("FlowPassword can't be null when we it's not master/slave flow network")
		}
		accessToken, err := client.GetFlowToken(*body.FlowIP, *body.FlowPort, *body.FlowUsername, *body.FlowPassword)
		if err != nil {
			return false, nil, false, nil, err
		}
		body.FlowToken = accessToken
	}

	cli := client.NewFlowClientCliFromFN(body)
	remoteDeviceInfo, err := cli.DeviceInfo()
	if err != nil {
		return false, nil, false, nil, err
	} else {
		if deviceInfo.GlobalUUID == remoteDeviceInfo.GlobalUUID {
			body.IsRemote = boolean.NewFalse()
			if !isMasterSlave {
				d.resetHostAndCredential(body)
				body.IsRemote = boolean.NewFalse()
			}
		} else {
			body.IsRemote = boolean.NewTrue()
		}
	}
	isRemote := boolean.IsTrue(body.IsRemote)
	// rollback is needed only when flow-network is remote,
	// if we make it true in local it blocks the next transaction of clone creation which leads deadlock
	var tx *gorm.DB
	if tx = d.DB; isRemote {
		tx = d.DB.Begin()
	}
	return isMasterSlave, cli, isRemote, tx, nil
}

func (d *GormDatabase) resetHostAndCredential(body *model.FlowNetwork) {
	body.FlowHTTPS = nil
	body.FlowIP = nil
	body.FlowPort = nil
	body.FlowUsername = nil
	body.FlowPassword = nil
	body.FlowToken = nil
}

func (d *GormDatabase) afterCreateUpdateFlowNetwork(body *model.FlowNetwork, isMasterSlave bool, cli *client.FlowClient, isRemote bool, tx *gorm.DB) (*model.FlowNetwork, error) {
	bodyToSync := *body
	if !isMasterSlave && isRemote {
		var localStorageFlowNetwork *model.LocalStorageFlowNetwork
		if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		bodyToSync.FlowHTTPS = localStorageFlowNetwork.FlowHTTPS
		bodyToSync.FlowIP = nstring.NewStringAddress(localStorageFlowNetwork.FlowIP)
		bodyToSync.FlowPort = integer.New(localStorageFlowNetwork.FlowPort)
		bodyToSync.FlowUsername = nstring.NewStringAddress(localStorageFlowNetwork.FlowUsername)
		bodyToSync.FlowPassword = nstring.NewStringAddress(localStorageFlowNetwork.FlowPassword)
		bodyToSync.FlowToken = nstring.NewStringAddress(localStorageFlowNetwork.FlowToken)
	}
	err := d.syncAndEditFlowNetwork(cli, body, &bodyToSync)
	if err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	if isRemote {
		tx.Commit()
	}
	d.DB.Model(&body).Updates(body)
	return body, nil
}

func (d *GormDatabase) syncAndEditFlowNetwork(cli *client.FlowClient, body *model.FlowNetwork, bodyToSync *model.FlowNetwork) error {
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return err
	}
	bodyToSync.GlobalUUID = deviceInfo.GlobalUUID
	bodyToSync.ClientId = deviceInfo.ClientId
	bodyToSync.ClientName = deviceInfo.ClientName
	bodyToSync.SiteId = deviceInfo.SiteId
	bodyToSync.SiteName = deviceInfo.SiteName
	bodyToSync.DeviceId = deviceInfo.DeviceId
	bodyToSync.DeviceName = deviceInfo.DeviceName
	res, err := cli.SyncFlowNetwork(bodyToSync)
	if err != nil {
		return err
	}
	body.SyncUUID = res.SyncUUID
	body.GlobalUUID = res.GlobalUUID
	body.ClientId = res.ClientId
	body.ClientName = res.ClientName
	body.SiteId = res.SiteId
	body.SiteName = res.SiteName
	body.DeviceId = res.DeviceId
	body.DeviceName = res.DeviceName
	return nil
}

func (d *GormDatabase) updateStreamsOnFlowNetwork(fn *model.FlowNetwork, streams []*model.Stream) (*model.FlowNetwork, error) {
	if err := d.DB.Model(&fn).Association("Streams").Replace(streams); err != nil {
		return nil, err
	}
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	for _, stream := range fn.Streams {
		_ = d.SyncStreamFunction(fn, stream, deviceInfo)
	}
	return fn, nil
}
