package core

import (
	"fmt"
	"log"

	"github.com/NubeDev/flow-framework/floweng/mqtt_helper"
)

func MQTTClient() SourceSpec {
	return SourceSpec{
		Name: "MQTTClient",
		Type: MQTTCLIENT,
		New:  NewMQTT,
	}
}
func (mc OSC) GetType() SourceType {
	return MQTTCLIENT
}

type OSCMsg struct {
	address  string
	argument interface{} // TODO think harder
	err      chan error
}
type OSC struct {
	quit   chan chan error
	conf   chan OSCConf
	toSend chan OSCMsg
}
type OSCConf struct {
	url  string
	port string
	err  chan error
}

func NewMQTT() Source {
	OSC := &OSC{
		quit:   make(chan chan error),
		toSend: make(chan OSCMsg),
		conf:   make(chan OSCConf),
	}
	return OSC
}

type Broker struct {
	Host     string
	Port     string
	Topic    string
	User     string
	Password string
	ClientId string
}

func (mc *OSC) Stop() {
	m := make(chan error)
	mc.quit <- m
	// block until closed
	err := <-m
	if err != nil {
		log.Fatal(err)
	}
}

var client *mqtt_helper.MqttConnection
var connected = false

func MQTTClientConnect() Spec {
	return Spec{
		Name:    "MQTTClientConnect",
		Inputs:  []Pin{Pin{"url", STRING}, Pin{"port", STRING}},
		Outputs: []Pin{Pin{"connected", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			//OSC := s.(*OSC)
			url, ok := in[0].(string)
			if !ok {
				out[0] = NewError("OSCClientConnect requires string url")
				return nil
			}
			port, ok := in[1].(string)
			if !ok {
				out[0] = NewError("OSCClientConnect requires number port")
				return nil
			}
			fmt.Println("MQTTClientConnect")
			fmt.Println(connected)
			if connected == false {
				client, _ = mqtt_helper.NewMQTTConnection(url, port)
				out[0] = true
				connected = client.MQTTIsConnected()
			} else {
				out[0] = false
			}
			return nil
		},
	}
}

var lastMessage = ""

func MQTTClientSend() Spec {
	return Spec{
		Name:    "MQTTClientSend",
		Inputs:  []Pin{Pin{"address", STRING}},
		Outputs: []Pin{Pin{"sent", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			addr, ok := in[0].(string)
			if !ok {
				out[0] = NewError("OSCSend requires string address")
				return nil
			}
			if connected == true {
				out[0] = true
				if lastMessage != addr {
					client.PublishMessage(addr, "test44")
				}
				out[0] = false
				lastMessage = addr
			} else {
				out[0] = false
			}
			return nil
		},
	}
}
