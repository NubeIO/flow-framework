package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func main()  {

	c := client.NewFlowRestClient("admin", "admin", "0.0.0.0", "1660")

	remotePointUUID := "id_p_TEST_REMOTE"
	remoteRubixUUID := "RUBIX_REMOTE"
	localRubixUUID := "id_n_5569693251d743c8"

	addNet, err := c.ClientAddNetwork("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add network")
	fmt.Println(addNet.Status)
	fmt.Println(addNet.Response.UUID)
	fmt.Println(addNet.Response.Name)

	addDev, err := c.ClientAddDevice(addNet.Response.UUID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add device")
	fmt.Println(addDev.Status)
	fmt.Println(addDev.Response.UUID)
	fmt.Println(addDev.Response.Name)

	addPoint, err := c.ClientAddPoint(addDev.Response.UUID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add point")
	fmt.Println(addPoint.Status)
	fmt.Println(addPoint.Response.UUID)
	fmt.Println(addPoint.Response.Name)

	stream := new(model.Stream)
	stream.Name = "test"
	stream.IsRemote = true
	addGateway, err := c.ClientAddGateway(stream)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println("Add gateway")
	fmt.Println(addGateway.Status)
	fmt.Println(addGateway.Response.UUID)
	fmt.Println(addGateway.Response.Name)

	pointUUID := addPoint.Response.UUID
	gatewayUUID := addGateway.Response.UUID


	// point 2 to make a producer connection to point 1
	tSub := new(client.Producer)
	tSub.Name = "test"
	tSub.Enable = true
	tSub.IsRemote  = true
	tSub.RemoteRubixUUID  = remoteRubixUUID
	tSub.FromUUID = remotePointUUID //remote point
	tSub.ToUUID = pointUUID  //local point
	tSub.StreamUUID = gatewayUUID
	tSub.ProducerApplication = "mapping"
	tSub.ProducerType = "point"

	addProducer, err := c.ClientAddProducer(*tSub)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println("Add Producer")
	fmt.Println(addProducer.Status)
	fmt.Println(addProducer.Response.UUID)
	fmt.Println(addProducer.Response.Name)

	remoteClient := client.NewFlowRestClient("admin", "admin", "0.0.0.0", "1661")

	// point 2 to make a producer connection to point 1
	rSub := new(client.Subscription)
	rSub.Name = "test"
	rSub.Enable = true
	rSub.IsRemote  = true
	rSub.RemoteRubixUUID  = localRubixUUID //local device id
	rSub.ToUUID = pointUUID  //local point
	rSub.StreamUUID = gatewayUUID
	rSub.ProducerApplication = "mapping"
	rSub.ProducerType = "point"

	addSubscription, err := remoteClient.ClientAddSubscription(*rSub)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Add Producer")
	fmt.Println(addSubscription.Status)
	fmt.Println(addSubscription.Response.UUID)
	fmt.Println(addSubscription.Response.Name)

	fmt.Println("FLOW-FRAMEWORK-TOKEN", c.ClientToken)


	//// point 1 to have 1 subscription to point 2
	//tSub2 := new(client.Subscription)
	//tSub2.Name = "test"
	//tSub2.Enable = true
	//tSub2.ThingUuid = addPoint.Response.UUID //pass in point 1 UUID
	//tSub2.GatewayUuid = addGateway.Response.UUID
	//tSub2.ProducerApplication = "mapping"
	//tSub2.ProducerType = "point"
	//
	//addSubscription, err := c.ClientAddSubscription(*tSub2)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//
	//}
	//fmt.Println("Add Subscription")
	//fmt.Println(addSubscription.Status)
	//fmt.Println(addSubscription.Response.UUID)
	//fmt.Println(addSubscription.Response.Name)








}

