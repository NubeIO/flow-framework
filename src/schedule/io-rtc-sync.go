package schedule

import (
	"fmt"
	"time"
)

type IOModuleModbusRTC struct {
	ModbusNetwork      string
	ModbusDevice       string
	Timezone           string
	Unix               int64
	Offset             int
	ModbusRegister     int64
	ModbusRegisterType string
	ModbusDataType     string
}

func GetRTCValues(Timezone string) (Unix int64, Offset int, err error) {
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
