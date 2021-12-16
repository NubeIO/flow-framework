package schedule

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/src/utilstime"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

//CheckWeeklyScheduleEntryWithEntryTimezone checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period ignores entry.Timezone and uses Local timezone.
func CheckWeeklyScheduleEntryWithEntryTimezone(entry WeeklyScheduleEntry) ScheduleCheckerResult {
	return CheckWeeklyScheduleEntry(entry, entry.Timezone)
}

//CheckWeeklyScheduleEntry checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleEntry(entry WeeklyScheduleEntry, checkTimezone string) ScheduleCheckerResult {
	result := ScheduleCheckerResult{}
	result.Payload = entry.Value
	result.IsActive = false
	result.IsException = false

	//get time.Location for entry timezone and check timezone
	location, err := time.LoadLocation(checkTimezone)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
		return result
	}

	now := time.Now().In(location)
	result.CheckTime = now.Unix()
	//get day of week and compare with entry.DaysNums
	//nowHour, nowMinute, nowSecond := now.Clock()
	nowYear, nowMonth, nowDate := now.Date()
	//log.Println("nowYear: ", nowYear, "nowMonth: ", nowMonth, "nowDate: ", nowDate)
	nowDayOfWeek := DaysOfTheWeek(now.Weekday())
	//nowDayOfWeekString := now.String()

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

	//parse start and end time into current day timestamps
	startTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStartHour, entryStartMins, 0, 0, location)
	stopTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStopHour, entryStopMins, 59, 0, location)

	//Check if the schedule is active today
	scheduleActiveToday := false
	for _, day := range entry.DaysNums {
		if day == nowDayOfWeek {
			scheduleActiveToday = true
			break
		}
	}
	//log.Println("scheduleActiveToday: ", scheduleActiveToday)

	//find the next active schedule day
	//log.Println("nowDayOfWeek: ", nowDayOfWeek, "entry.DaysNums: ", entry.DaysNums)
	nextDay, err := getNextScheduleDay(nowDayOfWeek, entry.DaysNums)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
		return result
	}
	nextDayDuration := getDurationTillNextScheduleDay(nowDayOfWeek, nextDay)
	//log.Println("nextDay: ", nextDay, "nextDayDuration: ", nextDayDuration)

	if scheduleActiveToday && now.Before(stopTimestamp) { //scheduled today and hasn't finished yet
		//check if today's schedule is currently active
		if now.After(startTimestamp) {
			result.IsActive = true
		}
		//PeriodStartStop
		result.PeriodStart = startTimestamp.Unix()
		result.PeriodStop = stopTimestamp.Unix()
		//NextStartStop
		result.NextStart = startTimestamp.Add(nextDayDuration).Unix()
		result.NextStop = stopTimestamp.Add(nextDayDuration).Unix()
	} else { // not scheduled today, OR scheduled today, but period already passed
		//PeriodStartStop
		result.PeriodStart = startTimestamp.Add(nextDayDuration).Unix()
		result.PeriodStop = stopTimestamp.Add(nextDayDuration).Unix()
		//NextStartStop
		//next scheduled day is being used for PeriodStart and PeriodStop, so get the NEXT next scheduled day for NextStart and NextStop
		nextNextDay, err := getNextScheduleDay(nextDay, entry.DaysNums)
		if err != nil {
			result.ErrorFlag = true
			result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
			return result
		}
		nextNextDayDuration := getDurationTillNextScheduleDay(nextDay, nextNextDay)
		result.NextStart = startTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
		result.NextStop = stopTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
	}
	//log.Println("CheckWeeklyScheduleEntry RETURN: ", result)
	return result
}

//CheckWeeklyScheduleCollection checks if there is a WeeklyScheduleEntry in the provided WeeklyScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleCollection(scheduleMap TypeWeekly, scheduleName string) ScheduleCheckerResult {
	finalResult := ScheduleCheckerResult{}
	var singleResult ScheduleCheckerResult
	count := 0
	var err error
	for _, scheduleEntry := range scheduleMap {
		if scheduleName == "ANY" || scheduleName == "ALL" || scheduleEntry.Name == scheduleName {
			scheduleEntry = ConvertDaysStringsToInt(scheduleEntry)
			//fmt.Println("WEEKLY SCHEDULE ", i, ": ", scheduleEntry)
			if scheduleEntry.Timezone == "" { // If timezone field is not assigned, get timezone from System Time
				systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
				//fmt.Println("systemTimezone 2: ", systemTimezone)
				if systemTimezone == "" {
					zone, _ := utilstime.GetHardwareTZ()
					scheduleEntry.Timezone = zone
				} else {
					scheduleEntry.Timezone = systemTimezone
				}
			}
			//singleResult = CheckWeeklyScheduleEntry(scheduleEntry, "Australia/Sydney")
			singleResult = CheckWeeklyScheduleEntryWithEntryTimezone(scheduleEntry)
			singleResult.Name = scheduleName
			//fmt.Println("finalResult ", finalResult, "singleResult: ", singleResult)
			if count == 0 {
				finalResult = singleResult
			} else {
				finalResult, err = CombineScheduleCheckerResults(finalResult, singleResult)
				if err != nil {
					log.Errorf("CheckEventScheduleEntry %v\n", err)
				}
			}
			//fmt.Println("finalResult ", finalResult)
			count++
		}
	}
	return finalResult
}

//getNextScheduleDay() returns the next scheduled day.
func getNextScheduleDay(today DaysOfTheWeek, scheduleDays []DaysOfTheWeek) (DaysOfTheWeek, error) {
	if len(scheduleDays) == 0 {
		return 0, errors.New("NO DAYS SCHEDULED")
	}
	i := int(today)
	j := 0
	var nextDay DaysOfTheWeek = 0
	nextFound := false
	//check until the next scheduled day is found in scheduleDays
	for j < 7 {
		//check the next day (rollover at end of week [Saturday(6)]
		if i == 6 {
			i = 0
		} else {
			i++
		}
		//check each of the days in scheduleDays
		for x := 0; x < len(scheduleDays); x++ {
			if scheduleDays[x] == DaysOfTheWeek(i) {
				nextDay = DaysOfTheWeek(i)
				nextFound = true
				break
			}
		}
		if nextFound {
			break
		}
		j++
	}
	if !nextFound {
		return 0, errors.New("NO VALID DAYS IN SCHEDULE")
	} else {
		return nextDay, nil
	}
}

//getDurationTillNextScheduleDay() returns the time.Duration to the next scheduled day.
func getDurationTillNextScheduleDay(today DaysOfTheWeek, nextDay DaysOfTheWeek) time.Duration {
	//make Duration to next scheduled day
	daysTillNext := 0
	if nextDay > today {
		daysTillNext = int(nextDay) - int(today)
	} else {
		//rest of the week, plus the next day as integer
		daysTillNext = (7 - int(today)) + int(nextDay)
	}
	//log.Println("daysTillNext: ", daysTillNext, "strconv.Itoa(daysTillNext): ", strconv.Itoa(daysTillNext))
	tillNextDayDuration, _ := time.ParseDuration(strconv.Itoa(daysTillNext*24) + "h")
	return tillNextDayDuration
}

//ConvertDaysStringsToInt converts strings of weekdays to integers
func ConvertDaysStringsToInt(weeklyScheduleEntry WeeklyScheduleEntry) WeeklyScheduleEntry {
	var lowerCaseStringDay string
	for _, v := range weeklyScheduleEntry.Days {
		lowerCaseStringDay = strings.ToLower(v)
		dayInt := DaysMap[lowerCaseStringDay]
		weeklyScheduleEntry.DaysNums = append(weeklyScheduleEntry.DaysNums, dayInt)
	}
	return weeklyScheduleEntry
}

//GetNextStartStop gets the next start and stop times from a single ScheduleCheckerResult
func GetNextStartStop(weeklyResultObj ScheduleCheckerResult) (nextStart int64, nextStop int64) {
	var start, stop int64
	if weeklyResultObj.IsActive {
		stop = weeklyResultObj.PeriodStop
		start = weeklyResultObj.NextStart
	} else {
		stop = weeklyResultObj.PeriodStop
		start = weeklyResultObj.PeriodStart
	}
	return start, stop
}

//WeeklyCheck checks all Weekly Schedules in the payload for active periods. It returns a combined ScheduleCheckerResult of all Weekly Schedules.
func WeeklyCheck(weekly TypeWeekly, scheduleName string) (ScheduleCheckerResult, error) {
	results := CheckWeeklyScheduleCollection(weekly, scheduleName)
	return results, nil
}
