package main

import (
	"bufio"

	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type SerialSetting struct {
	SerialPort     string
	Enable         bool
	BaudRate       int
	StopBits       serial.StopBits
	Parity         serial.Parity
	DataBits       int
	Timeout        int
	ActivePortList []string
	Connected      bool
	Error          bool
	I              Instance
}

// SerialOpen open serial port
func (i *Instance) SerialOpen() error {
	s := new(SerialSetting)
	var arg api.Args
	arg.WithSerialConnection = true
	net, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return err
	}
	if net.SerialConnection == nil {
		return err
	}
	s.SerialPort = net.SerialConnection.SerialPort
	s.BaudRate = int(net.SerialConnection.BaudRate)
	connected := false
	go func() error {
		sc := New(s)
		connected, err = sc.NewSerialConnection()
		if err != nil {
			log.Errorf("lora: issue on SerialOpenAndRead: %v\n", err)
		}
		sc.Loop()
		return nil
	}()
	return nil

}

// SerialClose close serial port
func (i *Instance) SerialClose() error {
	err := Disconnect()
	if err != nil {
		return err
	}
	return nil
}

func New(s *SerialSetting) *SerialSetting {
	if s.SerialPort == "" {
		s.SerialPort = "/dev/ttyACM0"
	}
	if s.BaudRate == 0 {
		s.BaudRate = 38400
	}
	return &SerialSetting{
		SerialPort: s.SerialPort,
		BaudRate:   s.BaudRate,
	}
}

var Port serial.Port

func (s *SerialSetting) NewSerialConnection() (connected bool, err error) {
	portName := s.SerialPort
	baudRate := s.BaudRate
	parity := s.Parity
	stopBits := s.StopBits
	dataBits := s.DataBits
	if s.Connected {
		log.Info("Existing serial port connection by this app is open So! close existing connection")
		err := Disconnect()
		if err != nil {
			log.Info(err)
			s.Error = true
			return false, err
		}
	}
	log.Info("LORA: try and connect to:", portName)
	m := &serial.Mode{
		BaudRate: baudRate,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: stopBits,
	}

	ports, err := serial.GetPortsList()
	s.ActivePortList = ports
	log.Info("LORA: ports: ", ports)

	port, err := serial.Open(portName, m)
	if err != nil {
		s.Error = true
		log.Error("LORA: error on open port", " ", err)
		return false, err
	}
	Port = port
	s.Connected = true
	log.Info("LORA: Connected to serial port: ", " ", portName, " ", "connected: ", " ", s.Connected)
	return s.Connected, nil
}

func (s *SerialSetting) Loop() {
	if s.Error || !s.Connected || Port == nil {
		return
	}
	count := 0
	scanner := bufio.NewScanner(Port)
	for scanner.Scan() {
		var data = scanner.Text()
		if decoder.CheckPayloadLength(data) {
			count = count + 1
			commonData, fullData := decoder.DecodePayload(data)
			s.I.publishSensor(commonData, fullData)
		} else {
			log.Printf("LORA: serial messsage size %d", len(data))
		}
	}
}
func Disconnect() error {
	if Port != nil {
		err := Port.Close()
		if err != nil {
			log.Error("LORA: err on trying to close the port")
			return err
		}
	}
	return nil
}
