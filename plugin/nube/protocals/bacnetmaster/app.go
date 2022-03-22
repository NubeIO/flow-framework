package main

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) bacnetUpdate(body mqtt.Message) {
	payload := new(bacnet_model.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	if err != nil {
		return
	}
	t, _ := mqttclient.TopicParts(body.Topic())
	const pointUUID = 14
	if t.Size() >= pointUUID {
		pUUID := t.Get(pointUUID)
		_pUUID := pUUID.(string)
		getPnt, err := i.db.GetOnePointByArgs(api.Args{AddressUUID: &_pUUID})
		if err != nil || getPnt.UUID == "" {
			log.Error("bacnet-master-plugin: ERROR on get GetPointByField() failed to find point", err, _pUUID)
			return
		}
		var pri model.Priority
		pri.P16 = payload.Value
		getPnt.Priority = &pri
		getPnt.CommonFault.InFault = false
		getPnt.CommonFault.MessageLevel = model.MessageLevel.Info
		getPnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
		getPnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
		getPnt.CommonFault.LastOk = time.Now().UTC()
		_, err = i.db.UpdatePointValue(getPnt.UUID, getPnt, true)
		if err != nil {
			log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
			return
		}
		return
	}
}
