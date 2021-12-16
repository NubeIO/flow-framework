package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/src/utilstime"
)

func main() {

	tz, err := utilstime.GetHardwareTZ()
	if err != nil {
		return
	}
	fmt.Println(tz, err)

	//schedule.ScheduleTest()
	schedule.ModbusScheduleTest()
}
