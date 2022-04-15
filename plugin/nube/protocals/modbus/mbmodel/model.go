package mbmodel

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
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
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"100"`
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
		Options []string `json:"options" default:"[\"read_coil\",\"read_discrete_input\",\"write_coil\",\"read_register\",\"read_holding\",\"write_holding\"]"`
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
		Options  []string `json:"options" default:"[\"uint16\",\"int16\",\"uint32\",\"int32\",\"uint64\",\"int64\",\"float32\",\"float64\"]"`
		Default  string   `json:"default" default:"uint16"`
	} `json:"data_type"`
	ObjectEncoding struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"leb_lew\",\"leb_bew\",\"beb_lew\",\"beb_bew\"]"`
		Default  string   `json:"default" default:"beb_lew"`
	} `json:"object_encoding"`
	MathOnPresentValue struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"(x + 0) + 0"`
		DisplayName string `json:"display_name" default:"math expression on present value"`
	} `json:"math_on_present_value"`
	MathOnWriteValue struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"(x + 0) + 0"`
		DisplayName string `json:"display_name" default:"math expression on write value"`
	} `json:"math_on_write_value"`
	PointPriorityArrayMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"priority_array_to_present_value\",\"priority_array_to_write_value\",\"read_only_no_priority_array_required\"]"`
		Default  string   `json:"default" default:"priority_array_to_present_value"`
	} `json:"point_priority_use_type"`
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
