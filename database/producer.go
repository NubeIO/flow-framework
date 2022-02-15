package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
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

func (d *GormDatabase) GetOneProducerByArgs(args api.Args) (*model.Producer, error) {
	var producerModel *model.Producer
	query := d.buildProducerQuery(args)
	if err := query.First(&producerModel).Error; err != nil {
		return nil, query.Error
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

func (d *GormDatabase) ProducersPointWrite(point *model.Point) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = point.UUID // TODO: check current_writer_uuid is needed or not
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &point.UUID})
	for _, producer := range producers {
		err := d.producerPointWrite(producer.UUID, point, producerModelBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerPointWrite(uuid string, point *model.Point, producerModelBody *model.Producer) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody) // TODO: check current_writer_uuid is needed or not
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}

	syncCOV := model.SyncCOV{Priority: point.Priority}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(point.Priority)
	ph := new(model.ProducerHistory)
	ph.ProducerUUID = uuid
	ph.DataStore = data
	_, err = d.CreateProducerHistory(ph)
	if err != nil {
		log.Errorf("producer: issue on write history ProducerWriteHist: %v\n", err)
		return errors.New("issue on write history for point")
	}
	return nil
}

func (d *GormDatabase) ProducersScheduleWrite(uuid string, body *model.ScheduleData) error {
	producerModelBody := new(model.Producer)
	producerModelBody.CurrentWriterUUID = uuid // TODO: check current_writer_uuid is needed or not
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &uuid})
	for _, producer := range producers {
		err := d.producerScheduleWrite(producer.UUID, body, producerModelBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) producerScheduleWrite(uuid string, scheduleData *model.ScheduleData, producerModelBody *model.Producer) error {
	producerModel, err := d.UpdateProducer(uuid, producerModelBody) // TODO: check current_writer_uuid is needed or not
	if err != nil {
		log.Errorf("producer: issue on update producer err: %v\n", err)
		return errors.New("issue on update producer")
	}

	syncCOV := model.SyncCOV{Schedule: scheduleData}
	err = d.TriggerCOVToWriterClone(producerModel, &syncCOV)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) TriggerCOVToWriterClone(producer *model.Producer, body *model.SyncCOV) error {
	wcs, err := d.GetWriterClones(api.Args{ProducerUUID: utils.NewStringAddress(producer.UUID)})
	if err != nil {
		return errors.New("error on getting writer clones from producer_uuid")
	}
	for _, wc := range wcs {
		_ = d.TriggerCOVFromWriterCloneToWriter(producer, wc, body)
	}
	return nil
}

func (d *GormDatabase) TriggerCOVFromWriterCloneToWriter(producer *model.Producer, wc *model.WriterClone, body *model.SyncCOV) error {
	stream, _ := d.GetStream(producer.StreamUUID, api.Args{WithFlowNetworks: true})
	body.WriterUUID = wc.SourceUUID
	for _, fn := range stream.FlowNetworks {
		// TODO: wc.FlowFrameworkUUID == "" remove from condition; it's here coz old deployment doesn't used to have that value
		if wc.FlowFrameworkUUID == "" || fn.UUID == wc.FlowFrameworkUUID {
			cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))
			_ = cli.SyncCOV(body)
		}
	}
	return nil
}
