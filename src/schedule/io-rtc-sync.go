package schedule

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"time"
)

type IOModuleModbusRTCData struct {
	ModbusNetwork        string
	ModbusDevice         string
	Timezone             string
	Unix                 int64
	Offset               int
	ModbusRTCRegister    int64
	ModbusOffsetRegister int64
	ModbusRegisterType   string
	ModbusDataType       string
}

func GetRTCValuesByTimezone(Timezone string) (Unix int64, Offset int, err error) {
	//get time.Location for entry timezone and check timezone
	location, err := time.LoadLocation(Timezone)
	if err != nil {
		return 0, 0, fmt.Errorf("IO MODULE RTC SYNC: Invalid Timezone.  Error: %s", err)
	}
	now := time.Now().In(location)
	unix := now.Unix()
	_, offset := now.Zone()

	return unix, offset, nil
}

func AssignRTCValuesToModbusDevice(point *model.Point, Unix int64, Offset int, Timezone string) (IOModuleModbusRTCData, error) {
	result := IOModuleModbusRTCData{}
	result.Timezone = Timezone
	result.Unix = Unix
	result.Offset = Offset
	result.ModbusRTCRegister = 307
	result.ModbusOffsetRegister = 309
	result.ModbusRegisterType = "WRITE_REGISTERS"
	result.ModbusDataType = "INT32"

	//get IO Module data from API to fill in ModbusNetwork and ModbusDevice

	return result, nil
}

func WriteRTCModbusValuesToIOModule(RTCSettings IOModuleModbusRTCData) error {
	//Write values to IO Module via Modbus Service API

	return nil
}
