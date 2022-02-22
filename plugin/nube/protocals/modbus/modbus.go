package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/grid-x/modbus"
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

func (i *Instance) setClient(network *model.Network, cacheClient, isSerial bool) (mbClient smod.ModbusClient, err error) {

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 38400
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveID = 1
	handler.Timeout = 5 * time.Second

	//client.SerialPort = *network.SerialPort
	//client.BaudRate = *network.SerialBaudRate
	//client.DataBits = *network.SerialDataBits
	//client.StopBits = *network.SerialStopBits
	//client.Parity = *network.SerialParity

	handler.Connect()
	defer handler.Close()

	mc := modbus.NewClient(handler)
	var c smod.ModbusClient
	c.RTUClientHandler = handler
	c.Client = mc
	connected = true
	return c, nil
	//c.RegType = smod.HoldingRegister
	//c.Endianness = smod.BigEndian
	//c.WordOrder = smod.LowWordFirst

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

func setParity(in string) uint {
	if in == model.SerialParity.None {
		return smod.ParityNone
	} else if in == model.SerialParity.Odd {
		return smod.ParityOdd
	} else if in == model.SerialParity.Even {
		return smod.ParityEven
	} else {
		return smod.ParityNone
	}
}

func setSerial(port string) string {
	p := fmt.Sprintf("rtu:///%s", port)
	return p
}
