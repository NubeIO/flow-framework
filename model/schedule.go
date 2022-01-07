package model

import (
	"gorm.io/datatypes"
	"time"
)

type Schedule struct {
	CommonUUID
	CommonNameUnique
	CommonEnable
	CommonThingClass
	CommonThingType
	IsActive  *bool          `json:"is_active"`
	IsGlobal  *bool          `json:"is_global"`
	Schedules datatypes.JSON `json:"schedules"`
	CommonCreated
}

type Schedules struct {
	Events    map[string]Events    `json:"events"`
	Weekly    map[string]Weekly    `json:"weekly"`
	Exception map[string]Exception `json:"exception"`
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

type Exception struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Date  string `json:"date"`
	Value int    `json:"value"`
}
