package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	log "github.com/sirupsen/logrus"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

type Writer struct {
	*model.Writer
}

func (d *GormDatabase) GetWriters(args api.Args) ([]*model.Writer, error) {
	var writers []*model.Writer
	query := d.buildWriterQuery(args)
	if err := query.Find(&writers).Error; err != nil {
		return nil, query.Error
	}
	return writers, nil
}

func (d *GormDatabase) CreateWriter(body *model.Writer) (*model.Writer, error) {
	name := ""
	switch body.WriterThingClass {
	case model.ThingClass.Point:
		point, err := d.GetPoint(body.WriterThingUUID, api.Args{})
		if err != nil {
			return nil, errors.New("point not found, please supply a valid point writer_thing_uuid")
		}
		name = point.Name
	case model.ThingClass.Schedule:
		fmt.Println(body.WriterThingUUID)
		schedule, err := d.GetSchedule(body.WriterThingUUID)
		if err != nil {
			return nil, errors.New("schedule not found, please supply a valid point writer_thing_uuid")
		}
		name = schedule.Name
	default:
		return nil, errors.New("we are not supporting writer_thing_uuid other than point for now")
	}

	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Writer)
	body.WriterThingName = name
	body.SyncUUID, _ = utils.MakeUUID()
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	// ignore err coz CASCADE delete is there, so there must be data
	consumer, _ := d.GetConsumer(body.ConsumerUUID, api.Args{})
	streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
	fn, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))
	syncWriterBody := model.SyncWriter{
		Writer:       *body,
		ProducerUUID: consumer.ProducerUUID,
	}
	_, err := cli.SyncWriter(&syncWriterBody)
	if err != nil {
		log.Error(err)
	}
	return body, nil
}

// GetWriter get it
func (d *GormDatabase) GetWriter(uuid string) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

// GetWriterByThing get it by its thing uuid
func (d *GormDatabase) GetWriterByThing(producerThingUUID string) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil
}

// DeleteWriter deletes it
func (d *GormDatabase) DeleteWriter(uuid string) (bool, error) {
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&writerModel)
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

// UpdateWriter  update it
func (d *GormDatabase) UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error) {
	var writerModel *model.Writer
	query := d.DB.Where("uuid = ?", uuid).Find(&writerModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&writerModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return writerModel, nil

}

// DropWriters delete all.
func (d *GormDatabase) DropWriters() (bool, error) {
	var writerModel *model.Writer
	query := d.DB.Where("1 = 1").Delete(&writerModel)
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

// CreateWriterWizard add a new consumer to an existing producer and add a new writer and writer clone
// use the flow-network UUID
func (d *GormDatabase) CreateWriterWizard(body *api.WriterWizard) (bool, error) {
	/*TODO: Binod
	var consumerModel model.Consumer
	var writerModel model.Writer
	var writerCloneModel model.WriterClone
	var session *client.FlowClient
	producer := new(model.Producer)
	flow, err := d.GetFlowNetwork(body.ConsumerFlowUUID, api.Args{})
	if err != nil {
		return false, err
	}
	isRemote := utils.BoolIsNil(flow.IsRemote)
	if isRemote {
		session = client.NewSessionWithToken("", flow.FlowIP, flow.FlowPort)
		pro, err := session.GetProducer(body.ProducerUUID)
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to remote GetProducer")
		}
		producer.UUID = pro.UUID
		producer.ProducerThingClass = pro.ProducerThingClass
		producer.ProducerThingType = pro.ProducerThingType
		producer.ProducerThingName = pro.ProducerThingName
		producer.ProducerThingUUID = pro.ProducerThingUUID
	} else {
		pro, err := d.GetProducer(body.ProducerUUID, api.Args{})
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to local GetProducer")
		}
		producer.UUID = pro.UUID
		producer.ProducerThingClass = pro.ProducerThingClass
		producer.ProducerThingType = pro.ProducerThingType
		producer.ProducerThingName = pro.ProducerThingName
		producer.ProducerThingUUID = pro.ProducerThingUUID

	}

	consumerModel.StreamUUID = body.ConsumerStreamUUID
	consumerModel.Name = "consumer stream"
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ProducerThingClass = producer.ProducerThingClass
	consumerModel.ProducerThingType = producer.ProducerThingType
	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumerModel.ProducerThingUUID = producer.ProducerThingUUID
	consumerModel.ProducerThingName = producer.ProducerThingName
	_, err = d.CreateConsumer(&consumerModel)
	if err != nil {
		log.Errorf("wizard:  CreateConsumer: %v\n", err)
		return false, errors.New("CREATE-WRITER: failed to local CreateConsumer")
	}
	//// writer
	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.ConsumerThingUUID = consumerModel.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return false, errors.New("CREATE-WRITER: failed to local CreateWriter")
	}
	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = body.ProducerUUID
	writerCloneModel.WriterUUID = writer.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API
	if !isRemote {
		writerClone, err := d.CreateWriterClone(&writerCloneModel)
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to local CreateWriterClone")
		}
		writerModel.CloneUUID = writerClone.UUID
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to local UpdateWriter")
		}
	} else {
		clone, err := session.CreateWriterClone(writerCloneModel)
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to remote CreateWriterClone")
		}
		writerModel.CloneUUID = clone.UUID
		writerModel.CloneUUID = clone.UUID
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			return false, errors.New("CREATE-WRITER: failed to local UpdateWriter")
		}
	}*/
	return true, nil
}

/*
WriterAction read or write a value to the writer and onto the writer clone
1. update writer
2. write to writerClone
3. producer to decide if it's a valid cov
3. producer to write to history
4. return to consumer and update as required if it's a valid cov (update the writerConeUUID, this could be from another flow-framework instance)
*/
func (d *GormDatabase) WriterAction(uuid string, body *model.WriterBody) (*model.ProducerHistory, error) {
	askRefresh := body.AskRefresh
	writer, err := d.GetWriter(uuid)
	if err != nil {
		return nil, err
	}
	data, action, err := d.validateWriterBody(writer.WriterThingClass, body)
	if err != nil {
		return nil, err
	}
	consumer, _ := d.GetConsumer(writer.ConsumerUUID, api.Args{})
	streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
	fnc, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	producerUUID := consumer.ProducerUUID
	cli := client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
	var pHistory *model.ProducerHistory
	if action == model.CommonNaming.Write {
		pHistoryModel := model.ProducerHistory{
			ProducerUUID: producerUUID,
			CommonCurrentWriterUUID: model.CommonCurrentWriterUUID{
				CurrentWriterUUID: writer.UUID,
			},
			DataStore: data,
		}
		if _, e := cli.AddProducerHistory(pHistoryModel); e != nil {
			return nil, e
		}
		pHistory = &pHistoryModel
	} else {
		pHistory, err = cli.GetProducerHistory(producerUUID)
		if err != nil {
			return nil, err
		}
	}
	if askRefresh {
		writer.DataStore = pHistory.DataStore
		consumer.CurrentWriterUUID = writer.UUID
		if writer.WriterThingClass == model.ThingClass.Point {
			priority := new(model.Priority)
			_ = json.Unmarshal(writer.DataStore, &priority)
			highestPriorityValue := priority.GetHighestPriorityValue()
			d.DB.Model(&model.Point{}).Where("uuid = ?", writer.WriterThingUUID).
				Updates(map[string]interface{}{
					"present_value":  highestPriorityValue,
					"original_value": highestPriorityValue,
				})
		} else if writer.WriterThingClass == model.ThingClass.Schedule {
			log.Println("IS SCH", writer.WriterThingUUID)
			// TODO: for schedule and others
		}
		d.DB.Model(&writer).Updates(writer)
		d.DB.Model(&consumer).Updates(consumer)
	}
	return pHistory, nil
}

func (d *GormDatabase) WriterBulkAction(body []*model.WriterBulk) (*utils.Array, error) {
	arr := utils.NewArray()
	for _, wri := range body {
		b := new(model.WriterBody)
		b.Action = wri.Action
		b.AskRefresh = wri.AskRefresh
		b.Priority = wri.Priority
		action, err := d.WriterAction(wri.WriterUUID, b)
		if err == nil {
			arr.Add(action)
		}
	}
	return arr, nil

}

func consumerRefresh(producerFeedback *model.ProducerHistory) (*model.Consumer, error) {
	updateConsumer := new(model.Consumer)
	updateConsumer.DataStore = producerFeedback.DataStore
	updateConsumer.CurrentWriterUUID = producerFeedback.CurrentWriterUUID
	return updateConsumer, nil
}

func (d *GormDatabase) validateWriterBody(thingClass string, body *model.WriterBody) ([]byte, string, error) {
	if thingClass == model.ThingClass.Point {
		var bk model.WriterBody
		if body.Action == model.WriterActions.Write {
			if body.Priority == bk.Priority {
				return nil, body.Action, errors.New("error: invalid json on writer body")
			}
			b, err := json.Marshal(body.Priority)
			if err != nil {
				return nil, body.Action, errors.New("error: failed to marshal priority on write body")
			}
			return b, body.Action, err
		} else if body.Action == model.WriterActions.Read {
			return nil, body.Action, nil
		} else {
			return nil, body.Action, errors.New("error: invalid action, try read or write")
		}
	} else if thingClass == model.ThingClass.Schedule {
		if body.Action == model.WriterActions.Write {
			b, err := json.Marshal(body.Schedule.Schedules)
			if err != nil {
				return nil, body.Action, errors.New("error: failed to marshal schedule on write body")
			}
			return b, body.Action, err
		} else if body.Action == model.WriterActions.Read {
			return nil, body.Action, nil
		} else {
			return nil, body.Action, errors.New("error: invalid action, try read or write")
		}
	}
	return nil, body.Action, errors.New("error: invalid data type on writer body, i.e. type could be a point")
}
