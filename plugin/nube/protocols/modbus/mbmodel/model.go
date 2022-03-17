package mbmodel

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
	"github.com/NubeIO/null"
)

type Priority struct {
	P1  null.Float `json:"_1,omitempty"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2  null.Float `json:"_2,omitempty"`
	P3  null.Float `json:"_3,omitempty"`
	P4  null.Float `json:"_4,omitempty"`
	P5  null.Float `json:"_5,omitempty"`
	P6  null.Float `json:"_6,omitempty"`
	P7  null.Float `json:"_7,omitempty"`
	P8  null.Float `json:"_8,omitempty"`
	P9  null.Float `json:"_9,omitempty"`
	P10 null.Float `json:"_10,omitempty"`
	P11 null.Float `json:"_11,omitempty"`
	P12 null.Float `json:"_12,omitempty"`
	P13 null.Float `json:"_13,omitempty"`
	P14 null.Float `json:"_14,omitempty"`
	P15 null.Float `json:"_15,omitempty"`
	P16 null.Float `json:"_16"` //removed and added to the point to save one DB write
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"mb"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type TransportTypeStruct struct {
	Type     string   `json:"type" default:"array"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[\"serial\",\"ip\",\"LoRa\"]"`
	Default  string   `json:"default" default:"serial"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"modbus"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	SerialPort    struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
		Default  string `json:"default" default:""`
	} `json:"serial_port"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[9600, 38400, 57600, 115200]"`
		Default  int    `json:"default" default:""`
	} `json:"serial_baud_rate"`
	SerialParity struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"odd\",\"even\",\"none\"]"`
		Default  string   `json:"default" default:""`
	} `json:"serial_parity"`
	SerialDataBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[7, 8]"`
		Default  int    `json:"default" default:""`
	} `json:"serial_data_bits"`
	SerialStopBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2]"`
		Default  int    `json:"default" default:""`
	} `json:"serial_stop_bits"`
	MaxPollRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"max_poll_rate"`
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
	} `json:"enable"`
}

type Device struct {
	Name          NameStruct          `json:"name"`
	Description   DescriptionStruct   `json:"description"`
	TransportType TransportTypeStruct `json:"transport_type"`
	AddressId     struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_id"`
	Host struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Options  string `json:"options" default:"192.168.15.10"`
		Default  string `json:"default" default:""`
	} `json:"host"`
	Port struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  int    `json:"options" default:"502"`
		Default  int    `json:"default" default:""`
	} `json:"port"`
	FastPollRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"5"`
	} `json:"fast_poll_rate"`
	NormalPollRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"30"`
	} `json:"normal_poll_rate"`
	SlowPollRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"120"`
	} `json:"slow_poll_rate"`
	ZeroMode struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
	} `json:"zero_mode"`
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
	} `json:"enable"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		//Options  []string `json:"options" default:"[\"read_coil\",\"read_coils\",\"read_discrete_input\",\"read_discrete_inputs\",\"write_coil\",\"read_registers\",\"read_holding\",\"read_holdings\",\"write_holding\",\"write_holdings\",\"read_int_16\",\"write_int_16\",\"read_uint_16\",\"write_uint_16\",\"read_int_32\",\"write_int_32\",\"read_uint_32\",\"write_uint_32\",\"read_float_32\",\"write_float_32\",\"read_float_64\",\"write_float_64\"]"`
		Options []string `json:"options" default:"[\"read_coil\",\"write_coil\",\"read_discrete_input\",\"read_register\",\"read_holding\",\"write_holding\"]"`
		Default string   `json:"default" default:"read_coils"`
	} `json:"object_type"`
	AddressId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_id"`
	AddressLength struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_length"`
	DataType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"digital\",\"uint16\",\"int16\",\"uint32\",\"int32\",\"uint64\",\"int64\",\"float32\",\"float64\"]"`
		Default  string   `json:"default" default:"digital"`
	} `json:"data_type"`
	WriteMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"read_once\",\"read_only\",\"write_once\",\"write_once_read_once\",\"write_always\",\"write_once_then_read\",\"write_and_maintain\"]"`
		Default  string   `json:"default" default:"read_only"`
	} `json:"write_mode"`
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
	ObjectEncoding struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"leb_lew\",\"leb_bew\",\"beb_lew\",\"beb_bew\"]"`
		Default  string   `json:"default" default:"beb_lew"`
	} `json:"object_encoding"`
	IsOutput struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
	} `json:"is_output"`
	Enable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
	} `json:"enable"`
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
