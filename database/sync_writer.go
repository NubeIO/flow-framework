package database

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	writerClone := model.WriterClone{}
	producer, err := d.GetProducer(body.ProducerUUID, api.Args{})
	if err != nil {
		return nil, errors.New("producer does not exist")
	}
	writerClone.WriterThingName = producer.ProducerThingName
	writerClone.WriterThingUUID = producer.ProducerThingUUID
	writerClone.WriterThingClass = producer.ProducerThingClass
	writerClone.WriterThingType = producer.ProducerThingType
	writerClone.ProducerUUID = body.ProducerUUID
	writerClone.SourceUUID = body.WriterUUID
	writerClone.FlowFrameworkUUID = body.FlowFrameworkUUID
	var writerCloneModel []*model.WriterClone
	if err = d.DB.Where("source_uuid = ? ", writerClone.SourceUUID).Find(&writerCloneModel).Error; err != nil {
		return nil, err
	}
	if len(writerCloneModel) == 0 {
		writerClone.UUID = nuuid.MakeTopicUUID(model.CommonNaming.WriterClone)
		if err = d.DB.Create(&writerClone).Error; err != nil {
			return nil, err
		}
	} else {
		writerClone.UUID = writerCloneModel[0].UUID
		if err = d.DB.Model(&writerCloneModel[0]).Updates(writerClone).Error; err != nil {
			return nil, err
		}
	}
	return &writerClone, nil
}

func (d *GormDatabase) SyncCOV(writerUUID string, body *model.SyncCOV) error {
	writer, err := d.GetWriter(writerUUID)
	if err != nil {
		return err
	}
	uuid := writer.WriterThingUUID
	if writer.WriterThingClass == model.ThingClass.Point {
		pointModel := model.PointWriter{
			Priority: body.Priority,
		}
		_, _, _, _, err = d.PointWrite(uuid, &pointModel, false)
		return err
	} else {
		return d.ScheduleWrite(writer.WriterThingUUID, body.Schedule)
	}
}

func (d *GormDatabase) SyncWriterWriteAction(sourceUUID string, body *model.SyncWriterAction) error {
	writerClone, err := d.GetOneWriterCloneByArgs(api.Args{SourceUUID: &sourceUUID})
	if err != nil {
		return err
	}
	if writerClone.WriterThingClass == model.ThingClass.Point {
		data, _ := json.Marshal(body.Priority)
		writerCloneBody := model.WriterClone{CommonWriter: model.CommonWriter{DataStore: data}}
		err = d.updateWriterClone(writerClone.UUID, &writerCloneBody)
		if err != nil {
			return nil
		}
		pointWriter := model.PointWriter{Priority: body.Priority}
		// TODO: change this section by below commented section
		producer, _ := d.GetProducer(writerClone.ProducerUUID, api.Args{})
		_, _, _, _, err = d.PointWrite(producer.ProducerThingUUID, &pointWriter, false)
		// Currently, writerClone.WriterThingUUID has not valid `WriterThingUUID` on old deployments
		// _, err = d.PointWrite(writerClone.WriterThingUUID, &point, true)
		return err
	} else if writerClone.WriterThingClass == model.ThingClass.Schedule {
		data, _ := json.Marshal(body.Schedule)
		writerCloneBody := model.WriterClone{CommonWriter: model.CommonWriter{DataStore: data}}
		err = d.updateWriterClone(writerClone.UUID, &writerCloneBody)
		if err != nil {
			return nil
		}
		// TODO: change this section by below commented section
		producer, _ := d.GetProducer(writerClone.ProducerUUID, api.Args{})
		err = d.ScheduleWrite(producer.ProducerThingUUID, body.Schedule)
		// Currently, writerClone.WriterThingUUID has not valid `WriterThingUUID` on old deployments
		// err = d.ScheduleWrite(writerClone.WriterThingUUID, body.Schedule)
		return err
	} else {
		return errors.New("no match writer thing class")
	}
}

func (d *GormDatabase) SyncWriterReadAction(sourceUUID string) error {
	writerClone, err := d.GetOneWriterCloneByArgs(api.Args{SourceUUID: &sourceUUID})
	if err != nil {
		return nil
	}
	producer, err := d.GetProducer(writerClone.ProducerUUID, api.Args{})
	if err != nil {
		return nil
	}
	syncCOV := model.SyncCOV{}
	if writerClone.WriterThingClass == model.ThingClass.Point {
		point, err := d.GetPoint(writerClone.WriterThingUUID, api.Args{WithPriority: true})
		if err != nil {
			return err
		}
		priorityMap := priorityarray.ConvertToMap(*point.Priority)
		syncCOV.Priority = &priorityMap
	} else {
		schedule, err := d.GetSchedule(writerClone.WriterThingUUID)
		if err != nil {
			return err
		}
		var scheduleData *model.ScheduleData
		err = json.Unmarshal(schedule.Schedule, scheduleData)
		if err != nil {
			return nil
		}
		syncCOV.Schedule = scheduleData
	}
	return d.TriggerCOVFromWriterCloneToWriter(producer, writerClone, &syncCOV)
}
