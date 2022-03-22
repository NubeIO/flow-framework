package bacnet_model

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type ServerPing struct {
	Version        string `json:"version"`
	UpTimeDate     string `json:"up_time_date"`
	UpMin          string `json:"up_min"`
	UpHour         string `json:"up_hour"`
	DeploymentMode string `json:"deployment_mode"`
	Bacnet         struct {
		Enabled                    bool    `json:"enabled"`
		Status                     bool    `json:"status"`
		UsePreSetEthernetInterface bool    `json:"use_pre_set_ethernet_interface"`
		PreSetEthernetInterface    string  `json:"pre_set_ethernet_interface"`
		DefaultPointCov            float32 `json:"default_point_cov"`
	} `json:"bacnet"`
	Mqtt struct {
		Enabled bool `json:"enabled"`
		Status  bool `json:"status"`
	} `json:"mqtt"`
}

type Server struct {
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	DeviceId          string `json:"device_id"`
	LocalObjName      string `json:"local_obj_name"`
	ModelName         string `json:"model_name"`
	VendorId          string `json:"vendor_id"`
	VendorName        string `json:"vendor_name"`
	EnableIpByNicName bool   `json:"enable_ip_by_nic_name"`
	IpByNicName       string `json:"ip_by_nic_name"` //eth0
}

// MqttPayload payload from the bacnet server
type MqttPayload struct {
	Value    *float64
	Priority int
}

type BacnetPoint struct {
	AddressUUID          string          `json:"address_uuid,omitempty"`
	ObjectType           string          `json:"object_type,omitempty"`
	ObjectName           string          `json:"object_name,omitempty"`
	Address              int             `json:"address,omitempty"`
	RelinquishDefault    float64         `json:"relinquish_default"`
	EventState           string          `json:"event_state,omitempty"`
	Units                string          `json:"units,omitempty"`
	Description          string          `json:"description,omitempty"`
	Enable               bool            `json:"enable,omitempty"`
	Fault                bool            `json:"fault,omitempty"`
	DataRound            float64         `json:"data_round,omitempty"`
	DataOffset           float64         `json:"data_offset,omitempty"`
	UseNextAvailableAddr bool            `json:"use_next_available_address,omitempty"`
	COV                  float64         `json:"cov,omitempty"`
	Priority             *model.Priority `json:"priority_array_write,omitempty"`
}

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

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"bacnetserver"`
	} `json:"plugin_name"`
}

type Device struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"analogValue\",\"analogOutput\"]"`
		Default  string   `json:"default" default:"analogValue"`
	} `json:"object_type"`
	AddressID struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_id"`
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
