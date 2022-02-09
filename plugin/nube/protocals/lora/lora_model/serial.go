package lora_model

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type TransportTypeStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Default  string `json:"default" default:"serial"`
}

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"false"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lora"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	SerialPort    struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
	} `json:"serial_port"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Default  int    `json:"default" default:"9600"`
	} `json:"serial_baud_rate"`
}

type Device struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`

	AddressUUID struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"min" default:"8"`
		Max         int    `json:"max" default:"8"`
		DisplayName string `json:"display_name" default:"Address UUID"`
	} `json:"address_uuid"`
	Model struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"min" default:"4"`
		Max         int    `json:"max" default:"6"`
		DisplayName string `json:"display_name" default:"THLM"`
	} `json:"bacnet_model"`
	Enable EnableStruct `json:"enable"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	AddressUUID struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"min" default:"8"`
		Max         int    `json:"max" default:"8"`
		DisplayName string `json:"display_name" default:"Address UUID"`
	} `json:"address_uuid"`
	IoId struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"DI1\",\"UI1\",\"AO1\",\"temp\",\"pulse\",\"motion\",\"light\"]"`
	} `json:"io_id"`
	IoType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"0-10dc\",\"0-40ma\",\"thermistor\"]"`
	} `json:"io_type"`
	ThingClass struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_class"`
	ThingType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_type"`
}

func GetNetworkSchema() *Network {
	network := &Network{}
	defaults.Set(network)
	return network
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
