package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
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
	if body.ProducerThingUUID == "" {
		return nil, errors.New("please pass in a producer_thing_uuid i.e. uuid of that class")
	}
	if body.ProducerThingClass == "" {
		return nil, errors.New("please pass in a producer_thing_class i.e. point, job etc")
	}
	producerThingName := ""

	switch body.ProducerThingClass {
	case model.ThingClass.Point:
		pnt, err := d.GetPoint(body.ProducerThingUUID, api.Args{})
		if err != nil {
			return nil, errors.New("point not found, please supply a valid point producer_thing_uuid")
		}
		producerThingName = pnt.Name
	case model.ThingClass.Schedule:
		sch, err := d.GetSchedule(body.ProducerThingUUID)
		if err != nil {
			return nil, errors.New("schedule not found, please supply a valid point producer_thing_uuid")
		}
		producerThingName = sch.Name
	default:
		return nil, errors.New("we are not supporting producer_thing_class other than point for now")
	}
	_, err := d.GetStream(body.StreamUUID, api.Args{})
	if err != nil {
		return nil, newError("GetStream", "error on trying to get validate the gateway UUID")
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Producer)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID, _ = utils.MakeUUID()
	body.ProducerThingName = nameIsNil(producerThingName)
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, newError("CreateProducer", "error on trying to add a new Producer")
	}
	return body, nil
}

func (d *GormDatabase) UpdateProducer(uuid string, body *model.Producer) (*model.Producer, error) {
	var producerModel *model.Producer
	if err := d.DB.Where("uuid = ?", uuid).Find(&producerModel).Error; err != nil {
		return nil, err
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&producerModel, body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&producerModel).Updates(body).Error; err != nil {
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
		producerModel.CurrentWriterUUID = point.UUID
		pointUUID := point.UUID
		pro, err := d.GetProducerByField("producer_thing_uuid", pointUUID)
		if err != nil {
			return "", errors.New("no point for this producer was not found")
		}
		_, err = d.UpdateProducer(pro.UUID, &producerModel)
		if err != nil {
			log.Errorf("producer: issue on update producer err: %v\n", err)
			return "", errors.New("issue on update producer")
		}
		value := point.PresentValue
		point.Priority.P16 = value
		b, err := json.Marshal(point)
		if err != nil {
			log.Errorf("producer: on update write history for point err: %v\n", err)
			return "", errors.New("issue on update write history for point")
		}
		_, err = d.ProducerWriteHist(pro.UUID, b)
		if err != nil {
			log.Errorf("producer: issue on write history ProducerWriteHist: %v\n", err)
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
