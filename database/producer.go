package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"time"
)

type Producer struct {
	*model.Producer
}

func (d *GormDatabase) GetProducers(args api.Args) ([]*model.Producer, error) {
	var producersModel []*model.Producer
	query := d.buildProducerQuery(args)
	if err := query.Find(&producersModel).Error; err != nil {
		return nil, err
	}
	return producersModel, nil
}

func (d *GormDatabase) GetProducer(uuid string, args api.Args) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.buildProducerQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&producerModel).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) CreateProducer(body *model.Producer) (*model.Producer, error) {
	_, err := d.GetStream(body.StreamUUID, api.Args{})
	if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the gateway UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, errorMsg("CreateProducer", "error on trying to add a new Producer", nil)
	}
	return body, nil
}

func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer, updateHist bool) (*model.Producer, error) {
	var producerModel *model.Producer
	if err := d.DB.Where("uuid = ?", uuid).Find(&producerModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return producerModel, nil
}

func (d *GormDatabase) DeleteProducer(uuid string) (bool, error) {
	var producerModel *model.Producer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&producerModel)
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

type Point struct {
	Priority model.Priority `json:"priority"`
}

// ProducerWriteHist  update it
func (d *GormDatabase) ProducerWriteHist(uuid string, writeData datatypes.JSON) (*model.ProducerHistory, error) {
	ph := new(model.ProducerHistory)
	ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
	ph.ProducerUUID = uuid
	ph.DataStore = writeData
	ph.Timestamp = time.Now().UTC()
	_, err := d.CreateProducerHistory(ph)
	if err != nil {
		return nil, err
	}
	return ph, nil
}

// ProducerWrite  update it
func (d *GormDatabase) ProducerWrite(thingType string, payload interface{}) (string, error) {
	var producerModel model.Producer
	p, err := eventbus.DecodeBody(thingType, payload)
	if err != nil {
		return "", err
	}
	if thingType == model.ThingClass.Point {
		point := p.(*model.Point)
		producerModel.ThingWriterUUID = point.UUID
		pointUUID := point.UUID
		pro, err := d.GetProducerByField("producer_thing_uuid", pointUUID)
		if err != nil {
			return "", errors.New("no point for this producer was not found")
		}
		_, err = d.UpdateProducer(pointUUID, &producerModel, false)
		if err != nil {
			log.Errorf("UpdateProducer")
			return "", errors.New("issue on update producer")
		}
		value := point.PresentValue
		point.Priority.P1 = &value
		b, err := json.Marshal(point)
		if err != nil {
			return "", errors.New("issue on update write history for point")
		}
		_, err = d.ProducerWriteHist(pro.UUID, b)
		if err != nil {
			return "", errors.New("issue on write history for point")
		}
		var proBody model.ProducerBody
		proBody.StreamUUID = pro.StreamUUID
		proBody.ProducerUUID = pro.UUID
		proBody.ThingType = pro.ProducerThingType
		proBody.ThingType = point.ThingType
		proBody.Payload = point

		err = d.producerBroadcast(proBody)
		if err != nil {
			return "", errors.New("issue on producer broadcast")
		}

		return pro.UUID, err
	}
	return "", err
}

// GetProducerByField returns the point for the given field ie name or nil.
//for example get a producer by its producer_thing_uuid
func (d *GormDatabase) GetProducerByField(field string, value string) (*model.Producer, error) {
	var producerModel *model.Producer
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).First(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}

// UpdateProducerByField get by field and update.
//for example update a producer by its producer_thing_uuid
func (d *GormDatabase) UpdateProducerByField(field string, value string, body *model.Producer, updateHist bool) (*model.Producer, error) {
	var producerModel *model.Producer
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Find(&producerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}
