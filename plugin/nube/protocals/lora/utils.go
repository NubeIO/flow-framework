package main

import (
	"reflect"
	"strings"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"go.bug.st/serial"
)

//wizard make a network/dev/pnt
func (inst *Instance) wizardSerial(body wizard) (string, error) {
	sp := "/dev/ttyACM0"
	if body.SerialPort != "" {
		sp = body.SerialPort
	}
	id := "AAB296C4"
	if body.SensorID != "" {
		id = body.SensorID
	}
	st := "THLM"
	if body.SensorType != "" {
		st = body.SensorType
	}
	var net model.Network
	net.Name = model.TransProtocol.Lora
	net.TransportType = model.TransType.Serial
	net.PluginPath = model.TransProtocol.Lora
	net.SerialPort = utils.NewStringAddress(sp)
	net.SerialBaudRate = utils.NewUint(38400)

	var dev model.Device
	dev.Name = model.TransProtocol.Lora
	dev.AddressUUID = &id
	dev.Manufacture = model.CommonNaming.NubeIO
	dev.Model = st

	_, err = inst.db.WizardNewNetworkDevicePoint("lora", &net, &dev, nil)
	if err != nil {
		return "error: add lora serial network wizard", err
	}

	inst.Disable()
	inst.Enable()
	return "pass: added network and points", err
}

//listSerialPorts list all serial ports on host
func (inst *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}

func BoolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// TODO: move this to a more global project utils file
func getStructFieldJSONNameByIndex(thing interface{}, index int) string {
	field := reflect.TypeOf(thing).Field(index)
	return getReflectFieldJSONName(field)
}

// TODO: move this to a more global project utils file
func getStructFieldJSONNameByName(thing interface{}, name string) string {
	field, err := reflect.TypeOf(thing).FieldByName(name)
	if !err {
		panic(err)
	}
	return getReflectFieldJSONName(field)
}

// TODO: move this to a more global project utils file
func getReflectFieldJSONName(field reflect.StructField) string {
	fieldName := field.Name

	switch jsonTag := field.Tag.Get("json"); jsonTag {
	case "-":
		fallthrough
	case "":
		return fieldName
	default:
		parts := strings.Split(jsonTag, ",")
		name := parts[0]
		if name == "" {
			name = fieldName
		}
		return name
	}
}
