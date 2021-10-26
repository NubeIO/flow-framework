package database

import (
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) GetLocalStorageFlowNetwork() (*model.LocalStorageFlowNetwork, error) {
	var localStorageFlowNetwork *model.LocalStorageFlowNetwork
	if err := d.DB.First(&localStorageFlowNetwork).Error; err != nil {
		return nil, err
	}
	return localStorageFlowNetwork, nil
}

func (d *GormDatabase) UpdateLocalStorageFlowNetwork(body *model.LocalStorageFlowNetwork) (*model.LocalStorageFlowNetwork, error) {
	var lsfn *model.LocalStorageFlowNetwork
	if err := d.DB.First(&lsfn).Error; err != nil {
		return nil, err
	}
	token, err := GetFlowToken(body.FlowIP, body.FlowPort, body.FlowUsername, body.FlowPassword)
	if err != nil {
		return nil, err
	}
	body.FlowToken = *token
	if err := d.DB.Model(&lsfn).Updates(&body).Error; err != nil {
		return nil, err
	}
	return lsfn, nil
}

func (d *GormDatabase) RefreshLocalStorageFlowToken() (*bool, error) {
	var lsfn *model.LocalStorageFlowNetwork
	if err := d.DB.First(&lsfn).Error; err != nil {
		return nil, err
	}
	conf := config.Get()
	token, err := GetFlowToken(conf.Server.ListenAddr, conf.Server.Port, lsfn.FlowUsername, lsfn.FlowPassword)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Model(&lsfn).Updates(model.LocalStorageFlowNetwork{FlowToken: *token}).Error; err != nil {
		return nil, err
	}
	return utils.NewTrue(), nil
}
