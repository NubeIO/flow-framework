package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
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
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetwork)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = utils.MakeUUID()
	if !utils.IsTrue(body.IsRemote) || body.FlowIP == "0.0.0.0" || body.FlowIP == "127.0.0.0" || body.FlowIP == "localhost" {
		body.FlowHTTPS = utils.NewFalse()
		body.FlowIP = "0.0.0.0"
		body.FlowPort = 1660
		body.IsRemote = utils.NewFalse()
	}
	isRemote := utils.IsTrue(body.IsRemote)
	//rollback is needed only when flow-network is remote,
	//if we make it true in local it blocks the next transaction of clone creation which leads deadlock
	var tx *gorm.DB
	if tx = d.DB; isRemote {
		tx = d.DB.Begin()
	}
	if err := tx.Create(&body).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	fnb := *body
	var localStorageFlowNetwork *model.LocalStorageFlowNetwork
	if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	if utils.IsTrue(body.IsRemote) {
		fnb.FlowHTTPS = localStorageFlowNetwork.FlowHTTPS
		fnb.FlowIP = localStorageFlowNetwork.FlowIP
		fnb.FlowPort = localStorageFlowNetwork.FlowPort
	}
	fnb.FlowUsername = localStorageFlowNetwork.FlowUsername
	fnb.FlowPassword = localStorageFlowNetwork.FlowPassword
	fnb.FlowToken = localStorageFlowNetwork.FlowToken
	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	fnb.GlobalUUID = deviceInfo.GlobalUUID
	fnb.ClientId = deviceInfo.ClientId
	fnb.ClientName = deviceInfo.ClientName
	fnb.SiteId = deviceInfo.SiteId
	fnb.SiteName = deviceInfo.SiteName
	fnb.DeviceId = deviceInfo.DeviceId
	fnb.DeviceName = deviceInfo.DeviceName
	cli := client.NewSessionWithToken(body.FlowToken, body.FlowIP, body.FlowPort)
	res, err := cli.SyncFlowNetwork(&fnb)
	if err != nil {
		if isRemote {
			tx.Rollback()
		}
		return nil, err
	}
	body.SyncUUID = res.SyncUUID
	body.GlobalUUID = res.GlobalUUID
	body.ClientId = res.ClientId
	body.ClientName = res.ClientName
	body.SiteId = res.SiteId
	body.SiteName = res.SiteName
	body.DeviceId = res.DeviceId
	body.DeviceName = res.DeviceName
	d.DB.Model(&body).Updates(body)
	if isRemote {
		tx.Commit()
	}
	return body, nil
}

func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	if err := d.DB.Where("uuid = ?", uuid).Find(&flowNetworkModel).Error; err != nil {
		return nil, err
	}
	if len(body.Streams) > 0 {
		if err := d.DB.Model(&flowNetworkModel).Association("Streams").Replace(body.Streams); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&flowNetworkModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return flowNetworkModel, nil
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
