package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqtt_client"
	"github.com/mustafaturan/bus/v3"
	"github.com/mustafaturan/monoton/v2"
	"github.com/mustafaturan/monoton/v2/sequencer"
	"log"
)


func NewBus() *bus.Bus {
	// configure id generator
	node        := uint64(1)
	initialTime := uint64(1577865600000)
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime);if err != nil {
		panic(err)
	}
	// init an id generator
	var idGenerator bus.Next = m.Next
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}
	b.RegisterTopics("points")
	b.RegisterHandler("points", PointHandler)
	return b
}

type BusPayload struct {
	GatewayUUID  	string   		`json:"gateway_uuid"`
	ThingName   	string
	MessageString  	string   		`json:"message_string"`
	MessageTS  		string   		`json:"message_ts"`
	Action  		string   		`json:"action"`
}

var BUS = NewBus()
var BusBackground = context.Background()

func publishMQTT(sensorStruct *model.Point) {
	a := mqtt_client.NewClient(mqtt_client.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(sensorStruct.Name, 888888)
	fmt.Println(a.IsConnected())
	topic := fmt.Sprintf("rubix/%s", sensorStruct.UUID)
	data, err := json.Marshal(sensorStruct)
	if err != nil {
		log.Println(err)
	}

	err = a.Publish(topic, mqtt_client.AtMostOnce, false, string(data))
	if err != nil {
		fmt.Println(err)
	}
}

var PointHandler = bus.Handler {
	Handle: func(ctx context.Context, e bus.Event) {
		//NewAgent
		data, _ := e.Data.(*model.Point)
		fmt.Println(data, 99999)
		publishMQTT(data)
		fmt.Println(e.Topic)
		fmt.Println(e.Data)
	},
	Matcher: ".*", // matches all topics
}

