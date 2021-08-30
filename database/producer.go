package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"time"
)

type Producer struct {
	*model.Producer

}

// GetProducers get all of them
func (d *GormDatabase) GetProducers() ([]*model.Producer, error) {
	var producersModel []*model.Producer
	query := d.DB.Preload("ProducerSubscriptionList").Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}



// CreateProducer make it
func (d *GormDatabase) CreateProducer(body *model.Producer) (*model.Producer, error) {
	//call points and make it exists
	_, err := d.GetStreamGateway(body.StreamUUID);if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the gateway UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	err = d.DB.Create(&body).Error; if err != nil {
		return nil, errorMsg("CreateProducer", "error on trying to add a new Producer", nil)
	}
	return body, nil
}



// GetProducer get it
func (d *GormDatabase) GetProducer(uuid string) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// UpdateProducer  update it
func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&producerModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	//TODO move this out of here as we need to also store the subscription UUID
	if producerModel.EnableHistory {
		ph := new(model.ProducerHistory)
		ph.ProducerUUID = uuid
		ph.PresentValue = body.PresentValue
		ph.Timestamp = time.Now().UTC()
		_, err := d.CreateProducerHistory(ph)
		if err != nil {
			return nil, err
		}
	}
	return producerModel, nil
}

// GetProducerByThingUUID get it by its
func (d *GormDatabase) GetProducerByThingUUID(thingUUID string) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.DB.Where("producer_thing_uuid = ? ", thingUUID).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// DeleteProducer deletes it
func (d *GormDatabase) DeleteProducer(uuid string) (bool, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&producerModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// DropProducers delete all.
func (d *GormDatabase) DropProducers() (bool, error) {
	var producerModel *model.Producer
	query := d.DB.Where("1 = 1").Delete(&producerModel)
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