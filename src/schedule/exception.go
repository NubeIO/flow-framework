package schedule

import (
	"errors"
)

//TODO: Known Issue 1) If exception event has passed, it will not be included in the exception ScheduleCheckerResults, and therefore if a schedule period start was before the exception period, the final result will not show the correct PeriodStart

//ExceptionCheck checks all Exception Schedules in the payload for active periods. It returns a combined ScheduleCheckerResult of all Exception Schedules.
func ExceptionCheck(exceptions TypeEvents, scheduleName string) (ScheduleCheckerResult, error) {
	//treat Exception schedules as Event schedules until the final step.
	results := CheckEventScheduleCollection(exceptions, scheduleName)

	//once the ScheduleCheckerResult has been computed (as an Event type schedule) then set the IsException flag.  This result should later be combined with Event and Weekly schedules.
	results.IsException = true
	return results, nil
}

//ApplyExceptionSchedule Combines ScheduleCheckerResults with an Exception ScheduleCheckerResults.  The IsException flag must be set for the Exception ScheduleCheckerResult.
//This function will mask the Exception schedule period and return the updated ScheduleCheckerResult.
func ApplyExceptionSchedule(current ScheduleCheckerResult, exception ScheduleCheckerResult) (ScheduleCheckerResult, error) {
	//log.Println("ApplyExceptionSchedule()")
	result := ScheduleCheckerResult{}
	var err error

	//Check for Exception Schedules
	if current.IsException || !exception.IsException {
		err = errors.New("ApplyExceptionSchedule function must have one ScheduleCheckerResults with IsException flag false, and one with IsException flag true")
		return result, err
	}

	//Check for empty ScheduleCheckerResult periods
	if current.PeriodStart == 0 || current.PeriodStop == 0 {
		return current, err
	}
	if exception.PeriodStart == 0 || exception.PeriodStop == 0 {
		return current, err
	}

	//AlertFlag & ErrorFlag & ErrorStrings
	if current.ErrorFlag {
		err = errors.New("`current` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return current, err
	}
	if exception.ErrorFlag {
		err = errors.New("`exception` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return current, err
	}
	result.AlertFlag = current.AlertFlag || exception.AlertFlag
	result.ErrorStrings = append(current.ErrorStrings, exception.ErrorStrings...)

	//Name
	if current.Name != exception.Name {
		result.Name = "ANY"
	} else {
		result.Name = current.Name
	}

	//Find order of periods
	currentPeriod := 0
	exceptionPeriod := 0
	//Find which Period is first
	if current.PeriodStart <= exception.PeriodStart {
		currentPeriod = 1
		exceptionPeriod = 2
	} else {
		exceptionPeriod = 1
		currentPeriod = 2
	}

	//Check if the periods overlap.
	overlap := false
	if currentPeriod == 1 && exception.PeriodStart < current.PeriodStop {
		overlap = true
	} else if exceptionPeriod == 1 && current.PeriodStart < exception.PeriodStop {
		overlap = true
	}
	//If no overlap, then current schedule isn't modified by exception schedule
	if !overlap {
		return current, err
	}

	//PeriodStart and //PeriodStop
	//If the current period starts first, and the periods overlap, then the result start time is the current period start, and the stop time is the exception start time.
	if currentPeriod == 1 {
		result.PeriodStart = current.PeriodStart
		result.PeriodStop = exception.PeriodStart
	} else { //exception period is first
		if exception.PeriodStop < current.PeriodStop {
			result.PeriodStart = exception.PeriodStop
			result.PeriodStop = current.PeriodStop
		} else { //current period is entirely within the exception period, so there is no schedule period.
			result.PeriodStart = 0
			result.PeriodStop = 0
		}
	}

	//Payload
	result.Payload = current.Payload

	//CheckTime
	if current.CheckTime < exception.CheckTime {
		result.CheckTime = current.CheckTime
	} else {
		result.CheckTime = exception.CheckTime
	}

	//Check if period has already passed
	if result.PeriodStop <= result.CheckTime {
		result.PeriodStart = 0
		result.PeriodStop = 0
		result.Payload = 0
	}

	//IsActive
	if result.CheckTime >= result.PeriodStart && result.CheckTime < result.PeriodStop {
		result.IsActive = true
	} else {
		result.IsActive = false
	}

	return result, nil
}
