package mqtt_helper

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type MqttConnection struct {
	mqttClient mqtt.Client
}


func NewMQTTConnection(u string, p string) (conn *MqttConnection, err error)  {
	//c := mqtt_config.MqttConfig("na", false)
	fmt.Println(u, p, 1111111)
	//var br mqtt_config.Broker
	//br.Host = "0.0.0.0"
	//br.Port = "1883"

	//c := mqtt_config.GetMqttConfig()
	opts := mqtt.NewClientOptions()
	host := "tcp://" + u + ":" + p
	fmt.Println(host)
	opts.AddBroker(fmt.Sprintf(host))
	opts.SetClientID("c.ClientId")
	opts.AutoReconnect = true
	opts.OnConnectionLost = connectLostHandler
	opts.OnConnect = connectHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("Connect problem: ", token.Error())
		return nil, token.Error()
	}
	conn = &MqttConnection{client}
	return conn, nil
}



func (conn *MqttConnection) MQTTIsConnected() bool {
	connected := conn.mqttClient.IsConnected()
	if !connected {
		log.Println("Health check MQTT fails")
	}
	return connected
}

func (conn *MqttConnection) PublishMessage(message string, topic string) {
	token := conn.mqttClient.Publish(topic, 1, false, message)
	token.Wait()
	log.Println("Publish to topic: ", topic)
}


var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Println("Connection lost: ", err)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Mqtt connected")
}


