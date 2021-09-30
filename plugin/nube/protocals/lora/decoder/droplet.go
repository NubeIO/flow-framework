package decoder

import (
	"strconv"
)

type TDropletTH struct {
	CommonValues
	Voltage     int     `json:"voltage"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
}

type TDropletTHL struct {
	TDropletTH
	Light int `json:"light"`
}

type TDropletTHLM struct {
	TDropletTHL
	Motion int `json:"motion"`
}

func DropletTH(data string, sensor TSensorType) TDropletTH {
	d := Common(data, sensor)
	temperature := dropletTemp(data)
	humidity := dropletHumidity(data)
	voltage := dropletVoltage(data)
	v := TDropletTH{
		CommonValues: d,
		Voltage:      voltage,
		Temperature:  temperature,
		Humidity:     humidity,
	}
	return v
}

func DropletTHL(data string, sensor TSensorType) TDropletTHL {
	d := DropletTH(data, sensor)
	light := dropletLight(data)
	v := TDropletTHL{
		TDropletTH: d,
		Light:      light,
	}
	return v
}

func DropletTHLM(data string, sensor TSensorType) TDropletTHLM {
	d := DropletTHL(data, sensor)
	motion := dropletMotion(data)
	v := TDropletTHLM{
		TDropletTHL: d,
		Motion:      motion,
	}
	return v
}

// works with old and new
func dropletTemp(data string) float64 {
	v, _ := strconv.ParseInt(data[10:12]+data[8:10], 16, 0)
	v_ := float64(v) / 100
	return v_
}

//dropletHumidity new
func dropletHumidity(data string) int {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v_ := v & 127
	return int(v_)
}

//dropletHumidityOld older version
func dropletHumidityOld(data string) int {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v_ := v % 128
	return int(v_)
}

// works with old and new
func dropletVoltage(data string) int {
	v, _ := strconv.ParseInt(data[22:24], 16, 0)
	v_ := v / 50
	return int(v_)
}

// works with old and new
func dropletLight(data string) int {
	v := data[20:22] + data[18:20]
	v_, _ := strconv.ParseInt(v, 16, 0)
	return int(v_)
}

// works with old and new
func dropletMotion(data string) int {
	v_, _ := strconv.ParseInt(data[16:18], 16, 0)
	if v_ > 127 {
		return 1
	} else {
		return 0
	}
}

//20ABBC903E083B27457200ED000000574D00  20ABBC90
//11AA0203000000000D2E2D95000000E55900  11AA0203
//19ABAA516D084127417400EB000000944400  19ABAA51

//msg.nodeId = hexString.substring(0, 8);
//msg.temp = parseInt("0x" + hexString.substring(10, 12) + hexString.substring(8, 10)) / 100;
//msg.pressure = parseInt("0x" + hexString.substring(14, 16) + hexString.substring(12, 14)) / 10;
//msg.humidity = parseInt("0x" + hexString.substring(16, 18)) % 128;
//msg.movement = parseInt("0x" + hexString.substring(16, 18)) > 127 ? true : false;
//msg.light = parseInt("0x" + hexString.substring(20, 22) + hexString.substring(18, 20));
//msg.voltage = parseInt("0x" + hexString.substring(22, 24)) / 50;
//msg.checksum = "0x" + hexString.substring(30, 32);
//msg.rssi = parseInt("0x" + hexString.substring(32, 34)) * -1;
//msg.snr = parseInt("0x" + hexString.substring(34, 36)) / 10;
