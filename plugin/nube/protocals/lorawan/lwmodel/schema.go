package lwmodel

import "github.com/NubeIO/flow-framework/plugin/defaults"

func GetNetworkSchema() *Network {
	m := &Network{}
	defaults.Set(m)
	return m
}

func GetDeviceSchema() *Device {
	device := &Device{}
	defaults.Set(device)
	return device
}

func GetPointSchema() *Point {
	point := &Point{}
	defaults.Set(point)
	return point
}

type NameNet struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"lorawan"`
}

type NameDev struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"device"`
}

type NamePnt struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"point"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
}

type Interface struct {
	Type     string   `json:"type" default:"string"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[]"`
	Default  string   `json:"default" default:""`
}

type Network struct {
	Name        NameNet           `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lorawan"`
	} `json:"plugin_name"`
	Interface struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[]"`
		Default  string   `json:"default" default:""`
	} `json:"network_interface"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"self-mapping\",\"rubix-io-to-lorawan\"]"`
		Default  string   `json:"default" default:""`
	} `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:"local"`
	} `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:""`
	} `json:"auto_mapping_flow_network_uuid"`
	Enable EnableStruct `json:"enable"`
}

type FlowDevice struct {
	Name           NameDev           `json:"name"`
	Description    DescriptionStruct `json:"description"`
	DeviceObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"device_object_id"`
}

type Point struct {
	Name        NamePnt           `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\"]"`
		Default  string   `json:"default" default:"analog_value"`
	} `json:"object_type"`
	ObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"0"`
	} `json:"object_id"`
	WriteMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"read_only\",\"write_only\"]"`
		Default  string   `json:"default" default:"read_only"`
	} `json:"write_mode"`
	WritePriority struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16]"`
		Default  int    `json:"default" default:"16"`
	} `json:"write_priority"`
}
