package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	var flowNetworkClonesModel []*model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkClonesModel, nil
}

func (d *GormDatabase) GetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, err

	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) DeleteFlowNetworkClone(uuid string) (bool, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&flowNetworkCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) GetOneFlowNetworkCloneByArgs(args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) DeleteOneFlowNetworkCloneByArgs(args api.Args) (bool, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&flowNetworkCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) RefreshFlowNetworkClonesConnections() (*bool, error) {
	var flowNetworkClones []*model.FlowNetworkClone
	d.DB.Where("is_master_slave IS NOT TRUE AND is_remote IS TRUE").Find(&flowNetworkClones)
	for _, fnc := range flowNetworkClones {
		accessToken, err := client.GetFlowToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowUsername, *fnc.FlowPassword)
		fncModel := model.FlowNetworkClone{}
		if err != nil {
			fncModel.IsError = boolean.NewTrue()
			fncModel.ErrorMsg = nstring.NewStringAddress(err.Error())
			fncModel.FlowToken = fnc.FlowToken
		} else {
			fncModel.IsError = boolean.NewFalse()
			fncModel.ErrorMsg = nil
			fncModel.FlowToken = accessToken
		}
		// here `.Select` is needed because NULL value needs to set on is_error=false
		if err := d.DB.Model(&fnc).Select("IsError", "ErrorMsg", "FlowToken").Updates(fncModel).Error; err != nil {
			log.Error(err)
		}
	}
	return boolean.NewTrue(), nil
}
