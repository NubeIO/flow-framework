package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/utils"
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

var connected bool

func (i *Instance) setClient(network *model.Network, cacheClient bool) (mbClient smod.ModbusClient, err error) {

	if network.TransportType == model.TransType.Serial {
		if network.SerialPort == nil {
			log.Errorln("invalid serial connection details", "SerialPort")

		}
		if network.SerialBaudRate == nil {
			log.Errorln("invalid serial connection details", "SerialBaudRate")

		}
		if network.SerialDataBits == nil {
			log.Errorln("invalid serial connection details", "SerialDataBits")

		}
		if network.SerialStopBits == nil {
			log.Errorln("invalid serial connection details", "SerialStopBits")

		}
		if network.SerialParity == nil {
			log.Errorln("invalid serial connection details", "SerialParity")
		}
		handler := modbus.NewRTUClientHandler(*network.SerialPort)
		handler.BaudRate = int(*network.SerialBaudRate)
		handler.DataBits = int(*network.SerialDataBits)
		handler.Parity = setParity(*network.SerialParity)
		handler.StopBits = int(*network.SerialStopBits)
		handler.Timeout = 5 * time.Second

		handler.Connect()
		defer handler.Close()
		mc := modbus.NewClient(handler)

		mbClient.RTUClientHandler = handler
		mbClient.Client = mc
		connected = true
		return mbClient, nil

	} else {

		if network.Host == nil {
			log.Errorln("invalid ip connection details", "Host")

		}
		if network.Port == nil {
			log.Errorln("invalid serial connection details", "Port")

		}

		url, err := utils.JoinIPPort(cli)

		handler := modbus.NewTCPClientHandler("localhost:502")
		handler.Connect()
		defer handler.Close()
		mc := modbus.NewClient(handler)

		mbClient.TCPClientHandler = handler
		mbClient.Client = mc
		connected = true
		return mbClient, nil
	}

}

//func (i *Instance) setClient(client Client, networkUUID string, cacheClient, isSerial bool) error {
//	var c *modbus.ModbusClient
//	if isSerial {
//		parity := setParity(client.Parity)
//		serialPort := setSerial(client.SerialPort)
//		if client.Timeout < 10 {
//			client.Timeout = 500
//		}
//		//TODO add in a check if client with same details exists
//		c, err = modbus.NewClient(&modbus.ClientConfiguration{
//			URL:      serialPort,
//			Speed:    client.BaudRate,
//			DataBits: client.DataBits,
//			Parity:   parity,
//			StopBits: client.StopBits,
//			Timeout:  client.Timeout * time.Millisecond,
//		})
//	} else {
//		var cli utils.URLParts
//		cli.Transport = "tcp"
//		cli.Host = client.Host
//		cli.Port = client.Port
//		url, err := utils.JoinURL(cli)
//		if err != nil {
//			connected = false
//			return err
//		}
//		if client.Timeout < 10 {
//			client.Timeout = 500
//		}
//		//TODO add in a check if client with same details exists
//		c, err = modbus.NewClient(&modbus.ClientConfiguration{
//			URL:     url,
//			Timeout: client.Timeout * time.Millisecond,
//		})
//		if err != nil {
//			connected = false
//			return err
//		}
//	}
//	var getC interface{}
//	if cacheClient { //store modbus client in cache to reuse the instance
//		getC, _ = i.store.Get(networkUUID)
//		if getC == nil {
//			i.store.Set(networkUUID, c, -1)
//		} else {
//			c = getC.(*modbus.ModbusClient)
//		}
//	}
//	err = c.Open()
//	connected = true
//	restMB = c
//	if err != nil {
//		connected = false
//		return err
//	}
//	return nil
//}

func isConnected() bool {
	return connected
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

func setSerial(port string) string {
	p := fmt.Sprintf("rtu:///%s", port)
	return p
}
