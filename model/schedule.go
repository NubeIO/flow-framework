package model

import (
	"gorm.io/datatypes"
	"time"
)

//Schedule model
type Schedule struct {
	CommonUUID
	CommonNameUnique
	CommonEnable
	CommonThingClass
	CommonThingType
	FlowNetworkUUID string   `json:"flow_network_uuid,omitempty" gorm:"TYPE:varchar(255) REFERENCES flow_networks;not null;default:null"`
	IsGlobal  *bool          `json:"is_global"`
	Schedules datatypes.JSON `json:"schedule"`
	CommonCreated
}

type Schedules struct {
	Events  map[string]Events  `json:"events"`
	Weekly  map[string]Weekly  `json:"weekly"`
	Holiday map[string]Holiday `json:"holiday"`
}

type Events struct {
	Name  string `json:"name"`
	Dates []struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"dates"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

type Weekly struct {
	Name  string   `json:"name"`
	Days  []string `json:"days"`
	Start string   `json:"start"`
	End   string   `json:"end"`
	Value int      `json:"value"`
	Color string   `json:"color"`
}

type Holiday struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Date  string `json:"date"`
	Value int    `json:"value"`
}
