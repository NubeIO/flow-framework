package main

import (
	"fmt"
	"github.com/grid-x/modbus"
	"time"
)

func main() {

	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 38400
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveID = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	defer handler.Close()
	fmt.Println(err)
	client := modbus.NewClient(handler)
	results, err := client.WriteSingleCoil(0, 0xFF00)
	fmt.Println(err)
	fmt.Println(results)

}
