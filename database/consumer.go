package database

import (
	"errors"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
)

type Consumers struct {
	*model.Consumer
}

func (d *GormDatabase) GetConsumers(args api.Args) ([]*model.Consumer, error) {
	var consumersModel []*model.Consumer
	query := d.buildConsumerQuery(args)
	if err := query.Find(&consumersModel).Error; err != nil {
		return nil, err
	}
	return consumersModel, nil
}

func (d *GormDatabase) GetConsumer(uuid string, args api.Args) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	query := d.buildConsumerQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&consumerModel).Error; err != nil {
		return nil, err
	}
	return consumerModel, nil
}

func (d *GormDatabase) CreateConsumer(body *model.Consumer) (*model.Consumer, error) {
	streamClone, err := d.GetStreamClone(body.StreamCloneUUID, api.Args{})
	if err != nil {
		return nil, newError("GetStreamClone", "error on trying to get validate the stream_clone_uuid")
	}
	flowNetworkCloneUUID := streamClone.FlowNetworkCloneUUID
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{UUID: &flowNetworkCloneUUID})
	if err != nil {
		return nil, newError("GetOneFlowNetworkCloneByArgs", "error on trying to get validate flow_network_clone from stream_clone_uuid")
	}
	cli := client.NewSessionWithToken(fnc.FlowToken, fnc.FlowIP, fnc.FlowPort)
	producer, err := cli.GetProducer(body.ProducerUUID)
	if err != nil {
		return nil, newError("GetProducer", "error on finding producer by producer_uuid")
	}
	if streamClone.SourceUUID != producer.StreamUUID {
		return nil, newError("Validation failure", "consumer stream_clones & producer stream are different source of truth")
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Consumer)
	body.Name = nameIsNil(body.Name)
	body.SyncUUID = producer.SyncUUID
	body.ProducerThingName = producer.ProducerThingName
	body.ProducerThingUUID = producer.ProducerThingUUID
	body.ProducerThingClass = producer.ProducerThingClass
	body.ProducerThingType = producer.ProducerThingType
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

func (d *GormDatabase) DeleteConsumer(uuid string) (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("uuid = ? ", uuid).Delete(&consumerModel)
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

func (d *GormDatabase) UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	if err := d.DB.Where("uuid = ?", uuid).Find(&consumerModel).Error; err != nil {
		return nil, err
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&consumerModel, body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&consumerModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return consumerModel, nil

}

func (d *GormDatabase) DropConsumers() (bool, error) {
	var consumerModel *model.Consumer
	query := d.DB.Where("1 = 1").Delete(&consumerModel)
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

/*
add new consumer auto add, writer and writer clone for a remote network and local
body
-- consumerStreamUUID
-- producerUUID

- get the streamUUID from the producerUUID
- from the consumerStreamUUID get the flowUUID
- first make sure that the producer device is online and the streamUUID is valid
- work out if the producer is local or remote
- add the new consumer, writer and writerClone
*/

func (d *GormDatabase) AddConsumerWizard(consumerStreamUUID, producerUUID string, consumerModel *model.Consumer) (*model.Consumer, error) {
	streamUUID := consumerStreamUUID
	var writerModel model.Writer
	var writerCloneModel model.WriterClone

	if producerUUID == "" {
		return nil, errors.New("error: no producer uuid provided")
	}
	if consumerStreamUUID == "" {
		return nil, errors.New("error: no stream uuid provided")
	}
	stream, flow, err := d.GetFlowUUID(streamUUID)
	if err != nil || stream.UUID == "" {
		return nil, errors.New("error: invalid stream UUID")
	}
	isRemote := utils.BoolIsNil(flow.IsRemote)

	var producer *model.Producer
	if isRemote {
		cli := client.NewSessionWithToken(flow.FlowToken, flow.FlowIP, flow.FlowPort)
		p, err := cli.GetProducer(producerUUID)
		if err != nil {
			return nil, errors.New("error: issue on get producer over rest client")
		}
		producer = p
	} else {
		p, err := d.GetProducer(producerUUID, api.Args{})
		if err != nil {
			return nil, errors.New("error: issue on get producer")
		}
		producer = p
	}

	if producer.UUID == "" {
		return nil, errors.New("error: no producer producer found with that UUID")
	}

	consumerModel.StreamCloneUUID = stream.UUID // TODO: BIOND
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ProducerThingName = producer.ProducerThingName
	consumerModel.ProducerThingUUID = producer.ProducerThingUUID
	consumerModel.ProducerThingClass = producer.ProducerThingClass
	consumerModel.ProducerThingType = producer.ProducerThingType

	consumerModel.ConsumerApplication = model.CommonNaming.Mapping
	consumer, err := d.CreateConsumer(consumerModel)
	if err != nil {
		return nil, errors.New("error: issue on create consumer")
	}
	// writer
	writerModel.ConsumerUUID = consumer.UUID
	//writerModel.ConsumerThingUUID = consumerModel.UUID //TODO: BINOD
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return nil, errors.New("error: issue on create writer")
	}
	// add consumer to the writerClone
	writerCloneModel.ProducerUUID = producer.UUID
	writerCloneModel.SourceUUID = writer.UUID
	writerModel.WriterThingClass = model.ThingClass.Point
	writerModel.WriterThingType = model.ThingClass.API

	if !isRemote {
		_, err := d.CreateWriterClone(&writerCloneModel)
		if err != nil {
			return nil, errors.New("error: issue on create writer clone over rest")
		}
		//update writerCloneUUID to writer
		//writerModel.CloneUUID = writerClone.UUID //TODO: Binod
		_, err = d.UpdateWriter(writerModel.UUID, &writerModel)
		if err != nil {
			return nil, errors.New("error: issue on update writer over rest")
		}
	} else {
		cli := client.NewSessionWithToken(flow.FlowToken, flow.FlowIP, flow.FlowPort)
		_, err := cli.CreateWriterClone(writerCloneModel)
		if err != nil {
			return nil, errors.New("error: issue on create writer clone")
		}
		//writerModel.CloneUUID = clone.UUID //TODO: Binod
		_, err = cli.EditWriter(writerModel.UUID, writerModel, false)
		if err != nil {
			return nil, errors.New("error: issue on update writer")
		}
	}
	return consumerModel, nil
}
