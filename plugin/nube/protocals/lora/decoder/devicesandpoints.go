package decoder

import (
	"github.com/NubeIO/flow-framework/model"
)

type LoRaDeviceDescription struct {
	DeviceName      string
	Model           string
	SensorCode      string
	CheckLength     func(data string) bool
	Decode          func(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{})
	GetPointsStruct func() interface{}
}

var NilLoRaDeviceDescription LoRaDeviceDescription = LoRaDeviceDescription{
	DeviceName:      "",
	Model:           "",
	SensorCode:      "",
	CheckLength:     NilLoRaDeviceDescriptionCheckLength,
	Decode:          NilLoRaDeviceDescriptionDecode,
	GetPointsStruct: NilLoRaDeviceDescriptionGetPointsStruct,
}

func NilLoRaDeviceDescriptionCheckLength(data string) bool {
	return false
}

func NilLoRaDeviceDescriptionDecode(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{}) {
	return &CommonValues{}, struct{}{}
}

func NilLoRaDeviceDescriptionGetPointsStruct() interface{} {
	return struct{}{}
}

var LoRaDeviceDescriptions = [...]LoRaDeviceDescription{
	{
		DeviceName:      "MicroEdge",
		Model:           "MicroEdge",
		SensorCode:      "AA",
		CheckLength:     CheckPayloadLengthME,
		Decode:          DecodeME,
		GetPointsStruct: GetPointsStructME,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THLM",
		SensorCode:      "AB",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHLM,
		GetPointsStruct: GetPointsStructTHLM,
	},
	{
		DeviceName:      "Droplet",
		Model:           "TH",
		SensorCode:      "B0",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTH,
		GetPointsStruct: GetPointsStructTH,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THL",
		SensorCode:      "B1",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHL,
		GetPointsStruct: GetPointsStructTHL,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THLM",
		SensorCode:      "B2",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHLM,
		GetPointsStruct: GetPointsStructTHLM,
	},
	{
		DeviceName:      "ZipHydroTap",
		Model:           "ZipHydroTap",
		SensorCode:      "D1",
		CheckLength:     CheckPayloadLengthZHT,
		Decode:          DecodeZHT,
		GetPointsStruct: GetPointsStructZHT,
	},
}

func GetLoRaDeviceDescriptionFromID(devID string) *LoRaDeviceDescription {
	return GetLoRaDeviceDescription(devID[2:4])
}

func GetLoRaDeviceDescription(sensorCode string) *LoRaDeviceDescription {
	for _, dev := range LoRaDeviceDescriptions {
		if sensorCode == dev.SensorCode {
			return &dev
		}
	}
	return &NilLoRaDeviceDescription
}

func GetDevicePointsStruct(device *model.Device) interface{} {
	return getDeviceDescriptionFromPayload(device.AddressUUID).GetPointsStruct()
}
