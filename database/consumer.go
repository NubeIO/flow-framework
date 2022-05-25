package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
	fnc, err := d.GetFlowNetworkClone(flowNetworkCloneUUID, api.Args{})
	if err != nil {
		return nil, newError("GetFlowNetworkClone", "error on trying to get validate flow_network_clone from stream_clone_uuid")
	}
	cli := client.NewFlowClientCliFromFNC(fnc)
	rawProducer, err := cli.GetQueryMarshal(urls.SingularUrl(urls.ProducerUrl, body.ProducerUUID), model.Producer{})
	if err != nil {
		return nil, err
	}
	producer := rawProducer.(*model.Producer)
	if streamClone.SourceUUID != producer.StreamUUID {
		return nil, newError("Validation failure", "consumer stream_clones & producer stream are different source of truth")
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Consumer)
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
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error) {
	var consumerModel *model.Consumer
	if err := d.DB.Where("uuid = ?", uuid).First(&consumerModel).Error; err != nil {
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

func (d *GormDatabase) DeleteConsumers(args api.Args) (bool, error) {
	var consumerModel *model.Consumer
	query := d.buildConsumerQuery(args)
	query = query.Delete(&consumerModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncConsumerWriters(uuid string) []*interfaces.SyncModel {
	consumer, _ := d.GetConsumer(uuid, api.Args{WithWriters: true})
	var outputs []*interfaces.SyncModel
	for _, writer := range consumer.Writers {
		_, err := d.UpdateWriter(writer.UUID, writer)
		var output interfaces.SyncModel
		if err != nil {
			output = interfaces.SyncModel{UUID: writer.UUID, IsError: true, Message: nstring.New(err.Error())}
		} else {
			output = interfaces.SyncModel{UUID: writer.UUID, IsError: false}
		}
		outputs = append(outputs, &output)
	}
	return outputs
}
