package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


var gatewaySubscriberChildTable = "Subscriber"
var gatewaySubscriptionChildTable = "Subscription"

// GetGateways get all of them
func (d *GormDatabase) GetGateways(withChildren bool) ([]*model.Stream, error) {
	var gatewaysModel []*model.Stream
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload(gatewaySubscriberChildTable).Preload(gatewaySubscriptionChildTable).Find(&gatewaysModel);if query.Error != nil {
			return nil, query.Error
		}
		return gatewaysModel, nil
	} else {
		query := d.DB.Find(&gatewaysModel);if query.Error != nil {
			return nil, query.Error
		}
		return gatewaysModel, nil
	}

}


// CreateGateway make it
func (d *GormDatabase) CreateGateway(body *model.Stream)  error {
	var gatewayModel []model.Stream
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Stream)
	if !body.IsRemote  {
		query := d.DB.Where("is_remote = ?", 0).First(&gatewayModel) //if existing local network then don't create it
		r := query.RowsAffected
		if r != 0 {
			return errorMsg("network", "a local gateway exists", nil)
		}
		body.Name = "Local rubix"
		body.Enable = true
		body.Description = "Local rubix gateway for sending data between jobs and points"
	}
	n := d.DB.Create(body).Error
	return n
}

// GetGateway get it
func (d *GormDatabase) GetGateway(uuid string) (*model.Stream, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("uuid = ? ", uuid).First(&gatewayModel); if query.Error != nil {
		return nil, query.Error
	}
	return gatewayModel, nil
}

// DeleteGateway deletes it
func (d *GormDatabase) DeleteGateway(uuid string) (bool, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("uuid = ? ", uuid).Delete(&gatewayModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateGateway  update it
func (d *GormDatabase) UpdateGateway(uuid string, body *model.Stream) (*model.Stream, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("uuid = ?", uuid).Find(&gatewayModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&gatewayModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return gatewayModel, nil

}

// DropGateways delete all.
func (d *GormDatabase) DropGateways() (bool, error) {
	var gatewayModel *model.Stream
	query := d.DB.Where("1 = 1").Delete(&gatewayModel)
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
