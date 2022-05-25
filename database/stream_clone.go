package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) GetStreamClones(args api.Args) ([]*model.StreamClone, error) {
	var streamClonesModel []*model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetStreamCloneByArg(args api.Args) (*model.StreamClone, error) {
	var streamClonesModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetStreamClone(uuid string, args api.Args) (*model.StreamClone, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamCloneModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamCloneModel, nil
}

func (d *GormDatabase) DeleteStreamClone(uuid string) (bool, error) {
	var streamCloneModel *model.StreamClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneStreamCloneByArgs(args api.Args) (bool, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	if err := query.First(&streamCloneModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncStreamCloneConsumers(uuid string) ([]*interfaces.SyncModel, error) {
	streamClone, _ := d.GetStreamClone(uuid, api.Args{WithConsumers: true})
	if streamClone == nil {
		return nil, errors.New("no stream_clone")
	}
	flowNetworkClone, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
	cli := client.NewFlowClientCliFromFNC(flowNetworkClone)
	var outputs []*interfaces.SyncModel
	for _, consumer := range streamClone.Consumers {
		rawProducer, err := cli.GetQueryMarshal(urls.SingularUrl(urls.ProducerUrl, consumer.ProducerUUID), model.Producer{})
		var output interfaces.SyncModel
		if err != nil {
			output = interfaces.SyncModel{UUID: consumer.UUID, IsError: true, Message: nstring.New(err.Error())}
			consumer.Connection = connection.Broken.String()
		} else {
			output = interfaces.SyncModel{UUID: consumer.UUID, IsError: false}
			producer := rawProducer.(*model.Producer)
			consumer.Connection = connection.Connected.String()
			consumer.ProducerThingName = producer.ProducerThingName
			consumer.ProducerThingUUID = producer.ProducerThingUUID
			consumer.ProducerThingClass = producer.ProducerThingClass
			consumer.ProducerThingType = producer.ProducerThingType
		}
		_, _ = d.UpdateConsumer(consumer.UUID, consumer)
		outputs = append(outputs, &output)
	}
	return outputs, nil
}
