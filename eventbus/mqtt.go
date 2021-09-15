package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqttclient"
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
	fmt.Println(a.IsConnected())
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
	GetService().RegisterTopic(MQTTUpdated)
	err := GetService().Emit(CTX(), MQTTUpdated, msg)
	if err != nil {
		return
	}
	fmt.Println(msg.Topic(), 999999)
	fmt.Printf("MSG recieved pointsValue: %s\n", msg.Payload())
}

func RegisterMQTTBus() ()  {
	c, _ := mqttclient.GetMQTT()
	//TODO this needs to be removed as its for a plugin, the plugin needs to register the topics it wants the main framework to subscribe to
	err := c.Subscribe("+/+/+/+/+/+/rubix/bacnet_server/points/+/#",mqttclient.AtMostOnce, handle)
	if err != nil {

	}
}
