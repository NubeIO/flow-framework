package main

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
	"github.com/grid-x/modbus"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	Host       string        `json:"ip"`
	Port       string        `json:"port"`
	SerialPort string        `json:"serial_port"`
	BaudRate   uint          `json:"baud_rate"` //38400
	Parity     string        `json:"parity"`    //none, odd, even DEFAULT IS PARITY_NONE
	DataBits   uint          `json:"data_bits"` //7 or 8
	StopBits   uint          `json:"stop_bits"` //1 or 2
	Timeout    time.Duration `json:"device_timeout_in_ms"`
}

func (inst *Instance) setClient(network *model.Network, device *model.Device, cacheClient bool) (mbClient smod.ModbusClient, err error) {
	if network.TransportType == model.TransType.Serial || network.TransportType == model.TransType.LoRa {
		serialPort := "/dev/ttyUSB0"
		baudRate := 38400
		stopBits := 1
		dataBits := 8
		parity := "N"
		if network.SerialPort != nil {
			serialPort = nils.StringIsNil(network.SerialPort)
		}
		if network.SerialBaudRate != nil {
			baudRate = int(nils.UnitIsNil(network.SerialBaudRate))
		}
		if network.SerialDataBits != nil {
			dataBits = int(nils.UnitIsNil(network.SerialDataBits))
		}
		if network.SerialStopBits != nil {
			stopBits = int(nils.UnitIsNil(network.SerialStopBits))
		}
		if network.SerialParity != nil {
			parity = nils.StringIsNil(network.SerialParity)
		}
		handler := modbus.NewRTUClientHandler(serialPort)
		handler.BaudRate = baudRate
		handler.DataBits = dataBits
		handler.Parity = setParity(parity)
		handler.StopBits = stopBits
		handler.Timeout = 5 * time.Second
		err := handler.Connect()
		if err != nil {
			return smod.ModbusClient{}, err
		}
		defer handler.Close()
		mc := modbus.NewClient(handler)
		mbClient.RTUClientHandler = handler
		mbClient.Client = mc
		return mbClient, nil

	} else {
		url, err := uurl.JoinIpPort(device.Host, device.Port)
		if err != nil {
			log.Errorf("modbus: failed to validate device IP %s\n", url)
			return smod.ModbusClient{}, err
		}
		handler := modbus.NewTCPClientHandler(url)
		err = handler.Connect()
		if err != nil {
			return smod.ModbusClient{}, err
		}
		defer handler.Close()
		mc := modbus.NewClient(handler)
		mbClient.TCPClientHandler = handler
		mbClient.Client = mc
		return mbClient, nil
	}
}

func setParity(in string) string {
	if in == model.SerialParity.None {
		return "N"
	} else if in == model.SerialParity.Odd {
		return "O"
	} else if in == model.SerialParity.Even {
		return "E"
	} else {
		return "N"
	}
}
