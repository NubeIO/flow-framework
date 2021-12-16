package decoder

import (
	"strconv"
)

type TDropletTH struct {
	CommonValues
	Voltage     float64 `json:"voltage"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Humidity    int     `json:"humidity"`
}

type TDropletTHL struct {
	TDropletTH
	Light int `json:"light"`
}

type TDropletTHLM struct {
	TDropletTHL
	Motion bool `json:"motion"`
}

func GetPointsStructTH() interface{} {
	return TDropletTH{}
}

func GetPointsStructTHL() interface{} {
	return TDropletTHL{}
}

func GetPointsStructTHLM() interface{} {
	return TDropletTHLM{}
}

func CheckPayloadLengthDroplet(data string) bool {
	dl := len(data)
	return dl == 36 || dl == 32 || dl == 44
}

func DecodeDropletTH(data string, _ *LoRaDeviceDescription) (*CommonValues, interface{}) {
	temperature := dropletTemp(data)
	pressure := dropletPressure(data)
	humidity := dropletHumidity(data)
	voltage := dropletVoltage(data)
	v := TDropletTH{
		Voltage:     voltage,
		Temperature: temperature,
		Pressure:    pressure,
		Humidity:    humidity,
	}
	return &v.CommonValues, v
}

func DecodeDropletTHL(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{}) {
	_, d := DecodeDropletTH(data, devDesc)
	light := dropletLight(data)
	v := TDropletTHL{
		TDropletTH: d.(TDropletTH),
		Light:      light,
	}
	return &v.CommonValues, v
}

func DecodeDropletTHLM(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{}) {
	_, d := DecodeDropletTHL(data, devDesc)
	motion := dropletMotion(data)
	v := TDropletTHLM{
		TDropletTHL: d.(TDropletTHL),
		Motion:      motion,
	}
	return &v.CommonValues, v
}

func dropletTemp(data string) float64 {
	v, _ := strconv.ParseInt(data[10:12]+data[8:10], 16, 0)
	v_ := float64(v) / 100
	return v_
}

func dropletPressure(data string) float64 {
	v, _ := strconv.ParseInt(data[14:16]+data[12:14], 16, 0)
	v_ := float64(v) / 10
	return v_
}

func dropletHumidity(data string) int {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v = v & 127
	return int(v)
}

func dropletVoltage(data string) float64 {
	v, _ := strconv.ParseInt(data[22:24], 16, 0)
	v_ := float64(v) / 50
	return v_
}

func dropletLight(data string) int {
	v, _ := strconv.ParseInt(data[20:22]+data[18:20], 16, 0)
	return int(v)
}

func dropletMotion(data string) bool {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	return v > 127
}
