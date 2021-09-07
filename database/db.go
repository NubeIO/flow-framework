package database

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/patrickmn/go-cache"
)


// DropAllFlow networks, gateways, commandGroup, consumers, jobs and children.
func (d *GormDatabase) DropAllFlow() (bool, error) {

	//delete networks
	var networkModel *model.Network
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}

	//delete jobs
	var jobModel *model.Job
	query = d.DB.Where("1 = 1").Delete(&jobModel)
	if query.Error != nil {
		return false, query.Error
	}
	////delete producer
	var producerModel *model.Producer
	query = d.DB.Where("1 = 1").Delete(&producerModel)
	if query.Error != nil {
		return false, query.Error
	}
	////delete consumers
	var consumerModel *model.Consumer
	query = d.DB.Where("1 = 1").Delete(&consumerModel)
	if query.Error != nil {
		return false, query.Error
	}
	////delete consumersList
	var consumersList *model.Writer
	query = d.DB.Where("1 = 1").Delete(&consumersList)
	if query.Error != nil {
		return false, query.Error
	}
	//delete commands
	var commandsModel *model.CommandGroup
	query = d.DB.Where("1 = 1").Delete(&commandsModel)
	if query.Error != nil {
		return false, query.Error
	}

	//delete streams
	var streamsModel *model.Stream
	query = d.DB.Where("1 = 1").Delete(&streamsModel)
	if query.Error != nil {
		return false, query.Error
	}


	//delete networks
	var flowNetworkModel *model.FlowNetwork
	query = d.DB.Where("1 = 1").Delete(&flowNetworkModel)
	if query.Error != nil {
		return false, query.Error
	}

	return true, nil
}

//SyncTopics sync all the topics
func (d *GormDatabase) SyncTopics()  {

	d.B.RegisterTopicParent("aa", "aa")

	g, err := d.GetStreams(false)
	for _, obj := range g {
		d.B.RegisterTopicParent(model.CommonNaming.Stream, obj.UUID)
	}
	s, err := d.GetPlugins()
	for _, obj := range s {
		d.B.RegisterTopicParent(model.CommonNaming.Plugin, obj.UUID)
	}
	sub, err := d.GetProducers()
	for _, obj := range sub {
		d.B.RegisterTopicParent(model.CommonNaming.Producer, obj.UUID)
	}
	rip, err := d.GetConsumers()
	for _, obj := range rip {
		d.B.RegisterTopicParent(model.CommonNaming.Consumer, obj.UUID)
	}
	j, err := d.GetJobs()
	for _, obj := range j {
		d.B.RegisterTopicParent(model.CommonNaming.Job, obj.UUID)
	}
	n, err := d.GetNetworks(false, false)
	for _, obj := range n {
		d.B.RegisterTopicParent(model.CommonNaming.Network, obj.UUID)
	}
	de, err := d.GetDevices(false)
	for _, obj := range de {
		d.B.RegisterTopicParent(model.CommonNaming.Network, obj.UUID)
	}
	p, err := d.GetPoints(false)
	for _, obj := range p {
		d.B.RegisterTopicParent(model.CommonNaming.Point, obj.UUID)
	}
	node, err := d.GetNodesList()
	for _, obj := range node {
		eventbus.NodeContext.Set(obj.UUID, obj, cache.NoExpiration)
		d.B.RegisterTopicParent(model.CommonNaming.Node, obj.UUID)
	}



	if err != nil {

	}
}
