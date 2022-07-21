package master

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeIO/flow-framework/plugin/defaults"
	"github.com/NubeIO/lib-networking/networking"
	"github.com/gin-gonic/gin"
)

type NameNet struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"bacnet-master"`
}

type NameDev struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"device"`
}

type NamePnt struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
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

type TransportTypeStruct struct {
	Type     string   `json:"type" default:"array"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[\"rs485\",\"ip\"]"`
	Default  string   `json:"default" default:"rs485"`
}

type Network struct {
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"true"`
	} `json:"enable"`
	Name struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"bac_net"`
		DisplayName string `json:"display_name" default:"Network Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"bacnetmaster"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	Port          struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Default  int    `json:"default" default:"47808"`
	} `json:"port"`
	Interface struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[]"`
		Default  string   `json:"default" default:""`
	} `json:"network_interface"`
	MaxPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Max Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"0.1"`
	} `json:"max_poll_rate"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"self-mapping\",\"rubix-io-to-bacnetserver\"]"`
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
	AutoMappingEnableHistories struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"false"`
	} `json:"auto_mapping_enable_histories"`
}

type Device struct {
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"true"`
	} `json:"enable"`
	Name struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"bac_dev"`
		DisplayName string `json:"display_name" default:"Device Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	Host        struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Options  string `json:"options" default:"192.168.15.10"`
		Default  string `json:"default" default:"192.168.15.10"`
	} `json:"host"`
	Port struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  int    `json:"options" default:"47808"`
		Default  int    `json:"default" default:"47808"`
	} `json:"port"`
	DeviceObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"device_object_id"`
	NetworkNumber struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"0"`
	} `json:"network_number"`
	DeviceMac struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"255"`
		Default  int    `json:"default" default:"0"`
	} `json:"device_mac"`
	MaxADPU struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"true"`
		Options  []int  `json:"options" default:"[50, 128, 206, 480, 1024, 1476]"`
		Default  int    `json:"default" default:"1024"`
	} `json:"max_adpu"`
	Segmentation struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"segmentation_both\",\"no_segmentation\",\"segmentation_transmit\",\"segmentation_receive\"]"`
		Default  string   `json:"default" default:"no_segmentation"`
	} `json:"segmentation"`
	FastPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Fast Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"5"`
	} `json:"fast_poll_rate"`
	NormalPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Normal Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"30"`
	} `json:"normal_poll_rate"`
	SlowPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Slow Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"120"`
	} `json:"slow_poll_rate"`
}

type Point struct {
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"true"`
		Default  *bool  `json:"default" default:"true"`
	} `json:"enable"`
	Name struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"bac_pnt"`
		DisplayName string `json:"display_name" default:"Point Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\",\"multi_state_input\",\"multi_state_value\",\"multi_state_output\"]"`
		Default  string   `json:"default" default:"analog_output"`
	} `json:"object_type"`
	ObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"object_id"`
	WriteMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"read_only\",\"write_only\",\"write_once_then_read\"]"`
		Default  string   `json:"default" default:"read_only"`
	} `json:"write_mode"`
	WritePriority struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16]"`
		Default  int    `json:"default" default:"16"`
	} `json:"write_priority"`
	PollPriority struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"high\",\"normal\",\"low\"]"`
		Default  string   `json:"default" default:"normal"`
	} `json:"poll_priority"`
	PollRate struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"fast\",\"normal\",\"slow\"]"`
		Default  string   `json:"default" default:"normal"`
	} `json:"poll_rate"`
	Fallback struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"false"`
		Default     *float64 `json:"default" default:""`
		DisplayName string   `json:"display_name" default:"Fallback Value"`
		Nullable    bool     `json:"nullable" default:"true"`
	} `json:"fallback"`
}

var nets = networking.New()

func GetNetworkSchema() *Network {
	m := &Network{}
	defaults.Set(m)
	names, err := nets.GetInterfacesNames()
	if err != nil {
		return m
	}
	var out []string
	out = append(out, "eth0")
	out = append(out, "eth1")
	for _, name := range names.Names {
		if name != "lo" {
			out = append(out, name)
		}
	}
	m.Interface.Options = out
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

func BodyWhoIs(ctx *gin.Context) (dto *bacnet.WhoIsOpts, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}
