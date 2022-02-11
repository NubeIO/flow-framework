package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
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

func (d *GormDatabase) ProducerWriteHist(uuid string, writeData datatypes.JSON) (*model.ProducerHistory, error) {
	ph := new(model.ProducerHistory)
	ph.ProducerUUID = uuid
	ph.DataStore = writeData
	_, err := d.CreateProducerHistory(ph)
	if err != nil {
		return nil, err
	}
	return ph, nil
}

func (d *GormDatabase) ProducerWrite(point model.Point) error {
	producerModel := new(model.Producer)
	producerModel.CurrentWriterUUID = point.UUID
	// TODO: replace GetProducerByField by GetOneProducerByArgs
	producer, err := d.GetProducerByField("producer_thing_uuid", point.UUID)
	if err != nil {
		return errors.New("producer doesn't exist for this point")
	}
	producerModel, err = d.UpdateProducer(producer.UUID, producerModel)
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}

	syncWriterCOV := model.SyncWriterCOV{
		OriginalValue:   point.OriginalValue,
		PresentValue:    point.PresentValue,
		CurrentPriority: point.CurrentPriority,
		Priority:        point.Priority,
	}
	err = d.TriggerCOVToWriterClone(producerModel, syncWriterCOV)
	if err != nil {
		return err
	}

	ph := new(model.ProducerHistory)
	ph.ProducerUUID = producer.UUID
	b, err := json.Marshal(point.Priority)
	if err != nil {
		log.Errorf("producer: on update write history for point err: %v\n", err)
		return errors.New("issue on update write history for point")
	}
	ph.DataStore = b
	_, err = d.CreateProducerHistory(ph)
	if err != nil {
		log.Errorf("producer: issue on write history ProducerWriteHist: %v\n", err)
		return errors.New("issue on write history for point")
	}
	return nil
}

func (d *GormDatabase) TriggerCOVToWriterClone(producer *model.Producer, body model.SyncWriterCOV) error {
	wcs, err := d.GetWriterClones(api.Args{ProducerUUID: utils.NewStringAddress(producer.UUID)})
	if err != nil {
		return errors.New("error on getting writer clones from producer_uuid")
	}
	for _, wc := range wcs {
		_ = d.TriggerCOVFromWriterCloneToWriter(producer, *wc, body)
	}
	return nil
}

func (d *GormDatabase) TriggerCOVFromWriterCloneToWriter(producer *model.Producer, wc model.WriterClone, body model.SyncWriterCOV) error {
	stream, _ := d.GetStream(producer.StreamUUID, api.Args{WithFlowNetworks: true})
	body.WriterUUID = wc.SourceUUID
	for _, fn := range stream.FlowNetworks {
		// TODO: wc.FlowFrameworkUUID == "" remove from condition; it's here coz old deployment doesn't used to have that value
		if wc.FlowFrameworkUUID == "" || fn.UUID == wc.FlowFrameworkUUID {
			cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))
			_ = cli.SyncWriterCOV(&body)
		}
	}
	return nil
}

// GetProducerByField returns the point for the given field ie name or nil.
// For example get a producer by its producer_thing_uuid
func (d *GormDatabase) GetProducerByField(field string, value string) (*model.Producer, error) {
	var producerModel *model.Producer
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).First(&producerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}
