package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func publishMQTT(sensorStruct model.ProducerBody) {
	a, _ := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Fatal(err)
	}
	topic := fmt.Sprintf("rubix/%s", sensorStruct.ProducerUUID)
	data, err := json.Marshal(sensorStruct)
	if err != nil {
		log.Fatal(err)
	}
	err = a.Publish(topic, mqttclient.AtMostOnce, false, string(data))
	if err != nil {
		log.Fatal(err)
	}
}

//used for getting data into the plugins
var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println(msg.Topic(), " ", "NEW MQTT MES")
	GetService().RegisterTopic(MQTTUpdated)
	err := GetService().Emit(CTX(), MQTTUpdated, msg)
	if err != nil {
		return
	}
}

func RegisterMQTTBus() {
	c, _ := mqttclient.GetMQTT()
	//TODO this needs to be removed as its for a plugin, the plugin needs to register the topics it wants the main framework to subscribe to, also unsubscribe when the plugin is disabled
	err := c.Subscribe("+/+/+/+/+/+/rubix/bacnet_server/points/+/#", mqttclient.AtMostOnce, handle) //lorawan chirpstack
	//err = c.Subscribe("application/#", mqttclient.AtMostOnce, handle) //lorawan chirpstack
	err = c.Subscribe("application/+/device/+/rx", mqttclient.AtMostOnce, handle) //lorawan chirpstack
	if err != nil {

	}
}
