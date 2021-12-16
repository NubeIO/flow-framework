package schedule

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/src/utilstime"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
)

type ScheduleCheckerModbusResult struct {
	Name                      string                    `json:"name"`
	Periods                   [7]WeeklyPeriodCollection //one WeeklyPeriodCollection for each day
	EventScheduleIsActive     bool                      `json:"event_schedule_is_active"`
	ExceptionScheduleIsActive bool                      `json:"exception_schedule_is_active"`
	CheckTime                 int64                     `json:"check_time"` //unix timestamp in seconds
	ErrorFlag                 bool                      `json:"error_flag"`
	AlertFlag                 bool                      `json:"alert_flag"`
	ErrorStrings              []string                  `json:"error_strings"`
}

type WeeklyPeriodCollection [5]WeeklyPeriod //Currently modbus schedules support 5 periods per day

type WeeklyPeriod struct {
	Start int  `json:"start"`
	Stop  int  `json:"stop"`
	Set   bool `json:"set"`
}

func ModbusScheduleTest() {
	json, err := ioutil.ReadFile("/home/user/Documents/Nube/Flow_Framework/flow-framework/src/schedule/old/schTest4.json")
	if err != nil {
		log.Errorf("ReadFile %v\n", err)
	}
	decodeSchedule, err := DecodeSchedule(json)
	log.Println("decodeSchedule: ", decodeSchedule)

	scheduleNameToCheck := "HVAC" //TODO: we need a way to specify the schedule name that is being checked for.

	modbusScheduleResult := ConvertScheduleJsonToModbusSchedule(decodeSchedule, scheduleNameToCheck, true)

	fmt.Println("modbusScheduleResult")
	fmt.Printf("%+v\n", modbusScheduleResult)

}

func ConvertScheduleJsonToModbusSchedule(scheduleJSON SchTypes, scheduleName string, includeEventSchedules bool) ScheduleCheckerModbusResult {
	// CONVERT WEEKLY SCHEDULES
	result := GetWeeklyPeriodsFromWeeklyScheduleCollection(scheduleJSON.Weekly, scheduleName)
	/*
		if err != nil {
			log.Errorf("system-plugin-modbus-schedule: issue on GetWeeklyPeriodsFromWeeklyScheduleCollection %v\n", err)
		}
	*/
	fmt.Println("weeklyResult")
	fmt.Printf("%+v\n", result)
	if !includeEventSchedules {
		return result
	}

	// CHECK EVENT SCHEDULES
	eventResult, err := EventCheck(scheduleJSON.Events, scheduleName) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
	}
	fmt.Println("eventResult")
	fmt.Printf("%+v\n", eventResult)

	result.EventScheduleIsActive = eventResult.IsActive

	// CHECK EXCEPTION SCHEDULES
	//exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, "ANY")  //This will check for any active schedules with any name
	exceptionResult, err := ExceptionCheck(scheduleJSON.Exceptions, scheduleName) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
	}
	fmt.Println("exceptionResult")
	fmt.Printf("%+v\n", exceptionResult)
	if exceptionResult.CheckIfEmpty() {
		fmt.Println("Exception schedule is empty")
	} else {
		if exceptionResult.IsActive {
			result.ExceptionScheduleIsActive = true
		}
	}

	return result
}

//GetWeeklyPeriodsFromWeeklySchedule Converts a single WeeklyScheduleEntry into the scheduled periods to values to be sent via modbus.
func GetWeeklyPeriodsFromWeeklySchedule(entry WeeklyScheduleEntry) ScheduleCheckerModbusResult {
	result := ScheduleCheckerModbusResult{}

	//Get current time in schedule timezone or system timezone
	if entry.Timezone == "" { // If timezone field is not assigned, get timezone from System Time
		systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
		//fmt.Println("systemTimezone 2: ", systemTimezone)
		if systemTimezone == "" {
			zone, _ := utilstime.GetHardwareTZ()
			entry.Timezone = zone
		} else {
			entry.Timezone = systemTimezone
		}
	}
	location, _ := time.LoadLocation(entry.Timezone)
	now := time.Now().In(location)
	result.CheckTime = now.Unix()

	//parse start and stop times
	var entryStartHour, entryStartMins, entryStopHour, entryStopMins int
	n, err1 := fmt.Sscanf(entry.Start, "%d:%d", &entryStartHour, &entryStartMins)
	m, err2 := fmt.Sscanf(entry.End, "%d:%d", &entryStopHour, &entryStopMins)
	if n != 2 || m != 2 || err1 != nil || err2 != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/Stop Time")
		return result
	}
	//log.Println("entryStartHour: ", entryStartHour, "entryStartMins: ", entryStartMins, "entryStopHour: ", entryStopHour, "entryStopMins: ", entryStopMins)
	for _, v := range entry.DaysNums {
		startTime := (entryStartHour * 100) + entryStartMins
		stopTime := (entryStopHour * 100) + entryStopMins
		result.Periods[v][0] = WeeklyPeriod{Start: startTime, Stop: stopTime, Set: true}
	}
	return result
}

//CombineScheduleCheckerWeeklyModbusResult Combines ScheduleCheckerModbusResults and checks that there aren't too many periods per day.
//Only deals with the weekly properties; event properties are handled in another function.
//The 'new' ScheduleCheckerModbusResults argument must only have 1 scheduled period per day (not previously combined).
func CombineScheduleCheckerWeeklyModbusResult(current, new ScheduleCheckerModbusResult) (ScheduleCheckerModbusResult, error) {
	//log.Println("CombineScheduleCheckerModbusResult()")
	result := ScheduleCheckerModbusResult{}
	var err error

	//AlertFlag & ErrorFlag & ErrorStrings
	if new.ErrorFlag {
		err = errors.New("`new` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return current, err
	}
	if current.ErrorFlag {
		err = errors.New("`current` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return new, err
	}
	result.AlertFlag = current.AlertFlag || new.AlertFlag
	result.ErrorStrings = append(current.ErrorStrings, new.ErrorStrings...)

	//Name
	if current.Name != new.Name {
		result.Name = "ANY"
	} else {
		result.Name = current.Name
	}

	//CheckTime
	if current.CheckTime < new.CheckTime {
		result.CheckTime = current.CheckTime
	} else {
		result.CheckTime = new.CheckTime
	}

	//Periods
	result.Periods = current.Periods
	for i, day := range result.Periods {
		for j, period := range day {
			if period.Set == false && new.Periods[i][0].Set == true {
				result.Periods[i][j] = WeeklyPeriod{Start: new.Periods[i][0].Start, Stop: new.Periods[i][0].Stop, Set: true}
				break
			} else if j == 4 && new.Periods[i][0].Set == true {
				result.AlertFlag = true
				weekday := time.Weekday(i).String()
				result.ErrorStrings = append(result.ErrorStrings, `Warning: 5 scheduled periods already exist for `+weekday+`, cannot add more.`)
			}
		}
	}
	return result, nil
}

//GetWeeklyPeriodsFromWeeklyScheduleCollection checks if there is a WeeklyScheduleEntry in the provided WeeklyScheduleCollection that matches the specified schedule Name.  If so it will convert the scheduled periods to values to be sent via modbus.
func GetWeeklyPeriodsFromWeeklyScheduleCollection(scheduleMap TypeWeekly, scheduleName string) ScheduleCheckerModbusResult {
	finalResult := ScheduleCheckerModbusResult{}
	var singleResult ScheduleCheckerModbusResult
	count := 0
	var err error
	for _, scheduleEntry := range scheduleMap {
		if scheduleName == "ANY" || scheduleName == "ALL" || scheduleEntry.Name == scheduleName {
			/*
				if count >= 5 {
					finalResult.AlertFlag = true
					finalResult.ErrorStrings = append(finalResult.ErrorStrings, "Warning: Maximum (5) daily schedules exceeded.  Output limited to 5 daily schedules")
					break
				}
			*/
			scheduleEntry = ConvertDaysStringsToInt(scheduleEntry)
			//fmt.Println("WEEKLY SCHEDULE ", i, ": ", scheduleEntry)

			singleResult = GetWeeklyPeriodsFromWeeklySchedule(scheduleEntry)
			singleResult.Name = scheduleName

			//fmt.Println("finalResult ", finalResult, "singleResult: ", singleResult)
			if count == 0 {
				finalResult = singleResult
			} else {
				finalResult, err = CombineScheduleCheckerWeeklyModbusResult(finalResult, singleResult)
				if err != nil {
					log.Errorf("GetWeeklyPeriodsFromWeeklyScheduleCollection:  %v\n", err)
				}
			}
			//fmt.Println("finalResult ", finalResult)
			count++
		}
	}
	return finalResult
}
