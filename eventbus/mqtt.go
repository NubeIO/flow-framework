package eventbus

import (
	"encoding/json"
	"fmt"

	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func publishMQTT(sensorStruct model.ProducerBody) {
	a, _ := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Error(err)
	}
	topic := fmt.Sprintf("rubix/%s", sensorStruct.ProducerUUID)
	data, err := json.Marshal(sensorStruct)
	if err != nil {
		log.Error(err)
	}
	err = a.Publish(topic, mqttclient.AtMostOnce, false, string(data))
	if err != nil {
		log.Error(err)
	}
}

var debug = false

// used for getting data into the plugins
var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if debug {
		log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	}
	GetService().RegisterTopic(MQTTUpdated)
	err := GetService().Emit(CTX(), MQTTUpdated, msg)
	if err != nil {
		return
	}
}

func RegisterMQTTBus(enableDebug bool) {
	if enableDebug {
		debug = true
	}
	c, _ := mqttclient.GetMQTT()
	// TODO this needs to be removed as its for a plugin, the plugin needs to register the topics it wants the main framework to subscribe to, also unsubscribe when the plugin is disabled
	c.Subscribe("+/+/+/+/+/+/rubix/bacnet_server/points/+/#", mqttclient.AtMostOnce, handle) // bacnet-server
	c.Subscribe("+/+/+/+/+/+/rubix/bacnet_master/points/+/#", mqttclient.AtMostOnce, handle) // bacnet-bserver
	c.Subscribe("application/+/device/+/event/up", mqttclient.AtMostOnce, handle)            // lorawan
}
