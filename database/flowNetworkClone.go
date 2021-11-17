package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	var flowNetworkClonesModel []*model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkClonesModel, nil
}

func (d *GormDatabase) GetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkCloneByArgs(args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) RefreshFlowNetworkClonesConnections() (*bool, error) {
	var flowNetworkClones []*model.FlowNetworkClone
	d.DB.Where("is_master_slave IS NOT TRUE").Find(&flowNetworkClones)
	for _, fnc := range flowNetworkClones {
		accessToken, err := client.GetFlowToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowUsername, *fnc.FlowPassword)
		fncModel := model.FlowNetworkClone{}
		if err != nil {
			fncModel.IsError = utils.NewTrue()
			fncModel.ErrorMsg = utils.NewStringAddress(err.Error())
			fncModel.FlowToken = fnc.FlowToken
		} else {
			fncModel.IsError = utils.NewFalse()
			fncModel.ErrorMsg = nil
			fncModel.FlowToken = accessToken
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fnc).Select("IsError", "ErrorMsg", "FlowToken").Updates(fncModel).Error; err != nil {
			log.Error(err)
		}
	}
	return utils.NewTrue(), nil
}
