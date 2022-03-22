package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
)

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(bacnet_model.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	t, _ := mqttclient.TopicParts(body.Topic())
	const objectType = 10
	top := t.Get(objectType)
	tt := top.(string)
	objType, addr := getPointAddr(tt)
	var point model.Point
	var pri model.Priority
	pri.P16 = payload.Value
	point.Priority = &pri
	pnt, _ := i.db.GetOnePointByArgs(api.Args{ObjectType: &objType, AddressID: &addr}) //TODO check conversion if existing exists, as in the same addr
	if err != nil {
		log.Error("BACNET UPDATE POINT PointAndQuery")
		return nil, err
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	if pnt == nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	_, err = i.db.UpdatePointValue(pnt.UUID, &point, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*bacnet_model.Point, error) {
	var point bacnet_model.BacnetPoint
	point.ObjectName = body.Name
	point.Enable = true
	point.Description = body.Description
	point.Address = utils.IntIsNil(body.AddressID)
	point.ObjectType = body.ObjectType
	point.COV = utils.Float64IsNil(body.COV)
	point.EventState = "normal"
	point.Units = "noUnits"
	point.RelinquishDefault = utils.Float64IsNil(body.Fallback)

	rt.Method = nrest.POST
	rt.Path = "/api/bacnet/points"

	addPoint, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: point})

	if &point == nil {
		log.Error("BACNET ADD POINT issue on add")
		return nil, errors.New("BACNET ADD POINT issue on add")
	}
	if err != nil {
		log.Errorf("BACNET: ADD POINT issue on add rest: %v\n", err)
		return nil, err
	}
	bacPntUUID := gjson.Get(string(addPoint.Body), "uuid").String()
	body.AddressUUID = &bacPntUUID
	_, err = i.db.UpdatePoint(body.UUID, body, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on update point when getting bacnet point uuid")
		return nil, err
	}
	return nil, nil

}

//pointPatch from rest
func (i *Instance) pointPatch(body *model.Point) (*model.Point, error) {
	if body.AddressUUID == nil {
		return nil, errors.New("no address_uuid")
	}
	point := new(bacnet_model.BacnetPoint)
	point.Priority = new(model.Priority)
	if (*body.Priority).P16 != nil {
		(*point.Priority).P16 = (*body.Priority).P16
	}
	point.ObjectName = body.Name
	point.Address = utils.IntIsNil(body.AddressID)
	point.ObjectType = body.ObjectType
	//point.Units = body.Unit
	point.Description = body.Description

	rt.Method = nrest.PATCH
	rt.Path = fmt.Sprintf("/api/bacnet/points/uuid/%s", *body.AddressUUID)
	_, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: point})

	if err != nil {
		log.Errorf("BACNET: EDIT POINT issue on add rest: %v\n", err)
		return nil, err
	}
	return nil, nil

}

//deletePoint point make sure
func (i *Instance) deletePoint(body *model.Point) (bool, error) {
	if body.AddressUUID == nil {
		return false, errors.New("no address_uuid")
	}
	rt.Method = nrest.DELETE
	rt.Path = fmt.Sprintf("/api/bacnet/points/uuid/%s", *body.AddressUUID)
	_, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
	if err != nil {
		log.Errorf("BACNET: DELETE POINT issue on add rest: %v\n", err)
		return false, err
	}
	return true, nil
}

//bacnetServerDeletePoint point make sure
func (i *Instance) bacnetServerDeletePoint(body *bacnet_model.BacnetPoint) (bool, error) {
	//cli := plgrest.NewNoAuth(ip, port)
	//_, err := cli.DeletePoint(body.ObjectType, body.Address)
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}
