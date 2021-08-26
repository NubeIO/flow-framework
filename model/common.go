package model

import (
	"time"
)


var CommonNaming = struct {
	Network   		string
	Device   		string
	Point   		string
	Stream   		string
	Job   			string
	Subscriber   	string
	Subscription   	string
	Alert   		string
	Mapping   		string
	CommandGroup   	string
	Rubix   		string
	RubixGlobal   	string
	FlowNetwork   	string

}{
	Network:   			"network",
	Device:   			"device",
	Point:   			"point",
	Stream:   			"stream",
	Job:   				"job",
	Subscriber:   		"subscriber",
	Subscription:   	"subscription",
	Alert:   			"alert",
	Mapping:   			"mapping",
    CommandGroup:   	"command_group",
	Rubix:   			"rubix",
	RubixGlobal:   		"rubix_global",
	FlowNetwork:   		"flow_network",
}


var CommonNamingCommandGroup = struct {
	PointWrite  			string
	MasterSchedule   		string
	SilenceAlarm   		    string


}{
	PointWrite:   		"point_write",
	MasterSchedule:   	"master_schedule",
	SilenceAlarm:   	"silence_alarm",

}




type CommonDescription struct {
	Description string `json:"description"`
}

type CommonName struct {
	Name string `json:"name"  valid:"required~name is required"`
}

type CommonNameUnique struct {
	Name  string `json:"name"  gorm:"type:varchar(255);unique;not null"`
}

type CommonEnable struct {
	Enable 	*bool `json:"enable"`
}

type CommonID struct {
	ID	string `json:"id"`
}

type CommonIDUnique struct {
	Name  string `json:"id"  gorm:"type:varchar(255);unique;not null"`
}

type CommonUUID struct {
	UUID	string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
}

type CommonFlowNetworkUUID struct {
	FlowNetworkUUID	string 	`json:"flow_network_uuid" gorm:"type:varchar(255);unique;not null"`
}

type CommonFlowNetworkName struct {
	FlowNetworkName	string 	`json:"flow_network_name" gorm:"type:varchar(255);unique;not null"`
}

type CommonRubixUUID struct {
	RubixUUID	string 	`json:"rubix_uuid"`
}

type CommonCreated struct {
	CreatedAt 	time.Time `json:"created_on"`
	UpdatedAt 	time.Time  `json:"updated_on"`
}

type CommonHistory struct {
	EnableHistory 	bool   `json:"history_enable"`
}

type CommonValue struct {
	Value		float64 `json:"value"`
	ValueRaw	string `json:"value_raw"`
}

type CommonFault struct {
	Fault 			bool `json:"fault"`
	FaultMessage 	bool `json:"fault_message"`
}


type CommonIP struct {
	IP		string `json:"ip"`
	Port 	int `json:"port"`
	HTTP 	bool `json:"http"`
	HTTPS 	bool `json:"https"`
}


type CommonStore struct {
	CommonValue
	CommonFault
}

//CommonProducer a point or job
type CommonProducer struct {
	CommonUUID
	CommonNameUnique
	CommonDescription
	CommonEnable
	CommonValue
	CommonHistory
	CommonCreated
}


type CommonSubscriberPermissions struct {
	Blacklist 		bool  	`json:"blacklist"`
	ReadOnly  		bool 	`json:"read_only"`
	AllowCRUD  		bool 	`json:"allow_crud"` //not sure if this will be used, but it will allow the subscriber to update the producer
}


type CommonSubscriber struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`

}

type CommonSubscription struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`


}
