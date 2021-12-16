package schedule

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func ScheduleTest() {
	json, err := ioutil.ReadFile("/home/user/Documents/Nube/Flow_Framework/flow-framework/src/schedule/old/schTest4.json")
	if err != nil {
		log.Errorf("ReadFile %v\n", err)
	}
	decodeSchedule, err := DecodeSchedule(json)
	log.Println("decodeSchedule: ", decodeSchedule)

	scheduleNameToCheck := "HVAC" //TODO: we need a way to specify the schedule name that is being checked for.

	// CHECK WEEKLY SCHEDULES
	//result, err := schedule.WeeklyCheck(decodeSchedule.Weekly, "ANY")  //This will check for any active schedules with any name
	weeklyResult, err := WeeklyCheck(decodeSchedule.Weekly, scheduleNameToCheck) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
	}
	fmt.Println("weeklyResult")
	fmt.Printf("%+v\n", weeklyResult)

	// CHECK EVENT SCHEDULES
	//eventResult, err := schedule.EventCheck(decodeSchedule.Events, "ANY")  //This will check for any active schedules with any name
	eventResult, err := EventCheck(decodeSchedule.Events, scheduleNameToCheck) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
	}
	fmt.Println("eventResult")
	fmt.Printf("%+v\n", eventResult)

	//Combine Event and Weekly schedule results.
	weeklyAndEventResult, err := CombineScheduleCheckerResults(weeklyResult, eventResult)
	fmt.Println("weeklyAndEventResult")
	fmt.Printf("%+v\n", weeklyAndEventResult)

	// CHECK EXCEPTION SCHEDULES
	//exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, "ANY")  //This will check for any active schedules with any name
	exceptionResult, err := ExceptionCheck(decodeSchedule.Exceptions, scheduleNameToCheck) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
	}
	fmt.Println("exceptionResult")
	fmt.Printf("%+v\n", exceptionResult)
	if exceptionResult.CheckIfEmpty() {
		fmt.Println("Exception schedule is empty")
	}

	finalResult, err := ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult) //This applies the exception schedule to mask the combined weekly and event schedules.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on ApplyExceptionSchedule %v\n", err)
	}
	fmt.Println("finalResult")
	fmt.Printf("%+v\n", finalResult)

	fmt.Println("schedule is ", finalResult.IsActive)
}
