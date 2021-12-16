package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	baseModel "github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/plgrest"
	"github.com/NubeIO/flow-framework/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(model.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	fmt.Println(body.Payload())
	fmt.Println(body.Topic())
	t, _ := mqttclient.TopicParts(body.Topic())
	const objectType = 10
	top := t.Get(objectType)
	tt := top.(string)
	objType, addr := getPointAddr(tt)
	var point baseModel.Point
	var pri baseModel.Priority
	pri.P16 = payload.Value
	point.Priority = &pri
	pnt, _ := i.db.PointAndQuery(objType, addr) //TODO check conversion if existing exists, as in the same addr
	if err != nil {
		log.Error("BACNET UPDATE POINT PointAndQuery")
		return nil, err
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = baseModel.MessageLevel.Info
	point.CommonFault.MessageCode = baseModel.CommonFaultCode.Ok
	point.CommonFault.Message = baseModel.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	if pnt == nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	_, _ = i.db.UpdatePointValue(pnt.UUID, &point, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *baseModel.Point) (*model.Point, error) {
	var point model.BacnetPoint
	point.ObjectName = body.Name
	point.Enable = true
	point.Description = body.Description
	point.Address = utils.IntIsNil(body.AddressID)
	point.ObjectType = body.ObjectType
	point.COV = utils.Float64IsNil(body.COV)
	point.EventState = "normal"
	point.Units = "noUnits"
	point.RelinquishDefault = utils.Float64IsNil(body.Fallback)
	cli := plgrest.NewNoAuth(ip, port)
	if &point == nil {
		log.Error("BACNET ADD POINT issue on add")
		return nil, errors.New("BACNET ADD POINT issue on add")
	}
	_, err := cli.AddPoint(point)
	//TODO check conversion if existing exists, as in the same addr and also set the point in fault or out of fault
	if err != nil {
		log.Errorf("BACNET: ADD POINT issue on add rest: %v\n", err)
		return nil, err
	}
	return nil, nil

}

//pointPatch from rest
func (i *Instance) pointPatch(body *baseModel.Point) (*model.Point, error) {
	//point := new(model.BacnetPoint)
	point := new(model.BacnetPoint)
	//point.Priority.P1 = body.Priority.P1

	//point.Priority.P2 = body.Priority.P2
	//point.Priority.P3 = body.Priority.P3
	//point.Priority.P4 = body.Priority.P4
	//point.Priority.P5 = body.Priority.P5
	//point.Priority.P6 = body.Priority.P6
	//point.Priority.P7 = body.Priority.P7
	//point.Priority.P8 = body.Priority.P8
	//point.Priority.P9 = body.Priority.P9
	//point.Priority.P10 = body.Priority.P10
	//point.Priority.P11 = body.Priority.P11
	//point.Priority.P12 = body.Priority.P12
	//point.Priority.P13 = body.Priority.P13
	//point.Priority.P14 = body.Priority.P14
	//point.Priority.P15 = body.Priority.P15

	point.Priority = new(baseModel.Priority)
	if (*body.Priority).P16 != nil {
		(*point.Priority).P16 = (*body.Priority).P16
	}
	point.ObjectName = body.Name
	addr := body.AddressID
	obj := body.ObjectType

	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.EditPoint(*point, obj, utils.IntIsNil(addr))
	if err != nil {
		log.Errorf("BACNET: EDIT POINT issue on add rest: %v\n", err)
		return nil, err
	}
	return nil, nil

}

//deletePoint point make sure
func (i *Instance) deletePoint(body *baseModel.Point) (bool, error) {
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.DeletePoint(body.ObjectType, utils.IntIsNil(body.AddressID))
	if err != nil {
		return false, err
	}
	return true, nil
}

//bacnetServerDeletePoint point make sure
func (i *Instance) bacnetServerDeletePoint(body *model.BacnetPoint) (bool, error) {
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.DeletePoint(body.ObjectType, body.Address)
	if err != nil {
		return false, err
	}
	return true, nil
}
