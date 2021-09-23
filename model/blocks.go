package model

import (
	"encoding/json"
	"log"
)

// TODO: cleanup original Proto.. structs

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Block struct {
	ID       int      `json:"id"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Parent   int      `json:"parent"`
	Position Position `json:"position" gorm:"embedded"`
	IsSource bool     `json:"is_source"`
	// TODO: maybe add list of static block routes
}

type BlockStaticRoute struct {
	ID         int `json:"id"`
	BlockID    int `json:"block_id" gorm:"REFERENCES blocks;not null`
	RouteIndex int `json:"route_index"`
	Type       int `json:"type"`
}

type BlockRouteValueNumber struct {
	BlockRoute int     `json:"block_route" gorm:"REFERENCES block_static_route;not null`
	Value      float64 `json:"value"`
}

type BlockRouteValueString struct {
	BlockRoute int    `json:"block_route" gorm:"REFERENCES block_static_route;not null`
	Value      string `json:"value"`
}

type BlockRouteValueBool struct {
	BlockRoute int  `json:"block_route" gorm:"REFERENCES block_static_route;not null`
	Value      bool `json:"value"`
}

type ProtoBlockStaticRoute struct {
	BlockStaticRoute
	Value interface{} `json:"value"`
}

type ProtoBlockRouteNumber struct {
	BlockStaticRoute
	BlockRouteValueNumber
}

type ProtoBlockRouteString struct {
	BlockStaticRoute
	BlockRouteValueString
}

type ProtoBlockRouteBool struct {
	BlockStaticRoute
	BlockRouteValueBool
}

type Connection struct {
	ID          int `json:"id"`
	SourceID    int `json:"source_id" gorm:"REFERENCES blocks;not null`
	SourceRoute int `json:"source_route"`
	TargetID    int `json:"target_id" gorm:"REFERENCES blocks;not null`
	TargetRoute int `json:"target_route"`
}

type SourceParameter struct {
	BlockID    int    `json:"block_id" gorm:"REFERENCES blocks;not null`
	Parameters string `json:"params"`
}

type Link struct {
	ID       int `json:"id"`
	SourceID int `json:"source_id"`
	BlockID  int `json:"block_id"`
}

// type Group struct {
// 	ID       int     `json:"id"`
// 	Group    int     `json:"parent"`
// 	Children []int   `json:"children"`
// 	Label    string  `json:"label"`
// 	PosX     float64 `json:"x"`
// 	PosY     float64 `json:"y"`
// }

func (source *SourceParameter) MarshalParameters(params []map[string]string) {
	bytes, err := json.Marshal(params)
	if err != nil {
		log.Println("Error marshaling source paramaters")
		source.Parameters = ""
	}
	source.Parameters = string(bytes)
}

func (source *SourceParameter) UnmarshalParameters() []map[string]string {
	var params []map[string]string
	err := json.Unmarshal([]byte(source.Parameters), &params)
	if err != nil {
		log.Println("Error unmarshaling source paramaters")
	}
	return params
}
