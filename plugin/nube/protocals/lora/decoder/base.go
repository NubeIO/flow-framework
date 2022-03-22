// Package decoder contains all device specific payload decoding and checks
package decoder

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CommonValues struct {
	Sensor string `json:"sensor"`
	ID     string `json:"id"`
	Rssi   int    `json:"rssi"`
}

func getDeviceDescriptionFromPayload(data *string) *LoRaDeviceDescription {
	if data == nil {
		return nil
	}
	if *data == "!" {
		return nil
	}
	if len(*data) >= 4 {
		sensorCode := (*data)[2:4]
		return GetLoRaDeviceDescription(sensorCode)
	} else {
		return nil
	}
}

func checkPayloadLength(data string, dev *LoRaDeviceDescription) bool {
	log.Println("lora-decoder:", data)
	dl := len(data)
	if dl <= 4 {
		return false
	}
	if data == "!\r\n" || data == "!\n" {
		return false
	}
	if dev == &NilLoRaDeviceDescription {
		return false
	}
	return dev.CheckLength(data)
}

func DecodePayload(data string) (*CommonValues, interface{}) {
	devDesc := getDeviceDescriptionFromPayload(&data)
	if devDesc == &NilLoRaDeviceDescription {
		return &CommonValues{}, nil
	}
	if !checkPayloadLength(data, devDesc) {
		return &CommonValues{}, nil
	}

	cmn, payload := devDesc.Decode(data, devDesc)
	decodeCommonValues(cmn, data, devDesc.Model)
	return cmn, payload
}

func decodeCommonValues(payload *CommonValues, data string, sensor string) {
	payload.Sensor = sensor
	payload.ID = decodeID(data)
	payload.Rssi = decodeRSSI(data)
}

func decodeID(data string) string {
	id := data[0:8]
	log.Println("lora-decoder:", id)
	return id
}

func decodeRSSI(data string) int {
	dataLen := len(data)
	v, _ := strconv.ParseInt(data[dataLen-4:dataLen-2], 16, 0)
	v = v * -1
	return int(v)
}
